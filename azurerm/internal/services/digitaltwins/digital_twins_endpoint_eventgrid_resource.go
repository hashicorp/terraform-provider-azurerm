package digitaltwins

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/digitaltwins/mgmt/2020-10-31/digitaltwins"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/digitaltwins/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/digitaltwins/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDigitalTwinsEndpointEventGrid() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDigitalTwinsEndpointEventGridCreateUpdate,
		Read:   resourceArmDigitalTwinsEndpointEventGridRead,
		Update: resourceArmDigitalTwinsEndpointEventGridCreateUpdate,
		Delete: resourceArmDigitalTwinsEndpointEventGridDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DigitalTwinsEndpointID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DigitaltwinsInstanceName,
			},

			"digital_twins_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DigitaltwinsInstanceID,
			},

			"eventgrid_topic_endpoint": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPS,
			},

			"eventgrid_topic_primary_access_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"eventgrid_topic_secondary_access_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"dead_letter_storage_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}
func resourceArmDigitalTwinsEndpointEventGridCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DigitalTwins.EndpointClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	digitalTwinsId, err := parse.DigitalTwinsInstanceID(d.Get("digital_twins_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewDigitalTwinsEndpointID(subscriptionId, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, name).ID("")

	if d.IsNewResource() {
		existing, err := client.Get(ctx, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Digital Twins Endpoint %q (Resource Group %q / Instance %q): %+v", name, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_digital_twins_endpoint_eventgrid", id)
		}
	}

	properties := digitaltwins.EndpointResource{
		Properties: &digitaltwins.EventGrid{
			EndpointType:     digitaltwins.EndpointTypeEventGrid,
			TopicEndpoint:    utils.String(d.Get("eventgrid_topic_endpoint").(string)),
			AccessKey1:       utils.String(d.Get("eventgrid_topic_primary_access_key").(string)),
			AccessKey2:       utils.String(d.Get("eventgrid_topic_secondary_access_key").(string)),
			DeadLetterSecret: utils.String(d.Get("dead_letter_storage_secret").(string)),
		},
	}

	future, err := client.CreateOrUpdate(ctx, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, name, properties)
	if err != nil {
		return fmt.Errorf("creating/updating Digital Twins EventGrid Endpoint %q (Resource Group %q / Instance %q): %+v", name, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of the Digital Twins EventGrid Endpoint %q (Resource Group %q / Instance %q): %+v", name, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, err)
	}

	if _, err := client.Get(ctx, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, name); err != nil {
		return fmt.Errorf("retrieving Digital Twins EventGrid Endpoint %q (Resource Group %q / Instance %q): %+v", name, digitalTwinsId.ResourceGroup, digitalTwinsId.Name, err)
	}

	d.SetId(id)

	return resourceArmDigitalTwinsEndpointEventGridRead(d, meta)
}

func resourceArmDigitalTwinsEndpointEventGridRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[INFO] Digital Twins EventGrid Endpoint %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Digital Twins EventGrid Endpoint %q (Resource Group %q / Instance %q): %+v", id.EndpointName, id.ResourceGroup, id.DigitalTwinsInstanceName, err)
	}
	d.Set("name", id.EndpointName)
	d.Set("digital_twins_id", parse.NewDigitalTwinsInstanceID(subscriptionId, id.ResourceGroup, id.DigitalTwinsInstanceName).ID(""))
	if resp.Properties != nil {
		if _, ok := resp.Properties.AsEventGrid(); !ok {
			return fmt.Errorf("retrieving Digital Twins Endpoint %q (Resource Group %q / Instance %q) is not type Event Grid", id.EndpointName, id.ResourceGroup, id.DigitalTwinsInstanceName)
		}
	}

	return nil
}

func resourceArmDigitalTwinsEndpointEventGridDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DigitalTwins.EndpointClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DigitalTwinsEndpointID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.DigitalTwinsInstanceName, id.EndpointName)
	if err != nil {
		return fmt.Errorf("deleting Digital Twins EventGrid Endpoint %q (Resource Group %q / Instance %q): %+v", id.EndpointName, id.ResourceGroup, id.DigitalTwinsInstanceName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the Digital Twins EventGrid Endpoint %q (Resource Group %q / Instance %q): %+v", id.EndpointName, id.ResourceGroup, id.DigitalTwinsInstanceName, err)
	}
	return nil
}
