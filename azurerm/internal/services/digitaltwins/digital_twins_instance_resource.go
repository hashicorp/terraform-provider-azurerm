package digitaltwins

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/digitaltwins/mgmt/2020-10-31/digitaltwins"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/digitaltwins/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/digitaltwins/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDigitalTwinsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDigitalTwinsInstanceCreate,
		Read:   resourceArmDigitalTwinsInstanceRead,
		Update: resourceArmDigitalTwinsInstanceUpdate,
		Delete: resourceArmDigitalTwinsInstanceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DigitalTwinsInstanceID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DigitaltwinsInstanceName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}
func resourceArmDigitalTwinsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DigitalTwins.InstanceClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewDigitalTwinsInstanceID(resourceGroup, name)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing Digital Twins Instance %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_digital_twins_instance", *existing.ID)
	}

	properties := digitaltwins.Description{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, properties)
	if err != nil {
		return fmt.Errorf("creating Digital Twins Instance %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Digital Twins Instance %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	_, err = client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Digital Twins Instance %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id.ID(subscriptionId))

	return resourceArmDigitalTwinsInstanceRead(d, meta)
}

func resourceArmDigitalTwinsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DigitalTwins.InstanceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DigitalTwinsInstanceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Digital Twins Instance %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Digital Twins Instance %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.Properties; props != nil {
		d.Set("host_name", props.HostName)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDigitalTwinsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DigitalTwins.InstanceClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DigitalTwinsInstanceID(d.Id())
	if err != nil {
		return err
	}

	props := digitaltwins.PatchDescription{}

	if d.HasChange("tags") {
		props.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, props); err != nil {
		return fmt.Errorf("updating Digital Twins Instance %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceArmDigitalTwinsInstanceRead(d, meta)
}

func resourceArmDigitalTwinsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DigitalTwins.InstanceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DigitalTwinsInstanceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Digital Twins Instance %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for Digital Twins Instance %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}
