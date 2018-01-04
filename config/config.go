package config

import (
	"log"
	"simplex/db"
	"simplex/opts"
	"github.com/naoina/toml"
	"github.com/intdxdt/fileutil"
)

type Server struct {
	ServerAddress          string  `toml:"ServerAddress"`
	DBHost                 string  `toml:"DBHost"`
	Password               string  `toml:"Password"`
	Database               string  `toml:"Database"`
	User                   string  `toml:"User"`
	Table                  string  `toml:"Table"`
	SimplficationType      string  `toml:"SimplficationType"`
	SRID                   int     `toml:"SRID"`
	Dim                    int     `toml:"Dim"`
	Threshold              float64 `toml:"Threshold"`
	MinDist                float64 `toml:"MinDist"`
	RelaxDist              float64 `toml:"RelaxDist"`
	AvoidNewSelfIntersects bool    `toml:"AvoidNewSelfIntersects"`
	GeomRelation           bool    `toml:"GeomRelation"`
	DistRelation           bool    `toml:"DistRelation"`
	DirRelation            bool    `toml:"DirRelation"`
}

func (cfg *Server) DBConfig() db.Config {
	return db.Config{
		Host:           cfg.DBHost,
		Password:       cfg.Password,
		Database:       cfg.Database,
		User:           cfg.User,
		Table:          cfg.Table,
		GeometryColumn: db.GeomColumn,
		IdColumn:       db.IdColumn,
	}
}

func (cfg *Server) DPOptions() *opts.Opts {
	return &opts.Opts{
		Threshold:              cfg.Threshold,
		MinDist:                cfg.MinDist,
		RelaxDist:              cfg.RelaxDist,
		KeepSelfIntersects:     false,
		AvoidNewSelfIntersects: cfg.AvoidNewSelfIntersects,
		GeomRelation:           cfg.GeomRelation,
		DistRelation:           cfg.DistRelation,
		DirRelation:            cfg.DirRelation,
	}
}

func (cfg *Server) Load(fileName string) *Server {
	if txt, err := fileutil.ReadAllOfFile(fileName); err == nil {
		if err = toml.Unmarshal([]byte(txt), cfg); err != nil {
			log.Panic(err)
		}
	} else {
		log.Panic(err)
	}
	return cfg
}
