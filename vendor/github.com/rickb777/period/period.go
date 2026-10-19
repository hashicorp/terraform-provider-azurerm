// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"time"

	"github.com/govalues/decimal"
)

// Period holds a period of time as a set of decimal numbers, one for each field in the ISO-8601
// period.
//
// By conventional, all the fields should have the same sign. However, this is not restricted,
// so each field after the first non-zero field can be independently positive or negative.
// Sometimes this makes sense, e.g. "P1YT-1S" is one second less than one year.
//
// The precision is large: all fields are scaled decimals using int64 internally for calculations.
// The value of each field can have up to 19 digits (the range of int64), of which up to 19 digits
// can be a decimal fraction. So the range is much wider than that of time.Duration; be aware that
// periods more than 292 years or less than one nanosecond are outside the convertible range.
//
// For convenience, the method inputs and outputs use int.
//
// Fractions are supported on the least significant non-zero field only. It is an error for
// more-significant fields to have fractional values too.
//
// Instances are immutable.
type Period struct {
	years, months, weeks, days, hours, minutes, seconds decimal.Decimal

	// neg indicates a negative period, which negates all fields (even if they are already negative)
	neg bool
}

// Zero is the zero period.
var Zero = Period{}

//-------------------------------------------------------------------------------------------------

// NewYMD creates a simple period without any fractional parts. The fields are initialised verbatim
// without any normalisation; e.g. 12 months will not become 1 year. Use [Period.Normalise] if you need to.
//
// This function is equivalent to NewYMWD(years, months, 0, days)
func NewYMD(years, months, days int) Period {
	return NewYMWD(years, months, 0, days)
}

// NewYMWD creates a simple period without any fractional parts. The fields are initialised verbatim
// without any normalisation; e.g. 12 months will not become 1 year. Use [Period.Normalise] if you need to.
func NewYMWD(years, months, weeks, days int) Period {
	return New(years, months, weeks, days, 0, 0, 0)
}

// NewHMS creates a simple period without any fractional parts. The fields are initialised verbatim
// without any normalisation; e.g. 120 seconds will not become 2 minutes. Use [Period.Normalise] if you need to.
func NewHMS(hours, minutes, seconds int) Period {
	return New(0, 0, 0, 0, hours, minutes, seconds)
}

// New creates a simple period without any fractional parts. The fields are initialised verbatim
// without any normalisation; e.g. 120 seconds will not become 2 minutes. Use [Period.Normalise] if you need to.
func New(years, months, weeks, days, hours, minutes, seconds int) Period {
	return Period{
		years:   decimal.MustNew(int64(years), 0),
		months:  decimal.MustNew(int64(months), 0),
		weeks:   decimal.MustNew(int64(weeks), 0),
		days:    decimal.MustNew(int64(days), 0),
		hours:   decimal.MustNew(int64(hours), 0),
		minutes: decimal.MustNew(int64(minutes), 0),
		seconds: decimal.MustNew(int64(seconds), 0),
	}.normaliseSign()
}

// MustNewDecimal creates a period from seven decimal values. The fields are trimmed but no normalisation
// is applied, e.g. 120 seconds will not become 2 minutes. Use [Period.Normalise] if you need to.
//
// Periods only allow the least-significant non-zero field to contain a fraction. If any of the
// more-significant fields is supplied with a fraction, this function panics.
func MustNewDecimal(years, months, weeks, days, hours, minutes, seconds decimal.Decimal) Period {
	p, err := NewDecimal(years, months, weeks, days, hours, minutes, seconds)
	if err != nil {
		panic(err)
	}
	return p
}

