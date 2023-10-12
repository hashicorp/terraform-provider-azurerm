// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"
)

type CertificateContactsId struct {
	KeyVaultBaseUrl string
}

func NewCertificateContactsID(keyVaultBaseUrl string) (*CertificateContactsId, error) {
	// example: https://example-keyvault.vault.azure.net/certificates/contacts
	keyVaultUrl, err := url.Parse(keyVaultBaseUrl)
	if err != nil || keyVaultBaseUrl == "" {
		return nil, fmt.Errorf("parsing %q: %+v", keyVaultBaseUrl, err)
	}

	if hostParts := strings.Split(keyVaultUrl.Host, ":"); len(hostParts) > 1 {
		keyVaultUrl.Host = hostParts[0]
	}

	return &CertificateContactsId{
		KeyVaultBaseUrl: keyVaultUrl.String(),
	}, nil
}

func (id CertificateContactsId) String() string {
	components := []string{
		fmt.Sprintf("Base Url %q", id.KeyVaultBaseUrl),
	}
	return fmt.Sprintf("Key Vault Certificate Contacts: (%s)", strings.Join(components, " / "))
}

func (id CertificateContactsId) ID() string {
	// example: https://example-keyvault.vault.azure.net/certificates/contacts
	segments := []string{
		strings.TrimSuffix(id.KeyVaultBaseUrl, "/"),
		"certificates",
		"contacts",
	}
	return strings.TrimSuffix(strings.Join(segments, "/"), "/")
}

func CertificateContactsID(input string) (*CertificateContactsId, error) {
	idURL, err := url.ParseRequestURI(input)
	if err != nil {
		return nil, fmt.Errorf("cannot parse Azure Key Vault Certificate Contacts Id: %s", err)
	}

	path := idURL.Path
	path = strings.TrimSuffix(path, "/")

	if path != "/certificates/contacts" {
		return nil, fmt.Errorf("keyVault Certificate Contacts ID path must be '/certificates/contacts', got %q", path)
	}

	id := CertificateContactsId{
		KeyVaultBaseUrl: fmt.Sprintf("%s://%s/", idURL.Scheme, idURL.Host),
	}

	return &id, nil
}
