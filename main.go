package main

import (
	"os"
	"flag"
	"runtime"
	"os/signal"
	"simplex/opts"
	"simplex/offset"
	"fmt"
)

var Port int
var Host string

const DebugMode = 0
const ReleaseMode = 1
const Error = 500
const Success = 200

var SimpleHistory *SimpleMap
var VesselHistory *History
var Options *opts.Opts
var SimplificationType = NOPW
var Offseter = offset.MaxOffset

func init() {
	VesselHistory = NewHistory()
	SimpleHistory = NewSimpleMap()
	Options = &opts.Opts{Threshold: 5000}

	flag.IntVar(&Port, "port", 8000, "host port")
	flag.StringVar(&Host, "host", "localhost", "host address")
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	var server = NewServer("localhost:8000", ReleaseMode)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<- c
		fmt.Println("singnaled exit ...")
		close(server.Exit)
		os.Exit(0)
	}()

	server.Run()
}
