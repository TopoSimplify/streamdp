package main

import (
	"flag"
	"runtime"
	"simplex/opts"
	"simplex/offset"
	"os"
	"os/signal"
)

var Port int
var Host string

const DebugMode = 0
const ReleaseMode = 0
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
	var server = NewServer("localhost:8000", DebugMode)
	defer close(server.Exit)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<- c
		server.Exit <- struct{}{}
		os.Exit(0)
	}()

	server.Run()
}
