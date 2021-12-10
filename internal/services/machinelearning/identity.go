package machinelearning

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2021-07-01/machinelearningservices"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	msiParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/parse"
	msiValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SystemAssignedUserAssigned struct{}

func (s SystemAssignedUserAssigned) Expand(input []interface{}) (*machinelearningservices.Identity, error) {
	if len(input) == 0 || input[0] == nil {
		return &machinelearningservices.Identity{
			Type: machinelearningservices.ResourceIdentityTypeNone,
		}, nil
	}

	v := input[0].(map[string]interface{})

	config := &machinelearningservices.Identity{
		Type: machinelearningservices.ResourceIdentityType(v["type"].(string)),
	}

	identityIds := v["identity_ids"].(*pluginsdk.Set).List()
	if len(identityIds) != 0 {
		if config.Type != machinelearningservices.ResourceIdentityTypeUserAssigned && config.Type != machinelearningservices.ResourceIdentityTypeSystemAssignedUserAssigned {
			return nil, fmt.Errorf("`identity_ids` can only be specified when `type` includes `UserAssigned`")
		}
		config.UserAssignedIdentities = map[string]*machinelearningservices.UserAssignedIdentity{}
		for _, id := range identityIds {
			config.UserAssignedIdentities[id.(string)] = &machinelearningservices.UserAssignedIdentity{}
		}
	}

	return config, nil
}

func (s SystemAssignedUserAssigned) Flatten(input *machinelearningservices.Identity) ([]interface{}, error) {
	if input == nil || input.Type == machinelearningservices.ResourceIdentityTypeNone {
		return []interface{}{}, nil
	}

	coalesce := func(input *string) string {
		if input == nil {
			return ""
		}

		return *input
	}

	var identityIds []string
	for id := range input.UserAssignedIdentities {
		parsedId, err := msiParse.UserAssignedIdentityIDInsensitively(id)
		if err != nil {
			return nil, err
		}
		identityIds = append(identityIds, parsedId.ID())
	}

	return []interface{}{
		map[string]interface{}{
			"type":         input.Type,
			"identity_ids": identityIds,
			"principal_id": coalesce(input.PrincipalID),
			"tenant_id":    coalesce(input.TenantID),
		},
	}, nil
}

func (s SystemAssignedUserAssigned) Schema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(machinelearningservices.ResourceIdentityTypeUserAssigned),
						string(machinelearningservices.ResourceIdentityTypeSystemAssigned),
						string(machinelearningservices.ResourceIdentityTypeSystemAssignedUserAssigned),
					}, false),
				},
				"identity_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					ForceNew: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: msiValidate.UserAssignedIdentityID,
					},
				},
				"principal_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"tenant_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func (s SystemAssignedUserAssigned) SchemaDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"identity_ids": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
				"principal_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"tenant_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}
