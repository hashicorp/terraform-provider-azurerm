package containers

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/accounts"
)

func TestGetResourceManagerResourceID(t *testing.T) {
	actual := Client{}.GetResourceManagerResourceID("11112222-3333-4444-5555-666677778888", "group1", "account1", "container1")
	expected := "/subscriptions/11112222-3333-4444-5555-666677778888/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1"
	if actual != expected {
		t.Fatalf("Expected the Resource Manager Resource ID to be %q but got %q", expected, actual)
	}
}

func TestParseContainerIDStandard(t *testing.T) {
	input := "https://example1.blob.core.windows.net/container1"
	expected := ContainerId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.BlobSubDomainType,
			DomainSuffix:  "core.windows.net",
		},
		ContainerName: "container1",
	}
	actual, err := ParseContainerID(input, "core.windows.net")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if actual.AccountId.AccountName != expected.AccountId.AccountName {
		t.Fatalf("expected AccountName to be %q but got %q", expected.AccountId.AccountName, actual.AccountId.AccountName)
	}
	if actual.AccountId.SubDomainType != expected.AccountId.SubDomainType {
		t.Fatalf("expected SubDomainType to be %q but got %q", expected.AccountId.SubDomainType, actual.AccountId.SubDomainType)
	}
	if actual.AccountId.DomainSuffix != expected.AccountId.DomainSuffix {
		t.Fatalf("expected DomainSuffix to be %q but got %q", expected.AccountId.DomainSuffix, actual.AccountId.DomainSuffix)
	}
	if actual.ContainerName != expected.ContainerName {
		t.Fatalf("expected ContainerName to be %q but got %q", expected.ContainerName, actual.ContainerName)
	}
}

func TestParseContainerIDInADNSZone(t *testing.T) {
	input := "https://example1.zone1.blob.storage.azure.net/container1"
	expected := ContainerId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.BlobSubDomainType,
			DomainSuffix:  "storage.azure.net",
			ZoneName:      pointer.To("zone1"),
		},
		ContainerName: "container1",
	}
	actual, err := ParseContainerID(input, "storage.azure.net")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if actual.AccountId.AccountName != expected.AccountId.AccountName {
		t.Fatalf("expected AccountName to be %q but got %q", expected.AccountId.AccountName, actual.AccountId.AccountName)
	}
	if actual.AccountId.SubDomainType != expected.AccountId.SubDomainType {
		t.Fatalf("expected SubDomainType to be %q but got %q", expected.AccountId.SubDomainType, actual.AccountId.SubDomainType)
	}
	if actual.AccountId.DomainSuffix != expected.AccountId.DomainSuffix {
		t.Fatalf("expected DomainSuffix to be %q but got %q", expected.AccountId.DomainSuffix, actual.AccountId.DomainSuffix)
	}
	if pointer.From(actual.AccountId.ZoneName) != pointer.From(expected.AccountId.ZoneName) {
		t.Fatalf("expected ZoneName to be %q but got %q", pointer.From(expected.AccountId.ZoneName), pointer.From(actual.AccountId.ZoneName))
	}
	if actual.ContainerName != expected.ContainerName {
		t.Fatalf("expected ContainerName to be %q but got %q", expected.ContainerName, actual.ContainerName)
	}
}

func TestParseContainerIDInAnEdgeZone(t *testing.T) {
	input := "https://example1.blob.zone1.edgestorage.azure.net/container1"
	expected := ContainerId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.BlobSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			ZoneName:      pointer.To("zone1"),
			IsEdgeZone:    true,
		},
		ContainerName: "container1",
	}
	actual, err := ParseContainerID(input, "edgestorage.azure.net")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if actual.AccountId.AccountName != expected.AccountId.AccountName {
		t.Fatalf("expected AccountName to be %q but got %q", expected.AccountId.AccountName, actual.AccountId.AccountName)
	}
	if actual.AccountId.SubDomainType != expected.AccountId.SubDomainType {
		t.Fatalf("expected SubDomainType to be %q but got %q", expected.AccountId.SubDomainType, actual.AccountId.SubDomainType)
	}
	if actual.AccountId.DomainSuffix != expected.AccountId.DomainSuffix {
		t.Fatalf("expected DomainSuffix to be %q but got %q", expected.AccountId.DomainSuffix, actual.AccountId.DomainSuffix)
	}
	if pointer.From(actual.AccountId.ZoneName) != pointer.From(expected.AccountId.ZoneName) {
		t.Fatalf("expected ZoneName to be %q but got %q", pointer.From(expected.AccountId.ZoneName), pointer.From(actual.AccountId.ZoneName))
	}
	if !actual.AccountId.IsEdgeZone {
		t.Fatalf("expected the Account to be in an Edge Zone but it wasn't")
	}
	if actual.ContainerName != expected.ContainerName {
		t.Fatalf("expected ContainerName to be %q but got %q", expected.ContainerName, actual.ContainerName)
	}
}

func TestFormatContainerIDStandard(t *testing.T) {
	actual := ContainerId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.BlobSubDomainType,
			DomainSuffix:  "core.windows.net",
			IsEdgeZone:    false,
		},
		ContainerName: "container1",
	}.ID()
	expected := "https://example1.blob.core.windows.net/container1"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatContainerIDInDNSZone(t *testing.T) {
	actual := ContainerId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.BlobSubDomainType,
			DomainSuffix:  "storage.azure.net",
			IsEdgeZone:    false,
		},
		ContainerName: "container1",
	}.ID()
	expected := "https://example1.zone2.blob.storage.azure.net/container1"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatContainerIDInEdgeZone(t *testing.T) {
	actual := ContainerId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.BlobSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			IsEdgeZone:    true,
		},
		ContainerName: "container1",
	}.ID()
	expected := "https://example1.blob.zone2.edgestorage.azure.net/container1"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}
