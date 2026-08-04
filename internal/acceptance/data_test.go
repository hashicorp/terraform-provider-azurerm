// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"regexp"
	"testing"
)

func TestAccAzureRMTestDataRandomIntOfLength(t *testing.T) {
	td := TestData{
		RandomInteger: 112233445566779999,
	}

	cases := []struct {
		len      int
		expected int
	}{
		{
			len:      18,
			expected: 112233445566779999,
		},
		{
			len:      17,
			expected: 11223344556677999,
		},
		{
			len:      16,
			expected: 1122334455667799,
		},
		{
			len:      15,
			expected: 112233445566799,
		},
		{
			len:      14,
			expected: 11223344556699,
		},
		{
			len:      10,
			expected: 1122334499,
		},
		{
			len:      9,
			expected: 112233499,
		},
		{
			len:      8,
			expected: 11223399,
		},
	}

	for _, c := range cases {
		result := td.RandomIntOfLength(c.len)
		if result != c.expected {
			t.Fatalf("For length %d expected %d but got %d", c.len, c.expected, result)
		}
	}
}

const semVerRegex = `^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>alpha|beta|rc)(?:\.(?P<version>0|[1-9]\d*))?)?$`

func TestProviderRelease(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "major",
			input:       "1",
			expectError: true,
		},
		{
			name:        "major.minor",
			input:       "1.2",
			expectError: true,
		},
		{
			name:        "major.minor.patch",
			input:       "1.2.3",
			expectError: false,
		},
		{
			name:        "major.minor.patch.invalid",
			input:       "1.2.3.4",
			expectError: true,
		},
		{
			name:        "default to version file",
			expectError: false,
		},
	}

	for _, c := range cases {
		if !regexp.MustCompile(semVerRegex).MatchString(providerRelease([]string{c.input}...)) && !c.expectError {
			t.Fatalf("expected semver compliant version, got %v", c.input)
		}
		if regexp.MustCompile(semVerRegex).MatchString(providerRelease([]string{c.input}...)) && c.expectError {
			t.Fatalf("expected error but didn't get one for %v", c.input)
		}
	}
}
