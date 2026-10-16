// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestWindowsAdminPassword(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// bad: empty
			input:    "",
			expected: false,
		},
		{
			// bad: all lower characters
			input:    "juanjose",
			expected: false,
		},
		{
			// bad: only some caps
			input:    "JuanJose",
			expected: false,
		},
		{
			// bad: can't use reserved words
			input:    "P@$$w0rd",
			expected: false,
		},
		{
			// bad: can't be longer than 123 characters
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghiJklmno9q.stuvwxyzabcdefghijkjlmnopqrstabcdefghijklmnopqrstuvwxyzabcdefghiJklmno9q.stuvwxyzabcdefg",
			expected: false,
		},
		{
			// 72 characters its fine
			input:    "abcdefghijklmnopqrstuvwxyzabcD.9ghijklmnopqrstuvwxyzabcdefghijkjlmnopqrs",
			expected: true,
		},
		{
			// bad: can't be less than 8 characters
			input:    "A9.de",
			expected: false,
		},
		{
			// bad: "_" doesn't count as special character
			input:    "A9BC_7AB",
			expected: false,
		},
		{
			// "/" counts as special character
			input:    "A9BC/7AB",
			expected: true,
		},
		{
			// " " counts as special character
			input:    "A9BC 7AB",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := WindowsAdminPassword(v.input, "admin_password")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestWindowsAdminPasswordDoesNotExposePassword(t *testing.T) {
	for _, field := range []string{"admin_password", "admin_password_wo"} {
		t.Run(field, func(t *testing.T) {
			for _, password := range []string{"abc@123", "P@$$w0rd", "P@ssw0rd", "P@ssword123", "Pa$$word", "pass@word1", "Password!", "Password1", "Password22", "iloveyou!"} {
				_, errors := WindowsAdminPassword(password, field)
				if len(errors) == 0 {
					t.Fatal("expected the disallowed password to be rejected")
				}
				for _, err := range errors {
					if strings.Contains(err.Error(), password) {
						t.Fatal("validation diagnostic contains the supplied password")
					}
					if !strings.Contains(err.Error(), field) {
						t.Fatal("validation diagnostic does not identify the password attribute")
					}
				}
			}
		})
	}
}
