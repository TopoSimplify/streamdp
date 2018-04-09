package onlinedp

import (
	"simplex/db"
	"simplex/ctx"
	"simplex/opts"
)

func ByGeometricRelation(hull *db.Node, contexts *ctx.ContextGeometries) bool {
	return IsGeomRelateValid(hull, contexts)
}

func ByMinDistRelation(options *opts.Opts, hull *db.Node, contexts *ctx.ContextGeometries) bool {
	return IsDistRelateValid(options, hull, contexts)
}

func BySideRelation(hull *db.Node, contexts *ctx.ContextGeometries) bool {
	return IsDirRelateValid(hull, contexts)
}
