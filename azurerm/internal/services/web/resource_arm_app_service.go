package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceCreate,
		Read:   resourceArmAppServiceRead,
		Update: resourceArmAppServiceUpdate,
		Delete: resourceArmAppServiceDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AppServiceID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AppServiceName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"app_service_plan_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"app_settings": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"auth_settings": schemaAppServiceAuthSettings(),

			"backup": schemaAppServiceBackup(),

			"client_affinity_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"client_cert_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

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

						"value": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"identity": schemaAppServiceIdentity(),

			"https_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"logs": schemaAppServiceLogsConfig(),

			"site_config": schemaAppServiceSiteConfig(),

			"storage_account": schemaAppServiceStorageAccounts(),

			"source_control": schemaAppServiceSiteSourceControl(),

			"tags": tags.Schema(),

			"site_credential": {
				Type:     schema.TypeList,
				Computed: true,
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
		},
	}
}

func resourceArmAppServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	aspClient := meta.(*clients.Client).Web.AppServicePlansClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM App Service creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing App Service %q (Resource Group %q): %s", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_app_service", *existing.ID)
	}

	availabilityRequest := web.ResourceNameAvailabilityRequest{
		Name: utils.String(name),
		Type: web.CheckNameResourceTypesMicrosoftWebsites,
	}

	appServicePlanId := d.Get("app_service_plan_id").(string)
	aspID, err := ParseAppServicePlanID(appServicePlanId)
	if err != nil {
		return err
	}
	// Check if App Service Plan is part of ASE
	// If so, the name needs updating to <app name>.<ASE name>.appserviceenvironment.net and FQDN setting true for name availability check
	aspDetails, err := aspClient.Get(ctx, aspID.ResourceGroup, aspID.Name)
	if err != nil {
		return fmt.Errorf("App Service Environment %q (Resource Group %q) does not exist", aspID.Name, aspID.ResourceGroup)
	}
	if aspDetails.HostingEnvironmentProfile != nil {
		availabilityRequest.Name = utils.String(fmt.Sprintf("%s.%s.appserviceenvironment.net", name, aspID.Name))
		availabilityRequest.IsFqdn = utils.Bool(true)
	}
	available, err := client.CheckNameAvailability(ctx, availabilityRequest)
	if err != nil {
		return fmt.Errorf("Error checking if the name %q was available: %+v", name, err)
	}

	if !*available.NameAvailable {
		return fmt.Errorf("The name %q used for the App Service needs to be globally unique and isn't available: %s", name, *available.Message)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	enabled := d.Get("enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	t := d.Get("tags").(map[string]interface{})

	siteConfig, err := expandAppServiceSiteConfig(d.Get("site_config"))
	if err != nil {
		return fmt.Errorf("Error expanding `site_config` for App Service %q (Resource Group %q): %s", name, resourceGroup, err)
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
		appServiceIdentityRaw := d.Get("identity").([]interface{})
		appServiceIdentity := expandAppServiceIdentity(appServiceIdentityRaw)
		siteEnvelope.Identity = appServiceIdentity
	}

	siteEnvelope.SiteProperties.ClientAffinityEnabled = utils.Bool(d.Get("client_affinity_enabled").(bool))

	siteEnvelope.SiteProperties.ClientCertEnabled = utils.Bool(d.Get("client_cert_enabled").(bool))

	createFuture, err := client.CreateOrUpdate(ctx, resourceGroup, name, siteEnvelope)
	if err != nil {
		return fmt.Errorf("Error creating App Service %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	err = createFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for App Service %q (Resource Group %q) to be created: %s", name, resourceGroup, err)
	}

	if _, ok := d.GetOk("source_control"); ok {
		if siteConfig.ScmType != "" {
			return fmt.Errorf("cannot set source_control parameters when scm_type is set to %q", siteConfig.ScmType)
		}
		sourceControlProperties := expandAppServiceSiteSourceControl(d)
		sourceControl := &web.SiteSourceControl{}
		sourceControl.SiteSourceControlProperties = sourceControlProperties
		// TODO - Do we need to lock the app for updates?
		scFuture, err := client.CreateOrUpdateSourceControl(ctx, resourceGroup, name, *sourceControl)
		if err != nil {
			return fmt.Errorf("failed to create App Service Source Control for %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		err = scFuture.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("failed waiting for App Service Source Control configuration")
		}
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving App Service %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	if read.ID == nil || *read.ID == "" {
		return fmt.Errorf("Cannot read App Service %q (resource group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	authSettingsRaw := d.Get("auth_settings").([]interface{})
	authSettings := expandAppServiceAuthSettings(authSettingsRaw)

	auth := web.SiteAuthSettings{
		ID:                         read.ID,
		SiteAuthSettingsProperties: &authSettings,
	}

	if _, err := client.UpdateAuthSettings(ctx, resourceGroup, name, auth); err != nil {
		return fmt.Errorf("Error updating auth settings for App Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	logsConfig := expandAppServiceLogs(d.Get("logs"))

	logs := web.SiteLogsConfig{
		ID:                       read.ID,
		SiteLogsConfigProperties: &logsConfig,
	}

	if _, err := client.UpdateDiagnosticLogsConfig(ctx, resourceGroup, name, logs); err != nil {
		return fmt.Errorf("Error updating diagnostic logs config for App Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	backupRaw := d.Get("backup").([]interface{})
	if backup := expandAppServiceBackup(backupRaw); backup != nil {
		_, err = client.UpdateBackupConfiguration(ctx, resourceGroup, name, *backup)
		if err != nil {
			return fmt.Errorf("Error updating Backup Settings for App Service %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return resourceArmAppServiceUpdate(d, meta)
}

func resourceArmAppServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceID(d.Id())
	if err != nil {
		return err
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	appServicePlanId := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	t := d.Get("tags").(map[string]interface{})

	siteConfig, err := expandAppServiceSiteConfig(d.Get("site_config"))
	if err != nil {
		return fmt.Errorf("Error expanding `site_config` for App Service %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
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

	siteEnvelope.SiteProperties.ClientCertEnabled = utils.Bool(d.Get("client_cert_enabled").(bool))

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, siteEnvelope)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	// If `source_control` is defined, we need to set site_config.0.scm_type to "None" or we cannot update it
	_, hasSourceControl := d.GetOk("source_control.0.repo_url")

	scmType := web.ScmTypeNone

	if d.HasChange("site_config") || hasSourceControl {
		// update the main configuration
		siteConfig, err := expandAppServiceSiteConfig(d.Get("site_config"))
		if err != nil {
			return fmt.Errorf("Error expanding `site_config` for App Service %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
		}
		siteConfigResource := web.SiteConfigResource{
			SiteConfig: siteConfig,
		}

		scmType = siteConfig.ScmType
		// ScmType being set blocks the update of source_control in _most_ cases, ADO is an exception
		if hasSourceControl && scmType != web.ScmTypeVSTSRM {
			siteConfigResource.SiteConfig.ScmType = web.ScmTypeNone
		}

		if _, err := client.CreateOrUpdateConfiguration(ctx, id.ResourceGroup, id.Name, siteConfigResource); err != nil {
			return fmt.Errorf("Error updating Configuration for App Service %q: %+v", id.Name, err)
		}
	}

	// Don't send source_control changes for ADO controlled Apps
	if hasSourceControl && scmType != web.ScmTypeVSTSRM {
		sourceControlProperties := expandAppServiceSiteSourceControl(d)
		sourceControl := &web.SiteSourceControl{}
		sourceControl.SiteSourceControlProperties = sourceControlProperties
		scFuture, err := client.CreateOrUpdateSourceControl(ctx, id.ResourceGroup, id.Name, *sourceControl)
		if err != nil {
			return fmt.Errorf("failed to update App Service Source Control for %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		err = scFuture.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("failed waiting for App Service Source Control configuration: %+v", err)
		}

		sc, err := client.GetSourceControl(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("failed reading back App Service Source Control for %q", *sc.Name)
		}
	}

	if d.HasChange("auth_settings") {
		authSettingsRaw := d.Get("auth_settings").([]interface{})
		authSettingsProperties := expandAppServiceAuthSettings(authSettingsRaw)
		authSettings := web.SiteAuthSettings{
			ID:                         utils.String(d.Id()),
			SiteAuthSettingsProperties: &authSettingsProperties,
		}

		if _, err := client.UpdateAuthSettings(ctx, id.ResourceGroup, id.Name, authSettings); err != nil {
			return fmt.Errorf("Error updating Authentication Settings for App Service %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("backup") {
		backupRaw := d.Get("backup").([]interface{})
		if backup := expandAppServiceBackup(backupRaw); backup != nil {
			_, err = client.UpdateBackupConfiguration(ctx, id.ResourceGroup, id.Name, *backup)
			if err != nil {
				return fmt.Errorf("Error updating Backup Settings for App Service %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
			}
		} else {
			_, err = client.DeleteBackupConfiguration(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("Error removing Backup Settings for App Service %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
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

		if _, err := client.Update(ctx, id.ResourceGroup, id.Name, sitePatchResource); err != nil {
			return fmt.Errorf("Error updating App Service ARR Affinity setting %q: %+v", id.Name, err)
		}
	}

	// app settings updates have a side effect on logging settings. See the note below
	if d.HasChange("app_settings") {
		// update the AppSettings
		appSettings := expandAppServiceAppSettings(d)
		settings := web.StringDictionary{
			Properties: appSettings,
		}

		if _, err := client.UpdateApplicationSettings(ctx, id.ResourceGroup, id.Name, settings); err != nil {
			return fmt.Errorf("Error updating Application Settings for App Service %q: %+v", id.Name, err)
		}
	}

	// the logging configuration has a dependency on the app settings in Azure
	// e.g. configuring logging to blob storage will add the DIAGNOSTICS_AZUREBLOBCONTAINERSASURL
	// and DIAGNOSTICS_AZUREBLOBRETENTIONINDAYS app settings to the app service.
	// If the app settings are updated, also update the logging configuration if it exists, otherwise
	// updating the former will clobber the log settings
	hasLogs := len(d.Get("logs").([]interface{})) > 0
	if d.HasChange("logs") || (hasLogs && d.HasChange("app_settings")) {
		logs := expandAppServiceLogs(d.Get("logs"))
		logsResource := web.SiteLogsConfig{
			ID:                       utils.String(d.Id()),
			SiteLogsConfigProperties: &logs,
		}

		if _, err := client.UpdateDiagnosticLogsConfig(ctx, id.ResourceGroup, id.Name, logsResource); err != nil {
			return fmt.Errorf("Error updating Diagnostics Logs for App Service %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("storage_account") {
		storageAccountsRaw := d.Get("storage_account").(*schema.Set).List()
		storageAccounts := expandAppServiceStorageAccounts(storageAccountsRaw)
		properties := web.AzureStoragePropertyDictionaryResource{
			Properties: storageAccounts,
		}

		if _, err := client.UpdateAzureStorageAccounts(ctx, id.ResourceGroup, id.Name, properties); err != nil {
			return fmt.Errorf("Error updating Storage Accounts for App Service %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("connection_string") {
		// update the ConnectionStrings
		connectionStrings := expandAppServiceConnectionStrings(d)
		properties := web.ConnectionStringDictionary{
			Properties: connectionStrings,
		}

		if _, err := client.UpdateConnectionStrings(ctx, id.ResourceGroup, id.Name, properties); err != nil {
			return fmt.Errorf("Error updating Connection Strings for App Service %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("identity") {
		site, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("Error getting configuration for App Service %q: %+v", id.Name, err)
		}

		appServiceIdentityRaw := d.Get("identity").([]interface{})
		appServiceIdentity := expandAppServiceIdentity(appServiceIdentityRaw)
		site.Identity = appServiceIdentity

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, site)
		if err != nil {
			return fmt.Errorf("Error updating Managed Service Identity for App Service %q: %+v", id.Name, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error updating Managed Service Identity for App Service %q: %+v", id.Name, err)
		}
	}

	return resourceArmAppServiceRead(d, meta)
}

func resourceArmAppServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service %q (resource group %q) was not found - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service %q: %+v", id.Name, err)
	}

	configResp, err := client.GetConfiguration(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(configResp.Response) {
			log.Printf("[DEBUG] Configuration of App Service %q (resource group %q) was not found", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service Configuration %q: %+v", id.Name, err)
	}

	authResp, err := client.GetAuthSettings(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving the AuthSettings for App Service %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	backupResp, err := client.GetBackupConfiguration(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(backupResp.Response) {
			return fmt.Errorf("Error retrieving the BackupConfiguration for App Service %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	logsResp, err := client.GetDiagnosticLogsConfiguration(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving the DiagnosticsLogsConfiguration for App Service %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	appSettingsResp, err := client.ListApplicationSettings(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(appSettingsResp.Response) {
			log.Printf("[DEBUG] Application Settings of App Service %q (resource group %q) were not found", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service AppSettings %q: %+v", id.Name, err)
	}

	storageAccountsResp, err := client.ListAzureStorageAccounts(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service Storage Accounts %q: %+v", id.Name, err)
	}

	connectionStringsResp, err := client.ListConnectionStrings(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service ConnectionStrings %q: %+v", id.Name, err)
	}

	scmResp, err := client.GetSourceControl(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service Source Control %q: %+v", id.Name, err)
	}

	siteCredFuture, err := client.ListPublishingCredentials(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return err
	}
	err = siteCredFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return err
	}
	siteCredResp, err := siteCredFuture.Result(*client)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service Site Credential %q: %+v", id.Name, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
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

	// remove DIAGNOSTICS*, WEBSITE_HTTPLOGGING* settings - Azure will sync these, so just maintain the logs block equivalents in the state
	delete(appSettings, "DIAGNOSTICS_AZUREBLOBCONTAINERSASURL")
	delete(appSettings, "DIAGNOSTICS_AZUREBLOBRETENTIONINDAYS")
	delete(appSettings, "WEBSITE_HTTPLOGGING_CONTAINER_URL")
	delete(appSettings, "WEBSITE_HTTPLOGGING_RETENTION_DAYS")

	if err := d.Set("app_settings", appSettings); err != nil {
		return fmt.Errorf("Error setting `app_settings`: %s", err)
	}

	if err := d.Set("backup", flattenAppServiceBackup(backupResp.BackupRequestProperties)); err != nil {
		return fmt.Errorf("Error setting `backup`: %s", err)
	}

	if err := d.Set("storage_account", flattenAppServiceStorageAccounts(storageAccountsResp.Properties)); err != nil {
		return fmt.Errorf("Error setting `storage_account`: %s", err)
	}

	if err := d.Set("connection_string", flattenAppServiceConnectionStrings(connectionStringsResp.Properties)); err != nil {
		return fmt.Errorf("Error setting `connection_string`: %s", err)
	}

	siteConfig := flattenAppServiceSiteConfig(configResp.SiteConfig)
	if err := d.Set("site_config", siteConfig); err != nil {
		return err
	}

	authSettings := flattenAppServiceAuthSettings(authResp.SiteAuthSettingsProperties)
	if err := d.Set("auth_settings", authSettings); err != nil {
		return fmt.Errorf("Error setting `auth_settings`: %s", err)
	}

	logs := flattenAppServiceLogs(logsResp.SiteLogsConfigProperties)
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

	identity := flattenAppServiceIdentity(resp.Identity)
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("Error setting `identity`: %s", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmAppServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting App Service %q (resource group %q)", id.Name, id.ResourceGroup)

	deleteMetrics := true
	deleteEmptyServerFarm := false
	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name, &deleteMetrics, &deleteEmptyServerFarm)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return err
		}
	}

	return nil
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
