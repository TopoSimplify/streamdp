package main

import (
	"fmt"
	"time"
	"flag"
	zmq "github.com/pebbe/zmq4"
)

var Port int

func init() {
	flag.IntVar(&Port, "port", 5555, "listening port")
}

func main() {
	flag.Parse()
	server, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		panic(err)
	}
	defer server.Close()

	server.Bind(fmt.Sprintf("tcp://*:%v", Port))

	for {
		msg, _ := server.Recv(0)
		fmt.Println("Received ", msg)
		time.Sleep(time.Second)
	}
}
