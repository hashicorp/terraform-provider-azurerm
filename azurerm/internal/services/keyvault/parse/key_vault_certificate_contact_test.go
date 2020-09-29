package parse

import (
	"testing"
)

func TestKeyVaultCertificateContactID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *KeyVaultCertificateContactId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "invalid url",
			Input:    "https://",
			Expected: nil,
		},
		{
			Name:     "missing path part",
			Input:    "https://example-keyvault.vault.azure.net",
			Expected: nil,
		},
		{
			Name:     "missing contacts part",
			Input:    "https://example-keyvault.vault.azure.net/certificates",
			Expected: nil,
		},
		{
			Name:  "Certificate Contact ID",
			Input: "https://example-keyvault.vault.azure.net/certificates/contacts",
			Expected: &KeyVaultCertificateContactId{
				KeyVaultBaseUrl: "https://example-keyvault.vault.azure.net/",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "https://example-keyvault.vault.azure.net/certificates/Contacts",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := KeyVaultCertificateContactID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.KeyVaultBaseUrl != v.Expected.KeyVaultBaseUrl {
			t.Fatalf("Expected %q but got %q for Key Vault Base Url", v.Expected.KeyVaultBaseUrl, actual.KeyVaultBaseUrl)
		}
	}
}
