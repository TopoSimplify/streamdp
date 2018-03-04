package sidedness

import (
	"github.com/intdxdt/geom"
	"simplex/side"
)

//Direction Relate

func Homotopy(coordinates []*geom.Point, g geom.Geometry) bool {
	var bln = true
	var ac, bc []*geom.Point
	var n = len(coordinates) - 1
	var segment = geom.NewSegment(coordinates[0], coordinates[n])
	var i, j = homoSplit(segment, coordinates)

	if i < 0 && j < 0 {
		var gac = geom.NewPolygon(coordinates)
		if gac.Intersects(g) && segment.Intersects(g) {
			bln = true
		} else {
			bln = !gac.Intersects(g)
		}
	} else {
		ln := geom.NewSegment(coordinates[i], coordinates[j])
		inters := segment.Intersection(ln)
		ac = append([]*geom.Point{}, coordinates[:i+1]...)
		ac = append(ac, inters[0])
		bc = append([]*geom.Point{inters[0]}, coordinates[j:]...)

		var gac, gbc = geom.NewPolygon(ac), geom.NewPolygon(bc)
		var segIntersects = segment.Intersects(g)

		if (segIntersects && gac.Intersects(g)) || (segIntersects && gbc.Intersects(g)) {
			bln = true
		} else {
			bln = !gac.Intersects(g) && !gbc.Intersects(g)
		}
	}
	return bln
}

func homoSplit(segment *geom.Segment, coordinates []*geom.Point) (int, int) {
	var idx = 1
	var i, j = -1, -1
	var n = len(coordinates) - 1
	var subcoords = coordinates[idx:n]
	var curSide, prevSide *side.Side

	for _, c := range subcoords {
		curSide = segment.SideOf(c)
		if (prevSide != nil) && !(curSide.IsSameSide(prevSide)) {
			ln := geom.NewSegment(coordinates[idx-1], coordinates[idx])
			if segment.Intersects(ln) {
				i, j = idx-1, idx
			}
		}
		prevSide = curSide
		idx += 1
	}
	return i, j
}
