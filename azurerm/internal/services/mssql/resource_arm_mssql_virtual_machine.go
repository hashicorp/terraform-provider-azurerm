package mssql

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sqlvirtualmachine/mgmt/2017-03-01-preview/sqlvirtualmachine"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMsSqlVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMsSqlVirtualMachineCreateUpdate,
		Read:   resourceArmMsSqlVirtualMachineRead,
		Update: resourceArmMsSqlVirtualMachineCreateUpdate,
		Delete: resourceArmMsSqlVirtualMachineDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.MssqlVmID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"virtual_machine_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"sql_license_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sqlvirtualmachine.PAYG),
					string(sqlvirtualmachine.AHUB),
				}, false),
			},

			"sql_virtual_machine_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"auto_patching": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day_of_week": {
							Type:     schema.TypeString,
							Required: true,
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
							Required:     true,
							ValidateFunc: validation.IntBetween(30, 180),
						},

						"maintenance_window_starting_hour": {
							Type:         schema.TypeInt,
							Required:     true,
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
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
							//api will add updated credential name, and return "sqlvmName:name1,sqlvmName:name2"
							DiffSuppressFunc: mssqlVMCredentialNameDiffSuppressFunc,
						},

						"azure_key_vault_url": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							Sensitive:    true,
							ValidateFunc: validate.URLIsHTTPS,
						},

						"service_principal_name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"service_principal_secret": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
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
							Default:  string(sqlvirtualmachine.PRIVATE),
							ValidateFunc: validation.StringInSlice([]string{
								string(sqlvirtualmachine.LOCAL),
								string(sqlvirtualmachine.PRIVATE),
								string(sqlvirtualmachine.PUBLIC),
							}, false),
						},

						"sql_connectivity_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1433,
							ValidateFunc: validation.IntBetween(1024, 65535),
						},

						"sql_connectivity_update_password": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"sql_connectivity_update_username": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMsSqlVirtualMachineCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.SQLVirtualMachinesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vmId := d.Get("virtual_machine_id").(string)
	id, err := compute.ParseVirtualMachineID(vmId)
	if err != nil {
		return err
	}
	name := id.Name
	resourceGroupName := id.ResourceGroup

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, name, "*")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_mssql_virtual_machine", *existing.ID)
		}
	}

	// get location from vm
	vmclient := meta.(*clients.Client).Compute.VMClient
	respvm, err := vmclient.Get(ctx, resourceGroupName, name, "")
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure Virtual Machine %s: %+v", name, err)
	}

	if *respvm.Location == "" {
		return fmt.Errorf("Error location is empty from making Read request on Azure Virtual Machine %s: %+v", name, err)
	}

	properties := sqlvirtualmachine.Properties{
		VirtualMachineResourceID:               utils.String(d.Get("virtual_machine_id").(string)),
		SQLServerLicenseType:                   sqlvirtualmachine.SQLServerLicenseType(d.Get("sql_license_type").(string)),
		SQLManagement:                          sqlvirtualmachine.Full,
		AutoPatchingSettings:                   expandArmSqlVirtualMachineAutoPatchingSettings(d.Get("auto_patching").([]interface{})),
		KeyVaultCredentialSettings:             expandArmSqlVirtualMachineKeyVaultCredential(d.Get("key_vault_credential").([]interface{})),
		ServerConfigurationsManagementSettings: expandArmSqlVirtualMachineServerConfigurationsManagement(d.Get("server_configuration").([]interface{})),
	}

	if sqlVirtualMachineGroupResourceID, ok := d.GetOk("sql_virtual_machine_group_id"); ok {
		properties.SQLVirtualMachineGroupResourceID = utils.String(sqlVirtualMachineGroupResourceID.(string))
	}

	parameters := sqlvirtualmachine.SQLVirtualMachine{
		Location:   utils.String(*respvm.Location),
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

	resp, err := client.Get(ctx, resourceGroupName, name, "*")
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
	client := meta.(*clients.Client).MSSQL.SQLVirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MssqlVmID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	name := id.Name

	resp, err := client.Get(ctx, resourceGroupName, name, "*")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Sql Virtual Machine %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	if properties := resp.Properties; properties != nil {
		d.Set("virtual_machine_id", properties.VirtualMachineResourceID)
		d.Set("sql_license_type", string(properties.SQLServerLicenseType))
		flattenedAutoPatching := flattenArmSqlVirtualMachineAutoPatching(properties.AutoPatchingSettings)
		if err := d.Set("auto_patching", flattenedAutoPatching); err != nil {
			return fmt.Errorf("Error setting `auto_patching`: %+v", err)
		}
		flattenedKeyVaultCredential := flattenArmSqlVirtualMachineKeyVaultCredential(properties.KeyVaultCredentialSettings, d)
		if err := d.Set("key_vault_credential", flattenedKeyVaultCredential); err != nil {
			return fmt.Errorf("Error setting `key_vault_credential`: %+v", err)
		}
		flattenedServerConfiguration := flattenArmSqlVirtualMachineServerConfigurationsManagement(properties.ServerConfigurationsManagementSettings, d)
		if err := d.Set("server_configuration", flattenedServerConfiguration); err != nil {
			return fmt.Errorf("Error setting `server_configuration`: %+v", err)
		}
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMsSqlVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.SQLVirtualMachinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MssqlVmID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandArmSqlVirtualMachineAutoPatchingSettings(input []interface{}) *sqlvirtualmachine.AutoPatchingSettings {
	if len(input) == 0 {
		return nil
	}
	autoPatchingSetting := input[0].(map[string]interface{})

	return &sqlvirtualmachine.AutoPatchingSettings{
		Enable:                        utils.Bool(true),
		MaintenanceWindowDuration:     utils.Int32(int32(autoPatchingSetting["maintenance_window_duration_in_minutes"].(int))),
		MaintenanceWindowStartingHour: utils.Int32(int32(autoPatchingSetting["maintenance_window_starting_hour"].(int))),
		DayOfWeek:                     sqlvirtualmachine.DayOfWeek(autoPatchingSetting["day_of_week"].(string)),
	}
}

func flattenArmSqlVirtualMachineAutoPatching(autoPatching *sqlvirtualmachine.AutoPatchingSettings) []interface{} {
	if autoPatching == nil || !*autoPatching.Enable {
		return []interface{}{}
	}

	var startHour int32
	if autoPatching.MaintenanceWindowStartingHour != nil {
		startHour = *autoPatching.MaintenanceWindowStartingHour
	}

	var duration int32
	if autoPatching.MaintenanceWindowDuration != nil {
		duration = *autoPatching.MaintenanceWindowDuration
	}

	return []interface{}{
		map[string]interface{}{
			"day_of_week":                            string(autoPatching.DayOfWeek),
			"maintenance_window_starting_hour":       startHour,
			"maintenance_window_duration_in_minutes": duration,
		},
	}
}

func expandArmSqlVirtualMachineKeyVaultCredential(input []interface{}) *sqlvirtualmachine.KeyVaultCredentialSettings {
	if len(input) == 0 {
		return nil
	}
	keyVaultCredentialSetting := input[0].(map[string]interface{})

	return &sqlvirtualmachine.KeyVaultCredentialSettings{
		Enable:                 utils.Bool(true),
		CredentialName:         utils.String(keyVaultCredentialSetting["name"].(string)),
		AzureKeyVaultURL:       utils.String(keyVaultCredentialSetting["azure_key_vault_url"].(string)),
		ServicePrincipalName:   utils.String(keyVaultCredentialSetting["service_principal_name"].(string)),
		ServicePrincipalSecret: utils.String(keyVaultCredentialSetting["service_principal_secret"].(string)),
	}
}

func flattenArmSqlVirtualMachineKeyVaultCredential(keyVault *sqlvirtualmachine.KeyVaultCredentialSettings, d *schema.ResourceData) []interface{} {
	if keyVault == nil || !*keyVault.Enable {
		return []interface{}{}
	}

	name := ""
	if keyVault.CredentialName != nil {
		name = *keyVault.CredentialName
	}

	keyVaultUrl := ""
	if v, ok := d.GetOk("key_vault_credential.0.azure_key_vault_url"); ok {
		keyVaultUrl = v.(string)
	}

	servicePrincipalName := ""
	if v, ok := d.GetOk("key_vault_credential.0.service_principal_name"); ok {
		servicePrincipalName = v.(string)
	}

	servicePrincipalSecret := ""
	if v, ok := d.GetOk("key_vault_credential.0.service_principal_secret"); ok {
		servicePrincipalSecret = v.(string)
	}

	return []interface{}{
		map[string]interface{}{
			"name":                     name,
			"azure_key_vault_url":      keyVaultUrl,
			"service_principal_name":   servicePrincipalName,
			"service_principal_secret": servicePrincipalSecret,
		},
	}
}

func expandArmSqlVirtualMachineServerConfigurationsManagement(input []interface{}) *sqlvirtualmachine.ServerConfigurationsManagementSettings {
	if len(input) == 0 {
		return nil
	}
	serverConfigMM := input[0].(map[string]interface{})

	result := sqlvirtualmachine.ServerConfigurationsManagementSettings{
		SQLConnectivityUpdateSettings: &sqlvirtualmachine.SQLConnectivityUpdateSettings{
			ConnectivityType: sqlvirtualmachine.ConnectivityType(serverConfigMM["sql_connectivity_type"].(string)),
			Port:             utils.Int32(int32(serverConfigMM["sql_connectivity_port"].(int))),
		},
	}
	if userName, ok := serverConfigMM["sql_connectivity_update_username"]; ok {
		result.SQLConnectivityUpdateSettings.SQLAuthUpdateUserName = utils.String(userName.(string))
	}
	if pwd, ok := serverConfigMM["sql_connectivity_update_password"]; ok {
		result.SQLConnectivityUpdateSettings.SQLAuthUpdatePassword = utils.String(pwd.(string))
	}

	//additional feature
	if v, ok := serverConfigMM["is_r_services_enabled"]; ok {
		result.AdditionalFeaturesServerConfigurations = &sqlvirtualmachine.AdditionalFeaturesServerConfigurations{IsRServicesEnabled: utils.Bool(v.(bool))}
	}
	return &result
}

func flattenArmSqlVirtualMachineServerConfigurationsManagement(serverConfig *sqlvirtualmachine.ServerConfigurationsManagementSettings, d *schema.ResourceData) []interface{} {
	// if the structure of sqlvirtualmachine.ServerConfigurationsManagementSettings changes, we should update this
	if serverConfig == nil || reflect.DeepEqual(*serverConfig, sqlvirtualmachine.ServerConfigurationsManagementSettings{
		SQLConnectivityUpdateSettings:          &sqlvirtualmachine.SQLConnectivityUpdateSettings{},
		SQLWorkloadTypeUpdateSettings:          &sqlvirtualmachine.SQLWorkloadTypeUpdateSettings{},
		SQLStorageUpdateSettings:               &sqlvirtualmachine.SQLStorageUpdateSettings{},
		AdditionalFeaturesServerConfigurations: &sqlvirtualmachine.AdditionalFeaturesServerConfigurations{},
	}) {
		return []interface{}{}
	}

	//additional feature
	var isRServiceEnabled bool
	if serverConfig.AdditionalFeaturesServerConfigurations.IsRServicesEnabled != nil {
		isRServiceEnabled = *serverConfig.AdditionalFeaturesServerConfigurations.IsRServicesEnabled
	}
	//connectivity
	connectivityType := serverConfig.SQLConnectivityUpdateSettings.ConnectivityType

	var connectivityPort int32
	if serverConfig.SQLConnectivityUpdateSettings.Port != nil {
		connectivityPort = *serverConfig.SQLConnectivityUpdateSettings.Port
	}
	pwd := ""
	if v, ok := d.GetOk("server_configuration.0.sql_connectivity_update_password"); ok {
		pwd = v.(string)
	}
	userName := ""
	if v, ok := d.GetOk("server_configuration.0.sql_connectivity_update_username"); ok {
		userName = v.(string)
	}

	return []interface{}{
		map[string]interface{}{
			"is_r_services_enabled":            isRServiceEnabled,
			"sql_connectivity_type":            connectivityType,
			"sql_connectivity_port":            connectivityPort,
			"sql_connectivity_update_password": pwd,
			"sql_connectivity_update_username": userName,
		},
	}
}

func mssqlVMCredentialNameDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	oldNamelist := strings.Split(old, ",")
	for _, n := range oldNamelist {
		cur := strings.Split(n, ":")
		if len(cur) > 1 && cur[1] == new {
			return true
		}
	}
	return false
}
