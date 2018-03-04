package sidedness

import (
	"github.com/intdxdt/geom"
	"simplex/side"
)

func homoRegions(coordinates []*geom.Point) [][]*geom.Point {
	var segment *geom.Segment
	var i, j = 0, len(coordinates)-1
	var gln = geom.NewSegment(coordinates[i], coordinates[j])

	var cache []*geom.Point
	var regions = make([][]*geom.Point, 0)

	for i, c := range coordinates {
		cache = append(cache, c)
		if len(cache) > 2 {
			n := len(cache) - 1
			segment = geom.NewSegment(cache[n-1], cache[n])
			if i == j {
				regions = append(regions, cache)
			} else if segment.Intersects(gln) {
				inters := segment.Intersection(gln)
				nth := cache[n]
				cache = cache[:n:n]
				cache = append(cache, inters[0])
				regions = append(regions, cache)
				cache = []*geom.Point{inters[0], nth}
			}
		}
	}
	return regions
}

//Direction Relate
func IsHomotopic(coordinates []*geom.Point, g geom.Geometry) bool {
	var bln = true
	var envelope []*geom.Point
	var lastEnvelope []*geom.Point

	var regions = homoRegions(coordinates)
	var n = len(regions)
	if n > 1 {
		lastEnvelope = regions[n-1]
	}
	var subregions = regions[:n-1]

	for i, n := 0, len(subregions); i < n; i++ {
		envelope = append(envelope, subregions[i]...)
	}
	if len(envelope) > 0 {
		bln = !(geom.NewPolygon(envelope).Intersects(g))
	}
	if bln && len(lastEnvelope) > 0 {
		bln = !(geom.NewPolygon(lastEnvelope).Intersects(g))
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
		if prevSide == nil {
			prevSide = curSide
		}
		if curSide.Value() != prevSide.Value() {
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
	} else if i > 0 && j > 0 {
		ln := geom.NewSegment(coordinates[i], coordinates[j])
		inters := segment.Intersection(ln)
		ac = append([]*geom.Point{}, coordinates[:i+1]...)
		ac = append(ac, inters[0])
		bc = append([]*geom.Point{inters[0]}, coordinates[j:]...)
		var gac, gbc = geom.NewPolygon(ac), geom.NewPolygon(bc)

		var blnSegInters = segment.Intersects(g)
		if (gac.Intersects(g) && blnSegInters) || (gbc.Intersects(g) && blnSegInters) {
			bln = true
		} else {
			bln = !gac.Intersects(g) && !gbc.Intersects(g)
		}

	} else {
		bln = false
		panic("unhandled condition ")
	}
	return bln
}
