package validate

import "testing"

func TestValidatePipelineAndTriggerName(t *testing.T) {
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
			input:    "validname",
			expected: true,
		},
		{
			// can contain numbers in the middle
			input:    "valid02name",
			expected: true,
		},
		{
			// can contain numbers in the end
			input:    "validName1",
			expected: true,
		},
		{
			// can't contain `.`
			input:    "invalid.",
			expected: false,
		},
		{
			// can't contain `:@£`
			input:    ":@£",
			expected: false,
		},
		{
			// can't contain `>`
			input:    ">invalid",
			expected: false,
		},
		{
			// can't contain `&`
			input:    "invalid&name",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := SynapsePipelineAndTriggerName()(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
