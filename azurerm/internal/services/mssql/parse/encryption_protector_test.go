package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = EncryptionProtectorId{}

func TestEncryptionProtectorIDFormatter(t *testing.T) {
	actual := NewEncryptionProtectorID("12345678-1234-9876-4563-123456789012", "group1", "server1", "current").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/encryptionProtector/current"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestEncryptionProtectorID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *EncryptionProtectorId
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
			// missing ServerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/",
			Error: true,
		},

		{
			// missing value for ServerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/",
			Error: true,
		},

		{
			// missing EncryptionProtectorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/",
			Error: true,
		},

		{
			// missing value for EncryptionProtectorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/encryptionProtector/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/encryptionProtector/current",
			Expected: &EncryptionProtectorId{
				SubscriptionId:          "12345678-1234-9876-4563-123456789012",
				ResourceGroup:           "group1",
				ServerName:              "server1",
				EncryptionProtectorName: "current",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/GROUP1/PROVIDERS/MICROSOFT.SQL/SERVERS/SERVER1/ENCRYPTIONPROTECTOR/CURRENT",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := EncryptionProtectorID(v.Input)
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
		if actual.ServerName != v.Expected.ServerName {
			t.Fatalf("Expected %q but got %q for ServerName", v.Expected.ServerName, actual.ServerName)
		}
		if actual.EncryptionProtectorName != v.Expected.EncryptionProtectorName {
			t.Fatalf("Expected %q but got %q for EncryptionProtectorName", v.Expected.EncryptionProtectorName, actual.EncryptionProtectorName)
		}
	}
}
