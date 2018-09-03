package main

import (
	"log"
	"fmt"
	"sync"
	"spinner"
	"github.com/TopoSimplify/data/store"
	"github.com/intdxdt/fan"
	"github.com/TopoSimplify/streamdp/mtrafic"
	"github.com/intdxdt/fileglob"
)

//vessel pings
func vesselPings(dir string, filter, ignoreDirs []string, batchSize int) {
	var datafileStream = make(chan interface{})
	var exit = make(chan struct{})
	defer close(exit)

	fmt.Println("\033c")
	go func() {
		var s = spinner.NewSpinner("vessel pings ...", exit)
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
		return mtrafic.ReadMTraj(filepath)
	}

	var dataSourceStream = fan.Stream(
		datafileStream, worker, concurProcs, exit,
	)

	var done = make(chan struct{})

	var vessel = func(v *store.MTraj, wg *sync.WaitGroup) {
		var id = v.MMSI
		var expected = len(v.Traj)
		var count = 0

		for _, loc := range v.Traj {
			var ping = mtrafic.Ping{
				MMSI:   id,
				Time:   loc.Time,
				X:      loc.X,
				Y:      loc.Y,
				Speed:  loc.Speed,
				Status: loc.Status,
			}

			var token, err = mtrafic.Serialize(ping)
			if err != nil {
				log.Panic(err)
			}

			postToServer(token, id, true)
			count += 1
			//time.Sleep(60 * time.Millisecond)
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
		var buf = make([]*store.MTraj, 0)

		var flush = func() {
			var wg = &sync.WaitGroup{}
			wg.Add(len(buf))
			for _, v := range buf {
				go vessel(v, wg)
			}
			buf = make([]*store.MTraj, 0)
			wg.Wait()
		}

		for vs := range dataSourceStream {
			buf = append(buf, vs.(*store.MTraj))
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
