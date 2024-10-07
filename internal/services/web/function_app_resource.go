// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	webValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// Azure Function App shares the same infrastructure with Azure App Service.
// So this resource will reuse most of the App Service code, but remove the configurations which are not applicable for Function App.
func resourceFunctionApp() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFunctionAppCreate,
		Read:   resourceFunctionAppRead,
		Update: resourceFunctionAppUpdate,
		Delete: resourceFunctionAppDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FunctionAppID(id)
			return err
		}),

		DeprecationMessage: "The `azurerm_function_app` resource has been superseded by the `azurerm_linux_function_app` and `azurerm_windows_function_app` resources. Whilst this resource will continue to be available in the 2.x and 3.x releases it is feature-frozen for compatibility purposes, will no longer receive any updates and will be removed in a future major release of the Azure Provider.",

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
				ValidateFunc: webValidate.AppServiceName,
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

			"client_cert_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Required",
					"Optional",
				}, false),
			},

			"daily_memory_time_quota": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"enable_builtin_logging": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"https_only": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"os_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
				ValidateFunc: validation.StringInSlice([]string{
					"linux",
					"",
				}, false),
			},

			"key_vault_reference_identity_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			},

			"site_config": schemaAppServiceFunctionAppSiteConfig(),

			"source_control": schemaAppServiceSiteSourceControl(),

			"storage_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageAccountName,
			},

			"storage_account_access_key": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.NoZeroValues,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "~1",
			},

			"tags": tags.Schema(),

			// Computed Only

			"custom_domain_verification_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"default_hostname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"kind": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"outbound_ip_addresses": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"possible_outbound_ip_addresses": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

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
		},
	}
}

func resourceFunctionAppCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageDomainSuffix, ok := meta.(*clients.Client).Account.Environment.Storage.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Storage domain suffix for environment %q", meta.(*clients.Client).Account.Environment.Name)
	}

	log.Printf("[INFO] preparing arguments for AzureRM Function App creation.")

	id := parse.NewFunctionAppID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_function_app", id.ID())
	}

	availabilityRequest := web.ResourceNameAvailabilityRequest{
		Name: utils.String(id.SiteName),
		Type: web.CheckNameResourceTypesMicrosoftWebsites,
	}
	available, err := client.CheckNameAvailability(ctx, availabilityRequest)
	if err != nil {
		return fmt.Errorf("checking if the name %q was available: %+v", id.SiteName, err)
	}

	if !*available.NameAvailable {
		return fmt.Errorf("the name %q used for the Function App needs to be globally unique and isn't available: %s", id.SiteName, *available.Message)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	kind := "functionapp"
	if osTypeRaw, ok := d.GetOk("os_type"); ok {
		osType := osTypeRaw.(string)
		if osType == "linux" {
			kind = "functionapp,linux"
		}
	}

	appServicePlanID := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	clientCertMode := d.Get("client_cert_mode").(string)
	clientCertEnabled := clientCertMode != ""
	httpsOnly := d.Get("https_only").(bool)
	dailyMemoryTimeQuota := d.Get("daily_memory_time_quota").(int)
	t := d.Get("tags").(map[string]interface{})
	appServiceTier, err := getFunctionAppServiceTier(ctx, appServicePlanID, meta)
	if err != nil {
		return err
	}

	basicAppSettings, err := getBasicFunctionAppAppSettings(d, appServiceTier, *storageDomainSuffix, nil)
	if err != nil {
		return err
	}

	appSettings := expandFunctionAppAppSettings(d, basicAppSettings)

	siteConfig, err := expandFunctionAppSiteConfig(d)
	if err != nil {
		return fmt.Errorf("expanding `site_config` for %s: %s", id, err)
	}

	siteConfig.AppSettings = appSettingsMapToNameValuePair(appSettings)

	siteEnvelope := web.Site{
		Kind:     &kind,
		Location: &location,
		Tags:     tags.Expand(t),
		SiteProperties: &web.SiteProperties{
			ServerFarmID:         utils.String(appServicePlanID),
			Enabled:              utils.Bool(enabled),
			ClientCertEnabled:    utils.Bool(clientCertEnabled),
			HTTPSOnly:            utils.Bool(httpsOnly),
			DailyMemoryTimeQuota: utils.Int32(int32(dailyMemoryTimeQuota)),
			SiteConfig:           &siteConfig,
		},
	}

	if v, ok := d.GetOk("key_vault_reference_identity_id"); ok {
		siteEnvelope.SiteProperties.KeyVaultReferenceIdentity = utils.String(v.(string))
	}

	if clientCertMode != "" {
		siteEnvelope.SiteProperties.ClientCertMode = web.ClientCertMode(clientCertMode)
	}

	if _, ok := d.GetOk("identity"); ok {
		appServiceIdentity, err := expandAppServiceIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		siteEnvelope.Identity = appServiceIdentity
	}

	createFuture, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, siteEnvelope)
	if err != nil {
		return err
	}

	err = createFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("source_control"); ok {
		if siteConfig.ScmType != "" {
			return fmt.Errorf("cannot set source_control parameters when scm_type is set to %q", siteConfig.ScmType)
		}
		sourceControlProperties := expandAppServiceSiteSourceControl(d)
		sourceControl := &web.SiteSourceControl{}
		sourceControl.SiteSourceControlProperties = sourceControlProperties
		// TODO - Do we need to lock the function app for updates?
		future, err := client.CreateOrUpdateSourceControl(ctx, id.ResourceGroup, id.SiteName, *sourceControl)
		if err != nil {
			return fmt.Errorf("failed to create %s: %+v", id, err)
		}
		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for creation of %q: %+v", id, err)
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

	return resourceFunctionAppUpdate(d, meta)
}

func resourceFunctionAppUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageDomainSuffix, ok := meta.(*clients.Client).Account.Environment.Storage.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Storage domain suffix for environment %q", meta.(*clients.Client).Account.Environment.Name)
	}

	id, err := parse.FunctionAppID(d.Id())
	if err != nil {
		return err
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	kind := "functionapp"
	if osTypeRaw, ok := d.GetOk("os_type"); ok {
		osType := osTypeRaw.(string)
		if osType == "linux" {
			kind = "functionapp,linux"
		}
	}
	appServicePlanID := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	clientCertMode := d.Get("client_cert_mode").(string)
	clientCertEnabled := clientCertMode != ""
	httpsOnly := d.Get("https_only").(bool)
	dailyMemoryTimeQuota := d.Get("daily_memory_time_quota").(int)
	t := d.Get("tags").(map[string]interface{})

	appServiceTier, err := getFunctionAppServiceTier(ctx, appServicePlanID, meta)
	if err != nil {
		return err
	}

	var currentAppSettings map[string]*string
	appSettingsList, err := client.ListApplicationSettings(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("reading App Settings for %s: %+v", id, err)
	}
	if appSettingsList.Properties != nil {
		currentAppSettings = appSettingsList.Properties
	}

	basicAppSettings, err := getBasicFunctionAppAppSettings(d, appServiceTier, *storageDomainSuffix, currentAppSettings)
	if err != nil {
		return err
	}

	siteConfig, err := expandFunctionAppSiteConfig(d)
	if err != nil {
		return fmt.Errorf("expanding `site_config` for Function App %q (Resource Group %q): %s", id.SiteName, id.ResourceGroup, err)
	}

	siteConfig.AppSettings = &basicAppSettings

	// WEBSITE_VNET_ROUTE_ALL is superseded by a setting in site_config that defaults to false from 2021-02-01
	appSettings := expandFunctionAppAppSettings(d, basicAppSettings)
	if vnetRouteAll, ok := appSettings["WEBSITE_VNET_ROUTE_ALL"]; ok {
		if !d.HasChange("site_config.0.vnet_route_all_enabled") { // Only update the property if it's not set explicitly
			vnetRouteAllEnabled, _ := strconv.ParseBool(*vnetRouteAll)
			siteConfig.VnetRouteAllEnabled = &vnetRouteAllEnabled
		}
	}

	siteEnvelope := web.Site{
		Kind:     &kind,
		Location: &location,
		Tags:     tags.Expand(t),
		SiteProperties: &web.SiteProperties{
			ServerFarmID:         utils.String(appServicePlanID),
			Enabled:              utils.Bool(enabled),
			ClientCertEnabled:    utils.Bool(clientCertEnabled),
			HTTPSOnly:            utils.Bool(httpsOnly),
			DailyMemoryTimeQuota: utils.Int32(int32(dailyMemoryTimeQuota)),
			SiteConfig:           &siteConfig,
		},
	}

	if v, ok := d.GetOk("key_vault_reference_identity_id"); ok {
		siteEnvelope.SiteProperties.KeyVaultReferenceIdentity = utils.String(v.(string))
	}

	if clientCertMode != "" {
		siteEnvelope.SiteProperties.ClientCertMode = web.ClientCertMode(clientCertMode)
	}

	if _, ok := d.GetOk("identity"); ok {
		appServiceIdentity, err := expandAppServiceIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		siteEnvelope.Identity = appServiceIdentity
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, siteEnvelope)
	if err != nil {
		return fmt.Errorf("updating Function App %q (Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of Function App %q (Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
	}

	settings := web.StringDictionary{
		Properties: appSettings,
	}

	if _, err = client.UpdateApplicationSettings(ctx, id.ResourceGroup, id.SiteName, settings); err != nil {
		return fmt.Errorf("updating Application Settings for Function App %q: %+v", id.SiteName, err)
	}

	// If `source_control` is defined, we need to set site_config.0.scm_type to "None" or we cannot update it
	// repo_url is required by the API
	_, hasSourceControl := d.GetOk("source_control.0.repo_url")

	scmType := web.ScmTypeNone

	if d.HasChange("site_config") || hasSourceControl {
		siteConfig, err := expandFunctionAppSiteConfig(d)
		if err != nil {
			return fmt.Errorf("expanding `site_config` for Function App %q (Resource Group %q): %s", id.SiteName, id.ResourceGroup, err)
		}
		siteConfigResource := web.SiteConfigResource{
			SiteConfig: &siteConfig,
		}

		scmType = siteConfig.ScmType
		// ScmType being set blocks the update of source_control in _most_ cases, ADO is an exception
		if hasSourceControl && scmType != web.ScmTypeVSTSRM {
			siteConfigResource.SiteConfig.ScmType = web.ScmTypeNone
		}

		if _, err := client.CreateOrUpdateConfiguration(ctx, id.ResourceGroup, id.SiteName, siteConfigResource); err != nil {
			return fmt.Errorf("updating Configuration for Function App %q: %+v", id.SiteName, err)
		}
	}

	// Don't send source_control changes for ADO controlled Apps
	if hasSourceControl && scmType != web.ScmTypeVSTSRM {
		sourceControlProperties := expandAppServiceSiteSourceControl(d)
		sourceControl := &web.SiteSourceControl{}
		sourceControl.SiteSourceControlProperties = sourceControlProperties
		scFuture, err := client.CreateOrUpdateSourceControl(ctx, id.ResourceGroup, id.SiteName, *sourceControl)
		if err != nil {
			return fmt.Errorf("failed to create App Service Source Control for %q (Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
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

	if d.HasChange("auth_settings") {
		authSettingsRaw := d.Get("auth_settings").([]interface{})
		authSettingsProperties := expandAppServiceAuthSettings(authSettingsRaw)
		authSettings := web.SiteAuthSettings{
			ID:                         utils.String(d.Id()),
			SiteAuthSettingsProperties: &authSettingsProperties,
		}

		if _, err := client.UpdateAuthSettings(ctx, id.ResourceGroup, id.SiteName, authSettings); err != nil {
			return fmt.Errorf("updating Authentication Settings for Function App %q: %+v", id.SiteName, err)
		}
	}

	if d.HasChange("connection_string") {
		// update the ConnectionStrings
		connectionStrings := expandFunctionAppConnectionStrings(d)
		properties := web.ConnectionStringDictionary{
			Properties: connectionStrings,
		}

		if _, err := client.UpdateConnectionStrings(ctx, id.ResourceGroup, id.SiteName, properties); err != nil {
			return fmt.Errorf("updating Connection Strings for App Service %q: %+v", id.SiteName, err)
		}
	}

	return resourceFunctionAppRead(d, meta)
}

func resourceFunctionAppRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FunctionAppID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Function App %q (resource group %q) was not found - removing from state", id.SiteName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on AzureRM Function App %q: %+v", id.SiteName, err)
	}

	appSettingsResp, err := client.ListApplicationSettings(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(appSettingsResp.Response) {
			log.Printf("[DEBUG] Application Settings of Function App %q (resource group %q) were not found", id.SiteName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on AzureRM Function App AppSettings %q: %+v", id.SiteName, err)
	}

	connectionStringsResp, err := client.ListConnectionStrings(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("making Read request on AzureRM Function App ConnectionStrings %q: %+v", id.SiteName, err)
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
	authResp, err := client.GetAuthSettings(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("retrieving the AuthSettings for Function App %q (Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
	}

	d.Set("name", id.SiteName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("kind", resp.Kind)
	osType := ""
	if v := resp.Kind; v != nil && strings.Contains(*v, "linux") {
		osType = "linux"
	}
	d.Set("os_type", osType)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	appServicePlanID := ""
	if props := resp.SiteProperties; props != nil {
		servicePlan, err := commonids.ParseAppServicePlanIDInsensitively(pointer.From(props.ServerFarmID))
		if err != nil {
			return err
		}
		appServicePlanID = servicePlan.ID()
		d.Set("app_service_plan_id", appServicePlanID)
		d.Set("enabled", props.Enabled)
		d.Set("default_hostname", props.DefaultHostName)
		d.Set("https_only", props.HTTPSOnly)
		d.Set("daily_memory_time_quota", props.DailyMemoryTimeQuota)
		d.Set("outbound_ip_addresses", props.OutboundIPAddresses)
		d.Set("possible_outbound_ip_addresses", props.PossibleOutboundIPAddresses)
		d.Set("custom_domain_verification_id", props.CustomDomainVerificationID)

		clientCertMode := ""
		if props.ClientCertEnabled != nil && *props.ClientCertEnabled {
			clientCertMode = string(props.ClientCertMode)
		}
		d.Set("client_cert_mode", clientCertMode)

		if props.KeyVaultReferenceIdentity != nil {
			d.Set("key_vault_reference_identity_id", props.KeyVaultReferenceIdentity)
		}
	}

	appServiceTier, err := getFunctionAppServiceTier(ctx, appServicePlanID, meta)
	if err != nil {
		return err
	}

	appSettings := flattenAppServiceAppSettings(appSettingsResp.Properties)

	connectionString := appSettings["AzureWebJobsStorage"]

	// This teases out the necessary attributes from the storage connection string
	connectionStringParts := strings.Split(connectionString, ";")
	for _, part := range connectionStringParts {
		if strings.HasPrefix(part, "AccountName") {
			accountNameParts := strings.Split(part, "AccountName=")
			if len(accountNameParts) > 1 {
				d.Set("storage_account_name", accountNameParts[1])
			}
		}
		if strings.HasPrefix(part, "AccountKey") {
			accountKeyParts := strings.Split(part, "AccountKey=")
			if len(accountKeyParts) > 1 {
				d.Set("storage_account_access_key", accountKeyParts[1])
			}
		}
	}

	d.Set("version", appSettings["FUNCTIONS_EXTENSION_VERSION"])

	dashboard, ok := appSettings["AzureWebJobsDashboard"]
	d.Set("enable_builtin_logging", ok && dashboard != "")

	if _, ok = d.GetOk("app_settings.AzureWebJobsDashboard"); !ok {
		delete(appSettings, "AzureWebJobsDashboard")
	}
	if _, ok = d.GetOk("app_settings.AzureWebJobsStorage"); !ok {
		delete(appSettings, "AzureWebJobsStorage")
	}
	if _, ok = d.GetOk("app_settings.FUNCTIONS_EXTENSION_VERSION"); !ok {
		delete(appSettings, "FUNCTIONS_EXTENSION_VERSION")
	}

	// From the docs:
	// Only used when deploying to a Premium plan or to a Consumption plan running on Windows. Not supported for Consumptions plans running Linux.
	if (strings.EqualFold(appServiceTier, "dynamic") && strings.EqualFold(d.Get("os_type").(string), "linux")) ||
		(strings.EqualFold(appServiceTier, "dynamic") || strings.HasPrefix(strings.ToLower(appServiceTier), "elastic")) {
		delete(appSettings, "WEBSITE_CONTENTSHARE")
		delete(appSettings, "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING")
	}

	if err = d.Set("app_settings", appSettings); err != nil {
		return err
	}
	if err = d.Set("connection_string", flattenFunctionAppConnectionStrings(connectionStringsResp.Properties)); err != nil {
		return err
	}

	identity, err := flattenAppServiceIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %s", err)
	}

	configResp, err := client.GetConfiguration(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("making Read request on AzureRM Function App Configuration %q: %+v", id.SiteName, err)
	}

	siteConfig := flattenFunctionAppSiteConfig(configResp.SiteConfig)
	if err = d.Set("site_config", siteConfig); err != nil {
		return err
	}

	authSettings := flattenAppServiceAuthSettings(authResp.SiteAuthSettingsProperties)
	if err := d.Set("auth_settings", authSettings); err != nil {
		return fmt.Errorf("setting `auth_settings`: %s", err)
	}

	scmResp, err := client.GetSourceControl(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("making Read request on Function App Source Control %q: %+v", id.SiteName, err)
	}
	scm := flattenAppServiceSourceControl(scmResp.SiteSourceControlProperties)
	if err := d.Set("source_control", scm); err != nil {
		return fmt.Errorf("setting `source_control`: %s", err)
	}

	siteCred := flattenFunctionAppSiteCredential(siteCredResp.UserProperties)
	if err = d.Set("site_credential", siteCred); err != nil {
		return err
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceFunctionAppDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FunctionAppID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Function App %q (resource group %q)", id.SiteName, id.ResourceGroup)

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
