package parse

import (
	"fmt"
	"net/url"
	"strings"
)

type KeyVaultCertificateContactId struct {
	KeyVaultBaseUrl string
}

func KeyVaultCertificateContactID(id string) (*KeyVaultCertificateContactId, error) {
	// example: https://example-keyvault.vault.azure.net/certificates/contacts
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("parse Azure KeyVault Certificate Contact Id: %s", err)
	}

	path := strings.TrimSuffix(idURL.Path, "/")
	if path != "/certificates/contacts" {
		return nil, fmt.Errorf("keyVault Certificate Contact ID path must be '/certificates/contacts', got %q", path)
	}

	contactId := KeyVaultCertificateContactId{
		KeyVaultBaseUrl: fmt.Sprintf("%s://%s/", idURL.Scheme, idURL.Host),
	}

	return &contactId, nil
}
