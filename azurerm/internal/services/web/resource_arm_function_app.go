package web

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
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
			_, err := ParseAppServiceID(id)
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

			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"app_service_plan_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "~1",
			},

			// TODO remove this in 3.0
			"storage_connection_string": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Sensitive:     true,
				Deprecated:    "Deprecated in favor of `storage_account_name` and `storage_account_access_key`",
				ConflictsWith: []string{"storage_account_name", "storage_account_access_key"},
			},

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

			"app_settings": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"enable_builtin_logging": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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

			"identity": azure.SchemaAppServiceIdentity(),

			"tags": tags.Schema(),

			"default_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"linux",
				}, false),
			},

			"outbound_ip_addresses": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"possible_outbound_ip_addresses": {
				Type:     schema.TypeString,
				Computed: true,
			},

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

			"daily_memory_time_quota": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"site_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"always_on": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"use_32_bit_worker_process": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"websockets_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"linux_fx_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"http2_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"ip_restriction": {
							Type:       schema.TypeList,
							Optional:   true,
							Computed:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validate.CIDR,
									},
									"subnet_id": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
						"min_tls_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(web.OneFullStopZero),
								string(web.OneFullStopOne),
								string(web.OneFullStopTwo),
							}, false),
						},
						"ftps_state": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(web.AllAllowed),
								string(web.Disabled),
								string(web.FtpsOnly),
							}, false),
						},
						"pre_warmed_instance_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 10),
						},
						"cors": azure.SchemaWebCorsSettings(),
					},
				},
			},

			"auth_settings": azure.SchemaAppServiceAuthSettings(),

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

	if features.ShouldResourcesBeImported() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Function App %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_function_app", *existing.ID)
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
		appServiceIdentity := azure.ExpandAppServiceIdentity(appServiceIdentityRaw)
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

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Function App %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	authSettingsRaw := d.Get("auth_settings").([]interface{})
	authSettings := azure.ExpandAppServiceAuthSettings(authSettingsRaw)

	auth := web.SiteAuthSettings{
		ID:                         read.ID,
		SiteAuthSettingsProperties: &authSettings}

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

	id, err := ParseAppServiceID(d.Id())
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
		appServiceIdentity := azure.ExpandAppServiceIdentity(appServiceIdentityRaw)
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

	if d.HasChange("site_config") {
		siteConfig, err := expandFunctionAppSiteConfig(d)
		if err != nil {
			return fmt.Errorf("Error expanding `site_config` for Function App %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
		}
		siteConfigResource := web.SiteConfigResource{
			SiteConfig: &siteConfig,
		}
		if _, err := client.CreateOrUpdateConfiguration(ctx, id.ResourceGroup, id.Name, siteConfigResource); err != nil {
			return fmt.Errorf("Error updating Configuration for Function App %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("auth_settings") {
		authSettingsRaw := d.Get("auth_settings").([]interface{})
		authSettingsProperties := azure.ExpandAppServiceAuthSettings(authSettingsRaw)
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

	id, err := ParseAppServiceID(d.Id())
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

	// Let the user have a final say whether he wants to keep these 2 or not on
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

	identity := azure.FlattenAppServiceIdentity(resp.Identity)
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

	authSettings := azure.FlattenAppServiceAuthSettings(authResp.SiteAuthSettingsProperties)
	if err := d.Set("auth_settings", authSettings); err != nil {
		return fmt.Errorf("Error setting `auth_settings`: %s", err)
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

	id, err := ParseAppServiceID(d.Id())
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

func getBasicFunctionAppAppSettings(d *schema.ResourceData, appServiceTier, endpointSuffix string) ([]web.NameValuePair, error) {
	// TODO: This is a workaround since there are no public Functions API
	// You may track the API request here: https://github.com/Azure/azure-rest-api-specs/issues/3750
	dashboardPropName := "AzureWebJobsDashboard"
	storagePropName := "AzureWebJobsStorage"
	functionVersionPropName := "FUNCTIONS_EXTENSION_VERSION"

	contentSharePropName := "WEBSITE_CONTENTSHARE"
	contentFileConnStringPropName := "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"

	// TODO 3.0 - remove this logic for determining which storage account connection string to use
	storageConnection := ""
	if v, ok := d.GetOk("storage_connection_string"); ok {
		storageConnection = v.(string)
	}

	storageAccount := ""
	if v, ok := d.GetOk("storage_account_name"); ok {
		storageAccount = v.(string)
	}

	connectionString := ""
	if v, ok := d.GetOk("storage_account_access_key"); ok {
		connectionString = v.(string)
	}

	if storageConnection == "" && storageAccount == "" && connectionString == "" {
		return nil, fmt.Errorf("one of `storage_connection_string` or `storage_account_name` and `storage_account_access_key` must be specified")
	}

	if (storageAccount == "" && connectionString != "") || (storageAccount != "" && connectionString == "") {
		return nil, fmt.Errorf("both `storage_account_name` and `storage_account_access_key` must be specified")
	}

	if connectionString != "" && storageAccount != "" {
		storageConnection = fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", storageAccount, connectionString, endpointSuffix)
	}

	functionVersion := d.Get("version").(string)
	contentShare := strings.ToLower(d.Get("name").(string)) + "-content"

	basicSettings := []web.NameValuePair{
		{Name: &storagePropName, Value: &storageConnection},
		{Name: &functionVersionPropName, Value: &functionVersion},
	}

	if d.Get("enable_builtin_logging").(bool) {
		basicSettings = append(basicSettings, web.NameValuePair{
			Name:  &dashboardPropName,
			Value: &storageConnection,
		})
	}

	consumptionSettings := []web.NameValuePair{
		{Name: &contentSharePropName, Value: &contentShare},
		{Name: &contentFileConnStringPropName, Value: &storageConnection},
	}

	// On consumption and premium plans include WEBSITE_CONTENT components, unless it's a Linux consumption plan
	// (see https://github.com/Azure/azure-functions-python-worker/issues/598)
	if (strings.EqualFold(appServiceTier, "dynamic") || strings.EqualFold(appServiceTier, "elasticpremium")) &&
		!strings.EqualFold(d.Get("os_type").(string), "linux") {
		return append(basicSettings, consumptionSettings...), nil
	}

	return basicSettings, nil
}

func getFunctionAppServiceTier(ctx context.Context, appServicePlanId string, meta interface{}) (string, error) {
	id, err := ParseAppServicePlanID(appServicePlanId)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Unable to parse App Service Plan ID %q: %+v", appServicePlanId, err)
	}

	log.Printf("[DEBUG] Retrieving App Service Plan %q (Resource Group %q)", id.Name, id.ResourceGroup)

	appServicePlansClient := meta.(*clients.Client).Web.AppServicePlansClient
	appServicePlan, err := appServicePlansClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Could not retrieve App Service Plan ID %q: %+v", appServicePlanId, err)
	}

	if sku := appServicePlan.Sku; sku != nil {
		if tier := sku.Tier; tier != nil {
			return *tier, nil
		}
	}
	return "", fmt.Errorf("No `sku` block was returned for App Service Plan ID %q", appServicePlanId)
}

func expandFunctionAppAppSettings(d *schema.ResourceData, appServiceTier, endpointSuffix string) (map[string]*string, error) {
	output := expandAppServiceAppSettings(d)

	basicAppSettings, err := getBasicFunctionAppAppSettings(d, appServiceTier, endpointSuffix)
	if err != nil {
		return nil, err
	}
	for _, p := range basicAppSettings {
		output[*p.Name] = p.Value
	}

	return output, nil
}

func expandFunctionAppSiteConfig(d *schema.ResourceData) (web.SiteConfig, error) {
	configs := d.Get("site_config").([]interface{})
	siteConfig := web.SiteConfig{}

	if len(configs) == 0 {
		return siteConfig, nil
	}

	config := configs[0].(map[string]interface{})

	if v, ok := config["always_on"]; ok {
		siteConfig.AlwaysOn = utils.Bool(v.(bool))
	}

	if v, ok := config["use_32_bit_worker_process"]; ok {
		siteConfig.Use32BitWorkerProcess = utils.Bool(v.(bool))
	}

	if v, ok := config["websockets_enabled"]; ok {
		siteConfig.WebSocketsEnabled = utils.Bool(v.(bool))
	}

	if v, ok := config["linux_fx_version"]; ok {
		siteConfig.LinuxFxVersion = utils.String(v.(string))
	}

	if v, ok := config["cors"]; ok {
		corsSettings := v.(interface{})
		expand := azure.ExpandWebCorsSettings(corsSettings)
		siteConfig.Cors = &expand
	}

	if v, ok := config["http2_enabled"]; ok {
		siteConfig.HTTP20Enabled = utils.Bool(v.(bool))
	}

	if v, ok := config["ip_restriction"]; ok {
		ipSecurityRestrictions := v.(interface{})
		restrictions, err := expandFunctionAppIpRestriction(ipSecurityRestrictions)
		if err != nil {
			return siteConfig, err
		}
		siteConfig.IPSecurityRestrictions = &restrictions
	}

	if v, ok := config["min_tls_version"]; ok {
		siteConfig.MinTLSVersion = web.SupportedTLSVersions(v.(string))
	}

	if v, ok := config["ftps_state"]; ok {
		siteConfig.FtpsState = web.FtpsState(v.(string))
	}

	if v, ok := config["pre_warmed_instance_count"]; ok {
		siteConfig.PreWarmedInstanceCount = utils.Int32(int32(v.(int)))
	}

	return siteConfig, nil
}

func flattenFunctionAppSiteConfig(input *web.SiteConfig) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{})

	if input == nil {
		log.Printf("[DEBUG] SiteConfig is nil")
		return results
	}

	if input.AlwaysOn != nil {
		result["always_on"] = *input.AlwaysOn
	}

	if input.Use32BitWorkerProcess != nil {
		result["use_32_bit_worker_process"] = *input.Use32BitWorkerProcess
	}

	if input.WebSocketsEnabled != nil {
		result["websockets_enabled"] = *input.WebSocketsEnabled
	}

	if input.LinuxFxVersion != nil {
		result["linux_fx_version"] = *input.LinuxFxVersion
	}

	if input.HTTP20Enabled != nil {
		result["http2_enabled"] = *input.HTTP20Enabled
	}

	if input.PreWarmedInstanceCount != nil {
		result["pre_warmed_instance_count"] = *input.PreWarmedInstanceCount
	}

	result["ip_restriction"] = flattenFunctionAppIpRestriction(input.IPSecurityRestrictions)

	result["min_tls_version"] = string(input.MinTLSVersion)
	result["ftps_state"] = string(input.FtpsState)

	result["cors"] = azure.FlattenWebCorsSettings(input.Cors)

	results = append(results, result)
	return results
}

func expandFunctionAppConnectionStrings(d *schema.ResourceData) map[string]*web.ConnStringValueTypePair {
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

func expandFunctionAppIpRestriction(input interface{}) ([]web.IPSecurityRestriction, error) {
	restrictions := make([]web.IPSecurityRestriction, 0)

	for i, r := range input.([]interface{}) {
		if r == nil {
			continue
		}

		restriction := r.(map[string]interface{})

		ipAddress := restriction["ip_address"].(string)
		vNetSubnetID := restriction["subnet_id"].(string)

		if vNetSubnetID != "" && ipAddress != "" {
			return nil, fmt.Errorf(fmt.Sprintf("only one of `ip_address` or `subnet_id` can set for `site_config.0.ip_restriction.%d`", i))
		}

		if vNetSubnetID == "" && ipAddress == "" {
			return nil, fmt.Errorf(fmt.Sprintf("one of `ip_address` or `subnet_id` must be set for `site_config.0.ip_restriction.%d`", i))
		}

		ipSecurityRestriction := web.IPSecurityRestriction{}
		if ipAddress == "Any" {
			continue
		}

		if ipAddress != "" {
			ipSecurityRestriction.IPAddress = &ipAddress
		}

		if vNetSubnetID != "" {
			ipSecurityRestriction.VnetSubnetResourceID = &vNetSubnetID
		}

		restrictions = append(restrictions, ipSecurityRestriction)
	}

	return restrictions, nil
}

func flattenFunctionAppConnectionStrings(input map[string]*web.ConnStringValueTypePair) interface{} {
	results := make([]interface{}, 0)

	for k, v := range input {
		result := make(map[string]interface{})
		result["name"] = k
		result["type"] = string(v.Type)
		result["value"] = *v.Value
		results = append(results, result)
	}

	return results
}

func flattenFunctionAppSiteCredential(input *web.UserProperties) []interface{} {
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

func flattenFunctionAppIpRestriction(input *[]web.IPSecurityRestriction) []interface{} {
	restrictions := make([]interface{}, 0)

	if input == nil {
		return restrictions
	}

	for _, v := range *input {
		ipAddress := ""
		if v.IPAddress != nil {
			ipAddress = *v.IPAddress
			if ipAddress == "Any" {
				continue
			}
		}

		subnetId := ""
		if v.VnetSubnetResourceID != nil {
			subnetId = *v.VnetSubnetResourceID
		}

		restrictions = append(restrictions, map[string]interface{}{
			"ip_address": ipAddress,
			"subnet_id":  subnetId,
		})
	}

	return restrictions
}
