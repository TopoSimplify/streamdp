package main

import (
	"simplex/db"
	_ "github.com/lib/pq"
	"simplex/streamdp/config"
	"simplex/streamdp/onlinedp"
	"sync"
)

const (
	InputBufferSize = 3
	Busy            = "busy"
	Done            = "done"
)

type Server struct {
	Config       *config.ServerConfig
	Address      string
	Mode         int
	Src          *db.DataSrc
	ConstSrc     *db.DataSrc
	OnlineDP     *onlinedp.OnlineDP
	InputStream  chan []*db.Node
	SimpleStream chan []int
	Exit         chan struct{}
	ExitWg       *sync.WaitGroup
	TaskMap      map[string]string
	CurTaskID    string
}

func NewServer(address string, mode int) *Server {
	var exit = make(chan struct{})
	var inputStream = make(chan []*db.Node, InputBufferSize)
	var simpleStream = make(chan []int)
	var exitWg = &sync.WaitGroup{}
	exitWg.Add(0)

	var server = &Server{
		Address:      address,
		Mode:         mode,
		Config:       &config.ServerConfig{},
		InputStream:  inputStream,
		SimpleStream: simpleStream,
		Exit:         exit,
		ExitWg:       exitWg,
		TaskMap:      make(map[string]string),
	}

	return server
}
