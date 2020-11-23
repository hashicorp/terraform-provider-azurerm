package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = DataShareDataSetId{}

func TestDataShareDataSetIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	actual := NewDataShareDataSetId(subscriptionId, "group1", "account1", "share1", "set1").ID("")
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.DataShare/accounts/account1/shares/share1/dataSets/set1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestDataShareDataSetID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *DataShareDataSetId
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
			Expected: &DataShareDataSetId{
				SubscriptionId: "00000000-0000-0000-0000-000000000000",
				Name:           "dataSet1",
				AccountName:    "account1",
				ResourceGroup:  "resGroup1",
				ShareName:      "share1",
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

		actual, err := DataShareDataSetID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for Subscription Id", v.Expected.SubscriptionId, actual.SubscriptionId)
		}

		if actual.ShareName != v.Expected.ShareName {
			t.Fatalf("Expected %q but got %q for Share Name", v.Expected.ShareName, actual.ShareName)
		}

		if actual.AccountName != v.Expected.AccountName {
			t.Fatalf("Expected %q but got %q for Account Name", v.Expected.AccountName, actual.AccountName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
