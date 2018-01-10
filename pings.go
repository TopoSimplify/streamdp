package main

import (
	"log"
	"fmt"
	"time"
	"simplex/db"
	"simplex/streamdp/common"
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

		id = ping.MMSI
		var node = VesselHistory.Update(id, ping)
		//node
		if node != nil {
			nds = append(nds, node)
		}
	} else if !msg.KeepAlive {
		id = msg.Id
		var nodes = VesselHistory.MarkDone(id)
		for _, n := range  nodes {
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

		//if not running simplify in background - start one
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
				time.Sleep(3 * time.Second)
			}
			s.OnlineDP.Simplify(id)
			fmt.Println("<< done >>")
		}()
	}

	return nil
}
