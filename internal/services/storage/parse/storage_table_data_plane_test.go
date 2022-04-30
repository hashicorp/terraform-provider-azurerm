package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = StorageTableDataPlaneId{}

func TestStorageTableDataPlaneIDFormatter(t *testing.T) {
	actual := NewStorageTableDataPlaneId("storageAccount1", "core.windows.net", "table1").ID()
	expected := "https://storageAccount1.table.core.windows.net/Tables('table1')"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestStorageTableDataPlaneID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *StorageTableDataPlaneId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing host
			Input: "https://table.core.windows.net/Tables('table1')",
			Error: true,
		},

		{
			// missing domain suffix
			Input: "https://storageAccount1/Tables('table1')",
			Error: true,
		},

		{
			// missing table1 name
			Input: "https://storageAccount1.table.core.windows.net/",
			Error: true,
		},
		{
			// valid
			Input: "https://storageAccount1.table.core.windows.net/Tables('table1')",
			Expected: &StorageTableDataPlaneId{
				AccountName:  "storageAccount1",
				DomainSuffix: "core.windows.net",
				Name:         "table1",
			},
		},

		{
			// upper-cased
			Input: "HTTPS://STORAGEACCOUNT1.TABLE.CORE.WINDOWS.NET/TABLES('TABLE')",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := StorageTableDataPlaneID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.AccountName != v.Expected.AccountName {
			t.Fatalf("Expected %q but got %q for AccountName", v.Expected.AccountName, actual.AccountName)
		}
		if actual.DomainSuffix != v.Expected.DomainSuffix {
			t.Fatalf("Expected %q but got %q for DomainSuffix", v.Expected.DomainSuffix, actual.DomainSuffix)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
