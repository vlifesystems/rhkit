/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package rulehunter

import (
	"fmt"
	"github.com/lawrencewoodman/dlit_go"
	"sort"
	"strings"
)

func TweakRules(
	sortedRules []*Rule,
	fieldDescriptions map[string]*FieldDescription,
) []*Rule {
	numRulesPerGroup := 3
	groupedRules :=
		groupTweakableRules(sortedRules, numRulesPerGroup)
	return tweakRules(groupedRules, fieldDescriptions)
}

func groupTweakableRules(
	sortedRules []*Rule,
	numPerGroup int,
) map[string][]*Rule {
	groups := make(map[string][]*Rule)
	for _, rule := range sortedRules {
		isTweakable, fieldName, operator, _ := rule.getTweakableParts()
		if isTweakable {
			groupID := fmt.Sprintf("%s^%s", fieldName, operator)
			if len(groups[groupID]) < numPerGroup {
				groups[groupID] = append(groups[groupID], rule)
			}
		}
	}
	return groups
}

func tweakRules(
	groupedRules map[string][]*Rule,
	fieldDescriptions map[string]*FieldDescription,
) []*Rule {
	newRules := make([]*Rule, 1)
	newRules[0] = mustNewRule("true()")
	for _, rules := range groupedRules {
		firstRule := rules[0]
		comparisonPoints := makeComparisonPoints(rules, fieldDescriptions)
		for _, point := range comparisonPoints {
			newRule, err := firstRule.cloneWithValue(point)
			if err != nil {
				panic(fmt.Sprintf("Can't tweak rule: %s - %s", firstRule, err))
			}
			newRules = append(newRules, newRule)
		}
	}
	return newRules
}

func dlitInSlices(needle *dlit.Literal, haystacks ...[]*dlit.Literal) bool {
	for _, haystack := range haystacks {
		for _, v := range haystack {
			if needle.String() == v.String() {
				return true
			}
		}
	}
	return false
}

// TODO: Share similar code with generaters such as generateInt
func makeComparisonPoints(
	rules []*Rule,
	fieldDescriptions map[string]*FieldDescription,
) []string {
	var minInt int64
	var maxInt int64
	var minFloat float64
	var maxFloat float64
	var field string
	var tweakableValue string

	numbers := make([]*dlit.Literal, len(rules))
	newPoints := make([]*dlit.Literal, 0)
	for i, rule := range rules {
		_, field, _, tweakableValue = rule.getTweakableParts()
		numbers[i] = dlit.MustNew(tweakableValue)
	}

	numNumbers := len(numbers)
	sortNumbers(numbers)

	if fieldDescriptions[field].Kind == INT {
		for numI, numJ := 0, 1; numJ < numNumbers; numI, numJ = numI+1, numJ+1 {
			vI := numbers[numI]
			vJ := numbers[numJ]
			vIint, _ := vI.Int()
			vJint, _ := vJ.Int()
			if vIint < vJint {
				minInt = vIint
				maxInt = vJint
			} else {
				minInt = vJint
				maxInt = vIint
			}

			diff := maxInt - minInt
			step := diff / 10
			if diff < 10 {
				step = 1
			}

			for i := step; i < diff; i += step {
				newNum := dlit.MustNew(minInt + i)
				if !dlitInSlices(newNum, numbers, newPoints) {
					newPoints = append(newPoints, newNum)
				}
			}
		}
	} else {
		maxDP := fieldDescriptions[field].MaxDP
		for numI, numJ := 0, 1; numJ < numNumbers; numI, numJ = numI+1, numJ+1 {
			vI := numbers[numI]
			vJ := numbers[numJ]
			vIfloat, _ := vI.Float()
			vJfloat, _ := vJ.Float()
			if vIfloat < vJfloat {
				minFloat = vIfloat
				maxFloat = vJfloat
			} else {
				minFloat = vJfloat
				maxFloat = vIfloat
			}

			diff := maxFloat - minFloat
			step := diff / 10.0
			for i := step; i < diff; i += step {
				sum := minFloat + i
				for dp := maxDP; dp >= 0; dp-- {
					newNum := dlit.MustNew(floatReduceDP(sum, dp))
					if !dlitInSlices(newNum, numbers, newPoints) {
						newPoints = append(newPoints, newNum)
					}
				}
			}
		}
	}
	return arrayDlitsToStrings(newPoints)
}

func arrayDlitsToStrings(lits []*dlit.Literal) []string {
	r := make([]string, len(lits))
	for i, l := range lits {
		r[i] = l.String()
	}
	return r
}

func floatReduceDP(f float64, dp int) string {
	s := fmt.Sprintf("%.*f", dp, f)
	i := strings.IndexByte(s, '.')
	if i > -1 {
		s = strings.TrimRight(s, "0")
		i = strings.IndexByte(s, '.')
		if i == len(s)-1 {
			s = strings.TrimRight(s, ".")
		}
	}
	return s
}

// byNumber implements sort.Interface for []*dlit.Literal
type byNumber []*dlit.Literal

func (l byNumber) Len() int { return len(l) }
func (l byNumber) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l byNumber) Less(i, j int) bool {
	lI := l[i]
	lJ := l[j]
	iI, lIisInt := lI.Int()
	iJ, lJisInt := lJ.Int()
	if lIisInt && lJisInt {
		if iI < iJ {
			return true
		}
		return false
	}

	fI, lIisFloat := lI.Float()
	fJ, lJisFloat := lJ.Float()

	if lIisFloat && lJisFloat {
		if fI < fJ {
			return true
		}
		return false
	}
	panic(fmt.Sprintf("Can't compare numbers: %s, %s", lI, lJ))
}

func sortNumbers(nums []*dlit.Literal) {
	sort.Sort(byNumber(nums))
}