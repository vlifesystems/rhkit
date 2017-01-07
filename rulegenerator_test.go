package rhkit

import (
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/rule"
	"regexp"
	"sort"
	"testing"
)

func TestGenerateRules_1(t *testing.T) {
	testPurpose := "Ensure generates correct rules for each field"
	inputDescription := &Description{
		map[string]*fieldDescription{
			"team": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"a": valueDescription{dlit.NewString("a"), 3},
					"b": valueDescription{dlit.NewString("b"), 3},
					"c": valueDescription{dlit.NewString("c"), 3},
				},
			},
			"teamOut": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"a": valueDescription{dlit.NewString("a"), 3},
					"c": valueDescription{dlit.NewString("c"), 1},
					"d": valueDescription{dlit.NewString("d"), 3},
					"e": valueDescription{dlit.NewString("e"), 3},
					"f": valueDescription{dlit.NewString("f"), 3},
				},
			},
			"teamBob": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"a": valueDescription{dlit.NewString("a"), 3},
					"b": valueDescription{dlit.NewString("b"), 3},
					"c": valueDescription{dlit.NewString("c"), 3},
				},
			},
			"camp": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"arthur":  valueDescription{dlit.NewString("arthur"), 3},
					"offa":    valueDescription{dlit.NewString("offa"), 3},
					"richard": valueDescription{dlit.NewString("richard"), 3},
					"owen":    valueDescription{dlit.NewString("owen"), 3},
				},
			},
			"level": &fieldDescription{
				kind:  ftInt,
				min:   dlit.MustNew(0),
				max:   dlit.MustNew(5),
				maxDP: 0,
				values: map[string]valueDescription{
					"0": valueDescription{dlit.NewString("0"), 3},
					"1": valueDescription{dlit.NewString("1"), 3},
					"2": valueDescription{dlit.NewString("2"), 1},
					"3": valueDescription{dlit.NewString("3"), 3},
					"4": valueDescription{dlit.NewString("4"), 3},
					"5": valueDescription{dlit.NewString("5"), 3},
				},
			},
			"levelBob": &fieldDescription{
				kind:  ftInt,
				min:   dlit.MustNew(0),
				max:   dlit.MustNew(5),
				maxDP: 0,
				values: map[string]valueDescription{
					"0": valueDescription{dlit.NewString("0"), 3},
					"1": valueDescription{dlit.NewString("1"), 3},
					"2": valueDescription{dlit.NewString("2"), 3},
					"3": valueDescription{dlit.NewString("3"), 3},
					"5": valueDescription{dlit.NewString("5"), 3},
				},
			},
			"flowA": &fieldDescription{
				kind:  ftFloat,
				min:   dlit.MustNew(0),
				max:   dlit.MustNew(10.5),
				maxDP: 2,
				values: map[string]valueDescription{
					"0.0":  valueDescription{dlit.NewString("0.0"), 3},
					"2.34": valueDescription{dlit.NewString("2.34"), 3},
					"10.5": valueDescription{dlit.NewString("10.5"), 3},
				},
			},
			"flowB": &fieldDescription{
				kind:  ftFloat,
				min:   dlit.MustNew(0),
				max:   dlit.MustNew(10.5),
				maxDP: 2,
				values: map[string]valueDescription{
					"0.0":  valueDescription{dlit.NewString("0.0"), 3},
					"2.34": valueDescription{dlit.NewString("2.34"), 3},
					"2.44": valueDescription{dlit.NewString("2.44"), 1},
					"10.5": valueDescription{dlit.NewString("10.5"), 3},
				},
			},
			"position": &fieldDescription{
				kind:  ftInt,
				min:   dlit.MustNew(1),
				max:   dlit.MustNew(13),
				maxDP: 0,
				values: map[string]valueDescription{
					"1":  valueDescription{dlit.NewString("1"), 3},
					"2":  valueDescription{dlit.NewString("2"), 3},
					"3":  valueDescription{dlit.NewString("3"), 3},
					"4":  valueDescription{dlit.NewString("4"), 3},
					"5":  valueDescription{dlit.NewString("5"), 3},
					"6":  valueDescription{dlit.NewString("6"), 3},
					"7":  valueDescription{dlit.NewString("7"), 3},
					"8":  valueDescription{dlit.NewString("8"), 3},
					"9":  valueDescription{dlit.NewString("9"), 3},
					"10": valueDescription{dlit.NewString("10"), 3},
					"11": valueDescription{dlit.NewString("11"), 3},
					"12": valueDescription{dlit.NewString("12"), 3},
					"13": valueDescription{dlit.NewString("13"), 3},
				},
			},
		}}
	ruleFields :=
		[]string{"team", "teamOut", "camp", "level", "flowA", "flowB", "position"}
	cases := []struct {
		field     string
		wantRules []rule.Rule
	}{
		{"team", []rule.Rule{
			rule.NewEQFVS("team", "a"),
			rule.NewEQFVS("team", "b"), rule.NewEQFVS("team", "c"),
			rule.NewNEFVS("team", "a"),
			rule.NewNEFVS("team", "b"),
			rule.NewNEFVS("team", "c"),
			rule.NewEQFF("team", "teamOut"),
			rule.NewNEFF("team", "teamOut"),
		}},
		{"teamOut", []rule.Rule{
			rule.NewEQFVS("teamOut", "a"),
			rule.NewEQFVS("teamOut", "d"),
			rule.NewEQFVS("teamOut", "e"),
			rule.NewEQFVS("teamOut", "f"),
			rule.NewNEFVS("teamOut", "a"),
			rule.NewNEFVS("teamOut", "d"),
			rule.NewNEFVS("teamOut", "e"),
			rule.NewNEFVS("teamOut", "f"),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "d")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "e")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "f")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("d", "e")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("d", "f")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("e", "f")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "d", "e")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "d", "f")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "e", "f")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("d", "e", "f")),
		}},
		{"level", []rule.Rule{
			rule.NewEQFVI("level", 0),
			rule.NewEQFVI("level", 1),
			rule.NewEQFVI("level", 3),
			rule.NewEQFVI("level", 4),
			rule.NewEQFVI("level", 5),
			rule.NewNEFVI("level", 0),
			rule.NewNEFVI("level", 1),
			rule.NewNEFVI("level", 3),
			rule.NewNEFVI("level", 4),
			rule.NewNEFVI("level", 5),
			rule.NewLTFF("level", "position"),
			rule.NewLEFF("level", "position"),
			rule.NewNEFF("level", "position"),
			rule.NewGEFF("level", "position"),
			rule.NewGTFF("level", "position"),
			rule.NewEQFF("level", "position"),
			rule.NewGEFVI("level", 0),
			rule.NewGEFVI("level", 1),
			rule.NewGEFVI("level", 2),
			rule.NewGEFVI("level", 3),
			rule.NewGEFVI("level", 4),
			rule.NewLEFVI("level", 1),
			rule.NewLEFVI("level", 2),
			rule.NewLEFVI("level", 3),
			rule.NewLEFVI("level", 4),
			rule.NewLEFVI("level", 5),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "3")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "3")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("3", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("3", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "3")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "3", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "3", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "3", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "3", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("3", "4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "3", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "3", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "3", "4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "3", "4", "5")),
		}},
		{"flowA", []rule.Rule{
			rule.NewEQFVI("flowA", 0),
			rule.NewEQFVF("flowA", 2.34),
			rule.NewEQFVF("flowA", 10.5),
			rule.NewNEFVI("flowA", 0),
			rule.NewNEFVF("flowA", 2.34),
			rule.NewNEFVF("flowA", 10.5),
			rule.NewLTFF("flowA", "level"),
			rule.NewLEFF("flowA", "level"),
			rule.NewNEFF("flowA", "level"),
			rule.NewGEFF("flowA", "level"),
			rule.NewGTFF("flowA", "level"),
			rule.NewEQFF("flowA", "level"),
			rule.NewLTFF("flowA", "flowB"),
			rule.NewLEFF("flowA", "flowB"),
			rule.NewNEFF("flowA", "flowB"),
			rule.NewGEFF("flowA", "flowB"),
			rule.NewGTFF("flowA", "flowB"),
			rule.NewEQFF("flowA", "flowB"),
			rule.NewLTFF("flowA", "position"),
			rule.NewLEFF("flowA", "position"),
			rule.NewNEFF("flowA", "position"),
			rule.NewGEFF("flowA", "position"),
			rule.NewGTFF("flowA", "position"),
			rule.NewEQFF("flowA", "position"),
			rule.NewGEFVF("flowA", 0.0),
			rule.NewGEFVF("flowA", 1.05),
			rule.NewGEFVF("flowA", 2.1),
			rule.NewGEFVF("flowA", 3.15),
			rule.NewGEFVF("flowA", 4.2),
			rule.NewGEFVF("flowA", 5.25),
			rule.NewGEFVF("flowA", 6.3),
			rule.NewGEFVF("flowA", 7.35),
			rule.NewGEFVF("flowA", 8.4),
			rule.NewGEFVF("flowA", 9.45),
			rule.NewLEFVF("flowA", 1.05),
			rule.NewLEFVF("flowA", 2.1),
			rule.NewLEFVF("flowA", 3.15),
			rule.NewLEFVF("flowA", 4.2),
			rule.NewLEFVF("flowA", 5.25),
			rule.NewLEFVF("flowA", 6.3),
			rule.NewLEFVF("flowA", 7.35),
			rule.NewLEFVF("flowA", 8.4),
			rule.NewLEFVF("flowA", 9.45),
		}},
		{"flowB", []rule.Rule{
			rule.NewEQFVI("flowB", 0),
			rule.NewEQFVF("flowB", 2.34),
			rule.NewEQFVF("flowB", 10.5),
			rule.NewNEFVI("flowB", 0),
			rule.NewNEFVF("flowB", 2.34),
			rule.NewNEFVF("flowB", 10.5),
			rule.NewLTFF("flowB", "level"),
			rule.NewLEFF("flowB", "level"),
			rule.NewNEFF("flowB", "level"),
			rule.NewGEFF("flowB", "level"),
			rule.NewGTFF("flowB", "level"),
			rule.NewEQFF("flowB", "level"),
			rule.NewLTFF("flowB", "position"),
			rule.NewLEFF("flowB", "position"),
			rule.NewNEFF("flowB", "position"),
			rule.NewGEFF("flowB", "position"),
			rule.NewGTFF("flowB", "position"),
			rule.NewEQFF("flowB", "position"),
			rule.NewGEFVF("flowB", 0.0),
			rule.NewGEFVF("flowB", 1.05),
			rule.NewGEFVF("flowB", 2.1),
			rule.NewGEFVF("flowB", 3.15),
			rule.NewGEFVF("flowB", 4.2),
			rule.NewGEFVF("flowB", 5.25),
			rule.NewGEFVF("flowB", 6.3),
			rule.NewGEFVF("flowB", 7.35),
			rule.NewGEFVF("flowB", 8.4),
			rule.NewGEFVF("flowB", 9.45),
			rule.NewLEFVF("flowB", 1.05),
			rule.NewLEFVF("flowB", 2.1),
			rule.NewLEFVF("flowB", 3.15),
			rule.NewLEFVF("flowB", 4.2),
			rule.NewLEFVF("flowB", 5.25),
			rule.NewLEFVF("flowB", 6.3),
			rule.NewLEFVF("flowB", 7.35),
			rule.NewLEFVF("flowB", 8.4),
			rule.NewLEFVF("flowB", 9.45),
			rule.NewInFV("flowB", makeStringsDlitSlice("0.0", "2.34")),
			rule.NewInFV("flowB", makeStringsDlitSlice("0.0", "10.5")),
			rule.NewInFV("flowB", makeStringsDlitSlice("10.5", "2.34")),
		}},
		{"position", []rule.Rule{
			rule.NewEQFVI("position", 1),
			rule.NewEQFVI("position", 2),
			rule.NewEQFVI("position", 3),
			rule.NewEQFVI("position", 4),
			rule.NewEQFVI("position", 5),
			rule.NewEQFVI("position", 6),
			rule.NewEQFVI("position", 7),
			rule.NewEQFVI("position", 8),
			rule.NewEQFVI("position", 9),
			rule.NewEQFVI("position", 10),
			rule.NewEQFVI("position", 11),
			rule.NewEQFVI("position", 12),
			rule.NewEQFVI("position", 13),
			rule.NewNEFVI("position", 1),
			rule.NewNEFVI("position", 2),
			rule.NewNEFVI("position", 3),
			rule.NewNEFVI("position", 4),
			rule.NewNEFVI("position", 5),
			rule.NewNEFVI("position", 6),
			rule.NewNEFVI("position", 7),
			rule.NewNEFVI("position", 8),
			rule.NewNEFVI("position", 9),
			rule.NewNEFVI("position", 10),
			rule.NewNEFVI("position", 11),
			rule.NewNEFVI("position", 12),
			rule.NewNEFVI("position", 13),
			rule.NewGEFVI("position", 1),
			rule.NewGEFVI("position", 2),
			rule.NewGEFVI("position", 3),
			rule.NewGEFVI("position", 4),
			rule.NewGEFVI("position", 5),
			rule.NewGEFVI("position", 6),
			rule.NewGEFVI("position", 7),
			rule.NewGEFVI("position", 8),
			rule.NewGEFVI("position", 9),
			rule.NewGEFVI("position", 10),
			rule.NewGEFVI("position", 11),
			rule.NewGEFVI("position", 12),
			rule.NewLEFVI("position", 2),
			rule.NewLEFVI("position", 3),
			rule.NewLEFVI("position", 4),
			rule.NewLEFVI("position", 5),
			rule.NewLEFVI("position", 6),
			rule.NewLEFVI("position", 7),
			rule.NewLEFVI("position", 8),
			rule.NewLEFVI("position", 9),
			rule.NewLEFVI("position", 10),
			rule.NewLEFVI("position", 11),
			rule.NewLEFVI("position", 12),
			rule.NewLEFVI("position", 13),
		}},
	}

	rules, err := GenerateRules(inputDescription, ruleFields)
	if err != nil {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("GenerateRules(%v, %v) err: %v", inputDescription, ruleFields, err)
	}

	for _, c := range cases {
		gotFieldRules := getFieldRules(c.field, rules)
		rulesMatch, msg := matchRulesUnordered(gotFieldRules, c.wantRules)
		if !rulesMatch {
			gotFieldRuleStrs := rulesToSortedStrings(gotFieldRules)
			wantRuleStrs := rulesToSortedStrings(c.wantRules)
			t.Errorf("Test: %s\n", testPurpose)
			t.Errorf("matchRulesUnordered() rules don't match for field: %s - %s\ngot: %s\nwant: %s\n",
				c.field, msg, gotFieldRuleStrs, wantRuleStrs)
		}
	}
}

