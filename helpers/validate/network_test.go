package validate

import (
	"strconv"
	"testing"
)

func TestCIDR(t *testing.T) {
	cases := []struct {
		CIDR   string
		Errors int
	}{
		{
			CIDR:   "",
			Errors: 1,
		},
		{
			CIDR:   "0.0.0.0",
			Errors: 0,
		},
		{
			CIDR:   "127.0.0.1/8",
			Errors: 0,
		},
		{
			CIDR:   "127.0.0.1/33",
			Errors: 1,
		},
		{
			CIDR:   "127.0.0.1/-1",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.CIDR, func(t *testing.T) {
			_, errors := CIDR(tc.CIDR, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected CIDR to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}

func TestIPv4Address(t *testing.T) {
	cases := []struct {
		IP     string
		Errors int
	}{
		{
			IP:     "",
			Errors: 1,
		},
		{
			IP:     "0.0.0.0",
			Errors: 0,
		},
		{
			IP:     "1.2.3.no",
			Errors: 1,
		},
		{
			IP:     "text",
			Errors: 1,
		},
		{
			IP:     "1.2.3.4",
			Errors: 0,
		},
		{
			IP:     "12.34.43.21",
			Errors: 0,
		},
		{
			IP:     "100.123.199.0",
			Errors: 0,
		},
		{
			IP:     "255.255.255.255",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.IP, func(t *testing.T) {
			_, errors := IPv4Address(tc.IP, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected IPv4Address to return %d error(s) not %d", len(errors), tc.Errors)
			}
		})
	}
}

func TestIPv4AddressOrEmpty(t *testing.T) {
	cases := []struct {
		IP     string
		Errors int
	}{
		{
			IP:     "",
			Errors: 0,
		},
		{
			IP:     "0.0.0.0",
			Errors: 0,
		},
		{
			IP:     "1.2.3.no",
			Errors: 1,
		},
		{
			IP:     "text",
			Errors: 1,
		},
		{
			IP:     "1.2.3.4",
			Errors: 0,
		},
		{
			IP:     "12.34.43.21",
			Errors: 0,
		},
		{
			IP:     "100.123.199.0",
			Errors: 0,
		},
		{
			IP:     "255.255.255.255",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.IP, func(t *testing.T) {
			_, errors := IPv4AddressOrEmpty(tc.IP, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected IPv4AddressOrEmpty to return %d error(s) not %d", len(errors), tc.Errors)
			}
		})
	}
}

func TestPortNumber(t *testing.T) {
	cases := []struct {
		Port   int
		Errors int
	}{
		{
			Port:   -1,
			Errors: 1,
		},
		{
			Port:   0,
			Errors: 1,
		},
		{
			Port:   1,
			Errors: 0,
		},
		{
			Port:   8477,
			Errors: 0,
		},
		{
			Port:   65535,
			Errors: 0,
		},
		{
			Port:   65536,
			Errors: 1,
		},
		{
			Port:   7000000,
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.Port), func(t *testing.T) {
			_, errors := PortNumber(tc.Port, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected PortNumber to return %d error(s) not %d", len(errors), tc.Errors)
			}
		})
	}
}

func TestPortNumberOrZero(t *testing.T) {
	cases := []struct {
		Port   int
		Errors int
	}{
		{
			Port:   -1,
			Errors: 1,
		},
		{
			Port:   0,
			Errors: 0,
		},
		{
			Port:   1,
			Errors: 0,
		},
		{
			Port:   8477,
			Errors: 0,
		},
		{
			Port:   65535,
			Errors: 0,
		},
		{
			Port:   65536,
			Errors: 1,
		},
		{
			Port:   7000000,
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.Port), func(t *testing.T) {
			_, errors := PortNumberOrZero(tc.Port, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected PortNumberOrZero to return %d error(s) not %d", len(errors), tc.Errors)
			}
		})
	}
}
