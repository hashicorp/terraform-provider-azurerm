package validate

import (
	"testing"
)

func TestStorageShareDataPlaneID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// missing host
			Input: "https://file.core.windows.net/share1",
			Valid: false,
		},

		{
			// missing domain suffix
			Input: "https://storageaccount1/share1",
			Valid: false,
		},

		{
			// missing share name
			Input: "https://storageaccount1.file.core.windows.net/",
			Valid: false,
		},
		{
			// valid
			Input: "https://storageaccount1.file.core.windows.net/share1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "HTTPS://STORAGEACCOUNT1.FILE.CORE.WINDOWS.NET/SHARE1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := StorageShareDataPlaneID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
