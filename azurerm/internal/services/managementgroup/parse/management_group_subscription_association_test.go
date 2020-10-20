package parse

import (
	"testing"
)

func TestManagementGroupSubscriptionAssociationID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *ManagementGroupSubscriptionAssociationId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "One Segment",
			Input: "hello",
			Error: true,
		},
		{
			Name:  "Two Segments Invalid ID's",
			Input: "hello|world",
			Error: true,
		},
		{
			Name:  "Bad Management Group",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup|/subscriptions/00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Name:  "Subscription GUID, not scope",
			Input: "/providers/Microsoft.Management/managementGroups/myManagementGroup|00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Name:  "Bad GUID",
			Input: "/providers/Microsoft.Management/managementGroups/myManagementGroup|/subscriptions/mySubscriptionAlias",
			Error: true,
		},
		{
			Name:  "Nat Gateway / Public IP Association ID",
			Input: "/providers/Microsoft.Management/managementGroups/myManagementGroup|/subscriptions/00000000-0000-0000-0000-000000000000",
			Error: false,
			Expect: &ManagementGroupSubscriptionAssociationId{
				ManagementGroupID:   "/providers/Microsoft.Management/managementGroups/myManagementGroup",
				ManagementGroupName: "myManagementGroup",
				SubscriptionScopeID: "/subscriptions/00000000-0000-0000-0000-000000000000",
				SubscriptionID:      "00000000-0000-0000-0000-000000000000",
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

		if actual.ManagementGroupID != v.Expect.ManagementGroupID {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.ManagementGroupID, actual.ManagementGroupID)
		}

		if actual.ManagementGroupName != v.Expect.ManagementGroupName {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.ManagementGroupName, actual.ManagementGroupName)
		}

		if actual.SubscriptionScopeID != v.Expect.SubscriptionScopeID {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.SubscriptionScopeID, actual.SubscriptionScopeID)
		}

		if actual.SubscriptionID != v.Expect.SubscriptionID {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.SubscriptionID, actual.SubscriptionID)
		}
	}
}
