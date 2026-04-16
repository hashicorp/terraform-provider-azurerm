// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/durabletask"
)

func TestValidateSchedulerName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{Input: "validname", Valid: true},
		{Input: "valid-name", Valid: true},
		{Input: "valid123", Valid: true},
		{Input: "name-with-numbers-123", Valid: true},
		{Input: "abc", Valid: true},
		{Input: strings.Repeat("a", 63), Valid: true},
		{Input: "scheduler-1", Valid: true},
		{Input: "my-scheduler", Valid: true},
		{Input: "test123scheduler", Valid: true},
		{Input: "", Valid: false},
		{Input: "ab", Valid: false},
		{Input: "-invalid", Valid: false},
		{Input: "invalid-", Valid: false},
		{Input: strings.Repeat("a", 64), Valid: false},
	}

	for _, tc := range cases {
		_, errors := durabletask.ValidateSchedulerName(tc.Input, "scheduler_name")
		valid := len(errors) == 0
		if tc.Valid != valid {
			t.Fatalf("expected valid=%t for %q, got valid=%t", tc.Valid, tc.Input, valid)
		}
	}
}

func TestValidateTaskHubName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{Input: "validhub", Valid: true},
		{Input: "valid-hub", Valid: true},
		{Input: "hub123", Valid: true},
		{Input: "my-task-hub", Valid: true},
		{Input: "abc", Valid: true},
		{Input: strings.Repeat("a", 63), Valid: true},
		{Input: "taskhub-1", Valid: true},
		{Input: "test123hub", Valid: true},
		{Input: "", Valid: false},
		{Input: "ab", Valid: false},
		{Input: "-invalid", Valid: false},
		{Input: "invalid-", Valid: false},
		{Input: strings.Repeat("a", 64), Valid: false},
	}

	for _, tc := range cases {
		_, errors := durabletask.ValidateTaskHubName(tc.Input, "task_hub_name")
		valid := len(errors) == 0
		if tc.Valid != valid {
			t.Fatalf("expected valid=%t for %q, got valid=%t", tc.Valid, tc.Input, valid)
		}
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
