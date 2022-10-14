package parse

import (
	"testing"
)

func TestGenericResourceID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *GenericResourceId
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
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Features/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Features/features/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Features/features/Feature1",
			Expected: &GenericResourceId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceProvider: "Microsoft.Features",
				ResourceType:     "features",
			},
		},

		{
			// longer resource type
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Features/parentFeature/Feature1/feature/feature1",
			Expected: &GenericResourceId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceProvider: "Microsoft.Features",
				ResourceType:     "parentFeature/feature",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.FEATURES/FEATURES/FEATURE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := GenericResourceID(v.Input)
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
		if actual.ResourceProvider != v.Expected.ResourceProvider {
			t.Fatalf("Expected %q but got %q for ResourceProvider", v.Expected.ResourceProvider, actual.ResourceProvider)
		}
		if actual.ResourceType != v.Expected.ResourceType {
			t.Fatalf("Expected %q but got %q for ResourceType", v.Expected.ResourceType, actual.ResourceType)
		}
	}
}
