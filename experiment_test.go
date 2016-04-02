package rulehunter

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/rulehunter/internal"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Ensure that correct number is returned
func TestLoadExperiment(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	expectedExperiments := []*Experiment{
		&Experiment{},
		&Experiment{
			FileFormatVersion: "0.1",
			Title:             "This is a jolly nice title",
			Input: mustNewCsvInput(
				fieldNames,
				filepath.Join("fixtures", "bank.csv"),
				rune(';'),
				true,
			),
			FieldNames:        fieldNames,
			ExcludeFieldNames: []string{"education"},
			Aggregators: []internal.Aggregator{
				mustNewCountAggregator("numSignedUp", "y == \"yes\""),
				mustNewCalcAggregator("cost", "numMatches * 4.5"),
				mustNewCalcAggregator("income", "numSignedUp * 24"),
				mustNewCalcAggregator("profit", "income - cost")},
			Goals: []*dexpr.Expr{mustNewDExpr("profit > 0")},
			SortOrder: []SortField{
				SortField{"profit", DESCENDING},
				SortField{"numSignedUp", DESCENDING}},
		},
	}
	cases := []struct {
		filename string
		want     *Experiment
		wantErr  error
	}{
		{"missingfile.json", expectedExperiments[0],
			&os.PathError{"open", "missingfile.json",
				errors.New("no such file or directory")}},
		{filepath.Join("fixtures", "missingFileFormatVersion.json"),
			expectedExperiments[0],
			&ErrInvalidField{"fileFormatVersion", "",
				errors.New("Must have a valid version number")}},
		{filepath.Join("fixtures", "singleFieldName.json"),
			expectedExperiments[0],
			&ErrInvalidField{"fieldNames", "[\"age\"]",
				errors.New("Must specify at least two field names")}},
		{filepath.Join("fixtures", "bank.json"),
			expectedExperiments[1], nil},
	}
	for _, c := range cases {
		got, err := LoadExperiment(c.filename)
		if !errorMatch(c.wantErr, err) ||
			(err == nil && !experimentMatch(got, c.want)) {
			t.Errorf("LoadExperiment(%q) err: %q, wantErr: %q - got: %q, want: %q",
				c.filename, err, c.wantErr, got, c.want)
		}
	}
}

/***********************
   Helper functions
************************/

func experimentMatch(e1 *Experiment, e2 *Experiment) bool {
	if e1.FileFormatVersion != e2.FileFormatVersion ||
		e1.Title != e2.Title ||
		!areStringArraysEqual(e1.FieldNames, e2.FieldNames) ||
		!areStringArraysEqual(e1.ExcludeFieldNames, e2.ExcludeFieldNames) ||
		!areGoalsEqual(e1.Goals, e2.Goals) ||
		!areAggregatorsEqual(e1.Aggregators, e2.Aggregators) ||
		!areSortOrdersEqual(e1.SortOrder, e2.SortOrder) {
		return false
	}
	for {
		e1Record, e1Err := e1.Input.Read()
		e2Record, e2Err := e2.Input.Read()
		if e1Err != e2Err {
			return false
		} else if e1Err == nil && e2Err == nil {
			if !reflect.DeepEqual(e1Record, e2Record) {
				return false
			}
		} else if e1Err == io.EOF || e2Err == io.EOF {
			break
		}
	}
	return true
}

func areStringArraysEqual(a1 []string, a2 []string) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i, e := range a1 {
		if e != a2[i] {
			return false
		}
	}
	return true
}

func areGoalsEqual(g1 []*dexpr.Expr, g2 []*dexpr.Expr) bool {
	if len(g1) != len(g2) {
		return false
	}
	for i, g := range g1 {
		if g.String() != g2[i].String() {
			return false
		}
	}
	return true

}

func areAggregatorsEqual(
	a1 []internal.Aggregator,
	a2 []internal.Aggregator,
) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i, e := range a1 {
		if !e.IsEqual(a2[i]) {
			return false
		}
	}
	return true
}

func areSortOrdersEqual(so1 []SortField, so2 []SortField) bool {
	if len(so1) != len(so2) {
		return false
	}
	for i, sf1 := range so1 {
		sf2 := so2[i]
		if sf1.Field != sf2.Field || sf1.Direction != sf2.Direction {
			return false
		}
	}
	return true
}

func mustNewCsvInput(
	fieldNames []string,
	filename string,
	separator rune,
	skipFirstLine bool,
) Input {
	input, err :=
		newCsvInput(fieldNames, filename, separator, skipFirstLine)
	if err != nil {
		panic(fmt.Sprintf("Couldn't create Csv Input: %s", err))
	}
	return input
}
