package web

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	webValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// Azure Function App shares the same infrastructure with Azure App Service.
// So this resource will reuse most of the App Service code, but remove the configurations which are not applicable for Function App.
func resourceArmFunctionApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFunctionAppCreate,
		Read:   resourceArmFunctionAppRead,
		Update: resourceArmFunctionAppUpdate,
		Delete: resourceArmFunctionAppDelete,
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
				ValidateFunc: webValidate.AppServiceName,
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

			"client_affinity_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"daily_memory_time_quota": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"enable_builtin_logging": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"https_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"identity": schemaAppServiceIdentity(),

			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"linux",
				}, false),
			},

			"site_config": schemaAppServiceFunctionAppSiteConfig(),

			"source_control": schemaAppServiceSiteSourceControl(),

			"storage_account_name": {
				Type: schema.TypeString,
				// Required: true, // Uncomment this in 3.0
				Optional:      true,
				Computed:      true, // Remove this in 3.0
				ForceNew:      true,
				ValidateFunc:  storage.ValidateArmStorageAccountName,
				ConflictsWith: []string{"storage_connection_string"},
			},

			"storage_account_access_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true, // Remove this in 3.0
				// Required: true, // Uncomment this in 3.0
				Sensitive:     true,
				ValidateFunc:  validation.NoZeroValues,
				ConflictsWith: []string{"storage_connection_string"},
			},

			// TODO remove this in 3.0
			"storage_connection_string": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Sensitive:     true,
				Deprecated:    "Deprecated in favour of `storage_account_name` and `storage_account_access_key`",
				ConflictsWith: []string{"storage_account_name", "storage_account_access_key"},
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "~1",
			},

			"tags": tags.Schema(),

			// Computed Only

			"default_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kind": {
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
		},
	}
}

func resourceArmFunctionAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	endpointSuffix := meta.(*clients.Client).Account.Environment.StorageEndpointSuffix
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Function App creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing Function App %q (Resource Group %q): %s", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_function_app", *existing.ID)
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
		return fmt.Errorf("The name %q used for the Function App needs to be globally unique and isn't available: %s", name, *available.Message)
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
	clientAffinityEnabled := d.Get("client_affinity_enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	dailyMemoryTimeQuota := d.Get("daily_memory_time_quota").(int)
	t := d.Get("tags").(map[string]interface{})
	appServiceTier, err := getFunctionAppServiceTier(ctx, appServicePlanID, meta)
	if err != nil {
		return err
	}

	basicAppSettings, err := getBasicFunctionAppAppSettings(d, appServiceTier, endpointSuffix)
	if err != nil {
		return err
	}

	siteConfig, err := expandFunctionAppSiteConfig(d)
	if err != nil {
		return fmt.Errorf("Error expanding `site_config` for Function App %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	siteConfig.AppSettings = &basicAppSettings

	siteEnvelope := web.Site{
		Kind:     &kind,
		Location: &location,
		Tags:     tags.Expand(t),
		SiteProperties: &web.SiteProperties{
			ServerFarmID:          utils.String(appServicePlanID),
			Enabled:               utils.Bool(enabled),
			ClientAffinityEnabled: utils.Bool(clientAffinityEnabled),
			HTTPSOnly:             utils.Bool(httpsOnly),
			DailyMemoryTimeQuota:  utils.Int32(int32(dailyMemoryTimeQuota)),
			SiteConfig:            &siteConfig,
		},
	}

	if _, ok := d.GetOk("identity"); ok {
		appServiceIdentityRaw := d.Get("identity").([]interface{})
		appServiceIdentity := expandAppServiceIdentity(appServiceIdentityRaw)
		siteEnvelope.Identity = appServiceIdentity
	}

	createFuture, err := client.CreateOrUpdate(ctx, resourceGroup, name, siteEnvelope)
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
		_, err := client.CreateOrUpdateSourceControl(ctx, resourceGroup, name, *sourceControl)
		if err != nil {
			return fmt.Errorf("failed to create App Service Source Control for %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil || *read.ID == "" {
		return fmt.Errorf("Cannot read Function App %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	authSettingsRaw := d.Get("auth_settings").([]interface{})
	authSettings := expandAppServiceAuthSettings(authSettingsRaw)

	auth := web.SiteAuthSettings{
		ID:                         read.ID,
		SiteAuthSettingsProperties: &authSettings,
	}

	if _, err := client.UpdateAuthSettings(ctx, resourceGroup, name, auth); err != nil {
		return fmt.Errorf("Error updating auth settings for Function App %q (resource group %q): %+v", name, resourceGroup, err)
	}

	return resourceArmFunctionAppUpdate(d, meta)
}

func resourceArmFunctionAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	endpointSuffix := meta.(*clients.Client).Account.Environment.StorageEndpointSuffix
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceID(d.Id())
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
	clientAffinityEnabled := d.Get("client_affinity_enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	dailyMemoryTimeQuota := d.Get("daily_memory_time_quota").(int)
	t := d.Get("tags").(map[string]interface{})

	appServiceTier, err := getFunctionAppServiceTier(ctx, appServicePlanID, meta)
	if err != nil {
		return err
	}

	basicAppSettings, err := getBasicFunctionAppAppSettings(d, appServiceTier, endpointSuffix)
	if err != nil {
		return err
	}

	siteConfig, err := expandFunctionAppSiteConfig(d)
	if err != nil {
		return fmt.Errorf("Error expanding `site_config` for Function App %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
	}

	siteConfig.AppSettings = &basicAppSettings

	siteEnvelope := web.Site{
		Kind:     &kind,
		Location: &location,
		Tags:     tags.Expand(t),
		SiteProperties: &web.SiteProperties{
			ServerFarmID:          utils.String(appServicePlanID),
			Enabled:               utils.Bool(enabled),
			ClientAffinityEnabled: utils.Bool(clientAffinityEnabled),
			HTTPSOnly:             utils.Bool(httpsOnly),
			DailyMemoryTimeQuota:  utils.Int32(int32(dailyMemoryTimeQuota)),
			SiteConfig:            &siteConfig,
		},
	}

	if _, ok := d.GetOk("identity"); ok {
		appServiceIdentityRaw := d.Get("identity").([]interface{})
		appServiceIdentity := expandAppServiceIdentity(appServiceIdentityRaw)
		siteEnvelope.Identity = appServiceIdentity
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, siteEnvelope)
	if err != nil {
		return fmt.Errorf("Error updating Function App %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Function App %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	appSettings, err := expandFunctionAppAppSettings(d, appServiceTier, endpointSuffix)
	if err != nil {
		return err
	}
	settings := web.StringDictionary{
		Properties: appSettings,
	}

	if _, err = client.UpdateApplicationSettings(ctx, id.ResourceGroup, id.Name, settings); err != nil {
		return fmt.Errorf("Error updating Application Settings for Function App %q: %+v", id.Name, err)
	}

	// If `source_control` is defined, we need to set site_config.0.scm_type to "None" or we cannot update it
	// repo_url is required by the API
	_, hasSourceControl := d.GetOk("source_control.0.repo_url")

	scmType := web.ScmTypeNone

	if d.HasChange("site_config") || hasSourceControl {
		siteConfig, err := expandFunctionAppSiteConfig(d)
		if err != nil {
			return fmt.Errorf("Error expanding `site_config` for Function App %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
		}
		siteConfigResource := web.SiteConfigResource{
			SiteConfig: &siteConfig,
		}

		scmType = siteConfig.ScmType
		// ScmType being set blocks the update of source_control in _most_ cases, ADO is an exception
		if hasSourceControl && scmType != web.ScmTypeVSTSRM {
			siteConfigResource.SiteConfig.ScmType = web.ScmTypeNone
		}

		if _, err := client.CreateOrUpdateConfiguration(ctx, id.ResourceGroup, id.Name, siteConfigResource); err != nil {
			return fmt.Errorf("Error updating Configuration for Function App %q: %+v", id.Name, err)
		}
	}

	// Don't send source_control changes for ADO controlled Apps
	if hasSourceControl && scmType != web.ScmTypeVSTSRM {
		sourceControlProperties := expandAppServiceSiteSourceControl(d)
		sourceControl := &web.SiteSourceControl{}
		sourceControl.SiteSourceControlProperties = sourceControlProperties
		scFuture, err := client.CreateOrUpdateSourceControl(ctx, id.ResourceGroup, id.Name, *sourceControl)
		if err != nil {
			return fmt.Errorf("failed to create App Service Source Control for %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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
			return fmt.Errorf("Error updating Authentication Settings for Function App %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("connection_string") {
		// update the ConnectionStrings
		connectionStrings := expandFunctionAppConnectionStrings(d)
		properties := web.ConnectionStringDictionary{
			Properties: connectionStrings,
		}

		if _, err := client.UpdateConnectionStrings(ctx, id.ResourceGroup, id.Name, properties); err != nil {
			return fmt.Errorf("Error updating Connection Strings for App Service %q: %+v", id.Name, err)
		}
	}

	return resourceArmFunctionAppRead(d, meta)
}

func resourceArmFunctionAppRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] Function App %q (resource group %q) was not found - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Function App %q: %+v", id.Name, err)
	}

	appSettingsResp, err := client.ListApplicationSettings(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(appSettingsResp.Response) {
			log.Printf("[DEBUG] Application Settings of Function App %q (resource group %q) were not found", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Function App AppSettings %q: %+v", id.Name, err)
	}

	connectionStringsResp, err := client.ListConnectionStrings(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM Function App ConnectionStrings %q: %+v", id.Name, err)
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
	authResp, err := client.GetAuthSettings(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving the AuthSettings for Function App %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
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

	if props := resp.SiteProperties; props != nil {
		d.Set("app_service_plan_id", props.ServerFarmID)
		d.Set("enabled", props.Enabled)
		d.Set("default_hostname", props.DefaultHostName)
		d.Set("https_only", props.HTTPSOnly)
		d.Set("daily_memory_time_quota", props.DailyMemoryTimeQuota)
		d.Set("outbound_ip_addresses", props.OutboundIPAddresses)
		d.Set("possible_outbound_ip_addresses", props.PossibleOutboundIPAddresses)
		d.Set("client_affinity_enabled", props.ClientAffinityEnabled)
	}

	appSettings := flattenAppServiceAppSettings(appSettingsResp.Properties)

	connectionString := appSettings["AzureWebJobsStorage"]
	d.Set("storage_connection_string", connectionString)

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

	// Let the user have a final say whether they want to keep these 2 or not on
	// Linux consumption plans. They shouldn't be there according to a bug
	// report (see // https://github.com/Azure/azure-functions-python-worker/issues/598)
	if !strings.EqualFold(d.Get("os_type").(string), "linux") {
		delete(appSettings, "WEBSITE_CONTENTSHARE")
		delete(appSettings, "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING")
	}

	if err = d.Set("app_settings", appSettings); err != nil {
		return err
	}
	if err = d.Set("connection_string", flattenFunctionAppConnectionStrings(connectionStringsResp.Properties)); err != nil {
		return err
	}

	identity := flattenAppServiceIdentity(resp.Identity)
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("Error setting `identity`: %s", err)
	}

	configResp, err := client.GetConfiguration(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM Function App Configuration %q: %+v", id.Name, err)
	}

	siteConfig := flattenFunctionAppSiteConfig(configResp.SiteConfig)
	if err = d.Set("site_config", siteConfig); err != nil {
		return err
	}

	authSettings := flattenAppServiceAuthSettings(authResp.SiteAuthSettingsProperties)
	if err := d.Set("auth_settings", authSettings); err != nil {
		return fmt.Errorf("Error setting `auth_settings`: %s", err)
	}

	scmResp, err := client.GetSourceControl(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Function App Source Control %q: %+v", id.Name, err)
	}
	scm := flattenAppServiceSourceControl(scmResp.SiteSourceControlProperties)
	if err := d.Set("source_control", scm); err != nil {
		return fmt.Errorf("Error setting `source_control`: %s", err)
	}

	siteCred := flattenFunctionAppSiteCredential(siteCredResp.UserProperties)
	if err = d.Set("site_credential", siteCred); err != nil {
		return err
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmFunctionAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Function App %q (resource group %q)", id.Name, id.ResourceGroup)

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
