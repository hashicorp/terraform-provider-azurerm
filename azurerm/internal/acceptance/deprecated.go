package acceptance

import "os"

// Functions in this file are provided for compatibility purposes - but will eventually be deprecated
// Developers should instead switch to using the `BuildTestData` method to obtain Test/Client Data

func Location() string {
	return os.Getenv("ARM_TEST_LOCATION")
}

func AltLocation() string {
	return os.Getenv("ARM_TEST_LOCATION_ALT")
}

func AltLocation2() string {
	return os.Getenv("ARM_TEST_LOCATION_ALT2")
}
