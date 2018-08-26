package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmCdnProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmCdnProfileRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

			"sku": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmCdnProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cdnProfilesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: CDN Profile %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on Azure CDN Profile %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}