// NewDecimal creates a period from seven decimal values. The fields are trimmed but no normalisation
// is applied, e.g. 120 seconds will not become 2 minutes. Use [Period.Normalise] if you need to.
//
// Periods only allow the least-significant non-zero field to contain a fraction. If any of the
// more-significant fields is supplied with a fraction, an error will be returned. This can be safely
// ignored for non-standard behaviour.
func NewDecimal(years, months, weeks, days, hours, minutes, seconds decimal.Decimal) (period Period, err error) {
	ymwd := make([]byte, 0, 4)
	hms := make([]byte, 0, 4)

	if years.Scale() > 0 {
		if months.Coef() != 0 || weeks.Coef() != 0 || days.Coef() != 0 || hours.Coef() != 0 || minutes.Coef() != 0 || seconds.Coef() != 0 {
			ymwd = append(ymwd, 'Y')
		}
	}

	if months.Scale() > 0 {
		if weeks.Coef() != 0 || days.Coef() != 0 || hours.Coef() != 0 || minutes.Coef() != 0 || seconds.Coef() != 0 {
			ymwd = append(ymwd, 'M')
		}
	}

	if weeks.Scale() > 0 {
		if days.Coef() != 0 || hours.Coef() != 0 || minutes.Coef() != 0 || seconds.Coef() != 0 {
			ymwd = append(ymwd, 'W')
		}
	}

	if days.Scale() > 0 {
		if hours.Coef() != 0 || minutes.Coef() != 0 || seconds.Coef() != 0 {
			ymwd = append(ymwd, 'D')
		}
	}

	if hours.Scale() > 0 {
		if minutes.Coef() != 0 || seconds.Coef() != 0 {
			hms = append(hms, 'H')
		}
	}

	if minutes.Scale() > 0 {
		if seconds.Coef() != 0 {
			hms = append(hms, 'M')
		}
	}

	if seconds.Scale() > 0 && len(ymwd)+len(hms) > 0 {
		hms = append(hms, 'S')
	}

	p := Period{
		years:   years.Trim(0),
		months:  months.Trim(0),
		weeks:   weeks.Trim(0),
		days:    days.Trim(0),
		hours:   hours.Trim(0),
		minutes: minutes.Trim(0),
		seconds: seconds.Trim(0),
	}.normaliseSign()

	if len(ymwd)+len(hms) > 0 {
		sep := ""
		if len(ymwd) > 0 && len(hms) > 0 {
			sep = "/"
		}
		err = fmt.Errorf("only the least significant field can have a fraction; found %s%s%s fractions in %s", string(ymwd), sep, string(hms), p)
	}

	return p, err
}

// NewOf converts a time duration to a Period.
// The result just a number of seconds, possibly including a fraction. It is not normalised; see [Period.Normalise].
func NewOf(duration time.Duration) Period {
	seconds := decimal.MustNew(int64(duration), 9).Trim(0)
	return Period{seconds: seconds}.normaliseSign()
}

//-------------------------------------------------------------------------------------------------

// Between converts the span between two times to a period. Based on the Gregorian conversion
// algorithms of `time.Time`, the resultant period is precise.
//
// If t2 is before t1, the result is a negative period.
//
// The result just a number of seconds, possibly including a fraction. It is not normalised; see [Period.Normalise].
//
// Remember that the resultant period does not retain any knowledge of the calendar, so any subsequent
// computations applied to the period can only be precise if they concern either the date (year, month,
// day) part, or the clock (hour, minute, second) part, but not both.
func Between(t1, t2 time.Time) Period {
	if t2.Before(t1) {
		return NewOf(t2.Sub(t1))
	}

	return NewOf(t1.Sub(t2)).Negate()
}

//-------------------------------------------------------------------------------------------------

// IsZero returns true if applied to a period of zero length.
func (period Period) IsZero() bool {
	period.neg = false
	return period == Zero
}

// Sign returns 1 if the period is positive, -1 if it is negative, or zero otherwise.
func (period Period) Sign() int {
	switch {
	case period.neg:
		return -1
	case period != Zero:
		return 1
	default:
		return 0
	}
}

// IsNegative returns true if the period is negative.
func (period Period) IsNegative() bool {
	return period.neg
}

// IsPositive returns true if the period is positive or zero.
func (period Period) IsPositive() bool {
	return !period.neg
}

// Abs converts a negative period to a positive period.
func (period Period) Abs() Period {
	period.neg = false
	return period
}

// Negate changes the sign of the period. Zero is not altered.
func (period Period) Negate() Period {
	if period.IsZero() {
		return Zero
	}
	period.neg = !period.neg
	return period
}

//-------------------------------------------------------------------------------------------------

// OnlyYMWD returns the period with only the year, month, week and day fields.
// The hour, minute and second fields are zeroed.
func (period Period) OnlyYMWD() Period {
	return Period{neg: period.neg, years: period.years, months: period.months, weeks: period.weeks, days: period.days}
}

// OnlyHMS returns the period with only the hour, minute and second fields.
// The year, month, week and day fields are zeroed.
func (period Period) OnlyHMS() Period {
	return Period{neg: period.neg, hours: period.hours, minutes: period.minutes, seconds: period.seconds}
}

//-------------------------------------------------------------------------------------------------

// Years gets the whole number of years in the period.
func (period Period) Years() int {
	i, _, _ := period.YearsDecimal().Int64(0)
	return int(i)
}

// YearsDecimal gets the number of years in the period, including any fraction present.
func (period Period) YearsDecimal() decimal.Decimal {
	return period.applySign(period.years)
}

// Months gets the whole number of months in the period.
func (period Period) Months() int {
	i, _, _ := period.MonthsDecimal().Int64(0)
	return int(i)
}

// MonthsDecimal gets the number of months in the period, including any fraction present.
func (period Period) MonthsDecimal() decimal.Decimal {
	return period.applySign(period.months)
}

