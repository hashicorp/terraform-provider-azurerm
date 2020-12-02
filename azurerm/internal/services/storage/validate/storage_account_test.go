package validate

import "testing"

func TestStorageAccountIPRule(t *testing.T) {
	cases := []struct {
		IPRule string
		Errors int
	}{
		{
			IPRule: "",
			Errors: 1,
		},
		{
			IPRule: "0.0.0.0",
			Errors: 0,
		},
		{
			IPRule: "23.45.0.1/30",
			Errors: 0,
		},
		{
			IPRule: "23.45.0.1/32",
			Errors: 1,
		},
		{
			IPRule: "10.0.0.0/29",
			Errors: 1,
		},
		{
			IPRule: "172.16.0.0/29",
			Errors: 1,
		},
		{
			IPRule: "192.168.0.0/29",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.IPRule, func(t *testing.T) {
			_, errors := StorageAccountIpRule(tc.IPRule, "ip_rules")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected StorageAccountIPRule to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}
