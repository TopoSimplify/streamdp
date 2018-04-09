package main

import (
	"github.com/intdxdt/geom"
	"fmt"
)

func main() {
	//var coords = []*geom.Point{{1.5, 4.5},{1.5, 4.5},{1.5, 4.5},{1.5, 4.5},{1.5, 4.5},{1.5, 4.5},{1.5, 4.5}}
	//var coords = []*geom.Point{{1.5, 4.5}, {4.5, 4.5}, {7.5, 4.5}, {9.5, 4.5}}
	//var chull = geom.ConvexHull(coords, false)
	var wkt = "LINESTRING ( 146 245, 132 317, 198 371, 267 397, 368 376, 447 339, 502 279, 523 214, 489 168, 375.86733476421915 130.9143720672096, 292 140, 228 197, 219.43513831516063 261.950201110032, 268 322, 322 310, 376 276, 421 209, 473 124, 572 117, 640.130181640881 140.63700179377503, 682 212, 647.5834615512252 255.89355628249535 )"
	var coords = geom.NewLineStringFromWKT(wkt).Coordinates()
	var chull = geom.NewPolygon(coords)

	fmt.Println(chull.WKT())

}
