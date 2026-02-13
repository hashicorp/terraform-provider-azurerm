// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/durabletask"
	"github.com/stretchr/testify/require"
)

func TestValidateSchedulerName(t *testing.T) {
	validNames := []string{
		"validname",
		"valid-name",
		"valid123",
		"name-with-numbers-123",
		"abc",                   // minimum length
		strings.Repeat("a", 63), // maximum length
		"scheduler-1",
		"my-scheduler",
		"test123scheduler",
	}

	invalidNames := []string{
		"",                      // empty
		"ab",                    // too short
		"-invalid",              // starts with hyphen
		"invalid-",              // ends with hyphen
		strings.Repeat("a", 64), // too long
	}

	for _, name := range validNames {
		t.Run("valid_"+name, func(t *testing.T) {
			warnings, errors := durabletask.ValidateSchedulerName(name, "scheduler_name")
			require.Empty(t, warnings, "expected no warnings for valid name: %s", name)
			require.Empty(t, errors, "expected no errors for valid name: %s", name)
		})
	}

	for _, name := range invalidNames {
		t.Run("invalid_"+name, func(t *testing.T) {
			_, errors := durabletask.ValidateSchedulerName(name, "scheduler_name")
			require.NotEmpty(t, errors, "expected errors for invalid name: %s", name)
		})
	}
}

func TestValidateTaskHubName(t *testing.T) {
	validNames := []string{
		"validhub",
		"valid-hub",
		"hub123",
		"my-task-hub",
		"abc",                   // minimum length
		strings.Repeat("a", 63), // maximum length
		"taskhub-1",
		"test123hub",
	}

	invalidNames := []string{
		"",                      // empty
		"ab",                    // too short
		"-invalid",              // starts with hyphen
		"invalid-",              // ends with hyphen
		strings.Repeat("a", 64), // too long
	}

	for _, name := range validNames {
		t.Run("valid_"+name, func(t *testing.T) {
			warnings, errors := durabletask.ValidateTaskHubName(name, "task_hub_name")
			require.Empty(t, warnings, "expected no warnings for valid name: %s", name)
			require.Empty(t, errors, "expected no errors for valid name: %s", name)
		})
	}

	for _, name := range invalidNames {
		t.Run("invalid_"+name, func(t *testing.T) {
			_, errors := durabletask.ValidateTaskHubName(name, "task_hub_name")
			require.NotEmpty(t, errors, "expected errors for invalid name: %s", name)
		})
	}
}

func TestValidateRetentionPolicyID(t *testing.T) {
	validIDs := []string{
		"/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg1/providers/Microsoft.DurableTask/schedulers/scheduler1/retentionPolicies/default",
		"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DurableTask/schedulers/test-scheduler/retentionPolicies/default",
	}

	invalidIDs := []string{
		"",               // empty
		"not-a-valid-id", // malformed
	}

	for _, id := range validIDs {
		t.Run("valid", func(t *testing.T) {
			warnings, errors := durabletask.ValidateRetentionPolicyID(id, "id")
			require.Empty(t, warnings)
			require.Empty(t, errors)
		})
	}

	for _, id := range invalidIDs {
		t.Run("invalid", func(t *testing.T) {
			_, errors := durabletask.ValidateRetentionPolicyID(id, "id")
			require.NotEmpty(t, errors)
		})
	}
}
