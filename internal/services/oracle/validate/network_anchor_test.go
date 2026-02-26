// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestIsCommaSeparatedCIDRs(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "10.0.0.0/24",
			expected: true,
		},
		{
			input:    "10.0.0.0/24,,10.0.1.0/24",
			expected: false,
		},
		{
			input:    "10.0.0",
			expected: false,
		},
		{
			input:    "10.0.0.10",
			expected: false,
		},
		{
			input:    "345.123.10.10",
			expected: false,
		},
		{
			input:    "10.0.0.0/24, 10.0.0.1/24 ,   10.0.1.0/16",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := IsCommaSeparatedCIDRs(v.input, "cidrs")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestDomainNames(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "test-domain",
			expected: false,
		},
		{
			input:    "good-domain.com,,example.net",
			expected: false,
		},
		{
			input:    "good-domain.com, test.example.com,",
			expected: false,
		},
		{
			input:    "bad-do#main.com, test.exam&ple.com",
			expected: false,
		},
		{
			input:    "-bad123domain.com",
			expected: false,
		},
		{
			input:    "good-domain.com",
			expected: true,
		},
		{
			input:    "good-domain123.com, abc.def.test.example.com, example.net",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DomainNames(v.input, "domain_names")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
