package main

import (
	"simplex/db"
	"simplex/streamdp/mtrafic"
)

func (s *Server) aggregatePings(msg *mtrafic.PingMsg) error {
	var id int
	var nds = make([]*db.Node, 0)
	if msg.KeepAlive && len(msg.Ping) > 0 {
		var ping = &mtrafic.Ping{}
		var err = mtrafic.Deserialize(msg.Ping, ping)
		if err != nil {
			return err
		}

		id = int(ping.MMSI)
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
		if _, err := s.Src.Exec(insertSQL); err != nil {
			panic(err)
		}
	}

	if !msg.KeepAlive {
		VesselHistory.Delete(msg.Id)
	}

	return nil
}
