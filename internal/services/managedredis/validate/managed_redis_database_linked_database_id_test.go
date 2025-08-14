// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestValidateLinkedDatabaseIncludesSelf(t *testing.T) {
	testCases := []struct {
		name                string
		linkedDatabaseIds   []string
		selfDbId            string
		expectError         bool
		expectedErrorString string
	}{
		{
			name:              "Empty linked database list should return nil",
			linkedDatabaseIds: []string{},
			selfDbId:          "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Cache/redisEnterprise/example-cluster/databases/database1",
			expectError:       false,
		},
		{
			name:              "Nil linked database list should return nil",
			linkedDatabaseIds: nil,
			selfDbId:          "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Cache/redisEnterprise/example-cluster/databases/database1",
			expectError:       false,
		},
		{
			name: "Self database included in linked list should return nil",
			linkedDatabaseIds: []string{
				"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Cache/redisEnterprise/example-cluster/databases/database1",
				"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Cache/redisEnterprise/example-cluster/databases/database2",
			},
			selfDbId:    "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Cache/redisEnterprise/example-cluster/databases/database1",
			expectError: false,
		},
		{
			name: "Single database in linked list matching self should return nil",
			linkedDatabaseIds: []string{
				"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Cache/redisEnterprise/example-cluster/databases/database1",
			},
			selfDbId:    "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Cache/redisEnterprise/example-cluster/databases/database1",
			expectError: false,
		},
		{
			name: "Self database not included in linked list should return error",
			linkedDatabaseIds: []string{
				"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Cache/redisEnterprise/example-cluster/databases/database2",
				"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Cache/redisEnterprise/example-cluster/databases/database3",
			},
			selfDbId:            "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Cache/redisEnterprise/example-cluster/databases/database1",
			expectError:         true,
			expectedErrorString: "linked_database_id must include this database ID: /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Cache/redisEnterprise/example-cluster/databases/database1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateLinkedDatabaseIncludesSelf(tc.linkedDatabaseIds, tc.selfDbId)

			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}
				if tc.expectedErrorString != "" && !strings.Contains(err.Error(), tc.expectedErrorString) {
					t.Fatalf("Expected error to contain '%s' but got '%s'", tc.expectedErrorString, err.Error())
				}
			} else if err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}
		})
	}
}
