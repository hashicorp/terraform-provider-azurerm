package validate

import (
	"testing"
)

func TestServiceEndpointPolicyAlias(t *testing.T) {
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
			Name:  "NoLeadingSlash",
			Input: "services/",
			Error: true,
		},
		{
			Name:  "NoTrailingSlash",
			Input: "/services",
			Error: true,
		},
		{
			Name:  "One word lone of valid char",
			Input: "a",
			Error: true,
		},
		{
			Name:  "Expected Input",
			Input: "/services/Azure/MachineLearning",
			Error: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %s: %q", v.Name, v.Input)

		_, errors := SubnetServiceEndpointStoragePolicyAlias(v.Input, "")
		isError := len(errors) != 0
		if v.Error != isError {
			t.Fatalf("Expected %t but got %t", v.Error, isError)
		}
	}
}
