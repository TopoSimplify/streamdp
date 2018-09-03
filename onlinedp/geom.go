package onlinedp

import (
	"github.com/TopoSimplify/db"
	"github.com/TopoSimplify/ctx"
)

//Checks geometric relation to other context geometries
func IsGeomRelateValid(hull *db.Node, contexts *ctx.ContextGeometries) bool {
	var seg = hull.Segment()
	var ln = hull.Polyline().Geometry()
	var lnGInter, segGInter bool
	var g *ctx.ContextGeometry

	var bln = true
	var geometries = contexts.DataView()

	for i, n := 0, contexts.Len(); bln && i < n; i++ {
		g = geometries[i]
		lnGInter = ln.Intersects(g.Geom)
		segGInter = seg.Intersects(g.Geom)

		bln = !((segGInter && !lnGInter) || (!segGInter && lnGInter) )
	}

	return bln
}
