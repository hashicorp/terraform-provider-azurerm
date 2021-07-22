package digitaltwins

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/digitaltwins/mgmt/2020-10-31/digitaltwins"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/digitaltwins/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/digitaltwins/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDigitalTwinsEndpointServiceBus() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDigitalTwinsEndpointServiceBusCreateUpdate,
		Read:   resourceDigitalTwinsEndpointServiceBusRead,
		Update: resourceDigitalTwinsEndpointServiceBusCreateUpdate,
		Delete: resourceDigitalTwinsEndpointServiceBusDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DigitalTwinsEndpointID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DigitalTwinsInstanceName,
			},

			"digital_twins_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DigitalTwinsInstanceID,
			},

			"servicebus_primary_connection_string": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"servicebus_secondary_connection_string": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"dead_letter_storage_secret": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}
func resourceDigitalTwinsEndpointServiceBusCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DigitalTwins.EndpointClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	digitalTwinsId, err := parse.DigitalTwinsInstanceID(d.Get("digital_twins_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewDigitalTwinsEndpointID(subscriptionId, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, name).ID()

	if d.IsNewResource() {
		existing, err := client.Get(ctx, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Digital Twins Endpoint %q (Resource Group %q / Instance %q): %+v", name, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_digital_twins_endpoint_servicebus", id)
		}
	}

	properties := digitaltwins.EndpointResource{
		Properties: &digitaltwins.ServiceBus{
			EndpointType:              digitaltwins.EndpointTypeServiceBus,
			PrimaryConnectionString:   utils.String(d.Get("servicebus_primary_connection_string").(string)),
			SecondaryConnectionString: utils.String(d.Get("servicebus_secondary_connection_string").(string)),
			DeadLetterSecret:          utils.String(d.Get("dead_letter_storage_secret").(string)),
		},
	}

	future, err := client.CreateOrUpdate(ctx, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, name, properties)
	if err != nil {
		return fmt.Errorf("creating/updating Digital Twins Endpoint ServiceBus %q (Resource Group %q / Instance %q): %+v", name, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of the Digital Twins Endpoint ServiceBus %q (Resource Group %q / Instance %q): %+v", name, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, err)
	}

	if _, err := client.Get(ctx, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, name); err != nil {
		return fmt.Errorf("retrieving Digital Twins Endpoint ServiceBus %q (Resource Group %q / Instance %q): %+v", name, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, err)
	}

	d.SetId(id)

	return resourceDigitalTwinsEndpointServiceBusRead(d, meta)
}

func resourceDigitalTwinsEndpointServiceBusRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DigitalTwins.EndpointClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DigitalTwinsEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.DigitalTwinsInstanceName, id.EndpointName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Digital Twins ServiceBus Endpoint %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Digital Twins Endpoint ServiceBus %q (Resource Group %q / Instance %q): %+v", id.EndpointName, id.ResourceGroup, id.DigitalTwinsInstanceName, err)
	}
	d.Set("name", id.EndpointName)
	d.Set("digital_twins_id", parse.NewDigitalTwinsInstanceID(subscriptionId, id.ResourceGroup, id.DigitalTwinsInstanceName).ID())
	if resp.Properties != nil {
		if _, ok := resp.Properties.AsServiceBus(); !ok {
			return fmt.Errorf("retrieving Digital Twins Endpoint %q (Resource Group %q / Instance %q) is not type ServiceBus", id.EndpointName, id.ResourceGroup, id.DigitalTwinsInstanceName)
		}
	}
	return nil
}

func resourceDigitalTwinsEndpointServiceBusDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DigitalTwins.EndpointClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DigitalTwinsEndpointID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.DigitalTwinsInstanceName, id.EndpointName)
	if err != nil {
		return fmt.Errorf("deleting Digital Twins Endpoint ServiceBus %q (Resource Group %q / Instance %q): %+v", id.EndpointName, id.ResourceGroup, id.DigitalTwinsInstanceName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the Digital Twins Endpoint ServiceBus %q (Resource Group %q / Instance %q): %+v", id.EndpointName, id.ResourceGroup, id.DigitalTwinsInstanceName, err)
	}
	return nil
}
