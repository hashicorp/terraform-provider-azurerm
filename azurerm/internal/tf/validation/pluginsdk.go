package validation

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// This file is intended to provide a transition from Plugin SDKv1 to Plugin SDKv2
// without introducing a merge conflict into every PR.

// Any returns a SchemaValidateFunc which tests if the provided value
// passes any of the provided SchemaValidateFunc
func Any(validators ...schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	return validation.Any(validators...)
}

// IntBetween returns a SchemaValidateFunc which tests if the provided value
// is of type int and is between min and max (inclusive)
func IntBetween(min, max int) schema.SchemaValidateFunc {
	return validation.IntBetween(min, max)
}

// IsCIDR is a SchemaValidateFunc which tests if the provided value is of type string and a valid CIDR
func IsCIDR(i interface{}, k string) (warnings []string, errors []error) {
	return validation.IsCIDR(i, k)
}

// IsIPAddress is a SchemaValidateFunc which tests if the provided value is of type string and is a single IP (v4 or v6)
func IsIPAddress(i interface{}, k string) (warnings []string, errors []error) {
	return validation.IsIPAddress(i, k)
}

// IsPortNumber is a SchemaValidateFunc which tests if the provided value is of type string and a valid TCP Port Number
func IsPortNumber(i interface{}, k string) (warnings []string, errors []error) {
	return validation.IsPortNumber(i, k)
}

// StringIsBase64 is a ValidateFunc that ensures a string can be parsed as Base64
func StringIsBase64(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringIsBase64(i, k)
}

// StringInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type string and matches the value of an element in the valid slice
// will test with in lower case if ignoreCase is true
func StringInSlice(valid []string, ignoreCase bool) schema.SchemaValidateFunc {
	return validation.StringInSlice(valid, ignoreCase)
}

// StringIsNotEmpty is a ValidateFunc that ensures a string is not empty
func StringIsNotEmpty(i interface{}, k string) ([]string, []error) {
	return validation.StringIsNotEmpty(i, k)
}
