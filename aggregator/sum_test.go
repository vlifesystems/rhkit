package aggregator

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
	"testing"
)

func TestNewSum_error(t *testing.T) {
	_, err := New("a", "sum", "3+4+{")
	wantErr := DescError{
		Name: "a",
		Kind: "sum",
		Err: dexpr.InvalidExprError{
			Expr: "3+4+{",
			Err:  dexpr.ErrSyntax,
		},
	}.Error()
	if err.Error() != wantErr {
		t.Errorf("New: gotErr: %s, wantErr: %s", err, wantErr)
	}
}

func TestSumResult(t *testing.T) {
	records := []map[string]*dlit.Literal{
		{
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(4.5),
			"band":   dlit.MustNew(4),
		},
		{
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(3.2),
			"band":   dlit.MustNew(7),
		},
		{
			"income": dlit.MustNew(2),
			"cost":   dlit.MustNew(1.2),
			"band":   dlit.MustNew(4),
		},
		{
			"income": dlit.MustNew(9),
			"cost":   dlit.MustNew(2),
			"band":   dlit.MustNew(9),
		},
	}
	goals := []*goal.Goal{}
	profitDesc := MustNew("profit", "sum", "income-cost")
	profit := profitDesc.New()
	instances := []Instance{profit}

	for i, record := range records {
		profit.NextRecord(record, i != 2)
	}
	want := 5.3
	numRecords := int64(len(records))
	got := profit.Result(instances, goals, numRecords)
	gotFloat, gotIsFloat := got.Float()
	if !gotIsFloat || gotFloat != want {
		t.Errorf("Result() got: %s, want: %f", got, want)
	}
}

func TestSumNextRecord_errors(t *testing.T) {
	cases := []struct {
		record map[string]*dlit.Literal
		arg    string
		want   error
	}{
		{record: map[string]*dlit.Literal{},
			arg: "cost + 2",
			want: dexpr.InvalidExprError{
				Expr: "cost + 2",
				Err:  dexpr.VarNotExistError("cost"),
			},
		},
		{record: map[string]*dlit.Literal{"cost": dlit.NewString("hello")},
			arg: "cost",
			want: dexpr.InvalidExprError{
				Expr: "sum+value",
				Err:  dexpr.ErrIncompatibleTypes,
			},
		},
	}
	for _, c := range cases {
		as := MustNew("a", "sum", c.arg)
		ai := as.New()
		got := ai.NextRecord(c.record, true)
		if got == nil || got.Error() != c.want.Error() {
			t.Errorf("NextRecord: got: %s, want: %s", got, c.want)
		}
	}
}

func TestSumSpecName(t *testing.T) {
	name := "a"
	as := MustNew(name, "sum", "income-cost")
	got := as.Name()
	if got != name {
		t.Errorf("Name - got: %s, want: %s", got, name)
	}
}

func TestSumSpecKind(t *testing.T) {
	kind := "sum"
	as := MustNew("a", kind, "income-cost")
	got := as.Kind()
	if got != kind {
		t.Errorf("Kind - got: %s, want: %s", got, kind)
	}
}

func TestSumSpecArg(t *testing.T) {
	arg := "income-cost"
	as := MustNew("a", "sum", arg)
	got := as.Arg()
	if got != arg {
		t.Errorf("Arg - got: %s, want: %s", got, arg)
	}
}

func TestSumInstanceName(t *testing.T) {
	as := MustNew("abc", "sum", "cost > 2")
	ai := as.New()
	got := ai.Name()
	want := "abc"
	if got != want {
		t.Errorf("Name: got: %s, want: %s", got, want)
	}
}
