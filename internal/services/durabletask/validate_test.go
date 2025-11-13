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
		"a",                           // minimum length
		strings.Repeat("a", 64),       // maximum length
		"scheduler-1",
		"my-scheduler",
		"test123scheduler",
	}

	invalidNames := []string{
		"",                            // empty
		"-invalid",                    // starts with hyphen
		"invalid-",                    // ends with hyphen
		"Invalid-Name",                // uppercase not allowed
		"invalid_underscore",          // underscores not allowed
		"invalid name",                // spaces not allowed
		"invalid.name",                // dots not allowed
		strings.Repeat("a", 65),       // too long
		"123invalid",                  // starts with number
		"invalid@name",                // special characters
		"invalid#name",
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
		"a",                           // minimum length
		strings.Repeat("a", 64),       // maximum length
		"taskhub-1",
		"test123hub",
	}

	invalidNames := []string{
		"",                            // empty
		"-invalid",                    // starts with hyphen
		"invalid-",                    // ends with hyphen
		"Invalid-Hub",                 // uppercase not allowed
		"invalid_underscore",          // underscores not allowed
		"invalid hub",                 // spaces not allowed
		"invalid.hub",                 // dots not allowed
		strings.Repeat("a", 65),       // too long
		"123invalid",                  // starts with number
		"invalid@hub",                 // special characters
		"hub!name",
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

func TestValidateRetentionPolicyName(t *testing.T) {
	validNames := []string{
		"validpolicy",
		"valid-policy",
		"policy123",
		"retention-policy-1",
		"a",                           // minimum length
		strings.Repeat("a", 64),       // maximum length
		"my-retention-policy",
		"test123policy",
	}

	invalidNames := []string{
		"",                            // empty
		"-invalid",                    // starts with hyphen
		"invalid-",                    // ends with hyphen
		"Invalid-Policy",              // uppercase not allowed
		"invalid_underscore",          // underscores not allowed
		"invalid policy",              // spaces not allowed
		"invalid.policy",              // dots not allowed
		strings.Repeat("a", 65),       // too long
		"123invalid",                  // starts with number
		"invalid@policy",              // special characters
		"policy#name",
	}

	for _, name := range validNames {
		t.Run("valid_"+name, func(t *testing.T) {
			warnings, errors := durabletask.ValidateRetentionPolicyName(name, "retention_policy_name")
			require.Empty(t, warnings, "expected no warnings for valid name: %s", name)
			require.Empty(t, errors, "expected no errors for valid name: %s", name)
		})
	}

	for _, name := range invalidNames {
		t.Run("invalid_"+name, func(t *testing.T) {
			_, errors := durabletask.ValidateRetentionPolicyName(name, "retention_policy_name")
			require.NotEmpty(t, errors, "expected errors for invalid name: %s", name)
		})
	}
}

func TestValidateSchedulerID(t *testing.T) {
	validIDs := []string{
		"/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg1/providers/Microsoft.DurableTask/schedulers/scheduler1",
		"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DurableTask/schedulers/test-scheduler",
	}

	invalidIDs := []string{
		"",                                                                                 // empty
		"not-a-valid-id",                                                                  // malformed
		"/subscriptions/invalid/resourceGroups/rg1/providers/Microsoft.DurableTask/schedulers/scheduler1", // invalid subscription
		"/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups//providers/Microsoft.DurableTask/schedulers/scheduler1", // empty resource group
		"/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/account1", // wrong provider
	}

	for _, id := range validIDs {
		t.Run("valid", func(t *testing.T) {
			warnings, errors := durabletask.ValidateSchedulerID(id, "id")
			require.Empty(t, warnings)
			require.Empty(t, errors)
		})
	}

	for _, id := range invalidIDs {
		t.Run("invalid", func(t *testing.T) {
			_, errors := durabletask.ValidateSchedulerID(id, "id")
			require.NotEmpty(t, errors)
		})
	}
}

func TestValidateTaskHubID(t *testing.T) {
	validIDs := []string{
		"/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg1/providers/Microsoft.DurableTask/taskHubs/hub1",
		"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DurableTask/taskHubs/test-hub",
	}

	invalidIDs := []string{
		"",                                                                                 // empty
		"not-a-valid-id",                                                                  // malformed
		"/subscriptions/invalid/resourceGroups/rg1/providers/Microsoft.DurableTask/taskHubs/hub1", // invalid subscription
		"/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups//providers/Microsoft.DurableTask/taskHubs/hub1", // empty resource group
		"/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg1/providers/Microsoft.Compute/virtualMachines/vm1", // wrong provider
	}

	for _, id := range validIDs {
		t.Run("valid", func(t *testing.T) {
			warnings, errors := durabletask.ValidateTaskHubID(id, "id")
			require.Empty(t, warnings)
			require.Empty(t, errors)
		})
	}

	for _, id := range invalidIDs {
		t.Run("invalid", func(t *testing.T) {
			_, errors := durabletask.ValidateTaskHubID(id, "id")
			require.NotEmpty(t, errors)
		})
	}
}

func TestValidateRetentionPolicyID(t *testing.T) {
	validIDs := []string{
		"/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg1/providers/Microsoft.DurableTask/retentionPolicies/policy1",
		"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DurableTask/retentionPolicies/retention-policy",
	}

	invalidIDs := []string{
		"",                                                                                 // empty
		"not-a-valid-id",                                                                  // malformed
		"/subscriptions/invalid/resourceGroups/rg1/providers/Microsoft.DurableTask/retentionPolicies/policy1", // invalid subscription
		"/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups//providers/Microsoft.DurableTask/retentionPolicies/policy1", // empty resource group
		"/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg1/providers/Microsoft.Network/policies/policy1", // wrong provider
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

func TestValidationErrorMessages(t *testing.T) {
	testCases := []struct {
		name          string
		validator     func(interface{}, string) ([]string, []error)
		input         string
		expectedError string
	}{
		{
			name:          "scheduler name empty",
			validator:     durabletask.ValidateSchedulerName,
			input:         "",
			expectedError: "cannot be empty",
		},
		{
			name:          "scheduler name too long",
			validator:     durabletask.ValidateSchedulerName,
			input:         strings.Repeat("a", 65),
			expectedError: "cannot be longer than 64 characters",
		},
		{
			name:          "scheduler name invalid characters",
			validator:     durabletask.ValidateSchedulerName,
			input:         "invalid_name",
			expectedError: "can only contain",
		},
		{
			name:          "task hub name starts with hyphen",
			validator:     durabletask.ValidateTaskHubName,
			input:         "-invalid",
			expectedError: "must start with a letter",
		},
		{
			name:          "retention policy name ends with hyphen",
			validator:     durabletask.ValidateRetentionPolicyName,
			input:         "invalid-",
			expectedError: "must end with a letter or number",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := tc.validator(tc.input, "test_field")
			require.NotEmpty(t, errors)
			
			// Check that at least one error message contains the expected text
			found := false
			for _, err := range errors {
				if strings.Contains(err.Error(), tc.expectedError) {
					found = true
					break
				}
			}
			require.True(t, found, "expected error message to contain '%s', but got: %v", tc.expectedError, errors)
		})
	}
}

func TestValidationFieldNameInErrors(t *testing.T) {
	testCases := []struct {
		name      string
		validator func(interface{}, string) ([]string, []error)
		input     string
		fieldName string
	}{
		{
			name:      "scheduler name validation includes field name",
			validator: durabletask.ValidateSchedulerName,
			input:     "",
			fieldName: "scheduler_name",
		},
		{
			name:      "task hub name validation includes field name",
			validator: durabletask.ValidateTaskHubName,
			input:     "",
			fieldName: "task_hub_name",
		},
		{
			name:      "retention policy name validation includes field name",
			validator: durabletask.ValidateRetentionPolicyName,
			input:     "",
			fieldName: "retention_policy_name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := tc.validator(tc.input, tc.fieldName)
			require.NotEmpty(t, errors)
			
			// Check that error message contains the field name in backticks
			errorStr := errors[0].Error()
			require.Contains(t, errorStr, "`"+tc.fieldName+"`", "error message should contain field name in backticks")
		})
	}
}
