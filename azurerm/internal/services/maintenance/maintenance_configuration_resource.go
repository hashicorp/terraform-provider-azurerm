package maintenance

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/maintenance/mgmt/2018-06-01-preview/maintenance"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMaintenanceConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMaintenanceConfigurationCreateUpdate,
		Read:   resourceArmMaintenanceConfigurationRead,
		Update: resourceArmMaintenanceConfigurationCreateUpdate,
		Delete: resourceArmMaintenanceConfigurationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.MaintenanceConfigurationID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/8653
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(maintenance.ScopeAll),
				ValidateFunc: validation.StringInSlice([]string{
					string(maintenance.ScopeAll),
					string(maintenance.ScopeHost),
					string(maintenance.ScopeInResource),
					string(maintenance.ScopeResource),
				}, false),
			},

			// There's a bug in the Azure API where the the key of tags is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/9075
			// use custom tags defition here to prevent inputting upper case key
			"tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				ValidateFunc: validate.TagsWithLowerCaseKey,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceArmMaintenanceConfigurationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("failure checking for present of existing MaintenanceConfiguration %q (Resource Group %q): %+v", name, resGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_maintenance_configuration", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	scope := d.Get("scope").(string)

	configuration := maintenance.Configuration{
		Name:     utils.String(name),
		Location: utils.String(location),
		ConfigurationProperties: &maintenance.ConfigurationProperties{
			MaintenanceScope: maintenance.Scope(scope),
			Namespace:        utils.String("Microsoft.Maintenance"),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.CreateOrUpdate(ctx, resGroup, name, configuration); err != nil {
		return fmt.Errorf("failure creating/updating MaintenanceConfiguration %q (Resource Group %q): %+v", name, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("failure retrieving MaintenanceConfiguration %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("cannot read MaintenanceConfiguration %q (Resource Group %q) ID", name, resGroup)
	}

	d.SetId(*resp.ID)
	return resourceArmMaintenanceConfigurationRead(d, meta)
}

func resourceArmMaintenanceConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] maintenance %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failure retrieving MaintenanceConfiguration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.ConfigurationProperties; props != nil {
		d.Set("scope", props.MaintenanceScope)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMaintenanceConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceConfigurationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("failure deleting MaintenanceConfiguration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return nil
}
