package onlinedp

import (
	"fmt"
	"time"
	"sort"
	"testing"
	"github.com/TopoSimplify/db"
	"github.com/TopoSimplify/rng"
	"database/sql"
	"github.com/TopoSimplify/opts"
	"github.com/TopoSimplify/streamdp/offset"
	"github.com/franela/goblin"
	"github.com/TopoSimplify/streamdp/common"
)

func TestOnline(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("test online", func() {
		g.It("should test online", func() {
			g.Timeout(1 * time.Hour)

			var coords = linearCoords("LINESTRING ( 780 600, 760 610, 740 620, 720 660, 720 700, 760 740, 820 760, 860 740, 880 720, 900 700, 880 660, 840 680, 820 700, 800 720, 760 710, 780 660, 820 640, 840 620, 860 580, 880 620, 830 660 )")
			var intRanges = [][]int{
				{0, 2}, {2, 4}, {4, 9}, {9, 14},
				{14, 15}, {15, 16}, {16, 17}, {17, 18},
				{18, 19}, {19, len(coords) - 1}}
			var hulls = createNodes(intRanges, coords)
			//printNodes(hulls)
			var fid = hulls[0].FID

			//var options = &opts.Opts{MinDist: 10}
			var serverCfg = loadConfig(ServerCfg)
			g.Assert(serverCfg != nil).IsTrue()
			var cfg = serverCfg.DBConfig()

			var sqlSrc, err = sql.Open("postgres", fmt.Sprintf(
				"user=%s password=%s dbname=%s sslmode=disable",
				cfg.User, cfg.Password, cfg.Database,
			))
			var src = &db.DataSrc{
				Src:    sqlSrc,
				Config: cfg,
				SRID:   serverCfg.SRID,
				Dim:    serverCfg.Dim,
				Table:  serverCfg.Table,
			}
			g.Assert(err == nil).IsTrue()
			g.Assert(sqlSrc != nil).IsTrue()
			var options = &opts.Opts{
				Threshold:              100,
				MinDist:                0,
				AvoidNewSelfIntersects: true,
			}
			var inst = NewOnlineDP(src, nil, options, offset.MaxOffset, true)
			g.Assert(inst != nil).IsTrue()
			g.Assert(inst.ScoreRelation(78)).IsTrue()
			g.Assert(inst.ScoreRelation(109)).IsFalse()

			//create online table
			err = db.CreateNodeTable(src)
			g.Assert(err == nil).IsTrue()
			insertNodesIntoOnlineTable(src, hulls)

			//-------------------snap----------------------------------------------
			inst.MarkSnapshot(fid, common.Snap)
			defer inst.MarkSnapshot(fid, common.UnSnap)
			//-------------------deformation---------------------------------------

			g.Assert(len(inst.selectDeformable(hulls[0]))).Equal(0)
			g.Assert(len(inst.selectDeformable(hulls[1]))).Equal(0)
			g.Assert(len(inst.selectDeformable(hulls[len(hulls)-1]))).Equal(0)
			g.Assert(len(inst.selectDeformable(hulls[2]))).Equal(1)
			var deformables = inst.selectDeformable(hulls[2])
			g.Assert(deformables[0].Range.Equals(hulls[2].Range)).IsTrue()

			//----------------------hull-------------------------------------------
			var hull = hulls[2]
			//----------------------find neighbours-------------------------------------------
			var neighbs = inst.FindNodeNeighbours(hull, true, hull.Range)
			sort.Sort(db.Nodes(neighbs))
			g.Assert(len(neighbs)).Equal(3)
			g.Assert(neighbs[0].Range.Equals(hulls[1].Range)).IsTrue()
			g.Assert(neighbs[1].Range.Equals(hulls[3].Range)).IsTrue()
			g.Assert(neighbs[2].Range.Equals(hulls[4].Range)).IsTrue()

			//----------------------self-inters-------------------------------------------
			var sideEffects = make([]*db.Node, 0)
			// self intersection constraint
			var collapsible = inst.SelectBySelfIntersection(options, hulls[1], &sideEffects)
			g.Assert(collapsible).IsTrue()
			g.Assert(len(sideEffects)).Equal(0)
			collapsible = inst.SelectBySelfIntersection(options, hulls[2], &sideEffects)
			g.Assert(collapsible).IsFalse()
			g.Assert(len(sideEffects)).Equal(1)

			g.Assert(inst.HasMoreDeformables(fid)).IsTrue()
			//----------------------simplify-------------------------------------------
			// 1.find and mark deformable nodes
			inst.MarkDeformables(fid)
			var nodes = queryNodesByStatus(src, SplitNode)
			g.Assert(len(nodes)).Equal(1)
			g.Assert(nodes[0].Range.Equals(hulls[2].Range)).IsTrue()

			//2.mark valid nodes as collapsible
			inst.MarkNullStateAsCollapsible(fid)
			nodes = queryNodesByStatus(src, Collapsible)
			g.Assert(len(nodes)).Equal(len(hulls) - 1)

			//3.find and split deformable nodes, set status as nullstate
			inst.SplitDeformables(fid)

			//4.remove deformable nodes
			inst.CleanUpDeformables(fid)
			nodes = queryNodesByStatus(src, NullState)
			sort.Sort(db.Nodes(nodes))
			g.Assert(len(nodes)).Equal(2)
			g.Assert(nodes[0].Range.Equals(rng.Range(4, 6))).IsTrue()
			g.Assert(nodes[1].Range.Equals(rng.Range(6, 9))).IsTrue()

			//-------------has more deformables--------------------------------
			g.Assert(inst.HasMoreDeformables(fid)).IsTrue()

			//1.find and mark deformable nodes
			inst.MarkDeformables(fid)
			nodes = queryNodesByStatus(src, SplitNode)
			g.Assert(len(nodes)).Equal(0)

			//2.mark valid nodes as collapsible
			inst.MarkNullStateAsCollapsible(fid)
			nodes = queryNodesByStatus(src, Collapsible)
			g.Assert(len(nodes)).Equal(len(hulls) + 1)

			//3.find and split deformable nodes, set status as nullstate
			inst.SplitDeformables(fid)

			//4.remove deformable nodes
			inst.CleanUpDeformables(fid)
			nodes = queryNodesByStatus(src, NullState)
			g.Assert(len(nodes)).Equal(0)

			inst.AggregateSimpleSegments(fid, 1)
			inst.AggregateSimpleSegments(fid, 2)
		})
	})
}
