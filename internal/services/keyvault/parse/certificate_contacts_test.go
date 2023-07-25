// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = CertificateContactsId{}

func TestCertificateContactsIDFormatter(t *testing.T) {
	actual, err := NewCertificateContactsID("https://example-keyvault.vault.azure.net")
	if err != nil {
		t.Fatalf("Error occurred when creating ID: %+v", err)
	}
	expected := "https://example-keyvault.vault.azure.net/certificates/contacts"
	if actual.ID() != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestCertificateContactsID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *CertificateContactsId
	}{
		{
			// valid
			Input: "https://example-keyvault.vault.azure.net/certificates/contacts",
			Expected: &CertificateContactsId{
				KeyVaultBaseUrl: "https://example-keyvault.vault.azure.net/",
			},
		},
		{
			// empty
			Input: "",
			Error: true,
		},
		{
			// missing suffix
			Input: "https://my-keyvault.vault.azure.net",
			Error: true,
		},
		{
			// wrong suffix
			Input: "https://my-keyvault.vault.azure.net/certificates/contact",
			Error: true,
		},
		{
			// additional item in path
			Input: "https://my-keyvault.vault.azure.net/x/certificates/contact",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := CertificateContactsID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.KeyVaultBaseUrl != v.Expected.KeyVaultBaseUrl {
			t.Fatalf("Expected %q but got %q for KeyVaultBaseUrl", v.Expected.KeyVaultBaseUrl, actual.KeyVaultBaseUrl)
		}
	}
}
