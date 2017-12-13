package onlinedp

import (
	"simplex/db"
	"simplex/rng"
	"simplex/opts"
	"github.com/intdxdt/geom"
)

func (self *OnlineDP) selectDeformable(hull *db.Node) []*db.Node {
	var selections = make([]*db.Node, 0)
	//if hull is segment
	if hull.Range.Size() == 1 {
		return selections
	}
	//if hull geometry is line then points are collinear
	if hull.HullType == geom.GeoType_LineString {
		return selections
	}

	// find hull neighbours
	var neighbs = self.FindNodeNeighbours(hull, self.Independent)

	// self intersection constraint
	// can self intersect with itself but not with other lines
	self.ByFeatureClassIntersection(hull, neighbs, &selections)

	// context_geom geometry constraint
	self.ValidateContextRelation(hull, &selections)
	return selections

}

//Constrain for self-intersection as a result of simplification
//returns boolean : is hull collapsible
func (self *OnlineDP) BySelfIntersection(options *opts.Opts, hull *db.Node, selections *[]*db.Node, excludeRanges ...*rng.Range) bool {
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

//select contiguous candidates
func contiguousCandidates(a, b *db.Node) (*db.Node, *db.Node) {
	//var selection = make([]*node.Node, 0)
	// compute sidedness relation between contiguous hulls to avoid hull flip

	// all hulls that are simple should be collapsible
	// if not collapsible -- add to selection for deformation
	// to reach collapsibility
	var sa, sb *db.Node
	//& the present should not affect the future
	if !a.Collapsible(b) {
		sa = a
	}

	// future should not affect the present
	if !b.Collapsible(a) {
		sb = b
	}
	return sa, sb
}

//select non-contiguous candidates
func nonContiguousCandidates(options *opts.Opts, a, b *db.Node) (*db.Node, *db.Node) {
	var aseg = a.Segment()
	var bseg = b.Segment()

	var aln = a.Polyline()
	var bln = b.Polyline()

	var asegGeom = aseg.Segment
	var bsegGeom = bseg.Segment

	var alnGeom = aln.Geometry
	var blnGeom = bln.Geometry

	var asegIntersBseg = asegGeom.Intersects(bsegGeom)
	var asegIntersBln = asegGeom.Intersects(blnGeom)
	var bsegIntersAln = bsegGeom.Intersects(alnGeom)
	var alnIntersBln = alnGeom.Intersects(blnGeom)
	var sa, sb *db.Node

	if asegIntersBseg && asegIntersBln && (!alnIntersBln) {
		sa = a
	} else if asegIntersBseg && bsegIntersAln && (!alnIntersBln) {
		sb = b
	} else if alnIntersBln {
		// find out whether is a shared vertex or overlap
		// is aseg inter bset  --- dist --- aln inter bln > relax dist
		var ptLns = alnGeom.Intersection(blnGeom)
		var atSeg = aseg.Intersection(bsegGeom)

		// if segs are disjoint but lines intersect, deform a&b
		if len(atSeg) == 0 && len(ptLns) > 0 {
			return a, b
		}

		for _, ptln := range ptLns {
			for _, ptseg := range atSeg {
				delta := ptln.Distance(ptseg)
				if delta > options.RelaxDist {
					return a, b
				}
			}
		}
	}
	return sa, sb
}
