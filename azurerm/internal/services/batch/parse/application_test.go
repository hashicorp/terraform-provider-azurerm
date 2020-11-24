package parse

import (
	"testing"
)

func TestBatchApplicationId(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *ApplicationId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Expected: nil,
		},
		{
			Name:     "Missing Accounts Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Batch/batchAccounts/",
			Expected: nil,
		},
		{
			Name:     "Batch Account ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Batch/batchAccounts/acctName1",
			Expected: nil,
		},
		{
			Name:     "Missing Applications Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Batch/batchAccounts/acctName1/applications/",
			Expected: nil,
		},
		{
			Name:  "Batch Application ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Batch/batchAccounts/acctName1/applications/Application1",
			Expected: &ApplicationId{
				ResourceGroup:    "resGroup1",
				BatchAccountName: "acctName1",
				ApplicationName:  "Application1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Batch/batchAccounts/acctName1/Applications/Application1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ApplicationID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ApplicationName != v.Expected.ApplicationName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.ApplicationName, actual.ApplicationName)
		}
		if actual.BatchAccountName != v.Expected.BatchAccountName {
			t.Fatalf("Expected %q but got %q for BatchAccountName", v.Expected.BatchAccountName, actual.BatchAccountName)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
