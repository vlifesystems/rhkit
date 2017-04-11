/*
	Copyright (C) 2017 vLife Systems Ltd <http://vlifesystems.com>
	This file is part of rhkit.

	rhkit is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	rhkit is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with rhkit; see the file COPYING.  If not, see
	<http://www.gnu.org/licenses/>.
*/

package rule

import (
	"github.com/lawrencewoodman/ddataset"
	"strconv"
)

// AddGEF represents a rule determining if fieldA + fieldB >= floatValue
type AddGEF struct {
	fieldA string
	fieldB string
	value  float64
}

func NewAddGEF(fieldA string, fieldB string, value float64) *AddGEF {
	return &AddGEF{fieldA: fieldA, fieldB: fieldB, value: value}
}

func (r *AddGEF) String() string {
	return r.fieldA + " + " + r.fieldB + " >= " +
		strconv.FormatFloat(r.value, 'f', -1, 64)
}

func (r *AddGEF) GetValue() float64 {
	return r.value
}

func (r *AddGEF) GetFields() []string {
	return []string{r.fieldA, r.fieldB}
}

// IsTrue returns whether the rule is true for this record.
// This rule relies on making shure that the two fields when
// added will not overflow, so this must have been checked
// before hand by looking at their max/min in the input description.
func (r *AddGEF) IsTrue(record ddataset.Record) (bool, error) {
	vA, ok := record[r.fieldA]
	if !ok {
		return false, InvalidRuleError{Rule: r}
	}

	vB, ok := record[r.fieldB]
	if !ok {
		return false, InvalidRuleError{Rule: r}
	}

	vAFloat, vAIsFloat := vA.Float()
	if !vAIsFloat {
		return false, IncompatibleTypesRuleError{Rule: r}
	}
	vBFloat, vBIsFloat := vB.Float()
	if !vBIsFloat {
		return false, IncompatibleTypesRuleError{Rule: r}
	}

	return vAFloat+vBFloat >= r.value, nil
}

// TODO: implement this by passing inputDescription
/*
func (r *AddGEF) Tweak(
	min *dlit.Literal,
	max *dlit.Literal,
	maxDP int,
	stage int,
) []Rule {
	rules := make([]Rule, 0)
	minFloat, _ := min.Float()
	maxFloat, _ := max.Float()
	step := (maxFloat - minFloat) / (10 * float64(stage))
	low := r.value - step
	high := r.value + step
	interStep := (high - low) / 20
	for n := low; n <= high; n += interStep {
		v := truncateFloat(n, maxDP)
		if v != r.value && v != low && v != high && v >= minFloat && v <= maxFloat {
			r := NewAddGEF(r.fieldA, r.fieldB, truncateFloat(n, maxDP))
			rules = append(rules, r)
		}
	}
	return rules
}
*/

func (r *AddGEF) Overlaps(o Rule) bool {
	switch x := o.(type) {
	case *AddGEF:
		oFields := x.GetFields()
		if r.fieldA == oFields[0] && r.fieldB == oFields[1] {
			return true
		}
	}
	return false
}
