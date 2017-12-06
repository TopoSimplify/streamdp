package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
)

var Port int
var Host string
var upgrader = websocket.Upgrader{} // use default options

func init() {
	flag.IntVar(&Port, "port", 8000, "host port")
	flag.StringVar(&Host, "host", "localhost", "host address")
}

func trafficHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
		}
		if err == nil {
			log.Printf("recv: %s", message)
			err = conn.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
			}
		}
	}
}

func main() {
	flag.Parse()
	var address = fmt.Sprintf("%v:%v", Host, Port)

	http.HandleFunc("/traffic", trafficHandler)

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalln(err)
	}
}
