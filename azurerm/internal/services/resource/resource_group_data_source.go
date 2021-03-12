package resource

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceResourceGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceResourceGroupRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name":     azure.SchemaResourceGroupNameForDataSource(),
			"location": azure.SchemaLocationForDataSource(),
			"tags":     tags.SchemaDataSource(),
		},
	}
}

func dataSourceResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.GroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Resource Group %q was not found", name)
		}
		return err
	}

	d.SetId(*resp.ID)

	return resourceResourceGroupRead(d, meta)
}
