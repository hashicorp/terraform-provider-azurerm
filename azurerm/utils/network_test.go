package utils

import (
	"testing"
)

func TestIPv6Compression(t *testing.T) {
	cases := []struct {
		Name   string
		Input  interface{}
		Output string
		Valid  bool
	}{
		{
			Name:   "input empty",
			Input:  "",
			Output: "",
			Valid:  true,
		},
		{
			Name:   "valid IPv6 input, invalid compression",
			Input:  "2001:0db8:85a3:0:0:8a2e:0370:7334",
			Output: "2001:0db8:85a3:0:0:8a2e:0370:7334",
			Valid:  false,
		},
		{
			Name:   "invalid IPv6 input",
			Input:  "2001::invalid",
			Output: "",
			Valid:  false,
		},
		{
			Name:   "valid IPv6 compression",
			Input:  "2001:0db8:85a3:0:0:8a2e:0370:7334",
			Output: "2001:db8:85a3::8a2e:370:7334",
			Valid:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			r := NormalizeIPv6Address(tc.Input)
			if (r != tc.Output) && tc.Valid {
				t.Fatalf("Expected NormalizeIPv6Address to return '%q' for '%q' (got '%q')", tc.Output, tc.Input, r)
			}
		})
	}
}
