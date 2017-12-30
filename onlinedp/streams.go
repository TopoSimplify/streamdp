package onlinedp

import (
	"log"
	"fmt"
	"bytes"
	"simplex/db"
	"simplex/dp"
)

func (self *OnlineDP) FindAndMarkDeformables() {
	var temp = self.tempNodeIDTableName()
	self.tempCreateNodeIdTable(temp)
	defer self.tempDropTable(temp)

	var worker = func(hull *db.Node) bool {
		// 0. find deformable node
		var selections = self.selectDeformable(hull)
		for _, o := range selections {
			self.tempInsertInNodeIdTable(temp, o.NID)
		}
		return true
	}

	var query = fmt.Sprintf(
		`SELECT id, fid, gob  FROM  %v WHERE status=%v;`,
		self.Src.NodeTable, NullState)
	var h, err = self.Src.Query(query)
	if err != nil {
		log.Panic(err)
	}

	var id, fid int
	var gob string
	for h.Next() {
		h.Scan(&id, &fid, &gob)
		var o = db.Deserialize(gob)
		o.NID, o.FID = id, fid
		worker(o)
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

	var id int
	for h.Next() {
		h.Scan(&id)
		buf = append(buf, id)
		if len(buf) > bufferSize {
			self.MarkNodesForDeformation(buf)
			buf = make([]int, 0) //reset
		}
	}

	//flush
	if len(buf) > 0 {
		self.MarkNodesForDeformation(buf)
	}
}

func (self *OnlineDP) MarkNodesForDeformation(selections []int) int {
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
		self.Src.NodeTable, buf.String(),
	)
	var r, err = self.Src.Exec(query)
	if err != nil {
		panic(err)
	}
	nrs, err := r.RowsAffected()
	if err != nil {
		panic(err)
	}
	return int(nrs)
}

func (self *OnlineDP) FindAndMarkNullStateAsCollapsible() {
	var query = fmt.Sprintf(
		`UPDATE %v SET status=%v WHERE status=%v;`,
		self.Src.NodeTable, Collapsible, NullState,
	)
	var _, err = self.Src.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}
}

func (self *OnlineDP) FindAndSplitDeformables() {
	var tempQ = self.tempQueryTableName()
	self.tempCreateTempQueryTable(tempQ)
	defer self.tempDropTable(tempQ)

	var worker = func(hull *db.Node) string {
		if hull.Range.Size() > 1 {
			var ha, hb = AtScoreSelection(hull, self.Score, dp.NodeGeometry)
			return ha.InsertSQL(self.Src.NodeTable, self.Src.SRID, hb)
		}
		return hull.UpdateSQL(self.Src.NodeTable, NullState)
	}

	var query = fmt.Sprintf(
		`SELECT id, fid, gob  FROM  %v WHERE status=%v;`,
		self.Src.NodeTable, SplitNode,
	)
	var h, err = self.Src.Query(query)
	if err != nil {
		panic(err)
	}
	var id, fid int
	var gob string
	var bufferSize = 100
	var buf = make([]string, 0)

	for h.Next() {
		h.Scan(&id, &fid, &gob)
		var o = db.Deserialize(gob)
		o.NID, o.FID = id, fid

		var selStr = worker(o)
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

func (self *OnlineDP) FindAndCleanUpDeformables() {
	var query = fmt.Sprintf(
		`DELETE FROM %v WHERE status=%v;`, self.Src.NodeTable, SplitNode,
	)
	var _, err = self.Src.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}
}

func (self *OnlineDP) HasMoreDeformables() bool {
	var query = fmt.Sprintf(
		`SELECT id FROM %v WHERE status=%v LIMIT 1;`,
		self.Src.NodeTable, NullState,
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
