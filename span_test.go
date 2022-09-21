package tnglib_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ipsusila/tnglib"
)

func doUnmarshal(out json.Unmarshaler) {
	// "ns", "us" (or "µs"), "ms", "s", "m", "h".
	vals := []interface{}{
		10, 0, 100, -1, 100.20,
		"11ns", "11us", "11µs", "11ms", "11s", "11m", "11h",
		"12.5s", "invalid", "12x", "",
	}
	tpl := `{"i":%v, "v": %#v }`
	for idx, val := range vals {
		data := fmt.Sprintf(tpl, idx, val)
		json.Unmarshal([]byte(data), out)
	}
}

func TestSpan(t *testing.T) {
	type spanT struct {
		Idx int         `json:"i"`
		Val tnglib.Span `json:"v"`
	}
	// "ns", "us" (or "µs"), "ms", "s", "m", "h".
	vals := []interface{}{
		10, 0, 100, -1, 100.20,
		"11ns", "11us", "11µs", "11ms", "11s", "11m", "11h",
		"12.5s", "invalid", "12x", "",
	}
	tpl := `{"i":%v, "v": %#v }`
	for idx, val := range vals {
		out := spanT{}
		data := fmt.Sprintf(tpl, idx, val)
		err := json.Unmarshal([]byte(data), &out)
		if err != nil {
			t.Logf("Error: %v\n", err)
		}
		s, _ := json.Marshal(out)
		t.Logf("In: %s => val=%v, out: %s\n", data, out.Val, s)
	}
}

func BenchmarkSpan(b *testing.B) {
	var s tnglib.Span
	for i := 0; i < b.N; i++ {
		doUnmarshal(&s)
	}
}
