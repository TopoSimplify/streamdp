package main

import (
	"time"
	"testing"
	"simplex/streamdp/data"
	"github.com/franela/goblin"
	"simplex/opts"
	"simplex/streamdp/offset"
	"simplex/db"
)

func generatePings(size int) []*data.Ping {
	var t = time.Now()
	var pts = make([]*data.Ping, 0)
	for i := 0; i < size; i++ {
		v := float64(i)
		t = t.Add(1 * time.Second)
		pts = append(pts, &data.Ping{X: v, Y: 0, Time: t})
	}
	return pts
}

func buildNodes(pts []*data.Ping, inst *OPW) []*db.Node {
	var nodes = make([]*db.Node, 0)
	for _, p := range pts {
		n := inst.Push(p)
		if n != nil {
			nodes = append(nodes, n)
		}
	}
	for _, n := range inst.Done() {
		nodes = append(nodes, n)
	}
	return nodes
}

func TestOPW(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("opw test", func() {
		g.It("should test opw nodes", func() {
			g.Timeout(1 * time.Hour)
			var pts = generatePings(1000)
			var options = &opts.Opts{
				Threshold:              5000.0,
				MinDist:                500.0,
				RelaxDist:              100.0,
				AvoidNewSelfIntersects: false,
				GeomRelation:           false,
				DistRelation:           false,
				DirRelation:            false,
			}

			var inst = NewOPW(options, NOPW, offset.MaxOffset, 300)
			var nodes = buildNodes(pts, inst)
			g.Assert(len(nodes)).Equal(3)

			inst = NewOPW(options, BOPW, offset.MaxOffset, 300)
			g.Assert(len(nodes)).Equal(3)

			inst = NewOPW(options, NOPW, offset.MaxOffset, 100)
			nodes = buildNodes(pts, inst)
			g.Assert(len(nodes)).Equal(10)

			inst = NewOPW(options, BOPW, offset.MaxOffset, 100)
			nodes = buildNodes(pts, inst)
			g.Assert(len(nodes)).Equal(10)
		})

	})
}