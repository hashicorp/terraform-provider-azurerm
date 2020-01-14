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

			"auto_patching": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day_of_week": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"maintenance_window_duration_in_minutes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"maintenance_window_starting_hour": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"key_vault_credential": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"credential_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"azure_key_vault_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_principal_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_principal_secret": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"server_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_r_services_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"sql_connectivity_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql_connectivity_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sql_connectivity_update_password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql_connectivity_update_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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

	expandSettings := "autoPatchingSettings,keyVaultCredentialSettings,serverConfigurationsManagementSettings"
	resp, err := client.Get(ctx, resourceGroupName, name, expandSettings)
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
		d.Set("auto_patching", flattenArmSqlVirtualMachineAutoPatching(properties.AutoPatchingSettings))
		d.Set("key_vault_credential", flattenArmSqlVirtualMachineKeyVaultCredential(properties.KeyVaultCredentialSettings, d))
		d.Set("server_configuration", flattenArmSqlVirtualMachineServerConfigurationsManagement(properties.ServerConfigurationsManagementSettings, d))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
