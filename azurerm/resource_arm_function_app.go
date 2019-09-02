package azurerm

import (
	"context"
	"fmt"
	"log"
	"strings"

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

// Azure Function App shares the same infrastructure with Azure App Service.
// So this resource will reuse most of the App Service code, but remove the configurations which are not applicable for Function App.
func resourceArmFunctionApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFunctionAppCreate,
		Read:   resourceArmFunctionAppRead,
		Update: resourceArmFunctionAppUpdate,
		Delete: resourceArmFunctionAppDelete,
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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"app_service_plan_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

			"storage_connection_string": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"app_settings": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"enable_builtin_logging": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"connection_string": {
				Type:     schema.TypeList,
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

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(web.ManagedServiceIdentityTypeSystemAssigned),
							}, true),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.Schema(),

			"default_hostname": {
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
						"virtual_network_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: false,
						},
						"cors": azure.SchemaWebCorsSettings(),
					},
				},
			},

			"auth_settings": azure.SchemaAppServiceAuthSettings(),

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
		},
	}
}

func resourceArmFunctionAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient
	ctx := meta.(*ArmClient).StopContext

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
	appServicePlanID := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	clientAffinityEnabled := d.Get("client_affinity_enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	t := d.Get("tags").(map[string]interface{})
	appServiceTier, err := getFunctionAppServiceTier(ctx, appServicePlanID, meta)
	if err != nil {
		return err
	}

	basicAppSettings := getBasicFunctionAppAppSettings(d, appServiceTier)

	siteConfig := expandFunctionAppSiteConfig(d)
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
			SiteConfig:            &siteConfig,
		},
	}

	if v, ok := d.GetOk("identity.0.type"); ok {
		siteEnvelope.Identity = &web.ManagedServiceIdentity{
			Type: web.ManagedServiceIdentityType(v.(string)),
		}
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
		return fmt.Errorf("Error updating auth settings for Function App %q (resource group %q): %+s", name, resourceGroup, err)
	}

	return resourceArmFunctionAppUpdate(d, meta)
}

func resourceArmFunctionAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["sites"]

	location := azure.NormalizeLocation(d.Get("location").(string))
	kind := "functionapp"
	appServicePlanID := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	clientAffinityEnabled := d.Get("client_affinity_enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	t := d.Get("tags").(map[string]interface{})

	appServiceTier, err := getFunctionAppServiceTier(ctx, appServicePlanID, meta)

	if err != nil {
		return err
	}
	basicAppSettings := getBasicFunctionAppAppSettings(d, appServiceTier)
	siteConfig := expandFunctionAppSiteConfig(d)

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
			SiteConfig:            &siteConfig,
		},
	}

	if v, ok := d.GetOk("identity.0.type"); ok {
		siteEnvelope.Identity = &web.ManagedServiceIdentity{
			Type: web.ManagedServiceIdentityType(v.(string)),
		}
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, siteEnvelope)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	appSettings := expandFunctionAppAppSettings(d, appServiceTier)
	settings := web.StringDictionary{
		Properties: appSettings,
	}

	_, err = client.UpdateApplicationSettings(ctx, resGroup, name, settings)
	if err != nil {
		return fmt.Errorf("Error updating Application Settings for Function App %q: %+v", name, err)
	}

	if d.HasChange("site_config") {
		siteConfig := expandFunctionAppSiteConfig(d)
		siteConfigResource := web.SiteConfigResource{
			SiteConfig: &siteConfig,
		}
		if _, err := client.CreateOrUpdateConfiguration(ctx, resGroup, name, siteConfigResource); err != nil {
			return fmt.Errorf("Error updating Configuration for Function App %q: %+v", name, err)
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
			return fmt.Errorf("Error updating Authentication Settings for Function App %q: %+v", name, err)
		}
	}

	if d.HasChange("connection_string") {
		// update the ConnectionStrings
		connectionStrings := expandFunctionAppConnectionStrings(d)
		properties := web.ConnectionStringDictionary{
			Properties: connectionStrings,
		}

		if _, err := client.UpdateConnectionStrings(ctx, resGroup, name, properties); err != nil {
			return fmt.Errorf("Error updating Connection Strings for App Service %q: %+v", name, err)
		}
	}

	return resourceArmFunctionAppRead(d, meta)
}

func resourceArmFunctionAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["sites"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Function App %q (resource group %q) was not found - removing from state", name, resGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Function App %q: %+v", name, err)
	}

	appSettingsResp, err := client.ListApplicationSettings(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(appSettingsResp.Response) {
			log.Printf("[DEBUG] Application Settings of Function App %q (resource group %q) were not found", name, resGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Function App AppSettings %q: %+v", name, err)
	}

	connectionStringsResp, err := client.ListConnectionStrings(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM Function App ConnectionStrings %q: %+v", name, err)
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
	authResp, err := client.GetAuthSettings(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving the AuthSettings for Function App %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("kind", resp.Kind)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.SiteProperties; props != nil {
		d.Set("app_service_plan_id", props.ServerFarmID)
		d.Set("enabled", props.Enabled)
		d.Set("default_hostname", props.DefaultHostName)
		d.Set("https_only", props.HTTPSOnly)
		d.Set("outbound_ip_addresses", props.OutboundIPAddresses)
		d.Set("possible_outbound_ip_addresses", props.PossibleOutboundIPAddresses)
		d.Set("client_affinity_enabled", props.ClientAffinityEnabled)
	}

	if err = d.Set("identity", flattenFunctionAppIdentity(resp.Identity)); err != nil {
		return err
	}

	appSettings := flattenAppServiceAppSettings(appSettingsResp.Properties)

	d.Set("storage_connection_string", appSettings["AzureWebJobsStorage"])
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
	if err = d.Set("connection_string", flattenFunctionAppConnectionStrings(connectionStringsResp.Properties)); err != nil {
		return err
	}

	configResp, err := client.GetConfiguration(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM Function App Configuration %q: %+v", name, err)
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
	client := meta.(*ArmClient).web.AppServicesClient

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["sites"]

	log.Printf("[DEBUG] Deleting Function App %q (resource group %q)", name, resGroup)

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

func getBasicFunctionAppAppSettings(d *schema.ResourceData, appServiceTier string) []web.NameValuePair {
	// TODO: This is a workaround since there are no public Functions API
	// You may track the API request here: https://github.com/Azure/azure-rest-api-specs/issues/3750
	dashboardPropName := "AzureWebJobsDashboard"
	storagePropName := "AzureWebJobsStorage"
	functionVersionPropName := "FUNCTIONS_EXTENSION_VERSION"
	contentSharePropName := "WEBSITE_CONTENTSHARE"
	contentFileConnStringPropName := "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"

	storageConnection := d.Get("storage_connection_string").(string)
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

	// If the application plan is NOT dynamic (consumption plan), we do NOT want to include WEBSITE_CONTENT components
	if !strings.EqualFold(appServiceTier, "dynamic") {
		return basicSettings
	}
	return append(basicSettings, consumptionSettings...)
}

func getFunctionAppServiceTier(ctx context.Context, appServicePlanId string, meta interface{}) (string, error) {
	id, err := azure.ParseAzureResourceID(appServicePlanId)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Unable to parse App Service Plan ID %q: %+v", appServicePlanId, err)
	}

	log.Printf("[DEBUG] Retrieving App Server Plan %s", id.Path["serverfarms"])

	appServicePlansClient := meta.(*ArmClient).web.AppServicePlansClient
	appServicePlan, err := appServicePlansClient.Get(ctx, id.ResourceGroup, id.Path["serverfarms"])
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

func expandFunctionAppAppSettings(d *schema.ResourceData, appServiceTier string) map[string]*string {
	output := expandAppServiceAppSettings(d)

	basicAppSettings := getBasicFunctionAppAppSettings(d, appServiceTier)
	for _, p := range basicAppSettings {
		output[*p.Name] = p.Value
	}

	return output
}

func expandFunctionAppSiteConfig(d *schema.ResourceData) web.SiteConfig {
	configs := d.Get("site_config").([]interface{})
	siteConfig := web.SiteConfig{}

	if len(configs) == 0 {
		return siteConfig
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

	if v, ok := config["virtual_network_name"]; ok {
		siteConfig.VnetName = utils.String(v.(string))
	}

	return siteConfig
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

	if input.VnetName != nil {
		result["virtual_network_name"] = *input.VnetName
	}

	result["cors"] = azure.FlattenWebCorsSettings(input.Cors)

	results = append(results, result)
	return results
}

func expandFunctionAppConnectionStrings(d *schema.ResourceData) map[string]*web.ConnStringValueTypePair {
	input := d.Get("connection_string").([]interface{})
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

func flattenFunctionAppIdentity(identity *web.ManagedServiceIdentity) interface{} {
	if identity == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	result["type"] = string(identity.Type)

	if identity.PrincipalID != nil {
		result["principal_id"] = *identity.PrincipalID
	}
	if identity.TenantID != nil {
		result["tenant_id"] = *identity.TenantID
	}

	return []interface{}{result}
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
