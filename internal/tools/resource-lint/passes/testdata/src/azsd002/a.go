package azsd002

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validCases() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Valid: Has a required field, so AtLeastOneOf not needed
		"config_with_required": {
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

		// Valid: Has AtLeastOneOf set on optional fields
		"config_with_atleastone": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"linux": {
						Type:         schema.TypeList,
						Optional:     true,
						AtLeastOneOf: []string{"config_with_atleastone.0.linux", "config_with_atleastone.0.windows"},
					},
					"windows": {
						Type:         schema.TypeList,
						Optional:     true,
						AtLeastOneOf: []string{"config_with_atleastone.0.linux", "config_with_atleastone.0.windows"},
					},
				},
			},
		},

		// Valid: Only one optional field, no need for AtLeastOneOf
		"config_single_optional": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"value": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},

		// Valid: Not a TypeList
		"simple_field": {
			Type:     schema.TypeString,
			Optional: true,
		},

		// Valid: Computed field, should be skipped
		"computed_field": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"property1": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"property2": {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},
	}
}

func invalidCases() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Invalid: Multiple optional fields without AtLeastOneOf
		"setting": { // want `AZSD002`
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"linux": {
						Type:     schema.TypeList,
						Optional: true,
					},
					"windows": {
						Type:     schema.TypeList,
						Optional: true,
					},
				},
			},
		},

		// Invalid: Three optional fields without AtLeastOneOf
		"platform": { // want `AZSD002`
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"linux": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"windows": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"macos": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}

func invalidStandaloneCases() *schema.Schema {
	return &schema.Schema{ // want `AZSD002`
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"linux": {
					Type:     schema.TypeList,
					Optional: true,
				},
				"windows": {
					Type:     schema.TypeList,
					Optional: true,
				},
			},
		},
	}
}