func TestGenerateRules_2(t *testing.T) {
	testPurpose := "Ensure generates a True rule"
	inputDescription := &Description{
		map[string]*fieldDescription{
			"team": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"a": valueDescription{dlit.MustNew("a"), 3},
					"b": valueDescription{dlit.MustNew("b"), 3},
					"c": valueDescription{dlit.MustNew("c"), 3},
				},
			},
			"teamOut": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"a": valueDescription{dlit.MustNew("a"), 3},
					"c": valueDescription{dlit.MustNew("c"), 3},
					"d": valueDescription{dlit.MustNew("d"), 3},
					"e": valueDescription{dlit.MustNew("e"), 3},
					"f": valueDescription{dlit.MustNew("f"), 3},
				},
			},
		}}
	ruleFields := []string{"team", "teamOut"}

	rules, err := GenerateRules(inputDescription, ruleFields)
	if err != nil {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("GenerateRules(%v, %v) err: %v", inputDescription, ruleFields, err)
	}

	trueRuleFound := false
	for _, r := range rules {
		if _, isTrueRule := r.(rule.True); isTrueRule {
			trueRuleFound = true
			break
		}
	}
	if !trueRuleFound {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("GenerateRules(%v, %v)  - True rule missing",
			inputDescription, ruleFields)
	}
}

