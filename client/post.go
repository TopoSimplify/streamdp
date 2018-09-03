package main

import (
	"bytes"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/TopoSimplify/streamdp/mtrafic"
	"github.com/pkg/errors"
	"fmt"
)

//post to server
func postToServer(token string, mmsi int, keepAlive bool) {
	var pmsg = mtrafic.PingMsg{Id: mmsi, Ping: token, KeepAlive: keepAlive}
	var msg, err = json.Marshal(pmsg)

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(msg)

	_, err = post(Address, msg)
	if err != nil {
		log.Panic(err)
	}
}

func post(address string, msg []byte) ([]byte, error) {
	var req, err = http.NewRequest("POST", address, bytes.NewBuffer(msg))

	req.Header.Set("X-Custom-Header", "ping")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil && (err != nil || resp.StatusCode == 500 ) {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(body)
		panic(err)
	}

	if resp != nil  {
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return body, err
	}
	return []byte{}, errors.New("invalid response")
}
