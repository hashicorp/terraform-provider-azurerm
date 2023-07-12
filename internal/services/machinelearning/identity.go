// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func expandIdentity(input []interface{}) (*identity.LegacySystemAndUserAssignedMap, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := identity.LegacySystemAndUserAssignedMap{
		Type: expanded.Type,
	}

	// 'Failed to perform resource identity operation. Status: 'BadRequest'. Response:
	// {"error":{"code":"BadRequest",
	//  "message":"The request format was unexpected, a non-UserAssigned identity type should not contain: userAssignedIdentities"
	// }}
	// Upstream issue: https://github.com/Azure/azure-rest-api-specs/issues/17650
	if len(expanded.IdentityIds) > 0 {
		userAssignedIdentities := make(map[string]identity.UserAssignedIdentityDetails)
		for id := range expanded.IdentityIds {
			userAssignedIdentities[id] = identity.UserAssignedIdentityDetails{}
		}
		out.IdentityIds = userAssignedIdentities
	}

	return &out, nil
}

func flattenIdentity(input *identity.LegacySystemAndUserAssignedMap) (*[]interface{}, error) {
	var config *identity.SystemAndUserAssignedMap

	if input != nil {
		config = &identity.SystemAndUserAssignedMap{
			Type:        input.Type,
			IdentityIds: nil,
		}

		if input.PrincipalId != "" {
			config.PrincipalId = input.PrincipalId
		}
		if input.TenantId != "nil" {
			config.TenantId = input.TenantId
		}
		identityIds := make(map[string]identity.UserAssignedIdentityDetails)
		if input.IdentityIds != nil {
			for k, v := range input.IdentityIds {
				details := identity.UserAssignedIdentityDetails{}

				if v.ClientId != nil {
					details.ClientId = utils.String(*v.ClientId)
				}
				if v.PrincipalId != nil {
					details.PrincipalId = utils.String(*v.PrincipalId)
				}

				identityIds[k] = details
			}
		}

		config.IdentityIds = identityIds
	}

	return identity.FlattenSystemAndUserAssignedMap(config)
}
