package tests

import (
	"os"
	"testing"
)

func checkIfShouldRunTestsCombined(t *testing.T) {
	if os.Getenv("TF_PROVIDER_SPLIT_COMBINED_TESTS") != "" {
		t.Skip("Skipping since this is being run Individually")
	}
}

func checkIfShouldRunTestsIndividually(t *testing.T) {
	if os.Getenv("TF_PROVIDER_SPLIT_COMBINED_TESTS") == "" {
		t.Skip("Skipping since this is being run as a Combined Test")
	}
}
