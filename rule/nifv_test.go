package rule

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestNiFVNew_panics(t *testing.T) {
	values := []*dlit.Literal{}
	wantPanic := "NewNiFV: Must contain at least one value"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("New() didn't panic")
		} else if r.(string) != wantPanic {
			t.Errorf("New() - got panic: %s, wanted: %s",
				r, wantPanic)
		}
	}()
	field := "station"
	NewNiFV(field, values)
}

func TestNiFVString(t *testing.T) {
	cases := []struct {
		values []*dlit.Literal
		want   string
	}{
		{values: []*dlit.Literal{
			dlit.MustNew(7.892),
			dlit.MustNew("harry"),
			dlit.MustNew(""),
			dlit.MustNew(" harry "),
			dlit.MustNew("fred and win"),
		},
			want: "ni(station,\"7.892\",\"harry\",\"\",\" harry \",\"fred and win\")",
		},
	}

	field := "station"
	for _, c := range cases {
		r := NewNiFV(field, c.values)
		got := r.String()
		if got != c.want {
			t.Errorf("String() got: %s, want: %s", got, c.want)
		}
	}
}

func TestNiFVGetInNiParts(t *testing.T) {
	field := "station"
	values := []*dlit.Literal{dlit.MustNew(7.892)}
	r := NewNiFV(field, values)
	a, b, c := r.GetInNiParts()
	if !a || b != "ni" || c != field {
		t.Errorf("GetInNiParts() got: %t,\"%s\",\"%s\" - want: %t,\"\",\"\"",
			a, b, c, true)
	}
}

func TestNiFVIsTrue(t *testing.T) {
	cases := []struct {
		field  string
		values []*dlit.Literal
		want   bool
	}{
		{field: "station1",
			values: []*dlit.Literal{
				dlit.MustNew(7.892),
				dlit.MustNew("harry"),
				dlit.MustNew(""),
				dlit.MustNew(" harry "),
				dlit.MustNew("fred and win"),
				dlit.MustNew("true"),
			},
			want: false,
		},
		{field: "station2",
			values: []*dlit.Literal{
				dlit.MustNew(7.892),
				dlit.MustNew("harry"),
				dlit.MustNew(""),
				dlit.MustNew(" harry "),
				dlit.MustNew("fred and win"),
				dlit.MustNew("true"),
			},
			want: false,
		},
		{field: "station3",
			values: []*dlit.Literal{
				dlit.MustNew(7.892),
				dlit.MustNew("harry"),
				dlit.MustNew(""),
				dlit.MustNew(" harry "),
				dlit.MustNew("fred and win"),
				dlit.MustNew("true"),
			},
			want: true,
		},
		{field: "flow1",
			values: []*dlit.Literal{
				dlit.MustNew(7.892),
				dlit.MustNew("harry"),
				dlit.MustNew(""),
				dlit.MustNew(" harry "),
				dlit.MustNew("fred and win"),
				dlit.MustNew("true"),
			},
			want: false,
		},
		{field: "flow2",
			values: []*dlit.Literal{
				dlit.MustNew(7.892),
				dlit.MustNew("harry"),
				dlit.MustNew(""),
				dlit.MustNew(" harry "),
				dlit.MustNew("fred and win"),
				dlit.MustNew("true"),
			},
			want: true,
		},
		{field: "success1",
			values: []*dlit.Literal{
				dlit.MustNew(7.892),
				dlit.MustNew("harry"),
				dlit.MustNew(""),
				dlit.MustNew(" harry "),
				dlit.MustNew("fred and win"),
				dlit.MustNew("true"),
			},
			want: false,
		},
		{field: "success2",
			values: []*dlit.Literal{
				dlit.MustNew(7.892),
				dlit.MustNew("harry"),
				dlit.MustNew(""),
				dlit.MustNew(" harry "),
				dlit.MustNew("fred and win"),
				dlit.MustNew("true"),
			},
			want: true,
		},
	}
	record := map[string]*dlit.Literal{
		"station1": dlit.MustNew("harry"),
		"station2": dlit.MustNew(" harry "),
		"station3": dlit.MustNew("  harry  "),
		"flow1":    dlit.MustNew(7.892),
		"flow2":    dlit.MustNew(7.893),
		"success1": dlit.MustNew("true"),
		"success2": dlit.MustNew("TRUE"),
		"band":     dlit.MustNew("alpha"),
	}
	for _, c := range cases {
		r := NewNiFV(c.field, c.values)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestNiFVIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		values  []*dlit.Literal
		wantErr error
	}{
		{field: "fred",
			values: []*dlit.Literal{dlit.NewString("hello")},
			wantErr: InvalidRuleError{
				Rule: NewNiFV("fred", []*dlit.Literal{dlit.NewString("hello")}),
			},
		},
		{field: "problem",
			values: []*dlit.Literal{dlit.NewString("hello")},
			wantErr: IncompatibleTypesRuleError{
				Rule: NewNiFV("problem", []*dlit.Literal{dlit.NewString("hello")}),
			},
		},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.NewString("alpha"),
		"problem": dlit.MustNew(errors.New("this is an error")),
	}
	for _, c := range cases {
		r := NewNiFV(c.field, c.values)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}
