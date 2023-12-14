package blobs

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/accounts"
)

func TestParseBlobIDStandard(t *testing.T) {
	input := "https://example1.blob.core.windows.net/container1/blob1.vhd"
	expected := BlobId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.BlobSubDomainType,
			DomainSuffix:  "core.windows.net",
		},
		ContainerName: "container1",
		BlobName:      "blob1.vhd",
	}
	actual, err := ParseBlobID(input, "core.windows.net")
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
	if actual.BlobName != expected.BlobName {
		t.Fatalf("expected BlobName to be %q but got %q", expected.BlobName, actual.BlobName)
	}
}

func TestParseBlobIDInADNSZone(t *testing.T) {
	input := "https://example1.zone1.blob.storage.azure.net/container1/blob1.vhd"
	expected := BlobId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.BlobSubDomainType,
			DomainSuffix:  "storage.azure.net",
			ZoneName:      pointer.To("zone1"),
		},
		ContainerName: "container1",
		BlobName:      "blob1.vhd",
	}
	actual, err := ParseBlobID(input, "storage.azure.net")
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
	if actual.BlobName != expected.BlobName {
		t.Fatalf("expected BlobName to be %q but got %q", expected.BlobName, actual.BlobName)
	}
}

func TestParseBlobIDInAnEdgeZone(t *testing.T) {
	input := "https://example1.blob.zone1.edgestorage.azure.net/container1/blob1.vhd"
	expected := BlobId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.BlobSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			ZoneName:      pointer.To("zone1"),
			IsEdgeZone:    true,
		},
		ContainerName: "container1",
		BlobName:      "blob1.vhd",
	}
	actual, err := ParseBlobID(input, "edgestorage.azure.net")
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
	if actual.BlobName != expected.BlobName {
		t.Fatalf("expected BlobName to be %q but got %q", expected.BlobName, actual.BlobName)
	}
}

func TestFormatBlobIDStandard(t *testing.T) {
	actual := BlobId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.BlobSubDomainType,
			DomainSuffix:  "core.windows.net",
			IsEdgeZone:    false,
		},
		ContainerName: "container1",
		BlobName:      "somefile.vhd",
	}.ID()
	expected := "https://example1.blob.core.windows.net/container1/somefile.vhd"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatBlobIDInDNSZone(t *testing.T) {
	actual := BlobId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.BlobSubDomainType,
			DomainSuffix:  "storage.azure.net",
			IsEdgeZone:    false,
		},
		ContainerName: "container1",
		BlobName:      "somefile.vhd",
	}.ID()
	expected := "https://example1.zone2.blob.storage.azure.net/container1/somefile.vhd"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatBlobIDInEdgeZone(t *testing.T) {
	actual := BlobId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.BlobSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			IsEdgeZone:    true,
		},
		ContainerName: "container1",
		BlobName:      "somefile.vhd",
	}.ID()
	expected := "https://example1.blob.zone2.edgestorage.azure.net/container1/somefile.vhd"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}
