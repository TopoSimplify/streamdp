package onlinedp

import "github.com/intdxdt/geom"

type LinearGeometry struct {
	Id          int
	Part        int
	Coordinates []*geom.Point
}

func NewLinearGeometry(id int, part int, coords [][]float64) *LinearGeometry {
	return &LinearGeometry{
		Id:          id,
		Part:        part,
		Coordinates: geom.AsPointArray(coords),
	}
}

