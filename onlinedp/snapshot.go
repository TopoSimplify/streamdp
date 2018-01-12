package onlinedp

import (
	"fmt"
	"log"
	"bytes"
	"simplex/db"
	"simplex/dp"
	"simplex/streamdp/common"
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
			UPDATE %v
			SET snapshot=%v
			WHERE fid=%v;
		`,
		self.Src.Table,
		snapState,
		fid,
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
		self.Src.Table,
		Collapsible,
		NullState, fid, common.Snap,
	)
	if _, err := self.Src.Exec(query); err != nil {
		log.Fatalln(err)
	}
}

func (self *OnlineDP) CleanUpDeformables(fid int) {
	var query = fmt.Sprintf(`
		DELETE FROM %v
		WHERE status=%v AND fid=%v AND snapshot=%v;
		`,
		self.Src.Table,
		SplitNode, fid, common.Snap,
	)
	if _, err := self.Src.Exec(query); err != nil {
		log.Panic(err)
	}
}

func (self *OnlineDP) MarkDeformables(fid int) {
	var temp = self.tempNodeIDTableName(fid)
	self.tempCreateNodeIdTable(temp)
	defer self.tempDropTable(temp)

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

	for h.Next() {
		var id int
		var gob string
		h.Scan(&id, &gob)
		var o = db.Deserialize(gob)
		o.NID, o.FID = id, fid

		var selections = self.selectDeformable(o)
		for _, s := range selections {
			self.tempInsertInNodeIdTable(temp, s.NID)
		}
	}

	self.MarkNullState(temp)
}

func (self *OnlineDP) MarkNullState(temp string) {
	const bufferSize = 200
	var buf = make([]int, 0)
	var query = fmt.Sprintf(`SELECT id  FROM  %v;`, temp)

	var h, err = self.Src.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer h.Close()

	var id int
	for h.Next() {
		h.Scan(&id)
		buf = append(buf, id)
		if len(buf) > bufferSize {
			self.markDeformableNodes(buf)
			buf = make([]int, 0) //reset
		}
	}

	//flush
	if len(buf) > 0 {
		self.markDeformableNodes(buf)
	}
}

func (self *OnlineDP) markDeformableNodes(selections []int) int {
	if len(selections) == 0 {
		return 0
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
			u2.id = u.id;`,
		self.Src.Table, buf.String(),
	)
	var r, err = self.Src.Exec(query)
	if err != nil {
		log.Panic(err)
	}
	nrs, err := r.RowsAffected()
	if err != nil {
		log.Panic(err)
	}
	return int(nrs)
}

func (self *OnlineDP) SplitDeformables(fid int) {
	var tempQ = fmt.Sprintf("%v_%v", self.tempQueryTableName(), fid)
	self.tempCreateTempQueryTable(tempQ)
	defer self.tempDropTable(tempQ)

	var worker = func(hull *db.Node) string {
		if hull.Range.Size() > 1 {
			var ha, hb = AtScoreSelection(hull, self.Score, dp.NodeGeometry)
			var vals = common.SnapshotNodeColumnValues(self.Src.SRID, common.Snap, ha, hb)
			return db.SQLInsertIntoTable(self.Src.Table, common.NodeColumnFields, vals)
		}
		return hull.UpdateSQL(self.Src.Table, NullState)
	}

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

	var bufferSize = 100
	var buf = make([]string, 0)
	for h.Next() {
		var id int
		var gob string

		h.Scan(&id, &gob)
		o := db.Deserialize(gob)
		o.NID, o.FID = id, fid

		selStr := worker(o)
		buf = append(buf, selStr)
		if len(buf) > bufferSize {
			self.tempInsertInTOTempQueryTable(tempQ, buf)
			buf = make([]string, 0) //reset
		}
	}

	//flush buf
	if len(buf) > 0 {
		self.tempInsertInTOTempQueryTable(tempQ, buf)
	}
	self.tempExecuteQueries(tempQ)
}
