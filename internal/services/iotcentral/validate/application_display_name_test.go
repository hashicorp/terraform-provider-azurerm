package validate

import "testing"

func TestApplicationDisplayName(t *testing.T) {
	testData := []struct {
		Value string
		Error bool
	}{
		{
			Value: "a",
			Error: false,
		},
		{
			Value: "A",
			Error: false,
		},
		{
			Value: "1",
			Error: false,
		},
		{
			Value: "1-",
			Error: false,
		},
		{
			Value: "a-",
			Error: false,
		},
		{
			Value: "A-",
			Error: false,
		},
		{
			Value: "a1-",
			Error: false,
		},
		{
			Value: "1a-",
			Error: false,
		},
		{
			Value: "aA-",
			Error: false,
		},
		{
			Value: "Aa-",
			Error: false,
		},
		{
			Value: "-",
			Error: false,
		},
		{
			Value: "-1",
			Error: false,
		},
		{
			Value: "_-a",
			Error: false,
		},
		{
			Value: "#$%$#!",
			Error: false,
		},
		{
			Value: "AA",
			Error: false,
		},
		{
			Value: "AA-1",
			Error: false,
		},
		{
			Value: "AA-a",
			Error: false,
		},
		{
			Value: "A1-",
			Error: false,
		},
		{
			Value: "AA-A",
			Error: false,
		},
		{
			Value: "AA-aA",
			Error: false,
		},
		{
			Value: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde1234",
			Error: false,
		},

		{
			Value: "",
			Error: true,
		},
		{
			Value: "adcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdssdavcadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdssdavcc",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Value)

		_, err := ApplicationDisplayName(v.Value, "unit test")
		if err != nil && !v.Error {
			t.Fatalf("Expected pass but got an error: %s", err)
		}
	}
}
