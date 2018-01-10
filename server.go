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
	"simplex/streamdp/offset"
	"simplex/streamdp/config"
	"simplex/streamdp/common"
)

func NewServer(address string, mode int) *Server {
	var pwd = common.ExecutionDir()
	var server = &Server{
		Address: address,
		Mode:    mode,
		Config:  &config.Server{},
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
	return server
}

type Server struct {
	Config   *config.Server
	Address  string
	Mode     int
	Src      *db.DataSrc
	ConstSrc *db.DataSrc
	OnlineDP *onlinedp.OnlineDP
}

func (s *Server) Run() {
	s.init()

	var router = gin.Default()
	if s.Mode == 0 {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router.POST("/ping", s.trafficRouter)
	router.POST("/history/clear", s.clearHistory)
	router.POST("/simplify", s.clearHistory)
	router.Run(s.Address)
}

func (s *Server) init() {
	var simpleType = strings.ToLower(s.Config.SimplficationType)

	if simpleType == "nopw" {
		SimplificationType = NOPW
	} else if simpleType == "bopw" {
		SimplificationType = BOPW
	}

	var dpOpts = s.Config.DPOptions()
	s.OnlineDP = onlinedp.NewOnlineDP(
		s.Src, s.ConstSrc, dpOpts, offset.MaxOffset,
		true,
	)

	//create online table
	if err := db.CreateNodeTable(s.Src); err != nil {
		log.Panic(err)
	}

	var simpleTable = common.SimpleTable(s.Src.Table)

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
		s.Src.SRID,
		simpleTable,
	)

	if _, err := s.Src.Exec(query); err != nil {
		log.Panic(err)
	}
	//o.Src.DuplicateTable(outputTable)
	s.Src.AlterAsMultiLineString(
		simpleTable, s.Src.Config.GeometryColumn, s.Src.SRID,
	)
}

func (s *Server) clearHistory(ctx *gin.Context) {
	VesselHistory.Clear()
	ctx.JSON(Success, gin.H{"message": "success"})
}

//func (s *Server) simplify(ctx *gin.Context) {
//	s.OnlineDP.Simplify()
//	ctx.JSON(Success, gin.H{"message": "success"})
//}

func (s *Server) trafficRouter(ctx *gin.Context) {
	var msg = &mtrafic.PingMsg{}
	var err = ctx.BindJSON(msg)

	if err != nil {
		log.Panic(err)
		ctx.JSON(Error, gin.H{"message": "error"})
		return
	}

	err = s.aggregatePings(msg)
	if err == nil {
		ctx.JSON(Success, gin.H{"message": "success"})
	} else {
		ctx.JSON(Error, gin.H{"message": "error"})
	}
}
