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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	webValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmFunctionAppSlot() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFunctionAppSlotCreate,
		Read:   resourceArmFunctionAppSlotRead,
		Update: resourceArmFunctionAppSlotUpdate,
		Delete: resourceArmFunctionAppSlotDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.FunctionAppSlotID(id)
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"identity": schemaAppServiceIdentity(),

			"function_app_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: webValidate.AppServiceName,
			},

			"app_service_plan_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateAppServicePlanID,
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "~1",
			},

			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storage.ValidateArmStorageAccountName,
			},

			"storage_account_access_key": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.NoZeroValues,
			},

			"app_settings": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"linux",
				}, false),
			},

			"client_affinity_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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

			"site_config": schemaAppServiceFunctionAppSiteConfig(),

			"auth_settings": schemaAppServiceAuthSettings(),

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

			"tags": tags.Schema(),
		},
	}
}

func resourceArmFunctionAppSlotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	endpointSuffix := meta.(*clients.Client).Account.Environment.StorageEndpointSuffix
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Function App Slot creation.")

	slot := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	functionAppName := d.Get("function_app_name").(string)

	if d.IsNewResource() {
		existing, err := client.GetSlot(ctx, resourceGroup, functionAppName, slot)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Slot %q (Function App %q / Resource Group %q): %s", slot, functionAppName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_function_app_slot", *existing.ID)
		}
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
	appServiceTier, err := getFunctionAppSlotServiceTier(ctx, appServicePlanID, meta)
	if err != nil {
		return err
	}

	basicAppSettings := getBasicFunctionAppSlotAppSettings(d, appServiceTier, endpointSuffix)

	siteConfig, err := expandFunctionAppSiteConfig(d)
	if err != nil {
		return fmt.Errorf("Error expanding `site_config` for Function App Slot %q (Resource Group %q): %s", slot, resourceGroup, err)
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

	createFuture, err := client.CreateOrUpdateSlot(ctx, resourceGroup, functionAppName, siteEnvelope, slot)
	if err != nil {
		return err
	}

	err = createFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.GetSlot(ctx, resourceGroup, functionAppName, slot)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for Slot %q (Function App %q / Resource Group %q) ID", slot, functionAppName, resourceGroup)
	}

	d.SetId(*read.ID)

	authSettingsRaw := d.Get("auth_settings").([]interface{})
	authSettings := expandAppServiceAuthSettings(authSettingsRaw)

	auth := web.SiteAuthSettings{
		ID:                         read.ID,
		SiteAuthSettingsProperties: &authSettings,
	}

	if _, err := client.UpdateAuthSettingsSlot(ctx, resourceGroup, functionAppName, auth, slot); err != nil {
		return fmt.Errorf("Error updating auth settings for Slot %q (Function App Slot %q / Resource Group %q): %+s", slot, functionAppName, resourceGroup, err)
	}

	return resourceArmFunctionAppSlotUpdate(d, meta)
}

func resourceArmFunctionAppSlotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	endpointSuffix := meta.(*clients.Client).Account.Environment.StorageEndpointSuffix
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FunctionAppSlotID(d.Id())
	if err != nil {
		return err
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	kind := "functionapp"
	if osTypeRaw, ok := d.GetOk("os_type"); ok {
		osType := osTypeRaw.(string)
		if osType == "Linux" {
			kind = "functionapp,linux"
		}
	}
	appServicePlanID := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	clientAffinityEnabled := d.Get("client_affinity_enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	dailyMemoryTimeQuota := d.Get("daily_memory_time_quota").(int)
	t := d.Get("tags").(map[string]interface{})

	appServiceTier, err := getFunctionAppSlotServiceTier(ctx, appServicePlanID, meta)
	if err != nil {
		return err
	}

	basicAppSettings := getBasicFunctionAppSlotAppSettings(d, appServiceTier, endpointSuffix)

	siteConfig, err := expandFunctionAppSiteConfig(d)
	if err != nil {
		return fmt.Errorf("Error expanding `site_config` for Slot %q (Function App %q / Resource Group %q): %s", id.Name, id.FunctionAppName, id.ResourceGroup, err)
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

	future, err := client.CreateOrUpdateSlot(ctx, id.ResourceGroup, id.FunctionAppName, siteEnvelope, id.Name)
	if err != nil {
		return fmt.Errorf("Error updating Slot %q (Function App %q / Resource Group %q): %s", id.Name, id.FunctionAppName, id.ResourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for update of Slot %q (Function App %q / Resource Group %q): %s", id.Name, id.FunctionAppName, id.ResourceGroup, err)
	}

	appSettings, err := expandFunctionAppSlotAppSettings(d, appServiceTier, endpointSuffix)
	if err != nil {
		return err
	}
	settings := web.StringDictionary{
		Properties: appSettings,
	}

	if _, err = client.UpdateApplicationSettingsSlot(ctx, id.ResourceGroup, id.FunctionAppName, settings, id.Name); err != nil {
		return fmt.Errorf("Error updating Application Settings for Function App Slot %q (Function App %q / Resource Group %q): %+v", id.Name, id.FunctionAppName, id.ResourceGroup, err)
	}

	if d.HasChange("site_config") {
		siteConfig, err := expandFunctionAppSiteConfig(d)
		if err != nil {
			return fmt.Errorf("Error expanding `site_config` for Slot %q (Function App %q / Resource Group %q): %s", id.Name, id.FunctionAppName, id.ResourceGroup, err)
		}
		siteConfigResource := web.SiteConfigResource{
			SiteConfig: &siteConfig,
		}
		if _, err := client.CreateOrUpdateConfigurationSlot(ctx, id.ResourceGroup, id.FunctionAppName, siteConfigResource, id.Name); err != nil {
			return fmt.Errorf("Error updating Configuration for Slot %q (Function App %q / Resource Group %q): %+v", id.Name, id.FunctionAppName, id.ResourceGroup, err)
		}
	}

	if d.HasChange("auth_settings") {
		authSettingsRaw := d.Get("auth_settings").([]interface{})
		authSettingsProperties := expandAppServiceAuthSettings(authSettingsRaw)
		authSettings := web.SiteAuthSettings{
			ID:                         utils.String(d.Id()),
			SiteAuthSettingsProperties: &authSettingsProperties,
		}

		if _, err := client.UpdateAuthSettingsSlot(ctx, id.ResourceGroup, id.FunctionAppName, authSettings, id.Name); err != nil {
			return fmt.Errorf("Error updating Authentication Settings for Slot %q (Function App %q / Resource Group %q): %+v", id.Name, id.FunctionAppName, id.ResourceGroup, err)
		}
	}

	if d.HasChange("connection_string") {
		// update the ConnectionStrings
		connectionStrings := expandFunctionAppSlotConnectionStrings(d)
		properties := web.ConnectionStringDictionary{
			Properties: connectionStrings,
		}

		if _, err := client.UpdateConnectionStringsSlot(ctx, id.ResourceGroup, id.FunctionAppName, properties, id.Name); err != nil {
			return fmt.Errorf("Error updating Connection Strings for Slot %q (Function App %q / Resource Group %q): %+v", id.Name, id.FunctionAppName, id.ResourceGroup, err)
		}
	}

	return resourceArmFunctionAppSlotRead(d, meta)
}

func resourceArmFunctionAppSlotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FunctionAppSlotID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetSlot(ctx, id.ResourceGroup, id.FunctionAppName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Function App Slot %q (Function App %q / Resource Group %q) was not found - removing from state", id.Name, id.FunctionAppName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error makeing read request on AzureRM Function App Slot %q (Function App %q / Resource Group %q): %s", id.Name, id.FunctionAppName, id.ResourceGroup, err)
	}

	appSettingsResp, err := client.ListApplicationSettingsSlot(ctx, id.ResourceGroup, id.FunctionAppName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(appSettingsResp.Response) {
			log.Printf("[DEBUG] Application Settings of AzureRM Function App Slot %q (Function App %q / Resource Group %q) were not found", id.Name, id.FunctionAppName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Function App Slot %q (Function App %q / Resource Group %q) AppSettings: %+v", id.Name, id.FunctionAppName, id.ResourceGroup, err)
	}

	connectionStringsResp, err := client.ListConnectionStringsSlot(ctx, id.ResourceGroup, id.FunctionAppName, id.Name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM Function App Slot %q (Function App %q / Resource Group %q) ConnectionStrings: %+v", id.Name, id.FunctionAppName, id.ResourceGroup, err)
	}

	siteCredFuture, err := client.ListPublishingCredentialsSlot(ctx, id.ResourceGroup, id.FunctionAppName, id.Name)
	if err != nil {
		return err
	}
	err = siteCredFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return err
	}
	siteCredResp, err := siteCredFuture.Result(*client)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM Function App Slot %q (Function App %q / Resource Group %q) Site Credentials: %+v", id.Name, id.FunctionAppName, id.ResourceGroup, err)
	}
	authResp, err := client.GetAuthSettingsSlot(ctx, id.ResourceGroup, id.FunctionAppName, id.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving the AuthSettings for AzureRM Function App Slot %q (Function App %q / Resource Group %q): %+v", id.Name, id.FunctionAppName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("function_app_name", id.FunctionAppName)
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

	delete(appSettings, "AzureWebJobsDashboard")
	delete(appSettings, "AzureWebJobsStorage")
	delete(appSettings, "FUNCTIONS_EXTENSION_VERSION")
	delete(appSettings, "WEBSITE_CONTENTSHARE")
	delete(appSettings, "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING")

	if err = d.Set("app_settings", appSettings); err != nil {
		return err
	}
	if err = d.Set("connection_string", flattenFunctionAppSlotConnectionStrings(connectionStringsResp.Properties)); err != nil {
		return err
	}

	identity := flattenAppServiceIdentity(resp.Identity)
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("Error setting `identity`: %s", err)
	}

	configResp, err := client.GetConfigurationSlot(ctx, id.ResourceGroup, id.FunctionAppName, id.Name)
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

	siteCred := flattenFunctionAppSlotSiteCredential(siteCredResp.UserProperties)
	if err = d.Set("site_credential", siteCred); err != nil {
		return err
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmFunctionAppSlotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FunctionAppSlotID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Function App Slot %q (Function App %q / Resource Group %q)", id.Name, id.FunctionAppName, id.ResourceGroup)

	deleteMetrics := true
	deleteEmptyServerFarm := false
	resp, err := client.DeleteSlot(ctx, id.ResourceGroup, id.FunctionAppName, id.Name, &deleteMetrics, &deleteEmptyServerFarm)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return err
		}
	}

	return nil
}

