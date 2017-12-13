package onlinedp

import (
	"simplex/db"
	"simplex/ctx"
)

//geometry relate
func IsGeomRelateValid(hull *db.Node, ctx *ctx.ContextGeometry) bool {
	var seg = hull.Segment()
	var lnGeom  = hull.Polyline().Geometry
	var segGeom = seg
	var ctxGeom = ctx.Geom

	var lnGInter = lnGeom.Intersects(ctxGeom)
	var segGInter = segGeom.Intersects(ctxGeom)

	var bln = true
	if (segGInter && !lnGInter) || (!segGInter && lnGInter) {
		bln = false
	}
	// both intersects & disjoint
	return bln
}
