// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"context"
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LogicAppResource struct{}

type LogicAppResourceModel struct {
	Name                       string                                     `tfschema:"name"`
	ResourceGroupName          string                                     `tfschema:"resource_group_name"`
	Location                   string                                     `tfschema:"location"`
	AppServicePlanId           string                                     `tfschema:"app_service_plan_id"`
	AppSettings                map[string]string                          `tfschema:"app_settings"`
	UseExtensionBundle         bool                                       `tfschema:"use_extension_bundle"`
	BundleVersion              string                                     `tfschema:"bundle_version"`
	ClientAffinityEnabled      bool                                       `tfschema:"client_affinity_enabled"`
	ClientCertificateEnabled   bool                                       `tfschema:"client_certificate_enabled"`
	ClientCertificateMode      string                                     `tfschema:"client_certificate_mode"`
	Enabled                    bool                                       `tfschema:"enabled"`
	FtpPublishBasicAuthEnabled bool                                       `tfschema:"ftp_publish_basic_authentication_enabled"`
	HTTPSOnly                  bool                                       `tfschema:"https_only"`
	Identity                   []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	SCMPublishBasicAuthEnabled bool                                       `tfschema:"scm_publish_basic_authentication_enabled"`
	SiteConfig                 []helpers.LogicAppSiteConfig               `tfschema:"site_config"`
	ConnectionStrings          []helpers.ConnectionString                 `tfschema:"connection_string"`
	StorageAccountName         string                                     `tfschema:"storage_account_name"`
	StorageAccountAccessKey    string                                     `tfschema:"storage_account_access_key"`
	PublicNetworkAccess        string                                     `tfschema:"public_network_access"`
	StorageAccountShareName    string                                     `tfschema:"storage_account_share_name"`
	Version                    string                                     `tfschema:"version"`
	VNETContentShareEnabled    bool                                       `tfschema:"vnet_content_share_enabled"`
	VirtualNetworkSubnetId     string                                     `tfschema:"virtual_network_subnet_id"`
	Tags                       map[string]string                          `tfschema:"tags"`

	CustomDomainVerificationId  string                           `tfschema:"custom_domain_verification_id"`
	DefaultHostname             string                           `tfschema:"default_hostname"`
	Kind                        string                           `tfschema:"kind"`
	OutboundIpAddresses         string                           `tfschema:"outbound_ip_addresses"`
	PossibleOutboundIpAddresses string                           `tfschema:"possible_outbound_ip_addresses"`
	SiteCredential              []helpers.SiteCredentialLogicApp `tfschema:"site_credential"`
}

var (
	LogicAppStdKind   = "functionapp,workflowapp"
	LogicAppLinuxKind = "functionapp,linux,container,workflowapp"

	storageConnectionStringFmt           = "DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s"
	storageAppSettingName                = "AzureWebJobsStorage"
	contentShareAppSettingName           = "WEBSITE_CONTENTSHARE"
	contentFileConnStringAppSettingName  = "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"
	functionVersionAppSettingName        = "FUNCTIONS_EXTENSION_VERSION"
	extensionBundleAppSettingName        = "AzureFunctionsJobHost__extensionBundle__id"
	extensionBundleAppSettingValue       = "Microsoft.Azure.Functions.ExtensionBundle.Workflows"
	extensionBundleVersionAppSettingName = "AzureFunctionsJobHost__extensionBundle__version"
)

func (r LogicAppResource) Arguments() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
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

		"client_certificate_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Should the logic app use Client Certificates",
		},

		"client_certificate_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      webapps.ClientCertModeRequired,
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForClientCertMode(), false),
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

		"site_config": helpers.SchemaLogicAppStandardSiteConfig(),

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
	}

	if !features.FivePointOh() {
		s["client_certificate_mode"].Default = nil
		s["public_network_access"].Default = nil
		s["public_network_access"].Computed = true
	}

	return s
}

func (r LogicAppResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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
	}
}

func (r LogicAppResource) ModelObject() interface{} {
	return &LogicAppResourceModel{}
}

func (r LogicAppResource) ResourceType() string {
	return "azurerm_logic_app_standard"
}

