package onlinedp

import (
	"simplex/db"
	"simplex/lnr"
	"simplex/opts"
)

const (
	NullState   = iota
	Collapsible
	SplitNode
)

const (
	concurProcs       = 8
	MergeFragmentSize = 1
	EpsilonDist       = 1.0e-5
)

type OnlineDP struct {
	Src         *db.DataSrc
	Const       *db.DataSrc
	Options     *opts.Opts
	Score       lnr.ScoreFn
	Independent bool
}


func (self *OnlineDP) ScoreRelation(val float64) bool {
	return val <= self.Options.Threshold
}


func NewOnlineDP(src, constraints *db.DataSrc, options *opts.Opts,
	offsetScore lnr.ScoreFn, independent bool) *OnlineDP {
	return &OnlineDP{
		Src:         src,
		Const:       constraints,
		Options:     options,
		Score:       offsetScore,
		Independent: independent,
	}
}

