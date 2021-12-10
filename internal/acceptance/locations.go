package acceptance

import (
	"math/rand"
	"os"
	"time"
)

// Regions is a list of Azure Regions which can be used for test purposes
type Regions struct {
	// Primary is the Primary/Default Azure Region which should be used for testing
	Primary string

	// Secondary is the Secondary Azure Region which should be used for testing
	Secondary string

	// Ternary is the Ternary Azure Region which should be used for testing
	Ternary string
}

// availableLocations returns a struct containing a random set of regions
// this will return a randomly ordered set of locations - and as such must be cached
// this allows us to distribute the test suite across Azure to provide more stable tests
func availableLocations() Regions {
	locations := []string{
		os.Getenv("ARM_TEST_LOCATION"),
		os.Getenv("ARM_TEST_LOCATION_ALT"),
		os.Getenv("ARM_TEST_LOCATION_ALT2"),
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(locations), func(i, j int) {
		locations[i], locations[j] = locations[j], locations[i]
	})

	return Regions{
		Primary:   locations[0],
		Secondary: locations[1],
		Ternary:   locations[2],
	}
}
