// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogicAppStandard() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceLogicAppStandardCreate,
		Read:   resourceLogicAppStandardRead,
		Update: resourceLogicAppStandardUpdate,
		Delete: resourceLogicAppStandardDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseLogicAppId(id)
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

			"ftp_publish_basic_authentication_enabled": {
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

			"scm_publish_basic_authentication_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

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
								string(webapps.ConnectionStringTypeApiHub),
								string(webapps.ConnectionStringTypeCustom),
								string(webapps.ConnectionStringTypeDocDb),
								string(webapps.ConnectionStringTypeEventHub),
								string(webapps.ConnectionStringTypeMySql),
								string(webapps.ConnectionStringTypeNotificationHub),
								string(webapps.ConnectionStringTypePostgreSQL),
								string(webapps.ConnectionStringTypeRedisCache),
								string(webapps.ConnectionStringTypeServiceBus),
								string(webapps.ConnectionStringTypeSQLAzure),
								string(webapps.ConnectionStringTypeSQLServer),
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

			"storage_account": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MaxItems: 5,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"type": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForAzureStorageType(), false),
						},
						"account_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"share_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"access_key": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"mount_path": {
							Type:     pluginsdk.TypeString,
							Optional: true,
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

			"public_network_access": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  helpers.PublicNetworkAccessEnabled,
				ValidateFunc: validation.StringInSlice([]string{
					helpers.PublicNetworkAccessEnabled,
					helpers.PublicNetworkAccessDisabled,
				}, false),
			},

			"storage_account_share_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "~4",
			},

			"vnet_content_share_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"virtual_network_subnet_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"tags": commonschema.Tags(),

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

	if !features.FivePointOh() {
		// Due to the way the `site_config.public_network_access_enabled` property and the `public_network_access` property
		// influence each other, the default needs to be handled in the Create for now until `site_config.public_network_access_enabled`
		// is removed in v5.0
		resource.Schema["public_network_access"].Default = nil
		resource.Schema["public_network_access"].Computed = true
	}
	return resource
}

func resourceLogicAppStandardCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppService.WebAppsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourcesClient := meta.(*clients.Client).AppService.ResourceProvidersClient

	env := meta.(*clients.Client).Account.Environment
	storageAccountDomainSuffix, ok := env.Storage.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine the domain suffix for storage accounts in environment %q: %+v", env.Name, env.Storage)
	}

	log.Printf("[INFO] preparing arguments for AzureRM Logic App Standard creation.")

	id := commonids.NewAppServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_logic_app_standard", id.ID())
	}

	availabilityRequest := resourceproviders.ResourceNameAvailabilityRequest{
		Name: id.SiteName,
		Type: resourceproviders.CheckNameResourceTypesMicrosoftPointWebSites,
	}

	available, err := resourcesClient.CheckNameAvailability(ctx, commonids.NewSubscriptionID(subscriptionId), availabilityRequest)
	if err != nil {
		return fmt.Errorf("checking if name %q was available: %+v", id.SiteName, err)
	}

	if available.Model == nil || available.Model.NameAvailable == nil {
		return fmt.Errorf("checking if name %q was available: `model` was nil", id.SiteName)
	}

	if !*available.Model.NameAvailable {
		return fmt.Errorf("the name %q used for the Logic App Standard needs to be globally unique and isn't available: %+v", id.SiteName, pointer.From(available.Model.Message))
	}

	clientCertMode := d.Get("client_certificate_mode").(string)
	clientCertEnabled := clientCertMode != ""

	basicAppSettings, err := getBasicLogicAppSettings(d, *storageAccountDomainSuffix)
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

	siteEnvelope := webapps.Site{
		Kind:     &kind,
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &webapps.SiteProperties{
			ServerFarmId:            pointer.To(d.Get("app_service_plan_id").(string)),
			Enabled:                 pointer.To(d.Get("enabled").(bool)),
			ClientAffinityEnabled:   pointer.To(d.Get("client_affinity_enabled").(bool)),
			ClientCertEnabled:       pointer.To(clientCertEnabled),
			HTTPSOnly:               pointer.To(d.Get("https_only").(bool)),
			SiteConfig:              &siteConfig,
			VnetContentShareEnabled: pointer.To(d.Get("vnet_content_share_enabled").(bool)),
		},
	}

	publicNetworkAccess := d.Get("public_network_access").(string)
	if !features.FivePointOh() {
		// if a user is still using `site_config.public_network_access_enabled` we should be setting `public_network_access` for them
		publicNetworkAccess = reconcilePNA(d)
		if v := siteEnvelope.Properties.SiteConfig.PublicNetworkAccess; v != nil && *v == helpers.PublicNetworkAccessDisabled {
			publicNetworkAccess = helpers.PublicNetworkAccessDisabled
		}
	}

	// conversely if `public_network_access` has been set it should take precedence, and we should be propagating the value for that to `site_config.public_network_access_enabled`
	if publicNetworkAccess == helpers.PublicNetworkAccessDisabled {
		siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessDisabled)
	} else if publicNetworkAccess == helpers.PublicNetworkAccessEnabled {
		siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessEnabled)
	}

	siteEnvelope.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)

	if clientCertEnabled {
		siteEnvelope.Properties.ClientCertMode = pointer.To(webapps.ClientCertMode(clientCertMode))
	}

	if v := d.Get("virtual_network_subnet_id").(string); v != "" {
		siteEnvelope.Properties.VirtualNetworkSubnetId = pointer.To(v)
	}

	if _, ok := d.GetOk("identity"); ok {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		siteEnvelope.Identity = expandedIdentity
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, siteEnvelope); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if d.GetOk("storage_account"); ok {
		storageConfig := expandLogicAppStandardStorageConfig(d)
		if _, err := client.UpdateAzureStorageAccounts(ctx, id, *storageConfig); err != nil {
			return fmt.Errorf("setting Storage Accounts for %s: %+v", id, err)
		}
	}

	// This setting is enabled by default on creation of a logic app, we only need to update if config sets this as `false`
	if ftpAuth := d.Get("ftp_publish_basic_authentication_enabled").(bool); !ftpAuth {
		policy := webapps.CsmPublishingCredentialsPoliciesEntity{
			Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
				Allow: ftpAuth,
			},
		}

		if _, err := client.UpdateFtpAllowed(ctx, id, policy); err != nil {
			return fmt.Errorf("updating FTP publish basic authentication policy for %s: %+v", id, err)
		}
	}

	// This setting is enabled by default on creation of a logic app, we only need to update if config sets this as `false`
	if scmAuth := d.Get("scm_publish_basic_authentication_enabled").(bool); !scmAuth {
		policy := webapps.CsmPublishingCredentialsPoliciesEntity{
			Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
				Allow: scmAuth,
			},
		}

		if _, err := client.UpdateScmAllowed(ctx, id, policy); err != nil {
			return fmt.Errorf("updating SCM publish basic authentication policy for %s: %+v", id, err)
		}
	}

	return resourceLogicAppStandardUpdate(d, meta)
}

func resourceLogicAppStandardUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppService.WebAppsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	env := meta.(*clients.Client).Account.Environment
	storageAccountDomainSuffix, ok := env.Storage.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine the domain suffix for storage accounts in environment %q: %+v", env.Name, env.Storage)
	}

	id, err := commonids.ParseLogicAppId(d.Id())
	if err != nil {
		return err
	}

	clientCertMode := d.Get("client_certificate_mode").(string)
	clientCertEnabled := clientCertMode != ""

	basicAppSettings, err := getBasicLogicAppSettings(d, *storageAccountDomainSuffix)
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
	appSettings, err := expandLogicAppStandardSettings(d, *storageAccountDomainSuffix)
	if err != nil {
		return fmt.Errorf("expanding `app_settings`: %+v", err)
	}
	if vnetRouteAll, ok := appSettings["WEBSITE_VNET_ROUTE_ALL"]; ok {
		if !d.HasChange("site_config.0.vnet_route_all_enabled") {
			vnetRouteAllEnabled, _ := strconv.ParseBool(vnetRouteAll)
			siteConfig.VnetRouteAllEnabled = &vnetRouteAllEnabled
		}
	}

	siteEnvelope := webapps.Site{
		Kind:     &kind,
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &webapps.SiteProperties{
			ServerFarmId:          pointer.To(d.Get("app_service_plan_id").(string)),
			Enabled:               pointer.To(d.Get("enabled").(bool)),
			ClientAffinityEnabled: pointer.To(d.Get("client_affinity_enabled").(bool)),
			ClientCertEnabled:     pointer.To(clientCertEnabled),
			HTTPSOnly:             pointer.To(d.Get("https_only").(bool)),
			PublicNetworkAccess:   pointer.To(d.Get("public_network_access").(string)),
			SiteConfig:            &siteConfig,
		},
	}

	if d.HasChange("public_network_access") {
		publicNetworkAccess := d.Get("public_network_access").(string)
		siteEnvelope.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)
		if publicNetworkAccess == helpers.PublicNetworkAccessEnabled {
			siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessEnabled)
		} else {
			siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessDisabled)
		}
	}

	if d.HasChange("vnet_content_share_enabled") {
		siteEnvelope.Properties.VnetContentShareEnabled = pointer.To(d.Get("vnet_content_share_enabled").(bool))
	}

	if !features.FivePointOh() { // Until 5.0 the site_config value of this must be reflected back into the top-level property if not set there
		siteConfig.PublicNetworkAccess = pointer.To(reconcilePNA(d))
	}

	if clientCertEnabled {
		siteEnvelope.Properties.ClientCertMode = pointer.To(webapps.ClientCertMode(clientCertMode))
	}

	if d.HasChange("virtual_network_subnet_id") {
		subnetId := d.Get("virtual_network_subnet_id").(string)
		if subnetId == "" {
			if _, err := client.DeleteSwiftVirtualNetwork(ctx, *id); err != nil {
				return fmt.Errorf("removing `virtual_network_subnet_id` association for %s: %+v", *id, err)
			}
			var empty *string
			siteEnvelope.Properties.VirtualNetworkSubnetId = empty
		} else {
			siteEnvelope.Properties.VirtualNetworkSubnetId = pointer.To(subnetId)
		}
	}

	if _, ok := d.GetOk("identity"); ok {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		siteEnvelope.Identity = expandedIdentity
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, siteEnvelope); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if d.HasChange("site_config") || (d.HasChange("public_network_access") && !features.FivePointOh()) { // update siteConfig before appSettings in case the appSettings get covered by basicAppSettings
		siteConfigResource := webapps.SiteConfigResource{
			Properties: &siteConfig,
		}

		if _, err := client.CreateOrUpdateConfiguration(ctx, *id, siteConfigResource); err != nil {
			return fmt.Errorf("updating Configuration for %s: %+v", *id, err)
		}
	}

	settings := webapps.StringDictionary{
		Properties: pointer.To(appSettings),
	}

	if _, err = client.UpdateApplicationSettings(ctx, *id, settings); err != nil {
		return fmt.Errorf("updating Application Settings for %s: %+v", *id, err)
	}

	if d.HasChange("storage_account") {
		storageConfig := expandLogicAppStandardStorageConfig(d)
		if _, err := client.UpdateAzureStorageAccounts(ctx, *id, *storageConfig); err != nil {
			return fmt.Errorf("setting Storage Accounts for %s: %+v", id, err)
		}
	}

	if d.HasChange("connection_string") {
		connectionStrings := expandLogicAppStandardConnectionStrings(d)
		properties := webapps.ConnectionStringDictionary{
			Properties: pointer.To(connectionStrings),
		}

		if _, err := client.UpdateConnectionStrings(ctx, *id, properties); err != nil {
			return fmt.Errorf("updating Connection Strings for %s: %+v", *id, err)
		}
	}

	// HasChange will return `true` when config specifies this argument as `true` during initial creation.
	// To avoid unnecessary updates, check if the resource is new.
	if d.HasChange("ftp_publish_basic_authentication_enabled") && !d.IsNewResource() {
		policy := webapps.CsmPublishingCredentialsPoliciesEntity{
			Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
				Allow: d.Get("ftp_publish_basic_authentication_enabled").(bool),
			},
		}

		if _, err := client.UpdateFtpAllowed(ctx, *id, policy); err != nil {
			return fmt.Errorf("updating FTP publish basic authentication policy for %s: %+v", id, err)
		}
	}

	// HasChange will return `true` when config specifies this argument as `true` during initial creation.
	// To avoid unnecessary updates, check if the resource is new.
	if d.HasChange("scm_publish_basic_authentication_enabled") && !d.IsNewResource() {
		policy := webapps.CsmPublishingCredentialsPoliciesEntity{
			Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
				Allow: d.Get("scm_publish_basic_authentication_enabled").(bool),
			},
		}

		if _, err := client.UpdateScmAllowed(ctx, *id, policy); err != nil {
			return fmt.Errorf("updating SCM publish basic authentication policy for %s: %+v", id, err)
		}
	}

	return resourceLogicAppStandardRead(d, meta)
}

func resourceLogicAppStandardRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppService.WebAppsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseLogicAppId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.SiteName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("kind", pointer.From(model.Kind))
		d.Set("location", location.Normalize(model.Location))

		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %s", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}

		if props := model.Properties; props != nil {
			servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(*props.ServerFarmId)
			if err != nil {
				return err
			}
			d.Set("app_service_plan_id", servicePlanId.ID())
			d.Set("enabled", pointer.From(props.Enabled))
			d.Set("default_hostname", pointer.From(props.DefaultHostName))
			d.Set("https_only", pointer.From(props.HTTPSOnly))
			d.Set("outbound_ip_addresses", pointer.From(props.OutboundIPAddresses))
			d.Set("possible_outbound_ip_addresses", pointer.From(props.PossibleOutboundIPAddresses))
			d.Set("client_affinity_enabled", pointer.From(props.ClientAffinityEnabled))
			d.Set("custom_domain_verification_id", pointer.From(props.CustomDomainVerificationId))
			d.Set("virtual_network_subnet_id", pointer.From(props.VirtualNetworkSubnetId))
			d.Set("vnet_content_share_enabled", pointer.From(props.VnetContentShareEnabled))
			d.Set("public_network_access", pointer.From(props.PublicNetworkAccess))

			clientCertMode := ""
			if props.ClientCertEnabled != nil && *props.ClientCertEnabled {
				clientCertMode = string(pointer.From(props.ClientCertMode))
			}
			d.Set("client_certificate_mode", clientCertMode)
		}
	}

	appSettingsResp, err := client.ListApplicationSettings(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing application settings for %s: %+v", *id, err)
	}

	if model := appSettingsResp.Model; model != nil {
		appSettings := pointer.From(model.Properties)

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
	}

	connectionStringsResp, err := client.ListConnectionStrings(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing connection strings for %s: %+v", *id, err)
	}

	if model := connectionStringsResp.Model; model != nil {
		if err = d.Set("connection_string", flattenLogicAppStandardConnectionStrings(model.Properties)); err != nil {
			return err
		}
	}

	ftpBasicAuth, err := client.GetFtpAllowed(ctx, *id)
	if err != nil || ftpBasicAuth.Model == nil {
		return fmt.Errorf("retrieving FTP publish basic authentication policy for %s: %+v", id, err)
	}

	if props := ftpBasicAuth.Model.Properties; props != nil {
		d.Set("ftp_publish_basic_authentication_enabled", props.Allow)
	}

	scmBasicAuth, err := client.GetScmAllowed(ctx, *id)
	if err != nil || scmBasicAuth.Model == nil {
		return fmt.Errorf("retrieving SCM publish basic authentication policy for %s: %+v", id, err)
	}

	if props := scmBasicAuth.Model.Properties; props != nil {
		d.Set("scm_publish_basic_authentication_enabled", props.Allow)
	}

	siteCredentials, err := helpers.ListPublishingCredentials(ctx, client, *id)
	if err != nil {
		return fmt.Errorf("listing publishing credentials for %s: %+v", *id, err)
	}

	if err = d.Set("site_credential", flattenLogicAppStandardSiteCredential(siteCredentials)); err != nil {
		return err
	}

	configResp, err := client.GetConfiguration(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving the configuration for %s: %+v", *id, err)
	}

	if model := configResp.Model; model != nil {
		siteConfig := flattenLogicAppStandardSiteConfig(model.Properties)
		if err = d.Set("site_config", siteConfig); err != nil {
			return err
		}
	}

	storageAccountsResp, err := client.ListAzureStorageAccounts(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing Azure Storage Accounts for %s: %+v", *id, err)
	}

	if model := storageAccountsResp.Model; model != nil {
		storageAccounts := flattenLogicAppStandardStorageConfig(storageAccountsResp.Model)
		if err = d.Set("storage_account", storageAccounts); err != nil {
			return err
		}
	}

	return nil
}

func resourceLogicAppStandardDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppService.WebAppsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseLogicAppId(d.Id())
	if err != nil {
		return err
	}

	opts := webapps.DefaultDeleteOperationOptions()
	opts.DeleteMetrics = pointer.To(true)
	opts.DeleteEmptyServerFarm = pointer.To(false)

	if _, err := client.Delete(ctx, *id, opts); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func getBasicLogicAppSettings(d *pluginsdk.ResourceData, endpointSuffix string) ([]webapps.NameValuePair, error) {
	storagePropName := "AzureWebJobsStorage"
	functionVersionPropName := "FUNCTIONS_EXTENSION_VERSION"
	contentSharePropName := "WEBSITE_CONTENTSHARE"
	contentFileConnStringPropName := "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"
	appKindPropName := "APP_KIND"
	appKindPropValue := "workflowApp"

	storageAccount := d.Get("storage_account_name").(string)
	accountKey := d.Get("storage_account_access_key").(string)
	storageConnection := fmt.Sprintf(
		"DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s",
		storageAccount,
		accountKey,
		endpointSuffix,
	)
	functionVersion := d.Get("version").(string)

	contentShare := strings.ToLower(d.Get("name").(string)) + "-content"
	if _, ok := d.GetOk("storage_account_share_name"); ok {
		contentShare = d.Get("storage_account_share_name").(string)
	}

	basicSettings := []webapps.NameValuePair{
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
			return nil, fmt.Errorf(
				"when `use_extension_bundle` is true, `bundle_version` must be specified",
			)
		}

		bundleSettings := []webapps.NameValuePair{
			{Name: &extensionBundlePropName, Value: &extensionBundleName},
			{Name: &extensionBundleVersionPropName, Value: &extensionBundleVersion},
		}

		return append(basicSettings, bundleSettings...), nil
	}

	return basicSettings, nil
}

func schemaLogicAppStandardSiteConfig() *pluginsdk.Schema {
	schema := &pluginsdk.Schema{
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
						string(webapps.FtpsStateAllAllowed),
						string(webapps.FtpsStateDisabled),
						string(webapps.FtpsStateFtpsOnly),
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
						string(webapps.SupportedTlsVersionsOnePointTwo),
					}, false),
				},

				"pre_warmed_instance_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(0, 20),
				},

				"scm_ip_restriction": schemaLogicAppStandardIpRestriction(),

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"scm_min_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(webapps.SupportedTlsVersionsOnePointTwo),
					}, false),
				},

				"scm_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(webapps.ScmTypeBitbucketGit),
						string(webapps.ScmTypeBitbucketHg),
						string(webapps.ScmTypeCodePlexGit),
						string(webapps.ScmTypeCodePlexHg),
						string(webapps.ScmTypeDropbox),
						string(webapps.ScmTypeExternalGit),
						string(webapps.ScmTypeExternalHg),
						string(webapps.ScmTypeGitHub),
						string(webapps.ScmTypeLocalGit),
						string(webapps.ScmTypeNone),
						string(webapps.ScmTypeOneDrive),
						string(webapps.ScmTypeTfs),
						string(webapps.ScmTypeVSO),
						string(webapps.ScmTypeVSTSRM),
					}, false),
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
						"v8.0",
					}, false),
				},

				"vnet_route_all_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
				},

				"auto_swap_slot_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}

	if !features.FivePointOh() {
		schema.Elem.(*pluginsdk.Resource).Schema["public_network_access_enabled"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeBool,
			Optional:   true,
			Computed:   true,
			Deprecated: "the `site_config.public_network_access_enabled` property has been superseded by the `public_network_access` property and will be removed in v5.0 of the AzureRM Provider.",
		}
		schema.Elem.(*pluginsdk.Resource).Schema["scm_min_tls_version"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(webapps.SupportedTlsVersionsOnePointZero),
				string(webapps.SupportedTlsVersionsOnePointOne),
				string(webapps.SupportedTlsVersionsOnePointTwo),
			}, false),
		}
		schema.Elem.(*pluginsdk.Resource).Schema["min_tls_version"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(webapps.SupportedTlsVersionsOnePointZero),
				string(webapps.SupportedTlsVersionsOnePointOne),
				string(webapps.SupportedTlsVersionsOnePointTwo),
			}, false),
		}
	}

	return schema
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
					ValidateFunc: validation.IntBetween(1, math.MaxInt32),
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

				// lintignore:XS003
				"headers": {
					Type:       pluginsdk.TypeList,
					Optional:   true,
					Computed:   true,
					MaxItems:   1,
					ConfigMode: pluginsdk.SchemaConfigModeAttr,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							// lintignore:S018
							"x_forwarded_host": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								MaxItems: 8,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},

							// lintignore:S018
							"x_forwarded_for": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								MaxItems: 8,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.IsCIDR,
								},
							},

							// lintignore:S018
							"x_azure_fdid": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								MaxItems: 8,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.IsUUID,
								},
							},

							// lintignore:S018
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

func flattenLogicAppStandardConnectionStrings(input *map[string]webapps.ConnStringValueTypePair) interface{} {
	results := make([]interface{}, 0)

	if input == nil || len(*input) == 0 {
		return results
	}

	for k, v := range *input {
		result := make(map[string]interface{})
		result["name"] = k
		result["type"] = string(v.Type)
		result["value"] = v.Value
		results = append(results, result)
	}

	return results
}

