package azure_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
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
		actual := azure.NormalizeLocation(v.input)
		if v.expected != actual {
			t.Fatalf("Expected %q but got %q", v.expected, actual)
		}
	}
}
