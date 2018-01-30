package onlinedp

import (
	"log"
	"simplex/db"
	"github.com/naoina/toml"
	"simplex/streamdp/config"
	"github.com/intdxdt/fileutil"
)

func loadConfig(filename string) *ServerConfig {
	var cfg = &ServerConfig{ServerConfig: config.ServerConfig{}}
	if err := cfg.Load(filename); err != nil {
		log.Panic(err)
	}
	cfg.Database = TestDBName
	cfg.Table = TestTable
	return cfg
}

type ServerConfig struct {
	config.ServerConfig
}

func (cfg *ServerConfig) DBConfig() db.Config {
	return db.Config{
		Host:           cfg.DBHost,
		Password:       cfg.Password,
		Database:       cfg.Database,
		User:           cfg.User,
		Table:          cfg.Table,
		GeometryColumn: "geom",
		IdColumn:       "id",
	}
}

func (cfg *ServerConfig) Load(fileName string) error {
	var txt, err = fileutil.ReadAllOfFile(fileName)
	if err != nil {
		return err
	}

	return toml.Unmarshal([]byte(txt), &cfg.ServerConfig)
}
