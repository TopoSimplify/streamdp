package main

import (
	"log"
	"sort"
	"simplex/db"
	"simplex/lnr"
	"simplex/rng"
	"simplex/opts"
	"simplex/streamdp/pt"
	"github.com/intdxdt/geom"
	"simplex/streamdp/offset"
	"simplex/streamdp/mtrafic"
)

type OPWType int

var ScoreFn = offset.MaxSEDOffset
var OPWScoreFn = offset.OPWMaxSEDOffset

const (
	NOPW              OPWType = iota
	BOPW
	MinimumCacheLimit = 3
	MaximumCacheLimit = 300
)

const (
	AtAnchor = 1
	Moored   = 5
	Aground  = 6
)

//SimplificationType OPW
type OPW struct {
	Id            int
	Nodes         DBNodes
	Options       *opts.Opts
	Score         lnr.ScoreFn
	Type          OPWType
	MaxCacheLimit int
	cache         Cache
	anchor        int
	float         int
}

//Creates a new constrained DP Simplification instance
func NewOPW(options *opts.Opts, opwType OPWType, offsetScore lnr.ScoreFn, maxCacheSize ...int) *OPW {
	var maxCacheLimit = MaximumCacheLimit
	if len(maxCacheSize) > 0 {
		maxCacheLimit = maxCacheSize[0]
	}
	var instance = &OPW{
		Nodes:         make(DBNodes, 0),
		Options:       options,
		Score:         offsetScore,
		Type:          opwType,
		MaxCacheLimit: maxCacheLimit,
		cache:         make(Cache, 0),
		anchor:        0,
		float:         -1,
	}
	return instance
}

func (self *OPW) ScoreRelation(val float64) bool {
	return val > self.Options.Threshold
}

func (self *OPW) Push(ping *mtrafic.Ping) *db.Node {
	var I = 0
	var node *db.Node
	var pnt = geom.NewPointXYZ(ping.X, ping.Y, float64(ping.Time.Unix()))
	if self.cache.size() > 0 {
		I = self.cache.lastIndex() + 1
	}

	var last *pt.Pt
	if self.cache.size() > 0 {
		last = self.cache.last()
	}
	var rmBool = (ping.Status == AtAnchor) ||
		(ping.Status == Moored) || (ping.Status == Aground)

	var eqBool = (last != nil) && last.Point.Equals2D(pnt)

	if last != nil {
		var rmLBool = (last.Ping.Status == AtAnchor) ||
			(last.Ping.Status == Moored) || (last.Ping.Status == Aground)

		if !rmLBool && rmBool {
			rmBool = false
		}
	}

	if eqBool || rmBool {
		return node
	}

	self.cache.append(&pt.Pt{
		Point: pnt, Ping: ping, I: I,
	})

	if self.cache.size() < MinimumCacheLimit {
		return node
	}

	var index, val = OPWScoreFn(self.cache)

	if self.ScoreRelation(val) {
		if self.Type == NOPW {
			node = self.aggregateNOPW(index)
		} else if self.Type == BOPW {
			node = self.aggregateBOPW(index)
		} else {
			log.Panic("unknown open window type")
		}
	}
	return node
}

func (self *OPW) Done() []*db.Node {
	self.updateFloatAnchor()
	if (self.Type == NOPW) && !self.cache.isEmpty() && !self.Nodes.IsEmpty() {
		self.Nodes.Append(self.drainCache(self.Nodes.Pop()))
	} else if (self.Type == BOPW) && !self.cache.isEmpty() && !self.Nodes.IsEmpty() {
		self.Nodes.Append(self.drainCache(self.Nodes.Pop()))
	} else if self.cache.size() > 1 && self.Nodes.IsEmpty() {
		self.Nodes.Append(self.createNode())
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

func (self *OPW) updateFloatAnchor() {
	if self.cache.size() > 0 {
		self.anchor = self.cache.firstIndex()
		self.float = self.cache.lastIndex()
	}
}

func (self *OPW) aggregateNOPW(index int) *db.Node {
	var stash Cache
	self.cache, stash = self.cache.split(index) //update cache
	self.updateFloatAnchor()                    //update float , anchor

	self.Nodes.Append(self.createNode())

	var nth = self.cache.last()
	self.cache.empty().append(nth).append(stash...)
	self.updateFloatAnchor()

	return self.popMaturedNode()
}

func (self *OPW) aggregateBOPW(index int) *db.Node {
	var flt = self.cache.pop() //pop float
	self.updateFloatAnchor()   //update: anchor, float
	self.Nodes.Append(self.createNode())

	var nth = self.cache.last()
	self.cache.empty().append(nth, flt)
	self.updateFloatAnchor()

	return self.popMaturedNode()
}

func (self *OPW) drainCache(nd *db.Node) *db.Node {
	var xrng = []int{
		nd.Range.I, nd.Range.J, self.cache.first().I, self.cache.last().I,
	}
	sort.Ints(xrng)

	var i, j = nd.Range.I, self.cache.lastIndex()

	//copy node coordinates
	var cache = self.nodeAsPoints(nd)
	//add rest to node coords
	for _, pnt := range self.cache[1:] {
		cache = append(cache, pnt.Point)
	}

	//new node
	return db.NewDBNode(cache, rng.NewRange(i, j), self.Id, NodeGeometry)
}

func (self *OPW) cacheAsPoints() []*geom.Point {
	var n = self.cache.size()
	var coords = make([]*geom.Point, n, n)
	for i := 0; i < n; i++ {
		coords[i] = self.cache[i].Point
	}
	return coords
}

func (self *OPW) nodeAsPoints(nd *db.Node) []*geom.Point {
	var lnrcoords = nd.Polyline().Coordinates
	var coords = make([]*geom.Point, len(lnrcoords))
	copy(coords, lnrcoords)
	return coords
}

func (self *OPW) createNode() *db.Node {
	return db.NewDBNode(
		self.cacheAsPoints(),
		rng.NewRange(self.anchor, self.float),
		self.Id, NodeGeometry,
	)
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
