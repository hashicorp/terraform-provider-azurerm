package identity

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	msivalidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var _ Identity = UserAssigned{}

type UserAssigned struct{}

func (u UserAssigned) Expand(input []interface{}) (*ExpandedConfig, error) {
	if len(input) == 0 || input[0] == nil {
		return &ExpandedConfig{
			Type: none,
		}, nil
	}

	v := input[0].(map[string]interface{})

	return &ExpandedConfig{
		Type:                    userAssigned,
		UserAssignedIdentityIds: utils.ExpandStringSlice(v["identity_ids"].(*pluginsdk.Set).List()),
	}, nil
}

func (u UserAssigned) Flatten(input *ExpandedConfig) []interface{} {
	if input == nil || input.Type == none {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"type":         input.Type,
			"identity_ids": utils.FlattenStringSlice(input.UserAssignedIdentityIds),
		},
	}
}

func (u UserAssigned) Schema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						userAssigned,
					}, false),
				},
				"identity_ids": {
					Type:     pluginsdk.TypeSet,
					Required: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: msivalidate.UserAssignedIdentityID,
					},
				},
			},
		},
	}
}

func (u UserAssigned) SchemaDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
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
			},
		},
	}
}
