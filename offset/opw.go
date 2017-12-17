package offset

import (
	"simplex/rng"
	"simplex/streamdp/pt"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/vect"
)

//euclidean offset distance from dp - anchor line [i, j] to maximum
//vertex at i < k <= j - not maximum offset is may not  be perpendicular
func OPWMaxOffset(polyline []*pt.Pt) (int, float64) {
	var xrange = rng.NewRange(0, len(polyline)-1)
	var seg = geom.NewSegment(polyline[xrange.I].Point, polyline[xrange.J].Point)
	var index, offset = xrange.J, 0.0

	if xrange.Size() > 1 {
		for _, k := range xrange.ExclusiveStride(1) {
			dist := seg.DistanceToPoint(polyline[k].Point)
			if dist >= offset {
				index, offset = k, dist
			}
		}
	}
	return index, offset
}

//computes Synchronized Euclidean Distance
func OPWMaxSEDOffset(polyline []*pt.Pt) (int, float64) {
	var t = 2
	var xrange = rng.NewRange(0, len(polyline)-1)
	var index, offset = xrange.J, 0.0
	var a, b = polyline[xrange.I].Point, polyline[xrange.J].Point
	var opts = &vect.Options{A: a, B: b, At: &a[t], Bt: &b[t]}
	var segvect = vect.NewVect(opts)

	if xrange.Size() > 1 {
		for _, k := range xrange.ExclusiveStride(1) {
			var pnt = polyline[k].Point
			sedvect := segvect.SEDVector(pnt, pnt[t])
			dist := sedvect.Magnitude()
			if dist >= offset {
				index, offset = k, dist
			}
		}
	}
	return index, offset
}

