package communication

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/communication/mgmt/2020-08-20/communication"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/communication/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/communication/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCommunicationService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCommunicationServiceCreateUpdate,
		Read:   resourceArmCommunicationServiceRead,
		Update: resourceArmCommunicationServiceCreateUpdate,
		Delete: resourceArmCommunicationServiceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CommunicationServiceID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CommunicationServiceName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"data_location": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "United States",
				ValidateFunc: validation.StringInSlice([]string{
					"Asia Pacific",
					"Australia",
					"Europe",
					"UK",
					"United States",
				}, false),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmCommunicationServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Communication.ServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewCommunicationServiceID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_communication_service", id.ID())
		}
	}

	parameter := communication.ServiceResource{
		// The location is always `global` from the Azure Portal
		Location: utils.String(location.Normalize("global")),
		ServiceProperties: &communication.ServiceProperties{
			DataLocation: utils.String(d.Get("data_location").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, &parameter)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmCommunicationServiceRead(d, meta)
}

func resourceArmCommunicationServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Communication.ServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CommunicationServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.ServiceProperties; props != nil {
		d.Set("data_location", props.DataLocation)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmCommunicationServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Communication.ServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CommunicationServiceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