func TestGenerateRules_3(t *testing.T) {
	testPurpose := "Ensure generates correct combination rules"
	inputDescription := &Description{
		map[string]*fieldDescription{
			"team": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"a": valueDescription{dlit.MustNew("a"), 3},
					"b": valueDescription{dlit.MustNew("b"), 3},
					"c": valueDescription{dlit.MustNew("c"), 3},
				},
			},
			"teamAlt": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"a": valueDescription{dlit.MustNew("a"), 3},
					"b": valueDescription{dlit.MustNew("b"), 3},
					"c": valueDescription{dlit.MustNew("c"), 3},
					"d": valueDescription{dlit.MustNew("d"), 3},
				},
			},
			"monthIn": &fieldDescription{
				kind: ftInt,
				min:  dlit.MustNew(1),
				max:  dlit.MustNew(3),
				values: map[string]valueDescription{
					"1": valueDescription{dlit.MustNew(1), 3},
					"2": valueDescription{dlit.MustNew(2), 3},
					"3": valueDescription{dlit.MustNew(3), 3},
				},
			},
			"monthOut": &fieldDescription{
				kind: ftInt,
				min:  dlit.MustNew(1),
				max:  dlit.MustNew(3),
				values: map[string]valueDescription{
					"1": valueDescription{dlit.MustNew(1), 3},
					"2": valueDescription{dlit.MustNew(2), 3},
					"3": valueDescription{dlit.MustNew(3), 3},
					"4": valueDescription{dlit.MustNew(4), 3},
				},
			},
			"win": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"t": valueDescription{dlit.MustNew("t"), 3},
					"f": valueDescription{dlit.MustNew("f"), 3},
				},
			},
		}}
	cases := []struct {
		field      string
		ruleFields []string
		wantRules  []rule.Rule
	}{
		{field: "team",
			ruleFields: []string{"team", "monthIn"},
			wantRules: []rule.Rule{
				rule.NewEQFVS("team", "a"),
				rule.NewEQFVS("team", "b"),
				rule.NewEQFVS("team", "c"),
				rule.NewNEFVS("team", "a"),
				rule.NewNEFVS("team", "b"),
				rule.NewNEFVS("team", "c"),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewNEFVI("monthIn", 3)),
			}},
		{field: "team",
			ruleFields: []string{"team", "monthOut"},
			wantRules: []rule.Rule{
				rule.NewEQFVS("team", "a"),
				rule.NewEQFVS("team", "b"),
				rule.NewEQFVS("team", "c"),
				rule.NewNEFVS("team", "a"),
				rule.NewNEFVS("team", "b"),
				rule.NewNEFVS("team", "c"),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewEQFVI("monthOut", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewEQFVI("monthOut", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewEQFVI("monthOut", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewEQFVI("monthOut", 4)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewNEFVI("monthOut", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewNEFVI("monthOut", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewNEFVI("monthOut", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewNEFVI("monthOut", 4)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewEQFVI("monthOut", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewEQFVI("monthOut", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewEQFVI("monthOut", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewEQFVI("monthOut", 4)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewNEFVI("monthOut", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewNEFVI("monthOut", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewNEFVI("monthOut", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewNEFVI("monthOut", 4)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewEQFVI("monthOut", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewEQFVI("monthOut", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewEQFVI("monthOut", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewEQFVI("monthOut", 4)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewNEFVI("monthOut", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewNEFVI("monthOut", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewNEFVI("monthOut", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewNEFVI("monthOut", 4)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewEQFVI("monthOut", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewEQFVI("monthOut", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewEQFVI("monthOut", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewEQFVI("monthOut", 4)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewNEFVI("monthOut", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewNEFVI("monthOut", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewNEFVI("monthOut", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewNEFVI("monthOut", 4)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewEQFVI("monthOut", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewEQFVI("monthOut", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewEQFVI("monthOut", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewEQFVI("monthOut", 4)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewNEFVI("monthOut", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewNEFVI("monthOut", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewNEFVI("monthOut", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewNEFVI("monthOut", 4)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewEQFVI("monthOut", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewEQFVI("monthOut", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewEQFVI("monthOut", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewEQFVI("monthOut", 4)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewNEFVI("monthOut", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewNEFVI("monthOut", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewNEFVI("monthOut", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewNEFVI("monthOut", 4)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "2"))),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "3"))),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "3"))),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("3", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "2"))),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "3"))),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "3"))),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("3", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "2"))),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "3"))),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "3"))),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("3", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "2"))),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "3"))),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "3"))),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("3", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "2"))),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "3"))),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "3"))),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("3", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "2"))),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "3"))),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "3"))),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("3", "4"))),
			}},
		{field: "teamAlt",
			ruleFields: []string{"teamAlt", "monthIn"},
			wantRules: []rule.Rule{
				rule.NewEQFVS("teamAlt", "a"),
				rule.NewEQFVS("teamAlt", "b"),
				rule.NewEQFVS("teamAlt", "c"),
				rule.NewEQFVS("teamAlt", "d"),
				rule.NewNEFVS("teamAlt", "a"),
				rule.NewNEFVS("teamAlt", "b"),
				rule.NewNEFVS("teamAlt", "c"),
				rule.NewNEFVS("teamAlt", "d"),
				rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "b")),
				rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "c")),
				rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "d")),
				rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "c")),
				rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "d")),
				rule.NewInFV("teamAlt", makeStringsDlitSlice("c", "d")),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "a"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "a"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "a"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "a"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "a"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "a"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "b"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "b"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "b"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "b"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "b"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "b"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "c"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "c"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "c"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "c"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "c"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "c"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "d"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "d"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "d"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "d"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "d"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "d"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "a"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "a"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "a"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "a"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "a"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "a"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "b"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "b"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "b"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "b"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "b"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "b"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "c"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "c"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "c"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "c"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "c"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "c"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "d"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "d"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "d"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "d"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "d"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "d"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "b")),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "b")),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "b")),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "b")),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "b")),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "b")),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "c")),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "c")),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "c")),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "c")),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "c")),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "c")),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "d")),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "d")),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "d")),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "d")),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "d")),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "d")),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "c")),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "c")),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "c")),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "c")),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "c")),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "c")),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "d")),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "d")),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "d")),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "d")),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "d")),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "d")),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("c", "d")),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("c", "d")),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("c", "d")),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("c", "d")),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("c", "d")),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("c", "d")),
					rule.NewNEFVI("monthIn", 3)),
			}},
	}

	for _, c := range cases {
		rules, err := GenerateRules(inputDescription, c.ruleFields)
		if err != nil {
			t.Errorf("Test: %s\n", testPurpose)
			t.Fatalf("GenerateRules(%v, %v) err: %v",
				inputDescription, c.ruleFields, err)
		}

		gotFieldRules := getFieldRules(c.field, rules)
		rulesMatch, msg := matchRulesUnordered(gotFieldRules, c.wantRules)
		if !rulesMatch {
			gotFieldRuleStrs := rulesToSortedStrings(gotFieldRules)
			wantRuleStrs := rulesToSortedStrings(c.wantRules)
			t.Errorf("Test: %s\n", testPurpose)
			t.Errorf("matchRulesUnordered() rules don't match for field: %s - %s\ngot: %s\nwant: %s\n",
				c.field, msg, gotFieldRuleStrs, wantRuleStrs)
		}
	}
}

