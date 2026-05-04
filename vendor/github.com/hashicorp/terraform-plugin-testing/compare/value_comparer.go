// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compare

// ValueComparer defines an interface that is implemented to run comparison logic on multiple values. Individual
// implementations determine how the comparison is performed (e.g., values differ, values equal).
type ValueComparer interface {
	// CompareValues should assert the given known values against any expectations.
	// Values are always ordered in the order they were added. Use the error
	// return to signal unexpected values or implementation errors.
	CompareValues(values ...any) error
}
