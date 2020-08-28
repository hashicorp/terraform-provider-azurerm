package migration

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAvailabilitySetV0ToV1(t *testing.T) {
	testData := []struct {
		name     string
		input    map[string]interface{}
		expected *string
	}{
		{
			name: "missing id",
			input: map[string]interface{}{
				"id": "",
			},
			expected: nil,
		},
		{
			name: "old id",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourcegroups/group1/providers/Microsoft.Compute/availabilitySets/set1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Compute/availabilitySets/set1"),
		},
		{
			name: "new id",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Compute/availabilitySets/set1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Compute/availabilitySets/set1"),
		},
	}

	for _, test := range testData {
		t.Logf("Testing %q...", test.name)
		result, err := AvailabilitySetV0ToV1(test.input, nil)
		if err != nil && test.expected == nil {
			continue
		}
		if err == nil && test.expected == nil {
			t.Fatal("expected an error but did not get one")
		}
		if err != nil && test.expected != nil {
			t.Fatalf("expected no error but got: %+v", err)
		}

		actualId := result["id"].(string)
		if *test.expected != actualId {
			t.Fatalf("expected %q but got %q", *test.expected, actualId)
		}
	}
}
