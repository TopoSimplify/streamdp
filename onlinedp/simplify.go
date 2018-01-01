package onlinedp

import (
	"fmt"
	"log"
)

func (self *OnlineDP) Simplify(fid int) {
	var snapshotTbl = fmt.Sprintf("temp_snapshot_%v", fid)
	self.tempCreateSnapshotTable(snapshotTbl)
	defer self.tempDropTable(snapshotTbl)

	self.updateSnapshot(fid, snapshotTbl)

	// 0.while has more deformables : loop
	for self.HasMoreDeformables(fid, snapshotTbl) {
		// 1.find and mark deformable nodes
		self.FindAndMarkDeformables(fid, snapshotTbl)
		// 2.mark valid nodes as collapsible
		self.FindAndMarkNullStateAsCollapsible(fid, snapshotTbl)
		// 3.find and split deformable nodes, set status as nullstate
		self.FindAndSplitDeformables(fid, snapshotTbl)
		// 4.remove deformable nodes
		self.FindAndCleanUpDeformables(fid, snapshotTbl)
	}

	self.FindAndProcessSimpleSegments(MergeFragmentSize, fid)
	//save simplification
	self.SaveSimplification(fid)
	//drop node table
	//self.Src.DeleteTable(self.Src.NodeTable)
}

func (self *OnlineDP) updateSnapshot(fid int, snapshotTbl string) {
	var query = fmt.Sprintf(`
		INSERT INTO %v (id, size, gob, status)
		SELECT id, size, gob, status
		FROM %v
		WHERE fid = %v;`,
		snapshotTbl, self.Src.Table, fid,
	)
	if _, err := self.Src.Exec(query); err != nil {
		log.Panic(err)
	}
}
