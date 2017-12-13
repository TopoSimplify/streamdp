package onlinedp

import (
	"simplex/db"
	"simplex/ctx"
	"simplex/opts"
)

func ByGeometricRelation(hull *db.Node, cg *ctx.ContextGeometry) bool {
	return IsGeomRelateValid(hull, cg)
}

func ByMinDistRelation(options *opts.Opts, hull *db.Node, cg *ctx.ContextGeometry) bool {
	return IsDistRelateValid(options, hull, cg)
}

func BySideRelation(hull *db.Node, cg *ctx.ContextGeometry) bool {
	return IsDirRelateValid(hull, cg)
}
