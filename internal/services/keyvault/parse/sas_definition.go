// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"
)

type SasDefinitionId struct {
	KeyVaultBaseUrl    string
	StorageAccountName string
	Name               string
}

func SasDefinitionID(id string) (*SasDefinitionId, error) {
	// example: https://example-keyvault.vault.azure.net/storage/exampleStorageAcc01/sas/exampleSasDefinition01
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("cannot parse Azure KeyVault Managed Storage Sas Definition Id: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	if len(components) != 4 {
		return nil, fmt.Errorf("key vault managed storage Sas Definition ID should have 4 segments, got %d: '%s'", len(components), path)
	}

	if components[0] != "storage" || components[2] != "sas" {
		return nil, fmt.Errorf("key vault managed storage Sas Definition ID path must begin with %q", "/storage/exampleStorageAcc01/sas")
	}

	sasDefinitionId := SasDefinitionId{
		KeyVaultBaseUrl:    fmt.Sprintf("%s://%s/", idURL.Scheme, idURL.Host),
		StorageAccountName: components[1],
		Name:               components[3],
	}

	return &sasDefinitionId, nil
}
