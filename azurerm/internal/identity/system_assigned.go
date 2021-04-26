package identity

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

var _ Identity = SystemAssigned{}

type SystemAssigned struct{}

func (s SystemAssigned) Expand(input []interface{}) (*ExpandedConfig, error) {
	if len(input) == 0 || input[0] == nil {
		return &ExpandedConfig{
			Type: none,
		}, nil
	}

	return &ExpandedConfig{
		Type: systemAssigned,
	}, nil
}

func (s SystemAssigned) Flatten(input *ExpandedConfig) []interface{} {
	if input == nil || input.Type == none {
		return []interface{}{}
	}

	var coalesce = func(input *string) string {
		if input == nil {
			return ""
		}

		return *input
	}

	return []interface{}{
		map[string]interface{}{
			"type":         input.Type,
			"principal_id": coalesce(input.PrincipalId),
			"tenant_id":    coalesce(input.TenantId),
		},
	}
}

func (s SystemAssigned) Schema() *pluginsdk.Schema {
	//lintignore:XS003
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						systemAssigned,
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
