package main

import (
	"log"
	"fmt"
	"sync"
	"time"
	"strings"
	"spinner"
	"simplex/db"
	"simplex/streamdp/common"
)

const (
	NullState      = 0
	batchSize      = 8
	simpleInterval = 1 //secs
	simpleIdTable  = "temp_simple_ids"
	IdleLimit      = 3
)

func (server *Server) goProcessInputStream() {
	for {
		select {
		//listen to results channel
		case <-server.Exit:
			server.ExitWg.Done()
			return
		case nodes := <-server.InputStream:
			if len(nodes) == 0 {
				continue
			}
			//var insertSQL = nds[0].InsertSQL(server.Config.Table, server.Config.SRID, nds...)
			var vals = common.SnapshotNodeColumnValues(
				server.Src.SRID,
				common.UnSnap,
				nodes...
			)
			var insertSQL = db.SQLInsertIntoTable(
				server.Src.Table,
				common.NodeColumnFields,
				vals,
			)
			if _, err := server.Src.Exec(insertSQL); err != nil {
				log.Panic(err)
			}
		}
	}
}

func (server *Server) goProcessSimpleStream() {
	var idleCount = 0
	server.dropSimpleIdTable() //drop if exist
	//listen to results channel
	for {
		select {
		case <-server.Exit:
			server.ExitWg.Done()
			return
		case <-time.After(simpleInterval * time.Second):
			server.createSimpleIdTable()
			server.copyIdsIntoSimpleIdTable()

			var buf = make([]int, 0)

			var simplify = func(fid int, wg *sync.WaitGroup) {
				server.OnlineDP.Simplify(fid)
				wg.Done()
			}

			var flush = func() {
				var n = len(buf)
				var exit = make(chan struct{})
				defer close(exit)

				var msg = fmt.Sprintf("processing ... %v", n)
				spinner.NewSpinner(msg, exit).Start()

				var wg = &sync.WaitGroup{}
				wg.Add(n)
				for _, fid := range buf {
					go simplify(fid, wg)
				}
				buf = make([]int, 0)
				wg.Wait()
			}

			var bln bool

			var query = fmt.Sprintf("SELECT fid FROM %v;", simpleIdTable)
			var h, err = server.Src.Query(query)
			if err != nil {
				log.Panic(err)
			}

			for h.Next() {
				idleCount = 0
				bln = true
				var fid int
				h.Scan(&fid)
				buf = append(buf, fid)
				if len(buf) >= batchSize {
					flush()
				}
			}
			//close handler
			h.Close()

			//flush
			if len(buf) > 0 {
				flush()
			}

			if !bln {
				idleCount += 1
				if idleCount > IdleLimit {
					server.TaskMap[server.CurTaskID] = Done
				}
				log.Println("...nothing to simplify...")
			}
			//drop table
			server.dropSimpleIdTable()
		}
	}
}

func (server *Server) createSimpleIdTable() {
	var query = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
		    fid  INT NOT NULL,
		    CONSTRAINT idx_%v PRIMARY KEY (fid)
		) WITH (OIDS=FALSE);`,
		simpleIdTable, simpleIdTable,
	)
	var _, err = server.Src.Exec(query)
	if err != nil {
		log.Panic(err)
	}
}

func (server *Server) dropSimpleIdTable() {
	var query = fmt.Sprintf(`
		DROP TABLE IF EXISTS %v CASCADE;
	`, simpleIdTable)
	if _, err := server.Src.Exec(query); err != nil {
		log.Panic(err)
	}
}

func (server *Server) copyIdsIntoSimpleIdTable() {
	var query = fmt.Sprintf(`
		SELECT DISTINCT fid  FROM %v WHERE status=%v;`,
		server.Src.Table, NullState,
	)

	var h, err = server.Src.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer h.Close()

	var bufferSize = 100
	var ids = make([]int, 0)
	for h.Next() {
		var id int
		h.Scan(&id)
		ids = append(ids, id)
		if len(ids) > bufferSize {
			server.insertInSimpleIdTable(ids)
			ids = make([]int, 0)
		}
	}
	if len(ids) > 0 {
		server.insertInSimpleIdTable(ids)
	}
}

func (server *Server) insertInSimpleIdTable(fids []int) {
	var qs = make([]string, 0)
	for _, id := range fids {
		qs = append(qs, fmt.Sprintf(`(%v)`, id))
	}
	var vals = strings.Join(qs, ",")
	var query = fmt.Sprintf(`
		INSERT INTO %v (fid) VALUES %v;
		`, simpleIdTable, vals,
	)
	if _, err := server.Src.Exec(query); err != nil {
		log.Panic(err)
	}
}
