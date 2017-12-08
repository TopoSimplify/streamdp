package main

import (
	"sort"
	"simplex/db"
	"simplex/lnr"
	"simplex/rng"
	"simplex/opts"
	"simplex/streamdp/data"
	"github.com/intdxdt/geom"
)

type OPWType int

const (
	NOPW              OPWType = iota
	BOPW
	MinimumCacheLimit = 3
	MaximumCacheLimit = 300
)

type Pt struct {
	*geom.Point
	Ping *data.Ping
	I    int
}

//Type OPW
type OPW struct {
	Id      int
	Part    int
	Nodes   DBNodes
	Options *opts.Opts
	Score   lnr.ScoreFn
	Cache   Cache
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
		Cache:   []*Pt{},
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
	var pt = geom.NewPointXYZ(ping.X, ping.Y, float64(ping.Time.Unix()))
	self.Cache = append(self.Cache, &Pt{Point: pt, Ping: ping, I: self.float})
	if len(self.Cache) < MinimumCacheLimit {
		return node
	}

	var _, val = MaxOffset(self.Cache)
	if self.ScoreRelation(val) || len(self.Cache) >= MaximumCacheLimit {
		if self.Type == NOPW {
			node = self.aggregateNOPW()
		} else if self.Type == BOPW {
			node = self.aggregateBOPW()
		} else {
			panic("unknown open window type")
		}
	}
	return node
}

func (self *OPW) Done() []*db.Node{
	var nd *db.Node
	if (self.Type == NOPW) && !self.Cache.IsEmpty() && !self.Nodes.IsEmpty() {
		nd = self.drainCache(self.Nodes.Pop())
		self.Nodes.Append(nd)
	} else if (self.Type == BOPW) && !self.Cache.IsEmpty() && !self.Nodes.IsEmpty() {
		nd = self.drainCache(self.Nodes.Pop())
		self.Nodes.Append(nd)
	} else if self.Cache.Size() > 1 && self.Nodes.IsEmpty() {
		nd = self.createNode(self.cacheAsPoints(), self.anchor, self.float)
		self.Nodes.Append(nd)
	}
	return self.Nodes.AsSlice()
}

func (self *OPW) lastVal() *Pt {
	return self.Cache[len(self.Cache)-1]
}

func (self *OPW) popMaturedNode() *db.Node {
	var node *db.Node
	if len(self.Nodes) >= 2 {
		node = self.Nodes.PopLeft()
	}
	return node
}

func (self *OPW) aggregateNOPW() *db.Node {
	var nth = self.lastVal()
	var coords = self.cacheAsPoints()
	var nd = self.createNode(coords, self.anchor, self.float)

	self.Nodes.Append(nd)

	self.emptyCache()
	self.Cache = append(self.Cache, nth)
	self.anchor = nth.I

	return self.popMaturedNode()
}

func (self *OPW) aggregateBOPW() *db.Node {
	var last, nth *Pt
	last, self.Cache = Pop(self.Cache)
	nth = self.lastVal()

	var coords = self.cacheAsPoints()
	var nd = self.createNode(coords, self.anchor, self.float)

	self.Nodes.Append(nd)

	self.emptyCache()
	self.Cache = append(self.Cache, nth)
	self.Cache = append(self.Cache, last)
	self.anchor = nth.I

	return self.popMaturedNode()
}

func (self *OPW) drainCache(nd *db.Node) *db.Node {
	var xrng = []int{nd.Range.I, nd.Range.J}
	var n = len(self.Cache)
	var rest = make([]*Pt, n, n)

	//copy Cache
	copy(rest, self.Cache)
	for _, pt := range rest {
		xrng = append(xrng, pt.I)
	}
	sort.Ints(xrng)

	//copy node coordinates
	var cache = make([]*geom.Point, len(nd.Polyline().Coordinates))
	copy(cache, nd.Polyline().Coordinates)

	//add rest to node coords
	for _, pt := range rest {
		cache = append(cache, pt.Point)
	}

	//new range
	var r = rng.NewRange(xrng[0], xrng[len(xrng)-1])
	//new node
	nd = db.New(cache, r, self.Id, self.Part, NodeGeometry)
	return nd
}

func (self *OPW) emptyCache() {
	self.Cache = []*Pt{}
}

func (self *OPW) cacheAsPoints() []*geom.Point {
	var n = len(self.Cache)
	var coords = make([]*geom.Point, n, n)
	for i := 0; i < n; i++ {
		coords[i] = self.Cache[i].Point
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

func Pop(a []*Pt) (*Pt, []*Pt) {
	var v *Pt
	var n int
	if len(a) == 0 {
		return nil, a
	}
	n = len(a) - 1
	v, a[n] = a[n], nil
	return v, a[:n]
}
