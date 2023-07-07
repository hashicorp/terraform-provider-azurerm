// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestDashboardProperties(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "{}",
			expected: false,
		},
		{
			input:    "{\"lenses\":{}}",
			expected: true,
		},
		{
			input:    "{\"lenses\": {\"0\": {\"order\": 0,\"parts\": {\"0\": {\"position\": {\"x\": 0,\"y\": 0,\"rowSpan\": 2,\"colSpan\": 3},\"metadata\": {\"inputs\": [],\"type\": \"Extension/HubsExtension/PartType/MarkdownPart\",\"settings\": {\"content\": {\"settings\": {\"content\": \"## This is only a test :)\",\"subtitle\": \"\",\"title\": \"Test MD Tile\"}}}}}}}}}",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DashboardProperties(v.input, "dashboard_properties")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
