package offset

import (
	"time"
	"testing"
	"github.com/TopoSimplify/streamdp/pt"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/math"
	"github.com/franela/goblin"
)

func TestOpw(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("cache test", func() {
		g.It("should test cache", func() {
			g.Timeout(1 * time.Hour)
			var data = []*pt.Pt{
				{Point: geom.PointXYZ(3.0, 1.6, 0.0), I: 0},
				{Point: geom.PointXYZ(3.0, 2.0, 1.0), I: 0},
				{Point: geom.PointXYZ(2.4, 2.8, 3.0), I: 0},
				{Point: geom.PointXYZ(0.5, 3.0, 4.5), I: 0},
				{Point: geom.PointXYZ(1.2, 3.2, 5.0), I: 0},
				{Point: geom.PointXYZ(1.4, 2.6, 6.0), I: 0},
				{Point: geom.PointXYZ(2.0, 3.5, 10.0), I: 0},
			}

			rootIndex, val := OPWMaxSEDOffset(data)
			var indx = rootIndex
			g.Assert(indx).Equal(3)
			g.Assert(math.Round(val, 5)).Equal(2.12121)

			indx, val = OPWMaxSEDOffset(data[:rootIndex+1])
			g.Assert(indx).Equal(2)
			g.Assert(math.Round(val, 5)).Equal(1.09949)

			indx, val = OPWMaxSEDOffset(data[rootIndex:])
			g.Assert(rootIndex + indx).Equal(5)
			g.Assert(math.Round(val, 5)).Equal(0.72710)

		})
	})
}

//node = tree.root.right;//root.left
//t.deepEqual(node[tree.key], [3, 6]);
//t.deepEqual(tree.int.index(node.int), 5);
//t.deepEqual(_.round(tree.int.val(node.int), 5), 0.72710);
//console.log(tree.print());
//t.end()
