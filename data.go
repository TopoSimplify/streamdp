package main

import (
	"fmt"
	"github.com/naoina/toml"
	"github.com/intdxdt/fileutil"
	"github.com/intdxdt/fileglob"
)

type Ping struct {
	Course float64 `toml:"course"`
	Time   string  `toml:"time"`
	X      float64 `toml:"x"`
	Y      float64 `toml:"y"`
	Speed  float64 `toml:"speed"`
}

type Vessel struct {
	MMSI       float64 `toml:"mmsi"`
	Type       float64 `toml:"type"`
	Geography  string  `toml:"geog"`
	Trajectory []*Ping `toml:"traj"`
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

func main() {
	var vessels, err = fileglob.Glob(
		"/home/titus/01/godev/src/simplex/streamdp/mmsis",
		[]string{"toml"}, false, []string{".git", ".idea"},
	)
	if err != nil {
		panic(err)
	}
	var dats = readAllVessels(vessels)
	fmt.Println(len(dats))
	//vs := readMMSIToml("/home/titus/01/godev/src/simplex/streamdp/mmsis/212773000.toml")
	//fmt.Println(vs)
}
