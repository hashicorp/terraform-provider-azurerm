package aznr001

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     schema.TypeString,
			Required: true,
		},

		"sku_name": {
			Type:     schema.TypeString,
			Required: true,
		},

		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
		},

		"tags": {
			Type:     schema.TypeMap,
			Optional: true,
		},
	}
}

// Test: Typed resource with wrong order
func argumentsWrong() map[string]*schema.Schema {
	return map[string]*schema.Schema{ // want `name, resource_group_name, location, sku_name, enabled, tags`
		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"tags": {
			Type:     schema.TypeMap,
			Optional: true,
		},

		"location": {
			Type:     schema.TypeString,
			Required: true,
		},

		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
		},

		"sku_name": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

// Test: Using schema from helper function in same package
func argumentsWithImportedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{ // want `name, enabled, not_inline`
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},

		"not_inline": GetPercentageSchema(),

		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}
