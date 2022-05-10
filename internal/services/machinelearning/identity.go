package machinelearning

import (
	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2021-07-01/machinelearningservices"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func expandIdentity(input []interface{}) (*machinelearningservices.Identity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := machinelearningservices.Identity{
		Type: machinelearningservices.ResourceIdentityType(string(expanded.Type)),
	}

	// work around the Swagger defining `SystemAssigned,UserAssigned` rather than `SystemAssigned, UserAssigned`
	if expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.Type = machinelearningservices.ResourceIdentityTypeSystemAssignedUserAssigned
	}

	// 'Failed to perform resource identity operation. Status: 'BadRequest'. Response:
	// {"error":{"code":"BadRequest",
	//  "message":"The request format was unexpected, a non-UserAssigned identity type should not contain: userAssignedIdentities"
	// }}
	// Upstream issue: https://github.com/Azure/azure-rest-api-specs/issues/17650
	if len(expanded.IdentityIds) > 0 {
		userAssignedIdentities := make(map[string]*machinelearningservices.UserAssignedIdentity)
		for id := range expanded.IdentityIds {
			userAssignedIdentities[id] = &machinelearningservices.UserAssignedIdentity{}
		}
		out.UserAssignedIdentities = userAssignedIdentities
	}

	return &out, nil
}

func flattenIdentity(input *machinelearningservices.Identity) (*[]interface{}, error) {
	var config *identity.SystemAndUserAssignedMap

	if input != nil {
		config = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: nil,
		}

		// work around the Swagger defining `SystemAssigned,UserAssigned` rather than `SystemAssigned, UserAssigned`
		if input.Type == machinelearningservices.ResourceIdentityTypeSystemAssignedUserAssigned {
			config.Type = identity.TypeSystemAssignedUserAssigned
		}

		if input.PrincipalID != nil {
			config.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			config.TenantId = *input.TenantID
		}
		identityIds := make(map[string]identity.UserAssignedIdentityDetails)
		for k, v := range input.UserAssignedIdentities {
			if v == nil {
				continue
			}

			details := identity.UserAssignedIdentityDetails{}

			if v.ClientID != nil {
				details.ClientId = utils.String(*v.ClientID)
			}
			if v.PrincipalID != nil {
				details.PrincipalId = utils.String(*v.PrincipalID)
			}

			identityIds[k] = details
		}

		config.IdentityIds = identityIds
	}

	return identity.FlattenSystemAndUserAssignedMap(config)
}
