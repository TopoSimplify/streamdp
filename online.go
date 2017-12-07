package main

import (
	"simplex/db"
	"simplex/opts"
	"simplex/lnr"
)


const (
	NullState   = iota
	Collapsible
	SplitNode
)

const (
	concurProcs       = 8
	MergeFragmentSize = 1
)

type OnlineDP struct {
	Src         *db.DataSrc
	Const       *db.DataSrc
	Options     *opts.Opts
	Score       lnr.ScoreFn
	Independent bool
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

