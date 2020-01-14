package mssqlvm

import (
	"fmt"
	"strings"
	"time"

	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/sqlvirtualmachine/mgmt/2017-03-01-preview/sqlvirtualmachine"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMsSqlVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMsSqlVirtualMachineCreateUpdate,
		Read:   resourceArmMsSqlVirtualMachineRead,
		Update: resourceArmMsSqlVirtualMachineCreateUpdate,
		Delete: resourceArmMsSqlVirtualMachineDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"location": azure.SchemaLocation(),

			"virtual_machine_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"sql_license_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sqlvirtualmachine.PAYG),
					string(sqlvirtualmachine.AHUB),
				}, false),
			},

			"sql_sku": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sqlvirtualmachine.Developer),
					string(sqlvirtualmachine.Express),
					string(sqlvirtualmachine.Standard),
					string(sqlvirtualmachine.Enterprise),
					string(sqlvirtualmachine.Web),
				}, false),
			},

			"sql_virtual_machine_group_resource_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"auto_patching": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"day_of_week": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(sqlvirtualmachine.Monday),
								string(sqlvirtualmachine.Tuesday),
								string(sqlvirtualmachine.Wednesday),
								string(sqlvirtualmachine.Thursday),
								string(sqlvirtualmachine.Friday),
								string(sqlvirtualmachine.Saturday),
								string(sqlvirtualmachine.Sunday),
							}, false),
						},
						"maintenance_window_duration_in_minutes": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(30, 180),
						},
						"maintenance_window_starting_hour": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 23),
						},
					},
				},
			},

			"key_vault_credential": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"credential_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"azure_key_vault_url": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validate.URLIsHTTPS,
						},
						"service_principal_name": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"service_principal_secret": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"server_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_r_services_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"sql_connectivity_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(sqlvirtualmachine.LOCAL),
								string(sqlvirtualmachine.PRIVATE),
								string(sqlvirtualmachine.PUBLIC),
							}, false),
						},
						"sql_connectivity_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validate.PortNumber,
						},
						"sql_connectivity_update_password": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"sql_connectivity_update_user_name": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"tags": tags.Schema(),

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmMsSqlVirtualMachineCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQLVM.SQLVirtualMachinesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vmResourceId := d.Get("virtual_machine_resource_id").(string)
	id, err := azure.ParseAzureResourceID(vmResourceId)
	if err != nil {
		return fmt.Errorf("Error creating Sql Virtual Machine from virtual machine %s: %+v", vmResourceId, err)
	}
	name := id.Path["virtualMachines"]

	resourceGroupName := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_sql_virtual_machine", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	properties := sqlvirtualmachine.Properties{
		SQLServerLicenseType:                   sqlvirtualmachine.SQLServerLicenseType(d.Get("sql_license_type").(string)),
		VirtualMachineResourceID:               utils.String(d.Get("virtual_machine_resource_id").(string)),
		AutoPatchingSettings:                   expandArmSqlVirtualMachineAutoPatchingSettings(d),
		KeyVaultCredentialSettings:             expandArmSqlVirtualMachineKeyVaultCredential(d),
		ServerConfigurationsManagementSettings: expandArmSqlVirtualMachineServerConfigurationsManagement(d),
	}

	if sqlVirtualMachineGroupResourceID, ok := d.GetOk("sql_virtual_machine_group_resource_id"); ok {
		properties.SQLVirtualMachineGroupResourceID = utils.String(sqlVirtualMachineGroupResourceID.(string))
	}

	if sqlImageSku, ok := d.GetOk("sql_sku"); ok {
		properties.SQLImageSku = sqlvirtualmachine.SQLImageSku(sqlImageSku.(string))
	}

	parameters := sqlvirtualmachine.SQLVirtualMachine{
		Location:   utils.String(location),
		Properties: &properties,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroupName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	expandSettings := getExpandSettings(d)
	resp, err := client.Get(ctx, resourceGroupName, name, expandSettings)
	if err != nil {
		return fmt.Errorf("Error retrieving Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q) ID", name, resourceGroupName)
	}
	d.SetId(*resp.ID)

	return resourceArmMsSqlVirtualMachineRead(d, meta)
}

func resourceArmMsSqlVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQLVM.SQLVirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	name := id.Path["sqlVirtualMachines"]

	expandSettings := getExpandSettings(d)
	resp, err := client.Get(ctx, resourceGroupName, name, expandSettings)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Sql Virtual Machine %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroupName)
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

func resourceArmMsSqlVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQLVM.SQLVirtualMachinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	name := id.Path["sqlVirtualMachines"]

	future, err := client.Delete(ctx, resourceGroupName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
		}
	}

	return nil
}

func getExpandSettings(d *schema.ResourceData) string {
	var result []string
	if _, ok := d.GetOk("auto_patching"); ok {
		result = append(result, "autoPatchingSettings")
	}
	if _, ok := d.GetOk("key_vault_credential"); ok {
		result = append(result, "keyVaultCredentialSettings")
	}
	if _, ok := d.GetOk("server_configuration"); ok {
		result = append(result, "serverConfigurationsManagementSettings")
	}
	return strings.Join(result, ",")
}

func expandArmSqlVirtualMachineAutoPatchingSettings(d *schema.ResourceData) *sqlvirtualmachine.AutoPatchingSettings {
	autoPatchingSettings := d.Get("auto_patching").([]interface{})
	if len(autoPatchingSettings) == 0 {
		return nil
	}
	autoPatchingSetting := autoPatchingSettings[0].(map[string]interface{})
	result := sqlvirtualmachine.AutoPatchingSettings{}

	if v, ok := autoPatchingSetting["enable"]; ok {
		result.Enable = utils.Bool(v.(bool))
	}
	if v, ok := autoPatchingSetting["maintenance_window_duration_in_minutes"]; ok {
		result.MaintenanceWindowDuration = utils.Int32(int32(v.(int)))
	}
	if v, ok := autoPatchingSetting["maintenance_window_starting_hour"]; ok {
		result.MaintenanceWindowStartingHour = utils.Int32(int32(v.(int)))
	}
	if dayOfWeek, ok := autoPatchingSetting["day_of_week"]; ok {
		result.DayOfWeek = sqlvirtualmachine.DayOfWeek(dayOfWeek.(string))
	}
	return &result
}

func flattenArmSqlVirtualMachineAutoPatching(autoPatching *sqlvirtualmachine.AutoPatchingSettings) []interface{} {
	if autoPatching == nil {
		return []interface{}{}
	}
	result := make(map[string]interface{})
	result["enable"] = autoPatching.Enable
	if maintenanceWindowDuration := autoPatching.MaintenanceWindowDuration; maintenanceWindowDuration != nil {
		result["maintenance_window_duration_in_minutes"] = autoPatching.MaintenanceWindowDuration
	}
	if maintenanceWindowStartingHour := autoPatching.MaintenanceWindowStartingHour; maintenanceWindowStartingHour != nil {
		result["maintenance_window_starting_hour"] = autoPatching.MaintenanceWindowStartingHour
	}
	result["day_of_week"] = string(autoPatching.DayOfWeek)
	return []interface{}{result}
}

func expandArmSqlVirtualMachineKeyVaultCredential(d *schema.ResourceData) *sqlvirtualmachine.KeyVaultCredentialSettings {
	keyVaultCredentialSettings := d.Get("key_vault_credential").([]interface{})
	if len(keyVaultCredentialSettings) == 0 {
		return nil
	}
	keyVaultCredentialSetting := keyVaultCredentialSettings[0].(map[string]interface{})
	result := sqlvirtualmachine.KeyVaultCredentialSettings{}

	if v, ok := keyVaultCredentialSetting["azure_key_vault_url"]; ok {
		result.AzureKeyVaultURL = utils.String(v.(string))
	}
	if v, ok := keyVaultCredentialSetting["credential_name"]; ok {
		result.CredentialName = utils.String(v.(string))
	}
	if v, ok := keyVaultCredentialSetting["service_principal_name"]; ok {
		result.ServicePrincipalName = utils.String(v.(string))
	}
	if v, ok := keyVaultCredentialSetting["service_principal_secret"]; ok {
		result.ServicePrincipalSecret = utils.String(v.(string))
	}
	if v, ok := keyVaultCredentialSetting["enable"]; ok {
		result.Enable = utils.Bool(v.(bool))
	}

	return &result
}

