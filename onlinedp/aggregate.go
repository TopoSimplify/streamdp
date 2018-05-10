package onlinedp

import (
	"fmt"
	"log"
	"github.com/TopoSimplify/db"
	"github.com/TopoSimplify/dp"
	"github.com/TopoSimplify/streamdp/common"
)

type LnrFeat struct{ FID int }

//Merge node fragments
func (self *OnlineDP) AggregateSimpleSegments(fid, fragmentSize int) {
	var temp = self.tempNodeIDTableName(fid)
	self.tempCreateNodeIdTable(temp)
	defer self.tempDropTable(temp)

	self.copyFragmentIdsIntoTempTable(fid, fragmentSize, temp)

	//read from temp table nids
	var query = fmt.Sprintf("SELECT id FROM %v;", temp)
	var h, err = self.Src.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer h.Close()

	for h.Next() {
		var nid int
		h.Scan(&nid)
		self.processNodeFragment(fid, nid)
	}
}

func (self *OnlineDP) processNodeFragment(fid, nid int) {
	var query = fmt.Sprintf(`
			SELECT id, node
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
	defer h.Close()

	for h.Next() {
		var id int
		var gob string

		h.Scan(&id, &gob)
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

	//mark any merges as collapsible
	self.MarkNullStateAsCollapsible(fid)
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
	defer h.Close()

	for h.Next() {
		var nid int
		h.Scan(&nid)
		self.tempInsertInNodeIdTable(temp, nid)
	}
}
