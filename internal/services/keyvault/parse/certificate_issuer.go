// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"
)

type IssuerId struct {
	KeyVaultBaseUrl string
	Name            string
}

func IssuerID(id string) (*IssuerId, error) {
	// example: https://example-keyvault.vault.azure.net/certificates/issuers/ExampleIssuer
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("cannot parse key vault certificate issuer ID: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	if len(components) != 3 {
		return nil, fmt.Errorf("key vault certificate issuer ID should have 3 segments, got %d: '%s'", len(components), path)
	}
	if components[0] != "certificates" || components[1] != "issuers" {
		return nil, fmt.Errorf("key vault certificate issuer ID path must begin with %q", "/certificates/issuers")
	}

	issuerId := IssuerId{
		KeyVaultBaseUrl: fmt.Sprintf("%s://%s/", idURL.Scheme, idURL.Host),
		Name:            components[2],
	}

	return &issuerId, nil
}
