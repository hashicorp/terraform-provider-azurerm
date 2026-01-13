package azsd001

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validCases() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Valid: MaxItems:1 with multiple properties (no need to flatten)
		"config_multi": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"property1": {
						Type:     schema.TypeString,
						Required: true,
					},
					"property2": {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},

		// Valid: MaxItems:1 with single property but has explanatory comment
		"config_with_comment": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{ // Additional properties will be added per service team confirmation
				Schema: map[string]*schema.Schema{
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},

		// Valid: Not MaxItems:1 (so no rule applies)
		"config_list": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},

		// Valid: Simple string field (not a block)
		"simple_field": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func invalidCases() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Invalid: MaxItems:1 with only one property and no comment
		"config_single": { // want `AZSD001`
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},

		// Invalid: Another case of MaxItems:1 with single property, no justification
		"settings": { // want `AZSD001`
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
	}
}

func invalidStandaloneCases() *schema.Schema {
	return &schema.Schema{ // want `AZSD001`
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}
