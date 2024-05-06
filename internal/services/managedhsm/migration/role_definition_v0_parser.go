// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
)

type legacyV0RoleDefinitionId struct {
	// @tombuildsstuff: For reasons I can't entirely ascertain it appears that we were making up a
	// Data Plane URI for this, rather than using the existing one, hence the need for this state migration.
	//
	// Example Old Value: `https://tharvey-keyvault.managedhsm.azure.net///RoleDefinition/uuid-idshifds-fks`
	//
	// Example New Value: `https://tharvey-keyvault.managedhsm.azure.net/{scope}/providers/Microsoft.Authorization/roleDefinitions/{roleDefinitionName}`
	//
	// NOTE: that for Role Definition IDs at this time the only supported value for scope is `/` - however
	// the Resource ID parser will handle any scope and the Data Source/Resource in question should limit as required.

	managedHSMName     string
	domainSuffix       string
	scope              string
	roleDefinitionName string
}

func parseLegacyV0RoleDefinitionId(input string) (*legacyV0RoleDefinitionId, error) {
	parsed, err := url.ParseRequestURI(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	endpoint, err := parse.ManagedHSMEndpoint(input, nil)
	if err != nil {
		return nil, fmt.Errorf("parsing endpoint from %q: %+v", input, err)
	}

	path := strings.TrimPrefix(parsed.Path, "/")
	split := strings.Split(path, "/RoleDefinition/")
	if len(split) != 2 {
		return nil, fmt.Errorf("expected a URI in the format `{scope}/RoleDefinition/{name}` but got %q", parsed.Path)
	}
	scope := split[0]
	name := split[1]
	if scope == "" || name == "" {
		return nil, fmt.Errorf("expected a URI in the format `{scope}/RoleDefinition/{name}` but got %q", parsed.Path)
	}

	return &legacyV0RoleDefinitionId{
		managedHSMName:     endpoint.ManagedHSMName,
		domainSuffix:       endpoint.DomainSuffix,
		scope:              scope,
		roleDefinitionName: name,
	}, nil
}
