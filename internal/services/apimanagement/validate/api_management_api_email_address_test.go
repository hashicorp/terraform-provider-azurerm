// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestEmailAddress(t *testing.T) {
	testData := []struct {
		Value string
		Error bool
	}{
		{
			Value: "a",
			Error: true,
		},
		{
			Value: "abc",
			Error: true,
		},
		{
			Value: "123",
			Error: true,
		},
		{
			Value: "test.com",
			Error: true,
		},
		{
			Value: "test@.com",
			Error: true,
		},
		{
			Value: "test.com",
			Error: true,
		},
		{
			Value: "test@test.com",
			Error: false,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Value)

		_, err := EmailAddress(v.Value, "unit test")
		if err != nil && !v.Error {
			t.Fatalf("Expected pass but got an error: %s", err)
		}
	}
}