func (r LogicAppResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			resourcesClient := metadata.Client.AppService.ResourceProvidersClient
			servicePlanClient := metadata.Client.AppService.ServicePlanClient
			aseClient := metadata.Client.AppService.AppServiceEnvironmentClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			env := metadata.Client.Account.Environment
			storageAccountDomainSuffix, ok := env.Storage.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine the domain suffix for storage accounts in environment %q: %+v", env.Name, env.Storage)
			}

			data := LogicAppResourceModel{}

			if err := metadata.Decode(&data); err != nil {
				return err
			}

			id := commonids.NewAppServiceID(subscriptionId, data.ResourceGroupName, data.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_logic_app_standard", id.ID())
			}

			servicePlanId, err := commonids.ParseAppServicePlanID(data.AppServicePlanId)
			if err != nil {
				return err
			}

			servicePlan, err := servicePlanClient.Get(ctx, *servicePlanId)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", servicePlanId, err)
			}

			availabilityRequest := resourceproviders.ResourceNameAvailabilityRequest{
				Name: data.Name,
				Type: resourceproviders.CheckNameResourceTypesMicrosoftPointWebSites,
			}
			if servicePlanModel := servicePlan.Model; servicePlanModel != nil {
				if ase := servicePlanModel.Properties.HostingEnvironmentProfile; ase != nil {
					// Attempt to check the ASE for the appropriate suffix for the name availability request.
					// This varies between internal and external ASE Types, and potentially has other names in other clouds
					// We use the "internal" as the fallback here, if we can read the ASE, we'll get the full one
					nameSuffix := "appserviceenvironment.net"
					if ase.Id != nil {
						aseId, err := commonids.ParseAppServiceEnvironmentIDInsensitively(*ase.Id)
						nameSuffix = fmt.Sprintf("%s.%s", aseId.HostingEnvironmentName, nameSuffix)
						if err != nil {
							metadata.Logger.Warnf("could not parse App Service Environment ID determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", data.Name, servicePlanId)
						} else {
							existingASE, err := aseClient.Get(ctx, *aseId)
							if err != nil || existingASE.Model == nil {
								metadata.Logger.Warnf("could not read App Service Environment to determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", data.Name, servicePlanId)
							} else if props := existingASE.Model.Properties; props != nil && props.DnsSuffix != nil && *props.DnsSuffix != "" {
								nameSuffix = *props.DnsSuffix
							}
						}
					}

					availabilityRequest.Name = fmt.Sprintf("%s.%s", data.Name, nameSuffix)
					availabilityRequest.IsFqdn = pointer.To(true)
				}
			}

			subId := commonids.NewSubscriptionID(subscriptionId)

			checkName, err := resourcesClient.CheckNameAvailability(ctx, subId, availabilityRequest)
			if err != nil {
				return fmt.Errorf("checking name availability for Linux %s: %+v", id, err)
			}
			if model := checkName.Model; model != nil && model.NameAvailable != nil && !*model.NameAvailable {
				return fmt.Errorf("the Site Name %q failed the availability check: %+v", id.SiteName, *model.Message)
			}

			clientCertEnabled := data.ClientCertificateEnabled
			if !features.FivePointOh() {
				clientCertEnabled = data.ClientCertificateMode != "" || data.ClientCertificateEnabled
			}

			basicAppSettings, err := getBasicLogicAppSettings(data, *storageAccountDomainSuffix)
			if err != nil {
				return err
			}

			siteConfig, err := expandLogicAppStandardSiteConfigForCreate(data.SiteConfig, metadata)
			if err != nil {
				return fmt.Errorf("expanding `site_config`: %+v", err)
			}

			kind := LogicAppStdKind
			if siteConfig.LinuxFxVersion != nil && len(*siteConfig.LinuxFxVersion) > 0 {
				kind = LogicAppLinuxKind
			}

			appSettings := expandAppSettings(data.AppSettings)
			appSettings = append(appSettings, basicAppSettings...)

			siteConfig.AppSettings = pointer.To(appSettings)

			if v, ok := data.AppSettings["WEBSITE_VNET_ROUTE_ALL"]; ok {
				// For compatibility between app_settings and site_config, we need to set the API property based on the presence of the app_setting map value if present.
				// a replacement of this resource should consider deprecating support for this.
				vnetRouteAll, _ := strconv.ParseBool(v)
				siteConfig.VnetRouteAllEnabled = pointer.To(vnetRouteAll)
			}

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(data.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			siteEnvelope := webapps.Site{
				Identity: expandedIdentity,
				Kind:     pointer.To(kind),
				Location: location.Normalize(data.Location),
				Properties: &webapps.SiteProperties{
					ServerFarmId:            pointer.To(data.AppServicePlanId),
					Enabled:                 pointer.To(data.Enabled),
					ClientAffinityEnabled:   pointer.To(data.ClientAffinityEnabled),
					ClientCertEnabled:       pointer.To(clientCertEnabled),
					HTTPSOnly:               pointer.To(data.HTTPSOnly),
					SiteConfig:              siteConfig,
					VnetContentShareEnabled: pointer.To(data.VNETContentShareEnabled),
					PublicNetworkAccess:     pointer.To(data.PublicNetworkAccess),
				},
				Tags: pointer.To(data.Tags),
			}

			if !features.FivePointOh() {
				publicNetworkAccess := data.PublicNetworkAccess
				// if a user is still using `site_config.public_network_access_enabled` we should be setting `public_network_access` for them
				publicNetworkAccess = reconcilePNA(metadata)
				if v := siteEnvelope.Properties.SiteConfig.PublicNetworkAccess; v != nil && *v == helpers.PublicNetworkAccessDisabled {
					publicNetworkAccess = helpers.PublicNetworkAccessDisabled
				}
				// conversely if `public_network_access` has been set it should take precedence, and we should be propagating the value for that to `site_config.public_network_access_enabled`
				if publicNetworkAccess == helpers.PublicNetworkAccessDisabled {
					siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessDisabled)
				} else if publicNetworkAccess == helpers.PublicNetworkAccessEnabled {
					siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessEnabled)
				}
				siteEnvelope.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)
			}

			if clientCertEnabled {
				siteEnvelope.Properties.ClientCertMode = pointer.ToEnum[webapps.ClientCertMode](data.ClientCertificateMode)
			}

			if data.VirtualNetworkSubnetId != "" {
				siteEnvelope.Properties.VirtualNetworkSubnetId = pointer.To(data.VirtualNetworkSubnetId)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, siteEnvelope); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			if !data.FtpPublishBasicAuthEnabled {
				policy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: data.FtpPublishBasicAuthEnabled,
					},
				}

				if _, err := client.UpdateFtpAllowed(ctx, id, policy); err != nil {
					return fmt.Errorf("updating FTP publish basic authentication policy for %s: %+v", id, err)
				}
			}

			if !data.SCMPublishBasicAuthEnabled {
				policy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: data.SCMPublishBasicAuthEnabled,
					},
				}

				if _, err := client.UpdateScmAllowed(ctx, id, policy); err != nil {
					return fmt.Errorf("updating FTP publish basic authentication policy for %s: %+v", id, err)
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(data.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStrings(ctx, id, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for Linux %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r LogicAppResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			var state LogicAppResourceModel

			id, err := commonids.ParseLogicAppId(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Linux %s: %+v", id, err)
			}

			state.Name = id.SiteName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Kind = pointer.From(model.Kind)
				state.Location = location.Normalize(model.Location)
				ident, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return err
				}
				state.Identity = pointer.From(ident)
				state.Tags = pointer.From(model.Tags)
				if props := model.Properties; props != nil {
					servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(*props.ServerFarmId)
					if err != nil {
						return err
					}

					state.AppServicePlanId = servicePlanId.ID()
					state.Enabled = pointer.From(props.Enabled)
					state.DefaultHostname = pointer.From(props.DefaultHostName)
					state.HTTPSOnly = pointer.From(props.HTTPSOnly)
					state.OutboundIpAddresses = pointer.From(props.OutboundIPAddresses)
					state.PossibleOutboundIpAddresses = pointer.From(props.PossibleOutboundIPAddresses)
					state.ClientAffinityEnabled = pointer.From(props.ClientAffinityEnabled)
					state.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
					state.VirtualNetworkSubnetId = pointer.From(props.VirtualNetworkSubnetId)
					state.VNETContentShareEnabled = pointer.From(props.VnetContentShareEnabled)
					state.PublicNetworkAccess = pointer.From(props.PublicNetworkAccess)
					if !features.FivePointOh() { // Maintaining the bugged implementation for compatibility until 5.0
						if pointer.From(props.ClientCertEnabled) {
							state.ClientCertificateMode = pointer.FromEnum(props.ClientCertMode)
						}
					} else {
						state.ClientCertificateMode = pointer.FromEnum(props.ClientCertMode)
					}
				}
			}

			appSettingsResp, err := client.ListApplicationSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing application settings for %s: %+v", *id, err)
			}

			if model := appSettingsResp.Model; model != nil {
				appSettings := pointer.From(model.Properties)

				connectionString := appSettings[storageAppSettingName]

				for _, part := range strings.Split(connectionString, ";") {
					if strings.HasPrefix(part, "AccountName") {
						accountNameParts := strings.Split(part, "AccountName=")
						if len(accountNameParts) > 1 {
							state.StorageAccountName = accountNameParts[1]
						}
					}
					if strings.HasPrefix(part, "AccountKey") {
						accountKeyParts := strings.Split(part, "AccountKey=")
						if len(accountKeyParts) > 1 {
							state.StorageAccountAccessKey = accountKeyParts[1]
						}
					}
				}

				if v, ok := appSettings[functionVersionAppSettingName]; ok {
					state.Version = v
				}

				if _, ok := appSettings["AzureFunctionsJobHost__extensionBundle__id"]; ok {
					state.UseExtensionBundle = true

					if val, ok := appSettings["AzureFunctionsJobHost__extensionBundle__version"]; ok {
						state.BundleVersion = val
					}
				} else {
					state.UseExtensionBundle = false
					state.BundleVersion = "[1.*, 2.0.0)"
				}

				state.StorageAccountShareName = appSettings[contentShareAppSettingName]
				delete(appSettings, contentFileConnStringAppSettingName)
				delete(appSettings, "APP_KIND")
				delete(appSettings, "AzureFunctionsJobHost__extensionBundle__id")
				delete(appSettings, "AzureFunctionsJobHost__extensionBundle__version")
				delete(appSettings, "AzureWebJobsDashboard")
				delete(appSettings, storageAppSettingName)
				delete(appSettings, functionVersionAppSettingName)
				delete(appSettings, contentShareAppSettingName)

				state.AppSettings = appSettings
			}

			connectionStringsResp, err := client.ListConnectionStrings(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing connection strings for %s: %+v", *id, err)
			}

			if model := connectionStringsResp.Model; model != nil {
				state.ConnectionStrings = helpers.FlattenConnectionStrings(model)
			}

			ftpBasicAuth, err := client.GetFtpAllowed(ctx, *id)
			if err != nil || ftpBasicAuth.Model == nil {
				return fmt.Errorf("retrieving FTP publish basic authentication policy for %s: %+v", id, err)
			}

			if props := ftpBasicAuth.Model.Properties; props != nil {
				state.FtpPublishBasicAuthEnabled = props.Allow
			}

			scmBasicAuth, err := client.GetScmAllowed(ctx, *id)
			if err != nil || scmBasicAuth.Model == nil {
				return fmt.Errorf("retrieving SCM publish basic authentication policy for %s: %+v", id, err)
			}

			if props := scmBasicAuth.Model.Properties; props != nil {
				state.SCMPublishBasicAuthEnabled = props.Allow
			}

			siteCredentials, err := helpers.ListPublishingCredentials(ctx, client, *id)
			if err != nil {
				return fmt.Errorf("listing publishing credentials for %s: %+v", *id, err)
			}

			state.SiteCredential = helpers.FlattenSiteCredentialsLogicApp(siteCredentials)

			configResp, err := client.GetConfiguration(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving the configuration for %s: %+v", *id, err)
			}

			if model := configResp.Model; model != nil {
				state.SiteConfig = flattenLogicAppStandardSiteConfig(model.Properties)
				// if !features.FivePointOh() {
				// 	if len(state.SiteConfig) > 0 {
				// 		if state.SiteConfig[0].PublicNetworkAccessEnabled {
				// 			state.PublicNetworkAccess = helpers.PublicNetworkAccessEnabled
				// 		} else {
				// 			state.PublicNetworkAccess = helpers.PublicNetworkAccessDisabled
				// 		}
				// 	}
				// }
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LogicAppResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := commonids.ParseLogicAppId(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting Linux %s", *id)

			delOptions := webapps.DeleteOperationOptions{
				DeleteMetrics:         pointer.To(true),
				DeleteEmptyServerFarm: pointer.To(false),
			}

			if _, err = client.Delete(ctx, *id, delOptions); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r LogicAppResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateLogicAppId
}

func (r LogicAppResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := commonids.ParseLogicAppId(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			data := LogicAppResourceModel{}

			if err := metadata.Decode(&data); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %v", *id, err)
			}
			if existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s for update", id)
			}

			siteEnvelope := *existing.Model.Properties

			sc, err := client.GetConfiguration(ctx, *id)
			if err != nil || sc.Model == nil {
				return fmt.Errorf("reading site_config for %s: %v", *id, err)
			}

			existingSiteConfig := sc.Model.Properties
			siteEnvelope.SiteConfig = existingSiteConfig

			appSettingsResp, err := client.ListApplicationSettings(ctx, *id)
			if err != nil || appSettingsResp.Model == nil {
				return fmt.Errorf("reading App Settings for Linux %s: %+v", id, err)
			}

			currentAppSettings := make([]webapps.NameValuePair, 0)
			if appSettingsResp.Model.Properties != nil {
				currentAppSettings = expandAppSettings(*appSettingsResp.Model.Properties)
			}
			existingSiteConfig.AppSettings = pointer.To(currentAppSettings)

			if metadata.ResourceData.HasChanges("site_config", "app_settings", "version", "storage_account_name", "storage_account_access_key") {
				existingSiteConfig, err = expandLogicAppStandardSiteConfigForUpdate(data.SiteConfig, metadata, existingSiteConfig)
				if err != nil {
					return fmt.Errorf("expanding site_config update for %s: %v", *id, err)
				}

				siteEnvelope.SiteConfig = existingSiteConfig
			}

			if metadata.ResourceData.HasChange("site_config.0.linux_fx_version") {
				kind := "functionapp,workflowapp"
				if metadata.ResourceData.Get("site_config.0.linux_fx_version").(string) != "" {
					kind = "functionapp,linux,container,workflowapp"
				}
				existing.Model.Kind = pointer.To(kind)
			}

			if metadata.ResourceData.HasChange("app_service_plan_id") {
				planId, err := commonids.ParseLogicAppIdInsensitively(metadata.ResourceData.Id())
				if err != nil {
					return err
				}

				siteEnvelope.ServerFarmId = pointer.To(planId.ID())
			}

			if metadata.ResourceData.HasChange("enabled") {
				siteEnvelope.Enabled = pointer.To(metadata.ResourceData.Get("enabled").(bool))
			}

			if metadata.ResourceData.HasChange("client_affinity_enabled") {
				siteEnvelope.ClientAffinityEnabled = pointer.To(metadata.ResourceData.Get("client_affinity_enabled").(bool))
			}

			if metadata.ResourceData.HasChanges("client_certificate_mode") {
				siteEnvelope.ClientCertMode = pointer.ToEnum[webapps.ClientCertMode](metadata.ResourceData.Get("client_certificate_mode").(string))
				siteEnvelope.ClientCertEnabled = pointer.To(metadata.ResourceData.Get("client_certificate_mode").(string) != "")
			}

			if metadata.ResourceData.HasChanges("https_only") {
				siteEnvelope.HTTPSOnly = pointer.To(metadata.ResourceData.Get("https_only").(bool))
			}

			if metadata.ResourceData.HasChange("public_network_access") {
				if strings.EqualFold(data.PublicNetworkAccess, helpers.PublicNetworkAccessEnabled) {
					siteEnvelope.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessEnabled)
					if !features.FivePointOh() {
						siteEnvelope.SiteConfig.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessEnabled)
					}
				} else {
					siteEnvelope.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessDisabled)
					if !features.FivePointOh() {
						siteEnvelope.SiteConfig.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessDisabled)
					}
				}
			}

			if metadata.ResourceData.HasChange("vnet_content_share_enabled") {
				siteEnvelope.VnetContentShareEnabled = pointer.To(data.VNETContentShareEnabled)
			}

			if metadata.ResourceData.HasChange("virtual_network_subnet_id") {
				subnetId := data.VirtualNetworkSubnetId
				if subnetId == "" {
					if _, err := client.DeleteSwiftVirtualNetwork(ctx, *id); err != nil {
						return fmt.Errorf("removing `virtual_network_subnet_id` association for %s: %+v", *id, err)
					}
					var empty *string
					siteEnvelope.VirtualNetworkSubnetId = empty
				} else {
					siteEnvelope.VirtualNetworkSubnetId = pointer.To(subnetId)
				}
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(data.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}

				existing.Model.Identity = expandedIdentity
			}

			existing.Model.Properties = pointer.To(siteEnvelope)

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = pointer.To(data.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if metadata.ResourceData.HasChange("connection_string") {
				connectionStrings := helpers.ExpandConnectionStrings(data.ConnectionStrings)
				if connectionStrings.Properties != nil {
					if _, err := client.UpdateConnectionStrings(ctx, *id, *connectionStrings); err != nil {
						return fmt.Errorf("setting Connection Strings for Linux %s: %+v", id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("ftp_publish_basic_authentication_enabled") {
				policy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: data.FtpPublishBasicAuthEnabled,
					},
				}

				if _, err := client.UpdateFtpAllowed(ctx, *id, policy); err != nil {
					return fmt.Errorf("updating FTP publish basic authentication policy for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("scm_publish_basic_authentication_enabled") {
				policy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: data.SCMPublishBasicAuthEnabled,
					},
				}

				if _, err := client.UpdateScmAllowed(ctx, *id, policy); err != nil {
					return fmt.Errorf("updating SCM publish basic authentication policy for %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

var _ sdk.ResourceWithUpdate = &LogicAppResource{}

func getBasicLogicAppSettings(d LogicAppResourceModel, endpointSuffix string) ([]webapps.NameValuePair, error) {
	appKindPropName := "APP_KIND"
	appKindPropValue := "workflowApp"

	storageAccount := d.StorageAccountName
	accountKey := d.StorageAccountAccessKey
	storageConnection := fmt.Sprintf(
		"DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s",
		storageAccount,
		accountKey,
		endpointSuffix,
	)
	functionVersion := d.Version

	contentShare := strings.ToLower(d.Name) + "-content"
	if d.StorageAccountShareName != "" {
		contentShare = d.StorageAccountShareName
	}

	basicSettings := []webapps.NameValuePair{
		{Name: &storageAppSettingName, Value: &storageConnection},
		{Name: &functionVersionAppSettingName, Value: &functionVersion},
		{Name: &appKindPropName, Value: &appKindPropValue},
		{Name: &contentShareAppSettingName, Value: &contentShare},
		{Name: &contentFileConnStringAppSettingName, Value: &storageConnection},
	}

	if d.UseExtensionBundle {
		extensionBundlePropName := "AzureFunctionsJobHost__extensionBundle__id"
		extensionBundleName := "Microsoft.Azure.Functions.ExtensionBundle.Workflows"
		extensionBundleVersionPropName := "AzureFunctionsJobHost__extensionBundle__version"
		extensionBundleVersion := d.BundleVersion

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
		Type:     pluginsdk.TypeList,
		Optional: true,
		// Computed:   true,
		// ConfigMode: pluginsdk.SchemaConfigModeAttr,
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

func flattenLogicAppStandardSiteConfig(input *webapps.SiteConfig) []helpers.LogicAppSiteConfig {
	results := make([]helpers.LogicAppSiteConfig, 0)
	result := helpers.LogicAppSiteConfig{}

	if input == nil {
		log.Printf("[DEBUG] SiteConfig is nil")
		return results
	}

	result.AlwaysOn = pointer.From(input.AlwaysOn)
	result.Use32BitWorkerProcess = pointer.From(input.Use32BitWorkerProcess)
	result.WebSocketsEnabled = pointer.From(input.WebSocketsEnabled)
	result.LinuxFxVersion = pointer.From(input.LinuxFxVersion)
	result.HTTP2Enabled = pointer.From(input.HTTP20Enabled)
	result.PreWarmedInstanceCount = pointer.From(input.PreWarmedInstanceCount)
	result.IpRestriction = helpers.FlattenIpRestrictions(input.IPSecurityRestrictions)
	result.SCMIPRestriction = helpers.FlattenIpRestrictions(input.ScmIPSecurityRestrictions)

	result.SCMUseMainIpRestriction = pointer.From(input.ScmIPSecurityRestrictionsUseMain)

	result.SCMType = pointer.FromEnum(input.ScmType)
	result.SCMMinTLSVersion = pointer.FromEnum(input.ScmMinTlsVersion)

	result.MinTLSVersion = pointer.FromEnum(input.MinTlsVersion)
	result.FTPSState = pointer.FromEnum(input.FtpsState)

	result.Cors = helpers.FlattenCorsSettings(input.Cors)

	result.AutoSwapSlotName = pointer.From(input.AutoSwapSlotName)

	result.HealthCheckPath = pointer.From(input.HealthCheckPath)

	result.ElasticInstanceMinimum = pointer.From(input.MinimumElasticInstanceCount)

	result.AppScaleLimit = pointer.From(input.FunctionAppScaleLimit)

	result.RuntimeScaleMonitoringEnabled = pointer.From(input.FunctionsRuntimeScaleMonitoringEnabled)

	result.DotnetFrameworkVersion = pointer.From(input.NetFrameworkVersion)

	result.VNETRouteAllEnabled = pointer.From(input.VnetRouteAllEnabled)

	if !features.FivePointOh() {
		result.PublicNetworkAccessEnabled = strings.EqualFold(pointer.From(input.PublicNetworkAccess), helpers.PublicNetworkAccessEnabled)
	}

	results = append(results, result)
	return results
}

func flattenLogicAppStandardIpRestriction(input *[]webapps.IPSecurityRestriction) []interface{} {
	restrictions := make([]interface{}, 0)

	if input == nil {
		return nil
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

func expandLogicAppStandardSiteConfigForCreate(d []helpers.LogicAppSiteConfig, metadata sdk.ResourceMetaData) (*webapps.SiteConfig, error) {
	siteConfig := &webapps.SiteConfig{}
	if len(d) == 0 {
		return siteConfig, nil
	}

	config := d[0]

	siteConfig.AlwaysOn = pointer.To(config.AlwaysOn)
	siteConfig.HTTP20Enabled = pointer.To(config.HTTP2Enabled)
	siteConfig.FunctionsRuntimeScaleMonitoringEnabled = pointer.To(config.RuntimeScaleMonitoringEnabled)
	siteConfig.Use32BitWorkerProcess = pointer.To(config.Use32BitWorkerProcess)
	siteConfig.WebSocketsEnabled = pointer.To(config.WebSocketsEnabled)

	if config.LinuxFxVersion != "" {
		siteConfig.LinuxFxVersion = pointer.To(config.LinuxFxVersion)
	}

	if len(config.Cors) > 0 {
		siteConfig.Cors = helpers.ExpandCorsSettings(config.Cors)
	}

	ipr, err := helpers.ExpandIpRestrictions(config.IpRestriction)
	if err != nil {
		return nil, err
	}
	siteConfig.IPSecurityRestrictions = ipr

	ipr, err = helpers.ExpandIpRestrictions(config.SCMIPRestriction)
	if err != nil {
		return nil, err
	}
	siteConfig.ScmIPSecurityRestrictions = ipr

	siteConfig.ScmIPSecurityRestrictionsUseMain = pointer.To(config.SCMUseMainIpRestriction)

	siteConfig.ScmMinTlsVersion = pointer.ToEnum[webapps.SupportedTlsVersions](config.SCMMinTLSVersion)
	siteConfig.ScmType = pointer.ToEnum[webapps.ScmType](config.SCMType)
	siteConfig.MinTlsVersion = pointer.ToEnum[webapps.SupportedTlsVersions](config.MinTLSVersion)

	siteConfig.FtpsState = pointer.ToEnum[webapps.FtpsState](config.FTPSState)

	if metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()[0].AsValueMap()["pre_warmed_instance_count"].IsKnown() {
		siteConfig.PreWarmedInstanceCount = pointer.To(config.PreWarmedInstanceCount)
	}
	siteConfig.HealthCheckPath = pointer.To(config.HealthCheckPath)

	if metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()[0].AsValueMap()["elastic_instance_minimum"].IsKnown() && config.ElasticInstanceMinimum > 0 {
		siteConfig.MinimumElasticInstanceCount = pointer.To(config.ElasticInstanceMinimum)
	}

	if metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()[0].AsValueMap()["app_scale_limit"].IsKnown() {
		siteConfig.FunctionAppScaleLimit = pointer.To(config.AppScaleLimit)
	}

	siteConfig.NetFrameworkVersion = pointer.To(config.DotnetFrameworkVersion)
	siteConfig.VnetRouteAllEnabled = pointer.To(config.VNETRouteAllEnabled)

	siteConfig.PublicNetworkAccess = pointer.To(metadata.ResourceData.Get("public_network_access").(string))
	if !features.FivePointOh() {
		siteConfig.PublicNetworkAccess = pointer.To(reconcilePNA(metadata))
	}

	return siteConfig, nil
}

func expandLogicAppStandardSiteConfigForUpdate(d []helpers.LogicAppSiteConfig, metadata sdk.ResourceMetaData, existing *webapps.SiteConfig) (*webapps.SiteConfig, error) {
	if len(d) == 0 {
		return nil, nil
	}

	siteConfig := &webapps.SiteConfig{}
	if existing != nil {
		siteConfig = existing
	}

	config := d[0]

	siteConfig.AlwaysOn = pointer.To(config.AlwaysOn)
	siteConfig.HTTP20Enabled = pointer.To(config.HTTP2Enabled)
	siteConfig.FunctionsRuntimeScaleMonitoringEnabled = pointer.To(config.RuntimeScaleMonitoringEnabled)
	siteConfig.Use32BitWorkerProcess = pointer.To(config.Use32BitWorkerProcess)
	siteConfig.ScmIPSecurityRestrictionsUseMain = pointer.To(config.SCMUseMainIpRestriction)
	siteConfig.WebSocketsEnabled = pointer.To(config.WebSocketsEnabled)
	siteConfig.VnetRouteAllEnabled = pointer.To(config.VNETRouteAllEnabled)

	if metadata.ResourceData.HasChange("site_config.0.linux_fx_version") {
		siteConfig.LinuxFxVersion = pointer.To(config.LinuxFxVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		siteConfig.Cors = helpers.ExpandCorsSettings(config.Cors)
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipr, err := helpers.ExpandIpRestrictions(config.IpRestriction)
		if err != nil {
			return nil, err
		}
		siteConfig.IPSecurityRestrictions = ipr
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		ipr, err := helpers.ExpandIpRestrictions(config.SCMIPRestriction)
		if err != nil {
			return nil, err
		}
		siteConfig.ScmIPSecurityRestrictions = ipr
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_min_tls_version") {
		siteConfig.ScmMinTlsVersion = pointer.ToEnum[webapps.SupportedTlsVersions](config.SCMMinTLSVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_type") {
		siteConfig.ScmType = pointer.ToEnum[webapps.ScmType](config.SCMType)
	}

	if metadata.ResourceData.HasChange("site_config.0.min_tls_version") {
		siteConfig.MinTlsVersion = pointer.ToEnum[webapps.SupportedTlsVersions](config.MinTLSVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		siteConfig.FtpsState = pointer.ToEnum[webapps.FtpsState](config.FTPSState)
	}

	if len(metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()) > 0 && metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()[0].AsValueMap()["pre_warmed_instance_count"].IsKnown() {
		siteConfig.PreWarmedInstanceCount = pointer.To(config.PreWarmedInstanceCount)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		siteConfig.HealthCheckPath = pointer.To(config.HealthCheckPath)
	}

	if len(metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()) > 0 && metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()[0].AsValueMap()["elastic_instance_minimum"].IsKnown() {
		siteConfig.MinimumElasticInstanceCount = pointer.To(config.ElasticInstanceMinimum)
	}

	if len(metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()) > 0 && metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()[0].AsValueMap()["app_scale_limit"].IsKnown() {
		siteConfig.FunctionAppScaleLimit = pointer.To(config.AppScaleLimit)
	}

	if metadata.ResourceData.HasChange("site_config.0.dotnet_framework_version") {
		siteConfig.NetFrameworkVersion = pointer.To(config.DotnetFrameworkVersion)
	}

	if !features.FivePointOh() && metadata.ResourceData.HasChanges("site_config.0.public_network_access_enabled") {
		siteConfig.PublicNetworkAccess = pointer.To(reconcilePNA(metadata))
	}

	if metadata.ResourceData.HasChanges("app_settings", "storage_account_name", "storage_account_share_name", "storage_account_access_key", "version") {
		o, n := metadata.ResourceData.GetChange("app_settings")

		appSettings := make([]webapps.NameValuePair, 0)
		if existing != nil {
			appSettings = *existing.AppSettings
		}

		siteConfig.AppSettings = mergeAppSettings(appSettings, o.(map[string]interface{}), n.(map[string]interface{}), metadata)
	}

	return siteConfig, nil
}

func expandAppSettings(input map[string]string) []webapps.NameValuePair {
	output := make([]webapps.NameValuePair, 0)

	for k, v := range input {
		nameValue := webapps.NameValuePair{
			Name:  pointer.To(k),
			Value: pointer.To(v),
		}
		output = append(output, nameValue)
	}

	return output
}

func mergeAppSettings(existing []webapps.NameValuePair, old, new map[string]interface{}, metadata sdk.ResourceMetaData) *[]webapps.NameValuePair {
	f := func(input map[string]interface{}) (result map[string]string) {
		result = make(map[string]string)
		for k, v := range input {
			result[k] = v.(string)
		}

		return
	}

	eMap := make(map[string]string)
	for _, i := range existing {
		n, v := pointer.From(i.Name), pointer.From(i.Value)
		eMap[n] = v
	}

	oMap := f(old)
	cMap := f(new)

	if metadata.ResourceData.HasChanges("storage_account_name", "storage_account_access_key") {
		accountName := metadata.ResourceData.Get("storage_account_name").(string)
		accountAccessKey := metadata.ResourceData.Get("storage_account_access_key").(string)
		suffix, _ := metadata.Client.Account.Environment.Storage.DomainSuffix()

		eMap[storageAppSettingName] = fmt.Sprintf(storageConnectionStringFmt, accountName, accountAccessKey, *suffix)
		eMap[contentFileConnStringAppSettingName] = fmt.Sprintf(storageConnectionStringFmt, accountName, accountAccessKey, *suffix)
	}

	if metadata.ResourceData.HasChange("storage_account_share_name") {
		n := metadata.ResourceData.Get("storage_account_share_name").(string)

		if n != "" {
			eMap[contentShareAppSettingName] = n
		} else {
			name := metadata.ResourceData.Get("name").(string)
			eMap[contentShareAppSettingName] = strings.ToLower(name) + "-content"
		}
	}

	if metadata.ResourceData.HasChange("version") {
		eMap[functionVersionAppSettingName] = metadata.ResourceData.Get("version").(string)
	}

	if metadata.ResourceData.HasChanges("bundle_version", "use_extension_bundle") {
		if bundleVersion := metadata.ResourceData.Get("bundle_version").(string); bundleVersion != "" {
			eMap[extensionBundleAppSettingName] = extensionBundleAppSettingValue
			eMap[extensionBundleVersionAppSettingName] = bundleVersion
		} else {
			delete(eMap, extensionBundleAppSettingName)
		}
	}

	remove := map[string]string{}
	addOrUpdate := map[string]string{}

	for k, v := range oMap {
		if _, ok := cMap[k]; !ok {
			remove[k] = v
			break
		}

		addOrUpdate[k] = v
	}

	for k, v := range cMap {
		addOrUpdate[k] = v
	}

	for k := range remove {
		delete(eMap, k)
	}

	for k, v := range addOrUpdate {
		eMap[k] = v
	}

	return pointer.To(expandAppSettings(eMap))
}

func reconcilePNA(d sdk.ResourceMetaData) string {
	pna := ""
	scPNASet := false
	d.ResourceData.GetRawConfig().AsValueMap()["public_network_access"].IsNull()
	if !d.ResourceData.GetRawConfig().AsValueMap()["public_network_access"].IsNull() { // is top level set, takes precedence
		pna = d.ResourceData.Get("public_network_access").(string)
	}
	if sc := d.ResourceData.GetRawConfig().AsValueMap()["site_config"]; !sc.IsNull() {
		if len(sc.AsValueSlice()) > 0 && !sc.AsValueSlice()[0].AsValueMap()["public_network_access_enabled"].IsNull() {
			scPNASet = true
		}
	}
	if pna == "" && scPNASet { // if not, or it's empty, is site_config value set
		pnaBool := d.ResourceData.Get("site_config.0.public_network_access_enabled").(bool)
		if pnaBool {
			pna = helpers.PublicNetworkAccessEnabled
		} else {
			pna = helpers.PublicNetworkAccessDisabled
		}
	}

	return pna
}
