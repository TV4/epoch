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
	for _, tst := range []struct {
		s string
		i int64
	}{
		{`{"ts":1}`, 1000000000},
		{`{"ts":1.03}`, 1030000000},
		{`{"ts":0.00000003}`, 30},
		{`{"ts":0.000000003}`, 3},
		{`{"ts":0.0000000003}`, 0},
		{`{"ts":"1000000001"}`, 1000000001000000000},
		{`{"ts":1000000001.030}`, 1000000001030000000},
		{`{"ts":1000000001.300}`, 1000000001300000000},
		{`{"ts":1000000001.000000001}`, 1000000001000000001},
		{`{"ts":1000000001.0300000}`, 1000000001030000000},
		{`{"ts":1000000001101}`, 1000000001101000000},
		{`{"ts":1000000001101.001}`, 1000000001101001000},
		{`{"ts":1000000001000001.1}`, 1000000001000001100},
		{`{"ts":1000000001000000001}`, 1000000001000000001},
	} {
		var tt tstruct
		err := json.Unmarshal([]byte(tst.s), &tt)
		if err != nil {
			t.Fatal(err, tst.s)
		}
		rt := time.Time(tt.E)
		if rt.UnixNano() != tst.i {
			t.Error("unexpected time:", rt.UnixNano(), tst.i)
		}
	}
}

func TestBadInput(t *testing.T) {
	for _, tst := range []struct {
		s string
	}{
		{`{"ts":"bad"}`},
		{`{"ts":"1.2.3"}`},
		{`{"ts":"1.a"}`},
		{`{"ts":"."}`},
	} {
		var tt tstruct
		err := json.Unmarshal([]byte(tst.s), &tt)
		if err == nil {
			t.Fatal("expected error but got nil", tst.s)
		}
	}
}
