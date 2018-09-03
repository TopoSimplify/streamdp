package onlinedp

import (
	"github.com/TopoSimplify/db"
	"github.com/TopoSimplify/rng"
	"github.com/TopoSimplify/opts"
	"github.com/intdxdt/geom"
)

func (self *OnlineDP) selectDeformable(hull *db.Node) []*db.Node {
	var selections = make([]*db.Node, 0)
	//if hull is segment
	if hull.Range.Size() == 1 {
		return selections
	}
	//collinear: hull is line
	if hull.HullType == geom.GeoTypeLineString {
		return selections
	}

	// find hull neighbours
	var neighbs = self.FindNodeNeighbours(hull, self.Independent)

	// self intersection constraint
	if self.Options.AvoidNewSelfIntersects {
		self.ByFeatureClassIntersection(hull, neighbs, &selections)
	}

	// context_geom geometry constraint
	self.ValidateContextRelation(hull, &selections)
	return selections

}

//Constrain for self-intersection as a result of simplification
//returns boolean : is hull collapsible
func (self *OnlineDP) SelectBySelfIntersection(options *opts.Opts, hull *db.Node, selections *[]*db.Node, excludeRanges ...rng.Rng) bool {
	//assume hull is valid and proof otherwise
	var bln = true
	// find hull neighbours
	var neighbs = self.FindNodeNeighbours(hull, self.Independent, excludeRanges...)

	var hulls = self.SelectFeatureClass(hull, neighbs)
	for _, h := range hulls {
		//if bln & selection contains current hull : bln : false
		if bln && (h == hull) {
			bln = false //cmp &
		}
		*selections = append(*selections, h)
	}

	return bln
}

//Constrain for self-intersection as a result of simplification
func (self *OnlineDP) ByFeatureClassIntersection(hull *db.Node, neighbs []*db.Node, selections *[]*db.Node) bool {
	var bln = true
	//find hull neighbours
	var hulls = self.SelectFeatureClass(hull, neighbs)
	for _, h := range hulls {
		//if bln & selection contains current hull : bln : false
		if bln && (h == hull) {
			bln = false // cmp ref
		}
		*selections = append(*selections, h)
	}
	return bln
}

//find context_geom deformable hulls
func (self *OnlineDP) SelectFeatureClass(hull *db.Node, ctxHulls []*db.Node) []*db.Node {
	var n int
	var inters, contig bool
	var options = self.Options
	var dict = make(map[[2]int]*db.Node, 0)

	// for each item in the context_geom list
	for _, h := range ctxHulls {
		n = 0
		var sameFeature = hull.FID == h.FID
		// find which item to deform against current hull
		if sameFeature { // check for contiguity
			inters, contig, n = IsContiguous(hull, h)
		} else {
			// contiguity is by default false for different features
			contig = false
			var ga, gb = hull.Geometry(), h.Geometry()

			inters = ga.Intersects(gb)
			if inters {
				var interpts = ga.Intersection(gb)
				inters = len(interpts) > 0
				n = len(interpts)
			}
		}

		if !inters { // disjoint : nothing to do, continue
			continue
		}

		var sa, sb *db.Node
		if contig && n > 1 { // contiguity with overlap greater than a vertex
			sa, sb = contiguousCandidates(hull, h)
		} else if !contig {
			sa, sb = nonContiguousCandidates(options, hull, h)
		}

		// add candidate deformation hulls to selection list
		if sa != nil {
			dict[sa.Range.AsArray()] = sa
		}
		if sb != nil {
			dict[sb.Range.AsArray()] = sb
		}
	}

	var items = make([]*db.Node, 0)
	for _, v := range dict {
		items = append(items, v)
	}
	return items
}
