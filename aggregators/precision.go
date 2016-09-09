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

package aggregators

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/goal"
	"github.com/vlifesystems/rulehunter/internal/dexprfuncs"
)

type precisionAggregator struct{}

type precisionSpec struct {
	name string
	expr *dexpr.Expr
}

type precisionInstance struct {
	spec  *precisionSpec
	numTP int64
	numFP int64
}

var precisionExpr = dexpr.MustNew("roundto(numTP/(numTP+numFP),4)")

func init() {
	Register("precision", &precisionAggregator{})
}

func (a *precisionAggregator) MakeSpec(
	name string,
	expr string,
) (AggregatorSpec, error) {
	dexpr, err := dexpr.New(expr)
	if err != nil {
		return nil, err
	}
	d := &precisionSpec{
		name: name,
		expr: dexpr,
	}
	return d, nil
}

func (ad *precisionSpec) New() AggregatorInstance {
	return &precisionInstance{
		spec:  ad,
		numTP: 0,
		numFP: 0,
	}
}

func (ad *precisionSpec) GetName() string {
	return ad.name
}

func (ad *precisionSpec) GetArg() string {
	return ad.expr.String()
}

func (ai *precisionInstance) GetName() string {
	return ai.spec.name
}

func (ai *precisionInstance) NextRecord(record map[string]*dlit.Literal,
	isRuleTrue bool) error {
	matchExprIsTrue, err := ai.spec.expr.EvalBool(record, dexprfuncs.CallFuncs)
	if err != nil {
		return err
	}
	if isRuleTrue {
		if matchExprIsTrue {
			ai.numTP++
		} else {
			ai.numFP++
		}
	}
	return nil
}

func (ai *precisionInstance) GetResult(
	aggregatorInstances []AggregatorInstance,
	goals []*goal.Goal,
	numRecords int64,
) *dlit.Literal {
	if ai.numTP == 0 && ai.numFP == 0 {
		return dlit.MustNew(0)
	}

	vars := map[string]*dlit.Literal{
		"numTP": dlit.MustNew(ai.numTP),
		"numFP": dlit.MustNew(ai.numFP),
	}
	return precisionExpr.Eval(vars, dexprfuncs.CallFuncs)
}
