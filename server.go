package main

import (
	"fmt"
	"log"
	"strings"
	"simplex/db"
	"database/sql"
	_ "github.com/lib/pq"
	"simplex/streamdp/data"
	"gopkg.in/gin-gonic/gin.v1"
	"simplex/streamdp/onlinedp"
	"simplex/streamdp/offset"
	"path/filepath"
)

func NewServer(address string, mode int) *Server {
	var pwd = ExecutionDir()
	var server = &Server{Address: address, Mode: mode, Config: &ServerConfig{}}
	if err := server.Config.Load(filepath.Join(pwd, "src.toml")); err != nil {
		log.Panic(err)
	}

	var cfg = server.Config.DBConfig()
	inputSrc, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Database,
	))

	if err != nil {
		log.Panic(err)
	}

	server.Src = &db.DataSrc{
		Src:       inputSrc,
		Config:    cfg,
		SRID:      server.Config.SRID,
		Dim:       server.Config.Dim,
		NodeTable: server.Config.Table,
	}
	server.ConstSrc = db.NewDataSrc(
		filepath.Join(pwd, "consts.toml"),
	)

	return server
}

type Server struct {
	Config   *ServerConfig
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
		s.Src, s.ConstSrc, dpOpts,
		offset.MaxOffset, true,
	)


	//create online table
	if err := s.initCreateOnlineTable(); err != nil {
		log.Fatalln(err)
	}
}

func (s *Server) clearHistory(ctx *gin.Context) {
	VesselHistory.Clear()
	ctx.JSON(Success, gin.H{"message": "success"})
}

func (s *Server) simplify(ctx *gin.Context) {
	s.OnlineDP.Simplify()
	s.SaveSimplification()
	ctx.JSON(Success, gin.H{"message": "success"})
}

func (s *Server) trafficRouter(ctx *gin.Context) {
	var msg = &data.PingMsg{}
	var err = ctx.BindJSON(msg)

	if err != nil {
		panic(err)
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
