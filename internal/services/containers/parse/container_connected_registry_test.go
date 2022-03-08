package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ContainerConnectedRegistryId{}

func TestContainerConnectedRegistryIDFormatter(t *testing.T) {
	actual := NewContainerConnectedRegistryID("12345678-1234-9876-4563-123456789012", "group1", "registry1", "registry1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/registry1/connectedRegistries/registry1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestContainerConnectedRegistryID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ContainerConnectedRegistryId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/",
			Error: true,
		},

		{
			// missing value for RegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/",
			Error: true,
		},

		{
			// missing ConnectedRegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/registry1/",
			Error: true,
		},

		{
			// missing value for ConnectedRegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/registry1/connectedRegistries/registry1",
			Expected: &ContainerConnectedRegistryId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "group1",
				RegistryName:          "registry1",
				ConnectedRegistryName: "registry1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/GROUP1/PROVIDERS/MICROSOFT.CONTAINERREGISTRY/REGISTRIES/REGISTRY1/CONNECTEDREGISTRIES/REGISTRY1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ContainerConnectedRegistryID(v.Input)
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
		if actual.ConnectedRegistryName != v.Expected.ConnectedRegistryName {
			t.Fatalf("Expected %q but got %q for ConnectedRegistryName", v.Expected.ConnectedRegistryName, actual.ConnectedRegistryName)
		}
	}
}
