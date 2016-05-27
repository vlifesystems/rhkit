package rulehunter

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/assessment"
	"github.com/vlifesystems/rulehunter/experiment"
	"github.com/vlifesystems/rulehunter/rule"
	"testing"
)

func TestAssessRules(t *testing.T) {
	rules := []*rule.Rule{
		rule.MustNew("band > 4"),
		rule.MustNew("band > 3"),
		rule.MustNew("cost > 1.2"),
	}
	aggregators := []*experiment.AggregatorDesc{
		&experiment.AggregatorDesc{"numIncomeGt2", "count", "income > 2"},
		&experiment.AggregatorDesc{"numBandGt4", "count", "band > 4"},
	}
	goals := []string{
		"numIncomeGt2 == 1",
		"numIncomeGt2 == 2",
		"numIncomeGt2 == 3",
		"numIncomeGt2 == 4",
		"numBandGt4 == 1",
		"numBandGt4 == 2",
		"numBandGt4 == 3",
		"numBandGt4 == 4",
	}
	fieldNames := []string{"income", "cost", "band"}
	records := [][]string{
		[]string{"3", "4.5", "4"},
		[]string{"3", "3.2", "7"},
		[]string{"2", "1.2", "4"},
		[]string{"0", "0", "9"},
	}
	dataset := NewLiteralDataset(fieldNames, records)
	experimentDesc := &experiment.ExperimentDesc{
		Title:         "",
		Dataset:       dataset,
		ExcludeFields: []string{},
		Aggregators:   aggregators,
		Goals:         goals,
		SortOrder:     []*experiment.SortDesc{},
	}
	experiment := mustNewExperiment(experimentDesc)
	wantIsSorted := false
	wantIsRefined := false
	wantNumRecords := int64(len(records))
	wantRuleAssessments := []*assessment.RuleAssessment{
		&assessment.RuleAssessment{
			Rule: rule.MustNew("band > 4"),
			Aggregators: map[string]*dlit.Literal{
				"numMatches":     dlit.MustNew("2"),
				"percentMatches": dlit.MustNew("50"),
				"numIncomeGt2":   dlit.MustNew("1"),
				"numBandGt4":     dlit.MustNew("2"),
				"goalsScore":     dlit.MustNew(1.001),
			},
			Goals: []*assessment.GoalAssessment{
				&assessment.GoalAssessment{"numIncomeGt2 == 1", true},
				&assessment.GoalAssessment{"numIncomeGt2 == 2", false},
				&assessment.GoalAssessment{"numIncomeGt2 == 3", false},
				&assessment.GoalAssessment{"numIncomeGt2 == 4", false},
				&assessment.GoalAssessment{"numBandGt4 == 1", false},
				&assessment.GoalAssessment{"numBandGt4 == 2", true},
				&assessment.GoalAssessment{"numBandGt4 == 3", false},
				&assessment.GoalAssessment{"numBandGt4 == 4", false},
			},
		},
		&assessment.RuleAssessment{
			Rule: rule.MustNew("band > 3"),
			Aggregators: map[string]*dlit.Literal{
				"numMatches":     dlit.MustNew("4"),
				"percentMatches": dlit.MustNew("100"),
				"numIncomeGt2":   dlit.MustNew("2"),
				"numBandGt4":     dlit.MustNew("2"),
				"goalsScore":     dlit.MustNew(0.002),
			},
			Goals: []*assessment.GoalAssessment{
				&assessment.GoalAssessment{"numIncomeGt2 == 1", false},
				&assessment.GoalAssessment{"numIncomeGt2 == 2", true},
				&assessment.GoalAssessment{"numIncomeGt2 == 3", false},
				&assessment.GoalAssessment{"numIncomeGt2 == 4", false},
				&assessment.GoalAssessment{"numBandGt4 == 1", false},
				&assessment.GoalAssessment{"numBandGt4 == 2", true},
				&assessment.GoalAssessment{"numBandGt4 == 3", false},
				&assessment.GoalAssessment{"numBandGt4 == 4", false},
			},
		},
		&assessment.RuleAssessment{
			Rule: rule.MustNew("cost > 1.2"),
			Aggregators: map[string]*dlit.Literal{
				"numMatches":     dlit.MustNew("2"),
				"percentMatches": dlit.MustNew("50"),
				"numIncomeGt2":   dlit.MustNew("2"),
				"numBandGt4":     dlit.MustNew("1"),
				"goalsScore":     dlit.MustNew(0.002),
			},
			Goals: []*assessment.GoalAssessment{
				&assessment.GoalAssessment{"numIncomeGt2 == 1", false},
				&assessment.GoalAssessment{"numIncomeGt2 == 2", true},
				&assessment.GoalAssessment{"numIncomeGt2 == 3", false},
				&assessment.GoalAssessment{"numIncomeGt2 == 4", false},
				&assessment.GoalAssessment{"numBandGt4 == 1", true},
				&assessment.GoalAssessment{"numBandGt4 == 2", false},
				&assessment.GoalAssessment{"numBandGt4 == 3", false},
				&assessment.GoalAssessment{"numBandGt4 == 4", false},
			},
		},
	}
	gotAssessment, err := AssessRules(rules, experiment)
	if err != nil {
		t.Errorf("AssessRules(%q, %q, %q, dataset) - err: %q",
			rules, aggregators, goals, err)
	}

	assessmentsMatch := areAssessmentsEqv(
		gotAssessment,
		wantNumRecords,
		wantIsSorted,
		wantIsRefined,
		wantRuleAssessments,
	)
	if !assessmentsMatch {
		t.Errorf("AssessRules(%q, %q, %q, dataset)\nassessments don't match\n - got: %s\n - wantRuleAssessments: %s, wantNumRecords: %d, wantIsSorted: %t, wantIsRefined: %t\n",
			rules, aggregators, goals, gotAssessment, wantRuleAssessments,
			wantNumRecords, wantIsSorted, wantIsRefined)
	}
}

