package onlinedp

import (
	"time"
	"testing"
	"github.com/TopoSimplify/rng"
	"github.com/franela/goblin"
)

func TestContig(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("test contiguous", func() {
		g.It("should test contig", func() {
			g.Timeout(1 * time.Hour)
			//var options = &opts.Opts{MinDist: 10}
			var coords = linearCoords("LINESTRING ( 780 600, 740 620, 720 660, 720 700, 760 740, 820 760, 860 740, 880 720, 900 700, 880 660, 840 680, 820 700, 800 720, 760 710, 780 660, 820 640, 840 620, 860 580, 880 620, 830 660 )")
			var intRanges = [][]int{{0, 1}, {1, 3}, {3, 8}, {8, 13}, {13, 17}, {17, len(coords) - 1}}
			var ranges = make([]*rng.Rng, 0)
			for _, r := range intRanges {
				ranges = append(ranges, rng.Range(r[0], r[1]))
			}

			var hulls = createNodes(intRanges, coords)


			g.Assert(len(ranges)).Equal(len(hulls))

			for i := range hulls {
				g.Assert(hulls[i].Range.Equals(ranges[i]))
			}
			inters, contig, count := IsContiguous(hulls[0], hulls[1])
			g.Assert(contig && inters).IsTrue()
			g.Assert(count == 1).IsTrue()

			inters, contig, count = IsContiguous(hulls[1], hulls[2])
			g.Assert(contig && inters).IsTrue()
			g.Assert(count == 1).IsTrue()

			inters, contig, count = IsContiguous(hulls[3], hulls[2])
			g.Assert(contig && inters).IsTrue()
			g.Assert(count > 1).IsTrue()

			inters, contig, count = IsContiguous(hulls[3], hulls[4])
			g.Assert(contig && inters).IsTrue()
			g.Assert(count == 1).IsTrue()

			inters, contig, count = IsContiguous(hulls[4], hulls[2])
			g.Assert(inters).IsTrue()
			g.Assert(contig).IsFalse()
			g.Assert(count > 1).IsTrue()
		})

	})
}
