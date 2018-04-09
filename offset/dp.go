package offset

import (
	"math"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/vect"
	"github.com/intdxdt/cart"
)

//euclidean offset distance from dp - anchor line [i, j] to maximum
//vertex at i < k <= j - not maximum offset is may not  be perpendicular
func MaxOffset(coordinates []*geom.Point) (int, float64) {
	var n = len(coordinates) - 1
	var index, offset = n, 0.0
	if n <= 1 {
		return index, offset
	}
	var dist float64
	var a, b = coordinates[0], coordinates[n]
	for k := 1; k < n; k++ { //exclusive range between 0 < k < n
		dist = geom.DistanceToPoint(a, b, coordinates[k])
		if dist >= offset {
			index, offset = k, dist
		}
	}
	return index, offset
}

//computes Synchronized Euclidean Distance
func MaxSEDOffset(coordinates []*geom.Point) (int, float64) {
	var n = len(coordinates) - 1
	var index, offset = n, 0.0
	if n <= 1 {
		return index, offset
	}

	var pt *geom.Point
	var dist, m, ptx, pty, ptt, px, py float64

	var a, b = coordinates[0], coordinates[n]
	var ax, ay = a[geom.X], a[geom.Y]
	var at, bt = a[geom.Z], b[geom.Z]

	var v = vect.NewVector(a, b)
	var mij = v.Magnitude()
	var fb = cart.Direction(v[0],v[1])
	var vx, vy = math.Cos(fb), math.Sin(fb)
	var dt = bt - at

	for k := 1; k < n; k++ { //exclusive range between 0 < k < n
		pt = coordinates[k]
		ptx, pty, ptt = pt[geom.X], pt[geom.Y], pt[geom.Z]

		m = (mij / dt) * (ptt - at)
		px, py = ax+(m*vx), ay+(m*vy)
		dist = math.Hypot(ptx-px, pty-py)

		if dist >= offset {
			index, offset = k, dist
		}
	}
	return index, offset
}
