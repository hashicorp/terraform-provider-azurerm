package features

import (
	"os"
	"strings"
)

// UseDynamicTestLocations returns whether or not the Acceptance Test data should use
// dynamic values for test locations
//
// In practice this means that of the available test locations, the primary, secondary
// and ternary locations will change dynamically for each test
//
// The primary benefit of this is to distribute the tests across Azure Regions,
// to improve the overall reliability. In time this'll become the default value.
//
// It's possible to opt into this by setting `ARM_PROVIDER_DYNAMIC_TEST` to `true`.
func UseDynamicTestLocations() bool {
	return strings.EqualFold(os.Getenv("ARM_PROVIDER_DYNAMIC_TEST"), "true")
}
