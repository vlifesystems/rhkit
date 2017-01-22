package rule

import (
	"github.com/lawrencewoodman/dlit"
	"reflect"
	"testing"
)

func TestNewBetweenFVI(t *testing.T) {
	min := int64(5)
	max := int64(6)
	r, err := NewBetweenFVI("flow", min, max)
	if err != nil {
		t.Errorf("NewBetweenFVI(%s, %d, %d) got err: %s", "flow", min, max, err)
	}
	if r == nil {
		t.Errorf("NewBetweenFVI(%s, %d, %d) got r: nil", "flow", min, max)
	}
}

func TestNewBetweenFVI_errors(t *testing.T) {
	cases := []struct {
		min        int64
		max        int64
		wantErrStr string
	}{
		{min: 5,
			max:        5,
			wantErrStr: "can't create Between rule where max: 5 <= min: 5",
		},
		{min: 6,
			max:        5,
			wantErrStr: "can't create Between rule where max: 5 <= min: 6",
		},
	}
	field := "flow"
	for _, c := range cases {
		r, err := NewBetweenFVI(field, c.min, c.max)
		if r != nil {
			t.Errorf("NewBetweenFVI(%s, %d, %d) rule got: %s, want: nil",
				field, c.min, c.max, r)
		}
		if err == nil {
			t.Errorf("NewBetweenFVI(%s, %d, %d) got err: nil, want: %s",
				field, c.min, c.max, c.wantErrStr)
		} else if err.Error() != c.wantErrStr {
			t.Errorf("NewBetweenFVI(%s, %d, %d) got err: %s, want: %s",
				field, c.min, c.max, err, c.wantErrStr)
		}
	}
}

func TestBetweenFVIString(t *testing.T) {
	field := "flow"
	min := int64(183)
	max := int64(287)
	want := "flow >= 183 && flow <= 287"
	r, err := NewBetweenFVI(field, min, max)
	if err != nil {
		t.Fatalf("NewBetweenFVI: %s", err)
	}
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestBetweenFVIIsTrue(t *testing.T) {
	cases := []struct {
		field string
		min   int64
		max   int64
		want  bool
	}{
		{field: "income", min: 18, max: 20, want: true},
		{field: "income", min: 19, max: 20, want: true},
		{field: "income", min: 18, max: 19, want: true},
		{field: "income", min: 10, max: 25, want: true},
		{field: "income", min: 10, max: 18, want: false},
		{field: "income", min: 20, max: 30, want: false},
		{field: "cost", min: 20, max: 30, want: true},
		{field: "cost", min: 20, max: 25, want: true},
		{field: "cost", min: 25, max: 30, want: true},
		{field: "cost", min: 20, max: 24, want: false},
		{field: "cost", min: 26, max: 30, want: false},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(25),
	}
	for _, c := range cases {
		r, err := NewBetweenFVI(c.field, c.min, c.max)
		if err != nil {
			t.Fatalf("NewBetweenFVI: %s", err)
		}
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestBetweenFVIIsTrue_errors(t *testing.T) {
	field := "rate"
	min := int64(18)
	max := int64(20)
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(18),
		"band":   dlit.NewString("alpha"),
	}
	r, err := NewBetweenFVI(field, min, max)
	if err != nil {
		t.Fatalf("NewBetweenFVI: %s", err)
	}
	wantErr := InvalidRuleError{Rule: r}
	_, err = r.IsTrue(record)
	if err != wantErr {
		t.Errorf("IsTrue(record) rule: %s, err: %v, want: %v", r, err, wantErr)
	}
}

func TestBetweenFVIGetFields(t *testing.T) {
	field := "rate"
	min := int64(18)
	max := int64(20)
	want := []string{"rate"}
	r, err := NewBetweenFVI(field, min, max)
	if err != nil {
		t.Fatalf("NewBetweenFVI: %s", err)
	}
	got := r.GetFields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetFields() got: %s, want: %s", got, want)
	}
}