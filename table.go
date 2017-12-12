package main

import (
	"log"
	"fmt"
	"bytes"
	"text/template"
)

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//TODO: dataset restarts beginning on start - change when real online
var onlineTblTemplate = `
CREATE TABLE IF NOT EXISTS {{.Table}} (
    id  SERIAL NOT NULL,
    i INT NOT NULL,
    j INT NOT NULL,
    size INT CHECK (size > 0),
    fid INT NOT NULL,
    part INT NOT NULL,
    gob TEXT NOT NULL,
    geom GEOMETRY(Geometry, {{.SRID}}) NOT NULL,
    status INT DEFAULT 0,
    CONSTRAINT pid_{{.Table}} PRIMARY KEY (id)
) WITH (OIDS=FALSE);
CREATE INDEX idx_i_{{.Table}} ON {{.Table}} (i);
CREATE INDEX idx_j_{{.Table}} ON {{.Table}} (j);
CREATE INDEX idx_size_{{.Table}} ON {{.Table}} (size);
CREATE INDEX idx_fid_{{.Table}} ON {{.Table}} (fid);
CREATE INDEX idx_part_{{.Table}} ON {{.Table}} (part);
CREATE INDEX idx_status_{{.Table}} ON {{.Table}} (status);
CREATE INDEX gidx_{{.Table}} ON {{.Table}} USING GIST (geom);
`

var onlineTemplate *template.Template

func init() {
	var err error
	onlineTemplate, err = template.New("online_table").Parse(onlineTblTemplate)
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *Server) initCreateOnlineTable() error {
	var query bytes.Buffer
	var err = onlineTemplate.Execute(&query, s.Config)
	if err != nil {
		log.Fatalln(err)
	}
	var tblSQl = fmt.Sprintf(`DROP TABLE IF EXISTS %v CASCADE;`, s.Config.Table)
	_, err = s.Src.Exec(tblSQl)
	if err != nil {
		log.Panic(err)
	}
	_, err = s.Src.Exec(query.String())
	return err
}
