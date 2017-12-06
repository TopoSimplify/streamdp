package main

import (
	"os"
	"fmt"
	"flag"
	"log"
	"net/url"
	"os/signal"
	"math/rand"
	"github.com/gorilla/websocket"
	"time"
)

var Port int
var Host string

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.IntVar(&Port, "port", 8080, "host port")
	flag.StringVar(&Host, "host", "localhost", "host address")
}

func main() {
	flag.Parse()
	var address = fmt.Sprintf("%v:%v", Host, Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var uri = &url.URL{Scheme: "ws", Host: address, Path: "/traffic"}
	log.Printf("connecting to %s", uri.String())

	var id = fmt.Sprintf("client id : %v", rand.Int())

	var conn = dialURI(uri)
	defer conn.Close()

	var done = make(chan struct{})
	defer close(done)

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			err := conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			if err != nil {
				log.Fatalln("write close:", err)
				return
			}
			return
		default:
			err := conn.WriteMessage(websocket.TextMessage, []byte(id))
			if err != nil {
				log.Fatalln("write:", err)
				return
			}
			_, message, err := conn.ReadMessage()

			if err != nil {
				log.Fatalln("read:", err)
				return
			}
			log.Printf("recv: %s", message)

			time.Sleep(1 * time.Second)
		}
	}
}

func dialURI(uri *url.URL) *websocket.Conn {
	log.Printf("connecting to %s", uri.String())
	conn, _, err := websocket.DefaultDialer.Dial(uri.String(), nil)
	if err != nil {
		log.Fatalln("dial:", err)
	}
	return conn
}
