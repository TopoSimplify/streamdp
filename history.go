package main

import (
	"sync"
	"github.com/TopoSimplify/db"
	"github.com/TopoSimplify/opts"
	"github.com/TopoSimplify/streamdp/mtrafic"
)

type History struct {
	sync.RWMutex
	m map[int]*OPW
}

func NewHistory() *History {
	return &History{m: make(map[int]*OPW, 0)}
}

func (h *History) Get(id int) *OPW {
	h.RLock()
	v := h.m[id]
	h.RUnlock()
	return v
}

func (h *History) MarkDone(id int) []*db.Node {
	var nodes []*db.Node
	h.Lock()
	// ----------------------------------------------
	if h.m[id] != nil {
		nodes = h.m[id].Done()
	}
	// ----------------------------------------------
	h.Unlock()
	return nodes
}

func (h *History) Update(id int, ping *mtrafic.Ping, options *opts.Opts) *db.Node {
	var node *db.Node
	h.Lock()
	// ----------------------------------------------
	if h.m[id] == nil {
		h.m[id] = NewOPW(options, SimplificationType, Offseter)
		h.m[id].Id = id
	}
	node = h.m[id].Push(ping)
	// ----------------------------------------------
	h.Unlock()
	return node
}

func (h *History) Delete(id int) {
	h.Lock()
	delete(h.m, id)
	h.Unlock()
}

func (h *History) Clear() {
	h.Lock()
	h.m = make(map[int]*OPW, 0)
	h.Unlock()
}
