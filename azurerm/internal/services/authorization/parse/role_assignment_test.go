package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = RoleAssignmentId{}

func TestRoleAssignmentIDFormatter(t *testing.T) {
	testData := []struct {
		SubscriptionId  string
		ResourceGroup   string
		ManagementGroup string
		Name            string
		Expected        string
	}{
		{
			SubscriptionId:  "",
			ResourceGroup:   "",
			ManagementGroup: "",
			Name:            "23456781-2349-8764-5631-234567890121",
		},
		{
			SubscriptionId:  "12345678-1234-9876-4563-123456789012",
			ResourceGroup:   "group1",
			ManagementGroup: "managementGroup1",
			Name:            "23456781-2349-8764-5631-234567890121",
		},
		{
			SubscriptionId:  "12345678-1234-9876-4563-123456789012",
			ResourceGroup:   "",
			ManagementGroup: "managementGroup1",
			Name:            "23456781-2349-8764-5631-234567890121",
		},
		{
			SubscriptionId:  "12345678-1234-9876-4563-123456789012",
			ResourceGroup:   "",
			ManagementGroup: "",
			Name:            "23456781-2349-8764-5631-234567890121",
			Expected:        "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
		},
		{
			SubscriptionId:  "12345678-1234-9876-4563-123456789012",
			ResourceGroup:   "group1",
			ManagementGroup: "",
			Name:            "23456781-2349-8764-5631-234567890121",
			Expected:        "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
		},
		{
			SubscriptionId:  "",
			ResourceGroup:   "",
			ManagementGroup: "12345678-1234-9876-4563-123456789012",
			Name:            "23456781-2349-8764-5631-234567890121",
			Expected:        "/providers/Microsoft.Management/managementGroups/12345678-1234-9876-4563-123456789012/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
		},
	}
	for _, v := range testData {
		t.Logf("testing %+v", v)
		actual, err := NewRoleAssignmentID(v.SubscriptionId, v.ResourceGroup, v.ManagementGroup, v.Name)
		if err != nil {
			if v.Expected == "" {
				continue
			}
		}
		actualId := actual.ID()
		if actualId != v.Expected {
			t.Fatalf("expected %q, got %q", v.Expected, actualId)
		}
	}
}

func TestRoleAssignmentID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *RoleAssignmentId
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
			// just subscription
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing ResourceGroup value
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing Management Group value
			Input: "/providers/Microsoft.Management/managementGroups/",
			Error: true,
		},

		{
			// missing Role Assignment value at Subscription Scope
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Authorization/roleAssignments/",
			Error: true,
		},

		{
			// missing Role Assignment value at Management Group scope
			Input: "/providers/Microsoft.Management/managementGroups/managementGroup1/providers/Microsoft.Authorization/roleAssignments/",
			Error: true,
		},

		{
			// valid at subscription scope
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Expected: &RoleAssignmentId{
				SubscriptionID:  "12345678-1234-9876-4563-123456789012",
				ResourceGroup:   "",
				ManagementGroup: "",
				Name:            "23456781-2349-8764-5631-234567890121",
			},
		},

		{
			// valid at resource group scope
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Expected: &RoleAssignmentId{
				SubscriptionID:  "12345678-1234-9876-4563-123456789012",
				ResourceGroup:   "group1",
				ManagementGroup: "",
				Name:            "23456781-2349-8764-5631-234567890121",
			},
		},

		{
			// valid at management group scope
			Input: "/providers/Microsoft.Management/managementGroups/managementGroup1/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Expected: &RoleAssignmentId{
				SubscriptionID:  "",
				ResourceGroup:   "",
				ManagementGroup: "managementGroup1",
				Name:            "23456781-2349-8764-5631-234567890121",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := RoleAssignmentID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("expected a value but got an error: %+v", err)
		}

		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Role Assignment Name", v.Expected.Name, actual.Name)
		}

		if actual.SubscriptionID != v.Expected.SubscriptionID {
			t.Fatalf("Expected %q but got %q for Role Assignment Name", v.Expected.SubscriptionID, actual.SubscriptionID)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Role Assignment Name", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.ManagementGroup != v.Expected.ManagementGroup {
			t.Fatalf("Expected %q but got %q for Role Assignment Name", v.Expected.ManagementGroup, actual.ManagementGroup)
		}
	}
}
