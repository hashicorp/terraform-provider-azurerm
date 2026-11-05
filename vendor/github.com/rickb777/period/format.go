// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"io"
	"strings"

	"github.com/govalues/decimal"
	"github.com/rickb777/plural/v2"
)

// Period converts the period to ISO-8601 string form.
// If there is a decimal fraction, it will be rendered using a decimal point separator
// (not a comma).
func (period Period) Period() ISOString {
	return ISOString(period.String())
}

// ISOString converts the period to ISO-8601 string form.
// If there is a decimal fraction, it will be rendered using a decimal point separator.
// (not a comma).
func (period Period) String() string {
	buf := &strings.Builder{}
	_, _ = period.WriteTo(buf)
	return buf.String()
}

// WriteTo converts the period to ISO-8601 form.
func (period Period) WriteTo(w io.Writer) (int64, error) {
	aw := adapt(w)

	if period == Zero {
		_, _ = aw.WriteString(string(CanonicalZero))
		return uwSum(aw)
	}

	if period.neg {
		_ = aw.WriteByte('-')
	}

	_ = aw.WriteByte('P')

	writeField(aw, period.years, Year)
	writeField(aw, period.months, Month)
	writeField(aw, period.weeks, Week)
	writeField(aw, period.days, Day)

	if period.hours.Coef() != 0 || period.minutes.Coef() != 0 || period.seconds.Coef() != 0 {
		_ = aw.WriteByte('T')

		writeField(aw, period.hours, Hour)
		writeField(aw, period.minutes, Minute)
		writeField(aw, period.seconds, Second)
	}

	return uwSum(aw)
}

func writeField(w usefulWriter, field decimal.Decimal, fieldDesignator Designator) {
	if field.Coef() != 0 {
		_, _ = w.WriteString(field.String())
		_ = w.WriteByte(fieldDesignator.Byte())
	}
}

//-------------------------------------------------------------------------------------------------

// Format converts the period to human-readable form using [DefaultFormatLocalisation].
// To adjust the result, see the [Period.Normalise], [Period.NormaliseDaysToYears], [Period.Simplify] and [Period.SimplifyWeeksToDays] methods.
func (period Period) Format() string {
	return period.FormatLocalised(DefaultFormatLocalisation)
}

// FormatLocalised converts the period to human-readable form in a localisable way.
// To adjust the result, see the [Period.Normalise], [Period.NormaliseDaysToYears], [Period.Simplify] and [Period.SimplifyWeeksToDays] methods.
func (period Period) FormatLocalised(config FormatLocalisation) string {
	if period.IsZero() {
		return config.ZeroValue
	}

	parts := make([]string, 0, 7)

	parts = appendNonBlank(parts, formatField(period.years, config.Negate, config.YearNames))
	parts = appendNonBlank(parts, formatField(period.months, config.Negate, config.MonthNames))
	parts = appendNonBlank(parts, formatField(period.weeks, config.Negate, config.WeekNames))
	parts = appendNonBlank(parts, formatField(period.days, config.Negate, config.DayNames))
	parts = appendNonBlank(parts, formatField(period.hours, config.Negate, config.HourNames))
	parts = appendNonBlank(parts, formatField(period.minutes, config.Negate, config.MinuteNames))
	parts = appendNonBlank(parts, formatField(period.seconds, config.Negate, config.SecondNames))

	return strings.Join(parts, ", ")
}

func formatField(field decimal.Decimal, negate func(string) string, names plural.Cases) string {
	number, _ := field.Float64()
	if number < 0 {
		return negate(names.Continuous(-number))
	}
	return names.Continuous(number)
}

func appendNonBlank(parts []string, s string) []string {
	if s == "" {
		return parts
	}
	return append(parts, s)
}

type FormatLocalisation struct {
	// ZeroValue is the string that represents a zero period "P0D".
	ZeroValue string

	// Negate alters a format string when the value is negative.
	Negate func(string) string

	// The plurals provide the localised format names for each field of the period.
	// Each is a sequence of plural cases where the first match is used, otherwise the last one is used.
	// The last one must include a "%v" placeholder for the number.
	YearNames, MonthNames, WeekNames, DayNames plural.Cases
	HourNames, MinuteNames, SecondNames        plural.Cases
}

// DefaultFormatLocalisation provides the formatting strings needed to format Period values in vernacular English.
var DefaultFormatLocalisation = FormatLocalisation{
	ZeroValue: "zero",
	Negate:    func(s string) string { return "minus " + s },

	// YearNames provides the English default format names for the years part of the period.
	// This is a sequence of plurals where the first match is used, otherwise the last one is used.
	YearNames:   plural.FromOne("%v year", "%v years"),
	MonthNames:  plural.FromOne("%v month", "%v months"),
	WeekNames:   plural.FromOne("%v week", "%v weeks"),
	DayNames:    plural.FromOne("%v day", "%v days"),
	HourNames:   plural.FromOne("%v hour", "%v hours"),
	MinuteNames: plural.FromOne("%v minute", "%v minutes"),
	SecondNames: plural.FromOne("%v second", "%v seconds"),
}
