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

func resourceArmDigitaltwinsEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDigitaltwinsEndpointCreateUpdate,
		Read:   resourceArmDigitaltwinsEndpointRead,
		Update: resourceArmDigitaltwinsEndpointCreateUpdate,
		Delete: resourceArmDigitaltwinsEndpointDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DigitaltwinsEndpointID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DigitaltwinsName(),
			},

			"digital_twins_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DigitaltwinsID,
			},

			"eventgrid_topic_endpoint": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPS,
			},

			"eventgrid_topic_primary_access_key": {
				Type:     schema.TypeString,
				Required: true,
			},

			"eventgrid_topic_secondary_access_key": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
func resourceArmDigitaltwinsEndpointCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Digitaltwins.EndpointClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	digitalTwinsID, _ := parse.DigitalTwinsID(d.Get("digital_twins_id").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, digitalTwinsID.ResourceGroup, digitalTwinsID.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Digital Twins Endpoint %q (Resource Group %q / Digital Twins Name %q): %+v", name, digitalTwinsID.ResourceGroup, digitalTwinsID.Name, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_digital_twins_endpoint_eventgrid", *existing.ID)
		}
	}

	endpointDescription := digitaltwins.EndpointResource{
		Properties: &digitaltwins.EventGrid{
			EndpointType:  digitaltwins.EndpointTypeEventGrid,
			TopicEndpoint: utils.String(d.Get("eventgrid_topic_endpoint").(string)),
			AccessKey1:    utils.String(d.Get("eventgrid_topic_primary_access_key").(string)),
			AccessKey2:    utils.String(d.Get("eventgrid_topic_secondary_access_key").(string)),
		},
	}

	future, err := client.CreateOrUpdate(ctx, digitalTwinsID.ResourceGroup, digitalTwinsID.Name, name, endpointDescription)
	if err != nil {
		return fmt.Errorf("creating/updating Digital Twins Endpoint %q (Resource Group %q / Digital Twins Name %q): %+v", name, digitalTwinsID.ResourceGroup, digitalTwinsID.Name, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating/updating future for Digital Twins Endpoint %q (Resource Group %q / Digital Twins Name %q): %+v", name, digitalTwinsID.ResourceGroup, digitalTwinsID.Name, err)
	}

	resp, err := client.Get(ctx, digitalTwinsID.ResourceGroup, digitalTwinsID.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Digital Twins Endpoint %q (Resource Group %q / Digital Twins Name %q): %+v", name, digitalTwinsID.ResourceGroup, digitalTwinsID.Name, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Digital Twins Endpoint %q (Resource Group %q / Digital Twins Name %q) ID", name, digitalTwinsID.ResourceGroup, digitalTwinsID.Name)
	}

	id, err := parse.DigitaltwinsEndpointID(*resp.ID)
	if err != nil {
		return err
	}
	d.SetId(id.ID(subscriptionId))

	return resourceArmDigitaltwinsEndpointRead(d, meta)
}

func resourceArmDigitaltwinsEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Digitaltwins.EndpointClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DigitaltwinsEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ResourceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] digitaltwins %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Digital Twins Endpoint %q (Resource Group %q / Digital Twins Name %q): %+v", id.Name, id.ResourceGroup, id.ResourceName, err)
	}
	d.Set("name", id.Name)
	d.Set("digital_twins_id", parse.NewDigitalTwinsID(id.ResourceGroup, id.ResourceName).ID(client.SubscriptionID))
	if resp.Properties != nil {
		if _, ok := resp.Properties.AsEventGrid(); !ok {
			return fmt.Errorf("retrieving Digital Twins Endpoint %q (Resource Group %q / Digital Twins Name %q) is not type Event Grid", id.Name, id.ResourceGroup, id.ResourceName)
		}
	}
	return nil
}

func resourceArmDigitaltwinsEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Digitaltwins.EndpointClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DigitaltwinsEndpointID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ResourceName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Digital Twins Endpoint %q (Resource Group %q / Digital Twins Name %q): %+v", id.Name, id.ResourceGroup, id.ResourceName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for Digital Twins Endpoint %q (Resource Group %q / Digital Twins Name %q): %+v", id.Name, id.ResourceGroup, id.ResourceName, err)
	}
	return nil
}
