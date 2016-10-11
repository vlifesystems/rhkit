package rhkit

import (
	"fmt"
	"testing"
)

func TestFieldTypeString(t *testing.T) {
	cases := []struct {
		in   fieldType
		want string
	}{
		{ftUnknown, "Unknown"},
		{ftIgnore, "Ignore"},
		{ftInt, "Int"},
		{ftFloat, "Float"},
		{ftString, "String"},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("String() c.in:%d got: %s, want: %s", c.in, got, c.want)
		}
	}
}

func TestFieldTypeString_panic(t *testing.T) {
	kind := fieldType(99)
	paniced := false
	wantPanic := fmt.Sprintf("Unsupported type: %d", kind)
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("String() - got panic: %s, wanted: %s", r, wantPanic)
			}
		}
	}()
	kind.String()
	if !paniced {
		t.Errorf("String() - failed to panic with: %s", wantPanic)
	}
}
