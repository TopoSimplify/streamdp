package main

import "simplex/streamdp/pt"

type Cache []*pt.Pt

func (self *Cache) first() *pt.Pt {
	return (*self)[0]
}

func (self *Cache) last() *pt.Pt {
	return (*self)[self.size()-1]
}

func (self *Cache) size() int {
	return len(*self)
}

func (self *Cache) isEmpty() bool {
	return self.size() == 0
}

func (self *Cache) empty() *Cache {
	*self = make(Cache, 0)
	return self
}

func (self *Cache) append(pts ...*pt.Pt) *Cache {
	*self = append(*self, pts...)
	return self
}

func (self *Cache) pop() *pt.Pt {
	if self.isEmpty() {
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
