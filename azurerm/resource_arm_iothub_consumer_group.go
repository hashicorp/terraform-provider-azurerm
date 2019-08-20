package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubConsumerGroupName,
			},

			"iothub_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"eventhub_endpoint_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),
		},
	}
}

func resourceArmIotHubConsumerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothub.ResourceClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM IoTHub Consumer Group creation.")

	name := d.Get("name").(string)
	iotHubName := d.Get("iothub_name").(string)
	endpointName := d.Get("eventhub_endpoint_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Consumer Group %q (Endpoint %q / IoTHub %q / Resource Group %q): %s", name, endpointName, iotHubName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_iothub_consumer_group", *existing.ID)
		}
	}

	if _, err := client.CreateEventHubConsumerGroup(ctx, resourceGroup, iotHubName, endpointName, name); err != nil {
		return fmt.Errorf("Error creating Consumer Group %q (Endpoint %q / IoTHub %q / Resource Group %q): %+v", name, endpointName, iotHubName, resourceGroup, err)
	}

	read, err := client.GetEventHubConsumerGroup(ctx, resourceGroup, iotHubName, endpointName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Consumer Group %q (Endpoint %q / IoTHub %q / Resource Group %q): %+v", name, endpointName, iotHubName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for Consumer Group %q (Endpoint %q / IoTHub %q / Resource Group %q): %+v", name, endpointName, iotHubName, resourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceArmIotHubConsumerGroupRead(d, meta)
}

func resourceArmIotHubConsumerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothub.ResourceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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

		return fmt.Errorf("Error making read request for Consumer Group %q (Endpoint %q / IoTHub %q / Resource Group %q): %+v", name, endpointName, iotHubName, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("iothub_name", iotHubName)
	d.Set("eventhub_endpoint_name", endpointName)
	d.Set("resource_group_name", resourceGroup)

	return nil
}

func resourceArmIotHubConsumerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothub.ResourceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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
			return fmt.Errorf("Error deleting Consumer Group %q (Endpoint %q / IoTHub %q / Resource Group %q): %+v", name, endpointName, iotHubName, resourceGroup, err)
		}
	}

	return nil
}
