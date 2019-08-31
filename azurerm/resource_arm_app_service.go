package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceCreate,
		Read:   resourceArmAppServiceRead,
		Update: resourceArmAppServiceUpdate,
		Delete: resourceArmAppServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAppServiceName,
			},

			"identity": azure.SchemaAppServiceIdentity(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"app_service_plan_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"site_config": azure.SchemaAppServiceSiteConfig(),

			"auth_settings": azure.SchemaAppServiceAuthSettings(),

			"logs": azure.SchemaAppServiceLogsConfig(),

			"backup": azure.SchemaAppServiceBackup(),

			"client_affinity_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"https_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"client_cert_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"app_settings": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},

			"storage_account": azure.SchemaAppServiceStorageAccounts(),

			"connection_string": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(web.APIHub),
								string(web.Custom),
								string(web.DocDb),
								string(web.EventHub),
								string(web.MySQL),
								string(web.NotificationHub),
								string(web.PostgreSQL),
								string(web.RedisCache),
								string(web.ServiceBus),
								string(web.SQLAzure),
								string(web.SQLServer),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
			},

			"tags": tags.Schema(),

			"site_credential": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"default_site_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"outbound_ip_addresses": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"possible_outbound_ip_addresses": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_control": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"branch": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmAppServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM App Service creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing App Service %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_service", *existing.ID)
		}
	}

	availabilityRequest := web.ResourceNameAvailabilityRequest{
		Name: utils.String(name),
		Type: web.CheckNameResourceTypesMicrosoftWebsites,
	}
	available, err := client.CheckNameAvailability(ctx, availabilityRequest)
	if err != nil {
		return fmt.Errorf("Error checking if the name %q was available: %+v", name, err)
	}

	if !*available.NameAvailable {
		return fmt.Errorf("The name %q used for the App Service needs to be globally unique and isn't available: %s", name, *available.Message)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	appServicePlanId := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	t := d.Get("tags").(map[string]interface{})

	siteConfig, err := azure.ExpandAppServiceSiteConfig(d.Get("site_config"))
	if err != nil {
		return fmt.Errorf("Error expanding `site_config` for App Service %q (Resource Group %q): %s", name, resGroup, err)
	}

	siteEnvelope := web.Site{
		Location: &location,
		Tags:     tags.Expand(t),
		SiteProperties: &web.SiteProperties{
			ServerFarmID: utils.String(appServicePlanId),
			Enabled:      utils.Bool(enabled),
			HTTPSOnly:    utils.Bool(httpsOnly),
			SiteConfig:   siteConfig,
		},
	}

	if _, ok := d.GetOk("identity"); ok {
		appServiceIdentity := azure.ExpandAppServiceIdentity(d)
		siteEnvelope.Identity = appServiceIdentity
	}

	if v, ok := d.GetOkExists("client_affinity_enabled"); ok {
		enabled := v.(bool)
		siteEnvelope.SiteProperties.ClientAffinityEnabled = utils.Bool(enabled)
	}

	if v, ok := d.GetOkExists("client_cert_enabled"); ok {
		certEnabled := v.(bool)
		siteEnvelope.SiteProperties.ClientCertEnabled = utils.Bool(certEnabled)
	}

	createFuture, err := client.CreateOrUpdate(ctx, resGroup, name, siteEnvelope)
	if err != nil {
		return fmt.Errorf("Error creating App Service %q (Resource Group %q): %s", name, resGroup, err)
	}

	err = createFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for App Service %q (Resource Group %q) to be created: %s", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving App Service %q (Resource Group %q): %s", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read App Service %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	authSettingsRaw := d.Get("auth_settings").([]interface{})
	authSettings := azure.ExpandAppServiceAuthSettings(authSettingsRaw)

	auth := web.SiteAuthSettings{
		ID:                         read.ID,
		SiteAuthSettingsProperties: &authSettings}

	if _, err := client.UpdateAuthSettings(ctx, resGroup, name, auth); err != nil {
		return fmt.Errorf("Error updating auth settings for App Service %q (Resource Group %q): %+s", name, resGroup, err)
	}

	logsConfig := azure.ExpandAppServiceLogs(d.Get("logs"))

	logs := web.SiteLogsConfig{
		ID:                       read.ID,
		SiteLogsConfigProperties: &logsConfig}

	if _, err := client.UpdateDiagnosticLogsConfig(ctx, resGroup, name, logs); err != nil {
		return fmt.Errorf("Error updating diagnostic logs config for App Service %q (Resource Group %q): %+s", name, resGroup, err)
	}

	backupRaw := d.Get("backup").([]interface{})
	if backup := azure.ExpandAppServiceBackup(backupRaw); backup != nil {
		_, err = client.UpdateBackupConfiguration(ctx, resGroup, name, *backup)
		if err != nil {
			return fmt.Errorf("Error updating Backup Settings for App Service %q (Resource Group %q): %s", name, resGroup, err)
		}
	}

	return resourceArmAppServiceUpdate(d, meta)
}

func resourceArmAppServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["sites"]

	location := azure.NormalizeLocation(d.Get("location").(string))

	appServicePlanId := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	t := d.Get("tags").(map[string]interface{})

	siteConfig, err := azure.ExpandAppServiceSiteConfig(d.Get("site_config"))
	if err != nil {
		return fmt.Errorf("Error expanding `site_config` for App Service %q (Resource Group %q): %s", name, resGroup, err)
	}

	siteEnvelope := web.Site{
		Location: &location,
		Tags:     tags.Expand(t),
		SiteProperties: &web.SiteProperties{
			ServerFarmID: utils.String(appServicePlanId),
			Enabled:      utils.Bool(enabled),
			HTTPSOnly:    utils.Bool(httpsOnly),
			SiteConfig:   siteConfig,
		},
	}

	if v, ok := d.GetOkExists("client_cert_enabled"); ok {
		certEnabled := v.(bool)
		siteEnvelope.SiteProperties.ClientCertEnabled = utils.Bool(certEnabled)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, siteEnvelope)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	if d.HasChange("site_config") {
		// update the main configuration
		siteConfig, err := azure.ExpandAppServiceSiteConfig(d.Get("site_config"))
		if err != nil {
			return fmt.Errorf("Error expanding `site_config` for App Service %q (Resource Group %q): %s", name, resGroup, err)
		}
		siteConfigResource := web.SiteConfigResource{
			SiteConfig: siteConfig,
		}

		if _, err := client.CreateOrUpdateConfiguration(ctx, resGroup, name, siteConfigResource); err != nil {
			return fmt.Errorf("Error updating Configuration for App Service %q: %+v", name, err)
		}
	}

	if d.HasChange("auth_settings") {
		authSettingsRaw := d.Get("auth_settings").([]interface{})
		authSettingsProperties := azure.ExpandAppServiceAuthSettings(authSettingsRaw)
		id := d.Id()
		authSettings := web.SiteAuthSettings{
			ID:                         &id,
			SiteAuthSettingsProperties: &authSettingsProperties,
		}

		if _, err := client.UpdateAuthSettings(ctx, resGroup, name, authSettings); err != nil {
			return fmt.Errorf("Error updating Authentication Settings for App Service %q: %+v", name, err)
		}
	}

	if d.HasChange("logs") {
		logs := azure.ExpandAppServiceLogs(d.Get("logs"))
		id := d.Id()
		logsResource := web.SiteLogsConfig{
			ID:                       &id,
			SiteLogsConfigProperties: &logs,
		}

		if _, err := client.UpdateDiagnosticLogsConfig(ctx, resGroup, name, logsResource); err != nil {
			return fmt.Errorf("Error updating Diagnostics Logs for App Service %q: %+v", name, err)
		}
	}

	if d.HasChange("backup") {
		backupRaw := d.Get("backup").([]interface{})
		if backup := azure.ExpandAppServiceBackup(backupRaw); backup != nil {
			_, err = client.UpdateBackupConfiguration(ctx, resGroup, name, *backup)
			if err != nil {
				return fmt.Errorf("Error updating Backup Settings for App Service %q (Resource Group %q): %s", name, resGroup, err)
			}
		} else {
			_, err = client.DeleteBackupConfiguration(ctx, resGroup, name)
			if err != nil {
				return fmt.Errorf("Error removing Backup Settings for App Service %q (Resource Group %q): %s", name, resGroup, err)
			}
		}
	}

	if d.HasChange("client_affinity_enabled") {

		affinity := d.Get("client_affinity_enabled").(bool)

		sitePatchResource := web.SitePatchResource{
			ID: utils.String(d.Id()),
			SitePatchResourceProperties: &web.SitePatchResourceProperties{
				ClientAffinityEnabled: &affinity,
			},
		}

		if _, err := client.Update(ctx, resGroup, name, sitePatchResource); err != nil {
			return fmt.Errorf("Error updating App Service ARR Affinity setting %q: %+v", name, err)
		}
	}

	if d.HasChange("app_settings") {
		// update the AppSettings
		appSettings := expandAppServiceAppSettings(d)
		settings := web.StringDictionary{
			Properties: appSettings,
		}

		if _, err := client.UpdateApplicationSettings(ctx, resGroup, name, settings); err != nil {
			return fmt.Errorf("Error updating Application Settings for App Service %q: %+v", name, err)
		}
	}

	if d.HasChange("storage_account") {
		storageAccounts := azure.ExpandAppServiceStorageAccounts(d)
		properties := web.AzureStoragePropertyDictionaryResource{
			Properties: storageAccounts,
		}

		if _, err := client.UpdateAzureStorageAccounts(ctx, resGroup, name, properties); err != nil {
			return fmt.Errorf("Error updating Storage Accounts for App Service %q: %+v", name, err)
		}
	}

	if d.HasChange("connection_string") {
		// update the ConnectionStrings
		connectionStrings := expandAppServiceConnectionStrings(d)
		properties := web.ConnectionStringDictionary{
			Properties: connectionStrings,
		}

		if _, err := client.UpdateConnectionStrings(ctx, resGroup, name, properties); err != nil {
			return fmt.Errorf("Error updating Connection Strings for App Service %q: %+v", name, err)
		}
	}

	if d.HasChange("identity") {
		site, err := client.Get(ctx, resGroup, name)
		if err != nil {
			return fmt.Errorf("Error getting configuration for App Service %q: %+v", name, err)
		}

		appServiceIdentity := azure.ExpandAppServiceIdentity(d)
		site.Identity = appServiceIdentity

		future, err := client.CreateOrUpdate(ctx, resGroup, name, site)

		if err != nil {
			return fmt.Errorf("Error updating Managed Service Identity for App Service %q: %+v", name, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error updating Managed Service Identity for App Service %q: %+v", name, err)
		}
	}

	return resourceArmAppServiceRead(d, meta)
}

func resourceArmAppServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["sites"]

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service %q (resource group %q) was not found - removing from state", name, resGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service %q: %+v", name, err)
	}

	configResp, err := client.GetConfiguration(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(configResp.Response) {
			log.Printf("[DEBUG] Configuration of App Service %q (resource group %q) was not found", name, resGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service Configuration %q: %+v", name, err)
	}

	authResp, err := client.GetAuthSettings(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving the AuthSettings for App Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	backupResp, err := client.GetBackupConfiguration(ctx, resGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(backupResp.Response) {
			return fmt.Errorf("Error retrieving the BackupConfiguration for App Service %q (Resource Group %q): %+v", name, resGroup, err)
		}
	}

	logsResp, err := client.GetDiagnosticLogsConfiguration(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving the DiagnosticsLogsConfiguration for App Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	appSettingsResp, err := client.ListApplicationSettings(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(appSettingsResp.Response) {
			log.Printf("[DEBUG] Application Settings of App Service %q (resource group %q) were not found", name, resGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service AppSettings %q: %+v", name, err)
	}

	storageAccountsResp, err := client.ListAzureStorageAccounts(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service Storage Accounts %q: %+v", name, err)
	}

	connectionStringsResp, err := client.ListConnectionStrings(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service ConnectionStrings %q: %+v", name, err)
	}

	scmResp, err := client.GetSourceControl(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service Source Control %q: %+v", name, err)
	}

	siteCredFuture, err := client.ListPublishingCredentials(ctx, resGroup, name)
	if err != nil {
		return err
	}
	err = siteCredFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return err
	}
	siteCredResp, err := siteCredFuture.Result(*client)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service Site Credential %q: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.SiteProperties; props != nil {
		d.Set("app_service_plan_id", props.ServerFarmID)
		d.Set("client_affinity_enabled", props.ClientAffinityEnabled)
		d.Set("enabled", props.Enabled)
		d.Set("https_only", props.HTTPSOnly)
		d.Set("client_cert_enabled", props.ClientCertEnabled)
		d.Set("default_site_hostname", props.DefaultHostName)
		d.Set("outbound_ip_addresses", props.OutboundIPAddresses)
		d.Set("possible_outbound_ip_addresses", props.PossibleOutboundIPAddresses)
	}

	appSettings := flattenAppServiceAppSettings(appSettingsResp.Properties)

	// remove DIAGNOSTICS* settings - Azure will sync these, so just maintain the logs block equivalents in the state
	delete(appSettings, "DIAGNOSTICS_AZUREBLOBCONTAINERSASURL")
	delete(appSettings, "DIAGNOSTICS_AZUREBLOBRETENTIONINDAYS")

	if err := d.Set("app_settings", appSettings); err != nil {
		return fmt.Errorf("Error setting `app_settings`: %s", err)
	}

	if err := d.Set("backup", azure.FlattenAppServiceBackup(backupResp.BackupRequestProperties)); err != nil {
		return fmt.Errorf("Error setting `backup`: %s", err)
	}

	if err := d.Set("storage_account", azure.FlattenAppServiceStorageAccounts(storageAccountsResp.Properties)); err != nil {
		return fmt.Errorf("Error setting `storage_account`: %s", err)
	}

	if err := d.Set("connection_string", flattenAppServiceConnectionStrings(connectionStringsResp.Properties)); err != nil {
		return fmt.Errorf("Error setting `connection_string`: %s", err)
	}

	siteConfig := azure.FlattenAppServiceSiteConfig(configResp.SiteConfig)
	if err := d.Set("site_config", siteConfig); err != nil {
		return err
	}

	authSettings := azure.FlattenAppServiceAuthSettings(authResp.SiteAuthSettingsProperties)
	if err := d.Set("auth_settings", authSettings); err != nil {
		return fmt.Errorf("Error setting `auth_settings`: %s", err)
	}

	logs := azure.FlattenAppServiceLogs(logsResp.SiteLogsConfigProperties)
	if err := d.Set("logs", logs); err != nil {
		return fmt.Errorf("Error setting `logs`: %s", err)
	}

	scm := flattenAppServiceSourceControl(scmResp.SiteSourceControlProperties)
	if err := d.Set("source_control", scm); err != nil {
		return fmt.Errorf("Error setting `source_control`: %s", err)
	}

	siteCred := flattenAppServiceSiteCredential(siteCredResp.UserProperties)
	if err := d.Set("site_credential", siteCred); err != nil {
		return fmt.Errorf("Error setting `site_credential`: %s", err)
	}

	identity := azure.FlattenAppServiceIdentity(resp.Identity)
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("Error setting `identity`: %s", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmAppServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["sites"]

	log.Printf("[DEBUG] Deleting App Service %q (resource group %q)", name, resGroup)

	deleteMetrics := true
	deleteEmptyServerFarm := false
	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Delete(ctx, resGroup, name, &deleteMetrics, &deleteEmptyServerFarm)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return err
		}
	}

	return nil
}

func flattenAppServiceSourceControl(input *web.SiteSourceControlProperties) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{})

	if input == nil {
		log.Printf("[DEBUG] SiteSourceControlProperties is nil")
		return results
	}

	if input.RepoURL != nil {
		result["repo_url"] = *input.RepoURL
	}
	if input.Branch != nil && *input.Branch != "" {
		result["branch"] = *input.Branch
	} else {
		result["branch"] = "master"
	}

	return append(results, result)
}

func expandAppServiceAppSettings(d *schema.ResourceData) map[string]*string {
	input := d.Get("app_settings").(map[string]interface{})
	output := make(map[string]*string, len(input))

	for k, v := range input {
		output[k] = utils.String(v.(string))
	}

	return output
}

func expandAppServiceConnectionStrings(d *schema.ResourceData) map[string]*web.ConnStringValueTypePair {
	input := d.Get("connection_string").(*schema.Set).List()
	output := make(map[string]*web.ConnStringValueTypePair, len(input))

	for _, v := range input {
		vals := v.(map[string]interface{})

		csName := vals["name"].(string)
		csType := vals["type"].(string)
		csValue := vals["value"].(string)

		output[csName] = &web.ConnStringValueTypePair{
			Value: utils.String(csValue),
			Type:  web.ConnectionStringType(csType),
		}
	}

	return output
}

func flattenAppServiceConnectionStrings(input map[string]*web.ConnStringValueTypePair) []interface{} {
	results := make([]interface{}, 0)

	for k, v := range input {
		result := make(map[string]interface{})
		result["name"] = k
		result["type"] = string(v.Type)
		if v.Value != nil {
			result["value"] = *v.Value
		}
		results = append(results, result)
	}

	return results
}

func flattenAppServiceAppSettings(input map[string]*string) map[string]string {
	output := make(map[string]string)
	for k, v := range input {
		output[k] = *v
	}

	return output
}

func validateAppServiceName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-]{1,60}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes and up to 60 characters in length", k))
	}

	return warnings, errors
}

func flattenAppServiceSiteCredential(input *web.UserProperties) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{})

	if input == nil {
		log.Printf("[DEBUG] UserProperties is nil")
		return results
	}

	if input.PublishingUserName != nil {
		result["username"] = *input.PublishingUserName
	}

	if input.PublishingPassword != nil {
		result["password"] = *input.PublishingPassword
	}

	return append(results, result)
}
