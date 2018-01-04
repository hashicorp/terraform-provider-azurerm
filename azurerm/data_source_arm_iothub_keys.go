package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmIotHubKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmIotHubKeysRead,

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

	key, err := iothubClient.GetKeysForKeyName(resourceGroup, iotHub, name)
	if err != nil {
		return fmt.Errorf("Error retrieving keys with name: %s", name)
	}

	var keys []map[string]interface{}
	keyMap := make(map[string]interface{})
	keyMap["key_name"] = *key.KeyName
	keyMap["primary_key"] = *key.PrimaryKey
	keyMap["secondary_key"] = *key.SecondaryKey
	keyMap["permissions"] = string(key.Rights)
	keys = append(keys, keyMap)

	d.SetId(name)
	d.Set("shared_access_policy", keys)

	return nil
}
