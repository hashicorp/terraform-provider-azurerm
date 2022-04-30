package validate

import (
	"testing"
)

func TestStorageTableDataPlaneID(t *testing.T) {
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
			Input: "https://table.core.windows.net/Tables('table1')",
			Valid: false,
		},

		{
			// missing domain suffix
			Input: "https://storageaccount1/Tables('table1')",
			Valid: false,
		},

		{
			// missing table name
			Input: "https://storageaccount1.table.core.windows.net/",
			Valid: false,
		},
		{
			// valid
			Input: "https://storageaccount1.table.core.windows.net/Tables('table1')",
			Valid: true,
		},

		{
			// upper-cased
			Input: "HTTPS://STORAGEACCOUNT1.TABLE.CORE.WINDOWS.NET/TABLES('TABLE1')",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := StorageTableDataPlaneID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
