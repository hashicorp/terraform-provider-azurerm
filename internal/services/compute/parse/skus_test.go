// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SkusId{}

func TestSkusIDFormatter(t *testing.T) {
	actual := NewSkusID("707acc15-6870-4327-99cb-acf3b7fd1633").ID()
	expected := "/subscriptions/707acc15-6870-4327-99cb-acf3b7fd1633"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSkusID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SkusId
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
			// valid
			Input: "/subscriptions/707acc15-6870-4327-99cb-acf3b7fd1633",
			Expected: &SkusId{
				SubscriptionId: "707acc15-6870-4327-99cb-acf3b7fd1633",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/707ACC15-6870-4327-99CB-ACF3B7FD1633",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SkusID(v.Input)
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
	}
}
