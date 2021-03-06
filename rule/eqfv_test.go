package rule

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/testhelpers"
	"math"
	"reflect"
	"testing"
)

func TestEQFVString(t *testing.T) {
	cases := []struct {
		field string
		value *dlit.Literal
		want  string
	}{
		{field: "income", value: dlit.MustNew(7.8903), want: "income == 7.8903"},
		{field: "income", value: dlit.MustNew(7.890300), want: "income == 7.8903"},
		{field: "income", value: dlit.MustNew(7.), want: "income == 7"},
		{field: "income", value: dlit.MustNew(7.00), want: "income == 7"},
		{field: "income", value: dlit.MustNew(7), want: "income == 7"},
		{field: "income", value: dlit.MustNew(0.34), want: "income == 0.34"},
		{field: "income", value: dlit.MustNew(0.3400), want: "income == 0.34"},
		{field: "income", value: dlit.MustNew(0.0), want: "income == 0"},
		{field: "income", value: dlit.MustNew(-53.9), want: "income == -53.9"},
		{field: "name", value: dlit.MustNew("borris"), want: "name == \"borris\""},
		{field: "name", value: dlit.MustNew("bo   rris"),
			want: "name == \"bo   rris\""},
		{field: "name", value: dlit.MustNew("  borris  "),
			want: "name == \"  borris  \""},
		{field: "name", value: dlit.MustNew(""), want: "name == \"\""},
	}
	for _, c := range cases {
		r := NewEQFV(c.field, c.value)
		got := r.String()
		if got != c.want {
			t.Errorf("String() got: %s, want: %s", got, c.want)
		}
	}
}

func TestEQFVIsTrue(t *testing.T) {
	cases := []struct {
		field string
		value *dlit.Literal
		want  bool
	}{
		{"income", dlit.MustNew(19.0), true},
		{"income", dlit.MustNew(int64(math.MaxInt64)), false},
		{"income", dlit.MustNew(-19.0), false},
		{"income", dlit.MustNew(20.0), false},
		{"flow", dlit.MustNew(124.564), true},
		{"flow", dlit.MustNew(-124.564), false},
		{"flow", dlit.MustNew(20.0), false},
		{"flow", dlit.MustNew(124.5645), false},
		{"flow", dlit.MustNew(124.565), false},
		{"band", dlit.MustNew("hello"), false},
		{"band", dlit.MustNew("alpha"), true},
		{"band", dlit.MustNew("ALPHA"), false},
		{"band", dlit.MustNew(8.9), false},
		{"success", dlit.MustNew("TRUE"), true},
		{"success", dlit.MustNew("true"), false},
		{"success", dlit.MustNew("1"), false},
		{"bigNums", dlit.MustNew(int64(math.MaxInt64)), true},
		{"bigNums", dlit.MustNew("1"), false},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.MustNew("alpha"),
		"success": dlit.MustNew("TRUE"),
		"bigNums": dlit.MustNew(int64(math.MaxInt64)),
	}
	for _, c := range cases {
		r := NewEQFV(c.field, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestEQFVIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		value   *dlit.Literal
		wantErr error
	}{
		{"fred", dlit.MustNew(8.9),
			InvalidRuleError{Rule: NewEQFV("fred", dlit.MustNew(8.9))}},
		{"problem", dlit.MustNew(8.9),
			IncompatibleTypesRuleError{Rule: NewEQFV("problem", dlit.MustNew(8.9))}},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.NewString("alpha"),
		"problem": dlit.MustNew(errors.New("this is an error")),
	}
	for _, c := range cases {
		r := NewEQFV(c.field, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestEQFVFields(t *testing.T) {
	r := NewEQFV("income", dlit.MustNew(5.5))
	want := []string{"income"}
	got := r.Fields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Fields() got: %s, want: %s", got, want)
	}
}

func TestGenerateEQFV(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"bandA": {
				Kind: description.Number,
				Min:  dlit.MustNew(1),
				Max:  dlit.MustNew(4),
				Values: map[string]description.Value{
					"1": {dlit.NewString("1"), 3},
					"2": {dlit.NewString("2"), 1},
					"3": {dlit.NewString("3"), 2},
					"4": {dlit.NewString("4"), 5},
				},
			},
			"bandB": {
				Kind: description.Number,
				Min:  dlit.MustNew(1),
				Max:  dlit.MustNew(4),
				Values: map[string]description.Value{
					"1": {dlit.NewString("1"), 3},
					"2": {dlit.NewString("2"), 1},
					"3": {dlit.NewString("3"), 2},
					"4": {dlit.NewString("4"), 5},
				},
			},
			"flow": {
				Kind:  description.Number,
				Min:   dlit.MustNew(1),
				Max:   dlit.MustNew(4),
				MaxDP: 2,
				Values: map[string]description.Value{
					"1":    {dlit.NewString("1"), 3},
					"2":    {dlit.NewString("2"), 1},
					"2.90": {dlit.NewString("2.90"), 1},
					"3.37": {dlit.NewString("3.37"), 2},
					"3.3":  {dlit.NewString("3.3"), 2},
					"4":    {dlit.NewString("4"), 5},
				},
			},
			"group": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Nelson":      {dlit.NewString("Nelson"), 3},
					"Collingwood": {dlit.NewString("Collingwood"), 1},
					"Mountbatten": {dlit.NewString("Mountbatten"), 1},
					"Drake":       {dlit.NewString("Drake"), 2},
				},
			},
			"month": {
				Kind: description.String,
				Values: map[string]description.Value{
					"May": {dlit.NewString("May"), 3},
				},
			},
		},
	}
	want := []Rule{
		NewEQFV("bandA", dlit.MustNew(1)),
		NewEQFV("bandA", dlit.MustNew(3)),
		NewEQFV("bandA", dlit.MustNew(4)),
		NewEQFV("flow", dlit.MustNew(1)),
		NewEQFV("flow", dlit.MustNew(3.37)),
		NewEQFV("flow", dlit.MustNew(3.3)),
		NewEQFV("flow", dlit.MustNew(4)),
		NewEQFV("group", dlit.MustNew("Nelson")),
		NewEQFV("group", dlit.MustNew("Drake")),
	}
	generationDesc := testhelpers.GenerationDesc{
		DFields:     []string{"bandA", "bandB", "flow", "group", "month"},
		DArithmetic: false,
		DDeny:       map[string][]string{"EQFV": []string{"bandB"}},
	}
	got := generateEQFV(inputDescription, generationDesc)
	if err := matchRulesUnordered(got, want); err != nil {
		t.Errorf("matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
			err, got, want)
	}
}
