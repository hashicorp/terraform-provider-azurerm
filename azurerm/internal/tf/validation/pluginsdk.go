package validation

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// This file is intended to provide a transition from Plugin SDKv1 to Plugin SDKv2
// without introducing a merge conflict into every PR.

// All returns a SchemaValidateFunc which tests if the provided value
// passes all provided SchemaValidateFunc
func All(validators ...schema.SchemaValidateFunc) schema.SchemaValidateFunc { //nolint:SA1019
	return validation.All(validators...)
}

// Any returns a SchemaValidateFunc which tests if the provided value
// passes any of the provided SchemaValidateFunc
func Any(validators ...schema.SchemaValidateFunc) schema.SchemaValidateFunc { //nolint:SA1019
	return validation.Any(validators...)
}

// FloatBetween returns a SchemaValidateFunc which tests if the provided value
// is of type float64 and is between min and max (inclusive).
func FloatBetween(min, max float64) func(interface{}, string) ([]string, []error) {
	return validation.FloatBetween(min, max)
}

// FloatAtLeast returns a SchemaValidateFunc which tests if the provided value
// is of type float and is at least min (inclusive)
func FloatAtLeast(min float64) func(interface{}, string) ([]string, []error) {
	return validation.FloatAtLeast(min)
}

// IntNotInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type int and matches the value of an element in the valid slice
func IntNotInSlice(valid []int) func(interface{}, string) ([]string, []error) {
	return validation.IntNotInSlice(valid)
}

// IntAtLeast returns a SchemaValidateFunc which tests if the provided value
// is of type int and is at least min (inclusive)
func IntAtLeast(min int) func(interface{}, string) ([]string, []error) {
	return validation.IntAtLeast(min)
}

// IntAtMost returns a SchemaValidateFunc which tests if the provided value
// is of type int and is at most max (inclusive)
func IntAtMost(max int) func(interface{}, string) ([]string, []error) {
	return validation.IntAtMost(max)
}

// IntBetween returns a SchemaValidateFunc which tests if the provided value
// is of type int and is between min and max (inclusive)
func IntBetween(min, max int) func(interface{}, string) ([]string, []error) {
	return validation.IntBetween(min, max)
}

// IntDivisibleBy returns a SchemaValidateFunc which tests if the provided value
// is of type int and is divisible by a given number
func IntDivisibleBy(divisor int) func(interface{}, string) ([]string, []error) {
	return validation.IntDivisibleBy(divisor)
}

// IntInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type int and matches the value of an element in the valid slice
func IntInSlice(valid []int) func(interface{}, string) ([]string, []error) {
	return validation.IntInSlice(valid)
}

// IsCIDR is a SchemaValidateFunc which tests if the provided value is of type string and a valid CIDR
func IsCIDR(i interface{}, k string) ([]string, []error) {
	return validation.IsCIDR(i, k)
}

// IsDayOfTheWeek id a SchemaValidateFunc which tests if the provided value is of type string and a valid english day of the week
func IsDayOfTheWeek(ignoreCase bool) func(interface{}, string) ([]string, []error) {
	return validation.IsDayOfTheWeek(ignoreCase)
}

// IsIPAddress is a SchemaValidateFunc which tests if the provided value is of type string and is a single IP (v4 or v6)
func IsIPAddress(i interface{}, k string) ([]string, []error) {
	return validation.IsIPAddress(i, k)
}

// IsIPv4Address is a SchemaValidateFunc which tests if the provided value is of type string and a valid IPv4 address
func IsIPv4Address(i interface{}, k string) ([]string, []error) {
	return validation.IsIPv4Address(i, k)
}

// IsIPv4Range is a SchemaValidateFunc which tests if the provided value is of type string, and in valid IP range
func IsIPv4Range(i interface{}, k string) ([]string, []error) {
	return validation.IsIPv4Range(i, k)
}

// IsIPv6Address is a SchemaValidateFunc which tests if the provided value is of type string and a valid IPv6 address
func IsIPv6Address(i interface{}, k string) ([]string, []error) {
	return validation.IsIPv6Address(i, k)
}

// IsMonth id a SchemaValidateFunc which tests if the provided value is of type string and a valid english month
func IsMonth(ignoreCase bool) func(interface{}, string) ([]string, []error) {
	return validation.IsMonth(ignoreCase)
}

// IsPortNumber is a SchemaValidateFunc which tests if the provided value is of type string and a valid TCP Port Number
func IsPortNumber(i interface{}, k string) ([]string, []error) {
	return validation.IsPortNumber(i, k)
}

// IsRFC3339Time is a SchemaValidateFunc which tests if the provided value is of type string and a valid RFC33349Time
func IsRFC3339Time(i interface{}, k string) ([]string, []error) {
	return validation.IsRFC3339Time(i, k)
}

