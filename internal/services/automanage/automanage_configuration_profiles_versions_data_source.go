package automanage

func dataSourceAutomanageConfigurationProfilesVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAutomanageConfigurationProfilesVersionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"configuration_profile_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceAutomanageConfigurationProfilesVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfilesVersionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	configurationProfileName := d.Get("configuration_profile_name").(string)

	resp, err := client.Get(ctx, configurationProfileName, name, resourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Automanage ConfigurationProfilesVersion %q (Resource Group %q / configurationProfileName %q) does not exist", name, resourceGroup, configurationProfileName)
		}
		return fmt.Errorf("retrieving Automanage ConfigurationProfilesVersion %q (Resource Group %q / configurationProfileName %q): %+v", name, resourceGroup, configurationProfileName, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Automanage ConfigurationProfilesVersion %q (Resource Group %q / configurationProfileName %q) ID", name, resourceGroup, configurationProfileName)
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("configuration_profile_name", configurationProfileName)
	d.Set("location", location.NormalizeNilable(resp.Location))
	return tags.FlattenAndSet(d, resp.Tags)
}
