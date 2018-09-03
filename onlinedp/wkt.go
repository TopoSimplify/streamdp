package onlinedp

import (
	"strings"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/math"
)

//str polyline
func wktLineString(coords []geom.Point, dim int) string {
	if len(coords) == 0 {
		panic("empty coordinates")
	}
	n := len(coords)
	lnstr := make([]string, n)
	for i := 0; i < n; i++ {
		lnstr[i] = coordStr(&coords[i], dim)
	}
	return "(" + strings.Join(lnstr, ", ") + ")"
}

//coordinate string
func coordStr(pt *geom.Point, dim int) string {
	var token = ""
	var coords = (*pt)[:dim]
	var n = len(coords) - 1
	for i, p := range coords {
		token += math.FloatToString(p)
		if i < n {
			token += " "
		}
	}
	return token
}
