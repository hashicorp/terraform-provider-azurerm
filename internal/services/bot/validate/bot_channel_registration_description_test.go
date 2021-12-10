package validate

import (
	"strings"
	"testing"
)

func TestBotChannelRegistrationDescription(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "Test123",
			Expected: true,
		},
		{
			Input:    strings.Repeat("t", 511),
			Expected: true,
		},
		{
			Input:    strings.Repeat("t", 512),
			Expected: true,
		},
		{
			Input:    strings.Repeat("t", 513),
			Expected: false,
		},
	}

	for _, v := range testCases {
		_, errors := BotChannelRegistrationDescription(v.Input, "description")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
