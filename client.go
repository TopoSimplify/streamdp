//
//  Hello World client.
//  Connects REQ socket to tcp://localhost:5555
//  Sends "Hello" to server, expects "World" back
//

package main

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
)

func main() {
	//  Socket to talk to server
	fmt.Println("Connecting to hello world server...")
	requester, _ := zmq.NewSocket(zmq.REQ)
	defer requester.Close()
	requester.Connect("tcp://localhost:5555")

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
