package main

import "simplex/db"

type DBNodes []*db.Node

func (self *DBNodes) Append(o *db.Node) {
	*self = append(*self, o)
}

func (self *DBNodes) IsEmpty() bool {
	return self.Size() == 0
}

func (self *DBNodes) Size() int {
	return len(*self)
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
