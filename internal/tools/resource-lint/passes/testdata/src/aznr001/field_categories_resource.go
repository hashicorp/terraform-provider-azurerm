package aznr001

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Test: Proper categorization - Required, Optional, Computed
func resourceFieldCategories() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// ID fields first
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

			// Location
			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// Required fields (alphabetical)
			"account_tier": {
				Type:     schema.TypeString,
				Required: true,
			},

			"sku": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Optional fields (alphabetical)
			"enable_https": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			// Computed fields (alphabetical)
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Test: Wrong category order
func resourceWrongCategoryOrder() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{ // want `name, resource_group_name, location, account_tier, sku, enable_https, tags, created_time, primary_key`
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// Computed field too early
			"primary_key": {
				Type:     schema.TypeString,
				Computed: true,
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

			// Optional before required
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"account_tier": {
				Type:     schema.TypeString,
				Required: true,
			},

			"sku": {
				Type:     schema.TypeString,
				Required: true,
			},

			"enable_https": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
