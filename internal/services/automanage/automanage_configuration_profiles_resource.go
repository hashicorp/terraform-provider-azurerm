package automanage

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/automanage/2022-05-04/automanage"
)

func resourceAutomanageConfigurationProfile() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomanageConfigurationProfileCreate,
		Read:   resourceAutomanageConfigurationProfileRead,
		Update: resourceAutomanageConfigurationProfileUpdate,
		Delete: resourceAutomanageConfigurationProfileDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AutomanageConfigurationProfileID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"configuration_json": {
				Type:         pluginsdk.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
			},

			"tags": tags.Schema(),
		},
	}
}
func resourceAutomanageConfigurationProfileCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Automanage.ConfigurationProfileClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewAutomanageConfigurationProfileID(subscriptionId, resourceGroup, name).ID()

	existing, err := client.Get(ctx, name, resourceGroup)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_automanage_configuration_profile", id)
	}

	configuration, err := structure.ExpandJsonFromString(d.Get("configuration_json").(string))
	if err != nil {
		return fmt.Errorf("creating azurerm_automanage_configuration_profile failed for expand json from configuration string with error msg %s", err)
	}

	parameters := automanage.ConfigurationProfile{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &automanage.ConfigurationProfileProperties{
			Configuration: configuration,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if _, err := client.CreateOrUpdate(ctx, name, resourceGroup, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id)
	return resourceAutomanageConfigurationProfileRead(d, meta)
}

func resourceAutomanageConfigurationProfileRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfileID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ConfigurationProfileName, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] automanage %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	d.Set("name", id.ConfigurationProfileName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.Properties; props != nil {
		if props.Configuration != nil {
			configurationValue := props.Configuration.(map[string]interface{})
			configurationStr, err := structure.FlattenJsonToString(configurationValue)
			if err != nil {
				return fmt.Errorf("read azurerm_automanage_configuration_profile failed for flattern json to configuration string with error msg %s", err)
			}
			d.Set("configuration_json", configurationStr)
		}
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceAutomanageConfigurationProfileUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfileID(d.Id())
	if err != nil {
		return err
	}

	parameters := automanage.ConfigurationProfileUpdate{}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ConfigurationProfileName, id.ResourceGroup, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}
	return resourceAutomanageConfigurationProfileRead(d, meta)
}

func resourceAutomanageConfigurationProfileDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfileID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.ConfigurationProfileName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}
