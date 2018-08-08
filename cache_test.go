package main

import (
	"time"
	"testing"
	"github.com/TopoSimplify/streamdp/pt"
	"github.com/intdxdt/geom"
	"github.com/franela/goblin"
)

func TestCmp(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("cache test", func() {
		g.It("should test cache", func() {
			defer func() {
				r := recover()
				g.Assert(r != nil).IsTrue()
			}()
			var coords = []*pt.Pt{
				{Point: geom.PointXY(1, 2), I: 0},
				{Point: geom.PointXY(3, 4), I: 1},
				{Point: geom.PointXY(5, 6), I: 2},
				{Point: geom.PointXY(7, 8), I: 3},
				{Point: geom.PointXY(9, 10), I: 4},
			}
			var cache = make(Cache, 0)
			g.Assert(cache.size()).Equal(0)
			cache.append(coords...)
			g.Assert(cache.size()).Equal(5)
			g.Assert(cache.firstIndex()).Equal(0)
			g.Assert(cache.lastIndex()).Equal(4)
			g.Assert(cache.first().I).Equal(0)
			g.Assert(cache.last().I).Equal(4)

			g.Assert(cache.first().Point.Equals2D(coords[0].Point)).IsTrue()
			g.Assert(cache.last().Point.Equals2D(coords[4].Point)).IsTrue()

			cache.pop()

			g.Assert(cache.size()).Equal(4)
			g.Assert(cache.last().Point.Equals2D(coords[3].Point)).IsTrue()
			g.Assert(cache.last().I).Equal(3)

			cache.pop()

			g.Assert(cache.size()).Equal(3)
			g.Assert(cache.last().Point.Equals2D(coords[2].Point)).IsTrue()
			g.Assert(cache.last().I).Equal(2)

			var cacheClone = cache.clone()
			g.Assert(cacheClone.size()).Equal(cache.size())

			cache = make(Cache, 0)
			cache.append(coords...)
			g.Assert(cache.size()).Equal(5)
			cache.empty().append(coords[0], coords[1])
			g.Assert(cache.size()).Equal(2)

			cache.empty()
			g.Assert(cache.size()).Equal(0)
			cache.pop()

		})
		g.It("should test coords", func() {
			g.Timeout(10 * time.Minute)

			var coords = []*pt.Pt{
				{Point: geom.PointXY(1, 2), I: 0},
				{Point: geom.PointXY(3, 4), I: 1},
				{Point: geom.PointXY(5, 6), I: 2},
				{Point: geom.PointXY(7, 8), I: 3},
				{Point: geom.PointXY(9, 10), I: 4},
			}
			var cache = make(Cache, 0)
			cache.append(coords...)
			before, after := cache.split(2)
			g.Assert(before.size()).Equal(3)
			g.Assert(after.size()).Equal(2)

			before, after = cache.split(4)
			g.Assert(before.size()).Equal(5)
			g.Assert(after.size()).Equal(0)

			before, after = cache.split(1)
			g.Assert(before.size()).Equal(2)
			g.Assert(after.size()).Equal(3)

			before, after = cache.split(0)
			g.Assert(before.size()).Equal(1)
			g.Assert(after.size()).Equal(4)

		})

	})
}
