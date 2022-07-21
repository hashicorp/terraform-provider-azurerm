package logic

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogicAppStandard() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppStandardCreate,
		Read:   resourceLogicAppStandardRead,
		Update: resourceLogicAppStandardUpdate,
		Delete: resourceLogicAppStandardDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LogicAppStandardID(id)
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
				ValidateFunc: validate.LogicAppStandardName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

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

			"use_extension_bundle": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"bundle_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "[1.*, 2.0.0)",
			},

			"client_affinity_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"client_certificate_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Required",
					"Optional",
				}, false),
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"https_only": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			// TODO: API supports UserAssigned & SystemAssignedUserAssigned too?
			"identity": commonschema.SystemAssignedIdentityOptional(),

			"site_config": schemaLogicAppStandardSiteConfig(),

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
							}, false),
						},

						"value": {
							Type:      pluginsdk.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},

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

			"storage_account_share_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "~3",
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

func resourceLogicAppStandardCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	endpointSuffix := meta.(*clients.Client).Account.Environment.StorageEndpointSuffix
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Logic App Standard creation.")

	id := parse.NewLogicAppStandardID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_logic_app_standard", id.ID())
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
		return fmt.Errorf("the name %q used for the Logic App Standard needs to be globally unique and isn't available: %+v", id.SiteName, *available.Message)
	}

	appServicePlanID := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	clientAffinityEnabled := d.Get("client_affinity_enabled").(bool)
	clientCertMode := d.Get("client_certificate_mode").(string)
	clientCertEnabled := clientCertMode != ""
	httpsOnly := d.Get("https_only").(bool)
	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	basicAppSettings, err := getBasicLogicAppSettings(d, endpointSuffix)
	if err != nil {
		return err
	}

	siteConfig, err := expandLogicAppStandardSiteConfig(d)
	if err != nil {
		return fmt.Errorf("expanding `site_config`: %+v", err)
	}

	kind := "functionapp,workflowapp"
	if siteConfig.LinuxFxVersion != nil && len(*siteConfig.LinuxFxVersion) > 0 {
		kind = "functionapp,linux,container,workflowapp"
	}

	// Some appSettings declared by user are required at creation time so we will combine both settings
	appSettings := expandAppSettings(d)
	appSettings = append(appSettings, basicAppSettings...)

	siteConfig.AppSettings = &appSettings

	siteEnvelope := web.Site{
		Kind:     &kind,
		Location: &location,
		Tags:     tags.Expand(t),
		SiteProperties: &web.SiteProperties{
			ServerFarmID:          utils.String(appServicePlanID),
			Enabled:               utils.Bool(enabled),
			ClientAffinityEnabled: utils.Bool(clientAffinityEnabled),
			ClientCertEnabled:     utils.Bool(clientCertEnabled),
			HTTPSOnly:             utils.Bool(httpsOnly),
			SiteConfig:            &siteConfig,
		},
	}

	if clientCertMode != "" {
		siteEnvelope.SiteProperties.ClientCertMode = web.ClientCertMode(clientCertMode)
	}

	if _, ok := d.GetOk("identity"); ok {
		appServiceIdentity, err := expandLogicAppStandardIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		siteEnvelope.Identity = appServiceIdentity
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, siteEnvelope)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppStandardUpdate(d, meta)
}

func resourceLogicAppStandardUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	endpointSuffix := meta.(*clients.Client).Account.Environment.StorageEndpointSuffix
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogicAppStandardID(d.Id())
	if err != nil {
		return err
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	appServicePlanID := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	clientAffinityEnabled := d.Get("client_affinity_enabled").(bool)
	clientCertMode := d.Get("client_certificate_mode").(string)
	clientCertEnabled := clientCertMode != ""
	httpsOnly := d.Get("https_only").(bool)
	t := d.Get("tags").(map[string]interface{})

	basicAppSettings, err := getBasicLogicAppSettings(d, endpointSuffix)
	if err != nil {
		return err
	}

	siteConfig, err := expandLogicAppStandardSiteConfig(d)
	if err != nil {
		return fmt.Errorf("expanding `site_config`: %+v", err)
	}

	kind := "functionapp,workflowapp"
	if siteConfig.LinuxFxVersion != nil && len(*siteConfig.LinuxFxVersion) > 0 {
		kind = "functionapp,linux,container,workflowapp"
	}

	siteConfig.AppSettings = &basicAppSettings

	// WEBSITE_VNET_ROUTE_ALL is superseded by a setting in site_config that defaults to false from 2021-02-01
	appSettings, err := expandLogicAppStandardSettings(d, endpointSuffix)
	if err != nil {
		return fmt.Errorf("expanding `app_settings`: %+v", err)
	}
	if vnetRouteAll, ok := appSettings["WEBSITE_VNET_ROUTE_ALL"]; ok {
		if !d.HasChange("site_config.0.vnet_route_all_enabled") {
			vnetRouteAllEnabled, _ := strconv.ParseBool(*vnetRouteAll)
			siteConfig.VnetRouteAllEnabled = &vnetRouteAllEnabled
		}
	}

	siteEnvelope := web.Site{
		Kind:     &kind,
		Location: &location,
		Tags:     tags.Expand(t),
		SiteProperties: &web.SiteProperties{
			ServerFarmID:          utils.String(appServicePlanID),
			Enabled:               utils.Bool(enabled),
			ClientAffinityEnabled: utils.Bool(clientAffinityEnabled),
			ClientCertEnabled:     utils.Bool(clientCertEnabled),
			HTTPSOnly:             utils.Bool(httpsOnly),
			SiteConfig:            &siteConfig,
		},
	}

	if clientCertMode != "" {
		siteEnvelope.SiteProperties.ClientCertMode = web.ClientCertMode(clientCertMode)
	}

	if _, ok := d.GetOk("identity"); ok {
		appServiceIdentity, err := expandLogicAppStandardIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		siteEnvelope.Identity = appServiceIdentity
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, siteEnvelope)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", id, err)
	}

	settings := web.StringDictionary{
		Properties: appSettings,
	}

	if _, err = client.UpdateApplicationSettings(ctx, id.ResourceGroup, id.SiteName, settings); err != nil {
		return fmt.Errorf("updating Application Settings for %s: %+v", *id, err)
	}

	if d.HasChange("site_config") {
		siteConfig, err := expandLogicAppStandardSiteConfig(d)
		if err != nil {
			return fmt.Errorf("expanding `site_config`: %+v", err)
		}
		siteConfigResource := web.SiteConfigResource{
			SiteConfig: &siteConfig,
		}

		if _, err := client.CreateOrUpdateConfiguration(ctx, id.ResourceGroup, id.SiteName, siteConfigResource); err != nil {
			return fmt.Errorf("updating Configuration for %s: %+v", *id, err)
		}
	}

	if d.HasChange("connection_string") {
		connectionStrings := expandLogicAppStandardConnectionStrings(d)
		properties := web.ConnectionStringDictionary{
			Properties: connectionStrings,
		}

		if _, err := client.UpdateConnectionStrings(ctx, id.ResourceGroup, id.SiteName, properties); err != nil {
			return fmt.Errorf("updating Connection Strings for %s: %+v", *id, err)
		}
	}

	return resourceLogicAppStandardRead(d, meta)
}

func resourceLogicAppStandardRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogicAppStandardID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	appSettingsResp, err := client.ListApplicationSettings(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("listing application settings for %s: %+v", *id, err)
	}

	connectionStringsResp, err := client.ListConnectionStrings(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("listing connection strings for %s: %+v", *id, err)
	}

	siteCredFuture, err := client.ListPublishingCredentials(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("listing publishing credentials for %s: %+v", *id, err)
	}
	if err = siteCredFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting to list the publishing credentials for %s: %+v", *id, err)
	}
	siteCredResp, err := siteCredFuture.Result(*client)
	if err != nil {
		return fmt.Errorf("retrieving the publishing credentials for %s: %+v", *id, err)
	}

	d.Set("name", id.SiteName)
	d.Set("resource_group_name", id.ResourceGroup)
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
		d.Set("custom_domain_verification_id", props.CustomDomainVerificationID)

		clientCertMode := ""
		if props.ClientCertEnabled != nil && *props.ClientCertEnabled {
			clientCertMode = string(props.ClientCertMode)
		}
		d.Set("client_certificate_mode", clientCertMode)
	}

	appSettings := flattenLogicAppStandardAppSettings(appSettingsResp.Properties)

	if err = d.Set("connection_string", flattenLogicAppStandardConnectionStrings(connectionStringsResp.Properties)); err != nil {
		return err
	}

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

	if _, ok := appSettings["AzureFunctionsJobHost__extensionBundle__id"]; ok {
		d.Set("use_extension_bundle", true)
		if val, ok := appSettings["AzureFunctionsJobHost__extensionBundle__version"]; ok {
			d.Set("bundle_version", val)
		}
	} else {
		d.Set("use_extension_bundle", false)
		d.Set("bundle_version", "[1.*, 2.0.0)")
	}

	d.Set("storage_account_share_name", appSettings["WEBSITE_CONTENTSHARE"])

	// Remove all the settings that are created by this resource so we don't to have to specify in app_settings
	// block whenever we use azurerm_logic_app_standard.
	delete(appSettings, "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING")
	delete(appSettings, "APP_KIND")
	delete(appSettings, "AzureFunctionsJobHost__extensionBundle__id")
	delete(appSettings, "AzureFunctionsJobHost__extensionBundle__version")
	delete(appSettings, "AzureWebJobsDashboard")
	delete(appSettings, "AzureWebJobsStorage")
	delete(appSettings, "FUNCTIONS_EXTENSION_VERSION")
	delete(appSettings, "WEBSITE_CONTENTSHARE")

	if err = d.Set("app_settings", appSettings); err != nil {
		return err
	}

	identity := flattenLogicAppStandardIdentity(resp.Identity)
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %s", err)
	}

	configResp, err := client.GetConfiguration(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("retrieving the configuration for %s: %+v", *id, err)
	}

	siteConfig := flattenLogicAppStandardSiteConfig(configResp.SiteConfig)
	if err = d.Set("site_config", siteConfig); err != nil {
		return err
	}

	siteCred := flattenLogicAppStandardSiteCredential(siteCredResp.UserProperties)
	if err = d.Set("site_credential", siteCred); err != nil {
		return err
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceLogicAppStandardDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogicAppStandardID(d.Id())
	if err != nil {
		return err
	}

	deleteMetrics := true
	deleteEmptyServerFarm := false
	if _, err := client.Delete(ctx, id.ResourceGroup, id.SiteName, &deleteMetrics, &deleteEmptyServerFarm); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func getBasicLogicAppSettings(d *pluginsdk.ResourceData, endpointSuffix string) ([]web.NameValuePair, error) {
	storagePropName := "AzureWebJobsStorage"
	functionVersionPropName := "FUNCTIONS_EXTENSION_VERSION"
	contentSharePropName := "WEBSITE_CONTENTSHARE"
	contentFileConnStringPropName := "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"
	appKindPropName := "APP_KIND"
	appKindPropValue := "workflowApp"

	storageAccount := d.Get("storage_account_name").(string)
	accountKey := d.Get("storage_account_access_key").(string)
	storageConnection := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", storageAccount, accountKey, endpointSuffix)
	functionVersion := d.Get("version").(string)

	contentShare := strings.ToLower(d.Get("name").(string)) + "-content"
	if _, ok := d.GetOk("storage_account_share_name"); ok {
		contentShare = d.Get("storage_account_share_name").(string)
	}

	basicSettings := []web.NameValuePair{
		{Name: &storagePropName, Value: &storageConnection},
		{Name: &functionVersionPropName, Value: &functionVersion},
		{Name: &appKindPropName, Value: &appKindPropValue},
		{Name: &contentSharePropName, Value: &contentShare},
		{Name: &contentFileConnStringPropName, Value: &storageConnection},
	}

	useExtensionBundle := d.Get("use_extension_bundle").(bool)
	if useExtensionBundle {
		extensionBundlePropName := "AzureFunctionsJobHost__extensionBundle__id"
		extensionBundleName := "Microsoft.Azure.Functions.ExtensionBundle.Workflows"
		extensionBundleVersionPropName := "AzureFunctionsJobHost__extensionBundle__version"
		extensionBundleVersion := d.Get("bundle_version").(string)

		if extensionBundleVersion == "" {
			return nil, fmt.Errorf("when `use_extension_bundle` is true, `bundle_version` must be specified")
		}

		bundleSettings := []web.NameValuePair{
			{Name: &extensionBundlePropName, Value: &extensionBundleName},
			{Name: &extensionBundleVersionPropName, Value: &extensionBundleVersion},
		}

		return append(basicSettings, bundleSettings...), nil
	}

	return basicSettings, nil
}

func schemaLogicAppStandardSiteConfig() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"cors": schemaLogicAppCorsSettings(),

				"ftps_state": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.FtpsStateAllAllowed),
						string(web.FtpsStateDisabled),
						string(web.FtpsStateFtpsOnly),
					}, false),
				},

				"http2_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"ip_restriction": schemaLogicAppStandardIpRestriction(),

				"linux_fx_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
				},

				"min_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SupportedTLSVersionsOneFullStopZero),
						string(web.SupportedTLSVersionsOneFullStopOne),
						string(web.SupportedTLSVersionsOneFullStopTwo),
					}, false),
				},

				"pre_warmed_instance_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(0, 20),
				},

				"use_32_bit_worker_process": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"websockets_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"health_check_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"elastic_instance_minimum": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(0, 20),
				},

				"app_scale_limit": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntAtLeast(0),
				},

				"runtime_scale_monitoring_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"dotnet_framework_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "v4.0",
					ValidateFunc: validation.StringInSlice([]string{
						"v4.0",
						"v5.0",
						"v6.0",
					}, false),
				},

				"vnet_route_all_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func schemaLogicAppCorsSettings() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"allowed_origins": {
					Type:     pluginsdk.TypeSet,
					Required: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},
				"support_credentials": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func schemaLogicAppStandardIpRestriction() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:       pluginsdk.TypeList,
		Optional:   true,
		Computed:   true,
		ConfigMode: pluginsdk.SchemaConfigModeAttr,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"ip_address": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"service_tag": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"virtual_network_subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"priority": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      65000,
					ValidateFunc: validation.IntBetween(1, 2147483647),
				},

				"action": {
					Type:     pluginsdk.TypeString,
					Default:  "Allow",
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"Allow",
						"Deny",
					}, false),
				},

				//lintignore:XS003
				"headers": {
					Type:       pluginsdk.TypeList,
					Optional:   true,
					Computed:   true,
					MaxItems:   1,
					ConfigMode: pluginsdk.SchemaConfigModeAttr,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							//lintignore:S018
							"x_forwarded_host": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								MaxItems: 8,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},

							//lintignore:S018
							"x_forwarded_for": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								MaxItems: 8,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.IsCIDR,
								},
							},

							//lintignore:S018
							"x_azure_fdid": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								MaxItems: 8,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.IsUUID,
								},
							},

							//lintignore:S018
							"x_fd_health_probe": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								MaxItems: 1,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
									ValidateFunc: validation.StringInSlice([]string{
										"1",
									}, false),
								},
							},
						},
					},
				},
			},
		},
	}
}

