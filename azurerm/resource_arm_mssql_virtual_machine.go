package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/sqlvirtualmachine/mgmt/2017-03-01-preview/sqlvirtualmachine"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	//"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssqlvm"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"log"
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

		Schema: map[string]*schema.Schema{
			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"virtual_machine_resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					_, err := azure.ParseAzureResourceID(v)
					errs = append(errs, err)
					return
				},
			},

			"sql_server_license_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sqlvirtualmachine.PAYG),
					string(sqlvirtualmachine.AHUB),
				}, false),
			},

			"sql_image_sku": {
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
				Type:     schema.TypeString,
				Optional: true,
			},

			"auto_patching_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"maintenance_window_duration": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"maintenance_window_starting_hour": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},

			"key_vault_credential_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"azure_key_vault_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"credential_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"service_principal_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"service_principal_secret": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},

			"server_configurations_management_settings": {
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
							Type:     schema.TypeInt,
							Optional: true,
						},
						"sql_connectivity_auth_update_password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"sql_connectivity_auth_update_user_name": {
							Type:     schema.TypeString,
							Optional: true,
						},

					},
				},
			},

			"storage_configuration_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"sql_data_default_file_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sql_data_luns": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"sql_log_default_file_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sql_log_luns": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"sql_temp_db_default_file_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sql_temp_db_luns": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"storage_workload_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(sqlvirtualmachine.GENERAL),
								string(sqlvirtualmachine.OLTP),
								string(sqlvirtualmachine.DW),
							}, false),
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMsSqlVirtualMachineCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).MSSQLVM.SQLVirtualMachinesClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	resourceGroupName := d.Get("resource_group_name").(string)

	id, _ := azure.ParseAzureResourceID(d.Get("virtual_machine_resource_id").(string))
	name := id.Path["virtualMachines"]

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
	sqlServerLicenseType := d.Get("sql_server_license_type").(string)
	virtualMachineResourceName := d.Get("virtual_machine_resource_id").(string)

	properties := sqlvirtualmachine.Properties{
		SQLServerLicenseType:     sqlvirtualmachine.SQLServerLicenseType(sqlServerLicenseType),
		VirtualMachineResourceID: &virtualMachineResourceName,
	}

	if sqlVirtualMachineGroupResourceID, ok := d.GetOk("sql_virtual_machine_group_resource_id"); ok {
		SQLVirtualMachineGroupResourceID := sqlVirtualMachineGroupResourceID.(string)
		properties.SQLVirtualMachineGroupResourceID = &SQLVirtualMachineGroupResourceID
	}


	if sqlImageSku, ok := d.GetOk("sql_image_sku"); ok {
		SQLImageSku := sqlvirtualmachine.SQLImageSku(sqlImageSku.(string))
		properties.SQLImageSku = SQLImageSku
	}

	if _, ok := d.GetOk("auto_patching_settings"); ok {
		properties.AutoPatchingSettings = expandArmSqlVirtualMachineAutoPatchingSettings(d)
	}

	if _, ok := d.GetOk("key_vault_credential_settings"); ok {
		properties.KeyVaultCredentialSettings = expandArmSqlVirtualMachineKeyVaultCredentialSettings(d)
	}

	if _, ok := d.GetOk("server_configurations_management_settings"); ok {
		properties.ServerConfigurationsManagementSettings = expandArmSqlVirtualMachineServerConfigurationsManagementSettings(d)
	}

	if _, ok := d.GetOk("storage_configuration_settings"); ok {
		properties.StorageConfigurationSettings = expandArmSqlVirtualMachineStorageConfigurationSettings(d)
	}

	Tags := d.Get("tags").(map[string]interface{})

	parameters := sqlvirtualmachine.SQLVirtualMachine{
		Location:   utils.String(location),
		Properties: &properties,
		Tags: tags.Expand(Tags),
	}


	future, err := client.CreateOrUpdate(ctx, resourceGroupName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	resp, err := client.Get(ctx, resourceGroupName, name, "")
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
	client := meta.(*ArmClient).MSSQLVM.SQLVirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	name := id.Path["sqlVirtualMachines"]

	resp, err := client.Get(ctx, resourceGroupName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Sql Virtual Machine %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	d.Set("resource_group_name", resourceGroupName)
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

func resourceArmMsSqlVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).MSSQLVM.SQLVirtualMachinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
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

func expandArmSqlVirtualMachineAutoBackupSettings(d *schema.ResourceData) *sqlvirtualmachine.AutoBackupSettings {
	autoBackupSettings := d.Get("auto_backup_settings").([]interface{})
	autoBackupSetting := autoBackupSettings[0].(map[string]interface{})

	result := sqlvirtualmachine.AutoBackupSettings{}

	if backupScheduleType, ok := autoBackupSetting["backup_schedule_type"]; ok {
		result.BackupScheduleType = sqlvirtualmachine.BackupScheduleType(backupScheduleType.(string))
	}
	if backupSystemDbs, ok := autoBackupSetting["backup_system_dbs"]; ok {
		BackupSystemDbs := backupSystemDbs.(bool)
		result.BackupSystemDbs = &BackupSystemDbs
	}
	if enable, ok := autoBackupSetting["enable"]; ok {
		Enable := enable.(bool)
		result.Enable = &Enable
	}
	if enableEncryption, ok := autoBackupSetting["enable_encryption"]; ok {
		EnableEncryption := enableEncryption.(bool)
		result.EnableEncryption = &EnableEncryption
	}
	if fullBackupFrequency, ok := autoBackupSetting["full_backup_frequency"]; ok {
		result.FullBackupFrequency = sqlvirtualmachine.FullBackupFrequencyType(fullBackupFrequency.(string))
	}
	if fullBackupStartTime, ok := autoBackupSetting["full_backup_start_time"]; ok {
		FullBackupStartTime := int32(fullBackupStartTime.(int))
		result.FullBackupStartTime = &FullBackupStartTime
	}
	if fullBackupWindowHours, ok := autoBackupSetting["full_backup_window_hours"]; ok {
		FullBackupWindowHours := int32(fullBackupWindowHours.(int))
		result.FullBackupWindowHours = &FullBackupWindowHours
	}
	if logBackupFrequency, ok := autoBackupSetting["log_backup_frequency"]; ok {
		LogBackupFrequency := int32(logBackupFrequency.(int))
		result.LogBackupFrequency = &LogBackupFrequency
	}
	if password, ok := autoBackupSetting["password"]; ok {
		Password := password.(string)
		result.Password = &Password
	}
	if retentionPeriod, ok := autoBackupSetting["retention_period"]; ok {
		RetentionPeriod := int32(retentionPeriod.(int))
		result.RetentionPeriod = &RetentionPeriod
	}
	if storageAccessKey, ok := autoBackupSetting["storage_access_key"]; ok {
		StorageAccessKey := storageAccessKey.(string)
		result.StorageAccessKey = &StorageAccessKey
	}
	if storageAccountURL, ok := autoBackupSetting["storage_account_url"]; ok {
		StorageAccountURL := storageAccountURL.(string)
		result.StorageAccountURL = &StorageAccountURL
	}
	return &result
}

func expandArmSqlVirtualMachineAutoPatchingSettings(d *schema.ResourceData) *sqlvirtualmachine.AutoPatchingSettings {
	autoPatchingSettings := d.Get("auto_patching_settings").([]interface{})
	autoPatchingSetting := autoPatchingSettings[0].(map[string]interface{})
	result := sqlvirtualmachine.AutoPatchingSettings{}

	if enable, ok := autoPatchingSetting["enable"]; ok {
		Enable := enable.(bool)
		result.Enable = &Enable
	}
	if maintenanceWindowDuration, ok := autoPatchingSetting["maintenance_window_duration"]; ok {
		MaintenanceWindowDuration := int32(maintenanceWindowDuration.(int))
		result.MaintenanceWindowDuration = &MaintenanceWindowDuration
	}
	if maintenanceWindowStartingHour, ok := autoPatchingSetting["maintenance_window_starting_hour"]; ok {
		MaintenanceWindowStartingHour := int32(maintenanceWindowStartingHour.(int))
		result.MaintenanceWindowStartingHour = &MaintenanceWindowStartingHour
	}
	if dayOfWeek, ok := autoPatchingSetting["day_of_week"]; ok {
		result.DayOfWeek = sqlvirtualmachine.DayOfWeek(dayOfWeek.(string))
	}
	return &result
}

func expandArmSqlVirtualMachineKeyVaultCredentialSettings(d *schema.ResourceData) *sqlvirtualmachine.KeyVaultCredentialSettings {
	keyVaultCredentialSettings := d.Get("key_vault_credential_settings").([]interface{})
	keyVaultCredentialSetting := keyVaultCredentialSettings[0].(map[string]interface{})
	result := sqlvirtualmachine.KeyVaultCredentialSettings{}

	if azureKeyVaultURL, ok := keyVaultCredentialSetting["azure_key_vault_url"]; ok {
		AzureKeyVaultURL := azureKeyVaultURL.(string)
		result.AzureKeyVaultURL = &AzureKeyVaultURL
	}
	if credentialName, ok := keyVaultCredentialSetting["credential_name"]; ok {
		CredentialName := credentialName.(string)
		result.CredentialName = &CredentialName
	}
	if servicePrincipalName, ok := keyVaultCredentialSetting["service_principal_name"]; ok {
		ServicePrincipalName := servicePrincipalName.(string)
		result.ServicePrincipalName = &ServicePrincipalName
	}
	if servicePrincipalSecret, ok := keyVaultCredentialSetting["service_principal_secret"]; ok {
		ServicePrincipalSecret := servicePrincipalSecret.(string)
		result.ServicePrincipalSecret = &ServicePrincipalSecret
	}
	if enable, ok := keyVaultCredentialSetting["enable"]; ok {
		Enable := enable.(bool)
		result.Enable = &Enable
	}

	return &result
}

func expandArmSqlVirtualMachineServerConfigurationsManagementSettings(d *schema.ResourceData) *sqlvirtualmachine.ServerConfigurationsManagementSettings {
	serverConfigMMs := d.Get("server_configurations_management_settings").([]interface{})
	serverConfigMM := serverConfigMMs[0].(map[string]interface{})

	result := sqlvirtualmachine.ServerConfigurationsManagementSettings{}
	sqlConnectivityUpdateSettings := sqlvirtualmachine.SQLConnectivityUpdateSettings{}
	//additional feature
	if isRServicesEnabled, ok := serverConfigMM["is_r_services_enabled"]; ok {
		IsRServicesEnabled := isRServicesEnabled.(bool)
		result.AdditionalFeaturesServerConfigurations = &sqlvirtualmachine.AdditionalFeaturesServerConfigurations{IsRServicesEnabled: &IsRServicesEnabled}
	}
	//connectivity
	if connectivityType, ok := serverConfigMM["sql_connectivity_type"]; ok {
		sqlConnectivityUpdateSettings.ConnectivityType = sqlvirtualmachine.ConnectivityType(connectivityType.(string))
	}
	if connectivityPort, ok := serverConfigMM["sql_connectivity_port"]; ok {
		ConnectivityPort := int32(connectivityPort.(int))
		sqlConnectivityUpdateSettings.Port = &ConnectivityPort
	}
	if sqlAuthUpdatePassword, ok := serverConfigMM["sql_connectivity_auth_update_password"]; ok {
		SQLAuthUpdatePassword := sqlAuthUpdatePassword.(string)
		sqlConnectivityUpdateSettings.SQLAuthUpdatePassword = &SQLAuthUpdatePassword
	}
	if sqlAuthUpdateUserName, ok := serverConfigMM["sql_connectivity_auth_update_user_name"]; ok {
		SQLAuthUpdateUserName := sqlAuthUpdateUserName.(string)
		sqlConnectivityUpdateSettings.SQLAuthUpdateUserName = &SQLAuthUpdateUserName
	}
	result.SQLConnectivityUpdateSettings = &sqlConnectivityUpdateSettings

	return &result
}

func expandArmSqlVirtualMachineStorageConfigurationSettings(d *schema.ResourceData) *sqlvirtualmachine.StorageConfigurationSettings {
	storageConfigs := d.Get("storage_configuration_settings").([]interface{})
	storageConfig := storageConfigs[0].(map[string]interface{})
	result := sqlvirtualmachine.StorageConfigurationSettings{}
	sqlDataSetting := sqlvirtualmachine.SQLStorageSettings{}
	sqlLogSetting := sqlvirtualmachine.SQLStorageSettings{}
	sqlTempDbSetting := sqlvirtualmachine.SQLStorageSettings{}

	if storageWorkloadType, ok := storageConfig["storage_workload_type"]; ok {
		result.StorageWorkloadType = sqlvirtualmachine.StorageWorkloadType(storageWorkloadType.(string))
	}
	//data setting
	if defaultFilePath, ok := storageConfig["sql_data_default_file_path"]; ok {
		DefaultFilePath := defaultFilePath.(string)
		sqlDataSetting.DefaultFilePath = &DefaultFilePath
	}
	if luns, ok := storageConfig["sql_data_luns"].(*schema.Set); ok {
		var Luns []int32
		for _, v := range luns.List() {
			n := int32(v.(int))
			Luns = append(Luns, n)
		}
		sqlDataSetting.Luns = &Luns
	}
	result.SQLDataSettings = &sqlDataSetting
	//log setting
	if defaultFilePath, ok := storageConfig["sql_log_default_file_path"]; ok {
		DefaultFilePath := defaultFilePath.(string)
		sqlLogSetting.DefaultFilePath = &DefaultFilePath
	}
	if luns, ok := storageConfig["sql_log_luns"].(*schema.Set); ok {
		var Luns []int32
		for _, v := range luns.List() {
			n := int32(v.(int))
			Luns = append(Luns, n)
		}
		sqlLogSetting.Luns = &Luns
	}
	result.SQLLogSettings = &sqlLogSetting
	//temp db setting
	if defaultFilePath, ok := storageConfig["sql_temp_db_default_file_path"]; ok {
		DefaultFilePath := defaultFilePath.(string)
		sqlTempDbSetting.DefaultFilePath = &DefaultFilePath
	}
	if luns, ok := storageConfig["sql_temp_db_luns"].(*schema.Set); ok {
		var Luns []int32
		for _, v := range luns.List() {
			n := int32(v.(int))
			Luns = append(Luns, n)
		}
		sqlTempDbSetting.Luns = &Luns
	}
	result.SQLTempDbSettings = &sqlTempDbSetting

	return &result
}


func flattenArmSqlVirtualMachineAutoPatchingSettings(autoPatchingSettings *sqlvirtualmachine.AutoPatchingSettings) []interface{} {
	if autoPatchingSettings == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["day_of_week"] = string(autoPatchingSettings.DayOfWeek)
	if enable := autoPatchingSettings.Enable; enable != nil {
		result["enable"] = *enable
	}
	if maintenanceWindowDuration := autoPatchingSettings.MaintenanceWindowDuration; maintenanceWindowDuration != nil {
		result["maintenance_window_duration"] = int(*maintenanceWindowDuration)
	}
	if maintenanceWindowStartingHour := autoPatchingSettings.MaintenanceWindowStartingHour; maintenanceWindowStartingHour != nil {
		result["maintenance_window_starting_hour"] = int(*maintenanceWindowStartingHour)
	}

	return []interface{}{result}
}

func flattenArmSqlVirtualMachineKeyVaultCredentialSettings(keyVaultCredentialSettings *sqlvirtualmachine.KeyVaultCredentialSettings) []interface{} {
	if keyVaultCredentialSettings == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if azureKeyVaultUrl := keyVaultCredentialSettings.AzureKeyVaultURL; azureKeyVaultUrl != nil {
		result["azure_key_vault_url"] = *azureKeyVaultUrl
	}
	if credentialName := keyVaultCredentialSettings.CredentialName; credentialName != nil {
		result["credential_name"] = *credentialName
	}
	if enable := keyVaultCredentialSettings.Enable; enable != nil {
		result["enable"] = *enable
	}
	if servicePrincipalName := keyVaultCredentialSettings.ServicePrincipalName; servicePrincipalName != nil {
		result["service_principal_name"] = *servicePrincipalName
	}
	if servicePrincipalSecret := keyVaultCredentialSettings.ServicePrincipalSecret; servicePrincipalSecret != nil {
		result["service_principal_secret"] = *servicePrincipalSecret
	}

	return []interface{}{result}
}

func flattenArmSqlVirtualMachineServerConfigurationsManagementSettings(serverConfigurationsManagementSettings *sqlvirtualmachine.ServerConfigurationsManagementSettings) []interface{} {
	if serverConfigurationsManagementSettings == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if additionalFeaturesServerConfigurations := serverConfigurationsManagementSettings.AdditionalFeaturesServerConfigurations; additionalFeaturesServerConfigurations != nil {
		if isRServicesEnabled := additionalFeaturesServerConfigurations.IsRServicesEnabled; isRServicesEnabled != nil {
			result["is_r_services_enabled"] = *isRServicesEnabled
		}
	}
	if sqlConnectivityUpdateSettings := serverConfigurationsManagementSettings.SQLConnectivityUpdateSettings; sqlConnectivityUpdateSettings != nil {
		result["sql_connectivity_type"] = string(sqlConnectivityUpdateSettings.ConnectivityType)
		if Port := sqlConnectivityUpdateSettings.Port; Port != nil {
			result["sql_connectivity_port"] = *Port
		}
		if SQLAuthUpdatePassword := sqlConnectivityUpdateSettings.SQLAuthUpdatePassword; SQLAuthUpdatePassword != nil {
			result["sql_connectivity_auth_update_password"] = *SQLAuthUpdatePassword
		}
		if SQLAuthUpdateUserName := sqlConnectivityUpdateSettings.SQLAuthUpdateUserName; SQLAuthUpdateUserName != nil {
			result["sql_connectivity_auth_update_user_name"] = *SQLAuthUpdateUserName
		}
	}

	return []interface{}{result}
}

func flattenArmSqlVirtualMachineStorageConfigurationSettings(storageConfigurationSettings *sqlvirtualmachine.StorageConfigurationSettings) []interface{} {
	if storageConfigurationSettings == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["storage_workload_type"] = string(storageConfigurationSettings.StorageWorkloadType)
	if sqlDataSettings := storageConfigurationSettings.SQLDataSettings; sqlDataSettings != nil {
		if DefaultFilePath := sqlDataSettings.DefaultFilePath; DefaultFilePath != nil {
			result["sql_data_default_file_path"] = *DefaultFilePath
		}
		if Luns := sqlDataSettings.Luns; Luns != nil {
			result["sql_data_luns"] = set.FromInt32Slice(*Luns)
		}
	}
	if sqlLogSettings := storageConfigurationSettings.SQLLogSettings; sqlLogSettings != nil {
		if DefaultFilePath := sqlLogSettings.DefaultFilePath; DefaultFilePath != nil {
			result["sql_log_default_file_path"] = *DefaultFilePath
		}
		if Luns := sqlLogSettings.Luns; Luns != nil {
			result["sql_log_luns"] = set.FromInt32Slice(*Luns)
		}
	}
	if sqlTempDbSettings := storageConfigurationSettings.SQLTempDbSettings; sqlTempDbSettings != nil {
		if DefaultFilePath := sqlTempDbSettings.DefaultFilePath; DefaultFilePath != nil {
			result["sql_temp_db_default_file_path"] = *DefaultFilePath
		}
		if Luns := sqlTempDbSettings.Luns; Luns != nil {
			result["sql_temp_db_luns"] = set.FromInt32Slice(*Luns)
		}
	}

	return []interface{}{result}
}
