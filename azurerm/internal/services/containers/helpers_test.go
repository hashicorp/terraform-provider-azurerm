package containers_test

import (
	"testing"
)

func checkIfShouldRunTestsIndividually(t *testing.T) {
	// NOTE: leaving this around so we can remove this gradually without
	// causing merge conflicts on open PR's
	//
	// This is no longer necessary since we limit the concurrent tests
	// for this package at the CI level
}
