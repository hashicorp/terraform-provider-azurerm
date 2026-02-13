// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/durabletask"
	"github.com/stretchr/testify/require"
)

func TestParseRetentionPolicyID(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *durabletask.RetentionPolicyID
		wantErr  bool
	}{
		{
			name:  "valid retention policy ID",
			input: "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg1/providers/Microsoft.DurableTask/schedulers/scheduler1/retentionPolicies/default",
			expected: &durabletask.RetentionPolicyID{
				SubscriptionId:    "12345678-1234-1234-1234-123456789012",
				ResourceGroupName: "rg1",
				SchedulerName:     "scheduler1",
			},
			wantErr: false,
		},
		{
			name:  "valid retention policy ID with hyphenated names",
			input: "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test-rg/providers/Microsoft.DurableTask/schedulers/test-scheduler/retentionPolicies/default",
			expected: &durabletask.RetentionPolicyID{
				SubscriptionId:    "12345678-1234-1234-1234-123456789012",
				ResourceGroupName: "test-rg",
				SchedulerName:     "test-scheduler",
			},
			wantErr: false,
		},
		{
			name:    "invalid retention policy ID - wrong provider",
			input:   "/subscriptions/12345/resourceGroups/rg1/providers/Microsoft.Network/policies/policy1",
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
			require.Equal(t, tc.expected.SchedulerName, result.SchedulerName)
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
	result := id.ID()

	require.Equal(t, expected, result)
}

func TestNewRetentionPolicyID(t *testing.T) {
	id := durabletask.NewRetentionPolicyID("sub-123", "rg-test", "scheduler-1")

	require.Equal(t, "sub-123", id.SubscriptionId)
	require.Equal(t, "rg-test", id.ResourceGroupName)
	require.Equal(t, "scheduler-1", id.SchedulerName)
}
