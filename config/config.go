package config

import (
	"log"
	"simplex/db"
	"simplex/opts"
	"github.com/naoina/toml"
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

func (scfg *ServerConfig) DBConfig() db.Config {
	return db.Config{
		Host:           scfg.DBHost,
		Password:       scfg.Password,
		Database:       scfg.Database,
		User:           scfg.User,
		Table:          scfg.Table,
		GeometryColumn: db.GeomColumn,
		IdColumn:       db.IdColumn,
	}
}

func (scfg *ServerConfig) DPOptions() *opts.Opts {
	return &opts.Opts{
		Threshold:              scfg.Threshold,
		MinDist:                scfg.MinDist,
		RelaxDist:              scfg.RelaxDist,
		KeepSelfIntersects:     false,
		AvoidNewSelfIntersects: scfg.AvoidNewSelfIntersects,
		GeomRelation:           scfg.GeomRelation,
		DistRelation:           scfg.DistRelation,
		DirRelation:            scfg.DirRelation,
	}
}

func (scfg *ServerConfig) Clone() *ServerConfig {
	return scfg.cloneCfg()
}

func (scfg ServerConfig) cloneCfg() *ServerConfig {
	return &scfg
}

func (scfg *ServerConfig) Load(tomlstring string) *ServerConfig {
	if err := toml.Unmarshal([]byte(tomlstring), scfg); err != nil {
		log.Panic(err)
	}
	return scfg
}