func flattenArmSqlVirtualMachineKeyVaultCredential(keyVault *sqlvirtualmachine.KeyVaultCredentialSettings, d *schema.ResourceData) []interface{} {
	if keyVault == nil {
		return []interface{}{}
	}
	result := make(map[string]interface{})
	result["enable"] = keyVault.Enable
	if credentialName := keyVault.CredentialName; credentialName != nil {
		result["credential_name"] = credentialName
	}
	if v, ok := d.GetOk("key_vault_credential.0.azure_key_vault_url"); ok {
		result["azure_key_vault_url"] = v
	}
	if v, ok := d.GetOk("key_vault_credential.0.service_principal_name"); ok {
		result["service_principal_name"] = v
	}
	if v, ok := d.GetOk("key_vault_credential.0.service_principal_secret"); ok {
		result["service_principal_secret"] = v
	}

	return []interface{}{result}
}

func expandArmSqlVirtualMachineServerConfigurationsManagement(d *schema.ResourceData) *sqlvirtualmachine.ServerConfigurationsManagementSettings {
	serverConfigMMs := d.Get("server_configuration").([]interface{})
	if len(serverConfigMMs) == 0 {
		return nil
	}
	serverConfigMM := serverConfigMMs[0].(map[string]interface{})

	result := sqlvirtualmachine.ServerConfigurationsManagementSettings{}
	sqlConnectivityUpdateSettings := sqlvirtualmachine.SQLConnectivityUpdateSettings{}
	//additional feature
	if v, ok := serverConfigMM["is_r_services_enabled"]; ok {
		result.AdditionalFeaturesServerConfigurations = &sqlvirtualmachine.AdditionalFeaturesServerConfigurations{IsRServicesEnabled: utils.Bool(v.(bool))}
	}
	//connectivity
	if connectivityType, ok := serverConfigMM["sql_connectivity_type"]; ok {
		sqlConnectivityUpdateSettings.ConnectivityType = sqlvirtualmachine.ConnectivityType(connectivityType.(string))
	}
	if v, ok := serverConfigMM["sql_connectivity_port"]; ok {
		sqlConnectivityUpdateSettings.Port = utils.Int32(int32(v.(int)))
	}
	result.SQLConnectivityUpdateSettings = &sqlConnectivityUpdateSettings

	return &result
}

func flattenArmSqlVirtualMachineServerConfigurationsManagement(serverConfig *sqlvirtualmachine.ServerConfigurationsManagementSettings, d *schema.ResourceData) []interface{} {
	if serverConfig == nil {
		return []interface{}{}
	}
	result := make(map[string]interface{})

	//additional feature
	result["is_r_services_enabled"] = serverConfig.AdditionalFeaturesServerConfigurations.IsRServicesEnabled

	//connectivity
	result["sql_connectivity_type"] = serverConfig.SQLConnectivityUpdateSettings.ConnectivityType

	if connectivityPort := serverConfig.SQLConnectivityUpdateSettings.Port; connectivityPort != nil {
		result["sql_connectivity_port"] = connectivityPort
	}
	if pwd, ok := d.GetOk("server_configuration.0.sql_connectivity_update_password"); ok {
		result["sql_connectivity_update_password"] = pwd
	}
	if userName, ok := d.GetOk("server_configuration.0.sql_connectivity_update_user_name"); ok {
		result["sql_connectivity_update_user_name"] = userName
	}

	return []interface{}{result}
}
