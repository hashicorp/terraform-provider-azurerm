package aznr001

import (
	"testdata/src/mockpkg/pluginsdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Test: Data source with correct order
func unTypedDataSourceValid() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"account_tier": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

// Test: Data source with wrong order
func unTypedDataSourceInvalid() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*schema.Schema{ // want `name, resource_group_name, account_tier, location, tags`
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"account_tier": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
