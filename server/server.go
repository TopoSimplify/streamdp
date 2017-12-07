package main

import (
	"fmt"
	"log"
	"flag"
	"simplex/streamdp/data"
	"gopkg.in/gin-gonic/gin.v1"
	"os"
)

var Port int
var Host string

const Error = 500
const Success = 200

type PingMsg struct {
	Ping string `json:"ping" binding:"required"`
}

func init() {
	flag.IntVar(&Port, "port", 8000, "host port")
	flag.StringVar(&Host, "host", "localhost", "host address")
}

func trafficHandler(ctx *gin.Context) {
	var message string
	var pingEv PingMsg
	if err := ctx.BindJSON(&pingEv); err == nil {
		message = pingEv.Ping
	}

	if len(message) > 0 {
		var ping = data.Pings{}
		var err = data.Deserialize(message, &ping)
		if err != nil {
			log.Println(err)
			ctx.JSON(Error, gin.H{"message": "error"})
			return
		}

		fmt.Println(ping)
		writeToFile(ping)

		ctx.JSON(Success, gin.H{"message": "success"})
	} else {
		ctx.JSON(Error, gin.H{"message": "error"})
	}
}

func writeToFile(p data.Pings) {
	var fname = fmt.Sprintf("/home/titus/01/godev/src/simplex/streamdp/mmdata/%v.txt", int(p.MMSI))
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.WriteString(p.String() + "\n"); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	//http.HandleFunc("/traffic", trafficHandler)
	router := gin.Default()
	router.POST("/ping", trafficHandler)
	router.Run(fmt.Sprintf(":%v", Port))
}
