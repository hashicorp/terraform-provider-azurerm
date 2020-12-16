package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = SubscriptionTemplateDeploymentId{}

func TestSubscriptionTemplateDeploymentIDFormatter(t *testing.T) {
	actual := NewSubscriptionTemplateDeploymentID("12345678-1234-9876-4563-123456789012", "deploy1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Resources/deployments/deploy1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSubscriptionTemplateDeploymentID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SubscriptionTemplateDeploymentId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing DeploymentName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Resources/",
			Error: true,
		},

		{
			// missing value for DeploymentName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Resources/deployments/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Resources/deployments/deploy1",
			Expected: &SubscriptionTemplateDeploymentId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				DeploymentName: "deploy1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.RESOURCES/DEPLOYMENTS/DEPLOY1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SubscriptionTemplateDeploymentID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.DeploymentName != v.Expected.DeploymentName {
			t.Fatalf("Expected %q but got %q for DeploymentName", v.Expected.DeploymentName, actual.DeploymentName)
		}
	}
}
