package main

import (
	"bytes"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"simplex/streamdp/mtrafic"
)

//post to server
func postToServer(token string, mmsi int, keepAlive bool) {
	var pmsg = mtrafic.PingMsg{Id: mmsi, Ping: token, KeepAlive: keepAlive}
	var msg, err = json.Marshal(pmsg)

	if err != nil {
		log.Panic(err)
	}

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

	body, err := ioutil.ReadAll(resp.Body)
	if resp != nil {
		resp.Body.Close()
	}

	return body, err
}
