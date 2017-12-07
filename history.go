package main

import (
	"sync"
)

type History struct {
	sync.RWMutex
	m map[int]*OPW
}

func NewHistory() *History {
	return &History{
		m: make(map[int]*OPW, 0),
	}
}

func (h *History) Get(id int) *OPW {
	h.RLock()
	v := h.m[id]
	h.RUnlock()
	return v
}

func (h *History) Set(id int, opw *OPW) {
	h.Lock()
	h.m[id] = opw
	h.Unlock()
}
