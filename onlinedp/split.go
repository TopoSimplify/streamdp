package onlinedp

import (
	"simplex/db"
	"simplex/rng"
	"simplex/lnr"
	"github.com/intdxdt/geom"
)

//split hull at vertex with
//maximum_offset offset -- k
func AtScoreSelection(hull *db.Node, scoreFn lnr.ScoreFn, gfn geom.GeometryFn) (*db.Node, *db.Node) {
	var coordinates = hull.Coordinates
	var rg = hull.Range
	var i, j = rg.I, rg.J
	var k, _ = scoreFn(coordinates)
	var rk = rg.Index(k)
	// ------------------------------------------------------------------------------------
	var fid, part = hull.FID, hull.Part
	var idA, idB = hull.SubNodeIds()
	// i..[ha]..k..[hb]..j
	var ha = db.New(coordinates[0:k+1], rng.NewRange(i, rk), fid, part, gfn, idA)
	var hb = db.New(coordinates[k:], rng.NewRange(rk, j), fid, part, gfn, idB)
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
	var fid, part = hull.FID, hull.Part
	var coords []*geom.Point
	for _, r := range ranges {
		i, j = r.I-I, r.J-I
		coords = coordinates[i:j+1]
		subHulls = append(subHulls, db.New(coords, r, fid, part, gfn))
	}
	return subHulls
}
