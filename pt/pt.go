package pt

import (
	"simplex/streamdp/mtrafic"
	"github.com/intdxdt/geom"
)

type Pt struct {
	*geom.Point
	Ping *mtrafic.Ping
	I    int
}