func TestAssessRules_errors(t *testing.T) {
	cases := []struct {
		rules       []*rule.Rule
		aggregators []*experiment.AggregatorDesc
		goals       []string
		wantErr     error
	}{
		{[]*rule.Rule{rule.MustNew("band ^^ 3")},
			[]*experiment.AggregatorDesc{
				&experiment.AggregatorDesc{"numIncomeGt2", "count", "income > 2"},
			},
			[]string{"numIncomeGt2 == 1"},
			errors.New("Invalid operator: \"^\"")},
		{[]*rule.Rule{rule.MustNew("hand > 3")},
			[]*experiment.AggregatorDesc{
				&experiment.AggregatorDesc{"numIncomeGt2", "count", "income > 2"},
			},
			[]string{"numIncomeGt2 == 1"},
			errors.New("Variable doesn't exist: hand")},
		{[]*rule.Rule{rule.MustNew("band > 3")},
			[]*experiment.AggregatorDesc{
				&experiment.AggregatorDesc{"numIncomeGt2", "count", "bincome > 2"},
			},
			[]string{"numIncomeGt2 == 1"},
			errors.New("Variable doesn't exist: bincome")},
		{[]*rule.Rule{rule.MustNew("band > 3")},
			[]*experiment.AggregatorDesc{
				&experiment.AggregatorDesc{"numIncomeGt2", "count", "income > 2"},
			},
			[]string{"numIncomeGt == 1"},
			errors.New("Variable doesn't exist: numIncomeGt")},
	}
	fieldNames := []string{"income", "cost", "band"}
	records := [][]string{
		[]string{"3", "4.5", "4"},
	}
	dataset := NewLiteralDataset(fieldNames, records)
	for _, c := range cases {
		experimentDesc := &experiment.ExperimentDesc{
			Title:         "",
			Dataset:       dataset,
			ExcludeFields: []string{},
			Aggregators:   c.aggregators,
			Goals:         c.goals,
			SortOrder:     []*experiment.SortDesc{},
		}
		experiment := mustNewExperiment(experimentDesc)
		_, err := AssessRules(c.rules, experiment)
		if err == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("AssessRules(%q, %q) - err: %s, wantErr: %s",
				c.rules, experiment, err, c.wantErr)
		}
	}
}

