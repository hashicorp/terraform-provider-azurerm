// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import "testing"

func TestFeatureIDFormatter(t *testing.T) {
	actual := NewFeatureID("12345678-1234-9876-4563-123456789012", "Microsoft.Test", "Feature1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Features/providers/Microsoft.Test/features/Feature1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestFeatureID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FeatureId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Features/providers/Microsoft.SomeAzureService/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Features/providers/Microsoft.SomeAzureService/features/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Features/providers/Microsoft.SomeAzureService/features/Feature1",
			Expected: &FeatureId{
				SubscriptionId:    "12345678-1234-9876-4563-123456789012",
				ProviderNamespace: "Microsoft.SomeAzureService",
				Name:              "Feature1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.FEATURES/PROVIDERS/MICROSOFT.SOMEAZURESERVICE/FEATURES/FEATURE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := FeatureID(v.Input)
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
		if actual.ProviderNamespace != v.Expected.ProviderNamespace {
			t.Fatalf("Expected %q but got %q for ProviderNamespace", v.Expected.ProviderNamespace, actual.ProviderNamespace)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
