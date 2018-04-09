package onlinedp

import (
	"simplex/ctx"
	"simplex/homotopy"
	"github.com/intdxdt/geom"
)

//DirectionRelate Relate
func DirectionRelate(coordinates []*geom.Point, contexts *ctx.ContextGeometries) bool {
	return homotopy.Homotopy(coordinates, contexts)
}
