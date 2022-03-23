package features

// nolint gocritic
// DeprecatedInThreePointOh returns the deprecation message if the provider
// is running in 3.0 mode - otherwise is returns an empty string (such that
// this deprecation should be ignored).
//
// This will be used to signify resources which will be Deprecated in 3.0,
// but not Removed (which will happen in a later, presumably 4.x release).
func DeprecatedInThreePointOh(deprecationMessage string) string {
	return deprecationMessage
}

// ThreePointOh returns whether this provider is running in 3.0 mode
// that is to say - the final 3.0 release
func ThreePointOh() bool {
	return true
}

// ThreePointOhBeta returns whether this provider is running in 3.0 mode
// or an opt-in Beta of the (non-breaking changes) coming in 3.0.
func ThreePointOhBeta() bool {
	return true
}
