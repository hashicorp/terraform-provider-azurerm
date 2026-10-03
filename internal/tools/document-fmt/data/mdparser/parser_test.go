// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mdparser

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var testDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	// Reference the markdown/testdata directory (shared with other tests)
	baseDir := filepath.Dir(file)
	testDir = filepath.Join(baseDir, "..", "..", "markdown", "testdata")
}

func TestParseMarkdownSection(t *testing.T) {
	testCases := []struct {
		file    string
		itemNum int // expected number of items after tokenization
		argsNum int // expected number of top-level arguments
	}{
		{"key_vault.html.markdown", 65, 16},
		{"media_transform.html.markdown", 270, 5},
	}

	for _, tc := range testCases {
		filePath := filepath.Join(testDir, tc.file)
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read test file %s: %v", tc.file, err)
		}

		// Test tokenization (item count)
		mark := newMarkFromString(string(content))
		if gotItems := len(mark.items); gotItems != tc.itemNum {
			t.Errorf("%s: expected %d items, got %d", tc.file, tc.itemNum, gotItems)
		}

		// Test argument parsing
		lines := strings.Split(string(content), "\n")
		properties := ParseMarkdownSection(lines)
		if gotArgs := len(properties.Names); gotArgs != tc.argsNum {
			t.Errorf("%s: expected %d arguments, got %d", tc.file, tc.argsNum, gotArgs)
		}
	}
}

func TestParseErrors(t *testing.T) {
	testCases := []struct {
		name          string
		content       []string
		expectedError string
		propertyName  string
	}{
		{
			name: "duplicate_property",
			content: []string{
				"* `name` - (Required) The name.",
				"* `name` - (Optional) Duplicate name.",
			},
			expectedError: DuplicateFieldsFound,
			propertyName:  "name",
		},
		{
			name: "block_not_defined",
			content: []string{
				"* `config` - (Required) A `config` block as defined below.",
			},
			expectedError: "block `config` not defined",
			propertyName:  "config",
		},
		{
			name: "misspelled_block_name",
			content: []string{
				"* `setting` - (Required) A `settings` block as defined below.",
				"",
				"---",
				"",
				"A `setting` block supports the following:",
				"",
				"* `value` - (Required) The value.",
			},
			expectedError: MisspelNameOfProperty,
			propertyName:  "setting",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			properties := ParseMarkdownSection(tc.content)

			prop, exists := properties.Objects[tc.propertyName]
			if !exists {
				t.Fatalf("property %s not found", tc.propertyName)
			}

			if len(prop.ParseErrors) == 0 {
				t.Fatalf("expected parse error but got none")
			}

			found := false
			for _, err := range prop.ParseErrors {
				if strings.Contains(err, tc.expectedError) {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected error containing %q, got %v", tc.expectedError, prop.ParseErrors)
			}
		})
	}
}
