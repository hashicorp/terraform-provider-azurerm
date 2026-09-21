// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"encoding/base64"
	"testing"
)

func TestBase64EncodeIfNot(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text",
			input:    "hello world",
			expected: base64.StdEncoding.EncodeToString([]byte("hello world")),
		},
		{
			name:     "already base64 encoded",
			input:    base64.StdEncoding.EncodeToString([]byte("hello world")),
			expected: base64.StdEncoding.EncodeToString([]byte("hello world")),
		},
		{
			name:     "empty string (valid base64)",
			input:    "",
			expected: "",
		},
		{
			name:     "invalid base64 string",
			input:    "this is not base64 !!!",
			expected: base64.StdEncoding.EncodeToString([]byte("this is not base64 !!!")),
		},
		{
			name:     "valid base64 string that happens to be an English word",
			input:    "face", // "face" is valid base64 since length is 4 and characters are valid
			expected: "face",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Base64EncodeIfNot(tc.input)
			if actual != tc.expected {
				t.Fatalf("expected: %q, got: %q", tc.expected, actual)
			}
		})
	}
}

func TestBase64IsEncoded(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid base64",
			input:    base64.StdEncoding.EncodeToString([]byte("hello world")),
			expected: true,
		},
		{
			name:     "invalid base64 characters",
			input:    "hello world!",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "valid base64 English word",
			input:    "face",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := base64IsEncoded(tc.input)
			if actual != tc.expected {
				t.Fatalf("expected: %v, got: %v", tc.expected, actual)
			}
		})
	}
}
