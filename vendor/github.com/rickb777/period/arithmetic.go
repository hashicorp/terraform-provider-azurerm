// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"errors"
	"math"
	"time"

	"github.com/govalues/decimal"
)

// AddTo adds the period to a time, returning the result.
// A flag is also returned that is true when the conversion was precise, and false otherwise.
//
// When the period specifies hours, minutes and seconds only, the result is precise.
//
// Similarly, when the period specifies whole years, months, weeks and days (i.e. without fractions),
// the result is precise.
//
// However, when years, months or days contains fractions, the result is only an approximation (it
// assumes that all days are 24 hours and every year is 365.2425 days, as per Gregorian calendar rules).
//
// Note: the use of AddDate has unintended consequences when considering the addition of a time only
// period. See <https://x.com/joelmcourtney/status/1803301619955904979>.
func (period Period) AddTo(t time.Time) (time.Time, bool) {
	if !zeroCalendarValues(period) && wholeCalendarValues(period) {
		// in this case, time.AddDate provides an exact solution

		years, _, ok1 := period.years.Int64(0)
		months, _, ok2 := period.months.Int64(0)
		weeks, _, ok3 := period.weeks.Int64(0)
		days, _, ok4 := period.days.Int64(0)

		hms, ok5 := totalHrMinSec(period)

		if period.neg {
			years = -years
			months = -months
			weeks = -weeks
			days = -days
			hms = -hms
		}

		t1 := t.AddDate(int(years), int(months), int(7*weeks+days)).Add(hms)
		return t1, ok1 && ok2 && ok3 && ok4 && ok5
	}

	// fractional years or months or weeks or days
	d, precise := period.Duration()
	return t.Add(d), precise
}

//-------------------------------------------------------------------------------------------------

// Add adds two periods together. Use this method along with [Period.Negate] in order to subtract periods.
// Arithmetic overflow will result in an error.
func (period Period) Add(other Period) (Period, error) {
	var left, right Period

	if period.neg {
		left = period.flipSign()
	} else {
		left = period
	}

	if other.neg {
		right = other.flipSign()
	} else {
		right = other
	}

	years, e1 := left.years.Add(right.years)
	months, e2 := left.months.Add(right.months)
	weeks, e3 := left.weeks.Add(right.weeks)
	days, e4 := left.days.Add(right.days)
	hours, e5 := left.hours.Add(right.hours)
	minutes, e6 := left.minutes.Add(right.minutes)
	seconds, e7 := left.seconds.Add(right.seconds)

	result := Period{years: years, months: months, weeks: weeks, days: days, hours: hours, minutes: minutes, seconds: seconds}.Normalise(true).normaliseSign()
	return result, errors.Join(e1, e2, e3, e4, e5, e6, e7)
}

// Subtract subtracts one period from another.
// Arithmetic overflow will result in an error.
func (period Period) Subtract(other Period) (Period, error) {
	return period.Add(other.Negate())
}

//-------------------------------------------------------------------------------------------------

// Mul multiplies a period by a factor. Obviously, this can both enlarge and shrink it,
// and change the sign if the factor is negative. The result is not normalised.
func (period Period) Mul(factor decimal.Decimal) (Period, error) {
	var years, months, weeks, days, hours, minutes, seconds decimal.Decimal
	var e1, e2, e3, e4, e5, e6, e7 error

	if period.years.Coef() != 0 {
		years, e1 = period.years.Mul(factor)
		years = years.Trim(0)
	}
	if period.months.Coef() != 0 {
		months, e2 = period.months.Mul(factor)
		months = months.Trim(0)
	}
	if period.weeks.Coef() != 0 {
		weeks, e3 = period.weeks.Mul(factor)
		weeks = weeks.Trim(0)
	}
	if period.days.Coef() != 0 {
		days, e4 = period.days.Mul(factor)
		days = days.Trim(0)
	}
	if period.hours.Coef() != 0 {
		hours, e5 = period.hours.Mul(factor)
		hours = hours.Trim(0)
	}
	if period.minutes.Coef() != 0 {
		minutes, e6 = period.minutes.Mul(factor)
		minutes = minutes.Trim(0)
	}
	if period.seconds.Coef() != 0 {
		seconds, e7 = period.seconds.Mul(factor)
		seconds = seconds.Trim(0)
	}

	result := Period{
		years:   years,
		months:  months,
		weeks:   weeks,
		days:    days,
		hours:   hours,
		minutes: minutes,
		seconds: seconds,
		neg:     period.neg,
	}

	return result.normaliseSign(), errors.Join(e1, e2, e3, e4, e5, e6, e7)
}

//-------------------------------------------------------------------------------------------------

// TotalDaysApprox gets the approximate total number of days in the period. The approximation assumes
// a year is 365.2425 days as per Gregorian calendar rules) and a month is 1/12 of that. Whole
// multiples of 24 hours are also included in the calculation.
func (period Period) TotalDaysApprox() int {
	sign := period.Sign()
	if sign == 0 {
		return 0
	}
	pn := period.Normalise(false)
	tdE9, _ := totalDaysApproxE9(pn)
	return sign * int(tdE9/1e9)
}

