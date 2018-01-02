package main

import (
	"simplex/db"
	"simplex/streamdp/mtrafic"
	"simplex/streamdp/common"
	"log"
	"time"
	"fmt"
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
		//var insertSQL = nds[0].InsertSQL(s.Config.Table, s.Config.SRID, nds...)
		var vals = common.SnapshotNodeColumnValues(s.Src.SRID, common.UnSnap, nds...)
		var insertSQL = db.SQLInsertIntoTable(s.Src.Table, common.NodeColumnFields, vals)
		if _, err := s.Src.Exec(insertSQL); err != nil {
			log.Panic(err)
		}

		if !SimpleHistory.Get(id) {
			SimpleHistory.Set(id)
			go func() {
				defer SimpleHistory.Done(id)
				s.OnlineDP.Simplify(id)
			}()
		}
	}

	if !msg.KeepAlive {
		VesselHistory.Delete(msg.Id)
		go func() {
			defer SimpleHistory.Done(id)
			for SimpleHistory.Get(id) { //wait for current snapshort to complete
				fmt.Println(">>> waiting for current snapshot")
				time.Sleep(1 * time.Second)
			}
			s.OnlineDP.Simplify(id)
			fmt.Println("<< done >>")
		}()
	}

	return nil
}
