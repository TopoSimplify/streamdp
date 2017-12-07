package main

import (
	"os"
	"fmt"
	"flag"
	"sync"
	"time"
	"log"
	"os/signal"
	"math/rand"
	"simplex/streamdp/data"
	"github.com/intdxdt/fan"
	"github.com/intdxdt/fileglob"
	"net/http"
	"bytes"
	"io/ioutil"
	"runtime"
	"strings"
)

var Port int
var Host string
var Address string

const concurProcs = 8

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.IntVar(&Port, "port", 8000, "host port")
	flag.StringVar(&Host, "host", "localhost", "host address")
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var msisDir = "/home/titus/01/godev/src/simplex/streamdp/mmsis"
	var ignoreDirs = []string{".git", ".idea"}
	var filter = []string{"toml"}
	Address = fmt.Sprintf("http://%v:%v/ping", Host, Port)
	vesselPings(msisDir, filter, ignoreDirs, concurProcs)
}

func vesselPings(dir string, filter, ignoreDirs []string, batchSize int) {
	var datafileStream = make(chan interface{})
	var exit = make(chan struct{})
	defer close(exit)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

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

	var wg sync.WaitGroup

	var done = make(chan struct{})

	//mmsi vessel
	vessel := func(v *data.Vessel) {
		defer wg.Done()

		for _, loc := range v.Trajectory {
			select {
			case <-exit:
				return
			default:
				dtm, err := time.Parse(time.RFC3339, loc.Time)
				if err != nil {
					panic(err)
				}

				token, err := data.Serialize(data.Pings{
					MMSI:   v.MMSI,
					Type:   v.Type,
					Course: loc.Course,
					Time:   dtm,
					X:      loc.X,
					Y:      loc.Y,
					Speed:  loc.Speed,
				})

				if err != nil {
					panic(err)
				}

				req, err := http.NewRequest("POST",
					Address, bytes.NewBuffer(
						[]byte(fmt.Sprintf(`{"ping":"%v"}`, token)),
					))
				req.Header.Set("X-Custom-Header", "ping")
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil || resp.StatusCode == 500 {
					body, _ := ioutil.ReadAll(resp.Body)
					log.Println(body)
					panic(err)
				}

				resp.Body.Close()
			}
		}
	}

	//now expand one worker into clones of workers
	go func() {
		defer close(done)
		var buf = make([]*data.Vessel, 0)
		var flush = func() {
			wg = sync.WaitGroup{}
			wg.Add(len(buf))
			for _, v := range buf {
				go vessel(v)
			}
			buf = make([]*data.Vessel, 0)
			wg.Wait()
		}
		for vs := range dataSourceStream {
			buf = append(buf, vs.(*data.Vessel))
			if len(buf) >= batchSize {
				flush()
				fmt.Println("size of buf :", len(buf))
				fmt.Println(strings.Repeat("-", 80))
				time.Sleep(5 * time.Second)
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
