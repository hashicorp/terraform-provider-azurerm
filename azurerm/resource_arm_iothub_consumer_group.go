package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmIotHubConsumerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIotHubConsumerGroupCreate,
		Read:   resourceArmIotHubConsumerGroupRead,
		Delete: resourceArmIotHubConsumerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// TODO: validation
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"iothub_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"eventhub_endpoint_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),
		},
	}
}

func resourceArmIotHubConsumerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothubResourceClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM IoTHub Consumer Group creation.")

	name := d.Get("name").(string)
	iotHubName := d.Get("iothub_name").(string)
	endpointName := d.Get("eventhub_endpoint_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	_, err := client.CreateEventHubConsumerGroup(ctx, resourceGroup, iotHubName, endpointName, name)
	if err != nil {
		return err
	}

	read, err := client.GetEventHubConsumerGroup(ctx, resourceGroup, iotHubName, endpointName, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read IoTHub Consumer Group %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmIotHubConsumerGroupRead(d, meta)
}

func resourceArmIotHubConsumerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothubResourceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	iotHubName := id.Path["IotHubs"]
	endpointName := id.Path["eventHubEndpoints"]
	name := id.Path["ConsumerGroups"]

	resp, err := client.GetEventHubConsumerGroup(ctx, resourceGroup, iotHubName, endpointName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on IoTHub Consumer Group %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("iothub_name", iotHubName)
	d.Set("eventhub_endpoint_name", endpointName)
	d.Set("resource_group_name", resourceGroup)

	return nil
}

func resourceArmIotHubConsumerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothubResourceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	iotHubName := id.Path["IotHubs"]
	endpointName := id.Path["eventHubEndpoints"]
	name := id.Path["ConsumerGroups"]

	resp, err := client.DeleteEventHubConsumerGroup(ctx, resourceGroup, iotHubName, endpointName, name)

	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error issuing delete request for IoTHub Consumer Group %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
