// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package markdown

import (
	"testing"
)

func TestNormalizeArgumentsContent(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expected   string
		wantChange bool
	}{
		{
			name:       "missing space after backtick",
			input:      "* `name`- (Required) The name.",
			expected:   "* `name` - (Required) The name.",
			wantChange: true,
		},
		{
			name:       "missing space after backtick with parenthesis",
			input:      "* `name`-(Required) The name.",
			expected:   "* `name` - (Required) The name.",
			wantChange: true,
		},
		{
			name:       "lowercase optional",
			input:      "* `name` - (optional) The name.",
			expected:   "* `name` - (Optional) The name.",
			wantChange: true,
		},
		{
			name:       "lowercase required",
			input:      "* `id` - (required) The ID.",
			expected:   "* `id` - (Required) The ID.",
			wantChange: true,
		},
		{
			name:       "wrong order optional",
			input:      "* `name` (Optional) - The name.",
			expected:   "* `name` - (Optional) The name.",
			wantChange: true,
		},
		{
			name:       "already correct",
			input:      "* `name` - (Required) The name.",
			expected:   "* `name` - (Required) The name.",
			wantChange: false,
		},
		{
			name:       "dash list to asterisk",
			input:      "- `name` - (Required) The name.",
			expected:   "* `name` - (Required) The name.",
			wantChange: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasChange := normalizeArgumentsContent([]string{tt.input})

			if len(result) != 1 {
				t.Fatalf("expected 1 line, got %d", len(result))
			}

			if result[0] != tt.expected {
				t.Errorf("normalizeArgumentsContent() = %q, want %q", result[0], tt.expected)
			}

			if hasChange != tt.wantChange {
				t.Errorf("normalizeArgumentsContent() hasChange = %v, want %v", hasChange, tt.wantChange)
			}
		})
	}
}
