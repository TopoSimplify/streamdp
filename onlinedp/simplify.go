package onlinedp

import (
	"log"
	"fmt"
)

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

	for i := 0; i < 2; i++ {
		log.Println(fmt.Sprintf("merging simple fragments:%v ... #%v", MergeFragmentSize, i))
		self.FindAndProcessSimpleSegments(MergeFragmentSize)
	}

	for i := 0; i < 2; i++ {
		log.Println(fmt.Sprintf("merging simple fragments:%v ... #%v", 2, i))
		self.FindAndProcessSimpleSegments(2)
	}

	for i := 0; i < 2; i++ {
		log.Println(fmt.Sprintf("merging simple fragments:%v ... #%v", 3, i))
		self.FindAndProcessSimpleSegments(3)
	}
	//save simplification
	self.SaveSimplification()
	//drop node table
	//self.Src.DeleteTable(self.Src.NodeTable)
}
