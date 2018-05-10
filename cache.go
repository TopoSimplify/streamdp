package main

import (
	"github.com/TopoSimplify/streamdp/pt"
)

type Cache []*pt.Pt

func (self *Cache) first() *pt.Pt {
	return (*self)[0]
}

func (self *Cache) last() *pt.Pt {
	return (*self)[self.size()-1]
}

func (self *Cache) firstIndex() int {
	return self.first().I
}

func (self *Cache) lastIndex() int {
	return self.last().I
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

func (self *Cache) clone() Cache {
	var list = *self
	var clone = make(Cache, 0, self.size())
	clone.append(list...)
	return clone
}

func (self *Cache) split(index int) (Cache, Cache) {
	var list = *self
	var before, after = make(Cache, 0), make(Cache, 0)
	for i := range list {
		if i <= index {
			before.append(list[i])
		} else if i > index {
			after.append(list[i])
		}
	}
	return before, after
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
