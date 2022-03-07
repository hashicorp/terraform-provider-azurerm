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
	if !ThreePointOhBeta() {
		return ""
	}

	return deprecationMessage
}

// ThreePointOh returns whether this provider is running in 3.0 mode
// that is to say - the final 3.0 release
//
// This exists to allow breaking changes to be piped through the provider
// during the development of 2.x until 3.0 is ready.
func ThreePointOh() bool {
	return false
}

// ThreePointOhBeta returns whether this provider is running in 3.0 mode
// which is an opt-in Beta of the (non-breaking changes) coming in 3.0.
//
// Any features behind this flag should be backwards-compatible to allow
// users to try out 3.0 functionality.
//
// This flag can be controlled by setting the Environment Variable
// `ARM_THREEPOINTZERO_BETA` to `true`.
func ThreePointOhBeta() bool {
	return ThreePointOh() || strings.EqualFold(os.Getenv("ARM_THREEPOINTZERO_BETA"), "true")
}

// ThreePointOhAppServiceResources returns whether this provider is opted into
// the Beta Resources coming in v3.0 - or explicitly opted into v3.0.
func ThreePointOhAppServiceResources() bool {
	if ThreePointOh() || ThreePointOhBeta() {
		return true
	}

	return strings.EqualFold(os.Getenv("ARM_THREEPOINTZERO_BETA_RESOURCES"), "true")
}
