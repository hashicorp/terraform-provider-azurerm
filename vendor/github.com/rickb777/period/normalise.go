package period

import (
	"github.com/govalues/decimal"
)

var (
	seven       = decimal.MustNew(7, 0)
	twelve      = decimal.MustNew(12, 0)
	twentyFour  = decimal.MustNew(24, 0)
	sixty       = decimal.MustNew(60, 0)
	daysPerYear = decimal.MustNew(3652425, 4) // by the Gregorian rule
)

// Normalise simplifies the fields by propagating large values towards the more significant fields.
//
// Because the number of hours per day is imprecise (due to daylight savings etc), and because
// the number of days per month is variable in the Gregorian calendar, there is a reluctance
// to transfer time to or from the days element. To give control over this, there are two modes:
// it operates in either precise or approximate mode.
//
//   - Multiples of 60 seconds become minutes - both modes.
//   - Multiples of 60 minutes become hours - both modes.
//   - Multiples of 24 hours become days - approximate mode only
//   - Multiples of 7 days become weeks - both modes.
//   - Multiples of 12 months become years - both modes.
//
// Note that leap seconds are disregarded: every minute is assumed to have 60 seconds.
//
// If the calculations would lead to arithmetic errors, the current values are kept unaltered.
//
// See also [Period.NormaliseDaysToYears].
func (period Period) Normalise(precise bool) Period {
	// first phase - ripple large numbers to the left
	period.minutes, period.seconds = moveWholePartsLeft(period.minutes, period.seconds, sixty)
	period.hours, period.minutes = moveWholePartsLeft(period.hours, period.minutes, sixty)
	if !precise {
		period.days, period.hours = moveWholePartsLeft(period.days, period.hours, twentyFour)
	}
	period.weeks, period.days = moveWholePartsLeft(period.weeks, period.days, seven)
	period.years, period.months = moveWholePartsLeft(period.years, period.months, twelve)
	return period
}

// NormaliseDaysToYears tries to propagate large numbers of days (and corresponding weeks)
// to the years field. Based on the Gregorian rule, there are assumed to be 365.2425 days per year.
//
//   - Multiples of 365.2425 days become years
//
// If the calculations would lead to arithmetic errors, the current values are kept unaltered.
//
// A common use pattern would be to chain this after [Period.Normalise], i.e.
//
//	p.Normalise(false).NormaliseDaysToYears()
func (period Period) NormaliseDaysToYears() Period {
	if period.neg {
		return period.Negate().NormaliseDaysToYears().Negate()
	}

	days := period.DaysIncWeeksDecimal()

	// the threshold is one Gregorian year, not 366 days
	if days.Cmp(daysPerYear) < 0 {
		return period
	}

	ey, rem, err := days.QuoRem(daysPerYear)
	if err != nil {
		return period
	}

	period.years, err = period.years.Add(ey)
	if err != nil {
		return period
	}

	period.weeks, period.days = moveWholePartsLeft(decimal.Zero, rem.Trim(0), seven)
	return period
}

func moveWholePartsLeft(larger, smaller, nd decimal.Decimal) (decimal.Decimal, decimal.Decimal) {
	if smaller.IsZero() {
		return larger, smaller
	}

	q, r, err := smaller.QuoRem(nd)
	if err != nil {
		return larger, smaller
	}

	if !r.IsZero() && r.Prec() <= r.Scale() {
		return larger, smaller // more complex so no change
	}

	l2, err := larger.Add(q)
	if err != nil {
		return larger, smaller
	}

	return l2, r
}

//-------------------------------------------------------------------------------------------------

// SimplifyWeeksToDays adds 7 * the weeks field to the days field, and sets the weeks field to zero.
// See also [Period.SimplifyWeeks].
func (period Period) SimplifyWeeksToDays() Period {
	wdays, _ := period.weeks.Mul(seven)
	days, _ := wdays.Add(period.days)
	period.days = days
	period.weeks = decimal.Zero
	return period
}

// SimplifyWeeks adds 7 * the weeks field to the days field, and sets the weeks field to zero,
// but only if some other fields are non-zero.
//
// This will increase compatibility with external systems that do not expect to receive a weeks
// component unless the other components are zero. This is because ISO-8601 periods contain either
// weeks or other fields but not both.
//
// See also [Period.SimplifyWeeksToDays].
func (period Period) SimplifyWeeks() Period {
	if period.years.Coef() != 0 || period.months.Coef() != 0 || period.days.Coef() != 0 ||
		period.hours.Coef() != 0 || period.minutes.Coef() != 0 || period.seconds.Coef() != 0 {

		return period.SimplifyWeeksToDays()
	}
	return period
}

