package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = ContainerRegistryScopeMapId{}

func TestContainerRegistryScopeMapIDFormatter(t *testing.T) {
	actual := NewContainerRegistryScopeMapID("12345678-1234-9876-4563-123456789012", "resGroup1", "registry1", "scopeMap1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/registries/registry1/scopeMaps/scopeMap1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestContainerRegistryScopeMapID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ContainerRegistryScopeMapId
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
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing RegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/",
			Error: true,
		},

		{
			// missing value for RegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/registries/",
			Error: true,
		},

		{
			// missing ScopeMapName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/registries/registry1/",
			Error: true,
		},

		{
			// missing value for ScopeMapName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/registries/registry1/scopeMaps/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/registries/registry1/scopeMaps/scopeMap1",
			Expected: &ContainerRegistryScopeMapId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "resGroup1",
				RegistryName:   "registry1",
				ScopeMapName:   "scopeMap1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.CONTAINERREGISTRY/REGISTRIES/REGISTRY1/SCOPEMAPS/SCOPEMAP1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ContainerRegistryScopeMapID(v.Input)
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
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.RegistryName != v.Expected.RegistryName {
			t.Fatalf("Expected %q but got %q for RegistryName", v.Expected.RegistryName, actual.RegistryName)
		}
		if actual.ScopeMapName != v.Expected.ScopeMapName {
			t.Fatalf("Expected %q but got %q for ScopeMapName", v.Expected.ScopeMapName, actual.ScopeMapName)
		}
	}
}
