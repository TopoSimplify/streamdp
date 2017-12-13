package onlinedp

import (
	"fmt"
	"log"
	"bytes"
	"simplex/db"
	"github.com/intdxdt/fan"
	"github.com/intdxdt/geom"
)

type Pt struct {
	pt *geom.Point
	i  int
}

//Find and merge simple segments
func (self *OnlineDP) SaveSimplification() {
	var stream = make(chan interface{}, 4*concurProcs)
	var exit = make(chan struct{})
	defer close(exit)

	var outputTable = self.Src.Config.Table + "_simple"
	self.Src.DuplicateTable(outputTable)
	self.Src.AlterAsMultiLineString(
		outputTable, self.Src.Config.GeometryColumn, self.Src.SRID,
	)

	go func() {
		var query = fmt.Sprintf(
			`SELECT %v FROM %v;`,
			self.Src.Config.IdColumn, self.Src.Config.Table,
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

	var out = fan.Stream(stream, worker, concurProcs, exit)
	for range out {
	}
}

func (self *OnlineDP) aggregateNodes(id int, outputTable string) {
	var query = fmt.Sprintf(`
		SELECT fid, part, gob
		FROM %v WHERE fid=%v ORDER BY fid asc, part asc, i asc;`,
		self.Src.NodeTable, id,
	)
	var h, err = self.Src.Query(query)
	if err != nil {
		panic(err)
	}

	var gob string
	var fid, part int
	var coordinates = make([][]*Pt, 0)

	var idx = -1
	var curPart = -1
	for h.Next() {
		h.Scan(&fid, &part, &gob)
		var o = db.Deserialize(gob)
		var i, j = 0, len(o.Coordinates)-1

		if idx == -1 {
			curPart = part
			coordinates = append(coordinates, []*Pt{})
		}

		idx = len(coordinates) - 1
		if curPart == part {
			var last *Pt
			if len(coordinates[idx]) > 0 {
				n := len(coordinates[idx]) - 1
				last = coordinates[idx][n]
			}
			if last == nil {
				coordinates[idx] = append(coordinates[idx],
					&Pt{pt: o.Coordinates[i], i: o.Range.I},
					&Pt{pt: o.Coordinates[j], i: o.Range.J},
				)
			} else if last.i == o.Range.I {
				coordinates[idx] = append(coordinates[idx],
					&Pt{pt: o.Coordinates[j], i: o.Range.J},
				)
			} else {
				panic("coordinates non contiguous")
			}

		} else {
			curPart = part
			coordinates = append(coordinates, []*Pt{
					{pt: o.Coordinates[i], i: o.Range.I},
					{pt: o.Coordinates[j], i: o.Range.J},
				})
		}
	}

	var buf bytes.Buffer
	buf.WriteString("MULTILINESTRING (")
	var n = len(coordinates) - 1
	for i, coords := range coordinates {
		var sub = make([]*geom.Point, len(coords))
		for idx, o := range coords {
			sub[idx] = o.pt
		}

		buf.WriteString(wktLineString(sub, self.Src.Dim))
		if i < n {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")

	var wkt = buf.String()
	var geomFromTxt = fmt.Sprintf(`st_geomfromtext('%v', %v)`, wkt, self.Src.SRID)
	query = fmt.Sprintf(
		`UPDATE %v SET %v=%v WHERE %v=%v;`,
		outputTable,
		self.Src.Config.GeometryColumn, geomFromTxt,
		self.Src.Config.IdColumn, id,
	)

	_, err = self.Src.Exec(query)
	if err != nil {
		panic(err)
	}
}
