package onlinedp

import (
	"simplex/db"
	"simplex/ctx"
)

//direction relate
func IsDirRelateValid(hull *db.Node, ctx *ctx.ContextGeometry) bool {
	return DirectionRelate(hull.Coordinates, ctx.Geom)
}
