package mssqlvm

import (
	"fmt"
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
				Computed: true,
				MaxItems: 1,
				ConfigMode: schema.SchemaConfigModeAttr,
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
						"maintenance_window_duration_in_minutes": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(30, 180),
						},
						"maintenance_window_starting_hour": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},

			"key_vault_credential": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"azure_key_vault_url": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.URLIsHTTPS,
						},
						"credential_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"service_principal_name": {
							Type:         schema.TypeString,
							Optional:     true,
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
				Computed: true,
				MaxItems: 1,
				ConfigMode: schema.SchemaConfigModeAttr,
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
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"storage_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sql_data_default_file_path": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"sql_data_luns": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntAtLeast(0),
							},
						},
						"sql_log_default_file_path": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"sql_log_luns": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntAtLeast(0),
							},
						},
						"sql_temp_db_default_file_path": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"sql_temp_db_luns": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntAtLeast(0),
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
	sqlServerLicenseType := d.Get("sql_license_type").(string)
	virtualMachineResourceName := d.Get("virtual_machine_resource_id").(string)

	properties := sqlvirtualmachine.Properties{
		SQLServerLicenseType:                   sqlvirtualmachine.SQLServerLicenseType(sqlServerLicenseType),
		VirtualMachineResourceID:               &virtualMachineResourceName,
		AutoPatchingSettings:                   expandArmSqlVirtualMachineAutoPatchingSettings(d),
		KeyVaultCredentialSettings:             expandArmSqlVirtualMachineKeyVaultCredential(d),
		ServerConfigurationsManagementSettings: expandArmSqlVirtualMachineServerConfigurationsManagement(d),
		StorageConfigurationSettings:           expandArmSqlVirtualMachineStorageConfiguration(d),
	}

	if sqlVirtualMachineGroupResourceID, ok := d.GetOk("sql_virtual_machine_group_resource_id"); ok {
		SQLVirtualMachineGroupResourceID := sqlVirtualMachineGroupResourceID.(string)
		properties.SQLVirtualMachineGroupResourceID = &SQLVirtualMachineGroupResourceID
	}

	if sqlImageSku, ok := d.GetOk("sql_sku"); ok {
		SQLImageSku := sqlvirtualmachine.SQLImageSku(sqlImageSku.(string))
		properties.SQLImageSku = SQLImageSku
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
	client := meta.(*clients.Client).MSSQLVM.SQLVirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroupName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if properties := resp.Properties; properties != nil {
		d.Set("sql_sku", string(properties.SQLImageSku))
		d.Set("sql_license_type", string(properties.SQLServerLicenseType))
		d.Set("virtual_machine_resource_id", properties.VirtualMachineResourceID)
	}
	if v := flattenArmSqlVirtualMachineAutoPatching(d); v != nil {
		if err := d.Set("auto_patching", v);err!=nil{
			return fmt.Errorf("Error setting `auto_patching`: %+v", err)
		}
	}
	if v := flattenArmSqlVirtualMachineKeyVaultCredential(d); v != nil {
		if err :=d.Set("key_vault_credential", v);err!=nil{
			return fmt.Errorf("Error setting `key_vault_credential`: %+v", err)
		}
	}
	if v := flattenArmSqlVirtualMachineServerConfigurationsManagement(d); v != nil {
		if err:=d.Set("server_configuration", v); err!=nil{
			return fmt.Errorf("Error setting `server_configuration`: %+v", err)
		}
	}
	if v := flattenArmSqlVirtualMachineStorageConfiguration(d); v != nil {
		if err:= d.Set("storage_configuration", v);err!=nil{
			return fmt.Errorf("Error setting `storage_configuration`: %+v", err)
		}
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

func expandArmSqlVirtualMachineAutoPatchingSettings(d *schema.ResourceData) *sqlvirtualmachine.AutoPatchingSettings {
	autoPatchingSettings := d.Get("auto_patching").([]interface{})
	if len(autoPatchingSettings) == 0 {
		return nil
	}
	autoPatchingSetting := autoPatchingSettings[0].(map[string]interface{})
	result := sqlvirtualmachine.AutoPatchingSettings{}

	if v, ok := autoPatchingSetting["enable"]; ok {
		enable := v.(bool)
		result.Enable = &enable
	}
	if v, ok := autoPatchingSetting["maintenance_window_duration_in_minutes"]; ok {
		maintenanceWindowDuration := int32(v.(int))
		result.MaintenanceWindowDuration = &maintenanceWindowDuration
	}
	if v, ok := autoPatchingSetting["maintenance_window_starting_hour"]; ok {
		maintenanceWindowStartingHour := int32(v.(int))
		result.MaintenanceWindowStartingHour = &maintenanceWindowStartingHour
	}
	if dayOfWeek, ok := autoPatchingSetting["day_of_week"]; ok {
		result.DayOfWeek = sqlvirtualmachine.DayOfWeek(dayOfWeek.(string))
	}
	return &result
}

func flattenArmSqlVirtualMachineAutoPatching(d *schema.ResourceData) []interface{} {
	// the API does not return it
	result := make(map[string]interface{})

	if enable, ok := d.GetOk("auto_patching.0.enable"); ok {
		result["enable"] = enable.(string)
	}
	if maintenanceWindowDuration, ok := d.GetOk("auto_patching.0.maintenance_window_duration_in_minutes"); ok {
		result["maintenance_window_duration_in_minutes"] = int32(maintenanceWindowDuration.(int))
	}
	if maintenanceWindowStartingHour, ok := d.GetOk("auto_patching.0.maintenance_window_starting_hour"); ok {
		result["maintenance_window_starting_hour"] = int32(maintenanceWindowStartingHour.(int))
	}
	if dayOfWeek, ok := d.GetOk("auto_patching.0.day_of_week"); ok {
		result["day_of_week"] = dayOfWeek.(string)
	}
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
		azureKeyVaultURL := v.(string)
		result.AzureKeyVaultURL = &azureKeyVaultURL
	}
	if v, ok := keyVaultCredentialSetting["credential_name"]; ok {
		credentialName := v.(string)
		result.CredentialName = &credentialName
	}
	if v, ok := keyVaultCredentialSetting["service_principal_name"]; ok {
		servicePrincipalName := v.(string)
		result.ServicePrincipalName = &servicePrincipalName
	}
	if v, ok := keyVaultCredentialSetting["service_principal_secret"]; ok {
		servicePrincipalSecret := v.(string)
		result.ServicePrincipalSecret = &servicePrincipalSecret
	}
	if v, ok := keyVaultCredentialSetting["enable"]; ok {
		enable := v.(bool)
		result.Enable = &enable
	}

	return &result
}

func flattenArmSqlVirtualMachineKeyVaultCredential(d *schema.ResourceData) []interface{} {
	// the API does not return it
	keyVaultCredentials := d.Get("key_vault_credential").([]interface{})
	if len(keyVaultCredentials) == 0 {
		return nil
	}
	keyVaultCredential := keyVaultCredentials[0].(map[string]interface{})
	result := make(map[string]interface{})

	if enable, ok := keyVaultCredential["enable"]; ok {
		result["enable"] = enable.(bool)
	}
	if azureKeyVaultURL, ok := keyVaultCredential["azure_key_vault_url"]; ok {
		result["azure_key_vault_url"] = azureKeyVaultURL.(string)
	}
	if credentialName, ok := keyVaultCredential["credential_name"]; ok {
		result["credential_name"] = credentialName.(string)
	}
	if servicePrincipalName, ok := keyVaultCredential["service_principal_name"]; ok {
		result["service_principal_name"] = servicePrincipalName.(string)
	}
	if servicePrincipalSecret, ok := keyVaultCredential["service_principal_secret"]; ok {
		result["service_principal_secret"] = servicePrincipalSecret.(string)
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
		isRServicesEnabled := v.(bool)
		result.AdditionalFeaturesServerConfigurations = &sqlvirtualmachine.AdditionalFeaturesServerConfigurations{IsRServicesEnabled: &isRServicesEnabled}
	}
	//connectivity
	if connectivityType, ok := serverConfigMM["sql_connectivity_type"]; ok {
		sqlConnectivityUpdateSettings.ConnectivityType = sqlvirtualmachine.ConnectivityType(connectivityType.(string))
	}
	if v, ok := serverConfigMM["sql_connectivity_port"]; ok {
		connectivityPort := int32(v.(int))
		sqlConnectivityUpdateSettings.Port = &connectivityPort
	}
	if v, ok := serverConfigMM["sql_connectivity_update_password"]; ok {
		sqlAuthUpdatePassword := v.(string)
		sqlConnectivityUpdateSettings.SQLAuthUpdatePassword = &sqlAuthUpdatePassword
	}
	if v, ok := serverConfigMM["sql_connectivity_update_user_name"]; ok {
		sqlAuthUpdateUserName := v.(string)
		sqlConnectivityUpdateSettings.SQLAuthUpdateUserName = &sqlAuthUpdateUserName
	}
	result.SQLConnectivityUpdateSettings = &sqlConnectivityUpdateSettings

	return &result
}

func flattenArmSqlVirtualMachineServerConfigurationsManagement(d *schema.ResourceData) []interface{} {
	// the API does not return it
	serverConfigMMs := d.Get("server_configuration").([]interface{})
	if len(serverConfigMMs) == 0 {
		return nil
	}
	serverConfigMM := serverConfigMMs[0].(map[string]interface{})

	result := make(map[string]interface{})

	//additional feature
	if isRServicesEnabled, ok := serverConfigMM["is_r_services_enabled"]; ok {
		result["is_r_services_enabled"] = isRServicesEnabled.(bool)
	}
	//connectivity
	if connectivityType, ok := serverConfigMM["sql_connectivity_type"]; ok {
		result["sql_connectivity_type"] = connectivityType.(string)
	}
	if connectivityPort, ok := serverConfigMM["sql_connectivity_port"]; ok {
		result["sql_connectivity_port"] = int32(connectivityPort.(int))
	}
	if sqlAuthUpdatePassword, ok := serverConfigMM["sql_connectivity_update_password"]; ok {
		result["sql_connectivity_update_password"] = sqlAuthUpdatePassword.(string)
	}
	if sqlAuthUpdateUserName, ok := serverConfigMM["sql_connectivity_update_user_name"]; ok {
		result["sql_connectivity_update_user_name"] = sqlAuthUpdateUserName.(string)
	}
	return []interface{}{result}
}

func expandArmSqlVirtualMachineStorageConfiguration(d *schema.ResourceData) *sqlvirtualmachine.StorageConfigurationSettings {
	storageConfigs := d.Get("storage_configuration").([]interface{})
	if len(storageConfigs) == 0 {
		return nil
	}
	storageConfig := storageConfigs[0].(map[string]interface{})
	result := sqlvirtualmachine.StorageConfigurationSettings{}
	sqlDataSetting := sqlvirtualmachine.SQLStorageSettings{}
	sqlLogSetting := sqlvirtualmachine.SQLStorageSettings{}
	sqlTempDbSetting := sqlvirtualmachine.SQLStorageSettings{}

	if storageWorkloadType, ok := storageConfig["storage_workload_type"]; ok {
		result.StorageWorkloadType = sqlvirtualmachine.StorageWorkloadType(storageWorkloadType.(string))
	}
	//data setting
	if v, ok := storageConfig["sql_data_default_file_path"]; ok {
		defaultFilePath := v.(string)
		sqlDataSetting.DefaultFilePath = &defaultFilePath
	}
	if s, ok := storageConfig["sql_data_luns"].(*schema.Set); ok {
		var luns []int32
		for _, v := range s.List() {
			n := int32(v.(int))
			luns = append(luns, n)
		}
		sqlDataSetting.Luns = &luns
	}
	result.SQLDataSettings = &sqlDataSetting
	//log setting
	if v, ok := storageConfig["sql_log_default_file_path"]; ok {
		defaultFilePath := v.(string)
		sqlLogSetting.DefaultFilePath = &defaultFilePath
	}
	if s, ok := storageConfig["sql_log_luns"].(*schema.Set); ok {
		var luns []int32
		for _, v := range s.List() {
			n := int32(v.(int))
			luns = append(luns, n)
		}
		sqlLogSetting.Luns = &luns
	}
	result.SQLLogSettings = &sqlLogSetting
	//temp db setting
	if v, ok := storageConfig["sql_temp_db_default_file_path"]; ok {
		defaultFilePath := v.(string)
		sqlTempDbSetting.DefaultFilePath = &defaultFilePath
	}
	if s, ok := storageConfig["sql_temp_db_luns"].(*schema.Set); ok {
		var luns []int32
		for _, v := range s.List() {
			n := int32(v.(int))
			luns = append(luns, n)
		}
		sqlTempDbSetting.Luns = &luns
	}
	result.SQLTempDbSettings = &sqlTempDbSetting

	return &result
}

func flattenArmSqlVirtualMachineStorageConfiguration(d *schema.ResourceData) []interface{} {
	// the API does not return it
	storageConfigs := d.Get("storage_configuration").([]interface{})
	if len(storageConfigs) == 0 {
		return nil
	}
	storageConfig := storageConfigs[0].(map[string]interface{})
	result := make(map[string]interface{})

	if storageWorkloadType, ok := storageConfig["storage_workload_type"]; ok {
		result["storage_workload_type"] = storageWorkloadType.(string)
	}
	if defaultFilePath, ok := storageConfig["sql_data_default_file_path"]; ok {
		result["sql_data_default_file_path"] = defaultFilePath.(string)
	}
	if luns, ok := storageConfig["sql_data_luns"]; ok {
		result["sql_data_luns"] = luns.(*schema.Set).List()
	}
	if defaultFilePath, ok := storageConfig["sql_log_default_file_path"]; ok {
		result["sql_log_default_file_path"] = defaultFilePath.(string)
	}
	if luns, ok := storageConfig["sql_log_luns"]; ok {
		result["sql_log_luns"] = luns.(*schema.Set).List()
	}
	if defaultFilePath, ok := storageConfig["sql_temp_db_default_file_path"]; ok {
		result["sql_temp_db_default_file_path"] = defaultFilePath.(string)
	}
	if luns, ok := storageConfig["sql_temp_db_luns"]; ok {
		result["sql_temp_db_luns"] = luns.(*schema.Set).List()
	}
	return []interface{}{result}
}
