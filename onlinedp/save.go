package onlinedp

import (
	"fmt"
	"bytes"
	"simplex/db"
	"simplex/streamdp/pt"
	"github.com/intdxdt/geom"
	"simplex/streamdp/common"
	"log"
)

//obj.createOutputOnlineTable()
//obj.SaveSimplification()

//Find and merge simple segments
func (self *OnlineDP) SaveSimplification(fid int) {
	var outputTable = common.SimpleTable(self.Src.Table)
	var query = fmt.Sprintf(`
		SELECT node
		FROM %v
		WHERE fid=%v AND snapshot=%v
		ORDER BY fid asc, i asc;
	`,
		self.Src.Table,
		fid, common.Snap,
	)
	var h, err = self.Src.Query(query)
	if err != nil {
		panic(err)
	}

	var gob string
	var coordinates = make([][]*pt.Pt, 0)

	var index = -1
	for h.Next() {
		h.Scan(&gob)
		var o = db.Deserialize(gob)
		var i, j = 0, len(o.Coordinates)-1

		if index == -1 {
			coordinates = append(coordinates, []*pt.Pt{})
		}

		index = len(coordinates) - 1
		var last *pt.Pt
		if len(coordinates[index]) > 0 {
			n := len(coordinates[index]) - 1
			last = coordinates[index][n]
		}

		if last == nil {
			coordinates[index] = append(coordinates[index],
				&pt.Pt{Point: o.Coordinates[i], I: o.Range.I},
				&pt.Pt{Point: o.Coordinates[j], I: o.Range.J},
			)
		} else if last.I == o.Range.I {
			coordinates[index] = append(coordinates[index],
				&pt.Pt{Point: o.Coordinates[j], I: o.Range.J},
			)
		} else {
			fmt.Println(query)
			log.Panic("coordinates non contiguous")
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
	query = fmt.Sprintf(`
		INSERT INTO %v (id, geom)
		VALUES (%v, %v)
		ON CONFLICT (id)
			DO UPDATE SET geom = %v;
		`,
		outputTable, fid, geomFromTxt, geomFromTxt,
	)

	_, err = self.Src.Exec(query)
	if err != nil {
		log.Panic(err)
	}
}
