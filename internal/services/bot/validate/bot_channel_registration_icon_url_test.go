package validate

import (
	"testing"
)

func TestBotChannelRegistrationIconUrl(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "test.png",
			Expected: true,
		},
		{
			Input:    "http://myicon.png",
			Expected: true,
		},
		{
			Input:    "test.jpg",
			Expected: false,
		},
	}

	for _, v := range testCases {
		_, errors := BotChannelRegistrationIconUrl(v.Input, "icon_url")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