// Weeks gets the whole number of weeks in the period.
func (period Period) Weeks() int {
	i, _, _ := period.WeeksDecimal().Int64(0)
	return int(i)
}

// WeeksDecimal gets the number of weeks in the period, including any fraction present.
func (period Period) WeeksDecimal() decimal.Decimal {
	return period.applySign(period.weeks)
}

// Days gets the whole number of days in the period.
func (period Period) Days() int {
	i, _, _ := period.DaysDecimal().Int64(0)
	return int(i)
}

// DaysDecimal gets the number of days in the period, including any fraction present.
func (period Period) DaysDecimal() decimal.Decimal {
	return period.applySign(period.days)
}

// DaysIncWeeks gets the whole number of days in the period, including all the weeks.
// The result is d + (w * 7), given d days and w weeks.
//
// See also [Period.SimplifyWeeksToDays].
func (period Period) DaysIncWeeks() int {
	i, _, _ := period.DaysIncWeeksDecimal().Int64(0)
	return int(i)
}

// DaysIncWeeksDecimal gets the number of days in the period, including all the weeks and including any
// fraction present. The result is d + (w * 7), given d days and w weeks.
//
// See also [Period.SimplifyWeeksToDays].
func (period Period) DaysIncWeeksDecimal() decimal.Decimal {
	wdays, _ := period.weeks.Mul(seven)
	days, _ := wdays.Add(period.days)
	return period.applySign(days)
}

// Hours gets the whole number of hours in the period.
func (period Period) Hours() int {
	i, _, _ := period.HoursDecimal().Int64(0)
	return int(i)
}

// HoursDecimal gets the number of hours in the period, including any fraction present.
func (period Period) HoursDecimal() decimal.Decimal {
	return period.applySign(period.hours)
}

// Minutes gets the whole number of minutes in the period.
func (period Period) Minutes() int {
	i, _, _ := period.MinutesDecimal().Int64(0)
	return int(i)
}

// MinutesDecimal gets the number of minutes in the period, including any fraction present.
func (period Period) MinutesDecimal() decimal.Decimal {
	return period.applySign(period.minutes)
}

// Seconds gets the whole number of seconds in the period.
func (period Period) Seconds() int {
	i, _, _ := period.SecondsDecimal().Int64(0)
	return int(i)
}

// SecondsDecimal gets the number of seconds in the period, including any fraction present.
func (period Period) SecondsDecimal() decimal.Decimal {
	return period.applySign(period.seconds)
}

func (period Period) applySign(field decimal.Decimal) decimal.Decimal {
	if period.neg {
		return field.Neg()
	}
	return field
}

//-------------------------------------------------------------------------------------------------

// GetInt gets one field as a whole number.
//
// A panic arises if the field is unknown.
func (period Period) GetInt(field Designator) int {
	value := period.GetField(field)
	i, _, _ := value.Int64(0)
	return int(i)
}

// GetField gets one field.
//
// A panic arises if the field is unknown.
func (period Period) GetField(field Designator) decimal.Decimal {
	switch field {
	case Year:
		return period.applySign(period.years)
	case Month:
		return period.applySign(period.months)
	case Week:
		return period.applySign(period.weeks)
	case Day:
		return period.applySign(period.days)
	case Hour:
		return period.applySign(period.hours)
	case Minute:
		return period.applySign(period.minutes)
	case Second:
		return period.applySign(period.seconds)
	}

	panic(field)
}

//-------------------------------------------------------------------------------------------------

// SetInt sets one field in the period as a whole number.
//
// A panic arises if the field is unknown.
func (period Period) SetInt(value int, field Designator) Period {
	d := decimal.MustNew(int64(value), 0)
	p, _ := period.SetField(d, field)
	return p
}

// SetField sets one field in the period. Like NewDecimal, an error arises if the new period
// would have multiple fields with fractions.
//
// A panic arises if the field is unknown.
func (period Period) SetField(value decimal.Decimal, field Designator) (Period, error) {
	switch field {
	case Year:
		return NewDecimal(value, period.months, period.weeks, period.days, period.hours, period.minutes, period.seconds)
	case Month:
		return NewDecimal(period.years, value, period.weeks, period.days, period.hours, period.minutes, period.seconds)
	case Week:
		return NewDecimal(period.years, period.months, value, period.days, period.hours, period.minutes, period.seconds)
	case Day:
		return NewDecimal(period.years, period.months, period.weeks, value, period.hours, period.minutes, period.seconds)
	case Hour:
		return NewDecimal(period.years, period.months, period.weeks, period.days, value, period.minutes, period.seconds)
	case Minute:
		return NewDecimal(period.years, period.months, period.weeks, period.days, period.hours, value, period.seconds)
	case Second:
		return NewDecimal(period.years, period.months, period.weeks, period.days, period.hours, period.minutes, value)
	}

	panic(field)
}
