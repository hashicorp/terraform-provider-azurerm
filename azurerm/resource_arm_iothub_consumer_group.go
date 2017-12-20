package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmIotHubConsumerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIotHubConsumerGroupCreateUpdate,
		Read:   resourceArmIotHubConsumerGroupRead,
		Update: resourceArmIotHubConsumerGroupCreateUpdate,
		Delete: resourceArmIotHubConsumerGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_name": resourceGroupNameSchema(),
			"iot_hub_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"event_hub_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmIotHubConsumerGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient
	log.Printf("[INFO} preparing arguments for AzureRM IoTHub Consumer Group creation.")

	groupName := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	iotHubName := d.Get("iot_hub_name").(string)
	eventhubEndpoint := d.Get("event_hub_endpoint").(string)

	_, err := iothubClient.CreateEventHubConsumerGroup(resourceGroup, iotHubName, eventhubEndpoint, groupName)
	if err != nil {
		return err
	}

	check, err := iothubClient.GetEventHubConsumerGroup(resourceGroup, iotHubName, eventhubEndpoint, groupName)
	if err != nil {
		return err
	}

	if check.ID == nil {
		return fmt.Errorf("Cannot read IoTHub Consumer Group %s (resource group %s) ID", groupName, resourceGroup)
	}

	d.SetId(*check.ID)

	return resourceArmIotHubConsumerGroupRead(d, meta)
}

func resourceArmIotHubConsumerGroupRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	groupName := id.Path["consumergroups"]
	iotHubName := id.Path["resourcename"]
	eventhubEndpoint := id.Path["eventhubs"]

	resp, err := iothubClient.GetEventHubConsumerGroup(resourceGroup, iotHubName, eventhubEndpoint, groupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making read request on Azure IoTHub Consumer Group %s: %+v", groupName, err)
	}

	d.Set("name", groupName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("iot_hub_name", iotHubName)
	d.Set("event_hub_endpoint", eventhubEndpoint)

	return nil
}

func resourceArmIotHubConsumerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	groupName := id.Path["consumergroups"]
	iotHubName := id.Path["resourcename"]
	eventhubEndpoint := id.Path["eventhubs"]

	resp, err := iothubClient.DeleteEventHubConsumerGroup(resourceGroup, iotHubName, eventhubEndpoint, groupName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error issuing Azure ARM delete request of IoTHub Consumer Group '%s' : %+v", groupName, err)
		}
	}

	return nil
}
