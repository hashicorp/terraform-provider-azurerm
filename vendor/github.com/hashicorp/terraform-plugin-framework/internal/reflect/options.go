// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package reflect

// Options provides configuration settings for how the reflection behavior
// works, letting callers tweak different behaviors based on their needs.
type Options struct {
	// UnhandledNullAsEmpty controls whether null values should be
	// translated into empty values without provider interaction, or if
	// they must be explicitly handled.
	UnhandledNullAsEmpty bool

	// UnhandledUnknownAsEmpty controls whether null values should be
	// translated into empty values without provider interaction, or if
	// they must be explicitly handled.
	UnhandledUnknownAsEmpty bool
}
