package validate

import (
	"testing"
)

func TestActionRuleName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// basic example, lead by lower case letter
			input:    "a",
			expected: true,
		},
		{
			// basic example, lead by upper case letter
			input:    "A",
			expected: true,
		},
		{
			// basic example, lead by number
			input:    "8",
			expected: true,
		},
		{
			// basic example, contain underscore
			input:    "a_b",
			expected: true,
		},
		{
			// basic example, end with underscore
			input:    "ab_",
			expected: true,
		},
		{
			// basic example, end with hyphen
			input:    "ab-",
			expected: true,
		},
		{
			// can not contain '+'
			input:    "a+",
			expected: false,
		},
		{
			// can't lead by '-'
			input:    "-a",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ActionRuleName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestActionRuleScheduleDate(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// basic example
			input:    "04/14/2020",
			expected: true,
		},
		{
			// month format is wrong
			input:    "4/14/2020",
			expected: false,
		},
		{
			// can't use '-'
			input:    "04-14-2020",
			expected: false,
		},
		{
			// the format can not be YYYY/MM/DD
			input:    "2020/04/14",
			expected: false,
		},
		{
			// year format is wrong
			input:    "04/14/20",
			expected: false,
		},
		{
			// can not contain time
			input:    "04/14/2020 08:08:08",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ActionRuleScheduleDate(v.input, "date")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestActionRuleScheduleTime(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// basic example
			input:    "06:00:00",
			expected: true,
		},
		{
			// basic example 2
			input:    "16:00:00",
			expected: true,
		},
		{
			// basic example 3
			input:    "6:00:00",
			expected: true,
		},
		{
			// can't use '-'
			input:    "06-00-00",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ActionRuleScheduleTime(v.input, "time")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
