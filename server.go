package main

import (
	"fmt"
	"log"
	"simplex/db"
	"database/sql"
	_ "github.com/lib/pq"
	"simplex/streamdp/data"
	"gopkg.in/gin-gonic/gin.v1"
)

func NewServer(address string, mode int) *Server {
	var server = (&Server{Address: address, Mode: mode}).loadConfig()
	var cfg = server.Config.DBConfig()
	var sqlsrc, err = sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Database,
	))
	if err != nil {
		panic(err)
	}

	server.Src = &db.DataSrc{
		Src:       sqlsrc,
		Config:    cfg,
		SRID:      server.Config.SRID,
		Dim:       server.Config.Dim,
		NodeTable: server.Config.Table,
	}
	return server
}

type Server struct {
	Config   *ServerConfig
	Address  string
	Mode     int
	Src      *db.DataSrc
	ConstSrc *db.DataSrc
}

func (s *Server) init() {
	//create online table
	if err := s.initCreateOnlineTable(); err != nil {
		log.Fatalln(err)
	}

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
	router.Run(s.Address)
}

func (s *Server) trafficRouter(ctx *gin.Context) {
	var msg = &data.PingMsg{}
	var err = ctx.BindJSON(msg)

	fmt.Println(">>> ping msg :", msg)

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

func (s *Server) loadConfig() *Server {
	var fileName = fmt.Sprintf("%v/src.toml", ExecutionDir())
	s.Config = &ServerConfig{}
	var err = s.Config.Load(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	return s
}
