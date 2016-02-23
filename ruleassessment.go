/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import (
	"errors"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
)

type RuleAssessment struct {
	rule        *dexpr.Expr
	aggregators []Aggregator
	goals       []*dexpr.Expr
}

// Note: This clones the aggregators to ensure the results are specific
//       to the rule.
func NewRuleAssessment(rule *dexpr.Expr, aggregators []Aggregator,
	goals []*dexpr.Expr) *RuleAssessment {
	cloneAggregators := make([]Aggregator, len(aggregators))
	for i, a := range aggregators {
		cloneAggregators[i] = a.CloneNew()
	}
	return &RuleAssessment{rule: rule, aggregators: cloneAggregators,
		goals: goals}
}

func (ra *RuleAssessment) NextRecord(record map[string]*dlit.Literal) error {
	var ruleIsTrue bool
	var err error
	for _, aggregator := range ra.aggregators {
		ruleIsTrue, err = ra.rule.EvalBool(record, map[string]dexpr.CallFun{})
		if err != nil {
			return err
		}
		err = aggregator.NextRecord(record, ruleIsTrue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ra *RuleAssessment) GetAggregatorValue(
	name string, numRecords int64) (*dlit.Literal, bool) {
	for _, aggregator := range ra.aggregators {
		if aggregator.GetName() == name {
			return aggregator.GetResult(ra.aggregators, numRecords), true
		}
	}
	// TODO: Test and create specific error type
	err := errors.New("Aggregator doesn't exist")
	l, _ := dlit.New(err)
	return l, false
}
