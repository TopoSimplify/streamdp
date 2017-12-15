package onlinedp

import (
	"fmt"
	"log"
	"simplex/db"
	"simplex/dp"
)

type LnrFeat struct {
	FID, Part int
}

//Find and merge simple segments
func (self *OnlineDP) FindAndProcessSimpleSegments(fragmentSize int) bool {
	var worker = func(o *LnrFeat) bool {
		//aggregate src into linear fid and parts
		self.AggregateSimpleSegments(o.FID, o.Part, fragmentSize)
		return true
	}

	var query = fmt.Sprintf(
		`SELECT DISTINCT fid, part  FROM %v ORDER BY fid asc, part asc;`,
		self.Src.NodeTable,
	)
	var h, err = self.Src.Query(query)
	if err != nil {
		log.Fatalln(err)
	}
	var fid, part int
	var bln bool
	for h.Next() {
		h.Scan(&fid, &part)
		o := worker(&LnrFeat{FID: fid, Part: part})
		bln = bln && o
	}
	return bln
}

//Merge segment fragments where possible
func (self *OnlineDP) AggregateSimpleSegments(fid, part, fragmentSize int) {

	var worker = func(hull *db.Node) []string {
		// find context neighbours
		var prev, nxt = self.FindContiguousNodeNeighbours(hull)
		// find mergeable neighbours contiguous
		var mergePrev, mergeNxt *db.Node

		if prev != nil {
			mergePrev = self.ContiguousFragmentsAtThreshold(
				self.Score, prev, hull, dp.NodeGeometry,
			)
		}

		if nxt != nil {
			mergeNxt = self.ContiguousFragmentsAtThreshold(
				self.Score, hull, nxt, dp.NodeGeometry,
			)
		}

		var merged bool
		var queries []string
		//nxt, prev
		if !merged && mergeNxt != nil {
			if self.ValidateMerge(mergeNxt, hull.Range, nxt.Range) {
				queries = append(queries,
					nxt.DeleteSQL(self.Src.NodeTable),
					mergeNxt.InsertSQL(self.Src.NodeTable, self.Src.SRID),
				)
				merged = true
			}
		}

		if !merged && mergePrev != nil {
			//prev cannot exist since moving from left --- right
			if self.ValidateMerge(mergePrev, hull.Range, prev.Range) {
				queries = append(queries,
					prev.DeleteSQL(self.Src.NodeTable),
					mergePrev.InsertSQL(self.Src.NodeTable, self.Src.SRID),
				)
				merged = true
			}
		}

		if merged {
			queries = append(queries, hull.DeleteSQL(self.Src.NodeTable))
		}
		return queries
	}

	var query = fmt.Sprintf(
		"SELECT id, fid, gob FROM %v WHERE fid=%v AND part=%v AND size=%v;",
		self.Src.NodeTable, fid, part, fragmentSize,
	)
	var h, err = self.Src.Query(query)
	if err != nil {
		panic(err)
	}
	for h.Next() {
		var gob string
		var id, fid int
		h.Scan(&id, &fid, &gob)
		var o = db.Deserialize(gob)
		o.NID, o.FID = id, fid

		var queries = worker(o)
		for _, q := range queries {
			_, err := self.Src.Exec(q)
			if err != nil {
				panic(err)
			}
		}
	}

}
