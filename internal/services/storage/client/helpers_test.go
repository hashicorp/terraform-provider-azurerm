package client

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

func TestAccountDetails_DataPlaneEndpoint(t *testing.T) {
	StorageDomainSuffix = pointer.To("core.windows.net")
	ad := AccountDetails{
		StorageAccountId: commonids.StorageAccountId{
			StorageAccountName: "exampleacct1234",
		},
	}

	uri, err := ad.DataPlaneEndpoint(EndpointTypeQueue)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	expected := "https://exampleacct1234.queue.core.windows.net"
	if *uri != expected {
		t.Fatalf("expected URI to be %q, got %q", expected, *uri)
	}
}
