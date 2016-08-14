package epoch

import (
	"encoding/json"
	"testing"
	"time"
)

type tstruct struct {
	E Time `json:"ts"`
}

var (
	ttest = time.Unix(1, 1003000000)
)

func TestMarshal(t *testing.T) {
	ts := tstruct{Time(ttest)}

	b, err := json.Marshal(ts)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(b); got != `{"ts":2003}` {
		t.Errorf("unexpected json: %s\n", got)
	}
}

func TestUnmarshal(t *testing.T) {
	var tt tstruct
	sec := []struct {
		s string
		i int64
	}{
		{`{"ts":1000000001}`, 1000000001000000000},
		{`{"ts":1000000001101}`, 1000000001101000000},
		{`{"ts":1000000001000001}`, 1000000001000001000},
		{`{"ts":1000000001000000001}`, 1000000001000000001},
	}

	for _, tst := range sec {
		err := json.Unmarshal([]byte(tst.s), &tt)
		if err != nil {
			t.Fatal(err)
		}
		rt := time.Time(tt.E)
		if rt.UnixNano() != tst.i {
			t.Error("unexpected time:", rt.UnixNano())
		}
	}

	json.Unmarshal([]byte(`{"ts":123123}`), &tt)
	rt := time.Time(tt.E)
	if nano := rt.UnixNano(); nano != 123123000000000 {
		t.Error("unexpected time")
	}

	err := json.Unmarshal([]byte(`{"ts":123123123123}`), &tt)
	if err == nil {
		t.Error("expected error, but got nil")
	}
	if e := err.Error(); e != "unexpected number of digits in timestamp" {
		t.Error("unexpected error:", e)
	}
}
