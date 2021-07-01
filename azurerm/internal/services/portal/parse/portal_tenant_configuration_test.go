package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = PortalTenantConfigurationId{}

func TestPortalTenantConfigurationIDFormatter(t *testing.T) {
	actual := NewPortalTenantConfigurationID("default").ID()
	expected := "/providers/Microsoft.Portal/tenantConfigurations/default"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestPortalTenantConfigurationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *PortalTenantConfigurationId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing Name
			Input: "/providers/Microsoft.Portal/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/providers/Microsoft.Portal/tenantConfigurations/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/Microsoft.Portal/tenantConfigurations/default",
			Expected: &PortalTenantConfigurationId{
				Name: "default",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.PORTAL/TENANTCONFIGURATIONS/DEFAULT",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := PortalTenantConfigurationID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
