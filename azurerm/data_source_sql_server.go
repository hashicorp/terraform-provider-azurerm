package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
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

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"administrator_login": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmSqlServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sql.ServersClient
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

	if id := resp.ID; id != nil {
		d.SetId(*resp.ID)
	}

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ServerProperties; props != nil {
		d.Set("fqdn", props.FullyQualifiedDomainName)
		d.Set("version", props.Version)
		d.Set("administrator_login", props.AdministratorLogin)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
