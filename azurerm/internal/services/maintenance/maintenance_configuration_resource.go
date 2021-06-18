package maintenance

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/maintenance/mgmt/2021-05-01/maintenance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMaintenanceConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmMaintenanceConfigurationCreateUpdate,
		Read:   resourceArmMaintenanceConfigurationRead,
		Update: resourceArmMaintenanceConfigurationCreateUpdate,
		Delete: resourceArmMaintenanceConfigurationDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ConfigurationV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.MaintenanceConfigurationIDInsensitively(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			// TODO use `azure.SchemaResourceGroupName()` in version 3.0
			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/8653
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"scope": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "All",
				ValidateFunc: validation.StringInSlice([]string{
					"All",
					string(maintenance.ScopeExtension),
					string(maintenance.ScopeHost),
					string(maintenance.ScopeInGuestPatch),
					string(maintenance.ScopeOSImage),
					string(maintenance.ScopeSQLDB),
					string(maintenance.ScopeSQLManagedInstance),
				}, false),
			},

			"visibility": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(maintenance.VisibilityCustom),
				ValidateFunc: validation.StringInSlice([]string{
					string(maintenance.VisibilityCustom),
					string(maintenance.VisibilityPublic),
				}, false),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMaintenanceConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	scope := d.Get("scope").(string)
	visibility := d.Get("visibility").(string)
	if visibility == string(maintenance.VisibilityPublic) {
		if !(scope == string(maintenance.ScopeSQLDB) || scope == string(maintenance.ScopeSQLManagedInstance)) {
			return fmt.Errorf("`scope` must be set to %s or %s when `visibility` is %s", string(maintenance.ScopeSQLDB), string(maintenance.ScopeSQLManagedInstance), visibility)
		}
	}

	id := parse.NewMaintenanceConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_maintenance_configuration", id.ID())
		}
	}

	configuration := maintenance.Configuration{
		Name:     utils.String(id.Name),
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		ConfigurationProperties: &maintenance.ConfigurationProperties{
			MaintenanceScope: maintenance.Scope(scope),
			Visibility:       maintenance.Visibility(visibility),
			Namespace:        utils.String("Microsoft.Maintenance"),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, configuration); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceArmMaintenanceConfigurationRead(d, meta)
}

func resourceArmMaintenanceConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceConfigurationIDInsensitively(d.Id())
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
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.ConfigurationProperties; props != nil {
		d.Set("scope", props.MaintenanceScope)
		d.Set("visibility", props.Visibility)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMaintenanceConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceConfigurationIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}
