package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceSqlServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSqlServerRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmSqlServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlServersClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Sql Server %q was not found in Resource Group %q", name, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Sql Server %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if serverFqdn := resp.FullyQualifiedDomainName; serverFqdn != nil {
		d.Set("fqdn", serverFqdn)
	}

	return nil
}
