package main

import (
	zmq "github.com/pebbe/zmq4"
	"fmt"
	"flag"
)

var Port int

func init() {
	flag.IntVar(&Port, "port", 5555, "listening port")
}

func main() {
	requester, _ := zmq.NewSocket(zmq.REQ)
	defer requester.Close()
	requester.Connect(fmt.Sprintf("tcp://localhost:%v", Port))

	for r := 0; r < 10; r++ {
		// send hello
		msg := fmt.Sprintf("Hello %d", r)
		fmt.Println("Sending ", msg)
		requester.Send(msg, 0)

		// Wait for reply:
		reply, _ := requester.Recv(0)
		fmt.Println("Received ", reply)
	}
}

func streamGenerator(){

}

