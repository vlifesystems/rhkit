/*
	Copyright (C) 2016 vLife Systems Ltd <http://vlifesystems.com>
	This file is part of Rulehunter.

	Rulehunter is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	Rulehunter is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with Rulehunter; see the file COPYING.  If not, see
	<http://www.gnu.org/licenses/>.
*/

package rulehunter

import (
	"errors"
	"fmt"
	"github.com/vlifesystems/rulehunter/aggregators"
	"github.com/vlifesystems/rulehunter/assessment"
	"github.com/vlifesystems/rulehunter/experiment"
	"github.com/vlifesystems/rulehunter/input"
	"github.com/vlifesystems/rulehunter/internal/ruleassessor"
	"github.com/vlifesystems/rulehunter/rule"
)

// Assess the rules using a single thread
func AssessRules(
	rules []*rule.Rule,
	e *experiment.Experiment,
) (*assessment.Assessment, error) {
	var allAggregators []aggregators.Aggregator
	var numRecords int64
	var err error

	allAggregators, err = addDefaultAggregators(e.Aggregators)
	if err != nil {
		return &assessment.Assessment{}, err
	}

	ruleAssessors := make([]*ruleassessor.RuleAssessor, len(rules))
	for i, rule := range rules {
		ruleAssessors[i] = ruleassessor.New(rule, allAggregators, e.Goals)
	}

	// The input must be cloned to be thread safe when AssessRules called by
	// AssessRulesMP
	inputClone, err := e.Input.Clone()
	if err != nil {
		return &assessment.Assessment{}, err
	}
	defer e.Input.Close()
	numRecords, err = processInput(inputClone, ruleAssessors)
	if err != nil {
		return &assessment.Assessment{}, err
	}
	goodRuleAssessors, err := filterGoodRuleAssessors(ruleAssessors, numRecords)
	if err != nil {
		return &assessment.Assessment{}, err
	}

	assessment, err := assessment.New(numRecords, goodRuleAssessors, e.Goals)
	return assessment, err
}

type AssessRulesMPOutcome struct {
	Assessment *assessment.Assessment
	Err        error
	Progress   float64
	Finished   bool
}

// Goroutine to assess the rules using multiple processes and report on
// progress through 'ec' channel
func AssessRulesMP(
	rules []*rule.Rule,
	e *experiment.Experiment,
	maxProcesses int,
	ec chan *AssessRulesMPOutcome,
) {
	var assessment *assessment.Assessment
	var isError bool
	ic := make(chan *assessRulesCOutcome)
	numRules := len(rules)
	if numRules < 2 {
		assessment, err := AssessRules(rules, e)
		ec <- &AssessRulesMPOutcome{assessment, err, 1.0, true}
		close(ec)
		return
	}
	progressIntervals := 1000
	numProcesses := 0
	if numRules < progressIntervals {
		progressIntervals = numRules
	}
	step := numRules / progressIntervals
	collectedI := 0
	for i := 0; i < numRules; i += step {
		progress := float64(collectedI) / float64(numRules)
		nextI := i + step
		if nextI > numRules {
			nextI = numRules
		}
		rulesPartial := rules[i:nextI]
		go assessRulesC(rulesPartial, e, ic)
		numProcesses++

		if numProcesses >= maxProcesses {
			assessment, isError = getCOutcome(ic, ec, assessment, progress)
			if isError {
				return
			}
			collectedI += step
			numProcesses--
		}
	}

	for p := 0; p < numProcesses; p++ {
		progress := float64(collectedI) / float64(numRules)
		assessment, isError = getCOutcome(ic, ec, assessment, progress)
		if isError {
			return
		}
		collectedI += step
	}

	ec <- &AssessRulesMPOutcome{assessment, nil, 1.0, true}
	close(ec)
}

func getCOutcome(
	ic chan *assessRulesCOutcome,
	ec chan *AssessRulesMPOutcome,
	_assessment *assessment.Assessment,
	progress float64,
) (*assessment.Assessment, bool) {
	var retAssessment *assessment.Assessment
	var err error
	ec <- &AssessRulesMPOutcome{nil, nil, progress, false}
	assessmentOutcome := <-ic
	if assessmentOutcome.err != nil {
		ec <- &AssessRulesMPOutcome{nil, assessmentOutcome.err, progress, false}
		close(ec)
		return nil, true
	}
	if _assessment == nil {
		retAssessment = assessmentOutcome.assessment
	} else {
		retAssessment, err = _assessment.Merge(assessmentOutcome.assessment)
		if err != nil {
			ec <- &AssessRulesMPOutcome{nil, err, progress, false}
			close(ec)
			return nil, true
		}
	}
	return retAssessment, false
}

type assessRulesCOutcome struct {
	assessment *assessment.Assessment
	err        error
}

func assessRulesC(rules []*rule.Rule,
	experiment *experiment.Experiment,
	c chan *assessRulesCOutcome,
) {
	assessment, err := AssessRules(rules, experiment)
	c <- &assessRulesCOutcome{assessment, err}
}

func filterGoodRuleAssessors(
	ruleAssessments []*ruleassessor.RuleAssessor,
	numRecords int64,
) ([]*ruleassessor.RuleAssessor, error) {
	goodRuleAssessors := make([]*ruleassessor.RuleAssessor, 0)
	for _, ruleAssessment := range ruleAssessments {
		numMatches, exists :=
			ruleAssessment.GetAggregatorValue("numMatches", numRecords)
		if !exists {
			// TODO: Create a proper error for this?
			err := errors.New("numMatches doesn't exist in aggregators")
			return goodRuleAssessors, err
		}
		numMatchesInt, isInt := numMatches.Int()
		if !isInt {
			// TODO: Create a proper error for this?
			err := errors.New(fmt.Sprintf("Can't cast to Int: %q", numMatches))
			return goodRuleAssessors, err
		}
		if numMatchesInt > 0 {
			goodRuleAssessors = append(goodRuleAssessors, ruleAssessment)
		}
	}
	return goodRuleAssessors, nil
}

func processInput(input input.Input,
	ruleAssessors []*ruleassessor.RuleAssessor,
) (int64, error) {
	numRecords := int64(0)
	// TODO: test this rewinds properly
	if err := input.Rewind(); err != nil {
		return numRecords, err
	}

	for input.Next() {
		record, err := input.Read()
		if err != nil {
			return numRecords, err
		}
		numRecords++
		for _, ruleAssessor := range ruleAssessors {
			err := ruleAssessor.NextRecord(record)
			if err != nil {
				return numRecords, err
			}
		}
	}

	return numRecords, input.Err()
}

func addDefaultAggregators(
	_aggregators []aggregators.Aggregator,
) ([]aggregators.Aggregator, error) {
	newAggregators := make([]aggregators.Aggregator, 2)
	numMatchesAggregator, err := aggregators.New("numMatches", "count", "1==1")
	if err != nil {
		return newAggregators, err
	}
	percentMatchesAggregator, err :=
		aggregators.New("percentMatches", "calc",
			"roundto(100.0 * numMatches / numRecords, 2)")
	if err != nil {
		return newAggregators, err
	}
	goalsPassedScoreAggregator, err :=
		aggregators.New("numGoalsPassed", "goalspassedscore")
	if err != nil {
		return newAggregators, err
	}
	newAggregators[0] = numMatchesAggregator
	newAggregators[1] = percentMatchesAggregator
	newAggregators = append(newAggregators, _aggregators...)
	newAggregators = append(newAggregators, goalsPassedScoreAggregator)
	return newAggregators, nil
}