// TotalMonthsApprox gets the approximate total number of months in the period. The days component
// is included by approximation, assuming a year is 365.2425 days (as per Gregorian calendar rules)
// and a month is 1/12 of that. Whole multiples of 24 hours are also included in the calculation.
func (period Period) TotalMonthsApprox() int {
	sign := period.Sign()
	if sign == 0 {
		return 0
	}
	pn := period.Normalise(false)
	tdE9, _ := totalDaysApproxE9(pn)
	return sign * int((tdE9/daysPerMonthE6)/1e3)
}

//-------------------------------------------------------------------------------------------------

// DurationApprox converts a period to the equivalent duration in nanoseconds.
// When the period specifies hours, minutes and seconds only, the result is precise.
// however, when the period specifies years, months, weeks and days, it is impossible to
// be precise because the result may depend on knowing date and timezone information. So
// the duration is estimated on the basis of a year being 365.2425 days (as per Gregorian
// calendar rules) and a month being 1/12 of a that; days are all assumed to be 24 hours long.
//
// Note that time.Duration is limited to the range 1 nanosecond to about 292 years maximum.
func (period Period) DurationApprox() time.Duration {
	d, _ := period.Duration()
	return d
}

// Duration converts a period to the equivalent duration in nanoseconds.
// A flag is also returned that is true when the conversion was precise, and false otherwise.
//
// When the period specifies hours, minutes and seconds only, the result is precise.
// However, when the period specifies years, months, weeks and days, it is impossible to
// be precise because the result may depend on knowing date and timezone information. So
// the duration is estimated on the basis of a year being 365.2425 days (as per Gregorian
// calendar rules) and a month being 1/12 of a that; days are all assumed to be 24 hours long.
//
// For periods shorter than one nanosecond, the duration will be zero and the precise flag
// will be returned false.
//
// Note that time.Duration is limited to the range 1 nanosecond to about 292 years maximum.
func (period Period) Duration() (time.Duration, bool) {
	sign := time.Duration(period.Sign())
	if sign == 0 {
		return 0, true
	}
	daysE9, ok1 := totalDaysApproxE9(period)
	ymwd := time.Duration(daysE9 * secondsPerDay)
	hms, ok2 := totalHrMinSec(period)
	return sign * (ymwd + hms), ymwd == 0 && ok1 && ok2
}

func totalDaysApproxE9(period Period) (int64, bool) {
	dd, okd := fieldDuration(period.days, 1e9)
	ww, okw := fieldDuration(period.weeks, 7*1e9)
	mm, okm := fieldDuration(period.months, daysPerMonthE6*1e3)
	yy, oky := fieldDuration(period.years, daysPerYearE6*1e3)
	return dd + ww + mm + yy, okd && okw && okm && oky
}

func totalHrMinSec(period Period) (time.Duration, bool) {
	hh, okh := fieldDuration(period.hours, int64(time.Hour))
	mm, okm := fieldDuration(period.minutes, int64(time.Minute))
	ss, oks := fieldDuration(period.seconds, int64(time.Second))
	total, okt := sumTotal(hh, mm, ss)
	return time.Duration(total), okh && okm && oks && okt
}

// sumTotal adds the values, returning false if the total is outside the int64 range.
func sumTotal(values ...int64) (int64, bool) {
	var total int64
	for _, v := range values {
		sum := total + v
		if (v > 0 && sum < total) || (v < 0 && sum > total) {
			return 0, false
		}
		total = sum
	}
	return total, true
}

func fieldDuration(field decimal.Decimal, factor int64) (int64, bool) {
	if field.Coef() == 0 {
		return 0, true
	}

	for i := field.Scale(); i > 0; i-- {
		factor /= 10
	}
	if factor <= 0 {
		return 0, false
	}

	coef := field.Coef()
	if coef > uint64(math.MaxInt64)/uint64(factor) {
		return 0, false
	}

	return int64(field.Sign()) * int64(coef) * factor, true
}

func wholeCalendarValues(period Period) bool {
	zeroCalValue := zeroCalendarValues(period)
	wholeYears := period.years.Scale() == 0
	wholeMonths := period.months.Scale() == 0
	wholeWeeks := period.weeks.Scale() == 0
	wholeDays := period.days.Scale() == 0

	return !zeroCalValue && wholeYears && wholeMonths && wholeWeeks && wholeDays
}

func zeroCalendarValues(period Period) bool {
	zeroYears := period.years == decimal.Zero
	zeroMonths := period.months == decimal.Zero
	zeroWeeks := period.weeks == decimal.Zero
	zeroDays := period.days == decimal.Zero

	return zeroYears && zeroMonths && zeroWeeks && zeroDays
}

const (
	secondsPerDay = 24 * 60 * 60 // assuming 24-hour day

	daysPerYearE6  = 365242500          // 365.2425 days by the Gregorian rule
	daysPerMonthE6 = daysPerYearE6 / 12 // 30.436875 days per month
)
