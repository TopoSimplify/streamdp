package main

import (
	"fmt"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/math"
)

const (
	X    = iota
	Y
	Z
	null = -9
)

 //Test whether a point lies inside a ring.
 //The ring may be oriented in either direction.
 //If the point lies on the ring boundary the result of this method is unspecified.
 //This algorithm does not attempt to first check the point against the envelope of the ring.
func completely_in_ring(poly []*geom.Point, pnt *geom.Point) bool {
	var i, i1 int
	var p1, p2 *geom.Point
	var x1, y1, x2, y2, xInt float64
	var p = *pnt
	// for each segment l = (i-1, i), see if it crosses ray from test point in positive x direction.
	var crossings = 0 // number of segment/ray crossings
	var n = len(poly)
	for i = 1; i < n; i++ {
		i1 = i - 1
		p1 = poly[i]
		p2 = poly[i1]

		if ((p1[Y] > p[Y]) && (p2[Y] <= p[Y])) || ((p2[Y] > p[Y]) && (p1[Y] <= p[Y])) {
			x1, y1 = p1[X]-p[X], p1[Y]-p[Y]
			x2, y2 = p2[X]-p[X], p2[Y]-p[Y]
			//segment straddles x axis, so compute intersection with x-axis.
			xInt = float64(math.SignOfDet2(x1, y1, x2, y2)) / (y2 - y1)
			//xsave = xInt
			//crosses ray if strictly positive intersection.
			if xInt > 0.0 {
				crossings++
			}
		}
	}
	//  p is inside if number of crossings is odd.
	return (crossings % 2) == 1
}


//pnp for arbitrary ring
func pointInPolygon(poly []*geom.Point, pt *geom.Point) bool {
	var polyCorners = len(poly)
	var i, j = 0, polyCorners-1
	var oddNodes = false
	var x, y = pt[0], pt[1]

	for i = 0; i < polyCorners; i++ {
		if (poly[i][1] < y && poly[j][1] >= y || poly[j][1] < y && poly[i][1] >= y) && (poly[i][0] <= x || poly[j][0] <= x) {
			if poly[i][0]+(y-poly[i][1])/(poly[j][1]-poly[i][1])*(poly[j][0]-poly[i][0]) < x {
				oddNodes = !oddNodes
			}
		}
		j = i
	}

	return oddNodes
}

func main() {
	//var cwkt = "POLYGON (( 278 307, 270 298, 274 286, 279 272, 301 274, 308 288, 311 304, 296 308, 278 307 ))"
	var cwkt = "POLYGON (( 510 430, 520 420, 530 430, 520 440, 530 440, 520 430, 510 430 ))"
	var coords = geom.NewPolygonFromWKT(cwkt).Coordinates()[0]
	var bln = pointInPolygon(coords, geom.NewPointXY(  525.5869623409229 ,438.78965670087933))
	var bln2 = completely_in_ring(coords, geom.NewPointXY(  525.5869623409229 ,438.78965670087933))
	fmt.Println(bln)
	fmt.Println("in ring ", bln2)
}
