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
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WindowsWebAppSlotResource struct{}

type WindowsWebAppSlotModel struct {
	Name                             string                                     `tfschema:"name"`
	AppServiceId                     string                                     `tfschema:"app_service_id"`
	ServicePlanID                    string                                     `tfschema:"service_plan_id"`
	AppSettings                      map[string]string                          `tfschema:"app_settings"`
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
	SiteConfig                       []helpers.SiteConfigWindowsWebAppSlot      `tfschema:"site_config"`
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

var _ sdk.ResourceWithUpdate = WindowsWebAppSlotResource{}

var _ sdk.ResourceWithStateMigration = WindowsWebAppSlotResource{}

func (r WindowsWebAppSlotResource) ModelObject() interface{} {
	return &WindowsWebAppSlotModel{}
}

func (r WindowsWebAppSlotResource) ResourceType() string {
	return "azurerm_windows_web_app_slot"
}

func (r WindowsWebAppSlotResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webapps.ValidateSlotID
}

func (r WindowsWebAppSlotResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
		},

		"app_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateWebAppID,
		},

		// Optional

		"service_plan_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateAppServicePlanID,
		},

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

		"site_config": helpers.SiteConfigSchemaWindowsWebAppSlot(),

		"storage_account": helpers.StorageAccountSchemaWindows(),

		"zip_deploy_file": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The local path and filename of the Zip packaged application to deploy to this Windows Web App. **Note:** Using this value requires `WEBSITE_RUN_FROM_PACKAGE=1` on the App in `app_settings`.",
		},

		"tags": tags.Schema(),

		"virtual_network_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},
	}
}

