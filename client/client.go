package main

import (
	"log"
	"fmt"
	"flag"
	"time"
	"runtime"
	"math/rand"
	"simplex/db"
	"database/sql"
	"path/filepath"
	"simplex/streamdp/config"
	"simplex/streamdp/common"
)

var Port int
var Host string
var Address string
var ClearHistoryAddress string

const concurProcs = 8

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.IntVar(&Port, "port", 8000, "host port")
	flag.StringVar(&Host, "host", "localhost", "host address")
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var pwd = common.ExecutionDir()
	var dataDir = filepath.Join(pwd, "../data")
	var srcFile = filepath.Join(pwd, "../resource/src.toml")
	var ignoreDirs = []string{".git", ".idea"}
	var filter = []string{"toml"}

	ClearHistoryAddress = fmt.Sprintf("http://%v:%v/history/clear", Host, Port)
	Address = fmt.Sprintf("http://%v:%v/ping", Host, Port)

	var serverCfg = (&config.Server{}).Load(srcFile)
	var dbCfg = serverCfg.DBConfig()

	var sqlsrc, err = sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.User, dbCfg.Password, dbCfg.Database,
	))
	if err != nil {
		log.Panic(err)
	}
	var src = &db.DataSrc{
		Src:    sqlsrc,
		Config: dbCfg,
		SRID:   serverCfg.SRID,
		Dim:    serverCfg.Dim,
		Table:  serverCfg.Table,
	}
	db.CreateNodeTable(src)

	//clear history
	runProcess(ClearHistoryAddress)
	//vessel pings
	vesselPings(dataDir, filter, ignoreDirs, concurProcs)
	//simplify
	//runProcess(SimplifyAddress)
}
