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
)

// BetweenFVI represents a rule determining if:
// field >= intValue && field <= intValue
type BetweenFVI struct {
	field string
	min   int64
	max   int64
}

func NewBetweenFVI(field string, min int64, max int64) (Rule, error) {
	if max <= min {
		return nil,
			fmt.Errorf("can't create Between rule where max: %d <= min: %d", max, min)
	}
	return &BetweenFVI{field: field, min: min, max: max}, nil
}

func MustNewBetweenFVI(field string, min int64, max int64) Rule {
	r, err := NewBetweenFVI(field, min, max)
	if err != nil {
		panic(err)
	}
	return r
}

func (r *BetweenFVI) String() string {
	return fmt.Sprintf("%s >= %d && %s <= %d", r.field, r.min, r.field, r.max)
}

func (r *BetweenFVI) IsTrue(record ddataset.Record) (bool, error) {
	value, ok := record[r.field]
	if !ok {
		return false, InvalidRuleError{Rule: r}
	}

	valueInt, valueIsInt := value.Int()
	if valueIsInt {
		return valueInt >= r.min && valueInt <= r.max, nil
	}

	return false, IncompatibleTypesRuleError{Rule: r}
}

func (r *BetweenFVI) GetFields() []string {
	return []string{r.field}
}
