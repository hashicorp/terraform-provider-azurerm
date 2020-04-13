package location

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestNormalizeLocation(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "West US",
			expected: "westus",
		},
		{
			input:    "South East Asia",
			expected: "southeastasia",
		},
		{
			input:    "southeastasia",
			expected: "southeastasia",
		},
	}

	for _, v := range cases {
		actual := Normalize(v.input)
		if v.expected != actual {
			t.Fatalf("Expected %q but got %q", v.expected, actual)
		}
	}
}

func TestNormalizeNilableLocation(t *testing.T) {
	cases := []struct {
		input    *string
		expected string
	}{
		{
			input:    utils.String("West US"),
			expected: "westus",
		},
		{
			input:    utils.String("South East Asia"),
			expected: "southeastasia",
		},
		{
			input:    utils.String("southeastasia"),
			expected: "southeastasia",
		},
		{
			input:    nil,
			expected: "",
		},
	}

	for _, v := range cases {
		actual := NormalizeNilable(v.input)
		if v.expected != actual {
			t.Fatalf("Expected %q but got %q", v.expected, actual)
		}
	}
}
