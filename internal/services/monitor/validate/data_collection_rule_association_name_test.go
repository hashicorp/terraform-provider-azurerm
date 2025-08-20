// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestDataCollectionRuleAssociationName(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "valid-name_123",
			ExpectError: false,
		},
		{
			Input:       "configurationAccessEndpoint",
			ExpectError: false,
		},
		{
			Input:       "test.name",
			ExpectError: false,
		},
		{
			Input:       "test name with spaces",
			ExpectError: false,
		},
		{
			Input:       "test-name-with-dashes",
			ExpectError: false,
		},
		// Test control characters
		{
			Input:       "test\x00name",
			ExpectError: true,
		},
		{
			Input:       "test\x1fname",
			ExpectError: true,
		},
		{
			Input:       "test\x7fname",
			ExpectError: true,
		},
		{
			Input:       "test\nname",
			ExpectError: true,
		},
		{
			Input:       "test\tname",
			ExpectError: true,
		},
		{
			Input:       "test\rname",
			ExpectError: true,
		},
		// Test forbidden characters: < > % & : \ ? /
		{
			Input:       "test<name",
			ExpectError: true,
		},
		{
			Input:       "test>name",
			ExpectError: true,
		},
		{
			Input:       "test%name",
			ExpectError: true,
		},
		{
			Input:       "test&name",
			ExpectError: true,
		},
		{
			Input:       "test:name",
			ExpectError: true,
		},
		{
			Input:       "test\\name",
			ExpectError: true,
		},
		{
			Input:       "test?name",
			ExpectError: true,
		},
		{
			Input:       "test/name",
			ExpectError: true,
		},
		// Test multiple forbidden characters
		{
			Input:       "test<>name",
			ExpectError: true,
		},
		{
			Input:       "test%&name",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			_, errors := DataCollectionRuleAssociationName(tc.Input, "test")

			if tc.ExpectError {
				if len(errors) == 0 {
					t.Fatalf("Expected an error for input %q but got none", tc.Input)
				}
			} else {
				if len(errors) > 0 {
					t.Fatalf("Expected no errors for input %q but got: %v", tc.Input, errors)
				}
			}
		})
	}
}
