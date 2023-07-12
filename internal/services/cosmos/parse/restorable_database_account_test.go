// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = RestorableDatabaseAccountId{}

func TestRestorableDatabaseAccountIDFormatter(t *testing.T) {
	actual := NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "location1", "account1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/locations/location1/restorableDatabaseAccounts/account1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestRestorableDatabaseAccountID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *RestorableDatabaseAccountId
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
			// missing LocationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/",
			Error: true,
		},

		{
			// missing value for LocationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/locations/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/locations/location1/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/locations/location1/restorableDatabaseAccounts/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/locations/location1/restorableDatabaseAccounts/account1",
			Expected: &RestorableDatabaseAccountId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				LocationName:   "location1",
				Name:           "account1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.DOCUMENTDB/LOCATIONS/LOCATION1/RESTORABLEDATABASEACCOUNTS/ACCOUNT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := RestorableDatabaseAccountID(v.Input)
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
		if actual.LocationName != v.Expected.LocationName {
			t.Fatalf("Expected %q but got %q for LocationName", v.Expected.LocationName, actual.LocationName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
