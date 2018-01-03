package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmIotHubKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmIotHubSkuRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_name": resourceGroupNameForDataSourceSchema(),
			"iot_hub_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceArmIotHubKeysRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient
	log.Printf("[INFO] Acquiring Keys for Authentication")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	iotHub := d.Get("iot_hub_name").(string)

	keys, err := iothubClient.GetKeysForKeyName(resourceGroup, iotHub, name)
	if err != nil {
		return fmt.Errorf("Error retrieving keys with name: %s", name)
	}

	d.Set("key_name", keys.KeyName)
	d.Set("primary_key", keys.PrimaryKey)
	d.Set("secondary_key", keys.SecondaryKey)
	d.Set("permissions", string(keys.Rights))

	return nil
}
