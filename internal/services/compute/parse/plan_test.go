// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = PlanId{}

func TestPlanIDFormatter(t *testing.T) {
	actual := NewPlanID("12345678-1234-9876-4563-123456789012", "agreement1", "offer1", "hourly").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.MarketplaceOrdering/agreements/agreement1/offers/offer1/plans/hourly"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestPlanID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *PlanId
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
			// missing AgreementName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.MarketplaceOrdering/",
			Error: true,
		},

		{
			// missing value for AgreementName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.MarketplaceOrdering/agreements/",
			Error: true,
		},

		{
			// missing OfferName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.MarketplaceOrdering/agreements/agreement1/",
			Error: true,
		},

		{
			// missing value for OfferName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.MarketplaceOrdering/agreements/agreement1/offers/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.MarketplaceOrdering/agreements/agreement1/offers/offer1/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.MarketplaceOrdering/agreements/agreement1/offers/offer1/plans/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.MarketplaceOrdering/agreements/agreement1/offers/offer1/plans/hourly",
			Expected: &PlanId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				AgreementName:  "agreement1",
				OfferName:      "offer1",
				Name:           "hourly",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.MARKETPLACEORDERING/AGREEMENTS/AGREEMENT1/OFFERS/OFFER1/PLANS/HOURLY",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := PlanID(v.Input)
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
		if actual.AgreementName != v.Expected.AgreementName {
			t.Fatalf("Expected %q but got %q for AgreementName", v.Expected.AgreementName, actual.AgreementName)
		}
		if actual.OfferName != v.Expected.OfferName {
			t.Fatalf("Expected %q but got %q for OfferName", v.Expected.OfferName, actual.OfferName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
