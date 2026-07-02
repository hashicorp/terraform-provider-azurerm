// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
)

func TestFrontDoorRouteStatusesSettled(t *testing.T) {
	tests := []struct {
		name        string
		input       *cdn.RouteProperties
		expected    bool
		expectError bool
	}{
		{
			name:     "missing route properties are not settled",
			input:    nil,
			expected: false,
		},
		{
			name: "provisioning in progress blocks",
			input: &cdn.RouteProperties{
				ProvisioningState: cdn.AfdProvisioningStateCreating,
				DeploymentStatus:  cdn.DeploymentStatusInProgress,
			},
			expected: false,
		},
		{
			name: "deployment status not started does not block once provisioning succeeds",
			input: &cdn.RouteProperties{
				ProvisioningState: cdn.AfdProvisioningStateSucceeded,
				DeploymentStatus:  cdn.DeploymentStatusNotStarted,
			},
			expected: true,
		},
		{
			name: "deployment failure errors",
			input: &cdn.RouteProperties{
				ProvisioningState: cdn.AfdProvisioningStateSucceeded,
				DeploymentStatus:  cdn.DeploymentStatusFailed,
			},
			expectError: true,
		},
		{
			name: "provisioning failure errors",
			input: &cdn.RouteProperties{
				ProvisioningState: cdn.AfdProvisioningStateFailed,
				DeploymentStatus:  cdn.DeploymentStatusSucceeded,
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := frontDoorRouteStatusesSettled(test.input)
			if test.expectError {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error but got %q", err)
			}

			if actual != test.expected {
				t.Fatalf("expected %t but got %t", test.expected, actual)
			}
		})
	}
}
