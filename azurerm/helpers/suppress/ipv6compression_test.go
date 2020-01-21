package suppress

import "testing"

func TestIPv6Compression(t *testing.T) {
	cases := []struct {
		Name     string
		IPv6A    string
		IPv6B    string
		Suppress bool
	}{
		{
			Name:     "both empty",
			IPv6A:    "",
			IPv6B:    "",
			Suppress: false,
		},
		{
			Name:     "empty vs IPv6",
			IPv6A:    "2001:0db8:85a3:0:0:8a2e:0370:7334",
			IPv6B:    "",
			Suppress: false,
		},
		{
			Name:     "invalid IPv6",
			IPv6A:    "2001::invalid",
			IPv6B:    "",
			Suppress: false,
		},
		{
			Name:     "different IPv6",
			IPv6A:    "2001:0db8:85a3:0:0:8a2e:0370:7334",
			IPv6B:    "2001:0db8:85a3:0:0:0:1:2",
			Suppress: false,
		},
		{
			Name:     "same IPv6",
			IPv6A:    "2001:0db8:85a3:0:0:8a2e:0370:7334",
			IPv6B:    "2001:0db8:85a3:0:0:8a2e:0370:7334",
			Suppress: true,
		},
		{
			Name:     "same compressed IPv6",
			IPv6A:    "2001:0db8:85a3:0:0:8a2e:0370:7334",
			IPv6B:    "2001:db8:85a3::8a2e:370:7334",
			Suppress: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if IPv6Compression("test", tc.IPv6A, tc.IPv6B, nil) != tc.Suppress {
				t.Fatalf("Expected IPv6Compression to return %t for '%q' == '%q'", tc.Suppress, tc.IPv6A, tc.IPv6B)
			}
		})
	}
}
