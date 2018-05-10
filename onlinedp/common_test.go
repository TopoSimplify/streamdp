package onlinedp

import (
	"time"
	"log"
	"fmt"
	"math/rand"
	"github.com/intdxdt/geom"
	"github.com/TopoSimplify/db"
	"github.com/TopoSimplify/dp"
	"github.com/TopoSimplify/rng"
	"github.com/TopoSimplify/pln"
	"github.com/TopoSimplify/streamdp/common"
)

const ServerCfg = "../resource/src.toml"
const TestDBName = "test_online_db"
const TestTable = "node_tbl"

func init() {
	var err error
	if err != nil {
		log.Fatalln(err)
	}
	rand.Seed(time.Now().UnixNano())
}

func linearCoords(wkt string) []*geom.Point {
	return geom.NewLineStringFromWKT(wkt).Coordinates()
}

//hull geom
func hullGeom(coords []*geom.Point) geom.Geometry {
	var g geom.Geometry

	if len(coords) > 2 {
		g = geom.NewPolygon(coords)
	} else if len(coords) == 2 {
		g = geom.NewLineString(coords)
	} else {
		g = coords[0].Clone()
	}
	return g
}

func createHulls(indxs [][]int, coords []*geom.Point) []*db.Node {
	poly := pln.New(coords)
	hulls := make([]*db.Node, 0)
	for _, o := range indxs {
		r := rng.NewRange(o[0], o[1])
		n := db.NewDBNode(poly.SubCoordinates(r), r, 1, hullGeom)
		hulls = append(hulls, n)
	}
	return hulls
}

func createNodes(indxs [][]int, coords []*geom.Point) []*db.Node {
	poly := pln.New(coords)
	hulls := make([]*db.Node, 0)
	var fid = rand.Intn(100)
	for _, o := range indxs {
		var r = rng.NewRange(o[0], o[1])
		//var dpnode = newNodeFromPolyline(poly, r, dp.NodeGeometry)
		var n = db.NewDBNode(poly.SubCoordinates(r), r, fid, dp.NodeGeometry, "x7")
		hulls = append(hulls, n)
	}
	return hulls
}

func insertNodesIntoOnlineTable(src *db.DataSrc, nds []*db.Node) {
	var vals = common.SnapshotNodeColumnValues(src.SRID, common.UnSnap, nds...)
	var insertSQL = db.SQLInsertIntoTable(src.Table, common.NodeColumnFields, vals)
	if _, err := src.Exec(insertSQL); err != nil {
		panic(err)
	}
}

func queryNodesByStatus(src *db.DataSrc, status int) []*db.Node {
	var query = fmt.Sprintf(
		`SELECT id, fid, node FROM %v WHERE status=%v;`,
		src.Table, status,
	)

	var h, err = src.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer h.Close()

	var id, fid int
	var gob string
	var nodes = make([]*db.Node, 0)
	for h.Next() {
		h.Scan(&id, &fid, &gob)
		var o = db.Deserialize(gob)
		o.NID, o.FID = id, fid
		nodes = append(nodes, o)
	}
	return nodes
}
