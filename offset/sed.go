package offset

import (
	"github.com/intdxdt/math"
	"github.com/intdxdt/vect"
	"github.com/intdxdt/geom"
)

//Maximum SED offset distance
func MaxSEDOffset(coordinates geom.Coords) (int, float64) {
	return maxSEDOffset(coordinates, hypot)
}

//Square Maximum SED offset distance
func SqureMaxSEDOffset(coordinates geom.Coords) (int, float64) {
	return maxSEDOffset(coordinates, squareHypot)
}

//@formatter:off
//computes Synchronized Euclidean Distance
func maxSEDOffset(coordinates geom.Coords, hypotFn func(float64, float64) float64) (int, float64) {
	var n = coordinates.Len() - 1
	var index, offset = n, 0.0
	if n <= 1 {
		return index, offset
	}

	var pt *geom.Point
	var dist, m, ptx, pty, ptt, px, py float64

	var a, b = coordinates.Pt(0), coordinates.Pt(n)
	var ax, ay = a[geom.X], a[geom.Y]
	var at, bt = a[geom.Z], b[geom.Z]

	var v = vect.NewVector(*a, *b)
	var mij = v.Magnitude()
	var fb = direction(v[geom.X], v[geom.Y])
	var vx, vy = math.Cos(fb), math.Sin(fb)
	var dt = bt - at

	for k := 1; k < n; k++ { //exclusive range between 0 < k < n
		pt = coordinates.Pt(k)
		ptx, pty, ptt = pt[geom.X], pt[geom.Y], pt[geom.Z]

		m = (mij / dt) * (ptt - at)
		px, py = ax+(m*vx), ay+(m*vy)
		dist = hypotFn(ptx-px, pty-py)

		if dist >= offset {
			index, offset = k, dist
		}
	}
	return index, offset
}

//Dir computes direction in radians - counter clockwise from x-axis.
func direction(x, y float64) float64 {
	var d = math.Atan2(y, x)
	if d < 0 {
		d += math.Tau
	}
	return d
}

func squareHypot(p, q float64) float64 {
	return (p * p) + (q * q)
}

func hypot(p, q float64) float64 {
	if p < 0 {
		p = -p
	}
	if q < 0 {
		q = -q
	}
	if p < q {
		p, q = q, p
	}
	if p == 0 {
		return 0
	}
	q = q / p
	return p * math.Sqrt(1+q*q)
}
