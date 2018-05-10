package onlinedp

import (
	"github.com/TopoSimplify/db"
	"github.com/TopoSimplify/ctx"
	"github.com/TopoSimplify/opts"
)

//distance relate
func IsDistRelateValid(options *opts.Opts, hull *db.Node, contexts *ctx.ContextGeometries) bool {
	var minDistance = options.MinDist
	var seg = hull.Segment()
	var lnGeom = hull.Polyline().Geometry
	var original, simple float64
	var segGeom = seg
	var g *ctx.ContextGeometry

	var bln = true
	var geometries = contexts.DataView()

	for i, n := 0, contexts.Len(); bln && i < n; i++ {
		g = geometries[i]
		original = lnGeom.Distance(g.Geom)
		simple = segGeom.Distance(g.Geom)

		//if original violates constraint, then simple can
		// >= than original or <= original, either way should be true
		// [original & simple] <= minDistance, then simple cannot be  simple >= minDistance no matter
		// how many vertices introduced
		bln = (simple >= minDistance) || (original < minDistance)
	}
	return bln
}
