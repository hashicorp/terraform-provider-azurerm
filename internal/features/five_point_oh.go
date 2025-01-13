// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package features

import (
	"os"
	"strings"
)

// nolint gocritic
// DeprecatedInFivePointOh returns the deprecation message if the provider
// is running in 5.0 mode - otherwise returns an empty string (such that
// this deprecation should be ignored).
// This can be used for the following scenarios:
//   - Signify resources which will be Deprecated in 5.0, but not Removed (which will happen in a later release).
//   - For properties undergoing a rename, where the renamed property will only be introduced in the next release
func DeprecatedInFivePointOh(deprecationMessage string) string {
	if !FivePointOhBeta() {
		return ""
	}

	return deprecationMessage
}

// FivePointOh returns whether this provider is running in 5.0 mode
// that is to say - the final 5.0 release
//
// This exists to allow breaking changes to be piped through the provider
// during the development of 4.x until 5.0 is ready.
func FivePointOh() bool {
	return false
}

// FivePointOhBeta returns whether this provider is running in 5.0 mode
// which is an opt-in Beta of the changes coming in 5.0.
//
// This exists to allow breaking changes to be piped through the provider
// during the development of 4.x until 5.0 is ready.
//
// The environment variable `ARM_FIVEPOINTZERO_BETA` has been added
// to facilitate testing. But it should be noted that
// `ARM_FIVEPOINTZERO_BETA` is ** NOT READY FOR PUBLIC USE ** and
// ** SHOULD NOT BE SET IN PRODUCTION ENVIRONMENTS **
// Setting `ARM_FIVEPOINTZERO_BETA` will cause irreversible changes
// to your state.
func FivePointOhBeta() bool {
	return FivePointOh() || strings.EqualFold(os.Getenv("ARM_FIVEPOINTZERO_BETA"), "true")
}
