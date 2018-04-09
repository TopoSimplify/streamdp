package onlinedp

import (
	"simplex/db"
	"simplex/ctx"
)

//direction relate
func IsDirRelateValid(hull *db.Node, contexts *ctx.ContextGeometries) bool {
	return DirectionRelate(hull.Coordinates, contexts)
}
