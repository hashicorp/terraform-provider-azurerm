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

func resourceArmDigitaltwinsDigitalTwin() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDigitaltwinsDigitalTwinCreate,
		Read:   resourceArmDigitaltwinsDigitalTwinRead,
		Update: resourceArmDigitaltwinsDigitalTwinUpdate,
		Delete: resourceArmDigitaltwinsDigitalTwinDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DigitaltwinsDigitalTwinID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DigitaltwinsName(),
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
func resourceArmDigitaltwinsDigitalTwinCreate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Digitaltwins.DigitalTwinClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Digitaltwins DigitalTwin %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_digital_twins", *existing.ID)
	}

	digitalTwinsCreate := digitaltwins.Description{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, digitalTwinsCreate)
	if err != nil {
		return fmt.Errorf("creating Digitaltwins DigitalTwin %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating future for Digitaltwins DigitalTwin %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Digitaltwins DigitalTwin %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Digitaltwins DigitalTwin %q (Resource Group %q) ID", name, resourceGroup)
	}

	id, err := parse.DigitaltwinsDigitalTwinID(*resp.ID)
	if err != nil {
		return err
	}
	d.SetId(id.ID(subscriptionId))

	return resourceArmDigitaltwinsDigitalTwinRead(d, meta)
}

func resourceArmDigitaltwinsDigitalTwinRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Digitaltwins.DigitalTwinClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DigitaltwinsDigitalTwinID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] digitaltwins %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Digitaltwins DigitalTwin %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.Properties; props != nil {
		d.Set("host_name", props.HostName)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDigitaltwinsDigitalTwinUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Digitaltwins.DigitalTwinClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DigitaltwinsDigitalTwinID(d.Id())
	if err != nil {
		return err
	}

	digitalTwinsPatchDescription := digitaltwins.PatchDescription{}

	if d.HasChange("tags") {
		digitalTwinsPatchDescription.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, digitalTwinsPatchDescription); err != nil {
		return fmt.Errorf("updating Digitaltwins DigitalTwin %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return resourceArmDigitaltwinsDigitalTwinRead(d, meta)
}

func resourceArmDigitaltwinsDigitalTwinDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Digitaltwins.DigitalTwinClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DigitaltwinsDigitalTwinID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Digitaltwins DigitalTwin %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for Digitaltwins DigitalTwin %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return nil
}
