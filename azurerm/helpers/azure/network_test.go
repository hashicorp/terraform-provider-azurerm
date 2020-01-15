package azure

import (
	"testing"
)

func TestCompressIPv6Address(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "2607:f8b0:4005:0800:0000:0000:0000:1003",
			expected: "2607:f8b0:4005:800::1003",
		},
		{
			input:    "2001:0db8:1234:0000:0000:0000:0000:0000",
			expected: "2001:db8:1234::",
		},
	}

	for _, v := range cases {
		actual, _ := CompressIPv6Address(v.input)
		if v.expected != actual {
			t.Fatalf("Expected %q but got %q", v.expected, actual)
		}
	}
}
