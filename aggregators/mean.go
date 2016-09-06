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
	"fmt"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/goal"
	"github.com/vlifesystems/rulehunter/internal/dexprfuncs"
)

type meanAggregator struct{}

type meanSpec struct {
	name string
	expr *dexpr.Expr
}

type meanInstance struct {
	spec       *meanSpec
	sum        *dlit.Literal
	numRecords int64
}

var meanExpr = dexpr.MustNew("sum/n")
var meanSumExpr = dexpr.MustNew("sum+value")

func init() {
	Register("mean", &meanAggregator{})
}

func (a *meanAggregator) MakeSpec(
	name string,
	expr string,
) (AggregatorSpec, error) {
	dexpr, err := dexpr.New(expr)
	if err != nil {
		return nil, err
	}
	d := &meanSpec{
		name: name,
		expr: dexpr,
	}
	return d, nil
}

func (ad *meanSpec) New() AggregatorInstance {
	return &meanInstance{
		spec:       ad,
		sum:        dlit.MustNew(0),
		numRecords: 0,
	}
}

func (ad *meanSpec) GetName() string {
	return ad.name
}

func (ad *meanSpec) GetArg() string {
	return ad.expr.String()
}

func (ai *meanInstance) GetName() string {
	return ai.spec.name
}

func (ai *meanInstance) NextRecord(
	record map[string]*dlit.Literal,
	isRuleTrue bool,
) error {
	if isRuleTrue {
		ai.numRecords++
		exprValue := ai.spec.expr.Eval(record, dexprfuncs.CallFuncs)
		_, valueIsFloat := exprValue.Float()
		if !valueIsFloat {
			return fmt.Errorf("mean aggregator: value isn't a float: %s", exprValue)
		}

		vars := map[string]*dlit.Literal{
			"sum":   ai.sum,
			"value": exprValue,
		}
		ai.sum = meanSumExpr.Eval(vars, dexprfuncs.CallFuncs)
	}
	return nil
}

func (ai *meanInstance) GetResult(
	aggregatorInstances []AggregatorInstance,
	goals []*goal.Goal,
	numRecords int64,
) *dlit.Literal {
	if ai.numRecords == 0 {
		return dlit.MustNew(0)
	}

	vars := map[string]*dlit.Literal{
		"sum": ai.sum,
		"n":   dlit.MustNew(ai.numRecords),
	}
	return meanExpr.Eval(vars, dexprfuncs.CallFuncs)
}