func flattenLogicAppStandardStorageConfig(input *webapps.AzureStoragePropertyDictionaryResource) []interface{} {
	results := make([]interface{}, 0)

	if input == nil {
		return results
	}

	for k, v := range *input.Properties {
		result := make(map[string]interface{})
		result["name"] = k
		result["type"] = string(pointer.From(v.Type))
		result["account_name"] = pointer.From(v.AccountName)
		result["share_name"] = v.ShareName
		result["access_key"] = v.AccessKey
		result["mount_path"] = v.MountPath
		results = append(results, result)
	}

	return results
}

func flattenLogicAppStandardSiteConfig(input *webapps.SiteConfig) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{})

	if input == nil {
		log.Printf("[DEBUG] SiteConfig is nil")
		return results
	}

	result["always_on"] = pointer.From(input.AlwaysOn)
	result["use_32_bit_worker_process"] = pointer.From(input.Use32BitWorkerProcess)
	result["websockets_enabled"] = pointer.From(input.WebSocketsEnabled)
	result["linux_fx_version"] = pointer.From(input.LinuxFxVersion)
	result["http2_enabled"] = pointer.From(input.HTTP20Enabled)
	result["pre_warmed_instance_count"] = pointer.From(input.PreWarmedInstanceCount)

	result["ip_restriction"] = flattenLogicAppStandardIpRestriction(input.IPSecurityRestrictions)

	result["scm_ip_restriction"] = flattenLogicAppStandardIpRestriction(input.ScmIPSecurityRestrictions)

	result["scm_use_main_ip_restriction"] = pointer.From(input.ScmIPSecurityRestrictionsUseMain)

	result["scm_type"] = string(pointer.From(input.ScmType))
	result["scm_min_tls_version"] = string(pointer.From(input.ScmMinTlsVersion))

	result["min_tls_version"] = string(pointer.From(input.MinTlsVersion))
	result["ftps_state"] = string(pointer.From(input.FtpsState))

	result["cors"] = flattenLogicAppStandardCorsSettings(input.Cors)

	result["auto_swap_slot_name"] = pointer.From(input.AutoSwapSlotName)

	result["health_check_path"] = pointer.From(input.HealthCheckPath)

	result["elastic_instance_minimum"] = pointer.From(input.MinimumElasticInstanceCount)

	result["app_scale_limit"] = pointer.From(input.FunctionAppScaleLimit)

	result["runtime_scale_monitoring_enabled"] = pointer.From(input.FunctionsRuntimeScaleMonitoringEnabled)

	result["dotnet_framework_version"] = pointer.From(input.NetFrameworkVersion)

	result["vnet_route_all_enabled"] = pointer.From(input.VnetRouteAllEnabled)

	publicNetworkAccessEnabled := true
	if input.PublicNetworkAccess != nil {
		publicNetworkAccessEnabled = !strings.EqualFold(pointer.From(input.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled)
	}

	if !features.FivePointOh() {
		result["public_network_access_enabled"] = publicNetworkAccessEnabled
	}

	results = append(results, result)
	return results
}

func flattenLogicAppStandardSiteCredential(input *webapps.User) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{})

	if input == nil || input.Properties == nil {
		log.Printf("[DEBUG] UserProperties is nil")
		return results
	}

	result["username"] = input.Properties.PublishingUserName

	result["password"] = pointer.From(input.Properties.PublishingPassword)

	return append(results, result)
}

func flattenLogicAppStandardIpRestriction(input *[]webapps.IPSecurityRestriction) []interface{} {
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
				switch pointer.From(v.Tag) {
				case webapps.IPFilterTagServiceTag:
					restriction["service_tag"] = *ip
				default:
					restriction["ip_address"] = *ip
				}
			}
		}

		subnetId := ""
		if subnetIdRaw := v.VnetSubnetResourceId; subnetIdRaw != nil {
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
			restriction["headers"] = flattenHeaders(*headers)
		}

		restrictions = append(restrictions, restriction)
	}

	return restrictions
}

