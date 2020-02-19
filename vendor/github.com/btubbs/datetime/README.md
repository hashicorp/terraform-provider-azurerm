# datetime [![Build Status](https://travis-ci.org/btubbs/datetime.svg?branch=master)](https://travis-ci.org/btubbs/datetime) [![Coverage Status](https://coveralls.io/repos/github/btubbs/datetime/badge.svg?branch=master)](https://coveralls.io/github/btubbs/datetime?branch=master)

`datetime` provides a Parse function for turning commonly-used 
[ISO 8601](https://www.iso.org/iso-8601-date-and-time-format.html) date/time formats into
Golang time.Time variables.  `datetime.Parse` takes two arguments:

- the string you want to parse
- the timezone location to be used if there's not one specified inside the string

Unlike Go's built-in RFC-3339 time format, this package automatically supports ISO 8601 date and
time stamps with varying levels of granularity.  Examples:

```go
package main

import (
	"fmt"
	"time"

	"github.com/btubbs/datetime"
)

func main() {
	// just a year, defaulting to the time.UTC timezone
	fmt.Println(datetime.Parse("2007", time.UTC)) // 2007-01-01 00:00:00 +0000 UTC <nil>

	// a year and a month, this time defaulting to time.Local timezone
	fmt.Println(datetime.Parse("2007-11", time.Local)) // 2007-11-01 00:00:00 -0600 MDT <nil>

	// a full date
	fmt.Println(datetime.Parse("2007-11-22", time.UTC)) // 2007-11-22 00:00:00 +0000 UTC <nil>

	// adding time
	fmt.Println(datetime.Parse("2007-11-22T12:30:22", time.UTC)) // 2007-11-22 12:30:22 -0700 MST <nil>

	// fractions of a second
	fmt.Println(datetime.Parse("2007-11-22T12:30:22.321", time.UTC)) // 2007-11-22 12:30:22.321 -0700 MST <nil>

	// omitting dashes and colons, as ISO 8601 allows
	fmt.Println(datetime.Parse("20071122T123022", time.UTC)) // 2007-11-22 12:30:22 -0700 MST <nil>

	// a timezone offset inside the input will override the default provided to datetime.Parse
	fmt.Println(datetime.Parse("2007-11-22T12:30:22+0800", time.Local)) // 2007-11-22 12:30:22 +0800 +0800 <nil>

	// adding separators to the offset too
	fmt.Println(datetime.Parse("2007-11-22T12:30:22+08:00", time.UTC)) // 2007-11-22 12:30:22 +0800 +08:00 <nil>

	// using a shorthand for UTC
	fmt.Println(datetime.Parse("2007-11-22T12:30:22Z", time.Local)) // 2007-11-22 12:30:22 +0000 UTC <nil>
}
```

`DefaultUTC` and `DefaultLocal` types are also provided.  Used as struct fields, their Scan, Value,
and UnmarshalJSON methods support easy parsing of ISO 8601 timestamps from external systems.
