package main

import (
	"simplex/db"
	"simplex/streamdp/data"
)

func (s *Server) aggregatePings(msg *data.PingMsg) error {
	var id int
	var nds = make([]*db.Node, 0)
	if msg.KeepAlive && len(msg.Ping) > 0 {
		var ping = &data.Ping{}
		var err = data.Deserialize(msg.Ping, ping)
		if err != nil {
			return err
		}

		id = int(ping.MMSI)
		//*db.Node
		if n := VesselHistory.Update(id, ping); n != nil {
			nds = append(nds, n)
		}
	} else if !msg.KeepAlive {
		id = msg.Id
		for _, n := range VesselHistory.MarkDone(id) {
			if n != nil {
				nds = append(nds, n)
			}
		}
	}

	if len(nds) > 0 {
		var insertSQL = nds[0].InsertSQL(s.Config.Table, s.Config.SRID, nds...)
		_, err := s.Src.Exec(insertSQL)
		if err != nil {
			panic(err)
		}
	}

	return nil
}
