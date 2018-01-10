package mtrafic

import (
	"fmt"
	"log"
	"time"
	"io/ioutil"
	"simplex/data/store"
	"github.com/naoina/toml"
)

type PingMsg struct {
	Id        int    `json:"id"`
	Ping      string `json:"ping"`
	KeepAlive bool   `json:"keepalive"`
}

type Ping struct {
	MMSI   int   `toml:"mmsi"`
	Time   time.Time `toml:"time"`
	X      float64   `toml:"x"`
	Y      float64   `toml:"y"`
	Speed  float64   `toml:"speed"`
	Status int       `toml:"status"`
}

func (p *Ping) String() string {
	return fmt.Sprintf(
		`{ MMSI:%v, Time:%v, X:%v, Y:%v, Speed:%v, Status:%v }`,
		p.MMSI,  p.Time.Unix(), p.X, p.Y, p.Speed, p.Status)
}

func ReadMTraj(fname string) *store.MTraj {
	var mtraj = &store.MTraj{}
	var dat, err = ioutil.ReadFile(fname)
	if err != nil {
		log.Panic(err)
	}
	err = toml.Unmarshal(dat, mtraj)
	if err != nil {
		log.Panic(err)
	}
	return mtraj
}

//func ReadMTraj(fileName string) *Vessel {
//	var vsl = &Vessel{}
//	var txt, err = fileutil.ReadAllOfFile(fileName)
//	if err != nil {
//		panic(err)
//	}
//	err = toml.Unmarshal([]byte(txt), vsl)
//	if err != nil {
//		panic(err)
//	}
//	return vsl
//}

func ReadAllVessels(srcs []string) []*store.MTraj {
	var vessels = make([]*store.MTraj, 0)
	for _, src := range srcs {
		vs := ReadMTraj(src)
		vessels = append(vessels, vs)
	}
	return vessels
}
