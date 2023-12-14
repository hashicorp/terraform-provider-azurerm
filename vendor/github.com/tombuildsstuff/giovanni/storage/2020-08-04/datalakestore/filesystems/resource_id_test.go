package filesystems

import (
	"testing"
)

func TestGetResourceManagerResourceID(t *testing.T) {
	actual := Client{}.GetResourceManagerResourceID("11112222-3333-4444-5555-666677778888", "group1", "account1", "container1")
	expected := "/subscriptions/11112222-3333-4444-5555-666677778888/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1"
	if actual != expected {
		t.Fatalf("Expected the Resource Manager Resource ID to be %q but got %q", expected, actual)
	}
}
