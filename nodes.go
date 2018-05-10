package main

import "github.com/TopoSimplify/db"

type DBNodes []*db.Node

func (self *DBNodes) Append(o *db.Node) {
	*self = append(*self, o)
}

func (self *DBNodes) AsSlice() []*db.Node{
	return []*db.Node(*self)
}

func (self *DBNodes) IsEmpty() bool {
	return self.Size() == 0
}

func (self *DBNodes) Size() int {
	return len(*self)
}

func (self *DBNodes) PopLeft() *db.Node {
	if self.IsEmpty() {
		panic("attempt to pop from an empty slice")
	}
	var list = *self
	var o    = list[0]
	list[0]  = nil
	list     = list[1:]
	var n    = len(list)
	*self    = list[0:n:n]
	return   o
}

func (self *DBNodes) Pop() *db.Node {
	if self.IsEmpty() {
		panic("attempt to pop from an empty slice")
	}
	var list = *self
	var n = len(list) - 1
	var o = list[n]
	list[n] = nil
	list = list[:n]
	*self = list
	return o
}
