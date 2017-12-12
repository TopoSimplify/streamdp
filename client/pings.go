package main

import (
	"log"
	"sync"
	"time"
	"simplex/streamdp/data"
	"github.com/intdxdt/fan"
	"github.com/intdxdt/fileglob"
	"spinner"
	"fmt"
)

//vessel pings
func vesselPings(dir string, filter, ignoreDirs []string, batchSize int) {
	var datafileStream = make(chan interface{})
	var exit = make(chan struct{})
	defer close(exit)
	fmt.Println("\033c")
	go func() {
		var s = spinner.NewSpinner("mmsi pings ...", exit)
		s.Start()
	}()

	go func() {
		var vessels, err = fileglob.Glob(
			dir, filter, false, ignoreDirs,
		)
		if err != nil {
			log.Fatalln(err)
		}
		for _, o := range vessels {
			datafileStream <- o
		}
		close(datafileStream)
	}()

	var worker = func(v interface{}) interface{} {
		var filepath = v.(string)
		return data.ReadMMSIToml(filepath)
	}
	var dataSourceStream = fan.Stream(datafileStream, worker, concurProcs, exit)

	var done = make(chan struct{})

	vessel := func(v *data.Vessel, wg *sync.WaitGroup) {
		var id = -9
		var expected = len(v.Trajectory)
		var count = 0

		for _, loc := range v.Trajectory {
			dtm, err := time.Parse(time.RFC3339, loc.Time)
			if err != nil {
				panic(err)
			}

			ping := data.Ping{
				MMSI:   v.MMSI,
				Type:   v.Type,
				Course: loc.Course,
				Time:   dtm,
				X:      loc.X,
				Y:      loc.Y,
				Speed:  loc.Speed,
			}

			token, err := data.Serialize(ping)
			if err != nil {
				panic(err)
			}

			if id < 0 {
				id = int(ping.MMSI)
			}
			postToServer(token, id, true)
			count += 1
		}

		postToServer("", id, false)

		if count != expected {
			log.Panic("invalid size")
		}
		wg.Done()
	}

	//now expand one worker into clones of workers
	go func() {
		defer close(done)
		var buf = make([]*data.Vessel, 0)

		var flush = func() {
			var wg = &sync.WaitGroup{}
			wg.Add(len(buf))
			for _, v := range buf {
				go vessel(v, wg)
			}
			buf = make([]*data.Vessel, 0)
			wg.Wait()
		}

		for vs := range dataSourceStream {
			buf = append(buf, vs.(*data.Vessel))
			if len(buf) >= batchSize {
				flush()
			}
		}

		//flush
		if len(buf) > 0 {
			flush()
		}
	}()

	//wait for all the clones to be done
	//in a new go routine
	<-done
}
