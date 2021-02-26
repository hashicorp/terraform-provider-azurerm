package validate

import (
	"strings"
	"testing"
)

func TestManagedApplicationDefinitionName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "h",
			expected: false,
		},
		{
			input:    "he",
			expected: false,
		},
		{
			input:    "hel",
			expected: true,
		},
		{
			input:    "hel2",
			expected: true,
		},
		{
			input:    "_hello",
			expected: false,
		},
		{
			input:    "hello-",
			expected: false,
		},
		{
			input:    "malcolm-in!the-middle",
			expected: false,
		},
		{
			input:    strings.Repeat("a", 63),
			expected: true,
		},
		{
			input:    strings.Repeat("a", 64),
			expected: true,
		},
		{
			input:    strings.Repeat("a", 65),
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ApplicationDefinitionName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestManagedApplicationDefinitionDisplayName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "hel",
			expected: false,
		},
		{
			input:    "hell",
			expected: true,
		},
		{
			input:    "hello",
			expected: true,
		},
		{
			input:    strings.Repeat("a", 59),
			expected: true,
		},
		{
			input:    strings.Repeat("a", 60),
			expected: true,
		},
		{
			input:    strings.Repeat("a", 61),
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ApplicationDefinitionDisplayName(v.input, "display_name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestManagedApplicationDefinitionDescription(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: true,
		},
		{
			input:    strings.Repeat("a", 199),
			expected: true,
		},
		{
			input:    strings.Repeat("a", 200),
			expected: true,
		},
		{
			input:    strings.Repeat("a", 201),
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ApplicationDefinitionDescription(v.input, "description")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
