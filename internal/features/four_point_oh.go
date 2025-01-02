// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package features

import (
	"os"
	"strings"
)

// nolint gocritic
// DeprecatedInFourPointOh returns the deprecation message if the provider
// is running in 4.0 mode - otherwise is returns an empty string (such that
// this deprecation should be ignored).
// This can be used for the following scenarios:
//   - Signify resources which will be Deprecated in 4.0, but not Removed (which will happen in a later release).
//   - For properties undergoing a rename, where the renamed property will only be introduced in the next release
func DeprecatedInFourPointOh(deprecationMessage string) string {
	if !FourPointOhBeta() {
		return ""
	}

	return deprecationMessage
}

// FourPointOh returns whether this provider is running in 4.0 mode
// that is to say - the final 4.0 release
//
// This exists to allow breaking changes to be piped through the provider
// during the development of 3.x until 4.0 is ready.
func FourPointOh() bool {
	return true
}

// FourPointOhBeta returns whether this provider is running in 4.0 mode
// which is an opt-in Beta of the changes coming in 4.0.
//
// This exists to allow breaking changes to be piped through the provider
// during the development of 3.x until 4.0 is ready.
//
// The environment variable `ARM_FOURPOINTZERO_BETA` has been added
// to facilitate testing. But it should be noted that
// `ARM_FOURPOINTZERO_BETA` is ** NOT READY FOR PUBLIC USE ** and
// ** SHOULD NOT BE SET IN PRODUCTION ENVIRONMENTS **
// Setting `ARM_FOURPOINTZERO_BETA` will cause irreversible changes
// to your state.
func FourPointOhBeta() bool {
	return FourPointOh() || strings.EqualFold(os.Getenv("ARM_FOURPOINTZERO_BETA"), "true")
}
