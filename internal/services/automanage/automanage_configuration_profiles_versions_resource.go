package automanage

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/mgmt/2022-05-04/automanage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"log"
	"time"
)

func resourceAutomanageConfigurationProfilesVersion() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomanageConfigurationProfilesVersionCreateUpdate,
		Read:   resourceAutomanageConfigurationProfilesVersionRead,
		Update: resourceAutomanageConfigurationProfilesVersionCreateUpdate,
		Delete: resourceAutomanageConfigurationProfilesVersionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AutomanageConfigurationProfilesVersionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"configuration_profile_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"configuration": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}
func resourceAutomanageConfigurationProfilesVersionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Automanage.ConfigurationProfilesVersionClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	configurationProfileName := d.Get("configuration_profile_name").(string)

	id := parse.NewAutomanageConfigurationProfilesVersionID(subscriptionId, resourceGroup, configurationProfileName, name).ID()

	if d.IsNewResource() {
		existing, err := client.Get(ctx, configurationProfileName, name, resourceGroup)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing Automanage ConfigurationProfilesVersion %q (Resource Group %q / configurationProfileName %q): %+v", name, resourceGroup, configurationProfileName, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_automanage_configuration_profiles_version", id)
		}
	}

	configuration, err := structure.ExpandJsonFromString(d.Get("configuration").(string))
	if err != nil {
		return fmt.Errorf("creating/updating Automanage ConfigurationProfilesVersion failed for expand json from configuration string with error msg %s", err)
	}

	parameters := automanage.ConfigurationProfile{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &automanage.ConfigurationProfileProperties{
			Configuration: configuration,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if _, err := client.CreateOrUpdate(ctx, configurationProfileName, name, resourceGroup, parameters); err != nil {
		return fmt.Errorf("creating/updating Automanage ConfigurationProfilesVersion %q (Resource Group %q / configurationProfileName %q): %+v", name, resourceGroup, configurationProfileName, err)
	}

	d.SetId(id)
	return resourceAutomanageConfigurationProfilesVersionRead(d, meta)
}

func resourceAutomanageConfigurationProfilesVersionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfilesVersionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfilesVersionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ConfigurationProfileName, id.VersionName, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] automanage %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Automanage ConfigurationProfilesVersion %q (Resource Group %q / configurationProfileName %q): %+v", id.VersionName, id.ResourceGroup, id.ConfigurationProfileName, err)
	}
	d.Set("name", id.VersionName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("configuration_profile_name", id.ConfigurationProfileName)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.Properties; props != nil {
		if props.Configuration != nil {
			configurationValue := props.Configuration.(map[string]interface{})
			configurationStr, _ := structure.FlattenJsonToString(configurationValue)
			d.Set("configuration", configurationStr)
		}
	}
	d.Set("type", resp.Type)
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceAutomanageConfigurationProfilesVersionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfilesVersionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfilesVersionID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.VersionName, id.ConfigurationProfileName); err != nil {
		return fmt.Errorf("deleting Automanage ConfigurationProfilesVersion %q (Resource Group %q / configurationProfileName %q): %+v", id.ConfigurationProfileName, id.ResourceGroup, id.ConfigurationProfileName, err)
	}
	return nil
}
