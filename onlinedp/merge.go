package onlinedp

import (
	"log"
	"github.com/TopoSimplify/db"
	"github.com/TopoSimplify/rng"
	"github.com/TopoSimplify/common"
	"github.com/intdxdt/geom"
)

//Merge contiguous fragments based combined score
func (self *OnlineDP) ContiguousFragmentsAtThreshold(
	scoreFn func(geom.Coords) (int, float64),
	ha, hb *db.Node,
	gfn func(geom.Coords) geom.Geometry,
) *db.Node {

	if !ha.Range.Contiguous(hb.Range) {
		log.Panic("node are not contiguous")
	}
	var coordinates = ContiguousCoordinates(ha, hb)

	_, val := scoreFn(coordinates)
	if self.ScoreRelation(val) {
		return contiguousFragments(coordinates, ha, hb, gfn)
	}
	return nil
}

//Merge two ranges
func Range(ra, rb rng.Rng) rng.Rng {
	var ranges = common.SortInts(append(ra.AsSlice(), rb.AsSlice()...))
	// i...[ra]...k...[rb]...j
	return rng.Range(ranges[0], ranges[len(ranges)-1])
}

func ContiguousCoordinates(prev, next *db.Node) geom.Coords {
	if !prev.Range.Contiguous(next.Range) {
		panic("node are not contiguous")
	}

	if next.Range.I < prev.Range.J && next.Range.J == prev.Range.I {
		prev, next = next, prev
	}

	var coordinates = prev.Coordinates.Points()
	var n = len(coordinates) - 1
	coordinates = append(coordinates[:n:n], next.Coordinates.Points()...)
	return geom.Coordinates(coordinates)
}

//Merge contiguous hulls
func contiguousFragments(coordinates geom.Coords, ha, hb *db.Node, gfn func(geom.Coords) geom.Geometry) *db.Node {
	var r = Range(ha.Range, hb.Range)
	// i...[ha]...k...[hb]...j
	return db.NewDBNode(coordinates, r, ha.FID, gfn)
}
