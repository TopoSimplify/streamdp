package main

import (
	"github.com/TopoSimplify/db"
	"github.com/TopoSimplify/streamdp/mtrafic"
)

func (server *Server) aggregatePings(msg *mtrafic.PingMsg) error {
	var id int
	var nds = make([]*db.Node, 0)
	var options = server.OnlineDP.Options

	if msg.KeepAlive && len(msg.Ping) > 0 {
		var ping = &mtrafic.Ping{}
		var err = mtrafic.Deserialize(msg.Ping, ping)
		if err != nil {
			return err
		}

		id = ping.MMSI
		var node = VesselHistory.Update(id, ping, options)
		//node
		if node != nil {
			nds = append(nds, node)
		}

	} else if !msg.KeepAlive {
		id = msg.Id
		var nodes = VesselHistory.MarkDone(id)
		for _, n := range nodes {
			if n != nil {
				nds = append(nds, n)
			}
		}
	}

	//send to input stream
	if len(nds) > 0 {
		server.InputStream <- nds
	}

	return nil
}
