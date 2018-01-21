package main

import (
	"log"
	"fmt"
	"flag"
	"time"
	"math/rand"
	"io/ioutil"
	"path/filepath"
	"github.com/naoina/toml"
	"simplex/streamdp/config"
	"simplex/streamdp/common"
	"simplex/streamdp/mtrafic"
)

var Port int
var Host string
var PingAddress string
var ClearHistoryAddress string
var UpdateServerCfg string

const concurProcs = 8

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.IntVar(&Port, "port", 8000, "host port")
	flag.StringVar(&Host, "host", "localhost", "host address")

	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	PingAddress = fmt.Sprintf("http://%v:%v/ping", Host, Port)
	ClearHistoryAddress = fmt.Sprintf("http://%v:%v/history/clear", Host, Port)
	UpdateServerCfg = fmt.Sprintf("http://%v:%v/update/server/config", Host, Port)

	var pwd = common.ExecutionDir()
	//var dataDir = filepath.Join(pwd, "../data")
	//var ignoreDirs = []string{".git", ".idea"}
	//var filter = []string{"toml"}

	var srcFile = filepath.Join(pwd, "../resource/src.toml")
	var constFile = filepath.Join(pwd, "../resource/consts.toml")

	var cfgMsg = mtrafic.CfgMsg{
		ConstraintToml: readTomlFile(constFile),
		ServerToml:     readTomlFile(srcFile),
	}

	for _, t := range []float64{1000, 2500} {
		var msg = cfgMsg.Clone()
		var cfg = config.ServerConfig{}
		var err = toml.Unmarshal([]byte(msg.ServerToml), &cfg)
		if err != nil {
			log.Panic(err)
		}
		cfg.Threshold = t
		stoml, err := toml.Marshal(cfg)
		msg.ServerToml = string(stoml)

		if err != nil {
			log.Panic(err)
		}

		res, err := post(UpdateServerCfg, []byte(msg.EncodeMsg().ToJSON()))
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(string(res))
		time.Sleep(5 * time.Second)
	}

	////clear history
	//runProcess(ClearHistoryAddress)
	//
	////vessel pings
	//vesselPings(dataDir, filter, ignoreDirs, concurProcs)
	////simplify
	////runProcess(SimplifyAddress)
}

func readTomlFile(fname string) string {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Panic(err)
	}
	return string(b)
}
