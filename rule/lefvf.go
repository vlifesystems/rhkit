/*
	Copyright (C) 2016 vLife Systems Ltd <http://vlifesystems.com>
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
	"fmt"
	"github.com/lawrencewoodman/ddataset"
	"strconv"
)

// LEFVF represents a rule determening if fieldA <= floatValue
type LEFVF struct {
	field string
	value float64
}

func NewLEFVF(field string, value float64) TweakableRule {
	return &LEFVF{field: field, value: value}
}

func (r *LEFVF) String() string {
	return r.field + " <= " + strconv.FormatFloat(r.value, 'f', -1, 64)
}

func (r *LEFVF) GetTweakableParts() (string, string, string) {
	return r.field, "<=", strconv.FormatFloat(r.value, 'f', -1, 64)
}

func (r *LEFVF) IsTrue(record ddataset.Record) (bool, error) {
	lh, ok := record[r.field]
	if !ok {
		return false, InvalidRuleError{Rule: r}
	}

	lhFloat, lhIsFloat := lh.Float()
	if lhIsFloat {
		return lhFloat <= r.value, nil
	}

	return false, IncompatibleTypesRuleError{Rule: r}
}

func (r *LEFVF) CloneWithValue(newValue interface{}) TweakableRule {
	f, ok := newValue.(float64)
	if ok {
		return NewLEFVF(r.field, f)
	}
	panic(fmt.Sprintf(
		"can't clone with newValue: %v of type %T, need type float64",
		newValue, newValue,
	))
}