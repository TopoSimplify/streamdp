package onlinedp

import (
	"fmt"
	"log"
	"bytes"
	"runtime"
	"simplex/db"
	"text/template"
	"simplex/streamdp/pt"
	"github.com/intdxdt/fan"
	"github.com/intdxdt/geom"
)

var onlineOutputTblTemplate = `
DROP TABLE IF EXISTS {{.NodeTable}} CASCADE;
CREATE TABLE IF NOT EXISTS {{.NodeTable}} (
    id  INT NOT NULL,
    geom GEOMETRY(Geometry, {{.SRID}}) NOT NULL,
    CONSTRAINT pid_{{.NodeTable}} PRIMARY KEY (id)
) WITH (OIDS=FALSE);`

var onlineOutputTemplate *template.Template

func init() {
	var err error
	onlineOutputTemplate, err = template.New("online_output_table").Parse(onlineOutputTblTemplate)
	if err != nil {
		log.Panic(err)
	}
}

//obj.createOutputOnlineTable()
//obj.SaveSimplification()

//Find and merge simple segments
func (self *OnlineDP) SaveSimplification() {
	var stream = make(chan interface{})
	var exit = make(chan struct{})
	defer close(exit)

	var query bytes.Buffer
	self.Src.NodeTable = fmt.Sprintf(`%v_simple`, self.Src.Config.Table)
	if err := onlineOutputTemplate.Execute(&query, self.Src); err != nil {
		log.Panic(err)
	}

	if _, err := self.Src.Src.Exec(query.String()); err != nil {
		log.Panic(err)
	}

	var outputTable = self.Src.NodeTable
	//o.Src.DuplicateTable(outputTable)
	self.Src.AlterAsMultiLineString(
		outputTable, self.Src.Config.GeometryColumn, self.Src.SRID,
	)

	go func() {
		var query = fmt.Sprintf(
			`SELECT DISTINCT %v FROM %v;`, "fid", self.Src.Config.Table,
		)
		var h, err = self.Src.Query(query)
		if err != nil {
			log.Fatalln(err)
		}

		var id int
		for h.Next() {
			h.Scan(&id)
			stream <- id
		}
		close(stream)
	}()

	var worker = func(v interface{}) interface{} {
		var id = v.(int)
		//aggregate src into linear fid and parts
		self.aggregateNodes(id, outputTable)
		return true
	}

	var out = fan.Stream(stream, worker, runtime.NumCPU(), exit)
	for range out {
	}
}

func (self *OnlineDP) aggregateNodes(id int, outputTable string) {
	var query = fmt.Sprintf(`
		SELECT fid, part, gob FROM %v WHERE fid=%v ORDER BY fid asc, part asc, i asc;`,
		self.Src.Config.Table, id,
	)
	var h, err = self.Src.Query(query)
	if err != nil {
		panic(err)
	}

	var gob string
	var fid, part int
	var coordinates = make([][]*pt.Pt, 0)

	var idx = -1
	var curPart = -1
	for h.Next() {
		h.Scan(&fid, &part, &gob)
		var o = db.Deserialize(gob)
		var i, j = 0, len(o.Coordinates)-1

		if idx == -1 {
			curPart = part
			coordinates = append(coordinates, []*pt.Pt{})
		}

		idx = len(coordinates) - 1
		if curPart == part {
			var last *pt.Pt
			if len(coordinates[idx]) > 0 {
				n := len(coordinates[idx]) - 1
				last = coordinates[idx][n]
			}
			if last == nil {
				coordinates[idx] = append(coordinates[idx],
					&pt.Pt{Point: o.Coordinates[i], I: o.Range.I},
					&pt.Pt{Point: o.Coordinates[j], I: o.Range.J},
				)
			} else if last.I == o.Range.I {
				coordinates[idx] = append(coordinates[idx],
					&pt.Pt{Point: o.Coordinates[j], I: o.Range.J},
				)
			} else {
				fmt.Println(query)
				panic("coordinates non contiguous")
			}
		} else {
			curPart = part
			coordinates = append(coordinates, []*pt.Pt{
				{Point: o.Coordinates[i], I: o.Range.I},
				{Point: o.Coordinates[j], I: o.Range.J},
			})
		}
	}

	var buf bytes.Buffer
	buf.WriteString("MULTILINESTRING (")
	var n = len(coordinates) - 1
	for i, coords := range coordinates {
		var sub = make([]*geom.Point, len(coords))
		for idx, o := range coords {
			sub[idx] = o.Point
		}

		buf.WriteString(wktLineString(sub, self.Src.Dim))
		if i < n {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")

	var wkt = buf.String()

	//fmt.Println(wkt)
	var geomFromTxt = fmt.Sprintf(
		`st_geomfromtext('%v', %v)`,
		wkt, self.Src.SRID,
	)
	query = fmt.Sprintf(
		`INSERT INTO %v (id, geom) VALUES (%v, %v);`,
		outputTable, id, geomFromTxt,
	)

	_, err = self.Src.Exec(query)
	if err != nil {
		panic(err)
	}
}
