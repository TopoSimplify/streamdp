package main

import (
	"os"
	"log"
	"fmt"
	"simplex/streamdp/data"
)

func aggregatePings(msg *data.PingMsg) error {
	if msg.KeepAlive && len(msg.Ping) > 0 {
		var ping = &data.Ping{}
		var err = data.Deserialize(msg.Ping, ping)
		if err != nil {
			return err
		}

		var id = int(ping.MMSI)
		if VesselHistory.Get(id) == nil {
			VesselHistory.Set(id, NewOPW(Options, Type, Offseter))
		}
		VesselHistory.Get(id).Push(ping)
	} else if !msg.KeepAlive {
		var id = msg.Id
		if id > 0 && VesselHistory.Get(id) != nil {
			VesselHistory.Get(id).Done()
			fmt.Println("writing to file : id : ", id)
			writeToFile(id, VesselHistory.Get(id))
		}
	}
	return nil
}

func writeToFile(id int, opw *OPW) {
	var fname = fmt.Sprintf("/home/titus/01/godev/src/simplex/streamdp/mmdata/%v.txt", id)
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Size of Nodes : ", len(opw.Nodes))

	for _, n := range opw.Nodes {
		if _, err := f.WriteString(n.Polyline().Geometry.WKT() + "\n"); err != nil {
			log.Fatal(err)
		}
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
