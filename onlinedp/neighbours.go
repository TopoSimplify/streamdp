package onlinedp

import (
	"fmt"
	"log"
	"strings"
	"simplex/db"
	"simplex/ctx"
	"simplex/rng"
	"simplex/streamdp/common"
	"github.com/intdxdt/geom"
	"github.com/paulmach/go.geojson"
)

type NeighbQ struct {
	ID, I, J, FID, Snap int
	NodeTable           string
}

func (self *OnlineDP) FindContiguousNodeNeighbours(node *db.Node) (*db.Node, *db.Node) {
	var query = fmt.Sprintf(`
		SELECT id, fid, node
		FROM  %v
		WHERE fid=%v AND snapshot=%v AND (j=%v OR i=%v) AND id <> %v
		ORDER BY i;
	`,
		self.Src.Table,
		node.FID,
		common.Snap,
		node.Range.I,
		node.Range.J,
		node.NID,
	)

	h, err := self.Src.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer h.Close()

	var idx = 0
	var gob string
	var id, fid int

	var nodes []*db.Node
	for h.Next() {
		h.Scan(&id, &fid, &gob)
		o := db.Deserialize(gob)
		o.NID, o.FID = id, fid
		if idx == 0 || idx == 1 {
			nodes = append(nodes, o)
		} else {
			fmt.Println(query)
			log.Panic("expects only two neighbours : prev and next")
		}
		idx++
	}
	return Neighbours(node, nodes)
}

func (self *OnlineDP) FindNodeNeighbours(node *db.Node, independentPlns bool, excludeRanges ...*rng.Range) []*db.Node {
	var query = `
		SELECT id, fid, node
		FROM  %v
		WHERE  ST_DWithin(ST_GeomFromText('%v', %v), geom, %v)
		AND id <> %v AND snapshot=%v `
	if independentPlns {
		//if idependent polylines then restrict neighbours to this fid
		query += fmt.Sprintf(` AND fid = %v`, node.FID)
	}
	query += ";"

	query = fmt.Sprintf(query,
		self.Src.Table,          //table
		node.WTK, self.Src.SRID, //geom(wkt, srid)
		EpsilonDist,
		node.NID, common.Snap, //id != nid
	)
	var h, err = self.Src.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer h.Close()

	var id, fid int
	var gob string
	var dn *db.Node
	var neighbs = make([]*db.Node, 0)

loopRow:
	for h.Next() {
		h.Scan(&id, &fid, &gob)
		dn = db.Deserialize(gob)
		dn.NID, dn.FID = id, fid

		for _, excl := range excludeRanges {
			if dn.Range.Equals(excl) {
				continue loopRow
			}
		}

		neighbs = append(neighbs, dn)
	}
	return neighbs
}

//find context neighbours
func (self *OnlineDP) FindContextNeighbours(queryWKT string, dist float64) []*ctx.ContextGeometry {
	var ctxs []*ctx.ContextGeometry
	if self.Const == nil {
		return ctxs
	}
	var query = fmt.Sprintf(`
			SELECT ST_AsGeoJson(%v)
			FROM  %v
			WHERE ST_DWithin(ST_GeomFromText('%v', %v), %v, %v);
		`,
		self.Const.Config.GeometryColumn, //as_geojson(geom)
		self.Const.Config.Table,          //table
		queryWKT, self.Const.SRID,        //geom(wkt, srid)
		self.Const.Config.GeometryColumn, //geom column
		dist,                             //dist
	)
	var h, err = self.Const.Query(query)
	if err != nil {
		log.Fatalln(err)
	}
	defer h.Close()

	for h.Next() {
		var g *geojson.Geometry
		h.Scan(&g)
		var gs = geometries(g)
		for _, o := range gs {
			ctxs = append(ctxs, ctx.New(o, 0, -1).AsContextNeighbour())
		}
	}
	return ctxs
}

func geometries(g *geojson.Geometry) []geom.Geometry {
	var gs = make([]geom.Geometry, 0)
	var gtype = strings.ToLower(string(g.Type))

	if gtype == "point" {
		gs = append(gs, point(g.Point))
	} else if gtype == "multipoint" {
		gs = append(gs, multiPoint(g.MultiPoint)...)
	} else if gtype == "linestring" {
		gs = append(gs, line(g.LineString))
	} else if gtype == "multilinestring" {
		gs = append(gs, multiLine(g.MultiLineString)...)
	} else if gtype == "polygon" {
		gs = append(gs, polygon(g.Polygon))
	} else if gtype == "multipolygon" {
		gs = append(gs, multiPolygon(g.MultiPolygon)...)
	}

	return gs
}

func point(pt []float64) *geom.Point {
	return geom.NewPoint(pt)
}

func line(ln [][]float64) *geom.LineString {
	return geom.NewLineString(geom.AsPointArray(ln))
}

func polygon(coords [][][]float64) *geom.Polygon {
	var shells = make([][]*geom.Point, 0)
	for _, ln := range coords {
		shells = append(shells, geom.AsPointArray(ln))
	}
	return geom.NewPolygon(shells...)
}

func multiLine(mlns [][][]float64) []geom.Geometry {
	var gs []geom.Geometry
	for _, ln := range mlns {
		gs = append(gs, line(ln))
	}
	return gs
}

func multiPoint(coords [][]float64) []geom.Geometry {
	var gs []geom.Geometry
	for _, pt := range coords {
		gs = append(gs, point(pt))
	}
	return gs
}

func multiPolygon(coords [][][][]float64) []geom.Geometry {
	var gs []geom.Geometry
	for _, ln := range coords {
		gs = append(gs, polygon(ln))
	}
	return gs
}
