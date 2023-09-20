// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SubscriptionPolicyRemediationId{}

func TestSubscriptionPolicyRemediationIDFormatter(t *testing.T) {
	actual := NewSubscriptionPolicyRemediationID("12345678-1234-9876-4563-123456789012", "remediation1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.PolicyInsights/remediations/remediation1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSubscriptionPolicyRemediationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SubscriptionPolicyRemediationId
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
			// missing RemediationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.PolicyInsights/",
			Error: true,
		},

		{
			// missing value for RemediationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.PolicyInsights/remediations/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.PolicyInsights/remediations/remediation1",
			Expected: &SubscriptionPolicyRemediationId{
				SubscriptionId:  "12345678-1234-9876-4563-123456789012",
				RemediationName: "remediation1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.POLICYINSIGHTS/REMEDIATIONS/REMEDIATION1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SubscriptionPolicyRemediationID(v.Input)
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
		if actual.RemediationName != v.Expected.RemediationName {
			t.Fatalf("Expected %q but got %q for RemediationName", v.Expected.RemediationName, actual.RemediationName)
		}
	}
}
