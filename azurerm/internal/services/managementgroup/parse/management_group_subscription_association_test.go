package parse

import "testing"

func TestManagementGroupSubscriptionAssociationID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *ManagementGroupSubscriptionAssociationId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "Missing Subscription",
			Input: "/managementGroup/MyManagementGroup",
			Error: true,
		},
		{
			Name:  "Missing Subscription Id",
			Input: "/managementGroup/MyManagementGroup/subscription/",
			Error: true,
		},
		{
			Name:  "Missing Management Group",
			Input: "/subscription/12345678-1234-1234-1234-123456789012",
			Error: true,
		},
		{
			Name:  "Missing Management Group Name",
			Input: "/managementGroup/subscription/12345678-1234-1234-1234-123456789012",
			Error: true,
		},
		{
			Name:  "Wrong Case",
			Input: "/MANAGEMENTGROUP/MyManagementGroup/SUBSCRIPTION/12345678-1234-1234-1234-123456789012",
			Error: true,
		},
		{
			Name:  "Valid",
			Input: "/managementGroup/MyManagementGroup/subscription/12345678-1234-1234-1234-123456789012",
			Expected: &ManagementGroupSubscriptionAssociationId{
				ManagementGroup: "MyManagementGroup",
				SubscriptionId:  "12345678-1234-1234-1234-123456789012",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ManagementGroupSubscriptionAssociationID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ManagementGroup != v.Expected.ManagementGroup {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.ManagementGroup, actual.ManagementGroup)
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
	}
}
