package main

import (
	"fmt"
	"log"
	"sync"
	"time"
	"strings"
	"runtime"
	"simplex/db"
	"database/sql"
	"simplex/streamdp/common"
	"simplex/streamdp/onlinedp"
	"simplex/streamdp/mtrafic"
)

func (server *Server) initSources() {
	if server.ConstSrc != nil {
		server.ConstSrc.Close()
	}
	if server.Src != nil {
		server.Src.Close()
	}
	var dbCfg = server.Config.DBConfig()
	inputSrc, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.User, dbCfg.Password, dbCfg.Database,
	))

	if err != nil {
		log.Panic(err)
	}

	server.Src = &db.DataSrc{
		Src:    inputSrc,
		Config: dbCfg,
		SRID:   server.Config.SRID,
		Dim:    server.Config.Dim,
		Table:  server.Config.Table,
	}

	var dpOpts = server.Config.DPOptions()
	server.OnlineDP = onlinedp.NewOnlineDP(
		server.Src,
		server.ConstSrc,
		dpOpts,
		ScoreFn,
		true,
	)

}

func (server *Server) closeSources() {
	//close sources
	if server.Src != nil {
		server.Src.Close()
	}
	if server.ConstSrc != nil {
		server.ConstSrc.Close()
	}
}

func (server *Server) loadConfig(msg *mtrafic.CfgMsg) {
	fmt.Println("threshold :", server.Config.Threshold)
	server.Config.Load(msg.ServerToml)
	server.ConstSrc = db.NewDataSrc(msg.ConstraintToml)
}

func (server *Server) closeStreams() {
	close(server.InputStream)
	close(server.SimpleStream)
}

func (server *Server) initExit() {
	close(server.Exit)
	server.ExitWg.Wait()
}

func (server *Server) init(msg *mtrafic.CfgMsg) {
	server.initExit()
	server.closeSources()
	server.loadConfig(msg)
	server.closeStreams()

	runtime.Gosched()
	log.Println("setting up server config ...")
	time.Sleep(3 * time.Second)

	server.initSources()

	server.Exit = make(chan struct{})
	server.InputStream = make(chan []*db.Node, InputBufferSize)
	server.SimpleStream = make(chan []int)

	var simpleType = strings.ToLower(server.Config.SimplficationType)

	if simpleType == "nopw" {
		SimplificationType = NOPW
	} else if simpleType == "bopw" {
		SimplificationType = BOPW
	} else {
		log.Panic("unknown simplification type: NOPW or BOPW")
	}

	//create online table
	if err := db.CreateNodeTable(server.Src); err != nil {
		log.Panic(err)
	}
	var simpleTable = common.SimpleTable(server.Src.Table)

	var query = fmt.Sprintf(`
		DROP TABLE IF EXISTS %v CASCADE;
		CREATE TABLE IF NOT EXISTS %v (
		    id          INT NOT NULL,
		    geom        GEOMETRY(Geometry, %v) NOT NULL,
			count       INT NOT NULL,
		    CONSTRAINT  pid_%v PRIMARY KEY (id)
		) WITH (OIDS=FALSE);`,
		simpleTable,
		simpleTable,
		server.Src.SRID,
		simpleTable,
	)

	if _, err := server.Src.Exec(query); err != nil {
		log.Panic(err)
	}

	//o.Src.DuplicateTable(outputTable)
	server.Src.AlterAsMultiLineString(
		simpleTable, server.Src.Config.GeometryColumn, server.Src.SRID,
	)

	server.ExitWg = &sync.WaitGroup{}
	server.ExitWg.Add(2)

	//launch input stream processing
	go server.goProcessInputStream()
	go server.goProcessSimpleStream()
}
