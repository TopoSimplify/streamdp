package onlinedp

import (
	"log"
	"fmt"
	"bytes"
	"simplex/db"
	"simplex/dp"
)

func (self *OnlineDP) FindAndMarkDeformables(fid int, snapshotTbl string) {
	var temp = self.tempNodeIDTableName(fid)
	self.tempCreateNodeIdTable(temp)
	defer self.tempDropTable(temp)

	var query = fmt.Sprintf(
		`SELECT id,  gob  FROM  %v WHERE status=%v`, snapshotTbl, NullState,
	)
	var h, err = self.Src.Query(query)
	if err != nil {
		log.Panic(err)
	}

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

	self.MarkNullState(temp, snapshotTbl)
}

func (self *OnlineDP) MarkNullState(temp, snapshotTbl string) {
	const bufferSize = 200
	var buf = make([]int, 0)
	var query = fmt.Sprintf(`SELECT id  FROM  %v;`, temp)
	var h, err = self.Src.Query(query)

	if err != nil {
		log.Panic(err)
	}

	var id int
	for h.Next() {
		h.Scan(&id)
		buf = append(buf, id)
		if len(buf) > bufferSize {
			self.MarkNodesForDeformation(buf, snapshotTbl)
			buf = make([]int, 0) //reset
		}
	}

	//flush
	if len(buf) > 0 {
		self.MarkNodesForDeformation(buf, snapshotTbl)
	}
}

func (self *OnlineDP) MarkNodesForDeformation(selections []int, snapshotTbl string) int {
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
		snapshotTbl, buf.String(),
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

func (self *OnlineDP) FindAndMarkNullStateAsCollapsible(fid int, snapshotTbl string) {
	var query = fmt.Sprintf(
		`UPDATE %v SET status=%v WHERE status=%v;`,
		snapshotTbl, Collapsible, NullState,
	)
	var _, err = self.Src.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}
}

func (self *OnlineDP) FindAndSplitDeformables(fid int, snapshotTbl string) {
	var tempQ = fmt.Sprintf("%v_%v", self.tempQueryTableName(), fid)
	self.tempCreateTempQueryTable(tempQ)
	defer self.tempDropTable(tempQ)

	var worker = func(hull *db.Node) string {
		if hull.Range.Size() > 1 {
			var ha, hb = AtScoreSelection(hull, self.Score, dp.NodeGeometry)
			return ha.InsertSQL(snapshotTbl, self.Src.SRID, hb)
		}
		return hull.UpdateSQL(snapshotTbl, NullState)
	}

	var query = fmt.Sprintf(`SELECT id, gob  FROM  %v WHERE status=%v;`, snapshotTbl, SplitNode)
	var h, err = self.Src.Query(query)
	if err != nil {
		log.Panic(err)
	}

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

func (self *OnlineDP) FindAndCleanUpDeformables(fid int, snapshotTbl string) {
	var query = fmt.Sprintf(`DELETE FROM %v WHERE status=%v;`, snapshotTbl, SplitNode)
	var _, err = self.Src.Exec(query)
	if err != nil {
		log.Panic(err)
	}
}

func (self *OnlineDP) HasMoreDeformables(fid int, tbl string) bool {
	var query = fmt.Sprintf(
		`SELECT id FROM %v WHERE status=%vLIMIT 1;`,
		tbl, NullState,
	)
	var h, err = self.Src.Query(query)
	if err != nil {
		log.Panic(err)
	}
	var bln = false
	for h.Next() {
		bln = true
	}
	return bln
}
