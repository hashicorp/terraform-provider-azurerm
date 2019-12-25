package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmMsSqlVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmMsSqlVirtualMachineRead,

		Schema: map[string]*schema.Schema{
			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"virtual_machine_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sql_server_license_type": {
				Type:     schema.TypeString,
				Computed: true,
			},


			"sql_image_sku": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sql_virtual_machine_group_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"auto_patching_settings": {
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
						"maintenance_window_duration": {
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

			"key_vault_credential_settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"azure_key_vault_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"credential_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"service_principal_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_principal_secret": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"server_configurations_management_settings": {
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
						"sql_connectivity_auth_update_password": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"sql_connectivity_auth_update_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

					},
				},
			},

			"storage_configuration_settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sql_data_default_file_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql_data_luns": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"sql_log_default_file_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql_log_luns": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"sql_temp_db_default_file_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql_temp_db_luns": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"storage_workload_type": {
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
	client := meta.(*ArmClient).MSSQLVM.SQLVirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	resourceGroupName := d.Get("resource_group").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroupName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q) was not found", name, resourceGroupName)
		}
		return fmt.Errorf("Error reading Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	d.SetId(*resp.ID)

	d.Set("resource_group", resourceGroupName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if properties := resp.Properties; properties != nil {
		if err := d.Set("auto_patching_settings", flattenArmSqlVirtualMachineAutoPatchingSettings(properties.AutoPatchingSettings)); err != nil {
			return fmt.Errorf("Error setting `auto_patching_settings`: %+v", err)
		}
		if err := d.Set("key_vault_credential_settings", flattenArmSqlVirtualMachineKeyVaultCredentialSettings(properties.KeyVaultCredentialSettings)); err != nil {
			return fmt.Errorf("Error setting `key_vault_credential_settings`: %+v", err)
		}
		d.Set("provisioning_state", properties.ProvisioningState)
		d.Set("sql_image_sku", string(properties.SQLImageSku))
		d.Set("sql_management", string(properties.SQLManagement))
		d.Set("sql_server_license_type", string(properties.SQLServerLicenseType))
		d.Set("sql_virtual_machine_group_resource_id", properties.SQLVirtualMachineGroupResourceID)
		d.Set("virtual_machine_resource_id", properties.VirtualMachineResourceID)
		if err := d.Set("server_configurations_management_settings", flattenArmSqlVirtualMachineServerConfigurationsManagementSettings(properties.ServerConfigurationsManagementSettings)); err != nil {
			return fmt.Errorf("Error setting `server_configurations_management_settings`: %+v", err)
		}
		if err := d.Set("storage_configuration_settings", flattenArmSqlVirtualMachineStorageConfigurationSettings(properties.StorageConfigurationSettings)); err != nil {
			return fmt.Errorf("Error setting `storage_configuration_settings`: %+v", err)
		}
	}
	d.Set("name", name)
	d.Set("id", resp.ID)

	return tags.FlattenAndSet(d, resp.Tags)
}
