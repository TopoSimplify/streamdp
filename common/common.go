package common

import (
	"os"
	"time"
	"math/rand"
	"simplex/ctx"
	"simplex/node"
	"path/filepath"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/math"
	"github.com/intdxdt/rtree"
	"github.com/intdxdt/deque"
	"simplex/db"
	"fmt"
	"simplex/streamdp/sandbox/s"
)

const (
	UnSnap   = iota
	Snap
)
const EpsilonDist   = 1.0e-5

func init(){
	rand.Seed(time.Now().UnixNano())
}

//Note : column fields corresponding to node.ColumnValues
const SnapNodeColumnFields = "fid, node, geom, i, j, size, snapshot"

func SnapshotNodeColumnValues(srid int, nodes ...*db.Node) [][]string {
	var colVals = func(n *db.Node) []string {
		return []string{
			fmt.Sprintf(`%v`, n.FID),
			fmt.Sprintf(`'%v'`, db.Serialize(n)),
			fmt.Sprintf(`ST_GeomFromText('%v', %v)`, n.WTK, srid),
			fmt.Sprintf(`%v`, n.Range.I),
			fmt.Sprintf(`%v`, n.Range.J),
			fmt.Sprintf(`%v`, n.Range.Size()),
			fmt.Sprintf(`%v`, Snap),
		}
	}
	var vals = make([][]string, 0)
	for _, n := range nodes {
		vals = append(vals, colVals(n))
	}
	return vals
}


func SimpleTable(tbl string) string {
	return fmt.Sprintf(`%v_simple`, tbl)
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
