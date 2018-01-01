package onlinedp

import (
	"fmt"
	"strings"
)

func (self *OnlineDP) tempNodeIDTableName(fid int) string {
	return fmt.Sprintf("temp_%v_%v", self.Src.NodeTable, fid)
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

func (self *OnlineDP) tempCreateSnapshotTable(temp string) {
	var query = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
		    id      INT NOT NULL,
		    size    INT,
			gob     TEXT NOT NULL,
		    status  INT DEFAULT 0,
		    CONSTRAINT pid_%v PRIMARY KEY (id)
 		) WITH (OIDS=FALSE);
		CREATE INDEX idx_size_%v ON %v (size);
		CREATE INDEX idx_status_%v ON %v (status);
`, temp, temp, temp, temp, temp, temp,
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
		panic(err)
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
	return fmt.Sprintf("tempQ_%v", self.Src.NodeTable)
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
		qs = append(qs, fmt.Sprintf(`('%v')`, encode64(q)))
	}
	var vals = strings.Join(qs, ",")
	var query = fmt.Sprintf(
		`INSERT INTO %v (query) VALUES %v;`,
		temp, vals,
	)
	var _, err = self.Src.Exec(query)
	if err != nil {
		panic(err)
	}
}

func (self *OnlineDP) tempExecuteQueries(tempQ string) {
	var query = fmt.Sprintf("SELECT query  FROM  %v;", tempQ)
	var h, err = self.Src.Query(query)
	if err != nil {
		panic(err)
	}

	const bufferSize = 100
	var q string
	var buf = make([]string, 0)
	for h.Next() {
		h.Scan(&q)
		buf = append(buf, decode64(q))
		if len(buf) > bufferSize {
			self.ExecuteTransaction(buf)
			buf = make([]string, 0)
		}
	}
	//flush
	if len(buf) > 0 {
		self.ExecuteTransaction(buf)
	}
}
