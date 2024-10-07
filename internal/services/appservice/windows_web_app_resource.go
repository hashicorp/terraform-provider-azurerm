// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
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
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WindowsWebAppResource struct{}

type WindowsWebAppModel struct {
	Name                             string                                     `tfschema:"name"`
	ResourceGroup                    string                                     `tfschema:"resource_group_name"`
	Location                         string                                     `tfschema:"location"`
	ServicePlanId                    string                                     `tfschema:"service_plan_id"`
	AppSettings                      map[string]string                          `tfschema:"app_settings"`
	StickySettings                   []helpers.StickySettings                   `tfschema:"sticky_settings"`
	AuthSettings                     []helpers.AuthSettings                     `tfschema:"auth_settings"`
	AuthV2Settings                   []helpers.AuthV2Settings                   `tfschema:"auth_settings_v2"`
	Backup                           []helpers.Backup                           `tfschema:"backup"`
	ClientAffinityEnabled            bool                                       `tfschema:"client_affinity_enabled"`
	ClientCertEnabled                bool                                       `tfschema:"client_certificate_enabled"`
	ClientCertMode                   string                                     `tfschema:"client_certificate_mode"`
	ClientCertExclusionPaths         string                                     `tfschema:"client_certificate_exclusion_paths"`
	Enabled                          bool                                       `tfschema:"enabled"`
	HttpsOnly                        bool                                       `tfschema:"https_only"`
	Identity                         []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	KeyVaultReferenceIdentityID      string                                     `tfschema:"key_vault_reference_identity_id"`
	LogsConfig                       []helpers.LogsConfig                       `tfschema:"logs"`
	PublicNetworkAccess              bool                                       `tfschema:"public_network_access_enabled"`
	PublishingDeployBasicAuthEnabled bool                                       `tfschema:"webdeploy_publish_basic_authentication_enabled"`
	PublishingFTPBasicAuthEnabled    bool                                       `tfschema:"ftp_publish_basic_authentication_enabled"`
	SiteConfig                       []helpers.SiteConfigWindows                `tfschema:"site_config"`
	StorageAccounts                  []helpers.StorageAccount                   `tfschema:"storage_account"`
	ConnectionStrings                []helpers.ConnectionString                 `tfschema:"connection_string"`
	CustomDomainVerificationId       string                                     `tfschema:"custom_domain_verification_id"`
	HostingEnvId                     string                                     `tfschema:"hosting_environment_id"`
	DefaultHostname                  string                                     `tfschema:"default_hostname"`
	Kind                             string                                     `tfschema:"kind"`
	OutboundIPAddresses              string                                     `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList            []string                                   `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses      string                                     `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList    []string                                   `tfschema:"possible_outbound_ip_address_list"`
	SiteCredentials                  []helpers.SiteCredential                   `tfschema:"site_credential"`
	ZipDeployFile                    string                                     `tfschema:"zip_deploy_file"`
	Tags                             map[string]string                          `tfschema:"tags"`
	VirtualNetworkSubnetID           string                                     `tfschema:"virtual_network_subnet_id"`
}

var _ sdk.ResourceWithCustomImporter = WindowsWebAppResource{}

var _ sdk.ResourceWithStateMigration = WindowsWebAppResource{}

func (r WindowsWebAppResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"service_plan_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateAppServicePlanID,
		},

		// Optional

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			ValidateFunc: validate.AppSettings,
		},

		"auth_settings": helpers.AuthSettingsSchema(),

		"auth_settings_v2": helpers.AuthV2SettingsSchema(),

		"backup": helpers.BackupSchema(),

		"client_affinity_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"client_certificate_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"client_certificate_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(webapps.ClientCertModeRequired),
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForClientCertMode(), false),
		},

		"client_certificate_exclusion_paths": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Description: "Paths to exclude when using client certificates, separated by ;",
		},

		"connection_string": helpers.ConnectionStringSchema(),

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

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"key_vault_reference_identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"logs": helpers.LogsConfigSchema(),

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"webdeploy_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"ftp_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"site_config": helpers.SiteConfigSchemaWindows(),

		"sticky_settings": helpers.StickySettingsSchema(),

		"storage_account": helpers.StorageAccountSchemaWindows(),

		"zip_deploy_file": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The local path and filename of the Zip packaged application to deploy to this Windows Web App. **Note:** Using this value requires either `WEBSITE_RUN_FROM_PACKAGE=1` or `SCM_DO_BUILD_DURING_DEPLOYMENT=true` to be set on the App in `app_settings`.",
		},

		"tags": tags.Schema(),

		"virtual_network_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},
	}
}

func (r WindowsWebAppResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"custom_domain_verification_id": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"default_hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hosting_environment_id": {
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

		"outbound_ip_address_list": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"possible_outbound_ip_addresses": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"possible_outbound_ip_address_list": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"site_credential": helpers.SiteCredentialSchema(),
	}
}

func (r WindowsWebAppResource) ModelObject() interface{} {
	return &WindowsWebAppModel{}
}

func (r WindowsWebAppResource) ResourceType() string {
	return "azurerm_windows_web_app"
}

func (r WindowsWebAppResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var webApp WindowsWebAppModel
			if err := metadata.Decode(&webApp); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			resourceProvidersClient := metadata.Client.AppService.ResourceProvidersClient
			servicePlanClient := metadata.Client.AppService.ServicePlanClient
			aseClient := metadata.Client.AppService.AppServiceEnvironmentClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			baseId := commonids.NewAppServiceID(subscriptionId, webApp.ResourceGroup, webApp.Name)
			id, err := commonids.ParseWebAppID(baseId.ID())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Windows %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			sc := webApp.SiteConfig[0]

			availabilityRequest := resourceproviders.ResourceNameAvailabilityRequest{
				Name: (webApp.Name),
				Type: resourceproviders.CheckNameResourceTypesMicrosoftPointWebSites,
			}

			servicePlanId, err := commonids.ParseAppServicePlanID(webApp.ServicePlanId)
			if err != nil {
				return err
			}

			servicePlan, err := servicePlanClient.Get(ctx, *servicePlanId)
			if err != nil {
				return fmt.Errorf("reading App %s: %+v", servicePlanId, err)
			}
			if servicePlan.Model != nil && servicePlan.Model.Properties != nil {
				if ase := servicePlan.Model.Properties.HostingEnvironmentProfile; ase != nil {
					// Attempt to check the ASE for the appropriate suffix for the name availability request. Not convinced
					// the `DNSSuffix` field is still valid and possibly should have been deprecated / removed as is legacy
					// setting from ASEv1? Hence the non-fatal approach here.
					nameSuffix := "appserviceenvironment.net"
					if ase.Id != nil {
						aseId, err := commonids.ParseAppServiceEnvironmentIDInsensitively(*ase.Id)
						nameSuffix = fmt.Sprintf("%s.%s", aseId.HostingEnvironmentName, nameSuffix)
						if err != nil {
							metadata.Logger.Warnf("could not parse App Service Environment ID determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", webApp.Name, servicePlanId)
						} else {
							existingASE, err := aseClient.Get(ctx, *aseId)
							if err != nil || existingASE.Model == nil {
								metadata.Logger.Warnf("could not read App Service Environment to determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", webApp.Name, servicePlanId)
							} else if props := existingASE.Model.Properties; props != nil && props.DnsSuffix != nil && *props.DnsSuffix != "" {
								nameSuffix = *props.DnsSuffix
							}
						}
					}
					availabilityRequest.Name = fmt.Sprintf("%s.%s", webApp.Name, nameSuffix)
					availabilityRequest.IsFqdn = pointer.To(true)
				}
				if servicePlan.Model.Sku != nil && servicePlan.Model.Sku.Name != nil {
					if helpers.IsFreeOrSharedServicePlan(*servicePlan.Model.Sku.Name) {
						if sc.AlwaysOn {
							return fmt.Errorf("always_on cannot be set to true when using Free, F1, D1 Sku")
						}
					}
				}
			}

			subscriptionID := commonids.NewSubscriptionID(subscriptionId)
			checkName, err := resourceProvidersClient.CheckNameAvailability(ctx, subscriptionID, availabilityRequest)
			if err != nil || checkName.Model == nil {
				return fmt.Errorf("checking name availability for %s: %+v", id, err)
			}
			if !*checkName.Model.NameAvailable {
				return fmt.Errorf("the Site Name %q failed the availability check: %+v", id.SiteName, *checkName.Model.Message)
			}

			siteConfig, err := sc.ExpandForCreate(webApp.AppSettings)
			if err != nil {
				return err
			}

			currentStack := ""
			if len(sc.ApplicationStack) == 1 {
				currentStack = sc.ApplicationStack[0].CurrentStack
			}

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(webApp.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			siteEnvelope := webapps.Site{
				Location: location.Normalize(webApp.Location),
				Tags:     pointer.To(webApp.Tags),
				Identity: expandedIdentity,
				Properties: &webapps.SiteProperties{
					ServerFarmId:          pointer.To(webApp.ServicePlanId),
					Enabled:               pointer.To(webApp.Enabled),
					HTTPSOnly:             pointer.To(webApp.HttpsOnly),
					SiteConfig:            siteConfig,
					ClientAffinityEnabled: pointer.To(webApp.ClientAffinityEnabled),
					ClientCertEnabled:     pointer.To(webApp.ClientCertEnabled),
					ClientCertMode:        pointer.To(webapps.ClientCertMode(webApp.ClientCertMode)),
					VnetRouteAllEnabled:   siteConfig.VnetRouteAllEnabled,
				},
			}

			pna := helpers.PublicNetworkAccessEnabled
			if !webApp.PublicNetworkAccess {
				pna = helpers.PublicNetworkAccessDisabled
			}

			// (@jackofallops) - Values appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
			siteEnvelope.Properties.PublicNetworkAccess = pointer.To(pna)
			siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = siteEnvelope.Properties.PublicNetworkAccess

			if webApp.KeyVaultReferenceIdentityID != "" {
				siteEnvelope.Properties.KeyVaultReferenceIdentity = pointer.To(webApp.KeyVaultReferenceIdentityID)
			}

			if webApp.VirtualNetworkSubnetID != "" {
				siteEnvelope.Properties.VirtualNetworkSubnetId = pointer.To(webApp.VirtualNetworkSubnetID)
			}

			if webApp.ClientCertExclusionPaths != "" {
				siteEnvelope.Properties.ClientCertExclusionPaths = pointer.To(webApp.ClientCertExclusionPaths)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, siteEnvelope); err != nil {
				return fmt.Errorf("creating Windows %s: %+v", id, err)
			}

			metadata.SetID(id)

			if currentStack != "" {
				siteMetadata := webapps.StringDictionary{Properties: &map[string]string{
					"CURRENT_STACK": currentStack,
				}}
				if _, err := client.UpdateMetadata(ctx, *id, siteMetadata); err != nil {
					return fmt.Errorf("setting Site Metadata for Current Stack on Windows %s: %+v", id, err)
				}
			}

			appSettings := helpers.ExpandAppSettingsForUpdate(siteConfig.AppSettings)
			appSettingsProps := *appSettings.Properties
			if metadata.ResourceData.HasChange("site_config.0.health_check_eviction_time_in_min") {
				appSettingsProps["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"] = strconv.FormatInt(webApp.SiteConfig[0].HealthCheckEvictionTime, 10)
				appSettings.Properties = &appSettingsProps
			}

			if appSettings != nil {
				if _, err := client.UpdateApplicationSettings(ctx, *id, *appSettings); err != nil {
					return fmt.Errorf("setting App Settings for Windows %s: %+v", id, err)
				}
			}

			stickySettings := helpers.ExpandStickySettings(webApp.StickySettings)

			if stickySettings != nil {
				stickySettingsUpdate := webapps.SlotConfigNamesResource{
					Properties: stickySettings,
				}
				if _, err := client.UpdateSlotConfigurationNames(ctx, *id, stickySettingsUpdate); err != nil {
					return fmt.Errorf("updating Sticky Settings for Windows %s: %+v", *id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(webApp.AuthSettings)
			if auth.Properties != nil {
				if _, err := client.UpdateAuthSettings(ctx, *id, *auth); err != nil {
					return fmt.Errorf("setting Authorisation Settings for %s: %+v", *id, err)
				}
			}

			authv2 := helpers.ExpandAuthV2Settings(webApp.AuthV2Settings)
			if authv2.Properties != nil {
				if _, err = client.UpdateAuthSettingsV2(ctx, *id, *authv2); err != nil {
					return fmt.Errorf("updating AuthV2 settings for Windows %s: %+v", *id, err)
				}
			}

			if metadata.ResourceData.HasChange("logs") {
				logsConfig := helpers.ExpandLogsConfig(webApp.LogsConfig)
				if logsConfig.Properties != nil {
					if _, err := client.UpdateDiagnosticLogsConfig(ctx, *id, *logsConfig); err != nil {
						return fmt.Errorf("setting Diagnostic Logs Configuration for Windows %s: %+v", id, err)
					}
				}
			}

			backupConfig, err := helpers.ExpandBackupConfig(webApp.Backup)
			if err != nil {
				return fmt.Errorf("expanding backup configuration for Windows %s: %+v", *id, err)
			}

			if backupConfig.Properties != nil {
				if _, err := client.UpdateBackupConfiguration(ctx, *id, *backupConfig); err != nil {
					return fmt.Errorf("adding Backup Settings for Windows %s: %+v", *id, err)
				}
			}

			storageConfig := helpers.ExpandStorageConfig(webApp.StorageAccounts)
			if storageConfig.Properties != nil {
				if _, err := client.UpdateAzureStorageAccounts(ctx, *id, *storageConfig); err != nil {
					if err != nil {
						return fmt.Errorf("setting Storage Accounts for Windows %s: %+v", *id, err)
					}
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(webApp.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStrings(ctx, *id, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for Windows %s: %+v", id, err)
				}
			}

			if webApp.ZipDeployFile != "" {
				if err = helpers.GetCredentialsAndPublish(ctx, client, *id, webApp.ZipDeployFile); err != nil {
					return err
				}
			}

			if !webApp.PublishingDeployBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: false,
					},
				}
				if _, err := client.UpdateScmAllowed(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", *id, err)
				}
			}

			if !webApp.PublishingFTPBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: false,
					},
				}
				if _, err := client.UpdateFtpAllowed(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for ftp publishing credentials for %s: %+v", *id, err)
				}
			}

			return nil
		},

		Timeout: 30 * time.Minute,
	}
}

func (r WindowsWebAppResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := commonids.ParseWebAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			webApp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(webApp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Windows %s: %+v", id, err)
			}

			// Despite being part of the defined `Get` response model, site_config is always nil so we get it explicitly
			webAppSiteConfig, err := client.GetConfiguration(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Site Config for Windows %s: %+v", id, err)
			}

			auth, err := client.GetAuthSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Windows %s: %+v", id, err)
			}

			var authV2 webapps.SiteAuthSettingsV2
			if strings.EqualFold(pointer.From(auth.Model.Properties.ConfigVersion), "v2") {
				authV2Resp, err := client.GetAuthSettingsV2(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading authV2 settings for Linux %s: %+v", id, err)
				}
				authV2 = *authV2Resp.Model
			}

			backup, err := client.GetBackupConfiguration(ctx, *id)
			if err != nil {
				if !response.WasNotFound(backup.HttpResponse) {
					return fmt.Errorf("reading Backup Settings for Windows %s: %+v", *id, err)
				}
			}

			logsConfig, err := client.GetDiagnosticLogsConfiguration(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Diagnostic Logs information for Windows %s: %+v", *id, err)
			}

			appSettings, err := client.ListApplicationSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading App Settings for Windows %s: %+v", *id, err)
			}

			stickySettings, err := client.ListSlotConfigurationNames(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Sticky Settings for Linux %s: %+v", *id, err)
			}

			storageAccounts, err := client.ListAzureStorageAccounts(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Storage Account information for Windows %s: %+v", *id, err)
			}

			connectionStrings, err := client.ListConnectionStrings(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Windows %s: %+v", id, err)
			}

			siteCredentials, err := helpers.ListPublishingCredentials(ctx, client, *id)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for %s: %+v", id, err)
			}

			siteMetadata, err := client.ListMetadata(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Site Metadata for Windows %s: %+v", id, err)
			}

			basicAuthFTP := true
			if basicAuthFTPResp, err := client.GetFtpAllowed(ctx, *id); err != nil || basicAuthFTPResp.Model == nil {
				return fmt.Errorf("retrieving state of FTP Basic Auth for %s: %+v", id, err)
			} else if csmProps := basicAuthFTPResp.Model.Properties; csmProps != nil {
				basicAuthFTP = csmProps.Allow
			}

			basicAuthWebDeploy := true
			if basicAuthWebDeployResp, err := client.GetScmAllowed(ctx, *id); err != nil || basicAuthWebDeployResp.Model == nil {
				return fmt.Errorf("retrieving state of WebDeploy Basic Auth for %s: %+v", id, err)
			} else if csmProps := basicAuthWebDeployResp.Model.Properties; csmProps != nil {
				basicAuthWebDeploy = csmProps.Allow
			}

			state := WindowsWebAppModel{
				Name:              id.SiteName,
				ResourceGroup:     id.ResourceGroupName,
				AuthSettings:      helpers.FlattenAuthSettings(auth.Model),
				AuthV2Settings:    helpers.FlattenAuthV2Settings(authV2),
				Backup:            helpers.FlattenBackupConfig(backup.Model),
				ConnectionStrings: helpers.FlattenConnectionStrings(connectionStrings.Model),
				LogsConfig:        helpers.FlattenLogsConfig(logsConfig.Model),
				StickySettings:    helpers.FlattenStickySettings(stickySettings.Model.Properties),
				SiteCredentials:   helpers.FlattenSiteCredentials(siteCredentials),
				StorageAccounts:   helpers.FlattenStorageAccounts(storageAccounts.Model),
			}

			if model := webApp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Kind = pointer.From(model.Kind)
				state.Tags = pointer.From(model.Tags)
				if props := model.Properties; props != nil {
					state.ClientAffinityEnabled = pointer.From(props.ClientAffinityEnabled)
					state.ClientCertEnabled = pointer.From(props.ClientCertEnabled)
					state.ClientCertMode = string(pointer.From(props.ClientCertMode))
					state.ClientCertExclusionPaths = pointer.From(props.ClientCertExclusionPaths)
					state.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
					state.DefaultHostname = pointer.From(props.DefaultHostName)
					state.Enabled = pointer.From(props.Enabled)
					state.HttpsOnly = pointer.From(props.HTTPSOnly)
					state.KeyVaultReferenceIdentityID = pointer.From(props.KeyVaultReferenceIdentity)
					state.OutboundIPAddresses = pointer.From(props.OutboundIPAddresses)
					state.OutboundIPAddressList = strings.Split(pointer.From(props.OutboundIPAddresses), ",")
					state.PossibleOutboundIPAddresses = pointer.From(props.PossibleOutboundIPAddresses)
					state.PossibleOutboundIPAddressList = strings.Split(pointer.From(props.PossibleOutboundIPAddresses), ",")
					state.PublicNetworkAccess = !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled)

					serverFarmId, err := commonids.ParseAppServicePlanIDInsensitively(pointer.From(props.ServerFarmId))
					if err != nil {
						return fmt.Errorf("parsing Service Plan ID for %s: %+v", id, err)
					}

					state.ServicePlanId = serverFarmId.ID()

					if hostingEnv := props.HostingEnvironmentProfile; hostingEnv != nil {
						hostingEnvId, err := parse.AppServiceEnvironmentIDInsensitively(pointer.From(hostingEnv.Id))
						if err != nil {
							return err
						}
						state.HostingEnvId = hostingEnvId.ID()
					}

					if subnetId := pointer.From(props.VirtualNetworkSubnetId); subnetId != "" {
						// some users have provisioned these without a prefixed `/` - as such we need to normalize these
						parsed, err := commonids.ParseSubnetIDInsensitively(subnetId)
						if err != nil {
							return err
						}
						state.VirtualNetworkSubnetID = parsed.ID()
					}

					state.PublishingFTPBasicAuthEnabled = basicAuthFTP
					state.PublishingDeployBasicAuthEnabled = basicAuthWebDeploy

					state.AppSettings = helpers.FlattenWebStringDictionary(appSettings.Model)

					currentStack := ""
					if m := siteMetadata.Model; m != nil && m.Properties != nil {
						p := *m.Properties
						if v, ok := p["CURRENT_STACK"]; ok {
							currentStack = v
						}
					}

					siteConfig := helpers.SiteConfigWindows{}
					if err := siteConfig.Flatten(webAppSiteConfig.Model.Properties, currentStack); err != nil {
						return err
					}
					siteConfig.SetHealthCheckEvictionTime(state.AppSettings)
					state.AppSettings = siteConfig.ParseNodeVersion(state.AppSettings)

					// For non-import cases we check for use of the deprecated docker settings - remove in 4.0
					_, usesDeprecatedDocker := metadata.ResourceData.GetOk("site_config.0.application_stack.0.docker_container_name")

					if helpers.FxStringHasPrefix(siteConfig.WindowsFxVersion, helpers.FxStringPrefixDocker) {
						if !features.FourPointOhBeta() {
							siteConfig.DecodeDockerDeprecatedAppStack(state.AppSettings, usesDeprecatedDocker)
						} else {
							siteConfig.DecodeDockerAppStack(state.AppSettings)
						}
					}

					state.SiteConfig = []helpers.SiteConfigWindows{siteConfig}

					// Filter out all settings we've consumed above
					if !features.FourPointOhBeta() && usesDeprecatedDocker {
						state.AppSettings = helpers.FilterManagedAppSettingsDeprecated(state.AppSettings)
					} else {
						state.AppSettings = helpers.FilterManagedAppSettings(state.AppSettings)
					}

					// Zip Deploys are not retrievable, so attempt to get from config. This doesn't matter for imports as an unexpected value here could break the deployment.
					if deployFile, ok := metadata.ResourceData.Get("zip_deploy_file").(string); ok {
						state.ZipDeployFile = deployFile
					}

					flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
					if err != nil {
						return fmt.Errorf("flattening `identity`: %+v", err)
					}

					state.Identity = pointer.From(flattenedIdentity)

					if err := metadata.Encode(&state); err != nil {
						return fmt.Errorf("encoding: %+v", err)
					}

				}
			}

			return nil
		},
	}
}

func (r WindowsWebAppResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := commonids.ParseWebAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			delOptions := webapps.DeleteOperationOptions{
				DeleteEmptyServerFarm: pointer.To(false),
				DeleteMetrics:         pointer.To(false),
			}
			if _, err = client.Delete(ctx, *id, delOptions); err != nil {
				return fmt.Errorf("deleting Windows %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r WindowsWebAppResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateAppServiceID
}

func (r WindowsWebAppResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			servicePlanClient := metadata.Client.AppService.ServicePlanClient

			id, err := commonids.ParseWebAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state WindowsWebAppModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("reading Windows %s: %v", *id, err)
			}

			model := *existing.Model

			// Despite being part of the defined `Get` response model, site_config is always nil so we get it explicitly
			webAppSiteConfig, err := client.GetConfiguration(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Site Config for Windows %s: %+v", id, err)
			}
			model.Properties.SiteConfig = webAppSiteConfig.Model.Properties

			var serviceFarmId string
			servicePlanChange := false
			if model.Properties.ServerFarmId != nil {
				serviceFarmId = *model.Properties.ServerFarmId
			}

			if metadata.ResourceData.HasChange("service_plan_id") {
				serviceFarmId = state.ServicePlanId
				model.Properties.ServerFarmId = pointer.To(serviceFarmId)
				servicePlanChange = true
			}
			servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(serviceFarmId)
			if err != nil {
				return err
			}

			servicePlan, err := servicePlanClient.Get(ctx, *servicePlanId)
			if err != nil {
				return fmt.Errorf("reading App %s: %+v", servicePlanId, err)
			}

			if metadata.ResourceData.HasChange("enabled") {
				model.Properties.Enabled = pointer.To(state.Enabled)
			}
			if metadata.ResourceData.HasChange("https_only") {
				model.Properties.HTTPSOnly = pointer.To(state.HttpsOnly)
			}
			if metadata.ResourceData.HasChange("client_affinity_enabled") {
				model.Properties.ClientAffinityEnabled = pointer.To(state.ClientAffinityEnabled)
			}
			if metadata.ResourceData.HasChange("client_certificate_enabled") {
				model.Properties.ClientCertEnabled = pointer.To(state.ClientCertEnabled)
			}
			if metadata.ResourceData.HasChange("client_certificate_mode") {
				model.Properties.ClientCertMode = pointer.To(webapps.ClientCertMode(state.ClientCertMode))
			}
			if metadata.ResourceData.HasChange("client_certificate_exclusion_paths") {
				model.Properties.ClientCertExclusionPaths = pointer.To(state.ClientCertExclusionPaths)
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(state.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				model.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("key_vault_reference_identity_id") {
				model.Properties.KeyVaultReferenceIdentity = pointer.To(state.KeyVaultReferenceIdentityID)
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(state.Tags)
			}

			if metadata.ResourceData.HasChange("virtual_network_subnet_id") {
				subnetId := metadata.ResourceData.Get("virtual_network_subnet_id").(string)
				if subnetId == "" {
					if _, err := client.DeleteSwiftVirtualNetwork(ctx, *id); err != nil {
						return fmt.Errorf("removing `virtual_network_subnet_id` association for %s: %+v", *id, err)
					}
					var empty *string
					model.Properties.VirtualNetworkSubnetId = empty
				} else {
					model.Properties.VirtualNetworkSubnetId = pointer.To(subnetId)
				}
			}

			currentStack := ""
			sc := state.SiteConfig[0]

			if servicePlan.Model != nil && servicePlan.Model.Sku != nil && servicePlan.Model.Sku.Name != nil {
				if helpers.IsFreeOrSharedServicePlan(*servicePlan.Model.Sku.Name) {
					if sc.AlwaysOn {
						return fmt.Errorf("always_on feature has to be turned off before switching to a free/shared Sku")
					}
				}
			}

			if len(sc.ApplicationStack) == 1 {
				currentStack = sc.ApplicationStack[0].CurrentStack
			}

			if metadata.ResourceData.HasChanges("site_config", "app_settings") || servicePlanChange {
				model.Properties.SiteConfig, err = sc.ExpandForUpdate(metadata, model.Properties.SiteConfig, state.AppSettings)
				if err != nil {
					return err
				}
				model.Properties.VnetRouteAllEnabled = existing.Model.Properties.SiteConfig.VnetRouteAllEnabled
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				pna := helpers.PublicNetworkAccessEnabled
				if !state.PublicNetworkAccess {
					pna = helpers.PublicNetworkAccessDisabled
				}

				// (@jackofallops) - Values appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
				model.Properties.PublicNetworkAccess = pointer.To(pna)
				model.Properties.SiteConfig.PublicNetworkAccess = model.Properties.PublicNetworkAccess
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating Windows %s: %+v", id, err)
			}

			siteMetadata := webapps.StringDictionary{Properties: &map[string]string{
				"CURRENT_STACK": currentStack,
			}}
			if _, err := client.UpdateMetadata(ctx, *id, siteMetadata); err != nil {
				return fmt.Errorf("setting Site Metadata for Current Stack on Windows %s: %+v", id, err)
			}

			updateLogs := false

			// sending App Settings updates can clobber logs configuration so must be updated before we send any Log updates
			if metadata.ResourceData.HasChanges("app_settings", "site_config") {
				appSettingsUpdate := helpers.ExpandAppSettingsForUpdate(model.Properties.SiteConfig.AppSettings)
				appSettingsProps := *appSettingsUpdate.Properties
				if state.SiteConfig[0].HealthCheckEvictionTime != 0 {
					appSettingsProps["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"] = strconv.FormatInt(state.SiteConfig[0].HealthCheckEvictionTime, 10)
					appSettingsUpdate.Properties = &appSettingsProps
				} else {
					delete(appSettingsProps, "WEBSITE_HEALTHCHECK_MAXPINGFAILURES")
					appSettingsUpdate.Properties = &appSettingsProps
				}
				if _, err := client.UpdateApplicationSettings(ctx, *id, *appSettingsUpdate); err != nil {
					return fmt.Errorf("updating App Settings for Linux %s: %+v", *id, err)
				}

				updateLogs = true
			}

			if metadata.ResourceData.HasChange("connection_string") {
				connectionStringUpdate := helpers.ExpandConnectionStrings(state.ConnectionStrings)
				if connectionStringUpdate.Properties == nil {
					connectionStringUpdate.Properties = &map[string]webapps.ConnStringValueTypePair{}
				}
				if _, err := client.UpdateConnectionStrings(ctx, *id, *connectionStringUpdate); err != nil {
					return fmt.Errorf("updating Connection Strings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("sticky_settings") {
				emptySlice := make([]string, 0)
				stickySettings := helpers.ExpandStickySettings(state.StickySettings)
				stickySettingsUpdate := webapps.SlotConfigNamesResource{
					Properties: &webapps.SlotConfigNames{
						AppSettingNames:       &emptySlice,
						ConnectionStringNames: &emptySlice,
					},
				}

				if stickySettings != nil {
					if stickySettings.AppSettingNames != nil {
						stickySettingsUpdate.Properties.AppSettingNames = stickySettings.AppSettingNames
					}
					if stickySettings.ConnectionStringNames != nil {
						stickySettingsUpdate.Properties.ConnectionStringNames = stickySettings.ConnectionStringNames
					}
				}

				if _, err := client.UpdateSlotConfigurationNames(ctx, *id, stickySettingsUpdate); err != nil {
					return fmt.Errorf("updating Sticky Settings for Linux %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings") {
				authUpdate := helpers.ExpandAuthSettings(state.AuthSettings)
				if authUpdate.Properties == nil {
					authUpdate.Properties = &webapps.SiteAuthSettingsProperties{
						Enabled:                           pointer.To(false),
						ClientSecret:                      pointer.To(""),
						ClientSecretSettingName:           pointer.To(""),
						ClientSecretCertificateThumbprint: pointer.To(""),
						GoogleClientSecret:                pointer.To(""),
						FacebookAppSecret:                 pointer.To(""),
						GitHubClientSecret:                pointer.To(""),
						TwitterConsumerSecret:             pointer.To(""),
						MicrosoftAccountClientSecret:      pointer.To(""),
					}
					updateLogs = true
				}
				if _, err := client.UpdateAuthSettings(ctx, *id, *authUpdate); err != nil {
					return fmt.Errorf("updating Auth Settings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings_v2") {
				authV2Update := helpers.ExpandAuthV2Settings(state.AuthV2Settings)
				if _, err := client.UpdateAuthSettingsV2(ctx, *id, *authV2Update); err != nil {
					return fmt.Errorf("updating AuthV2 Settings for Linux %s: %+v", id, err)
				}
				updateLogs = true
			}

			if metadata.ResourceData.HasChange("backup") {
				backupUpdate, err := helpers.ExpandBackupConfig(state.Backup)
				if err != nil {
					return fmt.Errorf("expanding backup configuration for Windows %s: %+v", *id, err)
				}

				if backupUpdate.Properties == nil {
					if _, err = client.DeleteBackupConfiguration(ctx, *id); err != nil {
						return fmt.Errorf("removing Backup Settings for Windows %s: %+v", id, err)
					}
				} else {
					if _, err = client.UpdateBackupConfiguration(ctx, *id, *backupUpdate); err != nil {
						return fmt.Errorf("updating Backup Settings for Windows %s: %+v", id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("logs") || updateLogs {
				logsUpdate := helpers.ExpandLogsConfig(state.LogsConfig)
				if logsUpdate.Properties == nil {
					logsUpdate = helpers.DisabledLogsConfig() // The API is update only, so we need to send an update with everything switched of when a user removes the "logs" block
				}
				if _, err = client.UpdateDiagnosticLogsConfig(ctx, *id, *logsUpdate); err != nil {
					return fmt.Errorf("updating Logs Config for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("storage_account") {
				storageAccountUpdate := helpers.ExpandStorageConfig(state.StorageAccounts)
				if _, err := client.UpdateAzureStorageAccounts(ctx, *id, *storageAccountUpdate); err != nil {
					return fmt.Errorf("updating Storage Accounts for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("zip_deploy_file") {
				if err = helpers.GetCredentialsAndPublish(ctx, client, *id, state.ZipDeployFile); err != nil {
					return err
				}
			}

			if metadata.ResourceData.HasChange("ftp_publish_basic_authentication_enabled") {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: state.PublishingFTPBasicAuthEnabled,
					},
				}
				if _, err := client.UpdateFtpAllowed(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for ftp publishing credentials for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("webdeploy_publish_basic_authentication_enabled") {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: state.PublishingDeployBasicAuthEnabled,
					},
				}
				if _, err := client.UpdateScmAllowed(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r WindowsWebAppResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.AppService.WebAppsClient
		servicePlanClient := metadata.Client.AppService.ServicePlanClient

		id, err := commonids.ParseWebAppID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}
		site, err := client.Get(ctx, *id)
		if err != nil || site.Model == nil || site.Model.Properties == nil {
			return fmt.Errorf("reading Windows %s: %+v", id, err)
		}
		props := site.Model.Properties
		if props.ServerFarmId == nil {
			return fmt.Errorf("determining Service Plan ID for Windows %s: %+v", id, err)
		}
		servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(*props.ServerFarmId)
		if err != nil {
			return err
		}

		sp, err := servicePlanClient.Get(ctx, *servicePlanId)
		if err != nil || sp.Model == nil || sp.Model.Kind == nil {
			return fmt.Errorf("reading Service Plan for Windows %s: %+v", id, err)
		}
		if strings.Contains(strings.ToLower(*sp.Model.Kind), "linux") {
			return fmt.Errorf("specified Service Plan is not a Windows plan")
		}

		return nil
	}
}

func (r WindowsWebAppResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.WindowsWebAppV0toV1{},
		},
	}
}
