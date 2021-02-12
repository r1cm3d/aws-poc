package protocol

import (
	"reflect"
	"testing"
	"time"
)

func TestUnmarshalJSON_Errors(t *testing.T) {
	errDate := "unparseableData"
	cases := []struct {
		name string
		in   []byte
		want error
	}{
		{"null", []byte("null"), nil},
		{"parseError", []byte(errDate), &time.ParseError{
			Layout:     "2006-01-02",
			Value:      errDate,
			LayoutElem: "2006",
			ValueElem:  errDate,
			Message:    "",
		}}}
	var d Date

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := d.UnmarshalJSON(c.in); !reflect.DeepEqual(got, c.want) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}
