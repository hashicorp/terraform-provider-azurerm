package mssqlvm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmMsSqlVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmMsSqlVirtualMachineRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"virtual_machine_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sql_license_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sql_sku": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceArmMsSqlVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQLVM.SQLVirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroupName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q) was not found", name, resourceGroupName)
		}
		return fmt.Errorf("Error reading Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group", resourceGroupName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if properties := resp.Properties; properties != nil {
		d.Set("sql_sku", string(properties.SQLImageSku))
		d.Set("sql_license_type", string(properties.SQLServerLicenseType))
		d.Set("virtual_machine_resource_id", properties.VirtualMachineResourceID)
	}
	if v := flattenArmSqlVirtualMachineAutoPatching(d); v != nil {
		d.Set("auto_patching", v)
	}
	if v := flattenArmSqlVirtualMachineKeyVaultCredential(d); v != nil {
		d.Set("key_vault_credential", v)
	}
	if v := flattenArmSqlVirtualMachineServerConfigurationsManagement(d); v != nil {
		d.Set("server_configuration", v)
	}
	if v := flattenArmSqlVirtualMachineStorageConfiguration(d); v != nil {
		d.Set("storage_configuration", v)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
