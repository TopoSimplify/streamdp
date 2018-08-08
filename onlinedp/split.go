package onlinedp

import (
	"github.com/TopoSimplify/db"
	"github.com/TopoSimplify/rng"
	"github.com/intdxdt/geom"
)

//split hull at vertex with
//maximum_offset offset -- k
func AtScoreSelection(hull *db.Node, scoreFn func([]*geom.Point) (int, float64),
	gfn geom.GeometryFn) (*db.Node, *db.Node) {

	var coordinates = hull.Coordinates
	var rg = hull.Range
	var i, j = rg.I, rg.J
	var k, _ = scoreFn(coordinates)
	var rk = rg.I + k

	// ------------------------------------------------------------------------------------
	var fid = hull.FID
	var idA, idB = hull.SubNodeIds()
	// i..[ha]..k..[hb]..j
	var ha = db.NewDBNode(coordinates[0:k+1], rng.Range(i, rk), fid, gfn, idA)
	var hb = db.NewDBNode(coordinates[k:], rng.Range(rk, j), fid, gfn, idB)
	// ------------------------------------------------------------------------------------

	return ha, hb
}

//split hull at indices (index, index, ...)
func AtIndex(hull *db.Node, indices []int, gfn geom.GeometryFn) []*db.Node {
	//formatter:off
	var coordinates = hull.Coordinates
	var ranges = hull.Range.Split(indices)
	var subHulls = make([]*db.Node, 0, len(ranges))
	var I = hull.Range.I
	var i, j int
	var fid = hull.FID
	var coords []*geom.Point
	for _, r := range ranges {
		i, j = r.I-I, r.J-I
		coords = coordinates[i:j+1]
		subHulls = append(subHulls, db.NewDBNode(coords, r, fid, gfn))
	}
	return subHulls
}
