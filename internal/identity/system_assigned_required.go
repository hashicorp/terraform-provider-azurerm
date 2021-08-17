package identity

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ Identity = SystemAssignedRequired{}

type SystemAssignedRequired struct{}

func (SystemAssignedRequired) Expand(input []interface{}) (*ExpandedConfig, error) {
	return SystemAssigned{}.Expand(input)
}

func (SystemAssignedRequired) Flatten(input *ExpandedConfig) []interface{} {
	return SystemAssigned{}.Flatten(input)
}

func (SystemAssignedRequired) Schema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(systemAssigned),
					}, false),
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
