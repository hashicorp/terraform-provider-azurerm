package parse

import (
	"testing"
)

func TestDataShareAccountID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *AccountId
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
			Name:     "Missing Account Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/",
			Expected: nil,
		},
		{
			Name:  "Datashare account ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/account1",
			Expected: &AccountId{
				Name:          "account1",
				ResourceGroup: "resGroup1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataShare/Accounts/account1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := AccountID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}

func TestDataShareID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *ShareId
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
			Name:     "Missing Account Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/",
			Expected: nil,
		},
		{
			Name:     "Missing Share",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/account1/",
			Expected: nil,
		},
		{
			Name:     "Missing Share Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/account1/shares/",
			Expected: nil,
		},
		{
			Name:  "Data Share ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/account1/shares/share1",
			Expected: &ShareId{
				Name:          "share1",
				AccountName:   "account1",
				ResourceGroup: "resGroup1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/account1/Shares/share1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := ShareID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.AccountName != v.Expected.AccountName {
			t.Fatalf("Expected %q but got %q for account name", v.Expected.AccountName, actual.AccountName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}

func TestDataShareDataSetID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *DataSetId
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
			Name:     "Missing Account Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/",
			Expected: nil,
		},
		{
			Name:     "Missing Share",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/account1/",
			Expected: nil,
		},
		{
			Name:     "Missing Share Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/account1/shares/",
			Expected: nil,
		},
		{
			Name:     "Missing DataSet",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/account1/shares/share1",
			Expected: nil,
		},
		{
			Name:     "Missing DataSet Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/account1/shares/share1/dataSets",
			Expected: nil,
		},
		{
			Name:  "DataSet ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/account1/shares/share1/dataSets/dataSet1",
			Expected: &DataSetId{
				Name:          "dataSet1",
				AccountName:   "account1",
				ResourceGroup: "resGroup1",
				ShareName:     "share1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/account1/shares/share1/DataSets/dataSet1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := DataSetID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ShareName != v.Expected.ShareName {
			t.Fatalf("Expected %q but got %q for account name", v.Expected.ShareName, actual.ShareName)
		}

		if actual.AccountName != v.Expected.AccountName {
			t.Fatalf("Expected %q but got %q for account name", v.Expected.AccountName, actual.AccountName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
