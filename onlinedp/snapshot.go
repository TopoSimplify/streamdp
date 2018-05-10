package onlinedp

import (
	"fmt"
	"log"
	"bytes"
	"github.com/TopoSimplify/db"
	"github.com/TopoSimplify/dp"
	"github.com/TopoSimplify/streamdp/common"
	"github.com/intdxdt/fan"
)

func (self *OnlineDP) HasMoreDeformables(fid int) bool {
	var query = fmt.Sprintf(`
			SELECT id
			FROM %v
			WHERE status=%v AND fid=%v AND snapshot=%v
			LIMIT 1;
		`,
		self.Src.Table,
		NullState, fid, common.Snap,
	)

	var h, err = self.Src.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer h.Close()

	var bln = false
	for h.Next() {
		bln = true
	}
	return bln
}

func (self *OnlineDP) MarkSnapshot(fid int, snapState int) {
	var query = fmt.Sprintf(`
			UPDATE %v SET snapshot=%v WHERE fid=%v;
		`,
		self.Src.Table, snapState, fid,
	)
	var _, err = self.Src.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}
}

func (self *OnlineDP) MarkNullStateAsCollapsible(fid int) {
	var query = fmt.Sprintf(`
			UPDATE %v
			SET status=%v
			WHERE status=%v AND fid=%v AND snapshot=%v;
		`,
		self.Src.Table, Collapsible, NullState, fid, common.Snap,
	)
	if _, err := self.Src.Exec(query); err != nil {
		log.Fatalln(err)
	}
}

func (self *OnlineDP) CleanUpDeformables(fid int) {
	var query = fmt.Sprintf(`
		DELETE FROM %v WHERE status=%v AND fid=%v AND snapshot=%v;
		`, self.Src.Table, SplitNode, fid, common.Snap,
	)
	if _, err := self.Src.Exec(query); err != nil {
		log.Panic(err)
	}
}

func (self *OnlineDP) MarkDeformables(fid int) {
	const concur = 4
	var query = fmt.Sprintf(`
			SELECT id, node
			FROM  %v
			WHERE status=%v AND fid=%v AND snapshot=%v;
		`,
		self.Src.Table, NullState, fid, common.Snap,
	)

	var h, err = self.Src.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer h.Close()

	var stream = make(chan interface{}, 4*concur)
	var exit = make(chan struct{})
	defer close(exit)

	go func() {
		for h.Next() {
			var id int
			var gob string
			h.Scan(&id, &gob)
			var o = db.Deserialize(gob)
			o.NID, o.FID = id, fid
			stream <- o
		}
		close(stream)
	}()

	var worker = func(v interface{}) interface{} {
		return self.selectDeformable(v.(*db.Node))
	}
	var out = fan.Stream(stream, worker, concur, exit)

	const bufferSize = 200
	var buf = make(map[int]struct{})

	for selections := range out {
		var nodes = selections.([]*db.Node)
		for _, node := range nodes {
			buf[node.NID] = struct{}{}
			if len(buf) > bufferSize {
				self.markDeformableNodes(buf)
				buf = make(map[int]struct{}) //reset
			}
		}
	}

	//flush
	if len(buf) > 0 {
		self.markDeformableNodes(buf)
	}

}

func (self *OnlineDP) markDeformableNodes(buffer map[int]struct{}) {
	var selections = make([]int, 0, len(buffer))
	for id := range buffer {
		selections = append(selections, id)
	}

	if len(selections) == 0 {
		return
	}
	var buf bytes.Buffer
	var k = len(selections) - 1
	for i, nid := range selections {
		buf.WriteString(fmt.Sprintf(`(%v, %v)`, nid, SplitNode))
		if i < k {
			buf.WriteString(`,`)
		}
	}

	var query = fmt.Sprintf(`
		UPDATE %v AS u
		SET status= u2.status
		FROM
			( VALUES %v ) AS u2 ( id, status )
		WHERE
			u2.id = u.id;
	`,
		self.Src.Table, buf.String(),
	)
	var _, err = self.Src.Exec(query)
	if err != nil {
		log.Panic(err)
	}
}

func (self *OnlineDP) SplitDeformables(fid int) {
	const concur = 4
	const bufferSize = 100

	var query = fmt.Sprintf(`
			SELECT id, node
			FROM  %v
			WHERE status=%v AND fid=%v AND snapshot=%v;
		`,
		self.Src.Table, SplitNode, fid, common.Snap,
	)

	var h, err = self.Src.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer h.Close()

	var stream = make(chan interface{}, 4*concur)
	var exit = make(chan struct{})
	defer close(exit)

	go func() {
		for h.Next() {
			var id int
			var gob string

			h.Scan(&id, &gob)
			o := db.Deserialize(gob)
			o.NID, o.FID = id, fid
			stream <- o
		}
		close(stream)
	}()

	var worker = func(v interface{}) interface{} {
		var hull = v.(*db.Node)
		if hull.Range.Size() > 1 {
			var ha, hb = AtScoreSelection(hull, self.Score, dp.NodeGeometry)
			var vals = common.SnapshotNodeColumnValues(self.Src.SRID, common.Snap, ha, hb)
			return db.SQLInsertIntoTable(self.Src.Table, common.NodeColumnFields, vals)
		}
		return hull.UpdateSQL(self.Src.Table, NullState)
	}

	var out = fan.Stream(stream, worker, concur, exit)

	var buf = make([]string, 0)

	var processQuery = func(buf []string) {
		for _, query := range buf {
			if _, err := self.Src.Exec(query); err != nil {
				fmt.Println(query)
				log.Panic(err)
			}
		}
	}

	for insQ := range out {
		buf = append(buf, insQ.(string))
		if len(buf) > bufferSize {
			processQuery(buf)
			buf = make([]string, 0) //reset
		}
	}

	//flush buf
	if len(buf) > 0 {
		processQuery(buf)
	}
}
