// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"regexp"
)

var _ Check = stringRegexp{}

type stringRegexp struct {
	regex *regexp.Regexp
}

// CheckValue determines whether the passed value is of type string, and
// contains a sequence of bytes that match the regular expression supplied
// to StringRegexp.
func (v stringRegexp) CheckValue(other any) error {
	otherVal, ok := other.(string)

	if !ok {
		return fmt.Errorf("expected string value for StringRegexp check, got: %T", other)
	}

	if !v.regex.MatchString(otherVal) {
		return fmt.Errorf("expected regex match %s for StringRegexp check, got: %s", v.regex.String(), otherVal)
	}

	return nil
}

// String returns the string representation of the value.
func (v stringRegexp) String() string {
	return v.regex.String()
}

// StringRegexp returns a Check for asserting equality between the
// supplied regular expression and a value passed to the CheckValue method.
func StringRegexp(regex *regexp.Regexp) stringRegexp {
	return stringRegexp{
		regex: regex,
	}
}
