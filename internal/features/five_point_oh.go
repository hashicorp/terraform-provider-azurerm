// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package features

// nolint gocritic
// DeprecatedInFivePointOh returns the deprecation message if the provider
// is running in 4.0 mode - otherwise is returns an empty string (such that
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

// FourPointOh returns whether this provider is running in 5.0 mode
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
func FivePointOhBeta() bool {
	return FivePointOh() || false
}
