package azurerm

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmIotHubSku() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmIotHubSkuRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_name": resourceGroupNameForDataSourceSchema(),
		},
	}
}

func dataSourceArmIotHubSkuRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient

	log.Printf("[INFO] Acquiring IoTHub SKU")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	return nil
}
