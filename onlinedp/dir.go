package onlinedp

import (
	"github.com/TopoSimplify/db"
	"github.com/TopoSimplify/ctx"
)

//direction relate
func IsDirRelateValid(hull *db.Node, contexts *ctx.ContextGeometries) bool {
	return DirectionRelate(hull.Coordinates, contexts)
}
