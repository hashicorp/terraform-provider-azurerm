package accounts

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func TestParseAccountIDStandard(t *testing.T) {
	input := "https://example.blob.core.windows.net"
	expected := AccountId{
		AccountName:   "example",
		SubDomainType: BlobSubDomainType,
		DomainSuffix:  "core.windows.net",
	}
	actual, err := ParseAccountID(input, "core.windows.net")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if actual.AccountName != expected.AccountName {
		t.Fatalf("expected AccountName to be %q but got %q", expected.AccountName, actual.AccountName)
	}
	if actual.SubDomainType != expected.SubDomainType {
		t.Fatalf("expected SubDomainType to be %q but got %q", expected.SubDomainType, actual.SubDomainType)
	}
	if actual.DomainSuffix != expected.DomainSuffix {
		t.Fatalf("expected DomainSuffix to be %q but got %q", expected.DomainSuffix, actual.DomainSuffix)
	}
}

func TestParseAccountIDInDNSZone(t *testing.T) {
	input := "https://example.zone.blob.storage.azure.net"
	expected := AccountId{
		AccountName:   "example",
		ZoneName:      pointer.To("zone"),
		SubDomainType: BlobSubDomainType,
		DomainSuffix:  "storage.azure.net",
	}
	actual, err := ParseAccountID(input, "storage.azure.net")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if actual.AccountName != expected.AccountName {
		t.Fatalf("expected AccountName to be %q but got %q", expected.AccountName, actual.AccountName)
	}
	if actual.ZoneName == nil {
		t.Fatalf("expected ZoneName to have a value but got nil")
	}
	if *actual.ZoneName != *expected.ZoneName {
		t.Fatalf("expected ZoneName to be %q but got %q", *expected.ZoneName, *actual.ZoneName)
	}
	if actual.SubDomainType != expected.SubDomainType {
		t.Fatalf("expected SubDomainType to be %q but got %q", expected.SubDomainType, actual.SubDomainType)
	}
	if actual.DomainSuffix != expected.DomainSuffix {
		t.Fatalf("expected DomainSuffix to be %q but got %q", expected.DomainSuffix, actual.DomainSuffix)
	}
	if actual.IsEdgeZone {
		t.Fatalf("expected IsEdgeZone to be false but got %t", actual.IsEdgeZone)
	}
}

func TestParseAccountIDInEdgeZone(t *testing.T) {
	input := "https://example.blob.danger.edgestorage.azure.net"
	expected := AccountId{
		AccountName:   "example",
		ZoneName:      pointer.To("danger"),
		IsEdgeZone:    true,
		SubDomainType: BlobSubDomainType,
		DomainSuffix:  "edgestorage.azure.net",
	}
	actual, err := ParseAccountID(input, "edgestorage.azure.net")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if actual.AccountName != expected.AccountName {
		t.Fatalf("expected AccountName to be %q but got %q", expected.AccountName, actual.AccountName)
	}
	if actual.ZoneName == nil {
		t.Fatalf("expected ZoneName to have a value but got nil")
	}
	if *actual.ZoneName != *expected.ZoneName {
		t.Fatalf("expected ZoneName to be %q but got %q", *expected.ZoneName, *actual.ZoneName)
	}
	if actual.SubDomainType != expected.SubDomainType {
		t.Fatalf("expected SubDomainType to be %q but got %q", expected.SubDomainType, actual.SubDomainType)
	}
	if !actual.IsEdgeZone {
		t.Fatalf("expected IsEdgeZone to be true but got %t", actual.IsEdgeZone)
	}
}

func TestFormatAccountIDStandard(t *testing.T) {
	actual := AccountId{
		AccountName:   "example1",
		ZoneName:      nil,
		SubDomainType: FileSubDomainType,
		DomainSuffix:  "core.windows.net",
		IsEdgeZone:    false,
	}.ID()
	expected := "https://example1.file.core.windows.net"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatAccountIDInDNSZone(t *testing.T) {
	actual := AccountId{
		AccountName:   "example1",
		ZoneName:      pointer.To("zone1"),
		SubDomainType: FileSubDomainType,
		DomainSuffix:  "storage.azure.net",
		IsEdgeZone:    false,
	}.ID()
	expected := "https://example1.zone1.file.storage.azure.net"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatAccountIDInEdgeZone(t *testing.T) {
	actual := AccountId{
		AccountName:   "example1",
		ZoneName:      pointer.To("zone1"),
		SubDomainType: FileSubDomainType,
		DomainSuffix:  "edgestorage.azure.net",
		IsEdgeZone:    true,
	}.ID()
	expected := "https://example1.file.zone1.edgestorage.azure.net"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}
