package main

import (
	"sort"
	"simplex/db"
	"simplex/lnr"
	"simplex/rng"
	"simplex/opts"
	"simplex/streamdp/pt"
	"simplex/streamdp/data"
	"github.com/intdxdt/geom"
	"simplex/streamdp/offset"
)

type OPWType int

const (
	NOPW              OPWType = iota
	BOPW
	MinimumCacheLimit = 3
	MaximumCacheLimit = 300
)

//SimplificationType OPW
type OPW struct {
	Id      int
	Part    int
	Nodes   DBNodes
	Options *opts.Opts
	Score   lnr.ScoreFn
	cache   Cache
	Type    OPWType

	anchor int
	float  int
}

//Creates a new constrained DP Simplification instance
func NewOPW(options *opts.Opts, opwType OPWType, offsetScore lnr.ScoreFn) *OPW {
	var instance = &OPW{
		Nodes:   make(DBNodes, 0),
		Options: options,
		Score:   offsetScore,
		cache:   make(Cache, 0),
		Type:    opwType,
		anchor:  0,
		float:   -1,
	}
	return instance
}

func (self *OPW) ScoreRelation(val float64) bool {
	return val > self.Options.Threshold
}

func (self *OPW) Push(ping *data.Ping) *db.Node {
	self.float += 1
	var node *db.Node
	var pnt = geom.NewPointXYZ(ping.X, ping.Y, float64(ping.Time.Unix()))
	self.cache = append(self.cache, &pt.Pt{Point: pnt, Ping: ping, I: self.float})
	if len(self.cache) < MinimumCacheLimit {
		return node
	}

	var index, val = offset.OPWMaxOffset(self.cache)
	if self.ScoreRelation(val) || (len(self.cache) >= MaximumCacheLimit) {
		if self.Type == NOPW {
			node = self.aggregateNOPW(index)
		} else if self.Type == BOPW {
			node = self.aggregateBOPW(index)
		} else {
			panic("unknown open window type")
		}
	}
	return node
}

func (self *OPW) Done() []*db.Node {
	var nd *db.Node
	if (self.Type == NOPW) && !self.cache.isEmpty() && !self.Nodes.IsEmpty() {
		nd = self.drainCache(self.Nodes.Pop())
		self.Nodes.Append(nd)
	} else if (self.Type == BOPW) && !self.cache.isEmpty() && !self.Nodes.IsEmpty() {
		nd = self.drainCache(self.Nodes.Pop())
		self.Nodes.Append(nd)
	} else if self.cache.size() > 1 && self.Nodes.IsEmpty() {
		nd = self.createNode(self.cacheAsPoints(), self.anchor, self.float)
		self.Nodes.Append(nd)
	}
	return self.Nodes.AsSlice()
}

func (self *OPW) popMaturedNode() *db.Node {
	var node *db.Node
	if len(self.Nodes) >= 2 {
		node = self.Nodes.PopLeft()
	}
	return node
}

func (self *OPW) floatAnchor() (int, int) {
	return self.cache.first().I, self.cache.last().I
}

func (self *OPW) aggregateNOPW(index int) *db.Node {
	var n int
	var stash = self.cache[index+1:]
	n = len(stash)
	stash = stash[:n:n]

	//restrict from 0 to index
	self.cache = self.cache[:index+1]
	n = len(self.cache)
	self.cache = self.cache[:n:n]

	self.anchor, self.float = self.floatAnchor()
	var nd = self.createNode(self.cacheAsPoints(), self.anchor, self.float)

	self.Nodes.Append(nd)

	var nth = self.cache.last()
	self.cache.empty().append(nth).append(stash...)
	self.anchor, self.float = self.floatAnchor()

	return self.popMaturedNode()
}

func (self *OPW) aggregateBOPW(index int) *db.Node {
	var last = self.cache.pop()                  //pop float
	self.anchor, self.float = self.floatAnchor() //update: anchor, float

	//create node
	var nd = self.createNode(self.cacheAsPoints(), self.anchor, self.float)

	self.Nodes.Append(nd)

	var nth = self.cache.last()
	self.cache.empty().append(nth).append(last)
	self.anchor = nth.I

	return self.popMaturedNode()
}

func (self *OPW) drainCache(nd *db.Node) *db.Node {
	var xrng = []int{nd.Range.I, nd.Range.J}
	var n = len(self.cache)
	var rest = make([]*pt.Pt, n, n)

	//copy cache
	copy(rest, self.cache)
	for _, pnt := range rest {
		xrng = append(xrng, pnt.I)
	}
	sort.Ints(xrng)

	//copy node coordinates
	var cache = make([]*geom.Point, len(nd.Polyline().Coordinates))
	copy(cache, nd.Polyline().Coordinates)

	//add rest to node coords
	for _, pnt := range rest {
		cache = append(cache, pnt.Point)
	}

	//new range
	var r = rng.NewRange(xrng[0], xrng[len(xrng)-1])
	//new node
	nd = db.New(cache, r, self.Id, self.Part, NodeGeometry)
	return nd
}

func (self *OPW) cacheAsPoints() []*geom.Point {
	var n = len(self.cache)
	var coords = make([]*geom.Point, n, n)
	for i := 0; i < n; i++ {
		coords[i] = self.cache[i].Point
	}
	return coords
}

func (self *OPW) createNode(coords []*geom.Point, i, j int) *db.Node {
	return db.New(coords, rng.NewRange(i, j), self.Id, self.Part, NodeGeometry)
}

//hull geom
func NodeGeometry(coordinates []*geom.Point) geom.Geometry {
	var g geom.Geometry
	if len(coordinates) > 2 {
		g = geom.NewPolygon(coordinates)
	} else if len(coordinates) == 2 {
		g = geom.NewLineString(coordinates)
	} else {
		g = coordinates[0].Clone()
	}
	return g
}
