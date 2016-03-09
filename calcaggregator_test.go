package main

import (
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"testing"
)

func TestCalcGetResult(t *testing.T) {
	aggregators := []Aggregator{
		mustNewCalcAggregator("a", "3 + 4"),
		mustNewCalcAggregator("b", "5 + 6"),
		mustNewCalcAggregator("c", "a + b"),
		mustNewCalcAggregator("2NumRecords", "numRecords * 2"),
		mustNewCalcAggregator("d", "a + e"),
	}
	want := []*dlit.Literal{
		mustNewLit(7),
		mustNewLit(11),
		mustNewLit(18),
		mustNewLit(24),
		mustNewLit(dexpr.ErrInvalidExpr("Variable doesn't exist: e")),
	}
	numRecords := int64(12)
	for i, aggregator := range aggregators {
		got := aggregator.GetResult(aggregators, numRecords)
		if got.String() != want[i].String() {
			t.Errorf("GetResult() i: %d got: %s, want: %s", i, got, want[i])
		}
	}
}

func TestCalcCloneNew(t *testing.T) {
	aggregators := []Aggregator{
		mustNewCalcAggregator("a", "3 + 4"),
		mustNewCalcAggregator("b", "5 + 6"),
		mustNewCalcAggregator("c", "a + b"),
	}
	numRecords := int64(12)
	aggregatorD := aggregators[2].CloneNew()
	gotC := aggregators[2].GetResult(aggregators, numRecords)
	gotD := aggregatorD.GetResult(aggregators, numRecords)

	if gotC.String() != gotD.String() && gotC.String() != "18" {
		t.Errorf("CloneNew() gotC: %s, gotD: %s", gotC, gotD)
	}
}
