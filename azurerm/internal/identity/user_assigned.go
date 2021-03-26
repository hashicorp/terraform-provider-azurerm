package identity

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var _ Identity = UserAssigned{}

type UserAssigned struct{}

func (u UserAssigned) Expand(input []interface{}) (*ExpandedConfig, error) {
	if len(input) == 0 || input[0] == nil {
		return &ExpandedConfig{
			Type: none,
		}, nil
	}

	return &ExpandedConfig{
		Type: systemAssigned,
	}, nil
}

func (u UserAssigned) Flatten(input *ExpandedConfig) []interface{} {
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

func (u UserAssigned) Schema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						userAssigned,
					}, false),
				},
				"identity_ids": {
					Type:     schema.TypeList,
					Required: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.NoZeroValues,
					},
				},
			},
		},
	}
}

func (u UserAssigned) SchemaDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"identity_ids": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}
