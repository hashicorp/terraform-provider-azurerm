package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = TenantTemplateDeploymentId{}

func TestTenantTemplateDeploymentIDFormatter(t *testing.T) {
	actual := NewTenantTemplateDeploymentID("deploy1").ID()
	expected := "/providers/Microsoft.Resources/deployments/deploy1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestTenantTemplateDeploymentID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *TenantTemplateDeploymentId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing DeploymentName
			Input: "/providers/Microsoft.Resources/",
			Error: true,
		},

		{
			// missing value for DeploymentName
			Input: "/providers/Microsoft.Resources/deployments/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/Microsoft.Resources/deployments/deploy1",
			Expected: &TenantTemplateDeploymentId{
				DeploymentName: "deploy1",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.RESOURCES/DEPLOYMENTS/DEPLOY1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := TenantTemplateDeploymentID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.DeploymentName != v.Expected.DeploymentName {
			t.Fatalf("Expected %q but got %q for DeploymentName", v.Expected.DeploymentName, actual.DeploymentName)
		}
	}
}
