package databoxedge

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceDataboxEdgeOrder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDataboxEdgeOrderRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
		},
	}
}

func dataSourceDataboxEdgeOrderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataboxEdge.OrderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, name, resourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Databox Edge Order does not exist")
		}
		return fmt.Errorf("retrieving Databox Edge Order (Resource Group %q / Name %q): %+v", resourceGroup, name, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Databox Edge Order (Resource Group %q / Name %q) ID", resourceGroup, name)
	}

	d.SetId(*resp.ID)
	d.Set("resource_group_name", resourceGroup)
	d.Set("name", name)
	return nil
}
