package main

import (
	"github.com/naoina/toml"
	"github.com/intdxdt/fileutil"
	"time"
)

type Pings struct {
	MMSI      int       `toml:"mmsi"`
	Type      int       `toml:"type"`
	Course    float64   `toml:"course"`
	Time      time.Time `toml:"time"`
	X         float64   `toml:"x"`
	Y         float64   `toml:"y"`
	Speed     float64   `toml:"speed"`
}

type Location struct {
	Course float64 `toml:"course"`
	Time   string  `toml:"time"`
	X      float64 `toml:"x"`
	Y      float64 `toml:"y"`
	Speed  float64 `toml:"speed"`
}

type Vessel struct {
	MMSI       float64     `toml:"mmsi"`
	Type       float64     `toml:"type"`
	Geography  string      `toml:"geog"`
	Trajectory []*Location `toml:"traj"`
}

func readMMSIToml(fileName string) *Vessel {
	var vsl = &Vessel{}
	var txt, err = fileutil.ReadAllOfFile(fileName)
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal([]byte(txt), vsl)
	if err != nil {
		panic(err)
	}
	return vsl
}

func readAllVessels(srcs []string) []*Vessel {
	var vessels = make([]*Vessel, 0)
	for _, src := range srcs {
		vs := readMMSIToml(src)
		vessels = append(vessels, vs)
	}
	return vessels
}
