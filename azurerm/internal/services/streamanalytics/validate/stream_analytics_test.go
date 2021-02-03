package validate

import (
	"testing"
)

func TestStreamAnalyticsJobStreamingUnits(t *testing.T) {
	cases := map[int]bool{
		0:  false,
		1:  true,
		2:  false,
		3:  true,
		4:  false,
		5:  false,
		6:  true,
		7:  false,
		8:  false,
		9:  false,
		10: false,
		11: false,
		12: true,
		18: true,
		24: true,
		30: true,
	}
	for i, shouldBeValid := range cases {
		_, errors := StreamAnalyticsJobStreamingUnits(i, "streaming_units")

		isValid := len(errors) == 0
		if shouldBeValid != isValid {
			t.Fatalf("Expected %d to be %t but got %t", i, shouldBeValid, isValid)
		}
	}
}