func (r WindowsWebAppSlotResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r WindowsWebAppSlotResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var webAppSlot WindowsWebAppSlotModel
			if err := metadata.Decode(&webAppSlot); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			appId, err := commonids.ParseWebAppID(webAppSlot.AppServiceId)
			if err != nil {
				return err
			}

			id := webapps.NewSlotID(appId.SubscriptionId, appId.ResourceGroupName, appId.SiteName, webAppSlot.Name)

			webApp, err := client.Get(ctx, *appId)

			if err != nil {
				return fmt.Errorf("reading parent Windows Web App for %s: %+v", id, err)
			}

			var servicePlanId *commonids.AppServicePlanId
			differentServicePlanToParent := false
			if webApp.Model == nil || webApp.Model.Properties == nil || webApp.Model.Properties.ServerFarmId == nil {
				return fmt.Errorf("could not determine Service Plan ID for %s: %+v", id, err)
			}

			servicePlanId, err = commonids.ParseAppServicePlanIDInsensitively(*webApp.Model.Properties.ServerFarmId)
			if err != nil {
				return err
			}

			if webAppSlot.ServicePlanID != "" {
				newServicePlanId, err := commonids.ParseAppServicePlanID(webAppSlot.ServicePlanID)
				if err != nil {
					return fmt.Errorf("parsing service_plan_id for %s: %+v", id, err)
				}
				// we only set `service_plan_id` when it differs from the parent `service_plan_id` which is causing issues
				// https://github.com/hashicorp/terraform-provider-azurerm/issues/21024
				// we'll error here if the `service_plan_id` equals the parent `service_plan_id`
				if strings.EqualFold(newServicePlanId.ID(), servicePlanId.ID()) {
					return fmt.Errorf("`service_plan_id` should only be specified when it differs from the `service_plan_id` of the associated Web App")
				}
				differentServicePlanToParent = true
				servicePlanId = newServicePlanId
			}

			existing, err := client.GetSlot(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Windows %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			sc := webAppSlot.SiteConfig[0]
			siteConfig, err := sc.ExpandForCreate(webAppSlot.AppSettings)
			if err != nil {
				return err
			}

			currentStack := ""
			if len(sc.ApplicationStack) == 1 {
				currentStack = sc.ApplicationStack[0].CurrentStack
			}

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(webAppSlot.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			siteEnvelope := webapps.Site{
				Location: location.Normalize(webApp.Model.Location),
				Tags:     pointer.To(webAppSlot.Tags),
				Identity: expandedIdentity,
				Properties: &webapps.SiteProperties{
					Enabled:                  pointer.To(webAppSlot.Enabled),
					HTTPSOnly:                pointer.To(webAppSlot.HttpsOnly),
					SiteConfig:               siteConfig,
					ClientAffinityEnabled:    pointer.To(webAppSlot.ClientAffinityEnabled),
					ClientCertEnabled:        pointer.To(webAppSlot.ClientCertEnabled),
					ClientCertMode:           pointer.To(webapps.ClientCertMode(webAppSlot.ClientCertMode)),
					ClientCertExclusionPaths: pointer.To(webAppSlot.ClientCertExclusionPaths),
					VnetRouteAllEnabled:      siteConfig.VnetRouteAllEnabled,
				},
			}

			if differentServicePlanToParent {
				siteEnvelope.Properties.ServerFarmId = pointer.To(servicePlanId.ID())
			}

			pna := helpers.PublicNetworkAccessEnabled
			if !webAppSlot.PublicNetworkAccess {
				pna = helpers.PublicNetworkAccessDisabled
			}

			// (@jackofallops) - Values appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
			siteEnvelope.Properties.PublicNetworkAccess = pointer.To(pna)
			siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = siteEnvelope.Properties.PublicNetworkAccess

			if webAppSlot.KeyVaultReferenceIdentityID != "" {
				siteEnvelope.Properties.KeyVaultReferenceIdentity = pointer.To(webAppSlot.KeyVaultReferenceIdentityID)
			}

			if webAppSlot.VirtualNetworkSubnetID != "" {
				siteEnvelope.Properties.VirtualNetworkSubnetId = pointer.To(webAppSlot.VirtualNetworkSubnetID)
				siteEnvelope.Properties.ServerFarmId = pointer.To(servicePlanId.ID())
			}

			if err := client.CreateOrUpdateSlotThenPoll(ctx, id, siteEnvelope); err != nil {
				return fmt.Errorf("creating Windows %s: %+v", id, err)
			}

			// (@jackofallops) - Windows Web App Slots need the siteConfig sending individually to actually accept the `windowsFxVersion` value or it's set as `DOCKER|` only.
			siteConfigUpdate := webapps.SiteConfigResource{
				Properties: siteConfig,
			}
			_, err = client.UpdateConfigurationSlot(ctx, id, siteConfigUpdate)
			if err != nil {
				return fmt.Errorf("updating %s site config: %+v", id, err)
			}

			metadata.SetID(id)

			if currentStack != "" {
				siteMetadata := webapps.StringDictionary{Properties: &map[string]string{
					"CURRENT_STACK": currentStack,
				}}
				if _, err := client.UpdateMetadataSlot(ctx, id, siteMetadata); err != nil {
					return fmt.Errorf("setting Site Metadata for Current Stack on Windows %s: %+v", id, err)
				}
			}

			appSettings := helpers.ExpandAppSettingsForUpdate(siteConfig.AppSettings)
			appSettingsProps := *appSettings.Properties
			if metadata.ResourceData.HasChange("site_config.0.health_check_eviction_time_in_min") {
				appSettingsProps["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"] = strconv.FormatInt(webAppSlot.SiteConfig[0].HealthCheckEvictionTime, 10)
				appSettings.Properties = &appSettingsProps
			}
			if appSettings != nil {
				if _, err := client.UpdateApplicationSettingsSlot(ctx, id, *appSettings); err != nil {
					return fmt.Errorf("setting App Settings for Windows %s: %+v", id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(webAppSlot.AuthSettings)
			if auth.Properties != nil {
				if _, err := client.UpdateAuthSettingsSlot(ctx, id, *auth); err != nil {
					return fmt.Errorf("setting Authorisation Settings for %s: %+v", id, err)
				}
			}

			authv2 := helpers.ExpandAuthV2Settings(webAppSlot.AuthV2Settings)
			if authv2.Properties != nil {
				if _, err = client.UpdateAuthSettingsV2Slot(ctx, id, *authv2); err != nil {
					return fmt.Errorf("updating AuthV2 settings for Linux %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("logs") {
				logsConfig := helpers.ExpandLogsConfig(webAppSlot.LogsConfig)
				if logsConfig.Properties != nil {
					if _, err := client.UpdateDiagnosticLogsConfigSlot(ctx, id, *logsConfig); err != nil {
						return fmt.Errorf("setting Diagnostic Logs Configuration for Windows %s: %+v", id, err)
					}
				}
			}

			backupConfig, err := helpers.ExpandBackupConfig(webAppSlot.Backup)
			if err != nil {
				return fmt.Errorf("expanding backup configuration for Windows %s: %+v", id, err)
			}

			if backupConfig.Properties != nil {
				if _, err := client.UpdateBackupConfigurationSlot(ctx, id, *backupConfig); err != nil {
					return fmt.Errorf("adding Backup Settings for Windows %s: %+v", id, err)
				}
			}

			storageConfig := helpers.ExpandStorageConfig(webAppSlot.StorageAccounts)
			if storageConfig.Properties != nil {
				if _, err := client.UpdateAzureStorageAccountsSlot(ctx, id, *storageConfig); err != nil {
					if err != nil {
						return fmt.Errorf("setting Storage Accounts for Windows %s: %+v", id, err)
					}
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(webAppSlot.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStringsSlot(ctx, id, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for Windows %s: %+v", id, err)
				}
			}

			if webAppSlot.ZipDeployFile != "" {
				if err = helpers.GetCredentialsAndPublishSlot(ctx, client, id, webAppSlot.ZipDeployFile); err != nil {
					return err
				}
			}

			if !webAppSlot.PublishingDeployBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: false,
					},
				}
				if _, err := client.UpdateScmAllowedSlot(ctx, id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", id, err)
				}
			}

			if !webAppSlot.PublishingFTPBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: false,
					},
				}
				if _, err := client.UpdateFtpAllowedSlot(ctx, id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for ftp publishing credentials for %s: %+v", id, err)
				}
			}

			return nil
		},

		Timeout: 30 * time.Minute,
	}
}

func (r WindowsWebAppSlotResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := webapps.ParseSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			webAppSlot, err := client.GetSlot(ctx, *id)
			if err != nil {
				if response.WasNotFound(webAppSlot.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Windows %s: %+v", *id, err)
			}

			// Despite being part of the defined `Get` response model, site_config is always nil so we get it explicitly
			webAppSiteSlotConfig, err := client.GetConfigurationSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Site Config for Windows %s: %+v", *id, err)
			}

			auth, err := client.GetAuthSettingsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Windows %s: %+v", *id, err)
			}

			var authV2 webapps.SiteAuthSettingsV2
			if strings.EqualFold(pointer.From(auth.Model.Properties.ConfigVersion), "v2") {
				authV2Resp, err := client.GetAuthSettingsV2Slot(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading authV2 settings for Linux %s: %+v", *id, err)
				}
				authV2 = *authV2Resp.Model
			}

			backup, err := client.GetBackupConfigurationSlot(ctx, *id)
			if err != nil {
				if !response.WasNotFound(backup.HttpResponse) {
					return fmt.Errorf("reading Backup Settings for Windows %s: %+v", *id, err)
				}
			}

			logsConfig, err := client.GetDiagnosticLogsConfigurationSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Diagnostic Logs information for Windows %s: %+v", *id, err)
			}

			appSettings, err := client.ListApplicationSettingsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading App Settings for Windows %s: %+v", *id, err)
			}

			storageAccounts, err := client.ListAzureStorageAccountsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Storage Account information for Windows %s: %+v", *id, err)
			}

			connectionStrings, err := client.ListConnectionStringsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Windows %s: %+v", *id, err)
			}

			siteCredentials, err := helpers.ListPublishingCredentialsSlot(ctx, client, *id)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for %s: %+v", id, err)
			}

			siteMetadata, err := client.ListMetadataSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Site Metadata for Windows %s: %+v", *id, err)
			}

			appId := commonids.NewAppServiceID(metadata.Client.Account.SubscriptionId, id.ResourceGroupName, id.SiteName)

			webApp, err := client.Get(ctx, appId)
			if err != nil {
				return fmt.Errorf("reading parent Web App for Linux %s: %+v", *id, err)
			}
			if webApp.Model.Properties == nil || webApp.Model.Properties.ServerFarmId == nil {
				return fmt.Errorf("reading parent Function App Service Plan information for Linux %s: %+v", *id, err)
			}

			basicAuthFTP := true
			if basicAuthFTPResp, err := client.GetFtpAllowedSlot(ctx, *id); err != nil || basicAuthFTPResp.Model == nil {
				return fmt.Errorf("retrieving state of FTP Basic Auth for %s: %+v", *id, err)
			} else if csmProps := basicAuthFTPResp.Model.Properties; csmProps != nil {
				basicAuthFTP = csmProps.Allow
			}

			basicAuthWebDeploy := true
			if basicAuthWebDeployResp, err := client.GetScmAllowedSlot(ctx, *id); err != nil || basicAuthWebDeployResp.Model == nil {
				return fmt.Errorf("retrieving state of WebDeploy Basic Auth for %s: %+v", *id, err)
			} else if csmProps := basicAuthWebDeployResp.Model.Properties; csmProps != nil {
				basicAuthWebDeploy = csmProps.Allow
			}

			state := WindowsWebAppSlotModel{
				Name:              id.SlotName,
				AppServiceId:      appId.ID(),
				AuthSettings:      helpers.FlattenAuthSettings(auth.Model),
				AuthV2Settings:    helpers.FlattenAuthV2Settings(authV2),
				Backup:            helpers.FlattenBackupConfig(backup.Model),
				ConnectionStrings: helpers.FlattenConnectionStrings(connectionStrings.Model),
				LogsConfig:        helpers.FlattenLogsConfig(logsConfig.Model),
				SiteCredentials:   helpers.FlattenSiteCredentials(siteCredentials),
				StorageAccounts:   helpers.FlattenStorageAccounts(storageAccounts.Model),
			}

			if model := webAppSlot.Model; model != nil {
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

					if hostingEnv := props.HostingEnvironmentProfile; hostingEnv != nil {
						hostingEnvId, err := parse.AppServiceEnvironmentIDInsensitively(*hostingEnv.Id)
						if err != nil {
							return err
						}
						state.HostingEnvId = hostingEnvId.ID()
					}

					parentAppFarmId, err := commonids.ParseAppServicePlanIDInsensitively(*webApp.Model.Properties.ServerFarmId)
					if err != nil {
						return fmt.Errorf("reading parent Service Plan ID: %+v", err)
					}
					if slotPlanIdRaw := props.ServerFarmId; slotPlanIdRaw != nil && *slotPlanIdRaw != "" && !strings.EqualFold(parentAppFarmId.ID(), *slotPlanIdRaw) {
						slotPlanId, err := commonids.ParseAppServicePlanIDInsensitively(pointer.From(slotPlanIdRaw))
						if err != nil {
							return fmt.Errorf("reading Slot Service Plan ID: %+v", err)
						}
						state.ServicePlanID = slotPlanId.ID()
					}

					if subnetId := pointer.From(props.VirtualNetworkSubnetId); subnetId != "" {
						state.VirtualNetworkSubnetID = subnetId
					}
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

				siteConfig := helpers.SiteConfigWindowsWebAppSlot{}
				siteConfig.Flatten(webAppSiteSlotConfig.Model.Properties, currentStack)
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

				state.SiteConfig = []helpers.SiteConfigWindowsWebAppSlot{siteConfig}

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

			return nil
		},
	}
}

