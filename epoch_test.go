package epoch

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMarshal(t *testing.T) {
	e := Time(time.Unix(1, 1003000000))
	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(b); got != `2003` {
		t.Errorf("unexpected json: %s\n", got)
	}
}

func TestUnmarshal(t *testing.T) {
	for _, tst := range []struct {
		s string
		i int64
	}{
		{`1`, 1000000000},
		{`1.03`, 1030000000},
		{`0.00000003`, 30},
		{`0.000000003`, 3},
		{`0.0000000003`, 0},
		{`"1000000001"`, 1000000001000000000},
		{`1000000001.030`, 1000000001030000000},
		{`1000000001.300`, 1000000001300000000},
		{`1000000001.000000001`, 1000000001000000001},
		{`1000000001.0300000`, 1000000001030000000},
		{`1000000001101`, 1000000001101000000},
		{`1000000001101.001`, 1000000001101001000},
		{`1000000001000001.1`, 1000000001000001100},
		{`1000000001000000001`, 1000000001000000001},
	} {
		var e Time
		err := json.Unmarshal([]byte(tst.s), &e)
		if err != nil {
			t.Fatal(err, tst.s)
		}
		rt := time.Time(e)
		if rt.UnixNano() != tst.i {
			t.Error("unexpected time:", rt.UnixNano(), tst.i)
		}
		if got, want := rt.Location(), time.UTC; got != want {
			t.Errorf("location = %q, want %q", got, want)
		}
	}
}

func TestBadInput(t *testing.T) {
	for _, s := range []string{
		`"bad"`,
		`"1.2.3"`,
		`"1.a"`,
		`"."`,
	} {
		var e Time
		err := json.Unmarshal([]byte(s), &e)
		if err == nil {
			t.Fatal("expected error but got nil", s)
		}
	}
}
