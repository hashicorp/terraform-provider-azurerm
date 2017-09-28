package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmIothub() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIothubCreate,
		Read:   resourceArmIothubRead,
		Update: resourceArmIothubUpdate,
		Delete: resourceArmIothubDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type: schema.TypeString,
			},

			"name": {
				Type: schema.TypeString,
			},

			"type": {
				Type: schema.TypeString,
			},

			"location": {
				Type: schema.TypeString,
			},

			"subscriptionid": {
				Type: schema.TypeString,
			},

			"resourcegroup": {
				Type: schema.TypeString,
			},

			"etag": {
				Type: schema.TypeString,
			},

			"skuinfo": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type: schema.TypeString,
						},

						"tier": {
							Type: schema.TypeString,
						},

						"capacity": {
							Type: schema.TypeInt,
						},
					},
				},
			},
		},
	}

}
func resourceArmIothubCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	iotHubCreateClient := client.iothubResourceClient
	resGroup := d.Get("resourceGroupName").(string)
	resName := d.Get("resourceName").(string)

	return nil

}
