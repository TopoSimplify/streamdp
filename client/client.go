package main

import (
	"log"
	"fmt"
	"flag"
	"time"
	"runtime"
	"math/rand"
	"database/sql"
	"path/filepath"
	"github.com/TopoSimplify/db"
	"github.com/TopoSimplify/streamdp/config"
	"github.com/TopoSimplify/streamdp/common"
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
	var pwd         = common.ExecutionDir()
	var dataDir     = filepath.Join(pwd, "../data")
	var srcFile     = filepath.Join(pwd, "../resource/src.toml")
	var ignoreDirs  = []string{".git", ".idea"}
	var filter      = []string{"toml"}

	ClearHistoryAddress = fmt.Sprintf("http://%v:%v/history/clear", Host, Port)
	Address = fmt.Sprintf("http://%v:%v/ping", Host, Port)

	var serverCfg = (&config.ServerConfig{}).Load(srcFile)
	var dbCfg = serverCfg.DBConfig()
	var connSettings = fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.User, dbCfg.Password, dbCfg.Database,
	)

	var sqlsrc, err = sql.Open("postgres", connSettings)
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
