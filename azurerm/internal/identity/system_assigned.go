package identity

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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

func (s SystemAssigned) Schema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						systemAssigned,
					}, false),
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
