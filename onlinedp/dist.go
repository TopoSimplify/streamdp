package onlinedp

import (
	"simplex/ctx"
	"simplex/opts"
	"simplex/db"
)

//distance relate
func IsDistRelateValid(options *opts.Opts, hull *db.Node, ctx *ctx.ContextGeometry) bool {
	var mindist = options.MinDist
	var seg = hull.Segment()
	var lnGeom = hull.Polyline().Geometry

	var segGeom = seg
	var ctxGeom = ctx.Geom

	var origDist = lnGeom.Distance(ctxGeom) // original relate
	var dr = segGeom.Distance(ctxGeom)      // new relate

	bln := dr >= mindist
	if (!bln) && origDist < mindist { //if not bln and origDist <= mindist:
		//if original violates constraint, then simple can
		// >= than original or <= original, either way should be true
		// [original & simple] <= mindist, then simple cannot be  simple >= mindist no matter
		// how many vertices introduced
		bln = true
	}
	return bln
}