// IsURLWithHTTPorHTTPS is a SchemaValidateFunc which tests if the provided value is of type string and a valid HTTP or HTTPS URL
func IsURLWithHTTPorHTTPS(i interface{}, k string) ([]string, []error) {
	return validation.IsURLWithHTTPorHTTPS(i, k)
}

// IsURLWithHTTPS is a SchemaValidateFunc which tests if the provided value is of type string and a valid HTTPS URL
func IsURLWithHTTPS(i interface{}, k string) ([]string, []error) {
	return validation.IsURLWithHTTPS(i, k)
}

// IsURLWithScheme is a SchemaValidateFunc which tests if the provided value is of type string and a valid URL with the provided schemas
func IsURLWithScheme(validSchemes []string) func(interface{}, string) ([]string, []error) {
	return validation.IsURLWithScheme(validSchemes)
}

// IsUUID is a ValidateFunc that ensures a string can be parsed as UUID
func IsUUID(i interface{}, k string) ([]string, []error) {
	return validation.IsUUID(i, k)
}

// None returns a SchemaValidateFunc which tests if the provided value
// returns errors for all of the provided SchemaValidateFunc
func None(validators map[string]func(interface{}, string) ([]string, []error)) func(interface{}, string) ([]string, []error) {
	return func(i interface{}, k string) ([]string, []error) {
		var allErrors []error
		var allWarnings []string
		for name, validator := range validators {
			validatorWarnings, validatorErrors := validator(i, k)
			if len(validatorWarnings) == 0 && len(validatorErrors) == 0 {
				allErrors = append(allErrors, fmt.Errorf("ID cannot be a %s", name))
			}
		}
		return allWarnings, allErrors
	}
}

// NoZeroValues is a SchemaValidateFunc which tests if the provided value is
// not a zero value. It's useful in situations where you want to catch
// explicit zero values on things like required fields during validation.
func NoZeroValues(i interface{}, k string) ([]string, []error) {
	return validation.NoZeroValues(i, k)
}

// StringDoesNotContainAny returns a SchemaValidateFunc which validates that the
// provided value does not contain any of the specified Unicode code points in chars.
func StringDoesNotContainAny(chars string) func(interface{}, string) ([]string, []error) {
	return validation.StringDoesNotContainAny(chars)
}

// StringInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type string and matches the value of an element in the valid slice
// will test with in lower case if ignoreCase is true
func StringInSlice(valid []string, ignoreCase bool) func(interface{}, string) ([]string, []error) {
	return func(i interface{}, k string) ([]string, []error) {
		return validation.StringInSlice(valid, ignoreCase)(i, k)
	}
}

// StringIsBase64 is a ValidateFunc that ensures a string can be parsed as Base64
func StringIsBase64(i interface{}, k string) ([]string, []error) {
	return validation.StringIsBase64(i, k)
}

// StringIsJSON is a SchemaValidateFunc which tests to make sure the supplied string is valid JSON.
func StringIsJSON(i interface{}, k string) ([]string, []error) {
	return validation.StringIsJSON(i, k)
}

// StringIsEmpty is a ValidateFunc that ensures a string has no characters
func StringIsEmpty(i interface{}, k string) ([]string, []error) {
	return validation.StringIsEmpty(i, k)
}

// StringIsNotEmpty is a ValidateFunc that ensures a string is not empty
func StringIsNotEmpty(i interface{}, k string) ([]string, []error) {
	return validation.StringIsNotEmpty(i, k)
}

// StringIsNotWhiteSpace is a ValidateFunc that ensures a string is not empty or consisting entirely of whitespace characters
func StringIsNotWhiteSpace(i interface{}, k string) ([]string, []error) {
	return validation.StringIsNotWhiteSpace(i, k)
}

// StringIsValidRegExp returns a SchemaValidateFunc which tests to make sure the supplied string is a valid regular expression.
func StringIsValidRegExp(i interface{}, k string) ([]string, []error) {
	return validation.StringIsValidRegExp(i, k)
}

// StringLenBetween returns a SchemaValidateFunc which tests if the provided value
// is of type string and has length between min and max (inclusive)
func StringLenBetween(min, max int) func(interface{}, string) ([]string, []error) {
	return validation.StringLenBetween(min, max)
}

// StringMatch returns a SchemaValidateFunc which tests if the provided value
// matches a given regexp. Optionally an error message can be provided to
// return something friendlier than "must match some globby regexp".
func StringMatch(r *regexp.Regexp, message string) func(interface{}, string) ([]string, []error) {
	return validation.StringMatch(r, message)
}

// StringNotInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type string and does not match the value of any element in the invalid slice
// will test with in lower case if ignoreCase is true
func StringNotInSlice(invalid []string, ignoreCase bool) func(interface{}, string) ([]string, []error) {
	return validation.StringNotInSlice(invalid, ignoreCase)
}
