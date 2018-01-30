package onlinedp

import (
	"fmt"
	"strings"
	"log"
	"simplex/streamdp/enc"
)

func (self *OnlineDP) tempNodeIDTableName(fid int) string {
	return fmt.Sprintf("temp_%v_%v", self.Src.Table, fid)
}

func (self *OnlineDP) tempCreateNodeIdTable(temp string) {
	var query = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
		    id  INT NOT NULL,
		    CONSTRAINT pid_%v PRIMARY KEY (id)
		) WITH (OIDS=FALSE);`, temp, temp,
	)
	var _, err = self.Src.Exec(query)
	if err != nil {
		panic(err)
	}
}

func (self *OnlineDP) tempInsertInNodeIdTable(temp string, nid int) {
	var query = fmt.Sprintf(
		`INSERT INTO %v (id) VALUES (%v) ON CONFLICT (id) DO NOTHING;`,
		temp, nid,
	)
	var _, err = self.Src.Exec(query)
	if err != nil {
		log.Panic(err)
	}
}

func (self *OnlineDP) tempDropTable(temp string) {
	var query = fmt.Sprintf("DROP TABLE IF EXISTS %v CASCADE;", temp)
	var _, err = self.Src.Exec(query)
	if err != nil {
		panic(err)
	}
}

//=============================================================

func (self *OnlineDP) tempQueryTableName() string {
	return fmt.Sprintf("tempQ_%v", self.Src.Table)
}

func (self *OnlineDP) tempCreateTempQueryTable(temp string) {
	var query = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
		    id      SERIAL NOT NULL,
		    query   text NOT NULL,
		    CONSTRAINT pid_%v PRIMARY KEY (id)
		) WITH (OIDS=FALSE);`, temp, temp,
	)
	var _, err = self.Src.Exec(query)
	if err != nil {
		panic(err)
	}
}

func (self *OnlineDP) tempInsertInTOTempQueryTable(temp string, queries []string) {
	if len(queries) == 0 {
		return
	}
	var qs = make([]string, 0)
	for _, q := range queries {
		qs = append(qs, fmt.Sprintf(`('%v')`, enc.Encode64(q)))
	}
	var vals = strings.Join(qs, ",")
	var query = fmt.Sprintf(
		`INSERT INTO %v (query) VALUES %v;`,
		temp, vals,
	)
	var _, err = self.Src.Exec(query)
	if err != nil {
		log.Panic(err)
	}
}

