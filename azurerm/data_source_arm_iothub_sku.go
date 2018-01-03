package azurerm

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmIotHubSku() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmIotHubSkuRead,

		Schema: map[string]*schema.Schema{
			"resource_group_name": resourceGroupNameForDataSourceSchema(),
			"iot_hub_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceArmIotHubSkuRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient
	log.Printf("[INFO] Acquiring IoTHub SKU")

	iothubName := d.Get("iot_hub_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	skuResp, err := iothubClient.GetValidSkus(resourceGroup, iothubName)
	if err != nil {
		return err
	}

	var skus []map[string]interface{}

	for _, sku := range *skuResp.Value {
		skuMap := make(map[string]interface{})
		skuMap["resource_type"] = *sku.ResourceType
		skuMap["capacity"] = *sku.Capacity
		skuMap["sku"] = *sku.Sku // SB todo: type SkuInfo Breakdown for map
		skus = append(skus, skuMap)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("iot_hub_name", iothubName)
	d.Set("sku", skus)

	return nil
}
