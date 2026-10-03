// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package markdown

import (
	"testing"
)

func TestRemoveRedundantSpace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "multiple spaces in content",
			input:    "* `abc` -  something  here.",
			expected: "* `abc` - something here.",
		},
		{
			name:     "no changes needed",
			input:    "* `abc` - something here.",
			expected: "* `abc` - something here.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeRedundantSpace(tt.input)
			if result != tt.expected {
				t.Errorf("removeRedundantSpace() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestTryFixProp(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "missing space before dash",
			input:    "* `name`- (Required) The name.",
			expected: "* `name` - (Required) The name.",
		},
		{
			name:     "missing space after dash",
			input:    "* `name` -(Required) The name.",
			expected: "* `name` - (Required) The name.",
		},
		{
			name:     "missing dash entirely",
			input:    "* `name` (Required) The name.",
			expected: "* `name` - (Required) The name.",
		},
		{
			name:     "attention marker conversion",
			input:    "** `note` - Important note.",
			expected: "~> ** `note` - Important note.",
		},
		{
			name:     "already correct",
			input:    "* `name` - (Required) The name.",
			expected: "* `name` - (Required) The name.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tryFixProp(tt.input)
			if result != tt.expected {
				t.Errorf("tryFixProp() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestTryBlockHeadDetect(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "typical block header",
			input:    "A `config` block supports the following:",
			expected: true,
		},
		{
			name:     "multiple blocks",
			input:    "The `management`, `portal`, and `scm` blocks support the following:",
			expected: true,
		},
		{
			name:     "not a block header",
			input:    "* `name` - (Required) The name.",
			expected: false,
		},
		{
			name:     "supports with colon",
			input:    "The following configuration supports:",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tryBlockHeadDetect(tt.input)
			if result != tt.expected {
				t.Errorf("tryBlockHeadDetect() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTryFixBlockHead(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "add 'The' prefix",
			input:    "`config` supports the following:",
			expected: "The `config` block supports the following:",
		},
		{
			name:     "add 'block' keyword",
			input:    "A `config` supports the following:",
			expected: "A `config` block supports the following:",
		},
		{
			name:     "remove asterisk prefix",
			input:    "* A `config` supports the following:",
			expected: "A `config` block supports the following:",
		},
		{
			name:     "already has block keyword",
			input:    "A `config` block supports the following:",
			expected: "A `config` block supports the following:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tryFixBlockHead(tt.input)
			if result != tt.expected {
				t.Errorf("tryFixBlockHead() = %q, want %q", result, tt.expected)
			}
		})
	}
}
