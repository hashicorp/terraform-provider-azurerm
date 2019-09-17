package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPublicIpPrefix() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPublicIpPrefixRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"sku": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"prefix_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"ip_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zones": azure.SchemaSingleZone(),

			"tags": tags.Schema(),
		},
	}
}

func dataSourceArmPublicIpPrefixRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PublicIPPrefixesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Public IP prefix %q was not found", name)
		}
		return err
	}

	d.SetId(*resp.ID)

	return resourceArmPublicIpPrefixRead(d, meta)
}
