package validate

import "testing"

func TestHelper_Validate_Ip4Address(t *testing.T) {
	cases := []struct {
		Ip     string
		Errors int
	}{
		{
			Ip:     "",
			Errors: 1,
		},
		{
			Ip:     "000.000.000.000",
			Errors: 1,
		},
		{
			Ip:     "1.2.3.no",
			Errors: 1,
		},
		{
			Ip:     "text",
			Errors: 1,
		},
		{
			Ip:     "1.2.3.4",
			Errors: 0,
		},
		{
			Ip:     "12.34.43.21",
			Errors: 0,
		},
		{
			Ip:     "100.123.199.0",
			Errors: 0,
		},
		{
			Ip:     "255.255.255.255",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		_, errors := Ip4Address(tc.Ip, "test")

		if len(errors) < tc.Errors {
			t.Fatalf("Expected AzureResourceId to have an error for %q", tc.Ip)
		}
	}
}

func TestHelper_Validate_MacAddress(t *testing.T) {
	cases := []struct {
		Ip     string
		Errors int
	}{
		{
			Ip:     "",
			Errors: 1,
		},
		{
			Ip:     "text",
			Errors: 1,
		},
		{
			Ip:     "12:34:no",
			Errors: 1,
		},
		{
			Ip:     "123:34:56:78:90:ab",
			Errors: 1,
		},
		{
			Ip:     "12:34:56:78:90:NO",
			Errors: 1,
		},
		{
			Ip:     "12:34:56:78:90:ab",
			Errors: 0,
		},
		{
			Ip:     "ab:cd:ef:AB:CD:EF",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		_, errors := Ip4Address(tc.Ip, "test")

		if len(errors) < tc.Errors {
			t.Fatalf("Expected AzureResourceId to have an error for %q", tc.Ip)
		}
	}
}
