package rulehunter

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/aggregators"
	"github.com/vlifesystems/rulehunter/goal"
	"github.com/vlifesystems/rulehunter/rule"
	"testing"
)

func TestNextRecord(t *testing.T) {
	// It is important for this test to reuse the goals
	// to ensure that they are cloned properly.
	inAggregators := []aggregators.AggregatorSpec{
		aggregators.MustNew("numIncomeGt2", "count", "income > 2"),
		aggregators.MustNew("numBandGt4", "count", "band > 4"),
		aggregators.MustNew("goalsScore", "goalsscore"),
	}
	records := [4]map[string]*dlit.Literal{
		map[string]*dlit.Literal{
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(4.5),
			"band":   dlit.MustNew(4),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(3.2),
			"band":   dlit.MustNew(7),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(2),
			"cost":   dlit.MustNew(1.2),
			"band":   dlit.MustNew(4),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(0),
			"cost":   dlit.MustNew(0),
			"band":   dlit.MustNew(9),
		},
	}
	goals := []*goal.Goal{
		goal.MustNew("numIncomeGt2 == 1"),
		goal.MustNew("numBandGt4 == 2"),
	}
	numRecords := int64(len(records))
	cases := []struct {
		rule             rule.Rule
		wantNumIncomeGt2 int64
		wantNumBandGt4   int64
		wantGoalsScore   float64
	}{
		{rule.NewGEFVI("band", 5), 1, 2, 2.0},
		{rule.NewGEFVI("band", 3), 2, 2, 0.001},
		{rule.NewGEFVF("cost", 1.3), 2, 1, 0},
	}
	for _, c := range cases {
		ra := newRuleAssessor(c.rule, inAggregators, goals)
		for _, record := range records {
			err := ra.NextRecord(record)
			if err != nil {
				t.Errorf("nextRecord(%q) rule: %s, aggregators: %q, goals: %q - err: %q",
					record, c.rule, inAggregators, goals, err)
			}
		}
		gotNumIncomeGt2, gt2Exists :=
			ra.GetAggregatorValue("numIncomeGt2", numRecords)
		if !gt2Exists {
			t.Errorf("numIncomeGt2 aggregator doesn't exist")
		}
		gotNumIncomeGt2Int, gt2IsInt := gotNumIncomeGt2.Int()
		if !gt2IsInt {
			t.Errorf("numIncomeGt2 aggregator can't be int")
		}
		if gotNumIncomeGt2Int != c.wantNumIncomeGt2 {
			t.Errorf("nextRecord() rule: %s, aggregators: %q, goals: %q - wantNumIncomeGt2: %d, got: %d",
				c.rule, inAggregators, goals, c.wantNumIncomeGt2, gotNumIncomeGt2Int)
		}
		gotNumBandGt4, gt4Exists :=
			ra.GetAggregatorValue("numBandGt4", numRecords)
		if !gt4Exists {
			t.Errorf("numBandGt4 aggregator doesn't exist")
		}
		gotNumBandGt4Int, gt4IsInt := gotNumBandGt4.Int()
		if !gt4IsInt {
			t.Errorf("numBandGt4 aggregator can't be int")
		}
		if gotNumBandGt4Int != c.wantNumBandGt4 {
			t.Errorf("nextRecord() rule: %s, aggregators: %q, goals: %q - wantNumBandGt4: %d, got: %d",
				c.rule, inAggregators, goals, c.wantNumBandGt4, gotNumBandGt4Int)
		}
		gotGoalsScore, goalsScoreExists :=
			ra.GetAggregatorValue("goalsScore", numRecords)
		if !goalsScoreExists {
			t.Errorf("goalsScore aggregator doesn't exist")
		}
		gotGoalsScoreFloat, goalsScoreIsFloat := gotGoalsScore.Float()
		if !goalsScoreIsFloat {
			t.Errorf("goalsScore aggregator can't be float")
		}
		if gotGoalsScoreFloat != c.wantGoalsScore {
			t.Errorf("nextRecord() rule: %s, aggregators: %q, goals: %q - wantGoalsScore: %f, got: %f",
				c.rule, inAggregators, goals, c.wantGoalsScore, gotGoalsScore)
		}
	}
}

func TestNextRecord_Errors(t *testing.T) {
	records := [4]map[string]*dlit.Literal{
		map[string]*dlit.Literal{"income": dlit.MustNew(3), "band": dlit.MustNew(4)},
		map[string]*dlit.Literal{"income": dlit.MustNew(3), "band": dlit.MustNew(7)},
		map[string]*dlit.Literal{"income": dlit.MustNew(2), "band": dlit.MustNew(4)},
		map[string]*dlit.Literal{"income": dlit.MustNew(0), "band": dlit.MustNew(9)},
	}
	goals := []*goal.Goal{goal.MustNew("numIncomeGt2 == 1")}
	cases := []struct {
		rule        rule.Rule
		aggregators []aggregators.AggregatorSpec
		wantErr     error
	}{
		{rule.NewGEFVI("band", 4),
			[]aggregators.AggregatorSpec{
				aggregators.MustNew("numIncomeGt2", "count", "fred > 2")},
			dexpr.ErrInvalidExpr{
				Expr: "fred > 2",
				Err:  dexpr.ErrVarNotExist("fred"),
			},
		},
		{rule.NewGEFVI("band", 4),
			[]aggregators.AggregatorSpec{
				aggregators.MustNew("numIncomeGt2", "count", "income > 2")}, nil},
		{rule.NewGEFVI("hand", 4),
			[]aggregators.AggregatorSpec{
				aggregators.MustNew("numIncomeGt2", "count", "income > 2")},
			rule.InvalidRuleError{Rule: rule.NewGEFVI("hand", 4)},
		},
	}
	for _, c := range cases {
		ra := newRuleAssessor(c.rule, c.aggregators, goals)
		for _, record := range records {
			err := ra.NextRecord(record)
			if !errorMatch(c.wantErr, err) {
				t.Errorf("NextRecord(%q) rule: %q, aggregators: %q, goals: %q err: %q, wantErr: %q",
					record, c.rule, c.aggregators, goals, err, c.wantErr)
				return
			}
		}
	}
}
