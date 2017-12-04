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
)

const concurProcs = 4

var Port int

func init() {
	flag.IntVar(&Port, "port", 5555, "listening port")
}

//var dats = readAllVessels(vessels)
//fmt.Println(len(dats))
//vs := readMMSIToml("/home/titus/01/godev/src/simplex/streamdp/mmsis/212773000.toml")
//fmt.Println(vs)

func main() {
	var msisDir = "/home/titus/01/godev/src/simplex/streamdp/mmsis"
	var ignoreDirs = []string{".git", ".idea"}
	var filter = []string{"toml"}
	streamGenerator(msisDir, filter, ignoreDirs)
}

func vesselClient() {

	// send hello
	requester.Send(msg, 0)
	// Wait for reply:
	reply, _ := requester.Recv(0)
	fmt.Println("Received ", reply)
}

func vesselPings(srcs <-chan interface{}, vesselBatchSize int) {
	var exit = make(chan struct{})
	defer close(exit)

	var vesselPings = func(v interface{}) interface{} {
		var vessel = v.(*Vessel)
		return vessel
	}

	//return fan.Stream(stream, vesselPings, concurProcs, exit)
	//for vessel := range out {}

	var wg sync.WaitGroup
	//set up number of of clones to wait for
	wg.Add(vesselBatchSize)
	var out = make(chan interface{}, vesselBatchSize)
	var onExit = false
	//assume only one worker reading from input chan
	vessel := func(v *Vessel) {
		defer wg.Done()
		requester, _ := zmq.NewSocket(zmq.REQ)
		defer requester.Close()
		requester.Connect(fmt.Sprintf("tcp://localhost:%v", Port))

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
				var p = Pings{
					MMSI:   int(v.MMSI),
					Type:   int(v.Type),
					Course: loc.Course,
					Time:   dtm,
					X:      loc.X,
					Y:      loc.Y,
					Speed:  loc.Speed,
				}
				var tokens = Serialize(p)
				res, err := requester.Send(tokens, 0)
				if err != nil {
					panic(err)
				}
				fmt.Println(res)
			}
		}
	}

	//now expand one worker into clones of workers
	go func() {
		for vs := range srcs {
			v := vs.(*Vessel)
			go vessel(v)
		}
	}()

	//wait for all the clones to be done
	//in a new go routine
	go func() {
		wg.Wait()
		close(out)
	}()

	//return out chan to whoever want to read from it
	return out

}

func streamGenerator(dir string, filter, ignoreDirs []string) <-chan interface{} {
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

	var readMMSI = func(v interface{}) interface{} {
		var src = v.(string)
		vessel := readMMSIToml(src)
		return vessel
	}

	return fan.Stream(stream, readMMSI, concurProcs, exit)
}
