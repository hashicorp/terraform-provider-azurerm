// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"

	"github.com/govalues/decimal"
)

// MustParse is as per [period.Parse] except that it panics if the string cannot be parsed.
// This is intended for setup code; don't use it for user inputs.
func MustParse[S ISOString | string](isoPeriod S) Period {
	p, err := Parse(isoPeriod)
	if err != nil {
		panic(err)
	}
	return p
}

// Parse parses strings that specify periods using ISO-8601 rules.
//
// In addition, a plus or minus sign can precede the period, e.g. "-P10D"
//
// It is possible to mix a number of weeks with other fields (e.g. P2M1W), although
// this would not be allowed by ISO-8601. See [Period.SimplifyWeeks].
//
// The zero value can be represented in several ways: all of the following
// are equivalent: "P0Y", "P0M", "P0W", "P0D", "PT0H", PT0M", PT0S", and "P0".
// The canonical zero is "P0D".
func Parse[S ISOString | string](isoPeriod S) (Period, error) {
	p := Period{}
	err := p.Parse(string(isoPeriod))
	return p, err
}

// Parse parses strings that specify periods using ISO-8601 rules.
//
// In addition, a plus or minus sign can precede the period, e.g. "-P10D"
//
// It is possible to mix a number of weeks with other fields (e.g. P2M1W), although
// this would not be allowed by ISO-8601. See [Period.SimplifyWeeks].
//
// The zero value can be represented in several ways: all of the following
// are equivalent: "P0Y", "P0M", "P0W", "P0D", "PT0H", PT0M", PT0S", and "P0".
// The canonical zero is "P0D".
func (period *Period) Parse(isoPeriod string) error {
	if isoPeriod == "" {
		return fmt.Errorf(`cannot parse a blank string as a period`)
	}

	p := Zero

	remaining := isoPeriod
	switch remaining[0] {
	case '-':
		p.neg = true
		remaining = remaining[1:]
	case '+':
		remaining = remaining[1:]
	}

	switch remaining {
	case "P0Y", "P0M", "P0W", "P0D", "PT0H", "PT0M", "PT0S":
		*period = Zero
		return nil // zero case
	case "":
		return fmt.Errorf(`cannot parse a blank string as a period`)
	}

	if remaining[0] != 'P' {
		return fmt.Errorf("%s: expected 'P' period mark at the start", isoPeriod)
	}
	remaining = remaining[1:]

	var haveFraction bool
	var number decimal.Decimal
	var years, months, weeks, days, hours, minutes, seconds itemState
	var des, previous Designator
	var err error
	nComponents := 0

	years, months, weeks, days = armed, armed, armed, armed

	isHMS := false
	for len(remaining) > 0 {
		if remaining[0] == 'T' {
			if isHMS {
				return fmt.Errorf("%s: 'T' designator cannot occur more than once", isoPeriod)
			}
			isHMS = true

			years, months, weeks, days = unready, unready, unready, unready
			hours, minutes, seconds = armed, armed, armed

			remaining = remaining[1:]

		} else {
			number, des, remaining, err = parseNextField(remaining, isoPeriod, isHMS)
			if err != nil {
				return err
			}

			if haveFraction && number.Coef() != 0 {
				return fmt.Errorf("%s: '%c' & '%c' only the last field can have a fraction", isoPeriod, previous.Byte(), des.Byte())
			}

			switch des {
			case Year:
				years, err = years.testAndSet(number, Year, &p.years, isoPeriod)
			case Month:
				months, err = months.testAndSet(number, Month, &p.months, isoPeriod)
			case Week:
				weeks, err = weeks.testAndSet(number, Week, &p.weeks, isoPeriod)
			case Day:
				days, err = days.testAndSet(number, Day, &p.days, isoPeriod)
			case Hour:
				hours, err = hours.testAndSet(number, Hour, &p.hours, isoPeriod)
			case Minute:
				minutes, err = minutes.testAndSet(number, Minute, &p.minutes, isoPeriod)
			case Second:
				seconds, err = seconds.testAndSet(number, Second, &p.seconds, isoPeriod)
			default:
				panic(fmt.Errorf("unreachable %s: '%c'", isoPeriod, des.Byte()))
			}
			nComponents++

			if err != nil {
				return err
			}

			if number.Scale() > 0 {
				haveFraction = true
				previous = des
			}
		}
	}

	if nComponents == 0 {
		return fmt.Errorf("%s: expected 'Y', 'M', 'W', 'D', 'H', 'M', or 'S' designator", isoPeriod)
	}

	*period = p.normaliseSign()
	return nil
}

//-------------------------------------------------------------------------------------------------

type itemState int

const (
	unready itemState = iota
	armed
	set
)

func (i itemState) testAndSet(number decimal.Decimal, des Designator, result *decimal.Decimal, original string) (itemState, error) {
	switch i {
	case unready:
		return i, fmt.Errorf("%s: '%c' designator cannot occur here", original, des.Byte())
	case set:
		return i, fmt.Errorf("%s: '%c' designator cannot occur more than once", original, des.Byte())
	}

	*result = number
	return set, nil
}

//-------------------------------------------------------------------------------------------------

func parseNextField(str, original string, isHMS bool) (decimal.Decimal, Designator, string, error) {
	number, i := scanDigits(str)
	switch i {
	case noNumberFound:
		return decimal.Zero, 0, "", fmt.Errorf("%s: expected a number but found '%c'", original, str[0])
	case stringIsAllNumeric:
		return decimal.Zero, 0, "", fmt.Errorf("%s: missing designator at the end", original)
	}

	dec, err := decimal.Parse(number)
	if err != nil {
		return decimal.Zero, 0, "", fmt.Errorf("%s: number invalid or out of range", original)
	}

	des, err := asDesignator(str[i], isHMS)
	if err != nil {
		return decimal.Zero, 0, "", fmt.Errorf("%s: %w", original, err)
	}

	return dec, des, str[i+1:], err
}

// scanDigits finds the index of the first non-digit character after some digits.
func scanDigits(s string) (string, int) {
	rs := []rune(s)
	number := make([]rune, 0, len(rs))

	for i, c := range rs {
		if i == 0 && c == '-' {
			number = append(number, c)
		} else if c == '.' || c == ',' {
			number = append(number, '.') // next step needs decimal point not comma
		} else if '0' <= c && c <= '9' {
			number = append(number, c)
		} else if len(number) > 0 {
			return string(number), i // index of the next non-digit character
		} else {
			return "", noNumberFound
		}
	}
	return "", stringIsAllNumeric
}

const (
	noNumberFound      = -1
	stringIsAllNumeric = -2
)
