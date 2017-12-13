package onlinedp

func (self *OnlineDP) Simplify() {
	// 0.while has more deformables : loop
	for self.HasMoreDeformables() {
		// 1.find and mark deformable nodes
		self.FindAndMarkDeformables()
		// 2.mark valid nodes as collapsible
		self.FindAndMarkNullStateAsCollapsible()
		// 3.find and split deformable nodes, set status as nullstate
		self.FindAndSplitDeformables()
		// 4.remove deformable nodes
		self.FindAndCleanUpDeformables()
	}
	self.FindAndProcessSimpleSegments(MergeFragmentSize)
	//save simplification
	self.SaveSimplification()
	//drop node table
	self.Src.DeleteTable(self.Src.NodeTable)
}