func (r WindowsWebAppSlotResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := webapps.ParseSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			delOpts := webapps.DeleteSlotOperationOptions{
				DeleteEmptyServerFarm: pointer.To(false),
				DeleteMetrics:         pointer.To(false),
			}

			if _, err := client.DeleteSlot(ctx, *id, delOpts); err != nil {
				return fmt.Errorf("deleting Windows %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r WindowsWebAppSlotResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			appId := commonids.NewAppServiceID(metadata.Client.Account.SubscriptionId, id.ResourceGroupName, id.SiteName)

			var state WindowsWebAppSlotModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.GetSlot(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("reading Windows %s: %v", id, err)
			}

			model := *existing.Model

			if metadata.ResourceData.HasChange("service_plan_id") {
				webApp, err := client.Get(ctx, appId)
				if err != nil {
					return fmt.Errorf("reading parent Windows Web App for %s: %+v", id, err)
				}
				if webApp.Model == nil || webApp.Model.Properties == nil || webApp.Model.Properties.ServerFarmId == nil {
					return fmt.Errorf("could not determine Service Plan ID for %s: %+v", id, err)
				}
				parentServicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(*webApp.Model.Properties.ServerFarmId)
				if err != nil {
					return err
				}

				o, n := metadata.ResourceData.GetChange("service_plan_id")
				oldPlan, err := commonids.ParseAppServicePlanID(o.(string))
				if err != nil {
					return err
				}

				newPlan, err := commonids.ParseAppServicePlanID(n.(string))
				if err != nil {
					return err
				}

				// we only set `service_plan_id` when it differs from the parent `service_plan_id` which is causing issues
				// https://github.com/hashicorp/terraform-provider-azurerm/issues/21024
				// we'll error here if the `service_plan_id` equals the parent `service_plan_id`
				if strings.EqualFold(newPlan.ID(), parentServicePlanId.ID()) {
					return fmt.Errorf("`service_plan_id` should only be specified when it differs from the `service_plan_id` of the associated Web App")
				}
				locks.ByID(oldPlan.ID())
				defer locks.UnlockByID(oldPlan.ID())
				locks.ByID(newPlan.ID())
				defer locks.UnlockByID(newPlan.ID())
				if model.Properties == nil {
					return fmt.Errorf("updating Service Plan for Windows %s: Slot SiteProperties was nil", *id)
				}
				model.Properties.ServerFarmId = pointer.To(newPlan.ID())
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

			currentStack := ""
			sc := state.SiteConfig[0]

			if len(sc.ApplicationStack) == 1 {
				currentStack = sc.ApplicationStack[0].CurrentStack
			}

			if metadata.ResourceData.HasChanges("site_config", "app_settings") {
				model.Properties.SiteConfig, err = sc.ExpandForUpdate(metadata, model.Properties.SiteConfig, state.AppSettings)
				if err != nil {
					return err
				}
				model.Properties.VnetRouteAllEnabled = model.Properties.SiteConfig.VnetRouteAllEnabled
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

			if metadata.ResourceData.HasChange("virtual_network_subnet_id") {
				subnetId := metadata.ResourceData.Get("virtual_network_subnet_id").(string)
				if subnetId == "" {
					if _, err := client.DeleteSwiftVirtualNetworkSlot(ctx, *id); err != nil {
						return fmt.Errorf("removing `virtual_network_subnet_id` association for %s: %+v", *id, err)
					}
					var empty *string
					model.Properties.VirtualNetworkSubnetId = empty
				} else {
					model.Properties.VirtualNetworkSubnetId = pointer.To(subnetId)
				}
			}

			if err := client.CreateOrUpdateSlotThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating Windows %s: %+v", *id, err)
			}

			// Windows Web App Slots need the siteConfig sending individually to actually accept the `windowsFxVersion` value or it's set as `DOCKER|` only.
			siteConfigUpdate := webapps.SiteConfigResource{
				Properties: model.Properties.SiteConfig,
			}
			_, err = client.UpdateConfigurationSlot(ctx, *id, siteConfigUpdate)
			if err != nil {
				return fmt.Errorf("updating %s site config: %+v", *id, err)
			}

			siteMetadata := webapps.StringDictionary{Properties: &map[string]string{
				"CURRENT_STACK": currentStack,
			}}
			if _, err := client.UpdateMetadataSlot(ctx, *id, siteMetadata); err != nil {
				return fmt.Errorf("setting Site Metadata for Current Stack on Windows %s: %+v", *id, err)
			}

			updateLogs := false

			// App Settings can clobber logs configuration so must be updated before we send any Log updates
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

				if _, err := client.UpdateApplicationSettingsSlot(ctx, *id, *appSettingsUpdate); err != nil {
					return fmt.Errorf("updating App Settings for Windows %s: %+v", *id, err)
				}

				updateLogs = true
			}

			if metadata.ResourceData.HasChange("connection_string") {
				connectionStringUpdate := helpers.ExpandConnectionStrings(state.ConnectionStrings)
				if connectionStringUpdate.Properties == nil {
					connectionStringUpdate.Properties = &map[string]webapps.ConnStringValueTypePair{}
				}
				if _, err := client.UpdateConnectionStringsSlot(ctx, *id, *connectionStringUpdate); err != nil {
					return fmt.Errorf("updating Connection Strings for Windows %s: %+v", *id, err)
				}

				updateLogs = true
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
				if _, err := client.UpdateAuthSettingsSlot(ctx, *id, *authUpdate); err != nil {
					return fmt.Errorf("updating Auth Settings for Windows %s: %+v", *id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings_v2") {
				authV2Update := helpers.ExpandAuthV2Settings(state.AuthV2Settings)
				if _, err := client.UpdateAuthSettingsV2Slot(ctx, *id, *authV2Update); err != nil {
					return fmt.Errorf("updating AuthV2 Settings for Linux %s: %+v", *id, err)
				}
				updateLogs = true
			}

			if metadata.ResourceData.HasChange("backup") {
				backupUpdate, err := helpers.ExpandBackupConfig(state.Backup)
				if err != nil {
					return fmt.Errorf("expanding backup configuration for Windows %s: %+v", *id, err)
				}

				if backupUpdate.Properties == nil {
					if _, err := client.DeleteBackupConfigurationSlot(ctx, *id); err != nil {
						return fmt.Errorf("removing Backup Settings for Windows %s: %+v", *id, err)
					}
				} else {
					if _, err := client.UpdateBackupConfigurationSlot(ctx, *id, *backupUpdate); err != nil {
						return fmt.Errorf("updating Backup Settings for Windows %s: %+v", *id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("logs") || updateLogs {
				logsUpdate := helpers.ExpandLogsConfig(state.LogsConfig)
				if logsUpdate.Properties == nil {
					logsUpdate = helpers.DisabledLogsConfig() // The API is update only, so we need to send an update with everything switched of when a user removes the "logs" block
				}
				if _, err := client.UpdateDiagnosticLogsConfigSlot(ctx, *id, *logsUpdate); err != nil {
					return fmt.Errorf("updating Logs Config for Windows %s: %+v", *id, err)
				}
			}

			if metadata.ResourceData.HasChange("storage_account") {
				storageAccountUpdate := helpers.ExpandStorageConfig(state.StorageAccounts)
				if _, err := client.UpdateAzureStorageAccountsSlot(ctx, *id, *storageAccountUpdate); err != nil {
					return fmt.Errorf("updating Storage Accounts for Windows %s: %+v", *id, err)
				}
			}

			if metadata.ResourceData.HasChange("zip_deploy_file") {
				if err = helpers.GetCredentialsAndPublishSlot(ctx, client, *id, state.ZipDeployFile); err != nil {
					return err
				}
			}

			if metadata.ResourceData.HasChange("ftp_publish_basic_authentication_enabled") {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: state.PublishingFTPBasicAuthEnabled,
					},
				}
				if _, err := client.UpdateFtpAllowedSlot(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for ftp publishing credentials for %s: %+v", *id, err)
				}
			}

			if metadata.ResourceData.HasChange("webdeploy_publish_basic_authentication_enabled") {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: state.PublishingDeployBasicAuthEnabled,
					},
				}
				if _, err := client.UpdateScmAllowedSlot(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (r WindowsWebAppSlotResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.WindowsWebAppSlotV0toV1{},
		},
	}
}