// Simplify simplifies the fields by propagating large values towards the less significant fields.
// This is akin to converting mixed fractions to improper fractions, across the group of fields.
// However, existing fields are not altered if they are a simple way of expressing their period already.
//
// For example, "P2Y1M" simplifies to "P25M" but "P2Y" remains "P2Y".
//
// Because the number of hours per day is imprecise (due to daylight savings etc), and because
// the number of days per month is variable in the Gregorian calendar, there is a reluctance
// to transfer time to or from the days element. To give control over this, there are two modes:
// it operates in either precise or approximate mode.
//
//   - Years may become multiples of 12 months if the number of months is non-zero - both modes.
//   - Weeks - see [Period.SimplifyWeeks] - both modes.
//   - Days may become multiples of 24 hours if the number of hours is non-zero - approximate mode only
//   - Hours may become multiples of 60 minutes if the number of minutes is non-zero - both modes.
//   - Minutes may become multiples of 60 seconds if the number of seconds is non-zero - both modes.
//
// If the calculations would lead to arithmetic errors, the current values are kept unaltered.
func (period Period) Simplify(precise bool) Period {
	period.years, period.months = moveToRight(period.years, period.months, twelve)
	p2 := period.SimplifyWeeks() // more thorough
	if !precise {
		p2.days, p2.hours = moveToRight(p2.days, p2.hours, twentyFour)
	}
	p2.hours, p2.minutes = moveToRight(p2.hours, p2.minutes, sixty)
	p2.minutes, p2.seconds = moveToRight(p2.minutes, p2.seconds, sixty)
	return p2
}

func moveToRight(larger, smaller, nd decimal.Decimal) (decimal.Decimal, decimal.Decimal) {
	if larger.IsZero() || isSimple(larger, smaller) {
		return larger, smaller
	}

	// first check whether it's actually simpler to keep things normalised
	lg1, sm1 := moveWholePartsLeft(larger, smaller, nd)
	if isSimple(lg1, sm1) {
		return lg1, sm1 // it's hard to beat this
	}

	extra, err := larger.Mul(nd)
	if err != nil {
		return larger, smaller
	}

	sm2, err := smaller.Add(extra)
	if err != nil {
		return larger, smaller
	}

	sm2 = sm2.Trim(0)

	originalDigits := larger.Prec() + smaller.Prec()
	if sm2.Prec() > originalDigits {
		return larger, smaller // because we would just add more digits
	}

	return decimal.Zero, sm2
}

func isSimple(larger, smaller decimal.Decimal) bool {
	return smaller.IsZero() && larger.Scale() == 0
}

// normaliseSign swaps the signs of all fields so that the largest non-zero field is positive and the overall sign
// indicates the original sign. Otherwise it has no effect.
func (period Period) normaliseSign() Period {
	if period.years.Sign() > 0 {
		return period
	} else if period.years.Sign() < 0 {
		return period.flipSign()
	}

	if period.months.Sign() > 0 {
		return period
	} else if period.months.Sign() < 0 {
		return period.flipSign()
	}

	if period.weeks.Sign() > 0 {
		return period
	} else if period.weeks.Sign() < 0 {
		return period.flipSign()
	}

	if period.days.Sign() > 0 {
		return period
	} else if period.days.Sign() < 0 {
		return period.flipSign()
	}

	if period.hours.Sign() > 0 {
		return period
	} else if period.hours.Sign() < 0 {
		return period.flipSign()
	}

	if period.minutes.Sign() > 0 {
		return period
	} else if period.minutes.Sign() < 0 {
		return period.flipSign()
	}

	if period.seconds.Sign() > 0 {
		return period
	} else if period.seconds.Sign() < 0 {
		return period.flipSign()
	}

	return Zero
}

func (period Period) flipSign() Period {
	period.neg = !period.neg
	return period.negateAllFields()
}

func (period Period) negateAllFields() Period {
	period.years = period.years.Neg()
	period.months = period.months.Neg()
	period.weeks = period.weeks.Neg()
	period.days = period.days.Neg()
	period.hours = period.hours.Neg()
	period.minutes = period.minutes.Neg()
	period.seconds = period.seconds.Neg()
	return period
}
