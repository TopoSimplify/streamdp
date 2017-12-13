package onlinedp

import (
	"simplex/db"
	"simplex/rng"
	"simplex/lnr"
	"simplex/common"
	"github.com/intdxdt/geom"
)

//Merge contiguous fragments based combined score
func (self *OnlineDP) ContiguousFragmentsAtThreshold(
	scoreFn lnr.ScoreFn, ha, hb *db.Node, gfn geom.GeometryFn,
) *db.Node {
	if !ha.Range.Contiguous(hb.Range) {
		panic("node are not contiguous")
	}
	var coordinates = ContiguousCoordinates(ha, hb)
	_, val := scoreFn(coordinates)
	if self.ScoreRelation(val) {
		return contiguousFragments(coordinates, ha, hb, gfn)
	}
	return nil
}

//Merge two ranges
func Range(ra, rb *rng.Range) *rng.Range {
	var ranges = common.SortInts(append(ra.AsSlice(), rb.AsSlice()...))
	// i...[ra]...k...[rb]...j
	return rng.NewRange(ranges[0], ranges[len(ranges)-1])
}

func ContiguousCoordinates(prev, next *db.Node) []*geom.Point {
	if !prev.Range.Contiguous(next.Range) {
		panic("node are not contiguous")
	}
	var coordinates = prev.Coordinates
	var n = len(coordinates) - 1
	coordinates = append(coordinates[:n:n], next.Coordinates...)
	return coordinates
}

//Merge contiguous hulls
func contiguousFragments(coordinates []*geom.Point, ha, hb *db.Node, gfn geom.GeometryFn) *db.Node {
	var r = Range(ha.Range, hb.Range)
	// i...[ha]...k...[hb]...j
	return db.New(coordinates, r, ha.FID, ha.Part, gfn)
}
