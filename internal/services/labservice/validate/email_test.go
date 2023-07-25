// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestEmail(t *testing.T) {
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
			Value: "testuser@contoso.com",
			Error: false,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Value)

		_, err := Email(v.Value, "email")
		if err != nil && !v.Error {
			t.Fatalf("Expected pass but got an error: %s", err)
		}
	}
}
