package aznr001

import (
	"testdata/src/mockpkg/pluginsdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceInValid() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*schema.Schema{ // want `name, resource_group_name, location, account_replication_type, account_tier, enable_https, tags, primary_key`
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

			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"account_replication_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"account_tier": {
				Type:     schema.TypeString,
				Required: true,
			},

			"enable_https": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"primary_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceValid() *pluginsdk.Resource {
	return &pluginsdk.Resource{
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

			"account_replication_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"account_tier": {
				Type:     schema.TypeString,
				Required: true,
			},

			"enable_https": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"primary_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
