package main

import (
	"fmt"
	"simplex/side"
	"github.com/intdxdt/geom"
)

func homoSplit(segment *geom.Segment , coordinates []*geom.Point) (int, int) {
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

func homotopy(coordinates []*geom.Point, gs ...geom.Geometry) bool {
	var bln = true
	var ac, bc []*geom.Point
	var n = len(coordinates) - 1
	var segment = geom.NewSegment(coordinates[0], coordinates[n])
	var i, j = homoSplit(segment, coordinates)

	if i < 0 && j < 0 {
		var gac = geom.NewPolygon(coordinates)
		for k, n := 0, len(gs); bln && k < n; k++ {
			bln = !gac.Intersects(gs[k])
		}
	} else if i > 0 && j > 0 {
		ln := geom.NewSegment(coordinates[i], coordinates[j])
		inters := segment.Intersection(ln)
		ac = append([]*geom.Point{}, coordinates[:i+1]...)
		ac = append(ac, inters[0])
		bc = append([]*geom.Point{inters[0]}, coordinates[j:]...)
		var gac, gbc = geom.NewPolygon(ac), geom.NewPolygon(bc)

		fmt.Println(gac.WKT())
		fmt.Println(gbc.WKT())

		for k, n := 0, len(gs); bln && k < n; k++ {
			bln = !gac.Intersects(gs[k]) && !gbc.Intersects(gs[k])
		}
	} else {
		bln = false
		panic("unhandled condition ")
	}
	return bln
}

func main() {
	//var cwkt = "POLYGON (( 278 307, 270 298, 274 286, 279 272, 301 274, 308 288, 311 304, 296 308, 278 307 ))"
	var cwkt = "POLYGON (( 221 347, 205 334, 221 322, 234 324, 237 342, 221 347 ))"
	//var wkt = "LINESTRING ( 155 171, 207 166, 253 175, 317 171, 367 182, 400 200, 428 249, 417 291, 383 324, 361 332, 333 347, 314 357, 257 383, 204 370, 176 337, 180 305, 214 295, 244 302, 281 332, 316 328, 331 306, 332 291, 315 265, 285 250, 247 261, 231 276, 195 264, 187 230, 216 215, 257 226, 273 217, 273 205, 240 197, 200 200, 178 193, 157 226, 156 246, 151 263, 120 264, 95 249, 89 261, 100 300, 116 359, 139 389, 172 413, 211 425, 256 430, 289 431, 348 427 )"
	var wkt = "LINESTRING ( 155 171, 207 166, 253 175, 317 171, 366.3526316090967 196.8011373660995, 403.24889098928804 217.07380735521562, 449.2399093677868 286.1290060157034, 451.783228702019 325.7276851722446, 415.8591984439269 355.71930167053, 418.6880722510337 402.11283210708103, 389.8335594185446 420.2176244725644, 350.9454024173684 441.6949908346222, 269.44926906112164 441.6949908346222, 233.36391648049494 444.9386180328808, 205.38763189551466 402.77146445551926, 204 370, 176 337, 180 305, 157 226, 95 249, 89 261, 100 300, 116 359, 139 389, 172 413, 211 425, 256 430, 289 431, 348 427 )"
	var g = geom.NewPolygonFromWKT(cwkt)
	var coords = geom.NewLineStringFromWKT(wkt).Coordinates()
	var bln = homotopy(coords, g)
	fmt.Println(bln)
}
