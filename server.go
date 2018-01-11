package main

import (
	"fmt"
	"log"
	"strings"
	"simplex/db"
	"database/sql"
	"path/filepath"
	_ "github.com/lib/pq"
	"simplex/streamdp/mtrafic"
	"gopkg.in/gin-gonic/gin.v1"
	"simplex/streamdp/onlinedp"
	"simplex/streamdp/config"
	"simplex/streamdp/common"
)

const (
	InputBufferSize = 3
)

type Server struct {
	Config       *config.Server
	Address      string
	Mode         int
	Src          *db.DataSrc
	ConstSrc     *db.DataSrc
	OnlineDP     *onlinedp.OnlineDP
	InputStream  chan []*db.Node
	SimpleStream chan []int
	Exit         chan struct{}
}

func NewServer(address string, mode int) *Server {
	var pwd = common.ExecutionDir()
	var exit = make(chan struct{})
	var inputStream = make(chan []*db.Node, InputBufferSize)
	var simpleStream = make(chan []int)

	var server = &Server{
		Address:      address,
		Mode:         mode,
		Config:       &config.Server{},
		InputStream:  inputStream,
		SimpleStream: simpleStream,
		Exit:         exit,
	}
	var fname = filepath.Join(pwd, "../resource/src.toml")
	server.Config.Load(fname)

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
	server.ConstSrc = db.NewDataSrc(filepath.Join(pwd, "../resource/consts.toml"))

	var dpOpts = server.Config.DPOptions()
	server.OnlineDP = onlinedp.NewOnlineDP(
		server.Src, server.ConstSrc, dpOpts,
		ScoreFn, true,
	)
	return server
}

func (server *Server) Run() {
	server.init()

	var router = gin.Default()
	if server.Mode == 0 {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router.POST("/ping", server.trafficRouter)
	router.POST("/history/clear", server.clearHistory)
	router.POST("/simplify", server.clearHistory)
	router.Run(server.Address)
}

func (server *Server) init() {
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

	//launch input stream processing
	go server.goProcessInputStream()
	go server.goProcessSimpleStream()
}

func (server *Server) clearHistory(ctx *gin.Context) {
	VesselHistory.Clear()
	ctx.JSON(Success, gin.H{"message": "success"})
}

//func (s *Server) simplify(ctx *gin.Context) {
//	s.OnlineDP.Simplify()
//	ctx.JSON(Success, gin.H{"message": "success"})
//}

func (server *Server) trafficRouter(ctx *gin.Context) {
	var msg = &mtrafic.PingMsg{}
	var err = ctx.BindJSON(msg)

	if err != nil {
		log.Panic(err)
		ctx.JSON(Error, gin.H{"message": "error"})
		return
	}

	err = server.aggregatePings(msg)
	if err == nil {
		ctx.JSON(Success, gin.H{"message": "success"})
	} else {
		ctx.JSON(Error, gin.H{"message": "error"})
	}
}
