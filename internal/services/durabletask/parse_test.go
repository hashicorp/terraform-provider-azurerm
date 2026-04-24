// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/durabletask"
)

func TestParseRetentionPolicyID(t *testing.T) {
	testCases := []struct {
		Input    string
		Valid    bool
		Expected *durabletask.RetentionPolicyID
	}{
		{
			Input: "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg1/providers/Microsoft.DurableTask/schedulers/scheduler1/retentionPolicies/default",
			Valid: true,
			Expected: &durabletask.RetentionPolicyID{
				SubscriptionId:    "12345678-1234-1234-1234-123456789012",
				ResourceGroupName: "rg1",
				SchedulerName:     "scheduler1",
			},
		},
		{
			Input: "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test-rg/providers/Microsoft.DurableTask/schedulers/test-scheduler/retentionPolicies/default",
			Valid: true,
			Expected: &durabletask.RetentionPolicyID{
				SubscriptionId:    "12345678-1234-1234-1234-123456789012",
				ResourceGroupName: "test-rg",
				SchedulerName:     "test-scheduler",
			},
		},
		{
			Input: "/subscriptions/12345/resourceGroups/rg1/providers/Microsoft.Network/policies/policy1",
			Valid: false,
		},
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "this-is-not-a-valid-resource-id",
			Valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Input, func(t *testing.T) {
			result, err := durabletask.ParseRetentionPolicyID(tc.Input)
			valid := err == nil
			if tc.Valid != valid {
				t.Fatalf("expected valid=%t for %q, got valid=%t: %+v", tc.Valid, tc.Input, valid, err)
			}
			if tc.Valid {
				if result == nil {
					t.Fatalf("expected non-nil result for %q", tc.Input)
				}
				if result.SubscriptionId != tc.Expected.SubscriptionId {
					t.Fatalf("expected SubscriptionId %q, got %q", tc.Expected.SubscriptionId, result.SubscriptionId)
				}
				if result.ResourceGroupName != tc.Expected.ResourceGroupName {
					t.Fatalf("expected ResourceGroupName %q, got %q", tc.Expected.ResourceGroupName, result.ResourceGroupName)
				}
				if result.SchedulerName != tc.Expected.SchedulerName {
					t.Fatalf("expected SchedulerName %q, got %q", tc.Expected.SchedulerName, result.SchedulerName)
				}
			}
		})
	}
}

func TestRetentionPolicyID_ID(t *testing.T) {
	id := durabletask.RetentionPolicyID{
		SubscriptionId:    "12345678-1234-1234-1234-123456789012",
		ResourceGroupName: "test-rg",
		SchedulerName:     "test-scheduler",
	}

	expected := "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test-rg/providers/Microsoft.DurableTask/schedulers/test-scheduler/retentionPolicies/default"
	actual := id.ID()

	if expected != actual {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func TestValidateRetentionPolicyID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{Input: "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg1/providers/Microsoft.DurableTask/schedulers/scheduler1/retentionPolicies/default", Valid: true},
		{Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DurableTask/schedulers/test-scheduler/retentionPolicies/default", Valid: true},
		{Input: "", Valid: false},
		{Input: "not-a-valid-id", Valid: false},
	}

	for _, tc := range cases {
		_, errors := durabletask.ValidateRetentionPolicyID(tc.Input, "id")
		valid := len(errors) == 0
		if tc.Valid != valid {
			t.Fatalf("expected valid=%t for %q, got valid=%t", tc.Valid, tc.Input, valid)
		}
	}
}

func TestNewRetentionPolicyID(t *testing.T) {
	id := durabletask.NewRetentionPolicyID("sub-123", "rg-test", "scheduler-1")

	if id.SubscriptionId != "sub-123" {
		t.Fatalf("expected SubscriptionId %q, got %q", "sub-123", id.SubscriptionId)
	}
	if id.ResourceGroupName != "rg-test" {
		t.Fatalf("expected ResourceGroupName %q, got %q", "rg-test", id.ResourceGroupName)
	}
	if id.SchedulerName != "scheduler-1" {
		t.Fatalf("expected SchedulerName %q, got %q", "scheduler-1", id.SchedulerName)
	}
}
