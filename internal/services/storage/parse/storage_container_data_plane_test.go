package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = StorageContainerDataPlaneId{}

func TestStorageContainerDataPlaneIDFormatter(t *testing.T) {
	actual := NewStorageContainerDataPlaneId("storageAccount1", "core.windows.net", "container1").ID()
	expected := "https://storageAccount1.blob.core.windows.net/container1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestStorageContainerDataPlaneID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *StorageContainerDataPlaneId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing host
			Input: "https://blob.core.windows.net/container1",
			Error: true,
		},

		{
			// missing domain suffix
			Input: "https://storageAccount1/container1",
			Error: true,
		},

		{
			// missing container name
			Input: "https://storageAccount1.blob.core.windows.net/",
			Error: true,
		},
		{
			// valid
			Input: "https://storageAccount1.blob.core.windows.net/container1",
			Expected: &StorageContainerDataPlaneId{
				AccountName:  "storageAccount1",
				DomainSuffix: "core.windows.net",
				Name:         "container1",
			},
		},

		{
			// upper-cased
			Input: "HTTPS://STORAGEACCOUNT1.BLOB.CORE.WINDOWS.NET/CONTAINER1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := StorageContainerDataPlaneID(v.Input)
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
