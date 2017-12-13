package onlinedp

import (
	"simplex/db"
)

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
