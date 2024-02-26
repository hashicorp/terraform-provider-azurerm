// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = RoleNestedItemId{}

type MHSMResourceType string

const (
	// TODO: this should be extended to support the other types of Nested Items available for a Managed HSM

	RoleDefinitionType MHSMResourceType = "RoleDefinition"
	RoleAssignmentType MHSMResourceType = "RoleAssignment"
)

type RoleNestedItemId struct {
	VaultBaseUrl string
	Scope        string
	Type         MHSMResourceType
	Name         string
}

func NewRoleNestedItemID(hsmBaseUrl, scope string, typ MHSMResourceType, name string) (*RoleNestedItemId, error) {
	keyVaultUrl, err := url.Parse(hsmBaseUrl)
	if err != nil || hsmBaseUrl == "" {
		return nil, fmt.Errorf("parsing managedHSM nested itemID %q: %+v", hsmBaseUrl, err)
	}
	// (@jackofallops) - Log Analytics service adds the port number to the API returns, so we strip it here
	if hostParts := strings.Split(keyVaultUrl.Host, ":"); len(hostParts) > 1 {
		keyVaultUrl.Host = hostParts[0]
	}

	return &RoleNestedItemId{
		VaultBaseUrl: keyVaultUrl.String(),
		Scope:        scope,
		Type:         typ,
		Name:         name,
	}, nil
}

func (n RoleNestedItemId) ID() string {
	// example: https://tharvey-keyvault.managedhsm.azure.net///RoleDefinition/uuid-idshifds-fks
	segments := []string{
		strings.TrimSuffix(n.VaultBaseUrl, "/"),
		n.Scope,
		string(n.Type),
		n.Name,
	}
	return strings.TrimSuffix(strings.Join(segments, "/"), "/")
}

func (n RoleNestedItemId) String() string {
	return n.ID()
}

func RoleNestedItemID(input string) (*RoleNestedItemId, error) {
	idURL, err := url.ParseRequestURI(input)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Azure KeyVault Child Id: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	nameSep := strings.LastIndex(path, "/")
	if nameSep <= 0 {
		return nil, fmt.Errorf("no name speparate exist in %s", input)
	}
	scope, name := path[:nameSep], path[nameSep+1:]

	typeSep := strings.LastIndex(scope, "/")
	if typeSep <= 0 {
		return nil, fmt.Errorf("no type speparate exist in %s", input)
	}
	scope, typ := scope[:typeSep], scope[typeSep+1:]

	childId := RoleNestedItemId{
		VaultBaseUrl: fmt.Sprintf("%s://%s/", idURL.Scheme, idURL.Host),
		Scope:        scope,
		Type:         MHSMResourceType(typ),
		Name:         name,
	}

	return &childId, nil
}
