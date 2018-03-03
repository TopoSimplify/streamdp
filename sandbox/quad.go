package main

import (
	"simplex/pln"
	"simplex/seg"
	"simplex/ctx"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/rtree"
	"fmt"
)

//Direction Relate
func DirectionRelate(polyline *pln.Polyline, g geom.Geometry) string {
	var coordinates = polyline.Coordinates
	var i, j = 0, len(coordinates)-1
	var segment *seg.Seg
	var gseg = seg.NewSeg(coordinates[i], coordinates[j], i, j)
	var gln = gseg.Geometry()
	var regions = make([][]*geom.Point, 0)
	var cache []*geom.Point
	for i, coord := range coordinates {
		cache = append(cache, coord)
		if len(cache) > 2 {
			n := len(cache) - 1
			segment = seg.NewSeg(cache[n-1], cache[n], i-1, i)
			if i == j {
				cache = append(cache, cache[0])
				regions = append(regions, cache)
			} else if segment.Intersects(gln) {
				inters := segment.Intersection(gln)
				nth := cache[n]
				cache = cache[:n:n]
				cache = append(cache, inters[0], cache[0])
				regions = append(regions, cache)
				cache = []*geom.Point{inters[0], nth}
			}
		}
	}
	for _, reg := range  regions{
		ply:= geom.NewPolygon(reg)
		fmt.Println(ply.WKT())
	}
	return ""
}

//find if intersects segment
func intersectsQuad(q geom.Geometry, res []*rtree.Node) bool {
	var bln = false
	for _, node := range res {
		c := node.GetItem().(*ctx.ContextGeometry)
		s := c.Geom.(*seg.Seg)
		if q.Intersects(s.Segment) {
			bln = true
			break
		}
	}
	return bln
}

func main() {
	var cwkt = "POLYGON (( 278 307, 270 298, 274 286, 279 272, 301 274, 308 288, 311 304, 296 308, 278 307 ))"
	var wkt = "LINESTRING ( 155 171, 207 166, 253 175, 317 171, 367 182, 400 200, 428 249, 417 291, 383 324, 361 332, 333 347, 314 357, 257 383, 204 370, 176 337, 180 305, 214 295, 244 302, 281 332, 316 328, 331 306, 332 291, 315 265, 285 250, 247 261, 231 276, 195 264, 187 230, 216 215, 257 226, 273 217, 273 205, 240 197, 200 200, 178 193, 157 226, 156 246, 151 263, 120 264, 95 249, 89 261, 100 300, 116 359, 139 389, 172 413, 211 425, 256 430, 289 431, 348 427 )"
	var coords = geom.NewLineStringFromWKT(wkt).Coordinates()
	var ln = pln.New(coords)
	var g = geom.NewPolygonFromWKT(cwkt)
	var relate = DirectionRelate(ln, g)
	fmt.Println(relate)
}