func TestCombinedRules(t *testing.T) {
	cases := []struct {
		inRules   []rule.Rule
		wantRules []rule.Rule
	}{
		{[]rule.Rule{
			rule.NewEQFVS("team", "a"),
			rule.NewGEFVI("band", 4),
			rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
		},
			[]rule.Rule{
				rule.NewAnd(rule.NewEQFVS("team", "a"), rule.NewGEFVI("band", 4)),
				rule.NewOr(rule.NewEQFVS("team", "a"), rule.NewGEFVI("band", 4)),
				rule.NewAnd(
					rule.NewEQFVS("team", "a"),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
				rule.NewOr(
					rule.NewEQFVS("team", "a"),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
				rule.NewAnd(
					rule.NewGEFVI("band", 4),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
				rule.NewOr(
					rule.NewGEFVI("band", 4),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
			}},
		{[]rule.Rule{rule.NewEQFVS("team", "a")}, []rule.Rule{}},
		{[]rule.Rule{}, []rule.Rule{}},
	}

	for _, c := range cases {
		gotRules := CombineRules(c.inRules)
		rulesMatch, msg := matchRulesUnordered(gotRules, c.wantRules)
		if !rulesMatch {
			gotRuleStrs := rulesToSortedStrings(gotRules)
			wantRuleStrs := rulesToSortedStrings(c.wantRules)
			t.Errorf("matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
				msg, gotRuleStrs, wantRuleStrs)
		}
	}
}

func TestMakeCompareValues(t *testing.T) {
	values1 := map[string]valueDescription{
		"a": valueDescription{dlit.MustNew("a"), 2},
		"c": valueDescription{dlit.MustNew("c"), 2},
		"d": valueDescription{dlit.MustNew("d"), 2},
		"e": valueDescription{dlit.MustNew("e"), 2},
		"f": valueDescription{dlit.MustNew("f"), 2},
	}
	values2 := map[string]valueDescription{
		"a": valueDescription{dlit.MustNew("a"), 2},
		"c": valueDescription{dlit.MustNew("c"), 1},
		"d": valueDescription{dlit.MustNew("d"), 2},
		"e": valueDescription{dlit.MustNew("e"), 2},
		"f": valueDescription{dlit.MustNew("f"), 2},
	}
	cases := []struct {
		values map[string]valueDescription
		i      int
		want   []*dlit.Literal
	}{
		{
			values: values1,
			i:      2,
			want:   []*dlit.Literal{dlit.NewString("c")},
		},
		{
			values: values2,
			i:      2,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      3,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("c")},
		},
		{
			values: values2,
			i:      3,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      4,
			want:   []*dlit.Literal{dlit.NewString("d")},
		},
		{
			values: values2,
			i:      4,
			want:   []*dlit.Literal{dlit.NewString("d")},
		},
		{
			values: values1,
			i:      5,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("d")},
		},
		{
			values: values2,
			i:      5,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("d")},
		},
		{
			values: values1,
			i:      6,
			want:   []*dlit.Literal{dlit.NewString("c"), dlit.NewString("d")},
		},
		{
			values: values2,
			i:      6,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      7,
			want: []*dlit.Literal{
				dlit.NewString("a"),
				dlit.NewString("c"),
				dlit.NewString("d"),
			},
		},
		{
			values: values2,
			i:      7,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      8,
			want:   []*dlit.Literal{dlit.NewString("e")},
		},
		{
			values: values2,
			i:      8,
			want:   []*dlit.Literal{dlit.NewString("e")},
		},
		{
			values: values1,
			i:      9,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("e")},
		},
		{
			values: values2,
			i:      9,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("e")},
		},
		{
			values: values1,
			i:      10,
			want:   []*dlit.Literal{dlit.NewString("c"), dlit.NewString("e")},
		},
		{
			values: values2,
			i:      10,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      11,
			want: []*dlit.Literal{
				dlit.NewString("a"),
				dlit.NewString("c"),
				dlit.NewString("e"),
			},
		},
		{
			values: values2,
			i:      11,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      12,
			want:   []*dlit.Literal{dlit.NewString("d"), dlit.NewString("e")},
		},
		{
			values: values2,
			i:      12,
			want:   []*dlit.Literal{dlit.NewString("d"), dlit.NewString("e")},
		},
		{
			values: values1,
			i:      13,
			want: []*dlit.Literal{
				dlit.NewString("a"),
				dlit.NewString("d"),
				dlit.NewString("e"),
			},
		},
		{
			values: values1,
			i:      14,
			want: []*dlit.Literal{
				dlit.NewString("c"),
				dlit.NewString("d"),
				dlit.NewString("e"),
			},
		},
		{
			values: values1,
			i:      15,
			want: []*dlit.Literal{
				dlit.NewString("a"),
				dlit.NewString("c"),
				dlit.NewString("d"),
				dlit.NewString("e"),
			},
		},
		{
			values: values1,
			i:      16,
			want:   []*dlit.Literal{dlit.NewString("f")},
		},
		{
			values: values2,
			i:      16,
			want:   []*dlit.Literal{dlit.NewString("f")},
		},
		{
			values: values1,
			i:      17,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("f")},
		},
		{
			values: values2,
			i:      17,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("f")},
		},
	}
	for _, c := range cases {
		got := makeCompareValues(c.values, c.i)
		if len(got) != len(c.want) {
			t.Errorf("makeCompareValues(%s, %d) got: %v, want: %v",
				c.values, c.i, got, c.want)
		}
		for j, v := range got {
			o := c.want[j]
			if o.String() != v.String() {
				t.Errorf("makeCompareValues(%s, %d) got: %v, want: %v",
					c.values, c.i, got, c.want)
			}
		}
	}
}

