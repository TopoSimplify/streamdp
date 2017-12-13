package onlinedp

import (
	"strings"
	"github.com/paulmach/go.geojson"
)

func FlattenLinearGeoms(id int, g *geojson.Geometry) []*LinearGeometry {
	var collection []*LinearGeometry
	var gtype = strings.ToLower(string(g.Type))
	if gtype == "linestring" {
		collection = append(collection, NewLinearGeometry(id, 0, g.LineString))
	} else if gtype == "multilinestring" {
		for i, coords := range g.MultiLineString {
			collection = append(collection, NewLinearGeometry(id, i, coords))
		}
	}
	return collection
}

