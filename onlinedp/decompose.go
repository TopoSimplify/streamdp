package onlinedp

import (
	"log"
	"fmt"
	"simplex/dp"
	"simplex/db"
	"github.com/intdxdt/fan"
	"github.com/paulmach/go.geojson"
)

type idG struct {
	Id int
	G  *geojson.Geometry
}

func (self *OnlineDP) Decompose() {
	var err = db.CreateNodeTable(self.Src)
	if err != nil {
		log.Fatalln(err)
	}

	var stream = make(chan interface{})
	var exit = make(chan struct{})
	defer close(exit)

	go func() {
		var query = fmt.Sprintf(
			"Select %v, ST_AsGeoJson(%v) as geom  from  %v;",
			self.Src.Config.IdColumn, self.Src.Config.GeometryColumn, self.Src.Config.Table,
		)
		var h, err = self.Src.Query(query)
		if err != nil {
			log.Fatalln(err)
		}

		for h.Next() {
			var o = &idG{}
			h.Scan(&o.Id, &o.G)
			stream <- o
		}
		close(stream)
	}()

	var worker = func(v interface{}) interface{} {
		var o = v.(*idG)
		var lns = FlattenLinearGeoms(o.Id, o.G)
		for _, ln := range lns {
			var nds = self.srcDecomposition(ln)
			db.BulkLoadNodes(self.Src, nds)
		}
		return true
	}
	var out = fan.Stream(stream, worker, concurProcs, exit)
	for range out {
	}
}

func (self *OnlineDP) srcDecomposition(lnr *LinearGeometry) []*db.Node {
	var nodes = dp.New(lnr.Coordinates, self.Options, self.Score).Decompose()
	var nds = make([]*db.Node, 0, len(nodes))

	for _, n := range nodes {
		n.SetId(fmt.Sprintf("id:%v/part:%v", lnr.Id, lnr.Part))
		var nd = db.NewDBNode(n)
		nd.FID, nd.Part = lnr.Id, lnr.Part
		nds = append(nds, nd)
	}
	return nds

}