/*************************************
 *    Helper Functions
 *************************************/
var matchFieldInRegexp = regexp.MustCompile("^((in\\()+)([^ ,]+)(.*)$")
var matchFieldMatchRegexp = regexp.MustCompile("^([^ (]+)( .*)$")

func getFieldRules(
	field string,
	rules []rule.Rule,
) []rule.Rule {
	fieldRules := make([]rule.Rule, 0)
	for _, rule := range rules {
		ruleStr := rule.String()
		ruleField := matchFieldMatchRegexp.ReplaceAllString(ruleStr, "$1")
		ruleField = matchFieldInRegexp.ReplaceAllString(ruleField, "$3")
		if field == ruleField {
			fieldRules = append(fieldRules, rule)
		}
	}
	return fieldRules
}

func rulesToSortedStrings(rules []rule.Rule) []string {
	r := make([]string, len(rules))
	for i, rule := range rules {
		r[i] = rule.String()
	}
	sort.Strings(r)
	return r
}

func matchRulesUnordered(
	rules1 []rule.Rule,
	rules2 []rule.Rule,
) (bool, string) {
	if len(rules1) != len(rules2) {
		return false, "rules different lengths"
	}
	for _, rule1 := range rules1 {
		found := false
		for _, rule2 := range rules2 {
			if rule1.String() == rule2.String() {
				found = true
				break
			}
		}
		if !found {
			return false, fmt.Sprintf("rule doesn't exist: %s", rule1)
		}
	}
	return true, ""
}