func TestAssessRulesMP(t *testing.T) {
	aggregators := []*experiment.AggregatorDesc{
		&experiment.AggregatorDesc{"numIncomeGt2", "count", "income > 2"},
		&experiment.AggregatorDesc{"numBandGt4", "count", "band > 4"},
	}
	goals := []string{
		"numIncomeGt2 == 1",
		"numIncomeGt2 == 2",
		"numIncomeGt2 == 3",
		"numIncomeGt2 == 4",
		"numBandGt4 == 1",
		"numBandGt4 == 2",
		"numBandGt4 == 3",
		"numBandGt4 == 4",
	}
	fieldNames := []string{"income", "cost", "band"}
	records := [][]string{
		[]string{"3", "4.5", "4"},
		[]string{"3", "3.2", "7"},
		[]string{"2", "1.2", "4"},
		[]string{"0", "0", "9"},
	}
	cases := []struct {
		rules []*rule.Rule
	}{
		{[]*rule.Rule{
			rule.MustNew("band > 4"),
			rule.MustNew("band > 3"),
			rule.MustNew("cost > 1.2"),
		}},
		{[]*rule.Rule{
			rule.MustNew("band > 4"),
			rule.MustNew("cost > 1.2"),
		}},
		{[]*rule.Rule{}},
	}

	dataset := NewLiteralDataset(fieldNames, records)
	experimentDesc := &experiment.ExperimentDesc{
		Title:         "",
		Dataset:       dataset,
		ExcludeFields: []string{},
		Aggregators:   aggregators,
		Goals:         goals,
		SortOrder:     []*experiment.SortDesc{},
	}
	experiment := mustNewExperiment(experimentDesc)
	maxProcesses := 4
	for _, cs := range cases {
		wantAssessment, err :=
			AssessRules(cs.rules, experiment)
		if err != nil {
			t.Errorf("AssessRules(%q, %q) - err: %q",
				cs.rules, experiment, err)
		}
		c := make(chan *AssessRulesMPOutcome)
		progress := 0.0
		var gotAssessment *assessment.Assessment
		go AssessRulesMP(cs.rules, experiment, maxProcesses, c)

		numRuns := 0
		lastProgress := -1.0
		for o := range c {
			numRuns++
			progress = o.Progress
			if o.Err != nil {
				t.Errorf("AssessRulesMP(%q, %q, ...) - err: %q",
					cs.rules, experiment, o.Err)
			}
			if progress <= lastProgress {
				t.Errorf("AssessRulesMP(%q, %q, ...) - progress not increasing in order: this: %f, last: %f",
					cs.rules, experiment, progress, lastProgress)
			}
			if o.Finished {
				gotAssessment = o.Assessment
			}
		}
		if progress != 1.0 {
			t.Errorf("AssessRulesMP(%q, %q, ...) - progress didn't finish at 100, but: %d",
				cs.rules, experiment, progress)
		}
		if numRuns < len(cs.rules) {
			t.Errorf("AssessRulesMP(%q, %q, ...) - only made %d runs",
				cs.rules, experiment, numRuns)
		}
		assessmentsMatch := areAssessmentsEqv(
			gotAssessment,
			wantAssessment.NumRecords,
			wantAssessment.IsSorted(),
			wantAssessment.IsRefined(),
			wantAssessment.RuleAssessments,
		)
		if !assessmentsMatch {
			t.Errorf("AssessRulesMP(%q, %q, ...)\nassessments don't match\n - got: %s\n - want: %s\n",
				cs.rules, experiment, gotAssessment, wantAssessment)
		}
	}
}

/******************************
 *  Helper functions
 ******************************/

func mustNewExperiment(ed *experiment.ExperimentDesc) *experiment.Experiment {
	e, err := experiment.New(ed)
	if err != nil {
		panic(fmt.Sprintf("Can't create Experiment: %s", err))
	}
	return e
}

// Are the assessments equivalent.  The ruleAssessments must match
// but don't have to be in the same order if both assessments are
// unsorted. If both are unsorted then this will sort the assessments
func areAssessmentsEqv(
	got *assessment.Assessment,
	wantNumRecords int64,
	wantIsSorted bool,
	wantIsRefined bool,
	wantRuleAssessments []*assessment.RuleAssessment,
) bool {
	if got.NumRecords != wantNumRecords {
		return false
	}
	if got.IsSorted() != wantIsSorted {
		return false
	}
	if got.IsRefined() != wantIsRefined {
		return false
	}
	for _, gotRuleAssesment := range got.RuleAssessments {
		found := false
		for _, wantRuleAssessment := range wantRuleAssessments {
			if gotRuleAssesment.IsEqual(wantRuleAssessment) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
