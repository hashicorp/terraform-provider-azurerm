package automanage

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/mgmt/2022-05-04/automanage"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
)

func resourceAutomanageConfigurationProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutomanageConfigurationProfileCreate,
		Read:   resourceAutomanageConfigurationProfileRead,
		Update: resourceAutomanageConfigurationProfileUpdate,
		Delete: resourceAutomanageConfigurationProfileDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AutomanageConfigurationProfileID(id)
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
func resourceAutomanageConfigurationProfileCreate(d *schema.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("checking for existing Automanage ConfigurationProfile %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_automanage_configuration_profile", id)
	}

	configuration, _ := structure.ExpandJsonFromString(d.Get("configuration").(string))

	parameters := automanage.ConfigurationProfile{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &automanage.ConfigurationProfileProperties{
			Configuration: configuration,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if _, err := client.CreateOrUpdate(ctx, name, resourceGroup, parameters); err != nil {
		return fmt.Errorf("creating Automanage ConfigurationProfile %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id)
	return resourceAutomanageConfigurationProfileRead(d, meta)
}

func resourceAutomanageConfigurationProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfileID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Name, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] automanage %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Automanage ConfigurationProfile %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
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

func resourceAutomanageConfigurationProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfileID(d.Id())
	if err != nil {
		return err
	}

	parameters := automanage.ConfigurationProfileUpdate{
		Properties: &automanage.ConfigurationProfileProperties{},
	}
	if d.HasChange("configuration") {
		configuration, _ := structure.ExpandJsonFromString(d.Get("configuration").(string))
		parameters.Properties.Configuration = configuration
	}
	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.Name, id.ResourceGroup, parameters); err != nil {
		return fmt.Errorf("updating Automanage ConfigurationProfile %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return resourceAutomanageConfigurationProfileRead(d, meta)
}

func resourceAutomanageConfigurationProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfileID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting Automanage ConfigurationProfile %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return nil
}
