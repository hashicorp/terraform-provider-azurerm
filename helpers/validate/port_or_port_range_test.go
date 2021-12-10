package validate

import (
	"testing"
)

func TestPortOrPortRangeWithin(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "0",
			expected: false,
		},
		{
			input:    "65536",
			expected: false,
		},
		{
			input:    "1",
			expected: true,
		},
		{
			input:    "65535",
			expected: true,
		},
		{
			input:    "634",
			expected: true,
		},
		{
			input:    "1000-50000",
			expected: true,
		},
		{
			input:    "1-65535",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := PortOrPortRangeWithin(1, 65535)(v.input, "port_or_port_range")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
