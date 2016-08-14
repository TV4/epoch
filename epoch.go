// package epoch provides a timestamp type that implements both json.Marshaler
// and json.Unmarshaler.
package epoch

import (
	"errors"
	"strconv"
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
		ts   = string(data)
	)

	if len(ts) < 10 {
		pad := "0000000000"[:10-len(ts)]
		ts = pad + ts
	}

	if q, err = strconv.ParseInt(ts[:10], 10, 64); err != nil {
		return err
	}

	if len(ts) > 10 {
		if n, err = strconv.ParseInt(ts[10:], 10, 64); err != nil {
			return err
		}
	}

	var t time.Time
	switch len(ts) {
	case 10: //sec
		t = time.Unix(q, 0)
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
