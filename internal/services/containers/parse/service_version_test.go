// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ServiceVersionId{}

func TestServiceVersionIDFormatter(t *testing.T) {
	actual := NewServiceVersionID("12345678-1234-9876-4563-123456789012", "location1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService/locations/location1/orchestrators"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestServiceVersionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ServiceVersionId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService/",
			Error: true,
		},

		{
			// missing value for LocationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService/locations/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService/locations/location1",
			Expected: &ServiceVersionId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				LocationName:   "location1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.CONTAINERSERVICE/LOCATIONS/LOCATION1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ServiceVersionID(v.Input)
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
	}
}
