package parse

import (
	"testing"
)

func TestSecurityCenterSubscriptionPricingID(t *testing.T) {
	testData := []struct {
		ResourceType string
		Input        string
		Error        bool
		Expect       *SecurityCenterSubscriptionPricingId
	}{
		{
			ResourceType: "Empty",
			Input:        "",
			Error:        true,
		},
		{
			ResourceType: "No Pricings Segment",
			Input:        "/subscriptions/00000000-0000-0000-0000-000000000000",
			Error:        true,
		},
		{
			ResourceType: "No Pricings Value",
			Input:        "/subscriptions/00000000-0000-0000-0000-000000000000/pricings/",
			Error:        true,
		},
		{
			ResourceType: "Security Center Subscription Pricing ID",
			Input:        "/subscriptions/00000000-0000-0000-0000-000000000000/pricings/VirtualMachines",
			Expect: &SecurityCenterSubscriptionPricingId{
				ResourceType: "VirtualMachines",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.ResourceType)

		actual, err := SecurityCenterSubscriptionPricingID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ResourceType != v.Expect.ResourceType {
			t.Fatalf("Expected %q but got %q for ResourceType", v.Expect.ResourceType, actual.ResourceType)
		}
	}
}
