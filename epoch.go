/*
Package epoch provides a timestamp type that implements json.Marshaler and
json.Unmarshaler.

When marshalled, epoch.Time will be converted to milliseconds since January 1,
1970. However, any timestamp resolution (seconds, milliseconds, microseconds,
nanoseconds) may be unmarshalled to epoch.Time.
*/
package epoch

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Time is an alias type for time.Time
type Time struct {
	time.Time
}

// MarshalJSON returns e as the string representation of the number of
// milliseconds since epoch
func (e Time) MarshalJSON() ([]byte, error) {
	t := e.UnixNano() / 1000000
	return []byte(strconv.FormatInt(t, 10)), nil
}

// UnmarshalJSON interprets data as an epoch timestamp and sets *e to a
// time.Time object with the corresponding value. The timestamp can include
// a fractional part which will be included in the resulting time object. The
// number of digits in the integer part determine the expected resolution of
// the timestamp:
//
// A timestamp with the integer part having less than 13 digits is interpreted
// as seconds since epoch.
//
// A timestamp with the integer part having less than 16 digits is interpreted
// as milliseconds since epoch.
//
// A timestamp with the integer part having less than 19 digits is interpreted
// as microseconds since epoch.
//
// A timestamp with the integer part having 19 or more digits is interpreted
// as nanoseconds since epoch.
func (e *Time) UnmarshalJSON(data []byte) error {
	var (
		intPart, fracPart string
		shiftLen          int
		t                 time.Time
		ts                = string(data)
	)

	ts = strings.Replace(ts, `"`, "", -1)

	// validate the timestamp as a number (integer or decimal)
	err := validateTimestamp(ts)
	if err != nil {
		return err
	}

	// proceed with the party
	intPart = ts
	p := strings.Split(ts, ".")
	if len(p) == 2 {
		intPart, fracPart = p[0], p[1]
	}

	// number of digits in a string timestamp that correspond to seconds given
	// the length of the timestamp itself
	shiftLen = timeIntLen(intPart)

	t, err = timeFromSecString(intPart, fracPart, shiftLen)
	if err != nil {
		return errors.New("could not parse timestamp")
	}

	//*(*time.Time)(e) = t
	(*e).Time = t
	return nil
}

func validateTimestamp(ts string) error {
	// check if we even have something
	l := len(ts)
	if l < 1 {
		return errors.New("data cannot be an empty number")
	}

	// check if positive number
	if string(ts[0]) == "-" {
		return errors.New("data is not a positive number")
	}

	// check if floating number
	// TODO(guillermo): Consider comma (,) separator as well?
	sep := strings.Count(ts, ".")
	if sep > 1 {
		return errors.New("data is not a valid decimal number")
	}

	// finally check if it is really a parsable number
	_, err := strconv.ParseFloat(ts, 64)
	if err != nil {
		return errors.New("data is not a parsable number")
	}

	return nil
}

// Measure the chunk of the integer part of a valid timestamp that will give seconds
func timeIntLen(ts string) int {
	var shiftLen int

	switch {
	case len(ts) < 13: // less than 13 integer digits => seconds
		shiftLen = len(ts)
	case len(ts) < 16: // between 13 and 15 integer digits => milliseconds
		shiftLen = len(ts) - 3
	case len(ts) < 19: // between 16 and 18 integer digits => microseconds
		shiftLen = len(ts) - 6
	case len(ts) >= 19: // 19 and above integer digits => nanoseconds
		shiftLen = len(ts) - 9
	}

	return shiftLen
}

func timeFromSecString(intPart, fracPart string, shiftLen int) (time.Time, error) {
	var (
		sec, nano int64
		err       error
	)

	mv := intPart[shiftLen:]
	fracPart = mv + fracPart
	intPart = intPart[:shiftLen]

	sec, err = strconv.ParseInt(intPart, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	if len(fracPart) > 0 {
		if len(fracPart) <= 9 {
			fracPart = fracPart + "0000000000"[:9-len(fracPart)]
		}
		if len(fracPart) > 9 {
			fracPart = fracPart[:9]
		}
		nano, err = strconv.ParseInt(fracPart, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
	}
	return time.Unix(sec, nano).UTC(), nil
}