func getBasicFunctionAppSlotAppSettings(d *schema.ResourceData, appServiceTier, endpointSuffix string) []web.NameValuePair {
	// TODO: This is a workaround since there are no public Functions API
	// You may track the API request here: https://github.com/Azure/azure-rest-api-specs/issues/3750
	dashboardPropName := "AzureWebJobsDashboard"
	storagePropName := "AzureWebJobsStorage"
	functionVersionPropName := "FUNCTIONS_EXTENSION_VERSION"
	contentSharePropName := "WEBSITE_CONTENTSHARE"
	contentFileConnStringPropName := "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"

	storageAccount := d.Get("storage_account_name").(string)
	connectionString := d.Get("storage_account_access_key").(string)
	storageConnection := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", storageAccount, connectionString, endpointSuffix)

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

	// On consumption and premium plans include WEBSITE_CONTENT components
	if strings.EqualFold(appServiceTier, "dynamic") || strings.EqualFold(appServiceTier, "elasticpremium") {
		return append(basicSettings, consumptionSettings...)
	}

	return basicSettings
}

func getFunctionAppSlotServiceTier(ctx context.Context, appServicePlanID string, meta interface{}) (string, error) {
	id, err := ParseAppServicePlanID(appServicePlanID)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Unable to parse App Service Plan ID %q: %+v", appServicePlanID, err)
	}

	log.Printf("[DEBUG] Retrieving App Service Plan %q (Resource Group %q)", id.Name, id.ResourceGroup)

	appServicePlansClient := meta.(*clients.Client).Web.AppServicePlansClient
	appServicePlan, err := appServicePlansClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Could not retrieve App Service Plan ID %q: %+v", appServicePlanID, err)
	}

	if sku := appServicePlan.Sku; sku != nil {
		if tier := sku.Tier; tier != nil {
			return *tier, nil
		}
	}
	return "", fmt.Errorf("No `sku` block was returned for App Service Plan ID %q", appServicePlanID)
}

func expandFunctionAppSlotAppSettings(d *schema.ResourceData, appServiceTier, endpointSuffix string) (map[string]*string, error) {
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

func expandFunctionAppSlotConnectionStrings(d *schema.ResourceData) map[string]*web.ConnStringValueTypePair {
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

func flattenFunctionAppSlotConnectionStrings(input map[string]*web.ConnStringValueTypePair) interface{} {
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

func flattenFunctionAppSlotSiteCredential(input *web.UserProperties) []interface{} {
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
