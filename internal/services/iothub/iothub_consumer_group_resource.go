package iothub

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIotHubConsumerGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubConsumerGroupCreate,
		Read:   resourceIotHubConsumerGroupRead,
		Delete: resourceIotHubConsumerGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ConsumerGroupID(id)
			return err
		}),

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
	subscriptionId := meta.(*clients.Client).IoTHub.DPSResourceClient.SubscriptionID
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM IoTHub Consumer Group creation.")

	id := parse.NewConsumerGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_name").(string), d.Get("eventhub_endpoint_name").(string), d.Get("name").(string))

	locks.ByName(id.IotHubName, IothubResourceName)
	defer locks.UnlockByName(id.IotHubName, IothubResourceName)

	if d.IsNewResource() {
		existing, err := client.GetEventHubConsumerGroup(ctx, id.ResourceGroup, id.IotHubName, id.EventHubEndpointName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Consumer Group %s: %+v", id.String(), err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_iothub_consumer_group", *existing.ID)
		}
	}

	if _, err := client.CreateEventHubConsumerGroup(ctx, id.ResourceGroup, id.IotHubName, id.EventHubEndpointName, id.Name); err != nil {
		return fmt.Errorf("creating %s: %+v", id.String(), err)
	}

	read, err := client.GetEventHubConsumerGroup(ctx, id.ResourceGroup, id.IotHubName, id.EventHubEndpointName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id.String(), err)
	}

	if read.ID == nil {
		return fmt.Errorf("cannot read %s: %+v", id.String(), err)
	}

	d.SetId(id.ID())

	return resourceIotHubConsumerGroupRead(d, meta)
}

func resourceIotHubConsumerGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConsumerGroupID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	iotHubName := id.IotHubName
	endpointName := id.EventHubEndpointName
	name := id.Name

	resp, err := client.GetEventHubConsumerGroup(ctx, resourceGroup, iotHubName, endpointName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request for %s: %+v", id.String(), err)
	}

	d.Set("name", id.Name)
	d.Set("iothub_name", id.IotHubName)
	d.Set("eventhub_endpoint_name", id.EventHubEndpointName)
	d.Set("resource_group_name", id.ResourceGroup)

	return nil
}

func resourceIotHubConsumerGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConsumerGroupID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	iotHubName := id.IotHubName
	endpointName := id.EventHubEndpointName
	name := id.Name

	locks.ByName(iotHubName, IothubResourceName)
	defer locks.UnlockByName(iotHubName, IothubResourceName)

	resp, err := client.DeleteEventHubConsumerGroup(ctx, resourceGroup, iotHubName, endpointName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", id.String(), err)
		}
	}

	return nil
}
