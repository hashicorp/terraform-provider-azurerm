package validate

import "testing"

func TestDataBoxJobDiskPassKey(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "     ",
			expected: false,
		},
		{
			input:    "hellohellohello!1",
			expected: true,
		},
		{
			input:    "123123123123123!1",
			expected: true,
		},
		{
			input:    "2@hellohellohello",
			expected: true,
		},
		{
			input:    "2@hellohe5$llohello3&",
			expected: true,
		},
		{
			input:    "hellohellohello2",
			expected: false,
		},
		{
			input:    "hellohellohello#",
			expected: false,
		},
		{
			input:    "#hellohellohello",
			expected: false,
		},
		{
			input:    "2hellohellohello",
			expected: false,
		},
		{
			input:    "malcolm-in-the-middle",
			expected: false,
		},
		{
			input:    "malcolm1in2the3middle",
			expected: false,
		},
		{
			input:    "hellohellohello2",
			expected: false,
		},
		{
			input:    "bsadfasdfsdf2)bsadfasdfsdf2)acc",
			expected: true,
		},
		{
			input:    "bsadfasdfsdf2)bsadfasdfsdf2)accs",
			expected: true,
		},
		{
			input:    "bsadfasdfsdf2)bsadfasdfsdf2)accss",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DataBoxJobDiskPassKey(v.input, "databox_disk_passkey")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
