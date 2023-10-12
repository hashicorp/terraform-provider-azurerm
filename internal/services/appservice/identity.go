// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
)

func expandIdentity(input []interface{}) (*web.ManagedServiceIdentity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := web.ManagedServiceIdentity{
		Type: web.ManagedServiceIdentityType(string(expanded.Type)),
	}
	if expanded.Type == identity.TypeUserAssigned || expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.UserAssignedIdentities = make(map[string]*web.UserAssignedIdentity)
		for k := range expanded.IdentityIds {
			out.UserAssignedIdentities[k] = &web.UserAssignedIdentity{
				// intentionally empty
			}
		}
	}
	return &out, nil
}

func flattenIdentity(input *web.ManagedServiceIdentity) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
		for k, v := range input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}
