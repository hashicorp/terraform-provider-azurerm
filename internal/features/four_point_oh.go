// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package features

// nolint gocritic
// DeprecatedInFourPointOh returns the deprecation message if the provider
// is running in 4.0 mode - otherwise is returns an empty string (such that
// this deprecation should be ignored).
//
// This will be used to signify resources which will be Deprecated in 4.0,
// but not Removed (which will happen in a later, presumably 5.0 release).
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
	return false
}

// FourPointOhBeta returns whether this provider is running in 4.0 mode
// which is an opt-in Beta of the (non-breaking changes) coming in 4.0.
//
// Any features behind this flag should be backwards-compatible to allow
// users to try out 4.0 functionality.
//
// This exists to allow breaking changes to be piped through the provider
// during the development of 3.x until 4.0 is ready.
func FourPointOhBeta() bool {
	return FourPointOh() || false
}
