package onlinedp

import (
	"simplex/db"
	"simplex/rng"
	"simplex/ctx"
)

func (self *OnlineDP) ValidateMerge(hull *db.Node, excludeRanges ...*rng.Range) bool {
	var bln = true
	var sideEffects = make([]*db.Node, 0)

	// self intersection constraint
	if self.Options.AvoidNewSelfIntersects {
		bln = self.SelectBySelfIntersection(self.Options, hull, &sideEffects, excludeRanges...)
	}

	if len(sideEffects) != 0 || !bln {
		return false
	}

	// context geometry constraint
	bln = self.ValidateContextRelation(hull, &sideEffects)
	return bln && (len(sideEffects) == 0)
}

//Constrain for context neighbours
// finds the collapsibility of hull with respect to context hull neighbours
// if hull is deformable, its added to selections
func (self *OnlineDP) ValidateContextRelation(hull *db.Node, selections *[]*db.Node) bool {
	if !self.checkContextRelations() {
		return true
	}
	var bln = true
	// find context neighbours - if valid
	var ctxs = self.FindContextNeighbours(hull.WTK, self.Options.MinDist)
	var ctxtgeoms = (&ctx.ContextGeometries{}).SetData(ctxs)

	if bln && self.Options.GeomRelation {
		bln = ByGeometricRelation(hull, ctxtgeoms)
	}

	if bln && self.Options.DistRelation {
		bln = ByMinDistRelation(self.Options, hull, ctxtgeoms)
	}

	if bln && self.Options.DirRelation {
		var neibs = self.FindNodeNeighbours(hull, self.Independent)
		bln = BySideRelation(hull, ctxtgeoms)
	}

	if !bln {
		*selections = append(*selections, hull)
	}

	return bln
}

func (self *OnlineDP) checkContextRelations() bool {
	return self.Options.GeomRelation || self.Options.DistRelation || self.Options.DirRelation
}
