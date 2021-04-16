package validation

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// This file is intended to provide a transition from Plugin SDKv1 to Plugin SDKv2
// without introducing a merge conflict into every PR.

// All returns a SchemaValidateFunc which tests if the provided value
// passes all provided SchemaValidateFunc
func All(validators ...schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	return validation.All(validators...)
}

// Any returns a SchemaValidateFunc which tests if the provided value
// passes any of the provided SchemaValidateFunc
func Any(validators ...schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	return validation.Any(validators...)
}

// FloatAtLeast returns a SchemaValidateFunc which tests if the provided value
// is of type float and is at least min (inclusive)
func FloatAtLeast(min float64) schema.SchemaValidateFunc {
	return validation.FloatAtLeast(min)
}

// IntAtLeast returns a SchemaValidateFunc which tests if the provided value
// is of type int and is at least min (inclusive)
func IntAtLeast(min int) schema.SchemaValidateFunc {
	return validation.IntAtLeast(min)
}

// IntAtMost returns a SchemaValidateFunc which tests if the provided value
// is of type int and is at most max (inclusive)
func IntAtMost(max int) schema.SchemaValidateFunc {
	return validation.IntAtMost(max)
}

// IntBetween returns a SchemaValidateFunc which tests if the provided value
// is of type int and is between min and max (inclusive)
func IntBetween(min, max int) schema.SchemaValidateFunc {
	return validation.IntBetween(min, max)
}

// IntInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type int and matches the value of an element in the valid slice
func IntInSlice(valid []int) schema.SchemaValidateFunc {
	return validation.IntInSlice(valid)
}

// IsCIDR is a SchemaValidateFunc which tests if the provided value is of type string and a valid CIDR
func IsCIDR(i interface{}, k string) ([]string, []error) {
	return validation.IsCIDR(i, k)
}

// IsDayOfTheWeek id a SchemaValidateFunc which tests if the provided value is of type string and a valid english day of the week
func IsDayOfTheWeek(ignoreCase bool) schema.SchemaValidateFunc {
	return validation.IsDayOfTheWeek(ignoreCase)
}

// IsIPAddress is a SchemaValidateFunc which tests if the provided value is of type string and is a single IP (v4 or v6)
func IsIPAddress(i interface{}, k string) ([]string, []error) {
	return validation.IsIPAddress(i, k)
}

// IsIPv4Address is a SchemaValidateFunc which tests if the provided value is of type string and a valid IPv4 address
func IsIPv4Address(i interface{}, k string) (warnings []string, errors []error) {
	return validation.IsIPv4Address(i, k)
}

// IsPortNumber is a SchemaValidateFunc which tests if the provided value is of type string and a valid TCP Port Number
func IsPortNumber(i interface{}, k string) ([]string, []error) {
	return validation.IsPortNumber(i, k)
}

// IsURLWithScheme is a SchemaValidateFunc which tests if the provided value is of type string and a valid URL with the provided schemas
func IsURLWithScheme(validSchemes []string) schema.SchemaValidateFunc {
	return validation.IsURLWithScheme(validSchemes)
}

// IsUUID is a ValidateFunc that ensures a string can be parsed as UUID
func IsUUID(i interface{}, k string) (warnings []string, errors []error) {
	return validation.IsUUID(i, k)
}

// NoZeroValues is a SchemaValidateFunc which tests if the provided value is
// not a zero value. It's useful in situations where you want to catch
// explicit zero values on things like required fields during validation.
func NoZeroValues(i interface{}, k string) ([]string, []error) {
	return validation.NoZeroValues(i, k)
}

// StringInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type string and matches the value of an element in the valid slice
// will test with in lower case if ignoreCase is true
func StringInSlice(valid []string, ignoreCase bool) schema.SchemaValidateFunc {
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

// StringIsNotEmpty is a ValidateFunc that ensures a string is not empty
func StringIsNotEmpty(i interface{}, k string) ([]string, []error) {
	return validation.StringIsNotEmpty(i, k)
}

// StringIsValidRegExp returns a SchemaValidateFunc which tests to make sure the supplied string is a valid regular expression.
func StringIsValidRegExp(i interface{}, k string) ([]string, []error) {
	return validation.StringIsValidRegExp(i, k)
}

// StringMatch returns a SchemaValidateFunc which tests if the provided value
// matches a given regexp. Optionally an error message can be provided to
// return something friendlier than "must match some globby regexp".
func StringMatch(r *regexp.Regexp, message string) schema.SchemaValidateFunc {
	return validation.StringMatch(r, message)
}

// StringLenBetween returns a SchemaValidateFunc which tests if the provided value
// is of type string and has length between min and max (inclusive)
func StringLenBetween(min, max int) schema.SchemaValidateFunc {
	return validation.StringLenBetween(min, max)
}
