package pt

import (
	"simplex/streamdp/data"
	"github.com/intdxdt/geom"
)

type Pt struct {
	*geom.Point
	Ping *data.Ping
	I    int
}

