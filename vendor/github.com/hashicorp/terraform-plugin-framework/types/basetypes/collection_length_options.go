// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package basetypes

// CollectionLengthOptions is a collection of toggles to control the behavior
// of the Length method on collection types (List, Set, Map, Tuple).
type CollectionLengthOptions struct {
	// UnhandledNullAsZero controls what happens when Length is called on a
	// null value. When set to true, zero will be returned. When set to false,
	// a panic will occur.
	UnhandledNullAsZero bool

	// UnhandledUnknownAsZero controls what happens when Length is called on
	// an unknown value. When set to true, zero will be returned. When set to
	// false, a panic will occur.
	UnhandledUnknownAsZero bool
}
