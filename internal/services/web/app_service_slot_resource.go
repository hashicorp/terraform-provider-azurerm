package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	webValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAppServiceSlot() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceSlotCreateUpdate,
		Read:   resourceAppServiceSlotRead,
		Update: resourceAppServiceSlotCreateUpdate,
		Delete: resourceAppServiceSlotDelete,

		DeprecationMessage: "The `azurerm_app_service_slot` resource has been superseded by the `azurerm_linux_web_app_slot` and `azurerm_windows_web_app_slot` resources. Whilst this resource will continue to be available in the 2.x and 3.x releases it is feature-frozen for compatibility purposes, will no longer receive any updates and will be removed in a future major release of the Azure Provider.",

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AppServiceSlotID(id)
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
				ValidateFunc: webValidate.AppServiceName,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"app_service_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: webValidate.AppServiceName,
			},

			"app_service_plan_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"site_config": schemaAppServiceSiteConfig(),

			"storage_account": schemaAppServiceStorageAccounts(),

			"auth_settings": schemaAppServiceAuthSettings(),

			"key_vault_reference_identity_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			},

			"logs": schemaAppServiceLogsConfig(),

			"client_affinity_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"https_only": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"app_settings": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
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
						"value": {
							Type:      pluginsdk.TypeString,
							Required:  true,
							Sensitive: true,
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
					},
				},
			},

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

			"default_site_hostname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppServiceSlotCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewAppServiceSlotID(subscriptionId, d.Get("resource_group_name").(string), d.Get("app_service_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.GetSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_app_service_slot", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	appServicePlanId := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	t := d.Get("tags").(map[string]interface{})
	affinity := d.Get("client_affinity_enabled").(bool)

	siteConfig, err := expandAppServiceSiteConfig(d.Get("site_config"))
	if err != nil {
		return fmt.Errorf("expanding `site_config` for %s: %s", id, err)
	}
	siteEnvelope := web.Site{
		Location: &location,
		Tags:     tags.Expand(t),
		SiteProperties: &web.SiteProperties{
			ServerFarmID:          utils.String(appServicePlanId),
			Enabled:               utils.Bool(enabled),
			HTTPSOnly:             utils.Bool(httpsOnly),
			SiteConfig:            siteConfig,
			ClientAffinityEnabled: &affinity,
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

	createFuture, err := client.CreateOrUpdateSlot(ctx, id.ResourceGroup, id.SiteName, siteEnvelope, id.SlotName)
	if err != nil {
		return fmt.Errorf("creating %s: %s", id, err)
	}

	err = createFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting for creation of %s: %s", id, err)
	}

	d.SetId(id.ID())

	return resourceAppServiceSlotUpdate(d, meta)
}

func resourceAppServiceSlotUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceSlotID(d.Id())
	if err != nil {
		return err
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	appServicePlanId := d.Get("app_service_plan_id").(string)
	siteConfig, err := expandAppServiceSiteConfig(d.Get("site_config"))
	if err != nil {
		return fmt.Errorf("expanding `site_config` for App Service Slot %q (Resource Group %q): %s", id.SlotName, id.ResourceGroup, err)
	}
	enabled := d.Get("enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	t := d.Get("tags").(map[string]interface{})

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
	if v, ok := d.GetOk("client_affinity_enabled"); ok {
		enabled := v.(bool)
		siteEnvelope.SiteProperties.ClientAffinityEnabled = utils.Bool(enabled)
	}

	if v, ok := d.GetOk("key_vault_reference_identity_id"); ok {
		siteEnvelope.SiteProperties.KeyVaultReferenceIdentity = utils.String(v.(string))
	}

	createFuture, err := client.CreateOrUpdateSlot(ctx, id.ResourceGroup, id.SiteName, siteEnvelope, id.SlotName)
	if err != nil {
		return fmt.Errorf("updating Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	err = createFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting for update of Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	if d.HasChange("site_config") {
		// update the main configuration
		siteConfig, err := expandAppServiceSiteConfig(d.Get("site_config"))
		if err != nil {
			return fmt.Errorf("expanding `site_config` for App Service Slot %q (Resource Group %q): %s", id.SlotName, id.ResourceGroup, err)
		}
		siteConfigResource := web.SiteConfigResource{
			SiteConfig: siteConfig,
		}
		if _, err := client.CreateOrUpdateConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, siteConfigResource, id.SlotName); err != nil {
			return fmt.Errorf("updating Configuration for App Service Slot %q/%q: %+v", id.SiteName, id.SlotName, err)
		}
	}

	if d.HasChange("auth_settings") {
		authSettingsRaw := d.Get("auth_settings").([]interface{})
		authSettingsProperties := expandAppServiceAuthSettings(authSettingsRaw)
		authSettings := web.SiteAuthSettings{
			ID:                         utils.String(d.Id()),
			SiteAuthSettingsProperties: &authSettingsProperties,
		}

		if _, err := client.UpdateAuthSettingsSlot(ctx, id.ResourceGroup, id.SiteName, authSettings, id.SlotName); err != nil {
			return fmt.Errorf("updating Authentication Settings for App Service %q: %+v", id.SiteName, err)
		}
	}

	if d.HasChange("app_settings") {
		// update the AppSettings
		appSettings := expandAppServiceAppSettings(d)
		settings := web.StringDictionary{
			Properties: appSettings,
		}

		if _, err := client.UpdateApplicationSettingsSlot(ctx, id.ResourceGroup, id.SiteName, settings, id.SlotName); err != nil {
			return fmt.Errorf("updating Application Settings for App Service Slot %q/%q: %+v", id.SiteName, id.SlotName, err)
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

		if _, err := client.UpdateDiagnosticLogsConfigSlot(ctx, id.ResourceGroup, id.SiteName, logsResource, id.SlotName); err != nil {
			return fmt.Errorf("updating Diagnostics Logs for App Service Slot %q/%q: %+v", id.SiteName, id.SlotName, err)
		}
	}

	if d.HasChange("storage_account") {
		storageAccountsRaw := d.Get("storage_account").(*pluginsdk.Set).List()
		storageAccounts := expandAppServiceStorageAccounts(storageAccountsRaw)
		properties := web.AzureStoragePropertyDictionaryResource{
			Properties: storageAccounts,
		}

		if _, err := client.UpdateAzureStorageAccountsSlot(ctx, id.ResourceGroup, id.SiteName, properties, id.SlotName); err != nil {
			return fmt.Errorf("updating Storage Accounts for App Service Slot %q/%q: %+v", id.SiteName, id.SlotName, err)
		}
	}

	if d.HasChange("connection_string") {
		// update the ConnectionStrings
		connectionStrings := expandAppServiceConnectionStrings(d)
		properties := web.ConnectionStringDictionary{
			Properties: connectionStrings,
		}

		if _, err := client.UpdateConnectionStringsSlot(ctx, id.ResourceGroup, id.SiteName, properties, id.SlotName); err != nil {
			return fmt.Errorf("updating Connection Strings for App Service Slot %q/%q: %+v", id.SiteName, id.SlotName, err)
		}
	}

	if d.HasChange("identity") {
		appServiceIdentity, err := expandAppServiceIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		sitePatchResource := web.SitePatchResource{
			ID:       utils.String(d.Id()),
			Identity: appServiceIdentity,
		}
		if _, err := client.UpdateSlot(ctx, id.ResourceGroup, id.SiteName, sitePatchResource, id.SlotName); err != nil {
			return fmt.Errorf("updating Managed Service Identity for App Service Slot %q/%q: %+v", id.SiteName, id.SlotName, err)
		}
	}

	return resourceAppServiceSlotRead(d, meta)
}

func resourceAppServiceSlotRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceSlotID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Slot %q (App Service %q / Resource Group %q) were not found - removing from state!", id.SlotName, id.SiteName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	configResp, err := client.GetConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		if utils.ResponseWasNotFound(configResp.Response) {
			log.Printf("[DEBUG] Configuration for Slot %q (App Service %q / Resource Group %q) were not found - removing from state!", id.SlotName, id.SiteName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Configuration for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	authResp, err := client.GetAuthSettingsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		return fmt.Errorf("reading Auth Settings for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	logsResp, err := client.GetDiagnosticLogsConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		return fmt.Errorf("retrieving the DiagnosticsLogsConfiguration for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	appSettingsResp, err := client.ListApplicationSettingsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		if utils.ResponseWasNotFound(appSettingsResp.Response) {
			log.Printf("[DEBUG] App Settings for Slot %q (App Service %q / Resource Group %q) were not found - removing from state!", id.SlotName, id.SiteName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading App Settings for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	storageAccountsResp, err := client.ListAzureStorageAccountsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		return fmt.Errorf("listing Storage Accounts for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	connectionStringsResp, err := client.ListConnectionStringsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		return fmt.Errorf("listing Connection Strings for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	siteCredFuture, err := client.ListPublishingCredentialsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		return fmt.Errorf("retrieving publishing credentials for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}
	err = siteCredFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting for publishing credentials for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}
	siteCredResp, err := siteCredFuture.Result(*client)
	if err != nil {
		return fmt.Errorf("reading publishing credentials for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	d.Set("name", id.SlotName)
	d.Set("app_service_name", id.SiteName)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.SiteProperties; props != nil {
		d.Set("app_service_plan_id", props.ServerFarmID)
		d.Set("client_affinity_enabled", props.ClientAffinityEnabled)
		d.Set("default_site_hostname", props.DefaultHostName)
		d.Set("enabled", props.Enabled)
		d.Set("https_only", props.HTTPSOnly)

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

	if err := d.Set("storage_account", flattenAppServiceStorageAccounts(storageAccountsResp.Properties)); err != nil {
		return fmt.Errorf("setting `storage_account`: %s", err)
	}

	if err := d.Set("connection_string", flattenAppServiceConnectionStrings(connectionStringsResp.Properties)); err != nil {
		return fmt.Errorf("setting `connection_string`: %s", err)
	}

	authSettings := flattenAppServiceAuthSettings(authResp.SiteAuthSettingsProperties)
	if err := d.Set("auth_settings", authSettings); err != nil {
		return fmt.Errorf("setting `auth_settings`: %s", err)
	}

	logs := flattenAppServiceLogs(logsResp.SiteLogsConfigProperties)
	if err := d.Set("logs", logs); err != nil {
		return fmt.Errorf("setting `logs`: %s", err)
	}

	identity, err := flattenAppServiceIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %s", err)
	}

	siteCred := flattenAppServiceSiteCredential(siteCredResp.UserProperties)
	if err := d.Set("site_credential", siteCred); err != nil {
		return fmt.Errorf("setting `site_credential`: %s", err)
	}

	siteConfig := flattenAppServiceSiteConfig(configResp.SiteConfig)
	if err := d.Set("site_config", siteConfig); err != nil {
		return fmt.Errorf("setting `site_config`: %s", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceAppServiceSlotDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceSlotID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Slot %q (App Service %q / Resource Group %q)", id.SlotName, id.SiteName, id.ResourceGroup)

	deleteMetrics := true
	deleteEmptyServerFarm := false
	resp, err := client.DeleteSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName, &deleteMetrics, &deleteEmptyServerFarm)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
		}
	}

	return nil
}
