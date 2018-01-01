package main

import (
	"bytes"
	"log"
	"fmt"
	"simplex/db"
	"text/template"
)

var onlineTblTemplate = `
CREATE TABLE IF NOT EXISTS {{.Table}} (
    id  SERIAL NOT NULL,
    i INT NOT NULL,
    j INT NOT NULL,
    size INT CHECK (size > 0),
    fid INT NOT NULL,
    gob TEXT NOT NULL,
    geom GEOMETRY(Geometry, {{.SRID}}) NOT NULL,
    status INT DEFAULT 0,
    CONSTRAINT pid_{{.Table}} PRIMARY KEY (id)
) WITH (OIDS=FALSE);
CREATE INDEX idx_i_{{.Table}} ON {{.Table}} (i);
CREATE INDEX idx_j_{{.Table}} ON {{.Table}} (j);
CREATE INDEX idx_size_{{.Table}} ON {{.Table}} (size);
CREATE INDEX idx_fid_{{.Table}} ON {{.Table}} (fid);
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

func initCreateOnlineTable(src *db.DataSrc, cfg *ServerConfig) error {
	var query bytes.Buffer
	var err = onlineTemplate.Execute(&query, cfg)
	if err != nil {
		log.Fatalln(err)
	}
	var tblSQl = fmt.Sprintf(`DROP TABLE IF EXISTS %v CASCADE;`, cfg.Table)
	_, err = src.Exec(tblSQl)
	if err != nil {
		log.Panic(err)
	}
	_, err = src.Exec(query.String())
	return err
}

