package validate

import (
	"testing"
)

func TestHealthbotBotsName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "valid123",
			Expected: true,
		},
		{
			Input:    "_invalid123",
			Expected: false,
		},
	}
	for _, v := range testCases {
		_, errors := HealthbotName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
