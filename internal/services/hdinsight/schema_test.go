// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package hdinsight

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

func TestHDInsightClusterVersionDiffSuppress(t *testing.T) {
	tests := []struct {
		name          string
		userInput     string
		azureResponse string
		suppressed    bool
	}{
		{
			name:          "empty name",
			userInput:     "",
			azureResponse: "",
			suppressed:    false,
		},
		{
			name:          "missing user input",
			userInput:     "",
			azureResponse: "1.2.3.4",
			suppressed:    false,
		},
		{
			name:          "missing api response",
			userInput:     "1.2",
			azureResponse: "",
			suppressed:    false,
		},
		{
			name:          "major minor user input",
			userInput:     "3.6",
			azureResponse: "3.6.1000.67",
			suppressed:    true,
		},
		{
			name:          "full version user input",
			userInput:     "3.6.1000.67",
			azureResponse: "3.6.1000.67",
			suppressed:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wasSuppressed := hdinsightClusterVersionDiffSuppressFunc("", tt.userInput, tt.azureResponse, nil)
			if tt.suppressed != wasSuppressed {
				t.Errorf("Expected %q to be %t but got %t", tt.name, tt.suppressed, wasSuppressed)
			}
		})
	}
}

func TestAddHDInsightUserAssignedIdentity(t *testing.T) {
	tests := []struct {
		name             string
		input            *identity.SystemAndUserAssignedMap
		expectedType     identity.Type
		expectedIdentity int
	}{
		{
			name:             "empty identity",
			expectedType:     identity.TypeUserAssigned,
			expectedIdentity: 1,
		},
		{
			name: "user assigned identity",
			input: &identity.SystemAndUserAssignedMap{
				Type: identity.TypeUserAssigned,
				IdentityIds: map[string]identity.UserAssignedIdentityDetails{
					"existing": {},
				},
			},
			expectedType:     identity.TypeUserAssigned,
			expectedIdentity: 2,
		},
		{
			name: "system assigned identity",
			input: &identity.SystemAndUserAssignedMap{
				Type: identity.TypeSystemAssigned,
			},
			expectedType:     identity.TypeSystemAssignedUserAssigned,
			expectedIdentity: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := addHDInsightUserAssignedIdentity(test.input, "added")
			if result.Type != test.expectedType {
				t.Fatalf("expected identity type %q, got %q", test.expectedType, result.Type)
			}
			if len(result.IdentityIds) != test.expectedIdentity {
				t.Fatalf("expected %d identities, got %d", test.expectedIdentity, len(result.IdentityIds))
			}
			if _, ok := result.IdentityIds["added"]; !ok {
				t.Fatal("expected added identity to be present")
			}
		})
	}
}
