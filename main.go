package main

import (
	"fmt"
	"flag"
	"sync"
	"time"
	"math/rand"
	"github.com/intdxdt/fan"
	zmq "github.com/pebbe/zmq4"
	"github.com/intdxdt/fileglob"
	"simplex/streamdp/data"
)

const concurProcs = 4


var Port int

func init() {
	flag.IntVar(&Port, "port", 5555, "listening port")
}

//var dats = read_all_vessels(vessels)
//fmt.Println(len(dats))
//vs := Read_MMSI_Toml("/home/titus/01/godev/src/simplex/streamdp/mmsis/212773000.toml")
//fmt.Println(vs)

func main() {
	var msisDir = "/home/titus/01/godev/src/simplex/streamdp/mmsis"
	var ignoreDirs = []string{".git", ".idea"}
	var filter = []string{"toml"}
	vesselPings(msisDir, filter, ignoreDirs, concurProcs)
}

func vesselPings(dir string, filter, ignoreDirs []string, vesselBatchSize int) {
	var stream = make(chan interface{}, 4*concurProcs)
	var exit = make(chan struct{})
	defer close(exit)

	go func() {
		var vessels, err = fileglob.Glob(
			dir, filter, false, ignoreDirs,
		)
		if err != nil {
			panic(err)
		}
		for _, o := range vessels {
			stream <- o
		}
		close(stream)
	}()
	var worker = func(v interface{}) interface{} {
		return data.ReadMMSIToml(v.(string))
	}
	var sources = fan.Stream(stream, worker, concurProcs, exit)

	var wg sync.WaitGroup
	//set up number of of clones to wait for
	wg.Add(vesselBatchSize)
	var onExit = false
	//assume only one worker reading from input chan
	vessel := func(v *data.Vessel) {
		defer wg.Done()

		client, _ := zmq.NewSocket(zmq.REQ)
		defer client.Close()

		client.Connect(fmt.Sprintf("tcp://localhost:%v",  Port))

		//perform fn here...
		for _, loc := range v.Trajectory {
			var delay = time.Duration(rand.Intn(5))
			time.Sleep(delay * time.Second)

			select {
			case <-exit:
				onExit = true
				return
			default:
				if onExit {
					return
				}
				dtm, err := time.Parse(time.RFC3339, loc.Time)
				if err != nil {
					panic(err)
				}

				var p = data.Pings{
					MMSI:   int(v.MMSI),
					Type:   int(v.Type),
					Course: loc.Course,
					Time:   dtm,
					X:      loc.X,
					Y:      loc.Y,
					Speed:  loc.Speed,
				}

				var tokens = data.Serialize(p)
				res, err := client.Send(tokens, 0)
				if err != nil {
					panic(err)
				}
				fmt.Println(res)
			}
		}
	}

	//now expand one worker into clones of workers
	go func() {
		var buf = make([]*data.Vessel, 0)
		for vs := range sources {
			buf = append(buf, vs.(*data.Vessel))
			if len(buf) == concurProcs {
				for _, v := range buf {
					go vessel(v)
				}
				buf = make([]*data.Vessel, 0)
			}
		}
	}()

	//wait for all the clones to be done
	//in a new go routine
	wg.Wait()
}
