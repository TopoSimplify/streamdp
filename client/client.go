package main

import (
	"fmt"
	"flag"
	"time"
	"runtime"
	"math/rand"
	"simplex/db"
	"database/sql"
)

var Port int
var Host string
var Address string
var ClearHistoryAddress string
var SimplifyAddress string

const concurProcs = 8
const GeomColumn = "geom"
const IdColumn = "id"

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.IntVar(&Port, "port", 8000, "host port")
	flag.StringVar(&Host, "host", "localhost", "host address")
	rand.Seed(time.Now().UTC().UnixNano())
}

const tomlpath = "/home/titus/01/godev/src/simplex/streamdp/sandbox/src.toml"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var msisDir = "/home/titus/01/godev/src/simplex/streamdp/mmsis"
	var ignoreDirs = []string{".git", ".idea"}
	var filter = []string{"toml"}

	ClearHistoryAddress = fmt.Sprintf("http://%v:%v/history/clear", Host, Port)
	SimplifyAddress = fmt.Sprintf("http://%v:%v/simplify", Host, Port)
	Address = fmt.Sprintf("http://%v:%v/ping", Host, Port)

	var serverCfg = loadConfig(tomlpath)
	var cfg = serverCfg.DBConfig()

	var sqlsrc, err = sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Database,
	))
	if err != nil {
		panic(err)
	}
	var src = &db.DataSrc{
		Src:       sqlsrc,
		Config:    cfg,
		SRID:      serverCfg.SRID,
		Dim:       serverCfg.Dim,
		NodeTable: serverCfg.Table,
	}
	initCreateOnlineTable(src, serverCfg)

	//clear history
	runProcess(ClearHistoryAddress)
	//vessel pings
	vesselPings(msisDir, filter, ignoreDirs, concurProcs)
	//simplify
	//runProcess(SimplifyAddress)
}