func flattenLogicAppStandardCorsSettings(input *webapps.CorsSettings) []interface{} {
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

func expandLogicAppStandardSiteConfig(d *pluginsdk.ResourceData) (webapps.SiteConfig, error) {
	configs := d.Get("site_config").([]interface{})
	siteConfig := webapps.SiteConfig{}

	if len(configs) == 0 {
		return siteConfig, nil
	}

	config := configs[0].(map[string]interface{})

	if v, ok := config["always_on"]; ok {
		siteConfig.AlwaysOn = pointer.To(v.(bool))
	}

	if v, ok := config["use_32_bit_worker_process"]; ok {
		siteConfig.Use32BitWorkerProcess = pointer.To(v.(bool))
	}

	if v, ok := config["websockets_enabled"]; ok {
		siteConfig.WebSocketsEnabled = pointer.To(v.(bool))
	}

	if v, ok := config["linux_fx_version"]; ok {
		siteConfig.LinuxFxVersion = pointer.To(v.(string))
	}

	if v, ok := config["cors"]; ok {
		expand := expandLogicAppStandardCorsSettings(v)
		siteConfig.Cors = &expand
	}

	if v, ok := config["http2_enabled"]; ok {
		siteConfig.HTTP20Enabled = pointer.To(v.(bool))
	}

	if v, ok := config["ip_restriction"]; ok {
		restrictions, err := expandLogicAppStandardIpRestriction(v)
		if err != nil {
			return siteConfig, err
		}
		siteConfig.IPSecurityRestrictions = &restrictions
	}

	if v, ok := config["scm_ip_restriction"]; ok {
		scmIPSecurityRestrictions := v.([]interface{})
		scmRestrictions, err := expandLogicAppStandardIpRestriction(scmIPSecurityRestrictions)
		if err != nil {
			return siteConfig, err
		}
		siteConfig.ScmIPSecurityRestrictions = &scmRestrictions
	}

	if v, ok := config["scm_use_main_ip_restriction"]; ok {
		siteConfig.ScmIPSecurityRestrictionsUseMain = pointer.To(v.(bool))
	}

	if v, ok := config["scm_min_tls_version"]; ok {
		siteConfig.ScmMinTlsVersion = pointer.To(webapps.SupportedTlsVersions(v.(string)))
	}

	if v, ok := config["scm_type"]; ok {
		siteConfig.ScmType = pointer.To(webapps.ScmType(v.(string)))
	}

	if v, ok := config["min_tls_version"]; ok {
		siteConfig.MinTlsVersion = pointer.To(webapps.SupportedTlsVersions(v.(string)))
	}

	if v, ok := config["ftps_state"]; ok {
		siteConfig.FtpsState = pointer.To(webapps.FtpsState(v.(string)))
	}

	// get value from `d` rather than the `config` map, or it will be covered by the zero-value "0" instead of nil.
	if v, ok := d.GetOk("site_config.0.pre_warmed_instance_count"); ok {
		siteConfig.PreWarmedInstanceCount = pointer.To(int64(v.(int)))
	}

	if v, ok := config["health_check_path"]; ok {
		siteConfig.HealthCheckPath = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("site_config.0.elastic_instance_minimum"); ok {
		siteConfig.MinimumElasticInstanceCount = pointer.To(int64(v.(int)))
	}

	if v, ok := d.GetOk("site_config.0.app_scale_limit"); ok {
		siteConfig.FunctionAppScaleLimit = pointer.To(int64(v.(int)))
	}

	if v, ok := config["runtime_scale_monitoring_enabled"]; ok {
		siteConfig.FunctionsRuntimeScaleMonitoringEnabled = pointer.To(v.(bool))
	}

	if v, ok := config["dotnet_framework_version"]; ok {
		siteConfig.NetFrameworkVersion = pointer.To(v.(string))
	}

	if v, ok := config["vnet_route_all_enabled"]; ok {
		siteConfig.VnetRouteAllEnabled = pointer.To(v.(bool))
	}

	if !features.FivePointOh() {
		siteConfig.PublicNetworkAccess = pointer.To(reconcilePNA(d))
	}

	return siteConfig, nil
}

func expandLogicAppStandardSettings(d *pluginsdk.ResourceData, endpointSuffix string) (map[string]string, error) {
	output := make(map[string]string)
	appSettings := expandAppSettings(d)
	basicAppSettings, err := getBasicLogicAppSettings(d, endpointSuffix)
	if err != nil {
		return nil, err
	}
	for _, p := range append(basicAppSettings, appSettings...) {
		output[*p.Name] = pointer.From(p.Value)
	}

	return output, nil
}

func expandAppSettings(d *pluginsdk.ResourceData) []webapps.NameValuePair {
	input := d.Get("app_settings").(map[string]interface{})
	output := make([]webapps.NameValuePair, 0)

	for k, v := range input {
		nameValue := webapps.NameValuePair{
			Name:  pointer.To(k),
			Value: pointer.To(v.(string)),
		}
		output = append(output, nameValue)
	}

	return output
}

func expandLogicAppStandardConnectionStrings(d *pluginsdk.ResourceData) map[string]webapps.ConnStringValueTypePair {
	input := d.Get("connection_string").(*pluginsdk.Set).List()
	output := make(map[string]webapps.ConnStringValueTypePair, len(input))

	for _, v := range input {
		vals := v.(map[string]interface{})

		csName := vals["name"].(string)
		csType := vals["type"].(string)
		csValue := vals["value"].(string)

		output[csName] = webapps.ConnStringValueTypePair{
			Value: csValue,
			Type:  webapps.ConnectionStringType(csType),
		}
	}

	return output
}

func expandLogicAppStandardCorsSettings(input interface{}) webapps.CorsSettings {
	settings := input.([]interface{})
	corsSettings := webapps.CorsSettings{}

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
		corsSettings.SupportCredentials = pointer.To(v.(bool))
	}

	return corsSettings
}

func expandLogicAppStandardStorageConfig(d *pluginsdk.ResourceData) *webapps.AzureStoragePropertyDictionaryResource {
	storageConfigs := d.Get("storage_account").(*pluginsdk.Set).List()
	storageAccounts := make(map[string]webapps.AzureStorageInfoValue)
	result := &webapps.AzureStoragePropertyDictionaryResource{}

	for _, v := range storageConfigs {
		config := v.(map[string]interface{})
		storageAccounts[config["name"].(string)] = webapps.AzureStorageInfoValue{
			Type:        pointer.To(webapps.AzureStorageType(config["type"].(string))),
			AccountName: pointer.To(config["account_name"].(string)),
			ShareName:   pointer.To(config["share_name"].(string)),
			AccessKey:   pointer.To(config["access_key"].(string)),
			MountPath:   pointer.To(config["mount_path"].(string)),
		}
	}
	result.Properties = &storageAccounts

	return result
}

func expandLogicAppStandardIpRestriction(input interface{}) ([]webapps.IPSecurityRestriction, error) {
	restrictions := make([]webapps.IPSecurityRestriction, 0)

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
			return nil, fmt.Errorf(
				"only one of `ip_address`, `service_tag` or `virtual_network_subnet_id` can be set for an IP restriction",
			)
		}

		if vNetSubnetID == "" && ipAddress == "" && serviceTag == "" {
			return nil, fmt.Errorf(
				"one of `ip_address`, `service_tag` or `virtual_network_subnet_id` must be set for an IP restriction",
			)
		}

		ipSecurityRestriction := webapps.IPSecurityRestriction{}
		if ipAddress == "Any" {
			continue
		}

		if ipAddress != "" {
			ipSecurityRestriction.IPAddress = &ipAddress
		}

		if serviceTag != "" {
			ipSecurityRestriction.IPAddress = &serviceTag
			ipSecurityRestriction.Tag = pointer.To(webapps.IPFilterTagServiceTag)
		}

		if vNetSubnetID != "" {
			ipSecurityRestriction.VnetSubnetResourceId = &vNetSubnetID
		}

		if name != "" {
			ipSecurityRestriction.Name = &name
		}

		if priority != 0 {
			ipSecurityRestriction.Priority = pointer.To(int64(priority))
		}

		if action != "" {
			ipSecurityRestriction.Action = &action
		}
		if headers, ok := restriction["headers"]; ok {
			ipSecurityRestriction.Headers = pointer.To(expandHeaders(headers.([]interface{})))
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

func reconcilePNA(d *pluginsdk.ResourceData) string {
	pna := ""
	scPNASet := true
	if !d.GetRawConfig().AsValueMap()["public_network_access"].IsNull() { // is top level set, takes precedence
		pna = d.Get("public_network_access").(string)
	}
	if sc := d.GetRawConfig().AsValueMap()["site_config"]; !sc.IsNull() {
		if len(sc.AsValueSlice()) > 0 && !sc.AsValueSlice()[0].AsValueMap()["public_network_access_enabled"].IsNull() {
			scPNASet = true
		}
	}
	if pna == "" && scPNASet { // if not, or it's empty, is site_config value set
		pnaBool := d.Get("site_config.0.public_network_access_enabled").(bool)
		if pnaBool {
			pna = helpers.PublicNetworkAccessEnabled
		} else {
			pna = helpers.PublicNetworkAccessDisabled
		}
	}

	return pna
}
