package iothub

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceIotHubConsumerGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubConsumerGroupCreate,
		Read:   resourceIotHubConsumerGroupRead,
		Delete: resourceIotHubConsumerGroupDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubConsumerGroupName,
			},

			"iothub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"eventhub_endpoint_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),
		},
	}
}

func resourceIotHubConsumerGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM IoTHub Consumer Group creation.")

	name := d.Get("name").(string)
	iotHubName := d.Get("iothub_name").(string)
	endpointName := d.Get("eventhub_endpoint_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(iotHubName, IothubResourceName)
	defer locks.UnlockByName(iotHubName, IothubResourceName)

	if d.IsNewResource() {
		existing, err := client.GetEventHubConsumerGroup(ctx, resourceGroup, iotHubName, endpointName, name)
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

	return resourceIotHubConsumerGroupRead(d, meta)
}

func resourceIotHubConsumerGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

func resourceIotHubConsumerGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	iotHubName := id.Path["IotHubs"]
	endpointName := id.Path["eventHubEndpoints"]
	name := id.Path["ConsumerGroups"]

	locks.ByName(iotHubName, IothubResourceName)
	defer locks.UnlockByName(iotHubName, IothubResourceName)

	resp, err := client.DeleteEventHubConsumerGroup(ctx, resourceGroup, iotHubName, endpointName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Consumer Group %q (Endpoint %q / IoTHub %q / Resource Group %q): %+v", name, endpointName, iotHubName, resourceGroup, err)
		}
	}

	return nil
}
