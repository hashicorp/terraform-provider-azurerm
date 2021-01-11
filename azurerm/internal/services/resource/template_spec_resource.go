package resource

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2019-06-01-preview/templatespecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceTemplateSpec() *schema.Resource {
	return &schema.Resource{
		Create: resourceTemplateSpecCreate,
		Read:   resourceTemplateSpecRead,
		Update: resourceTemplateSpecUpdate,
		Delete: resourceTemplateSpecDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.TemplateSpecID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.TemplateSpecName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.TemplateSpecDescription,
			},

			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.TemplateSpecDisplayName,
			},

			"tags": tags.Schema(),
		},
	}
}
func resourceTemplateSpecCreate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Resource.TemplateSpecClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewTemplateSpecID(subscriptionId, resourceGroup, name)

	existing, err := client.Get(ctx, resourceGroup, name, "versions")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Template Spec %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_template_spec", id.ID())
	}

	templateSpec := templatespecs.TemplateSpec{
		Location:   utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &templatespecs.Properties{},
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("description"); ok {
		templateSpec.Properties.Description = utils.String(v.(string))
	}

	if v, ok := d.GetOk("display_name"); ok {
		templateSpec.Properties.DisplayName = utils.String(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, templateSpec); err != nil {
		return fmt.Errorf("creating Template Spec %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceTemplateSpecRead(d, meta)
}

func resourceTemplateSpecRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.TemplateSpecClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TemplateSpecID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "versions")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Template Spec %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Template Spec %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.Properties; props != nil {
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceTemplateSpecUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.TemplateSpecClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TemplateSpecID(d.Id())
	if err != nil {
		return err
	}

	templateSpec := templatespecs.UpdateModel{}

	if d.HasChange("tags") {
		templateSpec.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, &templateSpec); err != nil {
		return fmt.Errorf("updating Template Spec %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceTemplateSpecRead(d, meta)
}

func resourceTemplateSpecDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.TemplateSpecClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TemplateSpecID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting Template Spec %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}
