package onlinedp

import (
	"fmt"
	"log"
	"simplex/db"
	"simplex/dp"
	"database/sql"
	"simplex/streamdp/common"
)

type LnrFeat struct{ FID int }

//Merge node fragments
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
	var query = fmt.Sprintf(`
			SELECT id, fid, node
			FROM %v
			WHERE id=%v;
		`,
		self.Src.Table,
		nid,
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
		merged = self.updateMergeQuery(ma, hull, na, &queries)
	}

	if !merged && mb != nil {
		merged = self.updateMergeQuery(mb, hull, nb, &queries)
	}

	if merged {
		queries = append(queries, hull.DeleteSQL(self.Src.Table))
	}

	return queries
}

func (self *OnlineDP) updateMergeQuery(merge, hull, neighb *db.Node, queries *[]string) bool {
	var merged = false
	if self.ValidateMerge(merge, hull.Range, neighb.Range) {
		var vals = common.SnapshotNodeColumnValues(self.Src.SRID, common.Snap, merge)
		*queries = append(*queries,
			neighb.DeleteSQL(self.Src.Table),
			db.SQLInsertIntoTable(self.Src.Table, common.NodeColumnFields, vals),
		)
		merged = true
	}
	return merged
}

func (self *OnlineDP) copyFragmentIdsIntoTempTable(fid, fragmentSize int, temp string) {
	var query = fmt.Sprintf(`
		SELECT id
		FROM %v
		WHERE fid=%v AND snapshot=%v AND size=%v;
	`,
		self.Src.Table,
		fid, common.Snap, fragmentSize,
	)

	var h, err = self.Src.Query(query)
	if err != nil {
		log.Panic(err)
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
