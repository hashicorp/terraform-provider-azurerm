package migration

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestCustomHttpsConfigurationV0ToV1(t *testing.T) {
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
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourcegroups/mygroup1/providers/Microsoft.Network/frontdoors/frontdoor1/frontendendpoints/exampleFrontendEndpoint2/customHttpsConfiguration/exampleFrontendEndpoint2",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/mygroup1/providers/Microsoft.Network/frontDoors/frontdoor1/frontendEndpoints/exampleFrontendEndpoint2"),
		},
		{
			name: "old id - mixed case",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourcegroups/mygroup1/providers/Microsoft.Network/Frontdoors/frontdoor1/frontendEndpoints/exampleFrontendEndpoint2/customHttpsConfiguration/exampleFrontendEndpoint2",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/mygroup1/providers/Microsoft.Network/frontDoors/frontdoor1/frontendEndpoints/exampleFrontendEndpoint2"),
		},
		{
			name: "new id",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/mygroup1/providers/Microsoft.Network/frontdoors/frontdoor1/frontendEndpoints/exampleFrontendEndpoint2",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/mygroup1/providers/Microsoft.Network/frontDoors/frontdoor1/frontendEndpoints/exampleFrontendEndpoint2"),
		},
	}
	for _, test := range testData {
		t.Logf("Testing %q..", test.name)
		result, err := CustomHttpsConfigurationV0ToV1(test.input, nil)
		if err != nil && test.expected == nil {
			continue
		} else {
			if err == nil && test.expected == nil {
				t.Fatalf("Expected an error but didn't get one")
			} else if err != nil && test.expected != nil {
				t.Fatalf("Expected no error but got: %+v", err)
			}
		}

		actualId := result["id"].(string)
		if *test.expected != actualId {
			t.Fatalf("expected %q but got %q!", *test.expected, actualId)
		}
	}
}
