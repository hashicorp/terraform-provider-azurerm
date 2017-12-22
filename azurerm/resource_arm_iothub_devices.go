package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmIotHubDevices() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Read:   resourceArmIotHubDevicesRead,
		Update: resourceArmIotHubDevicesUpdate,
		Delete: resourceArmIotHubDevicesDelete,

		Schema: map[string]*schema.Schema{
			"device_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_name": resourceGroupNameSchema(),
			"iotHub_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"import_device_parameters": {
				Type:     schema.TypeString,
				Required: true,
			},
			"import_device_request": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmIotHubDevicesRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient

	return nil
}

func resourceArmIotHubDevicesUpdate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient

	return resourceArmIotHubDevicesRead(d, meta)
}

func resourceArmIotHubDevicesDelete(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient

	return nil
}
