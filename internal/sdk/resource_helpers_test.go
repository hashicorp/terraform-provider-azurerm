// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func TestParseStructTags_Empty(t *testing.T) {
	actual, err := parseStructTags("")
	if err != nil {
		t.Fatalf("unexpected error %q", err.Error())
	}

	if actual != nil {
		t.Fatalf("expected actual to be nil but got %+v", *actual)
	}
}

func TestParseStructTags_WithValue(t *testing.T) {
	testData := []struct {
		input    reflect.StructTag
		expected *decodedStructTags
		error    *string
	}{
		{
			// empty hclPath
			input:    `tfschema:""`,
			expected: nil,
			error:    pointer.To("the `tfschema` struct tag was defined but empty"),
		},
		{
			// valid, no removedInNextMajorVersion
			input: `tfschema:"hello"`,
			expected: &decodedStructTags{
				hclPath:                   "hello",
				addedInNextMajorVersion:   false,
				removedInNextMajorVersion: false,
			},
		},
		{
			// valid, with removedInNextMajorVersion
			input: `tfschema:"hello,removedInNextMajorVersion"`,
			expected: &decodedStructTags{
				hclPath:                   "hello",
				addedInNextMajorVersion:   false,
				removedInNextMajorVersion: true,
			},
		},
		{
			// valid, with removedInNextMajorVersion and a space before the comma
			input: `tfschema:"hello, removedInNextMajorVersion"`,
			expected: &decodedStructTags{
				hclPath:                   "hello",
				addedInNextMajorVersion:   false,
				removedInNextMajorVersion: true,
			},
		},
		{
			// valid, with removedInNextMajorVersion and a space after the comma
			//
			// This would be caught in PR review, but would be a confusing error/experience
			// during development so it's worth being lenient here since it's non-impactful
			input: `tfschema:"hello ,removedInNextMajorVersion"`,
			expected: &decodedStructTags{
				hclPath:                   "hello",
				addedInNextMajorVersion:   false,
				removedInNextMajorVersion: true,
			},
		},
		{
			// valid, with removedInNextMajorVersion and a space either side
			//
			// This would be caught in PR review, but would be a confusing error/experience
			// during development so it's worth being lenient here since it's non-impactful
			input: `tfschema:"hello , removedInNextMajorVersion"`,
			expected: &decodedStructTags{
				hclPath:                   "hello",
				addedInNextMajorVersion:   false,
				removedInNextMajorVersion: true,
			},
		},
		{
			// valid, with addedInNextMajorVersion and a space before the comma
			input: `tfschema:"hello, addedInNextMajorVersion"`,
			expected: &decodedStructTags{
				hclPath:                   "hello",
				addedInNextMajorVersion:   true,
				removedInNextMajorVersion: false,
			},
		},
		{
			// valid, with addedInNextMajorVersion and a space before the comma
			input:    `tfschema:"hello, removedInNextMajorVersion, addedInNextMajorVersion"`,
			expected: nil,
			error:    pointer.To("the struct-tags `removedInNextMajorVersion` and `addedInNextMajorVersion` cannot be set together"),
		},
		{
			// valid, with addedInNextMajorVersion and a space before the comma
			input:    `tfschema:"hello, addedInNextMajorVersion, removedInNextMajorVersion"`,
			expected: nil,
			error:    pointer.To("the struct-tags `removedInNextMajorVersion` and `addedInNextMajorVersion` cannot be set together"),
		},
		{
			// invalid, unknown struct tags
			input:    `tfschema:"hello,world"`,
			expected: nil,
			error:    pointer.To(`internal-error: the struct-tag "world" is not implemented - struct tags are "hello,world"`),
		},
	}
	for i, data := range testData {
		t.Logf("Index %d - Input %q", i, data.input)
		actual, err := parseStructTags(data.input)
		if err != nil {
			if data.error != nil {
				if err.Error() == *data.error {
					continue
				}

				t.Fatalf("expected the error %q but got %q", *data.error, err.Error())
			}

			t.Fatalf("unexpected error %q", err.Error())
		}
		if data.error != nil {
			t.Fatalf("expected the error %q but didn't get one", *data.error)
		}

		if actual == nil {
			t.Fatalf("expected actual to have a value but got nil")
		}
		if !reflect.DeepEqual(*data.expected, *actual) {
			t.Fatalf("expected [%+v] and actual [%+v] didn't match", *data.expected, *actual)
		}
	}
}
