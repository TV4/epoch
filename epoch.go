// package epoch provides a timestamp type that implements both json.Marshaler
// and json.Unmarshaler.
package epoch

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Time time.Time

// MarshalJSON returns e as the string representation of the number of
// milliseconds since epoch
func (e Time) MarshalJSON() ([]byte, error) {
	t := time.Time(e).UnixNano() / 1000000
	return []byte(strconv.FormatInt(t, 10)), nil
}

// UnmarshalJSON interprets data as an epoch timestamp and sets *e to a
// time.Time object with the corresponding value. Seconds, milliseconds,
// microseconds, and nanoseconds since epoch are all accepted values.
func (e *Time) UnmarshalJSON(data []byte) error {
	var (
		q, n int64
		err  error
	)
	r := strings.Replace(string(data), `"`, ``, -1)
	if len(r) < 10 {
		pad := "0000000000"[:10-len(r)]
		r = pad + r
	}

	q, err = strconv.ParseInt(r[:10], 10, 64)
	if len(r) > 10 {
		n, err = strconv.ParseInt(r[11:], 10, 64)
	}

	var t time.Time
	switch len(r) {
	case 10: //sec
		t = time.Unix(q, n)
	case 13: //msec
		t = time.Unix(q, n*1000000)
	case 16: //musec
		t = time.Unix(q, n*1000)
	case 19: //nanosec
		t = time.Unix(q, n)
	default:
		return errors.New("unexpected number of digits in timestamp")
	}

	if err != nil {
		return err
	}
	*(*time.Time)(e) = t
	return nil
}
