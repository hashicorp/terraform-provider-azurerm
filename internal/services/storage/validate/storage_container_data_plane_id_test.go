package validate

import (
	"testing"
)

func TestStorageContainerDataPlaneID(t *testing.T) {
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
			Input: "https://blob.core.windows.net/container1",
			Valid: false,
		},

		{
			// missing domain suffix
			Input: "https://storageaccount1/container1",
			Valid: false,
		},

		{
			// missing container name
			Input: "https://storageaccount1.blob.core.windows.net/",
			Valid: false,
		},
		{
			// valid
			Input: "https://storageaccount1.blob.core.windows.net/container1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "HTTPS://STORAGEACCOUNT1.BLOB.CORE.WINDOWS.NET/CONTAINER1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := StorageContainerDataPlaneID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
