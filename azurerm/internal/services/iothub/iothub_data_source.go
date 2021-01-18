package iothub

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	iothubValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceIotHub() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIotHubRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: iothubValidate.IoTHubName,
			},
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"tags": tags.Schema(),
		},
	}
}

func dataSourceIotHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: IoTHub %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error retrieving IotHub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.SetId(*resp.ID)

	return tags.FlattenAndSet(d, resp.Tags)

}
