package onlinedp

import (
	"fmt"
	"log"
	"simplex/db"
	"simplex/dp"
	"database/sql"
)

type LnrFeat struct{ FID int }

//Find and merge simple segments
func (self *OnlineDP) FindAndProcessSimpleSegments(fragmentSize, fid int) bool {
	//aggregate src into linear fid
	var worker = func(fid int) bool {
		self.AggregateSimpleSegments(fid, fragmentSize)
		return true
	}

	var query = fmt.Sprintf(
		`SELECT DISTINCT fid  FROM %v ORDER BY fid asc;`, self.Src.Table,
	)
	var h, err = self.Src.Query(query)
	if err != nil {
		log.Panic(err)
	}

	var bln bool
	for h.Next() {
		var fid int
		h.Scan(&fid)
		o := worker(fid)
		bln = bln && o
	}
	return bln
}

//Merge segment fragments where possible
func (self *OnlineDP) AggregateSimpleSegments(fid, fragmentSize int) {
	var temp = self.tempNodeIDTableName(fid)
	self.tempCreateNodeIdTable(temp)
	defer self.tempDropTable(temp)

	self.copyFragmentIdsIntoTempTable(fid, fragmentSize, temp)
	//read from temp table nids
	var nidIter = self.iterTempNIDTable(temp)
	for nidIter.Next() {
		var nid int
		nidIter.Scan(&nid)
		self.processNodeFragment(nid)
	}
}

func (self *OnlineDP) processNodeFragment(nid int) {
	var query = fmt.Sprintf(
		"SELECT id, fid, gob FROM %v WHERE id=%v;", self.Src.Table, nid,
	)
	var h, err = self.Src.Query(query)
	if err != nil {
		log.Panic(err)
	}
	for h.Next() {
		var gob string
		var id, fid int
		h.Scan(&id, &fid, &gob)
		var o = db.Deserialize(gob)
		o.NID, o.FID = id, fid

		var queries = self.fragmentMerger(o)
		for _, q := range queries {
			//fmt.Println(q)
			if _, err := self.Src.Exec(q); err != nil {
				log.Panic(err)
			}
		}
	}
}

func (self *OnlineDP) fragmentMerger(hull *db.Node) []string {
	var queries []string
	var ma, mb *db.Node //mergeable neighbours

	// find context neighbours, swap a, b to favor b (next)
	var nb, na = self.FindContiguousNodeNeighbours(hull)

	//if na is bigger, swap na and nb
	if na != nil && nb != nil {
		if na.Range.Size() > nb.Range.Size() {
			na, nb = nb, na
		}
	}

	if na != nil {
		ma = self.ContiguousFragmentsAtThreshold(
			self.Score, hull, na, dp.NodeGeometry,
		)
	}

	if nb != nil {
		mb = self.ContiguousFragmentsAtThreshold(
			self.Score, hull, nb, dp.NodeGeometry,
		)
	}

	var merged = false
	if ma != nil {
		merged = self.checkMerge(ma, hull, na, &queries)
	}

	if !merged && mb != nil {
		merged = self.checkMerge(mb, hull, nb, &queries)
	}

	if merged {
		queries = append(queries, hull.DeleteSQL(self.Src.Table))
	}

	return queries
}

func (self *OnlineDP) checkMerge(merge, hull, neighb *db.Node, queries *[]string) bool {
	var merged = false
	if self.ValidateMerge(merge, hull.Range, neighb.Range) {
		*queries = append(*queries,
			neighb.DeleteSQL(self.Src.Table),
			merge.InsertSQL(self.Src.Table, self.Src.SRID),
		)
		merged = true
	}
	return merged
}

func (self *OnlineDP) copyFragmentIdsIntoTempTable(fid, fragmentSize int, temp string) {
	var query = fmt.Sprintf(
		"SELECT id FROM %v WHERE fid=%v AND size=%v;", self.Src.Table, fid, fragmentSize,
	)
	var h, err = self.Src.Query(query)
	if err != nil {
		panic(err)
	}
	for h.Next() {
		var nid int
		h.Scan(&nid)
		self.tempInsertInNodeIdTable(temp, nid)
	}
}

func (self *OnlineDP) iterTempNIDTable(temp string) *sql.Rows {
	var query = fmt.Sprintf("SELECT id FROM %v;", temp)
	var h, err = self.Src.Query(query)
	if err != nil {
		panic(err)
	}
	return h
}