func flattenLogicAppStandardAppSettings(input map[string]*string) map[string]string {
	output := make(map[string]string)
	for k, v := range input {
		output[k] = *v
	}

	return output
}

func flattenLogicAppStandardConnectionStrings(input map[string]*web.ConnStringValueTypePair) interface{} {
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

func flattenLogicAppStandardSiteConfig(input *web.SiteConfig) []interface{} {
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

	result["ip_restriction"] = flattenLogicAppStandardIpRestriction(input.IPSecurityRestrictions)

	result["min_tls_version"] = string(input.MinTLSVersion)
	result["ftps_state"] = string(input.FtpsState)

	result["cors"] = flattenLogicAppStandardCorsSettings(input.Cors)

	if input.AutoSwapSlotName != nil {
		result["auto_swap_slot_name"] = *input.AutoSwapSlotName
	}

	if input.HealthCheckPath != nil {
		result["health_check_path"] = *input.HealthCheckPath
	}

	if input.MinimumElasticInstanceCount != nil {
		result["elastic_instance_minimum"] = *input.MinimumElasticInstanceCount
	}

	if input.FunctionAppScaleLimit != nil {
		result["app_scale_limit"] = *input.FunctionAppScaleLimit
	}

	if input.FunctionsRuntimeScaleMonitoringEnabled != nil {
		result["runtime_scale_monitoring_enabled"] = *input.FunctionsRuntimeScaleMonitoringEnabled
	}

	if input.NetFrameworkVersion != nil {
		result["dotnet_framework_version"] = *input.NetFrameworkVersion
	}

	vnetRouteAllEnabled := false
	if input.VnetRouteAllEnabled != nil {
		vnetRouteAllEnabled = *input.VnetRouteAllEnabled
	}
	result["vnet_route_all_enabled"] = vnetRouteAllEnabled

	results = append(results, result)
	return results
}

func flattenLogicAppStandardSiteCredential(input *web.UserProperties) []interface{} {
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

func flattenLogicAppStandardIpRestriction(input *[]web.IPSecurityRestriction) []interface{} {
	restrictions := make([]interface{}, 0)

	if input == nil {
		return restrictions
	}

	for _, v := range *input {
		restriction := make(map[string]interface{})
		if ip := v.IPAddress; ip != nil {
			if *ip == "Any" {
				continue
			} else {
				switch v.Tag {
				case web.IPFilterTagServiceTag:
					restriction["service_tag"] = *ip
				default:
					restriction["ip_address"] = *ip
				}
			}
		}

		subnetId := ""
		if subnetIdRaw := v.VnetSubnetResourceID; subnetIdRaw != nil {
			subnetId = *subnetIdRaw
		}
		restriction["virtual_network_subnet_id"] = subnetId

		name := ""
		if nameRaw := v.Name; nameRaw != nil {
			name = *nameRaw
		}
		restriction["name"] = name

		priority := 0
		if priorityRaw := v.Priority; priorityRaw != nil {
			priority = int(*priorityRaw)
		}
		restriction["priority"] = priority

		action := ""
		if actionRaw := v.Action; actionRaw != nil {
			action = *actionRaw
		}
		restriction["action"] = action

		if headers := v.Headers; headers != nil {
			restriction["headers"] = flattenHeaders(headers)
		}

		restrictions = append(restrictions, restriction)
	}

	return restrictions
}

func flattenLogicAppStandardCorsSettings(input *web.CorsSettings) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	allowedOrigins := make([]interface{}, 0)
	if s := input.AllowedOrigins; s != nil {
		for _, v := range *s {
			allowedOrigins = append(allowedOrigins, v)
		}
	}
	result["allowed_origins"] = pluginsdk.NewSet(pluginsdk.HashString, allowedOrigins)

	if input.SupportCredentials != nil {
		result["support_credentials"] = *input.SupportCredentials
	}

	return append(results, result)
}

func flattenHeaders(input map[string][]string) []interface{} {
	output := make([]interface{}, 0)
	headers := make(map[string]interface{})
	if input == nil {
		return output
	}

	if forwardedHost, ok := input["x-forwarded-host"]; ok && len(forwardedHost) > 0 {
		headers["x_forwarded_host"] = forwardedHost
	}
	if forwardedFor, ok := input["x-forwarded-for"]; ok && len(forwardedFor) > 0 {
		headers["x_forwarded_for"] = forwardedFor
	}
	if fdids, ok := input["x-azure-fdid"]; ok && len(fdids) > 0 {
		headers["x_azure_fdid"] = fdids
	}
	if healthProbe, ok := input["x-fd-healthprobe"]; ok && len(healthProbe) > 0 {
		headers["x_fd_health_probe"] = healthProbe
	}

	return append(output, headers)
}

func expandLogicAppStandardSiteConfig(d *pluginsdk.ResourceData) (web.SiteConfig, error) {
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
		expand := expandLogicAppStandardCorsSettings(v)
		siteConfig.Cors = &expand
	}

	if v, ok := config["http2_enabled"]; ok {
		siteConfig.HTTP20Enabled = utils.Bool(v.(bool))
	}

	if v, ok := config["ip_restriction"]; ok {
		restrictions, err := expandLogicAppStandardIpRestriction(v)
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

	if v, ok := config["health_check_path"]; ok {
		siteConfig.HealthCheckPath = utils.String(v.(string))
	}

	if v, ok := config["elastic_instance_minimum"]; ok {
		siteConfig.MinimumElasticInstanceCount = utils.Int32(int32(v.(int)))
	}

	if v, ok := config["app_scale_limit"]; ok {
		siteConfig.FunctionAppScaleLimit = utils.Int32(int32(v.(int)))
	}

	if v, ok := config["runtime_scale_monitoring_enabled"]; ok {
		siteConfig.FunctionsRuntimeScaleMonitoringEnabled = utils.Bool(v.(bool))
	}

	if v, ok := config["dotnet_framework_version"]; ok {
		siteConfig.NetFrameworkVersion = utils.String(v.(string))
	}

	if v, ok := config["vnet_route_all_enabled"]; ok {
		siteConfig.VnetRouteAllEnabled = utils.Bool(v.(bool))
	}

	return siteConfig, nil
}

func expandLogicAppStandardIdentity(input []interface{}) (*web.ManagedServiceIdentity, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &web.ManagedServiceIdentity{
		Type: web.ManagedServiceIdentityType(expanded.Type),
	}, nil
}

func flattenLogicAppStandardIdentity(input *web.ManagedServiceIdentity) []interface{} {
	var transform *identity.SystemAssigned

	if input != nil {
		transform = &identity.SystemAssigned{
			Type: identity.Type(string(input.Type)),
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
	}

	return identity.FlattenSystemAssigned(transform)
}

func expandLogicAppStandardSettings(d *pluginsdk.ResourceData, endpointSuffix string) (map[string]*string, error) {
	output := make(map[string]*string)
	appSettings := expandAppSettings(d)
	basicAppSettings, err := getBasicLogicAppSettings(d, endpointSuffix)
	if err != nil {
		return nil, err
	}
	for _, p := range append(basicAppSettings, appSettings...) {
		output[*p.Name] = p.Value
	}

	return output, nil
}

func expandAppSettings(d *pluginsdk.ResourceData) []web.NameValuePair {
	input := d.Get("app_settings").(map[string]interface{})
	output := make([]web.NameValuePair, 0)

	for k, v := range input {
		nameValue := web.NameValuePair{
			Name:  utils.String(k),
			Value: utils.String(v.(string)),
		}
		output = append(output, nameValue)
	}

	return output
}

func expandLogicAppStandardConnectionStrings(d *pluginsdk.ResourceData) map[string]*web.ConnStringValueTypePair {
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

func expandLogicAppStandardCorsSettings(input interface{}) web.CorsSettings {
	settings := input.([]interface{})
	corsSettings := web.CorsSettings{}

	if len(settings) == 0 {
		return corsSettings
	}

	setting := settings[0].(map[string]interface{})

	if v, ok := setting["allowed_origins"]; ok {
		input := v.(*pluginsdk.Set).List()

		allowedOrigins := make([]string, 0)
		for _, param := range input {
			allowedOrigins = append(allowedOrigins, param.(string))
		}

		corsSettings.AllowedOrigins = &allowedOrigins
	}

	if v, ok := setting["support_credentials"]; ok {
		corsSettings.SupportCredentials = utils.Bool(v.(bool))
	}

	return corsSettings
}

func expandLogicAppStandardIpRestriction(input interface{}) ([]web.IPSecurityRestriction, error) {
	restrictions := make([]web.IPSecurityRestriction, 0)

	for _, r := range input.([]interface{}) {
		if r == nil {
			continue
		}

		restriction := r.(map[string]interface{})

		ipAddress := restriction["ip_address"].(string)
		vNetSubnetID := ""

		if subnetID, ok := restriction["virtual_network_subnet_id"]; ok && subnetID != "" {
			vNetSubnetID = subnetID.(string)
		}

		serviceTag := restriction["service_tag"].(string)

		name := restriction["name"].(string)
		priority := restriction["priority"].(int)
		action := restriction["action"].(string)

		if vNetSubnetID != "" && ipAddress != "" && serviceTag != "" {
			return nil, fmt.Errorf("only one of `ip_address`, `service_tag` or `virtual_network_subnet_id` can be set for an IP restriction")
		}

		if vNetSubnetID == "" && ipAddress == "" && serviceTag == "" {
			return nil, fmt.Errorf("one of `ip_address`, `service_tag` or `virtual_network_subnet_id` must be set for an IP restriction")
		}

		ipSecurityRestriction := web.IPSecurityRestriction{}
		if ipAddress == "Any" {
			continue
		}

		if ipAddress != "" {
			ipSecurityRestriction.IPAddress = &ipAddress
		}

		if serviceTag != "" {
			ipSecurityRestriction.IPAddress = &serviceTag
			ipSecurityRestriction.Tag = web.IPFilterTagServiceTag
		}

		if vNetSubnetID != "" {
			ipSecurityRestriction.VnetSubnetResourceID = &vNetSubnetID
		}

		if name != "" {
			ipSecurityRestriction.Name = &name
		}

		if priority != 0 {
			ipSecurityRestriction.Priority = utils.Int32(int32(priority))
		}

		if action != "" {
			ipSecurityRestriction.Action = &action
		}
		if headers, ok := restriction["headers"]; ok {
			ipSecurityRestriction.Headers = expandHeaders(headers.([]interface{}))
		}

		restrictions = append(restrictions, ipSecurityRestriction)
	}

	return restrictions, nil
}

func expandHeaders(input interface{}) map[string][]string {
	output := make(map[string][]string)

	for _, r := range input.([]interface{}) {
		if r == nil {
			continue
		}

		val := r.(map[string]interface{})
		if raw := val["x_forwarded_host"].(*pluginsdk.Set).List(); len(raw) > 0 {
			output["x-forwarded-host"] = *utils.ExpandStringSlice(raw)
		}
		if raw := val["x_forwarded_for"].(*pluginsdk.Set).List(); len(raw) > 0 {
			output["x-forwarded-for"] = *utils.ExpandStringSlice(raw)
		}
		if raw := val["x_azure_fdid"].(*pluginsdk.Set).List(); len(raw) > 0 {
			output["x-azure-fdid"] = *utils.ExpandStringSlice(raw)
		}
		if raw := val["x_fd_health_probe"].(*pluginsdk.Set).List(); len(raw) > 0 {
			output["x-fd-healthprobe"] = *utils.ExpandStringSlice(raw)
		}
	}

	return output
}
