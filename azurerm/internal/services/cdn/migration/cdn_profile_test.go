package migration

import (
	"context"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestCdnProfileV1ToV2(t *testing.T) {
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
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourcegroups/group1/providers/Microsoft.Cdn/profiles/profile1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Cdn/profiles/profile1"),
		},
		{
			name: "new id",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Cdn/profiles/profile1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Cdn/profiles/profile1"),
		},
	}
	for _, test := range testData {
		t.Logf("Testing %q..", test.name)
		result, err := CdnProfileV0ToV1{}.UpgradeFunc()(context.TODO(), test.input, nil)
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
