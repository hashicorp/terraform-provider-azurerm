package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = ResourceGroupTemplateDeploymentId{}

func TestResourceGroupTemplateDeploymentIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	actual := NewResourceGroupTemplateDeploymentID("group1", "deploy1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Resources/deployments/deploy1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestResourceGroupTemplateDeploymentIDParser(t *testing.T) {
	testData := []struct {
		input    string
		expected *ResourceGroupTemplateDeploymentId
	}{
		{
			// camel case
			input: "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Resources/deployments/deploy1",
			expected: &ResourceGroupTemplateDeploymentId{
				ResourceGroup: "group1",
				Name:          "deploy1",
			},
		},
		{
			// title case
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Resources/Deployments/deploy1",
			expected: nil,
		},
	}
	for _, test := range testData {
		t.Logf("Testing %q..", test.input)
		actual, err := ResourceGroupTemplateDeploymentID(test.input)
		if err != nil && test.expected == nil {
			continue
		} else {
			if err == nil && test.expected == nil {
				t.Fatalf("Expected an error but didn't get one")
			} else if err != nil && test.expected != nil {
				t.Fatalf("Expected no error but got: %+v", err)
			}
		}

		if actual.ResourceGroup != test.expected.ResourceGroup {
			t.Fatalf("Expected ResourceGroup to be %q but was %q", test.expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Name != test.expected.Name {
			t.Fatalf("Expected name to be %q but was %q", test.expected.Name, actual.Name)
		}
	}
}
