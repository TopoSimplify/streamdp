package onlinedp

import (
	"github.com/TopoSimplify/db"
	"github.com/TopoSimplify/ctx"
	"github.com/TopoSimplify/opts"
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
