package validate

import (
	"strings"
	"testing"
)

func TestServiceEndpointPolicyName(t *testing.T) {
	testData := []struct {
		Name  string
		Input string
		Error bool
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "Longest",
			Input: strings.Repeat("a", 80),
			Error: false,
		},
		{
			Name:  "Too long",
			Input: strings.Repeat("a", 81),
			Error: true,
		},
		{
			Name:  "One word lone of valid char",
			Input: "a",
			Error: false,
		},
		{
			Name:  "One word lone of invalid char - '_'",
			Input: "_",
			Error: true,
		},
		{
			Name:  "Invalid ending char - '-'",
			Input: "a-",
			Error: true,
		},
		{
			Name:  "Invalid middle char",
			Input: "a%a",
			Error: true,
		},
		{
			Name:  "Valid weired name",
			Input: "1.-_",
			Error: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %s: %q", v.Name, v.Input)

		_, errors := SubnetServiceEndpointPolicyName(v.Input, "")
		isError := len(errors) != 0
		if v.Error != isError {
			t.Fatalf("Expected %t but got %t", v.Error, isError)
		}
	}
}
