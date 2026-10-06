# period

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rickb777/period)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/period)](https://goreportcard.com/report/github.com/rickb777/period)
[![Build](https://github.com/rickb777/period/actions/workflows/go.yml/badge.svg)](https://github.com/rickb777/period/actions)
[![Coverage](https://coveralls.io/repos/github/rickb777/period/badge.svg?branch=main)](https://coveralls.io/github/rickb777/period?branch=main)
[![Issues](https://img.shields.io/github/issues/rickb777/period.svg)](https://github.com/rickb777/period/issues)

Package `period` has types that represent ISO-8601 periods of time.

The two core types are 

 * `ISOString` - an ISO-8601 string such as "P2.5Y"
 * `Period` - a struct with the seven numbers: years, months, weeks, days, hours, minutes and seconds.

These two can be converted to the other.

`Period` also allows various calculations to be made. Its fields each hold up to 19 digits precision.

## Status

The API is now stable for v1.

## See Also

* [date](https://pkg.go.dev/github.com/rickb777/date) handles date and wall clock representations, along with date ranges and time spans.

## Upgrading

The old version of this was `github.com/rickb777/date/period`, which had very limited number range and used fixed-point arithmetic.

The new version, here, depends instead on `github.com/govalues/decimal`, which gives a huge (but finite) number range. There is now a 'weeks' field, which the old version did not have (it followed `time.Time` API patterns, which don't have weeks).

These functions have changed:

 * `New` now needs one more input parameter for the weeks field (7 parameters in total)
 * `NewYMD` still exists; there is also `NewYMWD`, which will often be more appropriate.

These methods have changed:

 * `YearsDecimal`, `MonthsDecimal`, `WeeksDecimal`, `DaysDecimal`, `HoursDecimal`, `MinutesDecimal` and `SecondsDecimal` return the fields as `decimal.Decimal`. They replace the old `YearsFloat`, `MonthsFloat`, `DaysFloat`, `HoursFloat`, `MinutesFloat` and `SecondsFloat` methods. `Years`, `Months`, `Weeks`, `Days`, `Hours`, `Minutes` and `Seconds` still return `int` as before.
 * `DaysIncWeeks` and `DaysIncWeeksDecimal` were added to return d + w * 7, which provides the behaviour similar to the old `Days` and `DaysFloat` methods. 
 * The old `ModuloDays` was dropped now that weeks are implemented fully. 
 * `OnlyYMD` is now `OnlyYMWD`
 * `Scale` and `ScaleWithOverflowCheck` have been replaced with `Mul`, which returns the multiplication product and a possible `error`.
