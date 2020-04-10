package web

import (
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
			_, err := ParseFunctionAppSlotID(id)
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

			"identity": azure.SchemaAppServiceIdentity(),

			"function_app_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"app_service_plan_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"enable_builtin_logging": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"client_affinity_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
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

			"daily_memory_time_quota": {
				Type:     schema.TypeInt,
				Optional: true,
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

			"outbound_ip_addresses": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"possible_outbound_ip_addresses": {
				Type:     schema.TypeString,
				Computed: true,
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

			"tags": tags.Schema(),
		},
	}
}

func resourceArmFunctionAppSlotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	slot := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	functionAppName := d.Get("function_app_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
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

	appServicePlanId := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	clientAffinityEnabled := d.Get("client_affinity_enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	dailyMemoryTimeQuota := d.Get("daily_memory_time_quota").(int)
	t := d.Get("tags").(map[string]interface{})
	appServiceTier, err := getFunctionAppServiceTier(ctx, appServicePlanId, meta)
	if err != nil {
		return err
	}

	basicAppSettings := getBasicFunctionAppAppSettings(d, appServiceTier)
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
			ServerFarmID:          utils.String(appServicePlanId),
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

	createFuture, err := client.CreateOrUpdateSlot(ctx, resourceGroup, functionAppName, siteEnvelope, slot)
	if err != nil {
		return fmt.Errorf("Error creating Slot %q (Function App %q / Resource Group %q): %s", slot, functionAppName, resourceGroup, err)
	}

	err = createFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation of Slot %q (Function App %q / Resource Group %q): %s", slot, functionAppName, resourceGroup, err)
	}

	read, err := client.GetSlot(ctx, resourceGroup, functionAppName, slot)
	if err != nil {
		return fmt.Errorf("Error retrieving Slot %q (Function App %q / Resource Group %q): %s", slot, functionAppName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for Slot %q (Function App %q / Resource Group %q) ID", slot, functionAppName, resourceGroup)
	}

	d.SetId(*read.ID)

	authSettingsRaw := d.Get("auth_settings").([]interface{})
	authSettings := azure.ExpandAppServiceAuthSettings(authSettingsRaw)

	auth := web.SiteAuthSettings{
		ID:                         read.ID,
		SiteAuthSettingsProperties: &authSettings}

	if _, err := client.UpdateAuthSettingsSlot(ctx, resourceGroup, functionAppName, auth, slot); err != nil {
		return fmt.Errorf("Error updating auth settings for Slot %q (Function App Slot %q / Resource Group %q): %+s", slot, functionAppName, resourceGroup, err)
	}

	return resourceArmFunctionAppSlotUpdate(d, meta)
}

func resourceArmFunctionAppSlotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseFunctionAppSlotID(d.Id())
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
	appServicePlanId := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	clientAffinityEnabled := d.Get("client_affinity_enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	dailyMemoryTimeQuota := d.Get("daily_memory_time_quota").(int)
	t := d.Get("tags").(map[string]interface{})

	appServiceTier, err := getFunctionAppServiceTier(ctx, appServicePlanId, meta)

	if err != nil {
		return err
	}
	basicAppSettings := getBasicFunctionAppAppSettings(d, appServiceTier)
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
			ServerFarmID:          utils.String(appServicePlanId),
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

	future, err := client.CreateOrUpdateSlot(ctx, id.ResourceGroup, id.FunctionAppName, siteEnvelope, id.Name)
	if err != nil {
		return fmt.Errorf("Error updating Slot %q (Function App %q / Resource Group %q): %s", id.Name, id.FunctionAppName, id.ResourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for update of Slot %q (Function App %q / Resource Group %q): %s", id.Name, id.FunctionAppName, id.ResourceGroup, err)
	}

	appSettings := expandFunctionAppAppSettings(d, appServiceTier)
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
		authSettingsProperties := azure.ExpandAppServiceAuthSettings(authSettingsRaw)
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
		connectionStrings := expandFunctionAppConnectionStrings(d)
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

	id, err := ParseFunctionAppSlotID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetSlot(ctx, id.ResourceGroup, id.FunctionAppName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Slot %q (Function App %q / Resource Group %q) were not found - removing from state!", id.Name, id.FunctionAppName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Slot %q (Function App %q / Resource Group %q): %s", id.Name, id.FunctionAppName, id.ResourceGroup, err)
	}

	appSettingsResp, err := client.ListApplicationSettingsSlot(ctx, id.ResourceGroup, id.FunctionAppName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(appSettingsResp.Response) {
			log.Printf("[DEBUG] Application Settings of Slot %q (Function App %q / Resource Group %q) were not found", id.Name, id.FunctionAppName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Slot %q (Function App %q / Resource Group %q) AppSettings: %+v", id.Name, id.FunctionAppName, id.ResourceGroup, err)
	}

	connectionStringsResp, err := client.ListConnectionStringsSlot(ctx, id.ResourceGroup, id.FunctionAppName, id.Name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Slot %q (Function App %q / Resource Group %q) ConnectionStrings: %+v", id.Name, id.FunctionAppName, id.ResourceGroup, err)
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
		return fmt.Errorf("Error making Read request on Slot %q (Function App %q / Resource Group %q) Site Credentials: %+v", id.Name, id.FunctionAppName, id.ResourceGroup, err)
	}
	authResp, err := client.GetAuthSettingsSlot(ctx, id.ResourceGroup, id.FunctionAppName, id.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving the AuthSettings for Slot %q (Function App %q / Resource Group %q): %+v", id.Name, id.FunctionAppName, id.ResourceGroup, err)
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

	identity := azure.FlattenAppServiceIdentity(resp.Identity)
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

func resourceArmFunctionAppSlotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseFunctionAppSlotID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Slot %q (Function App %q / Resource Group %q)", id.Name, id.FunctionAppName, id.ResourceGroup)

	deleteMetrics := true
	deleteEmptyServerFarm := false
	resp, err := client.DeleteSlot(ctx, id.ResourceGroup, id.FunctionAppName, id.Name, &deleteMetrics, &deleteEmptyServerFarm)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Slot %q (Function App %q / Resource Group %q): %s", id.Name, id.FunctionAppName, id.ResourceGroup, err)
		}
	}

	return nil
}
