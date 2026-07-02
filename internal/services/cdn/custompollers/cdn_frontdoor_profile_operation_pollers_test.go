// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"errors"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
)

func TestFrontDoorProfileOperationInProgress(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error is not in progress",
			err:      nil,
			expected: false,
		},
		{
			name:     "matches another operation in progress",
			err:      errors.New("Conflict: The requested operation cannot be executed on the entity as another operation is in progress."),
			expected: true,
		},
		{
			name:     "matches generic operation in progress",
			err:      errors.New("operation is in progress"),
			expected: true,
		},
		{
			name:     "other error is not in progress",
			err:      errors.New("something else failed"),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := frontDoorProfileOperationInProgress(test.err)
			if actual != test.expected {
				t.Fatalf("expected %t but got %t", test.expected, actual)
			}
		})
	}
}

func TestFrontDoorProfileStatusesSettled(t *testing.T) {
	tests := []struct {
		name        string
		input       *profiles.ProfileProperties
		expected    bool
		expectError bool
	}{
		{
			name:     "missing properties are not settled",
			input:    nil,
			expected: false,
		},
		{
			name:     "missing provisioning state is not settled",
			input:    &profiles.ProfileProperties{},
			expected: false,
		},
		{
			name: "creating is not settled",
			input: &profiles.ProfileProperties{
				ProvisioningState: pointerToProfileProvisioningState(profiles.ProfileProvisioningStateCreating),
			},
			expected: false,
		},
		{
			name: "succeeded is settled",
			input: &profiles.ProfileProperties{
				ProvisioningState: pointerToProfileProvisioningState(profiles.ProfileProvisioningStateSucceeded),
			},
			expected: true,
		},
		{
			name: "failed errors",
			input: &profiles.ProfileProperties{
				ProvisioningState: pointerToProfileProvisioningState(profiles.ProfileProvisioningStateFailed),
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := frontDoorProfileStatusesSettled(test.input)
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

func pointerToProfileProvisioningState(input profiles.ProfileProvisioningState) *profiles.ProfileProvisioningState {
	return &input
}
