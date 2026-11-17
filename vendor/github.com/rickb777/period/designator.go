// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"strconv"
)

// Designator enumerates the seven fields in a Period.
type Designator int8

const (
	_ Designator = iota
	Second
	Minute
	Hour
	Day
	Week
	Month
	Year
)

func asDesignator(d byte, isHMS bool) (Designator, error) {
	switch d {
	case 'S':
		return Second, nil
	case 'H':
		return Hour, nil
	case 'D':
		return Day, nil
	case 'W':
		return Week, nil
	case 'Y':
		return Year, nil
	case 'M':
		if isHMS {
			return Minute, nil
		}
		return Month, nil
	}
	return 0, fmt.Errorf("expected a designator Y, M, W, D, H, or S not '%c'", d)
}

func (d Designator) Byte() byte {
	switch d {
	case Second:
		return 'S'
	case Minute:
		return 'M'
	case Hour:
		return 'H'
	case Day:
		return 'D'
	case Week:
		return 'W'
	case Month:
		return 'M'
	case Year:
		return 'Y'
	}
	panic(strconv.Itoa(int(d)))
}

//func (d designator) field() string {
//	switch d {
//	case second:
//		return "seconds"
//	case minute:
//		return "minutes"
//	case hour:
//		return "hours"
//	case Day:
//		return "days"
//	case week:
//		return "weeks"
//	case month:
//		return "months"
//	case year:
//		return "years"
//	}
//	panic(strconv.Itoa(int(d)))
//}
//
//func (d designator) min(other designator) designator {
//	if d < other {
//		return d
//	}
//	return other
//}
//
//func (d designator) IsOneOf(xx ...designator) bool {
//	for _, x := range xx {
//		if x == d {
//			return true
//		}
//	}
//	return false
//}
//
//func (d designator) IsNotOneOf(xx ...designator) bool {
//	for _, x := range xx {
//		if x == d {
//			return false
//		}
//	}
//	return true
//}
