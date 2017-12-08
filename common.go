package main

import (
	"os"
	"simplex/ctx"
	"simplex/node"
	"path/filepath"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/math"
	"github.com/intdxdt/rtree"
	"github.com/intdxdt/deque"
	"math/rand"
	"time"
)
const GeomColumn    = "geom"
const IdColumn      = "id"
const EpsilonDist   = 1.0e-5

func init(){
	rand.Seed(time.Now().UnixNano())
}

func randTableName(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(b)
}


//Convert slice of interface to ints
func asInts(iter []interface{}) []int {
	ints := make([]int, len(iter))
	for i, o := range iter {
		ints[i] = o.(int)
	}
	return ints
}

func castAsContextGeom(o interface{}) *ctx.ContextGeometry {
	return o.(*ctx.ContextGeometry)
}

func castAsNode(o interface{}) *node.Node {
	return o.(*node.Node)
}

func popLeftHull(que *deque.Deque) *node.Node {
	return que.PopLeft().(*node.Node)
}

//node.Nodes from Rtree boxes
func nodesFromBoxes(iter []rtree.BoxObj) []*node.Node {
	var nodes = make([]*node.Node, 0, len(iter))
	for _, h := range iter {
		nodes = append(nodes, h.(*node.Node))
	}
	return nodes
}

//node.Nodes from Rtree nodes
func nodesFromRtreeNodes(iter []*rtree.Node) []*node.Node {
	var nodes = make([]*node.Node, 0, len(iter))
	for _, h := range iter {
		nodes = append(nodes, h.GetItem().(*node.Node))
	}
	return nodes
}

//hull point compare
func PointIndexCmp(a interface{}, b interface{}) int {
	var self, other = a.(*geom.Point), b.(*geom.Point)
	var d = self[2] - other[2]
	if math.FloatEqual(d, 0.0) {
		return 0
	} else if d < 0 {
		return -1
	}
	return 1
}

func ExecutionDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}
