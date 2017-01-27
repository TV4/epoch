// Package epoch provides a timestamp type that implements both json.Marshaler
// and json.Unmarshaler.
package epoch

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Time is an alias type for time.Time
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
// Numbers with fractional parts are also accepted.
//
// A string of digits with the integer part equal to or shorter than 10
// digits is interpreted as seconds.
// A string of digits with the integer part equal to or shorter than 13
// digits is interpreted as seconds.
// A string of digits with the integer part equal to or shorter than 16
// digits is interpreted as seconds.
// A string of digits with the integer part equal to or shorter than 19
// digits is interpreted as seconds.
func (e *Time) UnmarshalJSON(data []byte) error {
	var (
		intPart, fracPart string
		t                 time.Time
		ts                = string(data)
	)

	ts = strings.Replace(ts, `"`, "", -1)

	// check if number
	_, err := strconv.ParseFloat(ts, 64)
	if err != nil {
		return errors.New("data is not a number")
	}

	intPart = ts
	p := strings.Split(ts, ".")
	if len(p) > 2 {
		return errors.New("data is not a number")
	}
	if len(p) == 2 {
		intPart, fracPart = p[0], p[1]
	}

	// shiftLen measures the chunk of the integer part of the timestamp
	// that will give seconds.
	var shiftLen int
	if len(intPart) < 13 { // less than 13 integer digits => seconds
		shiftLen = len(intPart)
	} else if len(intPart) < 16 { // between 13 and 15 integer digits => milliseconds
		shiftLen = len(intPart) - 3
	} else if len(intPart) < 19 { // between 16 and 18 integer digits => microseconds
		shiftLen = len(intPart) - 6
	} else if len(intPart) >= 19 { // 19 and above integer digits => nanoseconds
		shiftLen = len(intPart) - 9
	}

	t, err = timeFromSecString(intPart, fracPart, shiftLen)
	if err != nil {
		return errors.New("could not parse timestamp")
	}

	*(*time.Time)(e) = t
	return nil
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
		// put a '1' first to not lose any leading zeros
		fracPart = "1" + fracPart
		// add trailing zeros to make it nanosecond large
		if len(fracPart) <= 10 {
			fracPart = fracPart + "0000000000"[:10-len(fracPart)]
		}
		nano, err = strconv.ParseInt(fracPart, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		// retract the leading '1'
		nano = nano - 1000000000
	}
	return time.Unix(sec, nano), nil
}
