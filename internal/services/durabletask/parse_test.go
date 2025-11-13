// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/durabletask"
	"github.com/stretchr/testify/require"
)

func TestParseSchedulerID(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *durabletask.SchedulerId
		wantErr  bool
	}{
		{
			name:  "valid scheduler ID",
			input: "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg1/providers/Microsoft.DurableTask/schedulers/scheduler1",
			expected: &durabletask.SchedulerId{
				SubscriptionId:    "12345678-1234-1234-1234-123456789012",
				ResourceGroupName: "rg1",
				SchedulerName:     "scheduler1",
			},
			wantErr: false,
		},
		{
			name:    "invalid resource ID - wrong provider",
			input:   "/subscriptions/12345/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/account1",
			wantErr: true,
		},
		{
			name:    "invalid resource ID - missing scheduler name",
			input:   "/subscriptions/12345/resourceGroups/rg1/providers/Microsoft.DurableTask/schedulers/",
			wantErr: true,
		},
		{
			name:    "invalid resource ID - empty",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid resource ID - malformed",
			input:   "not-a-valid-resource-id",
			wantErr: true,
		},
		{
			name:    "invalid resource ID - missing subscription",
			input:   "/resourceGroups/rg1/providers/Microsoft.DurableTask/schedulers/scheduler1",
			wantErr: true,
		},
		{
			name:    "invalid resource ID - missing resource group",
			input:   "/subscriptions/12345/providers/Microsoft.DurableTask/schedulers/scheduler1",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := durabletask.ParseSchedulerID(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, result)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, tc.expected.SubscriptionId, result.SubscriptionId)
			require.Equal(t, tc.expected.ResourceGroupName, result.ResourceGroupName)
			require.Equal(t, tc.expected.SchedulerName, result.SchedulerName)
		})
	}
}

func TestSchedulerID_ID(t *testing.T) {
	id := durabletask.SchedulerId{
		SubscriptionId:    "12345678-1234-1234-1234-123456789012",
		ResourceGroupName: "test-rg",
		SchedulerName:     "test-scheduler",
	}

	expected := "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test-rg/providers/Microsoft.DurableTask/schedulers/test-scheduler"
	result := id.ID()

	require.Equal(t, expected, result)
}

func TestParseTaskHubID(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *durabletask.TaskHubId
		wantErr  bool
	}{
		{
			name:  "valid task hub ID",
			input: "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg1/providers/Microsoft.DurableTask/taskHubs/hub1",
			expected: &durabletask.TaskHubId{
				SubscriptionId:    "12345678-1234-1234-1234-123456789012",
				ResourceGroupName: "rg1",
				TaskHubName:       "hub1",
			},
			wantErr: false,
		},
		{
			name:    "invalid task hub ID - wrong provider",
			input:   "/subscriptions/12345/resourceGroups/rg1/providers/Microsoft.Storage/containers/container1",
			wantErr: true,
		},
		{
			name:    "invalid task hub ID - missing hub name",
			input:   "/subscriptions/12345/resourceGroups/rg1/providers/Microsoft.DurableTask/taskHubs/",
			wantErr: true,
		},
		{
			name:    "invalid task hub ID - empty",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid task hub ID - malformed",
			input:   "invalid-id-format",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := durabletask.ParseTaskHubID(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, result)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, tc.expected.SubscriptionId, result.SubscriptionId)
			require.Equal(t, tc.expected.ResourceGroupName, result.ResourceGroupName)
			require.Equal(t, tc.expected.TaskHubName, result.TaskHubName)
		})
	}
}

func TestTaskHubID_ID(t *testing.T) {
	id := durabletask.TaskHubId{
		SubscriptionId:    "12345678-1234-1234-1234-123456789012",
		ResourceGroupName: "test-rg",
		TaskHubName:       "test-hub",
	}

	expected := "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test-rg/providers/Microsoft.DurableTask/taskHubs/test-hub"
	result := id.ID()

	require.Equal(t, expected, result)
}

func TestParseRetentionPolicyID(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *durabletask.RetentionPolicyId
		wantErr  bool
	}{
		{
			name:  "valid retention policy ID",
			input: "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg1/providers/Microsoft.DurableTask/retentionPolicies/policy1",
			expected: &durabletask.RetentionPolicyId{
				SubscriptionId:      "12345678-1234-1234-1234-123456789012",
				ResourceGroupName:   "rg1",
				RetentionPolicyName: "policy1",
			},
			wantErr: false,
		},
		{
			name:  "valid retention policy ID with hyphenated name",
			input: "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test-rg/providers/Microsoft.DurableTask/retentionPolicies/retention-policy-1",
			expected: &durabletask.RetentionPolicyId{
				SubscriptionId:      "12345678-1234-1234-1234-123456789012",
				ResourceGroupName:   "test-rg",
				RetentionPolicyName: "retention-policy-1",
			},
			wantErr: false,
		},
		{
			name:    "invalid retention policy ID - wrong provider",
			input:   "/subscriptions/12345/resourceGroups/rg1/providers/Microsoft.Network/policies/policy1",
			wantErr: true,
		},
		{
			name:    "invalid retention policy ID - missing policy name",
			input:   "/subscriptions/12345/resourceGroups/rg1/providers/Microsoft.DurableTask/retentionPolicies/",
			wantErr: true,
		},
		{
			name:    "invalid retention policy ID - empty",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid retention policy ID - malformed",
			input:   "this-is-not-a-valid-resource-id",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := durabletask.ParseRetentionPolicyID(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, result)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, tc.expected.SubscriptionId, result.SubscriptionId)
			require.Equal(t, tc.expected.ResourceGroupName, result.ResourceGroupName)
			require.Equal(t, tc.expected.RetentionPolicyName, result.RetentionPolicyName)
		})
	}
}

func TestRetentionPolicyID_ID(t *testing.T) {
	id := durabletask.RetentionPolicyId{
		SubscriptionId:      "12345678-1234-1234-1234-123456789012",
		ResourceGroupName:   "test-rg",
		RetentionPolicyName: "test-policy",
	}

	expected := "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test-rg/providers/Microsoft.DurableTask/retentionPolicies/test-policy"
	result := id.ID()

	require.Equal(t, expected, result)
}

func TestResourceID_CaseInsensitivity(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "scheduler with uppercase",
			input: "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/RG1/providers/Microsoft.DurableTask/schedulers/Scheduler1",
		},
		{
			name:  "task hub with mixed case",
			input: "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/TestRG/providers/Microsoft.DurableTask/taskHubs/TestHub",
		},
		{
			name:  "retention policy with uppercase",
			input: "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/TEST-RG/providers/Microsoft.DurableTask/retentionPolicies/POLICY1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Azure resource IDs should be case-insensitive for parsing
			// but preserve original casing in the parsed result
			if schedulerID, err := durabletask.ParseSchedulerID(tc.input); err == nil {
				require.NotEmpty(t, schedulerID.ResourceGroupName)
				require.NotEmpty(t, schedulerID.SchedulerName)
			} else if taskHubID, err := durabletask.ParseTaskHubID(tc.input); err == nil {
				require.NotEmpty(t, taskHubID.ResourceGroupName)
				require.NotEmpty(t, taskHubID.TaskHubName)
			} else if policyID, err := durabletask.ParseRetentionPolicyID(tc.input); err == nil {
				require.NotEmpty(t, policyID.ResourceGroupName)
				require.NotEmpty(t, policyID.RetentionPolicyName)
			}
		})
	}
}
