package onlinedp

import (
	"log"
	"fmt"
	"bytes"
	"simplex/db"
	"simplex/dp"
	"simplex/rng"
	"simplex/pln"
	"simplex/node"
	"text/template"
	"github.com/intdxdt/geom"
	"math/rand"
	"time"
)



const ServerCfg = "/home/titus/01/godev/src/simplex/streamdp/test/src.toml"

var onlineTemplate *template.Template

var onlineTblTemplate = `
CREATE TABLE IF NOT EXISTS {{.Table}} (
    id  SERIAL NOT NULL,
    i INT NOT NULL,
    j INT NOT NULL,
    size INT CHECK (size > 0),
    fid INT NOT NULL,
    part INT NOT NULL,
    gob TEXT NOT NULL,
    geom GEOMETRY(Geometry, {{.SRID}}) NOT NULL,
    status INT DEFAULT 0,
    CONSTRAINT pid_{{.Table}} PRIMARY KEY (id),
	CONSTRAINT u_constraint UNIQUE (fid, i, j)
) WITH (OIDS=FALSE);
CREATE INDEX idx_i_{{.Table}} ON {{.Table}} (i);
CREATE INDEX idx_j_{{.Table}} ON {{.Table}} (j);
CREATE INDEX idx_size_{{.Table}} ON {{.Table}} (size);
CREATE INDEX idx_fid_{{.Table}} ON {{.Table}} (fid);
CREATE INDEX idx_part_{{.Table}} ON {{.Table}} (part);
CREATE INDEX idx_status_{{.Table}} ON {{.Table}} (status);
CREATE INDEX gidx_{{.Table}} ON {{.Table}} USING GIST (geom);
`

func init() {
	var err error
	onlineTemplate, err = template.New("online_table").Parse(onlineTblTemplate)
	if err != nil {
		log.Fatalln(err)
	}
	rand.Seed(time.Now().UnixNano())
}

func printNodes(nodes []*db.Node) {
	for _, h := range nodes {
		fmt.Println(h.WTK)
	}
}

func linearCoords(wkt string) []*geom.Point {
	return geom.NewLineStringFromWKT(wkt).Coordinates()
}

func createNodes(indxs [][]int, coords []*geom.Point)  []*db.Node{
	poly := pln.New(coords)
	hulls := make([]*db.Node, 0)
	var fid = rand.Intn(100)
	for _, o := range indxs {
		var r = rng.NewRange(o[0], o[1])
		//var dpnode = newNodeFromPolyline(poly, r, dp.NodeGeometry)
		var n = db.NewDBNode(poly.SubCoordinates(r), r, fid, 0, dp.NodeGeometry, "x7")
		hulls = append(hulls, n)
	}
	return hulls
}

//New Node
func newNodeFromPolyline(polyline *pln.Polyline, rng *rng.Range, gfn geom.GeometryFn) *node.Node {
	return node.New(polyline.SubCoordinates(rng), rng, gfn)
}

func createOnlineTable(src *db.DataSrc, cfg *ServerConfig) error {
	var query bytes.Buffer
	var err = onlineTemplate.Execute(&query, cfg)
	if err != nil {
		log.Fatalln(err)
	}
	var tblSQl = fmt.Sprintf(`DROP TABLE IF EXISTS %v CASCADE;`, cfg.Table)
	_, err = src.Exec(tblSQl)
	if err != nil {
		log.Panic(err)
	}
	_, err = src.Exec(query.String())
	return err
}

func insertNodesIntoOnlineTable(src *db.DataSrc, nds []*db.Node) {
	var insertSQL = nds[0].InsertSQL(src.Config.Table, src.SRID, nds...)
	if _, err := src.Exec(insertSQL); err != nil {
		panic(err)
	}
}


func queryNodesByStatus(src *db.DataSrc, status int) []*db.Node{
	var query = fmt.Sprintf(
		`SELECT id, fid, gob FROM %v WHERE status=%v;`,
		src.Table, status,
	)
	var h, err = src.Query(query)
	if err != nil {
		log.Panic(err)
	}
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
