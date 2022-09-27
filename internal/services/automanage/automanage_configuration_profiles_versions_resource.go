package automanage

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/mgmt/2022-05-04/automanage"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
)

func resourceAutomanageConfigurationProfilesVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutomanageConfigurationProfilesVersionCreateUpdate,
		Read:   resourceAutomanageConfigurationProfilesVersionRead,
		Update: resourceAutomanageConfigurationProfilesVersionCreateUpdate,
		Delete: resourceAutomanageConfigurationProfilesVersionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AutomanageConfigurationProfilesVersionID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"configuration_profile_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"configuration": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}
func resourceAutomanageConfigurationProfilesVersionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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

	configuration, _ := structure.ExpandJsonFromString(d.Get("configuration").(string))

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

func resourceAutomanageConfigurationProfilesVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfilesVersionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfilesVersionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ConfigurationProfileName, id.Name, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] automanage %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Automanage ConfigurationProfilesVersion %q (Resource Group %q / configurationProfileName %q): %+v", id.Name, id.ResourceGroup, id.ConfigurationProfileName, err)
	}
	d.Set("name", id.Name)
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

func resourceAutomanageConfigurationProfilesVersionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfilesVersionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfilesVersionID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.ConfigurationProfileName, id.Name); err != nil {
		return fmt.Errorf("deleting Automanage ConfigurationProfilesVersion %q (Resource Group %q / configurationProfileName %q): %+v", id.Name, id.ResourceGroup, id.ConfigurationProfileName, err)
	}
	return nil
}
