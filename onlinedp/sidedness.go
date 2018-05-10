package onlinedp

import (
	"github.com/TopoSimplify/ctx"
	"github.com/TopoSimplify/homotopy"
	"github.com/intdxdt/geom"
)

//DirectionRelate Relate
func DirectionRelate(coordinates []*geom.Point, contexts *ctx.ContextGeometries) bool {
	return homotopy.Homotopy(coordinates, contexts)
}
