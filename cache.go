package main

import "simplex/streamdp/pt"

type Cache []*pt.Pt

func (self *Cache) Append(o *pt.Pt) {
	*self = append(*self, o)
}

func (self *Cache) IsEmpty() bool {
	return self.Size() == 0
}

func (self *Cache) Size() int {
	return len(*self)
}

func (self *Cache) Pop() *pt.Pt {
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
