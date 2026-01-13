package azbp001

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func validCases() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Valid: Required String with ValidateFunc
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// Valid: Optional String with ValidateFunc
		"description": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 100),
		},

		// Valid: Computed-only String (no ValidateFunc needed)
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		// Valid: Non-String type (no ValidateFunc needed)
		"count": {
			Type:     schema.TypeInt,
			Required: true,
		},
	}
}

func invalidCases() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Invalid: Required String without ValidateFunc
		"resource_group_name": { // want `AZBP001`
			Type:     schema.TypeString,
			Required: true,
		},

		// Invalid: Optional String without ValidateFunc
		"location": { // want `AZBP001`
			Type:     schema.TypeString,
			Optional: true,
		},

		// Invalid: Required String with other fields but no ValidateFunc
		"sku": { // want `AZBP001`
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func invalidStandaloneCases() *schema.Schema {
	return &schema.Schema{ // want `AZBP001`
		Type:     schema.TypeString,
		Required: true,
	}
}
