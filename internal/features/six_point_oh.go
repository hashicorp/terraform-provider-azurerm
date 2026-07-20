// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package features

// DeprecatedInSixPointOh returns the deprecation message if the provider
// is running in 6.0 mode - otherwise returns an empty string (such that
// this deprecation should be ignored).
// This can be used for the following scenarios:
//   - Signify resources which will be Deprecated in 6.0, but not Removed (which will happen in a later release).
//   - For properties undergoing a rename, where the renamed property will only be introduced in the next release
func DeprecatedInSixPointOh(deprecationMessage string) string {
	if !SixPointOh() {
		return ""
	}

	return deprecationMessage
}

// SixPointOh returns whether this provider is running in 5.0 mode
// that is to say - the final 5.0 release
//
// This exists to allow breaking changes to be piped through the provider
// during the development of 5.x until 6.0 is ready.
// The environment variable `ARM_SIXPOINTZERO_BETA` has been added
// to facilitate testing. But it should be noted that
// `ARM_SIXPOINTZERO_BETA` is ** NOT READY FOR PUBLIC USE ** and
// ** SHOULD NOT BE SET IN PRODUCTION ENVIRONMENTS **
// Setting `ARM_SIXPOINTZERO_BETA` will cause irreversible changes
// to your state.
func SixPointOh() bool {
	return false // TODO - remove and uncomment after the release of 5.0
	// return strings.EqualFold(os.Getenv("ARM_SIXPOINTZERO_BETA"), "true")
}
