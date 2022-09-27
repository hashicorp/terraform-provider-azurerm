package automanage

func dataSourceAutomanageConfigurationProfileHCRPAssignment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAutomanageConfigurationProfileHCRPAssignmentRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"machine_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceAutomanageConfigurationProfileHCRPAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileHCRPAssignmentClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	machineName := d.Get("machine_name").(string)

	resp, err := client.Get(ctx, resourceGroup, machineName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Automanage ConfigurationProfileHCRPAssignment %q (Resource Group %q / machineName %q) does not exist", name, resourceGroup, machineName)
		}
		return fmt.Errorf("retrieving Automanage ConfigurationProfileHCRPAssignment %q (Resource Group %q / machineName %q): %+v", name, resourceGroup, machineName, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Automanage ConfigurationProfileHCRPAssignment %q (Resource Group %q / machineName %q) ID", name, resourceGroup, machineName)
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("machine_name", machineName)
	return nil
}
