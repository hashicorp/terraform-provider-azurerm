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
	webValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceSlot() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceSlotCreateUpdate,
		Read:   resourceArmAppServiceSlotRead,
		Update: resourceArmAppServiceSlotCreateUpdate,
		Delete: resourceArmAppServiceSlotDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AppServiceSlotID(id)
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
				ValidateFunc: webValidate.AppServiceName,
			},

			"identity": schemaAppServiceIdentity(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"app_service_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: webValidate.AppServiceName,
			},

			"app_service_plan_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"site_config": schemaAppServiceSiteConfig(),

			"auth_settings": schemaAppServiceAuthSettings(),

			"logs": schemaAppServiceLogsConfig(),

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

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"app_settings": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
		},
	}
}

func resourceArmAppServiceSlotCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	slot := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	appServiceName := d.Get("app_service_name").(string)

	if d.IsNewResource() {
		existing, err := client.GetSlot(ctx, resourceGroup, appServiceName, slot)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_slot", *existing.ID)
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
		return fmt.Errorf("Error expanding `site_config` for App Service Slot %q (Resource Group %q): %s", slot, resourceGroup, err)
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

	if _, ok := d.GetOk("identity"); ok {
		appServiceIdentityRaw := d.Get("identity").([]interface{})
		appServiceIdentity := expandAppServiceIdentity(appServiceIdentityRaw)
		siteEnvelope.Identity = appServiceIdentity
	}

	createFuture, err := client.CreateOrUpdateSlot(ctx, resourceGroup, appServiceName, siteEnvelope, slot)
	if err != nil {
		return fmt.Errorf("Error creating Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}

	err = createFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation of Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}

	read, err := client.GetSlot(ctx, resourceGroup, appServiceName, slot)
	if err != nil {
		return fmt.Errorf("Error retrieving Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for Slot %q (App Service %q / Resource Group %q) ID", slot, appServiceName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServiceSlotUpdate(d, meta)
}

func resourceArmAppServiceSlotUpdate(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error expanding `site_config` for App Service Slot %q (Resource Group %q): %s", id.SlotName, id.ResourceGroup, err)
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
	createFuture, err := client.CreateOrUpdateSlot(ctx, id.ResourceGroup, id.SiteName, siteEnvelope, id.SlotName)
	if err != nil {
		return fmt.Errorf("Error updating Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	err = createFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for update of Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	if d.HasChange("site_config") {
		// update the main configuration
		siteConfig, err := expandAppServiceSiteConfig(d.Get("site_config"))
		if err != nil {
			return fmt.Errorf("Error expanding `site_config` for App Service Slot %q (Resource Group %q): %s", id.SlotName, id.ResourceGroup, err)
		}
		siteConfigResource := web.SiteConfigResource{
			SiteConfig: siteConfig,
		}
		if _, err := client.CreateOrUpdateConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, siteConfigResource, id.SlotName); err != nil {
			return fmt.Errorf("Error updating Configuration for App Service Slot %q/%q: %+v", id.SiteName, id.SlotName, err)
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
			return fmt.Errorf("Error updating Authentication Settings for App Service %q: %+v", id.SiteName, err)
		}
	}

	if d.HasChange("app_settings") {
		// update the AppSettings
		appSettings := expandAppServiceAppSettings(d)
		settings := web.StringDictionary{
			Properties: appSettings,
		}

		if _, err := client.UpdateApplicationSettingsSlot(ctx, id.ResourceGroup, id.SiteName, settings, id.SlotName); err != nil {
			return fmt.Errorf("Error updating Application Settings for App Service Slot %q/%q: %+v", id.SiteName, id.SlotName, err)
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
			return fmt.Errorf("Error updating Diagnostics Logs for App Service Slot %q/%q: %+v", id.SiteName, id.SlotName, err)
		}
	}

	if d.HasChange("connection_string") {
		// update the ConnectionStrings
		connectionStrings := expandAppServiceConnectionStrings(d)
		properties := web.ConnectionStringDictionary{
			Properties: connectionStrings,
		}

		if _, err := client.UpdateConnectionStringsSlot(ctx, id.ResourceGroup, id.SiteName, properties, id.SlotName); err != nil {
			return fmt.Errorf("Error updating Connection Strings for App Service Slot %q/%q: %+v", id.SiteName, id.SlotName, err)
		}
	}

	if d.HasChange("identity") {
		appServiceIdentityRaw := d.Get("identity").([]interface{})
		appServiceIdentity := expandAppServiceIdentity(appServiceIdentityRaw)
		sitePatchResource := web.SitePatchResource{
			ID:       utils.String(d.Id()),
			Identity: appServiceIdentity,
		}
		_, err := client.UpdateSlot(ctx, id.ResourceGroup, id.SiteName, sitePatchResource, id.SlotName)
		if err != nil {
			return fmt.Errorf("Error updating Managed Service Identity for App Service Slot %q/%q: %+v", id.SiteName, id.SlotName, err)
		}
	}

	return resourceArmAppServiceSlotRead(d, meta)
}

func resourceArmAppServiceSlotRead(d *schema.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("Error reading Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	configResp, err := client.GetConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		if utils.ResponseWasNotFound(configResp.Response) {
			log.Printf("[DEBUG] Configuration for Slot %q (App Service %q / Resource Group %q) were not found - removing from state!", id.SlotName, id.SiteName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Configuration for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	authResp, err := client.GetAuthSettingsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		return fmt.Errorf("Error reading Auth Settings for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	logsResp, err := client.GetDiagnosticLogsConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		return fmt.Errorf("Error retrieving the DiagnosticsLogsConfiguration for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	appSettingsResp, err := client.ListApplicationSettingsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		if utils.ResponseWasNotFound(appSettingsResp.Response) {
			log.Printf("[DEBUG] App Settings for Slot %q (App Service %q / Resource Group %q) were not found - removing from state!", id.SlotName, id.SiteName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading App Settings for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	connectionStringsResp, err := client.ListConnectionStringsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		return fmt.Errorf("Error listing Connection Strings for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	siteCredFuture, err := client.ListPublishingCredentialsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		return fmt.Errorf("Error retrieving publishing credentials for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}
	err = siteCredFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for publishing credentials for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}
	siteCredResp, err := siteCredFuture.Result(*client)
	if err != nil {
		return fmt.Errorf("Error reading publishing credentials for Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
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

	if err := d.Set("connection_string", flattenAppServiceConnectionStrings(connectionStringsResp.Properties)); err != nil {
		return fmt.Errorf("Error setting `connection_string`: %s", err)
	}

	authSettings := flattenAppServiceAuthSettings(authResp.SiteAuthSettingsProperties)
	if err := d.Set("auth_settings", authSettings); err != nil {
		return fmt.Errorf("Error setting `auth_settings`: %s", err)
	}

	logs := flattenAppServiceLogs(logsResp.SiteLogsConfigProperties)
	if err := d.Set("logs", logs); err != nil {
		return fmt.Errorf("Error setting `logs`: %s", err)
	}

	identity := flattenAppServiceIdentity(resp.Identity)
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("Error setting `identity`: %s", err)
	}

	siteCred := flattenAppServiceSiteCredential(siteCredResp.UserProperties)
	if err := d.Set("site_credential", siteCred); err != nil {
		return fmt.Errorf("Error setting `site_credential`: %s", err)
	}

	siteConfig := flattenAppServiceSiteConfig(configResp.SiteConfig)
	if err := d.Set("site_config", siteConfig); err != nil {
		return fmt.Errorf("Error setting `site_config`: %s", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmAppServiceSlotDelete(d *schema.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("Error deleting Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
		}
	}

	return nil
}
