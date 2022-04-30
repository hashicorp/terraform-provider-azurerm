package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = StorageShareDataPlaneId{}

func TestStorageShareDataPlaneIDFormatter(t *testing.T) {
	actual := NewStorageShareDataPlaneId("storageAccount1", "core.windows.net", "share1").ID()
	expected := "https://storageAccount1.file.core.windows.net/share1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestStorageShareDataPlaneID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *StorageShareDataPlaneId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing host
			Input: "https://file.core.windows.net/share1",
			Error: true,
		},

		{
			// missing domain suffix
			Input: "https://storageAccount1/share1",
			Error: true,
		},

		{
			// missing share name
			Input: "https://storageAccount1.file.core.windows.net/",
			Error: true,
		},
		{
			// valid
			Input: "https://storageAccount1.file.core.windows.net/share1",
			Expected: &StorageShareDataPlaneId{
				AccountName:  "storageAccount1",
				DomainSuffix: "core.windows.net",
				Name:         "share1",
			},
		},

		{
			// upper-cased
			Input: "HTTPS://STORAGEACCOUNT1.FILE.CORE.WINDOWS.NET/SHARE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := StorageShareDataPlaneID(v.Input)
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
