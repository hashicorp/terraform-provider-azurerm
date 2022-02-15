package machinelearning

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2021-07-01/machinelearningservices"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func identityLegacySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(identity.TypeUserAssigned),
						string(identity.TypeSystemAssigned),
						string(identity.TypeSystemAssignedUserAssigned),
						"SystemAssigned,UserAssigned", // defined in the Swagger but should be normalized as above
					}, false),
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						// handle `SystemAssigned, UserAssigned` with and without the spaces being the same
						oldWithoutSpaces := strings.ReplaceAll(old, " ", "")
						newWithoutSpaces := strings.ReplaceAll(new, " ", "")
						return oldWithoutSpaces == newWithoutSpaces
					},
				},
				"identity_ids": {
					Type:     schema.TypeSet,
					Optional: true,
					ForceNew: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
				"principal_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"tenant_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func expandIdentity(input []interface{}) (*machinelearningservices.Identity, error) {
	if !features.ThreePointOhBeta() {
		// work around the Swagger defining `SystemAssigned,UserAssigned` rather than `SystemAssigned, UserAssigned`
		if len(input) > 0 {
			raw := input[0].(map[string]interface{})
			if identityType := raw["type"].(string); strings.EqualFold("SystemAssigned,UserAssigned", identityType) {
				raw["type"] = "SystemAssigned, UserAssigned"
			}
			input[0] = raw
		}
	}

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
