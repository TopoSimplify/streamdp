package onlinedp

import (
	"simplex/db"
	"simplex/ctx"
	"simplex/side"
	"github.com/intdxdt/geom"
)

//DirectionRelate Relate
func IsDirRelateValid(hull *db.Node, contexts *ctx.ContextGeometries) bool {
	var coordinates = hull.Coordinates
	var n = len(coordinates) - 1
	var g *ctx.ContextGeometry
	var a, b, inters []*geom.Point
	var ln, segment *geom.Segment
	var linestring = geom.NewLineString(coordinates)
	var segInters, lineInters, disjointA, disjointB bool

	var bln = true
	var geometries = contexts.DataView()
	segment = geom.NewSegment(coordinates[0], coordinates[n])
	var i, j = homotopicSplit(segment, coordinates)

	for idx, n := 0, contexts.Len(); bln && idx < n; idx++ {
		g = geometries[idx]
		segInters = segment.Intersects(g.Geom)
		lineInters = linestring.Intersects(g.Geom)

		if segInters && lineInters {
			continue //bln = true : continue
		}

		if i < 0 && j < 0 {
			bln = !geom.NewPolygon(coordinates).Intersects(g.Geom)
			continue
		}

		ln = geom.NewSegment(coordinates[i], coordinates[j])
		inters = segment.Intersection(ln)

		a = append([]*geom.Point{}, coordinates[:i+1]...)
		a = append(a, inters[0])
		b = append([]*geom.Point{inters[0]}, coordinates[j:]...)

		disjointA = !geom.NewPolygon(a).Intersects(g.Geom)
		disjointB = !geom.NewPolygon(b).Intersects(g.Geom)

		bln = disjointA && disjointB

	}
	return bln
}

func homotopicSplit(segment *geom.Segment, coordinates []*geom.Point) (int, int) {
	var i, j = -1, -1
	var curSide, prevSide *side.Side
	var ln *geom.Segment
	var c *geom.Point

	for idx, n := 1, len(coordinates)-1; idx < n; idx++ {
		c = coordinates[idx]
		curSide = segment.SideOf(c)
		if (prevSide != nil) && !(curSide.IsSameSide(prevSide)) {
			ln = geom.NewSegment(coordinates[idx-1], coordinates[idx])
			if segment.Intersects(ln) {
				i, j = idx-1, idx
			}
		}
		prevSide = curSide
	}
	return i, j
}
