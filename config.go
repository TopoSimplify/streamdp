package main

import (
	"simplex/db"
	"simplex/opts"
	"github.com/naoina/toml"
	"github.com/intdxdt/fileutil"
)

type ServerConfig struct {
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

func (cfg *ServerConfig) DBConfig() db.Config {
	return db.Config{
		Host:           cfg.DBHost,
		Password:       cfg.Password,
		Database:       cfg.Database,
		User:           cfg.User,
		Table:          cfg.Table,
		GeometryColumn: GeomColumn,
		IdColumn:       IdColumn,
	}
}

func (cfg *ServerConfig) DPOptions() *opts.Opts {
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

func (cfg *ServerConfig) Load(fileName string) error {
	var txt, err = fileutil.ReadAllOfFile(fileName)
	if err != nil {
		return err
	}
	return toml.Unmarshal([]byte(txt), cfg)
}
