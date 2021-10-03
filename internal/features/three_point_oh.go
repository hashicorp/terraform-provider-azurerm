package features

import (
	"os"
	"strings"
)

// nolint gocritic
// DeprecatedInThreePointOh returns the deprecation message if the provider
// is running in 3.0 mode - otherwise is returns an empty string (such that
// this deprecation should be ignored).
//
// This will be used to signify resources which will be Deprecated in 3.0,
// but not Removed (which will happen in a later, presumably 4.x release).
func DeprecatedInThreePointOh(deprecationMessage string) string {
	if !ThreePointOh() {
		return ""
	}

	return deprecationMessage
}

// ThreePointOh returns whether this provider is running in 3.0 mode
// that is to say - that functionality which requires/changes in 3.0
// should be conditionally toggled on
//
// At this point in time this exists just to be able to place this
// infrastructure as required - but in time we'll flip this through
// a Beta and then GA at 3.0 release.
func ThreePointOh() bool {
	return false
}

// ThreePointOhBetaResources returns whether this provider is opted into
// the Beta Resources coming in v3.0 - or explicitly opted into v3.0.
func ThreePointOhBetaResources() bool {
	if ThreePointOh() {
		return true
	}

	return strings.EqualFold(os.Getenv("ARM_THREEPOINTZERO_BETA_RESOURCES"), "true")
}
