package onlinedp

import (
	"fmt"
	"log"
	"simplex/streamdp/common"
)

func (self *OnlineDP) Simplify(fid int) {
	self.MarkSnapshot(fid, common.Snap)
	defer self.MarkSnapshot(fid, common.UnSnap)

	// 0.while has more deformables : loop
	for self.HasMoreDeformables(fid) {
		// 1.find and mark deformable nodes
		self.MarkDeformables(fid)
		// 2.mark valid nodes as collapsible
		self.MarkNullStateAsCollapsible(fid)
		// 3.find and split deformable nodes, set status as nullstate
		self.SplitDeformables(fid)
		// 4.remove deformable nodes
		self.CleanUpDeformables(fid)
	}

	//aggregate segments
	self.AggregateSimpleSegments(fid, MergeFragmentSize)

	//save simplification
	self.SaveSimplification(fid)

	//drop node table
	//self.Src.DeleteTable(self.Src.NodeTable)
}

func (self *OnlineDP) updateSnapshot(fid int, snapshotTbl string) {
	var query = fmt.Sprintf(`
		INSERT INTO %v (id, size, node, status)
		SELECT id, size, node, status
		FROM %v
		WHERE fid = %v;`,
		snapshotTbl, self.Src.Table, fid,
	)
	if _, err := self.Src.Exec(query); err != nil {
		log.Panic(err)
	}
}
