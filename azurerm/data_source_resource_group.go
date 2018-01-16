package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmResourceGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmResourceGroupRead,

		Schema: map[string]*schema.Schema{
			"name":     resourceGroupNameForDataSourceSchema(),
			"location": locationForDataSourceSchema(),
			"tags":     tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).resourceGroupsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resp, err := client.Get(ctx, name)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	if err := resourceArmResourceGroupRead(d, meta); err != nil {
		return err
	}

	return nil
}
