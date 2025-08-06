// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
)

type legacyV0RoleAssignmentId struct {
	// @tombuildsstuff: For reasons I can't entirely ascertain it appears that we were making up a
	// Data Plane URI for this, rather than using the existing one, hence the need for this state migration.
	//
	// Example Old Value: `https://tharvey-keyvault.managedhsm.azure.net///RoleAssignment/uuid-idshifds-fks`
	//
	// Example New Value: `https://tharvey-keyvault.managedhsm.azure.net/{scope}/providers/Microsoft.Authorization/roleAssignments/{roleAssignmentName}`

	managedHSMName     string
	domainSuffix       string
	scope              string
	roleAssignmentName string
}

func parseLegacyV0RoleAssignmentId(input string) (*legacyV0RoleAssignmentId, error) {
	parsed, err := url.ParseRequestURI(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	endpoint, err := parse.ManagedHSMEndpoint(input, nil)
	if err != nil {
		return nil, fmt.Errorf("parsing endpoint from %q: %+v", input, err)
	}

	path := strings.TrimPrefix(parsed.Path, "/")
	split := strings.Split(path, "/RoleAssignment/")
	if len(split) != 2 {
		return nil, fmt.Errorf("expected a URI in the format `{scope}/RoleAssignment/{name}` but got %q", parsed.Path)
	}
	scope := split[0]
	name := split[1]
	if scope == "" || name == "" {
		return nil, fmt.Errorf("expected a URI in the format `{scope}/RoleAssignment/{name}` but got %q", parsed.Path)
	}

	return &legacyV0RoleAssignmentId{
		managedHSMName:     endpoint.ManagedHSMName,
		domainSuffix:       endpoint.DomainSuffix,
		scope:              scope,
		roleAssignmentName: name,
	}, nil
}
