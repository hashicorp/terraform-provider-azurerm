// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAppService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceCreate,
		Read:   resourceAppServiceRead,
		Update: resourceAppServiceUpdate,
		Delete: resourceAppServiceDelete,

		DeprecationMessage: "The `azurerm_app_service` resource has been superseded by the `azurerm_linux_web_app` and `azurerm_windows_web_app` resources. Whilst this resource will continue to be available in the 2.x and 3.x releases it is feature-frozen for compatibility purposes, will no longer receive any updates and will be removed in a future major release of the Azure Provider.",

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AppServiceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AppServiceName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"app_service_plan_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"app_settings": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"auth_settings": schemaAppServiceAuthSettings(),

			"backup": schemaAppServiceBackup(),

			"client_affinity_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"client_cert_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"client_cert_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(web.ClientCertModeOptional),
					string(web.ClientCertModeRequired),
					string(web.ClientCertModeOptionalInteractiveUser),
				}, false),
			},

			"connection_string": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(web.ConnectionStringTypeAPIHub),
								string(web.ConnectionStringTypeCustom),
								string(web.ConnectionStringTypeDocDb),
								string(web.ConnectionStringTypeEventHub),
								string(web.ConnectionStringTypeMySQL),
								string(web.ConnectionStringTypeNotificationHub),
								string(web.ConnectionStringTypePostgreSQL),
								string(web.ConnectionStringTypeRedisCache),
								string(web.ConnectionStringTypeServiceBus),
								string(web.ConnectionStringTypeSQLAzure),
								string(web.ConnectionStringTypeSQLServer),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"value": {
							Type:      pluginsdk.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"https_only": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"key_vault_reference_identity_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			},

			"logs": schemaAppServiceLogsConfig(),

			"site_config": schemaAppServiceSiteConfig(),

			"storage_account": schemaAppServiceStorageAccounts(),

			"source_control": schemaAppServiceSiteSourceControl(),

			"tags": tags.Schema(),

			"site_credential": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"username": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"password": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"custom_domain_verification_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"default_site_hostname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"outbound_ip_addresses": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"outbound_ip_address_list": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"possible_outbound_ip_addresses": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"possible_outbound_ip_address_list": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func resourceAppServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	aspClient := meta.(*clients.Client).Web.AppServicePlansClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM App Service creation.")
	id := parse.NewAppServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_app_service", id.ID())
	}

	availabilityRequest := web.ResourceNameAvailabilityRequest{
		Name: utils.String(id.SiteName),
		Type: web.CheckNameResourceTypesMicrosoftWebsites,
	}

	appServicePlanId := d.Get("app_service_plan_id").(string)
	aspID, err := parse.AppServicePlanID(appServicePlanId)
	if err != nil {
		return err
	}
	// Check if App Service Plan is part of ASE
	// If so, the name needs updating to <app name>.<ASE name>.appserviceenvironment.net and FQDN setting true for name availability check
	aspDetails, err := aspClient.Get(ctx, aspID.ResourceGroup, aspID.ServerFarmName)
	// 404 is incorrectly being considered an acceptable response, issue tracked at https://github.com/Azure/azure-sdk-for-go/issues/15002
	if err != nil || utils.ResponseWasNotFound(aspDetails.Response) {
		return fmt.Errorf("App Service Environment %q or Resource Group %q does not exist", aspID.ServerFarmName, aspID.ResourceGroup)
	}
	if aspDetails.HostingEnvironmentProfile != nil {
		availabilityRequest.Name = utils.String(fmt.Sprintf("%s.%s.appserviceenvironment.net", id.SiteName, *aspDetails.HostingEnvironmentProfile.Name))
		availabilityRequest.IsFqdn = utils.Bool(true)
	}
	available, err := client.CheckNameAvailability(ctx, availabilityRequest)
	if err != nil {
		return fmt.Errorf("checking if the name %q was available: %+v", id.SiteName, err)
	}

	if !*available.NameAvailable {
		return fmt.Errorf("The name %q used for the App Service needs to be globally unique and isn't available: %s", id.SiteName, *available.Message)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	enabled := d.Get("enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	t := d.Get("tags").(map[string]interface{})

	siteConfig, err := expandAppServiceSiteConfig(d.Get("site_config"))
	if err != nil {
		return fmt.Errorf("expanding `site_config` for %s: %s", id, err)
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

	if v, ok := d.GetOk("key_vault_reference_identity_id"); ok {
		siteEnvelope.SiteProperties.KeyVaultReferenceIdentity = utils.String(v.(string))
	}

	if _, ok := d.GetOk("identity"); ok {
		appServiceIdentity, err := expandAppServiceIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		siteEnvelope.Identity = appServiceIdentity
	}

	siteEnvelope.SiteProperties.ClientAffinityEnabled = utils.Bool(d.Get("client_affinity_enabled").(bool))

	siteEnvelope.SiteProperties.ClientCertEnabled = utils.Bool(d.Get("client_cert_enabled").(bool))
	if *siteEnvelope.SiteProperties.ClientCertEnabled {
		if clientCertMode, ok := d.GetOk("client_cert_mode"); ok {
			siteEnvelope.SiteProperties.ClientCertMode = web.ClientCertMode(clientCertMode.(string))
		}
	}

	createFuture, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, siteEnvelope)
	if err != nil {
		return fmt.Errorf("creating %s: %s", id, err)
	}

	err = createFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting for %s to be created: %s", id, err)
	}

	if _, ok := d.GetOk("source_control"); ok {
		if siteConfig.ScmType != "" {
			return fmt.Errorf("cannot set source_control parameters when scm_type is set to %q", siteConfig.ScmType)
		}
		sourceControlProperties := expandAppServiceSiteSourceControl(d)
		sourceControl := &web.SiteSourceControl{}
		sourceControl.SiteSourceControlProperties = sourceControlProperties
		// TODO - Do we need to lock the app for updates?
		scFuture, err := client.CreateOrUpdateSourceControl(ctx, id.ResourceGroup, id.SiteName, *sourceControl)
		if err != nil {
			return fmt.Errorf("failed to create %s: %+v", id, err)
		}

		err = scFuture.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("failed waiting for App Service Source Control configuration")
		}
	}

	d.SetId(id.ID())

	authSettingsRaw := d.Get("auth_settings").([]interface{})
	authSettings := expandAppServiceAuthSettings(authSettingsRaw)

	auth := web.SiteAuthSettings{
		ID:                         utils.String(id.ID()),
		SiteAuthSettingsProperties: &authSettings,
	}

	if _, err := client.UpdateAuthSettings(ctx, id.ResourceGroup, id.SiteName, auth); err != nil {
		return fmt.Errorf("updating auth settings for %s: %+v", id, err)
	}

	logsConfig := expandAppServiceLogs(d.Get("logs"))

	logs := web.SiteLogsConfig{
		ID:                       utils.String(id.ID()),
		SiteLogsConfigProperties: &logsConfig,
	}

	if _, err := client.UpdateDiagnosticLogsConfig(ctx, id.ResourceGroup, id.SiteName, logs); err != nil {
		return fmt.Errorf("updating diagnostic logs config for %s: %+v", id, err)
	}

	backupRaw := d.Get("backup").([]interface{})
	if backup := expandAppServiceBackup(backupRaw); backup != nil {
		if _, err = client.UpdateBackupConfiguration(ctx, id.ResourceGroup, id.SiteName, *backup); err != nil {
			return fmt.Errorf("updating Backup Settings for %s: %+v", id, err)
		}
	}

	return resourceAppServiceUpdate(d, meta)
}

func resourceAppServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("expanding `site_config` for App Service %q (Resource Group %q): %s", id.SiteName, id.ResourceGroup, err)
	}

	// WEBSITE_VNET_ROUTE_ALL is superseded by a setting in site_config that defaults to false from 2021-02-01
	appSettings := expandAppServiceAppSettings(d)
	if vnetRouteAll, ok := appSettings["WEBSITE_VNET_ROUTE_ALL"]; ok {
		if !d.HasChange("site_config.0.vnet_route_all_enabled") { // Only update the property if it's not set explicitly
			siteConfig.VnetRouteAllEnabled = utils.Bool(strings.EqualFold(*vnetRouteAll, "true"))
		}
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

	if v, ok := d.GetOk("key_vault_reference_identity_id"); ok {
		siteEnvelope.SiteProperties.KeyVaultReferenceIdentity = utils.String(v.(string))
	}

	siteEnvelope.SiteProperties.ClientCertEnabled = utils.Bool(d.Get("client_cert_enabled").(bool))

	if *siteEnvelope.SiteProperties.ClientCertEnabled {
		if clientCertMode, ok := d.GetOk("client_cert_mode"); ok {
			siteEnvelope.SiteProperties.ClientCertMode = web.ClientCertMode(clientCertMode.(string))
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, siteEnvelope)
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
			return fmt.Errorf("expanding `site_config` for App Service %q (Resource Group %q): %s", id.SiteName, id.ResourceGroup, err)
		}
		siteConfigResource := web.SiteConfigResource{
			SiteConfig: siteConfig,
		}

		scmType = siteConfig.ScmType
		// ScmType being set blocks the update of source_control in _most_ cases, ADO is an exception
		if hasSourceControl && scmType != web.ScmTypeVSTSRM {
			siteConfigResource.SiteConfig.ScmType = web.ScmTypeNone
		}

		if _, err := client.CreateOrUpdateConfiguration(ctx, id.ResourceGroup, id.SiteName, siteConfigResource); err != nil {
			return fmt.Errorf("updating Configuration for App Service %q: %+v", id.SiteName, err)
		}
	}

	if d.HasChange("auth_settings") {
		authSettingsRaw := d.Get("auth_settings").([]interface{})
		authSettingsProperties := expandAppServiceAuthSettings(authSettingsRaw)
		authSettings := web.SiteAuthSettings{
			ID:                         utils.String(d.Id()),
			SiteAuthSettingsProperties: &authSettingsProperties,
		}

		if _, err := client.UpdateAuthSettings(ctx, id.ResourceGroup, id.SiteName, authSettings); err != nil {
			return fmt.Errorf("updating Authentication Settings for App Service %q: %+v", id.SiteName, err)
		}
	}

	if d.HasChange("backup") {
		backupRaw := d.Get("backup").([]interface{})
		if backup := expandAppServiceBackup(backupRaw); backup != nil {
			if _, err = client.UpdateBackupConfiguration(ctx, id.ResourceGroup, id.SiteName, *backup); err != nil {
				return fmt.Errorf("updating Backup Settings for App Service %q (Resource Group %q): %s", id.SiteName, id.ResourceGroup, err)
			}
		} else {
			if _, err = client.DeleteBackupConfiguration(ctx, id.ResourceGroup, id.SiteName); err != nil {
				return fmt.Errorf("removing Backup Settings for App Service %q (Resource Group %q): %s", id.SiteName, id.ResourceGroup, err)
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

		if _, err := client.Update(ctx, id.ResourceGroup, id.SiteName, sitePatchResource); err != nil {
			return fmt.Errorf("updating App Service ARR Affinity setting %q: %+v", id.SiteName, err)
		}
	}

	// app settings updates have a side effect on logging settings. See the note below
	if d.HasChange("app_settings") {
		// update the AppSettings
		appSettings = expandAppServiceAppSettings(d)

		settings := web.StringDictionary{
			Properties: appSettings,
		}

		if _, err := client.UpdateApplicationSettings(ctx, id.ResourceGroup, id.SiteName, settings); err != nil {
			return fmt.Errorf("updating Application Settings for App Service %q: %+v", id.SiteName, err)
		}
	}

	// Don't send source_control changes for ADO controlled Apps
	if hasSourceControl && scmType != web.ScmTypeVSTSRM {
		sourceControlProperties := expandAppServiceSiteSourceControl(d)
		sourceControl := &web.SiteSourceControl{}
		sourceControl.SiteSourceControlProperties = sourceControlProperties
		scFuture, err := client.CreateOrUpdateSourceControl(ctx, id.ResourceGroup, id.SiteName, *sourceControl)
		if err != nil {
			return fmt.Errorf("failed to update App Service Source Control for %q (Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
		}

		err = scFuture.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("failed waiting for App Service Source Control configuration: %+v", err)
		}

		sc, err := client.GetSourceControl(ctx, id.ResourceGroup, id.SiteName)
		if err != nil {
			return fmt.Errorf("failed reading back App Service Source Control for %q", *sc.Name)
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

		if _, err := client.UpdateDiagnosticLogsConfig(ctx, id.ResourceGroup, id.SiteName, logsResource); err != nil {
			return fmt.Errorf("updating Diagnostics Logs for App Service %q: %+v", id.SiteName, err)
		}
	}

	if d.HasChange("storage_account") {
		storageAccountsRaw := d.Get("storage_account").(*pluginsdk.Set).List()
		storageAccounts := expandAppServiceStorageAccounts(storageAccountsRaw)
		properties := web.AzureStoragePropertyDictionaryResource{
			Properties: storageAccounts,
		}

		if _, err := client.UpdateAzureStorageAccounts(ctx, id.ResourceGroup, id.SiteName, properties); err != nil {
			return fmt.Errorf("updating Storage Accounts for App Service %q: %+v", id.SiteName, err)
		}
	}

	if d.HasChange("connection_string") {
		// update the ConnectionStrings
		connectionStrings := expandAppServiceConnectionStrings(d)
		properties := web.ConnectionStringDictionary{
			Properties: connectionStrings,
		}

		if _, err := client.UpdateConnectionStrings(ctx, id.ResourceGroup, id.SiteName, properties); err != nil {
			return fmt.Errorf("updating Connection Strings for App Service %q: %+v", id.SiteName, err)
		}
	}

	if d.HasChange("identity") {
		site, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
		if err != nil {
			return fmt.Errorf("getting configuration for App Service %q: %+v", id.SiteName, err)
		}

		appServiceIdentity, err := expandAppServiceIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		site.Identity = appServiceIdentity
		site.SiteConfig = siteConfig

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, site)
		if err != nil {
			return fmt.Errorf("updating Managed Service Identity for App Service %q: %+v", id.SiteName, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("updating Managed Service Identity for App Service %q: %+v", id.SiteName, err)
		}
	}

	return resourceAppServiceRead(d, meta)
}

func resourceAppServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service %q (resource group %q) was not found - removing from state", id.SiteName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on AzureRM App Service %q: %+v", id.SiteName, err)
	}

	configResp, err := client.GetConfiguration(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(configResp.Response) {
			log.Printf("[DEBUG] Configuration of App Service %q (resource group %q) was not found", id.SiteName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on AzureRM App Service Configuration %q: %+v", id.SiteName, err)
	}

	authResp, err := client.GetAuthSettings(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("retrieving the AuthSettings for App Service %q (Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
	}

	backupResp, err := client.GetBackupConfiguration(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if !utils.ResponseWasNotFound(backupResp.Response) {
			return fmt.Errorf("retrieving the BackupConfiguration for App Service %q (Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
		}
	}

	logsResp, err := client.GetDiagnosticLogsConfiguration(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("retrieving the DiagnosticsLogsConfiguration for App Service %q (Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
	}

	appSettingsResp, err := client.ListApplicationSettings(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(appSettingsResp.Response) {
			log.Printf("[DEBUG] Application Settings of App Service %q (resource group %q) were not found", id.SiteName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on AzureRM App Service AppSettings %q: %+v", id.SiteName, err)
	}

	storageAccountsResp, err := client.ListAzureStorageAccounts(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("making Read request on AzureRM App Service Storage Accounts %q: %+v", id.SiteName, err)
	}

	connectionStringsResp, err := client.ListConnectionStrings(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("making Read request on AzureRM App Service ConnectionStrings %q: %+v", id.SiteName, err)
	}

	scmResp, err := client.GetSourceControl(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("making Read request on AzureRM App Service Source Control %q: %+v", id.SiteName, err)
	}

	siteCredFuture, err := client.ListPublishingCredentials(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return err
	}
	err = siteCredFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return err
	}
	siteCredResp, err := siteCredFuture.Result(*client)
	if err != nil {
		return fmt.Errorf("making Read request on AzureRM App Service Site Credential %q: %+v", id.SiteName, err)
	}

	d.Set("name", id.SiteName)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.SiteProperties; props != nil {
		servicePlan, err := commonids.ParseAppServicePlanIDInsensitively(pointer.From(props.ServerFarmID))
		if err != nil {
			return err
		}
		d.Set("app_service_plan_id", servicePlan.ID())
		d.Set("client_affinity_enabled", props.ClientAffinityEnabled)
		d.Set("enabled", props.Enabled)
		d.Set("https_only", props.HTTPSOnly)
		d.Set("client_cert_enabled", props.ClientCertEnabled)
		d.Set("client_cert_mode", props.ClientCertMode)
		d.Set("default_site_hostname", props.DefaultHostName)
		d.Set("outbound_ip_addresses", props.OutboundIPAddresses)
		if props.OutboundIPAddresses != nil {
			d.Set("outbound_ip_address_list", strings.Split(*props.OutboundIPAddresses, ","))
		}
		d.Set("possible_outbound_ip_addresses", props.PossibleOutboundIPAddresses)
		if props.PossibleOutboundIPAddresses != nil {
			d.Set("possible_outbound_ip_address_list", strings.Split(*props.PossibleOutboundIPAddresses, ","))
		}
		d.Set("custom_domain_verification_id", props.CustomDomainVerificationID)

		if props.KeyVaultReferenceIdentity != nil {
			d.Set("key_vault_reference_identity_id", props.KeyVaultReferenceIdentity)
		}
	}

	appSettings := flattenAppServiceAppSettings(appSettingsResp.Properties)

	// remove DIAGNOSTICS*, WEBSITE_HTTPLOGGING* settings - Azure will sync these, so just maintain the logs block equivalents in the state
	delete(appSettings, "DIAGNOSTICS_AZUREBLOBCONTAINERSASURL")
	delete(appSettings, "DIAGNOSTICS_AZUREBLOBRETENTIONINDAYS")
	delete(appSettings, "WEBSITE_HTTPLOGGING_CONTAINER_URL")
	delete(appSettings, "WEBSITE_HTTPLOGGING_RETENTION_DAYS")

	if err := d.Set("app_settings", appSettings); err != nil {
		return fmt.Errorf("setting `app_settings`: %s", err)
	}

	if err := d.Set("backup", flattenAppServiceBackup(backupResp.BackupRequestProperties)); err != nil {
		return fmt.Errorf("setting `backup`: %s", err)
	}

	if err := d.Set("storage_account", flattenAppServiceStorageAccounts(storageAccountsResp.Properties)); err != nil {
		return fmt.Errorf("setting `storage_account`: %s", err)
	}

	if err := d.Set("connection_string", flattenAppServiceConnectionStrings(connectionStringsResp.Properties)); err != nil {
		return fmt.Errorf("setting `connection_string`: %s", err)
	}

	siteConfig := flattenAppServiceSiteConfig(configResp.SiteConfig)
	if err := d.Set("site_config", siteConfig); err != nil {
		return err
	}

	authSettings := flattenAppServiceAuthSettings(authResp.SiteAuthSettingsProperties)
	if err := d.Set("auth_settings", authSettings); err != nil {
		return fmt.Errorf("setting `auth_settings`: %s", err)
	}

	logs := flattenAppServiceLogs(logsResp.SiteLogsConfigProperties)
	if err := d.Set("logs", logs); err != nil {
		return fmt.Errorf("setting `logs`: %s", err)
	}

	scm := flattenAppServiceSourceControl(scmResp.SiteSourceControlProperties)
	if err := d.Set("source_control", scm); err != nil {
		return fmt.Errorf("setting `source_control`: %s", err)
	}

	siteCred := flattenAppServiceSiteCredential(siteCredResp.UserProperties)
	if err := d.Set("site_credential", siteCred); err != nil {
		return fmt.Errorf("setting `site_credential`: %s", err)
	}

	identity, err := flattenAppServiceIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %s", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceAppServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting App Service %q (resource group %q)", id.SiteName, id.ResourceGroup)

	deleteMetrics := true
	deleteEmptyServerFarm := false
	resp, err := client.Delete(ctx, id.ResourceGroup, id.SiteName, &deleteMetrics, &deleteEmptyServerFarm)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return err
		}
	}

	return nil
}

func expandAppServiceAppSettings(d *pluginsdk.ResourceData) map[string]*string {
	input := d.Get("app_settings").(map[string]interface{})
	output := make(map[string]*string, len(input))

	for k, v := range input {
		output[k] = utils.String(v.(string))
	}

	return output
}

func expandAppServiceConnectionStrings(d *pluginsdk.ResourceData) map[string]*web.ConnStringValueTypePair {
	input := d.Get("connection_string").(*pluginsdk.Set).List()
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
