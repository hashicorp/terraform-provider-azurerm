package clusters

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = PrivateCloudId{}

func TestPrivateCloudIDFormatter(t *testing.T) {
	actual := NewPrivateCloudID("{subscriptionId}", "{resourceGroupName}", "{privateCloudName}").ID()
	expected := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/{privateCloudName}"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestParsePrivateCloudID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *PrivateCloudId
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
			Input: "/subscriptions/{subscriptionId}/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/{subscriptionId}/resourceGroups/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/{privateCloudName}",
			Expected: &PrivateCloudId{
				SubscriptionId: "{subscriptionId}",
				ResourceGroup:  "{resourceGroupName}",
				Name:           "{privateCloudName}",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/{SUBSCRIPTIONID}/RESOURCEGROUPS/{RESOURCEGROUPNAME}/PROVIDERS/MICROSOFT.AVS/PRIVATECLOUDS/{PRIVATECLOUDNAME}",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParsePrivateCloudID(v.Input)
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
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}

func TestParsePrivateCloudIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *PrivateCloudId
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
			Input: "/subscriptions/{subscriptionId}/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/{subscriptionId}/resourceGroups/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/{privateCloudName}",
			Expected: &PrivateCloudId{
				SubscriptionId: "{subscriptionId}",
				ResourceGroup:  "{resourceGroupName}",
				Name:           "{privateCloudName}",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateclouds/{privateCloudName}",
			Expected: &PrivateCloudId{
				SubscriptionId: "{subscriptionId}",
				ResourceGroup:  "{resourceGroupName}",
				Name:           "{privateCloudName}",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/PRIVATECLOUDS/{privateCloudName}",
			Expected: &PrivateCloudId{
				SubscriptionId: "{subscriptionId}",
				ResourceGroup:  "{resourceGroupName}",
				Name:           "{privateCloudName}",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/PrIvAtEcLoUdS/{privateCloudName}",
			Expected: &PrivateCloudId{
				SubscriptionId: "{subscriptionId}",
				ResourceGroup:  "{resourceGroupName}",
				Name:           "{privateCloudName}",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParsePrivateCloudIDInsensitively(v.Input)
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
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
