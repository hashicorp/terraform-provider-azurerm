package aznr001

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Test: Nested schemas should skip ID field ordering
func resourceWithNestedSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
				ForceNew: true,
			},

			// Nested schema should NOT require ID fields first
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Alphabetical order for nested - no ID field requirement
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

// Test: Nested schema with wrong order should be flagged
func resourceWithBadNestedSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{ // want `enabled, name, value`
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
		},
	}
}

// Helper function from same package but different file scope
func GetPercentageSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
	}
}
