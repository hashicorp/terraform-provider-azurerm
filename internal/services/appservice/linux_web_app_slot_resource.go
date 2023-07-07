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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
)

type LinuxWebAppSlotResource struct{}

type LinuxWebAppSlotModel struct {
	Name                          string                              `tfschema:"name"`
	AppServiceId                  string                              `tfschema:"app_service_id"`
	ServicePlanID                 string                              `tfschema:"service_plan_id"`
	AppSettings                   map[string]string                   `tfschema:"app_settings"`
	AuthSettings                  []helpers.AuthSettings              `tfschema:"auth_settings"`
	AuthV2Settings                []helpers.AuthV2Settings            `tfschema:"auth_settings_v2"`
	Backup                        []helpers.Backup                    `tfschema:"backup"`
	ClientAffinityEnabled         bool                                `tfschema:"client_affinity_enabled"`
	ClientCertEnabled             bool                                `tfschema:"client_certificate_enabled"`
	ClientCertMode                string                              `tfschema:"client_certificate_mode"`
	ClientCertExclusionPaths      string                              `tfschema:"client_certificate_exclusion_paths"`
	Enabled                       bool                                `tfschema:"enabled"`
	HttpsOnly                     bool                                `tfschema:"https_only"`
	KeyVaultReferenceIdentityID   string                              `tfschema:"key_vault_reference_identity_id"`
	LogsConfig                    []helpers.LogsConfig                `tfschema:"logs"`
	MetaData                      map[string]string                   `tfschema:"app_metadata"`
	SiteConfig                    []helpers.SiteConfigLinuxWebAppSlot `tfschema:"site_config"`
	StorageAccounts               []helpers.StorageAccount            `tfschema:"storage_account"`
	ConnectionStrings             []helpers.ConnectionString          `tfschema:"connection_string"`
	ZipDeployFile                 string                              `tfschema:"zip_deploy_file"`
	Tags                          map[string]string                   `tfschema:"tags"`
	CustomDomainVerificationId    string                              `tfschema:"custom_domain_verification_id"`
	DefaultHostname               string                              `tfschema:"default_hostname"`
	HostingEnvId                  string                              `tfschema:"hosting_environment_id"`
	Kind                          string                              `tfschema:"kind"`
	OutboundIPAddresses           string                              `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string                            `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string                              `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string                            `tfschema:"possible_outbound_ip_address_list"`
	PublicNetworkAccess           bool                                `tfschema:"public_network_access_enabled"`
	SiteCredentials               []helpers.SiteCredential            `tfschema:"site_credential"`
	VirtualNetworkSubnetID        string                              `tfschema:"virtual_network_subnet_id"`
}

var _ sdk.ResourceWithUpdate = LinuxWebAppSlotResource{}

func (r LinuxWebAppSlotResource) ModelObject() interface{} {
	return &LinuxWebAppSlotModel{}
}

func (r LinuxWebAppSlotResource) ResourceType() string {
	return "azurerm_linux_web_app_slot"
}

func (r LinuxWebAppSlotResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.WebAppSlotID
}

func (r LinuxWebAppSlotResource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: validate.WebAppID,
		},

		// Optional

		"service_plan_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.ServicePlanID,
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
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "Required",
			ValidateFunc: validation.StringInSlice([]string{
				string(web.ClientCertModeOptional),
				string(web.ClientCertModeRequired),
				string(web.ClientCertModeOptionalInteractiveUser),
			}, false),
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

		"virtual_network_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"key_vault_reference_identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"logs": helpers.LogsConfigSchema(),

		"site_config": helpers.SiteConfigSchemaLinuxWebAppSlot(),

		"storage_account": helpers.StorageAccountSchema(),

		"zip_deploy_file": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The local path and filename of the Zip packaged application to deploy to this Windows Web App. **Note:** Using this value requires `WEBSITE_RUN_FROM_PACKAGE=1` on the App in `app_settings`.",
		},

		"tags": tags.Schema(),
	}
}

func (r LinuxWebAppSlotResource) Attributes() map[string]*pluginsdk.Schema {
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

		"app_metadata": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
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

func (r LinuxWebAppSlotResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var webAppSlot LinuxWebAppSlotModel
			if err := metadata.Decode(&webAppSlot); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			appId, err := parse.WebAppID(webAppSlot.AppServiceId)
			if err != nil {
				return err
			}

			id := parse.NewWebAppSlotID(appId.SubscriptionId, appId.ResourceGroup, appId.SiteName, webAppSlot.Name)

			webApp, err := client.Get(ctx, appId.ResourceGroup, appId.SiteName)
			if err != nil {
				return fmt.Errorf("reading parent Linux Web App for %s: %+v", id, err)
			}

			if webApp.Location == nil {
				return fmt.Errorf("could not determine location for %s: %+v", id, err)
			}

			var servicePlanId *parse.ServicePlanId
			if webAppSlot.ServicePlanID != "" {
				servicePlanId, err = parse.ServicePlanID(webAppSlot.ServicePlanID)
				if err != nil {
					return err
				}
			} else {
				if props := webApp.SiteProperties; props == nil || props.ServerFarmID == nil {
					return fmt.Errorf("could not determine Service Plan ID for %s: %+v", id, err)
				} else {
					servicePlanId, err = parse.ServicePlanID(*props.ServerFarmID)
					if err != nil {
						return err
					}
				}
			}

			existing, err := client.GetSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Linux %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			sc := webAppSlot.SiteConfig[0]
			siteConfig, err := sc.ExpandForCreate(webAppSlot.AppSettings)
			if err != nil {
				return err
			}

			expandedIdentity, err := expandIdentity(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			siteEnvelope := web.Site{
				Location: webApp.Location,
				Identity: expandedIdentity,
				Tags:     tags.FromTypedObject(webAppSlot.Tags),
				SiteProperties: &web.SiteProperties{
					ServerFarmID:          pointer.To(servicePlanId.ID()),
					Enabled:               pointer.To(webAppSlot.Enabled),
					HTTPSOnly:             pointer.To(webAppSlot.HttpsOnly),
					SiteConfig:            siteConfig,
					ClientAffinityEnabled: pointer.To(webAppSlot.ClientAffinityEnabled),
					ClientCertEnabled:     pointer.To(webAppSlot.ClientCertEnabled),
					ClientCertMode:        web.ClientCertMode(webAppSlot.ClientCertMode),
					VnetRouteAllEnabled:   siteConfig.VnetRouteAllEnabled,
				},
			}

			pna := helpers.PublicNetworkAccessEnabled
			if !webAppSlot.PublicNetworkAccess {
				pna = helpers.PublicNetworkAccessDisabled
			}

			// (@jackofallops) - Values appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
			siteEnvelope.PublicNetworkAccess = pointer.To(pna)
			siteEnvelope.SiteConfig.PublicNetworkAccess = siteEnvelope.PublicNetworkAccess

			if webAppSlot.VirtualNetworkSubnetID != "" {
				siteEnvelope.SiteProperties.VirtualNetworkSubnetID = pointer.To(webAppSlot.VirtualNetworkSubnetID)
			}

			if webAppSlot.KeyVaultReferenceIdentityID != "" {
				siteEnvelope.SiteProperties.KeyVaultReferenceIdentity = pointer.To(webAppSlot.KeyVaultReferenceIdentityID)
			}

			if webAppSlot.ClientCertExclusionPaths != "" {
				siteEnvelope.ClientCertExclusionPaths = pointer.To(webAppSlot.ClientCertExclusionPaths)
			}

			future, err := client.CreateOrUpdateSlot(ctx, id.ResourceGroup, id.SiteName, siteEnvelope, id.SlotName)
			if err != nil {
				return fmt.Errorf("creating Linux %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of Linux %s: %+v", id, err)
			}

			metadata.SetID(id)

			appSettings := helpers.ExpandAppSettingsForUpdate(webAppSlot.AppSettings)
			if metadata.ResourceData.HasChange("site_config.0.health_check_eviction_time_in_min") {
				appSettings.Properties["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"] = pointer.To(strconv.Itoa(webAppSlot.SiteConfig[0].HealthCheckEvictionTime))
			}

			if appSettings.Properties != nil {
				if _, err := client.UpdateApplicationSettingsSlot(ctx, id.ResourceGroup, id.SiteName, *appSettings, id.SlotName); err != nil {
					return fmt.Errorf("setting App Settings for Linux %s: %+v", id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(webAppSlot.AuthSettings)
			if auth.SiteAuthSettingsProperties != nil {
				if _, err := client.UpdateAuthSettingsSlot(ctx, id.ResourceGroup, id.SiteName, *auth, id.SlotName); err != nil {
					return fmt.Errorf("setting Authorisation Settings for Linux %s: %+v", id, err)
				}
			}

			authv2 := helpers.ExpandAuthV2Settings(webAppSlot.AuthV2Settings)
			if authv2.SiteAuthSettingsV2Properties != nil {
				if _, err = client.UpdateAuthSettingsV2Slot(ctx, id.ResourceGroup, id.SiteName, *authv2, id.SlotName); err != nil {
					return fmt.Errorf("updating AuthV2 settings for Linux %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("logs") {
				logsConfig := helpers.ExpandLogsConfig(webAppSlot.LogsConfig)
				if logsConfig.SiteLogsConfigProperties != nil {
					if _, err := client.UpdateDiagnosticLogsConfigSlot(ctx, id.ResourceGroup, id.SiteName, *logsConfig, id.SlotName); err != nil {
						return fmt.Errorf("setting Diagnostic Logs Configuration for Linux %s: %+v", id, err)
					}
				}
			}

			backupConfig, err := helpers.ExpandBackupConfig(webAppSlot.Backup)
			if err != nil {
				return fmt.Errorf("expanding backup configuration for Linux %s: %+v", id, err)
			}
			if backupConfig.BackupRequestProperties != nil {
				if _, err := client.UpdateBackupConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, *backupConfig, id.SlotName); err != nil {
					return fmt.Errorf("adding Backup Settings for Linux %s: %+v", id, err)
				}
			}

			storageConfig := helpers.ExpandStorageConfig(webAppSlot.StorageAccounts)
			if storageConfig.Properties != nil {
				if _, err := client.UpdateAzureStorageAccountsSlot(ctx, id.ResourceGroup, id.SiteName, *storageConfig, id.SlotName); err != nil {
					if err != nil {
						return fmt.Errorf("setting Storage Accounts for Linux %s: %+v", id, err)
					}
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(webAppSlot.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStringsSlot(ctx, id.ResourceGroup, id.SiteName, *connectionStrings, id.SlotName); err != nil {
					return fmt.Errorf("setting Connection Strings for Linux %s: %+v", id, err)
				}
			}

			if webAppSlot.ZipDeployFile != "" {
				if err = helpers.GetCredentialsAndPublish(ctx, client, id.ResourceGroup, id.SiteName, webAppSlot.ZipDeployFile); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func (r LinuxWebAppSlotResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := parse.WebAppSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			webAppSlot, err := client.GetSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				if utils.ResponseWasNotFound(webAppSlot.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Linux %s: %+v", id, err)
			}

			// Despite being part of the defined `Get` response model, site_config is always nil so we get it explicitly
			webAppSiteConfig, err := client.GetConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading Site Config for Linux %s: %+v", id, err)
			}

			auth, err := client.GetAuthSettingsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Linux %s: %+v", id, err)
			}

			var authV2 web.SiteAuthSettingsV2
			if pointer.From(auth.ConfigVersion) == "v2" {
				authV2, err = client.GetAuthSettingsV2Slot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
				if err != nil {
					return fmt.Errorf("reading authV2 settings for Linux %s: %+v", *id, err)
				}
			}

			backup, err := client.GetBackupConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				if !utils.ResponseWasNotFound(backup.Response) {
					return fmt.Errorf("reading Backup Settings for Linux %s: %+v", id, err)
				}
			}

			logsConfig, err := client.GetDiagnosticLogsConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading Diagnostic Logs information for Linux %s: %+v", id, err)
			}

			appSettings, err := client.ListApplicationSettingsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading App Settings for Linux %s: %+v", id, err)
			}

			storageAccounts, err := client.ListAzureStorageAccountsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading Storage Account information for Linux %s: %+v", id, err)
			}

			connectionStrings, err := client.ListConnectionStringsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Linux %s: %+v", id, err)
			}

			siteCredentialsFuture, err := client.ListPublishingCredentialsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for Linux %s: %+v", id, err)
			}

			if err := siteCredentialsFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for Site Publishing Credential information for Linux %s: %+v", id, err)
			}
			siteCredentials, err := siteCredentialsFuture.Result(*client)
			if err != nil {
				return fmt.Errorf("reading Site Publishing Credential information for Linux %s: %+v", id, err)
			}

			webApp, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading parent Web App for Linux %s: %+v", *id, err)
			}
			if webApp.SiteProperties == nil || webApp.SiteProperties.ServerFarmID == nil {
				return fmt.Errorf("reading parent Function App Service Plan information for Linux %s: %+v", *id, err)
			}

			state := LinuxWebAppSlotModel{}
			if props := webAppSlot.SiteProperties; props != nil {
				state = LinuxWebAppSlotModel{
					Name:                          id.SlotName,
					AppServiceId:                  parse.NewWebAppID(id.SubscriptionId, id.ResourceGroup, id.SiteName).ID(),
					ClientAffinityEnabled:         pointer.From(props.ClientAffinityEnabled),
					ClientCertEnabled:             pointer.From(props.ClientCertEnabled),
					ClientCertMode:                string(props.ClientCertMode),
					ClientCertExclusionPaths:      pointer.From(props.ClientCertExclusionPaths),
					CustomDomainVerificationId:    pointer.From(props.CustomDomainVerificationID),
					DefaultHostname:               pointer.From(props.DefaultHostName),
					Kind:                          pointer.From(webAppSlot.Kind),
					KeyVaultReferenceIdentityID:   pointer.From(props.KeyVaultReferenceIdentity),
					Enabled:                       pointer.From(props.Enabled),
					HttpsOnly:                     pointer.From(props.HTTPSOnly),
					OutboundIPAddresses:           pointer.From(props.OutboundIPAddresses),
					OutboundIPAddressList:         strings.Split(pointer.From(props.OutboundIPAddresses), ","),
					PossibleOutboundIPAddresses:   pointer.From(props.PossibleOutboundIPAddresses),
					PossibleOutboundIPAddressList: strings.Split(pointer.From(props.PossibleOutboundIPAddresses), ","),
					PublicNetworkAccess:           !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled),
					Tags:                          tags.ToTypedObject(webAppSlot.Tags),
				}

				if hostingEnv := props.HostingEnvironmentProfile; hostingEnv != nil {
					state.HostingEnvId = pointer.From(hostingEnv.ID)
				}

				if subnetId := pointer.From(props.VirtualNetworkSubnetID); subnetId != "" {
					state.VirtualNetworkSubnetID = subnetId
				}

				parentAppFarmId, err := parse.ServicePlanIDInsensitively(*webApp.SiteProperties.ServerFarmID)
				if err != nil {
					return err
				}
				if slotPlanId := props.ServerFarmID; slotPlanId != nil && !strings.EqualFold(parentAppFarmId.ID(), *slotPlanId) {
					state.ServicePlanID = *slotPlanId
				}

				if subnetId := pointer.From(props.VirtualNetworkSubnetID); subnetId != "" {
					state.VirtualNetworkSubnetID = subnetId
				}
			}

			state.AppSettings = helpers.FlattenWebStringDictionary(appSettings)
			if err != nil {
				return fmt.Errorf("flattening app settings for Linux %s: %+v", id, err)
			}

			state.AuthSettings = helpers.FlattenAuthSettings(auth)

			state.AuthV2Settings = helpers.FlattenAuthV2Settings(authV2)

			state.Backup = helpers.FlattenBackupConfig(backup)

			state.LogsConfig = helpers.FlattenLogsConfig(logsConfig)

			siteConfig := helpers.SiteConfigLinuxWebAppSlot{}
			siteConfig.Flatten(webAppSiteConfig.SiteConfig)
			siteConfig.SetHealthCheckEvictionTime(state.AppSettings)

			// For non-import cases we check for use of the deprecated docker settings - remove in 4.0
			_, usesDeprecatedDocker := metadata.ResourceData.GetOk("site_config.0.application_stack.0.docker_image")

			if helpers.FxStringHasPrefix(siteConfig.LinuxFxVersion, helpers.FxStringPrefixDocker) {
				if !features.FourPointOhBeta() {
					siteConfig.DecodeDockerDeprecatedAppStack(state.AppSettings, usesDeprecatedDocker)
				} else {
					siteConfig.DecodeDockerAppStack(state.AppSettings)
				}
			}

			state.SiteConfig = []helpers.SiteConfigLinuxWebAppSlot{siteConfig}

			// Filter out all settings we've consumed above
			if !features.FourPointOhBeta() && usesDeprecatedDocker {
				state.AppSettings = helpers.FilterManagedAppSettingsDeprecated(state.AppSettings)
			} else {
				state.AppSettings = helpers.FilterManagedAppSettings(state.AppSettings)
			}

			state.StorageAccounts = helpers.FlattenStorageAccounts(storageAccounts)

			state.ConnectionStrings = helpers.FlattenConnectionStrings(connectionStrings)

			state.SiteCredentials = helpers.FlattenSiteCredentials(siteCredentials)

			// Zip Deploys are not retrievable, so attempt to get from config. This doesn't matter for imports as an unexpected value here could break the deployment.
			if deployFile, ok := metadata.ResourceData.Get("zip_deploy_file").(string); ok {
				state.ZipDeployFile = deployFile
			}

			if err := metadata.Encode(&state); err != nil {
				return fmt.Errorf("encoding: %+v", err)
			}

			flattenedIdentity, err := flattenIdentity(webAppSlot.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			if err := metadata.ResourceData.Set("identity", flattenedIdentity); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			return nil
		},
	}
}

func (r LinuxWebAppSlotResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := parse.WebAppSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			deleteMetrics := true
			deleteEmptyServerFarm := false
			if _, err := client.DeleteSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName, &deleteMetrics, &deleteEmptyServerFarm); err != nil {
				return fmt.Errorf("deleting Linux %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r LinuxWebAppSlotResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := parse.WebAppSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LinuxWebAppSlotModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.GetSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading Linux %s: %v", id, err)
			}

			if metadata.ResourceData.HasChange("service_plan_id") {
				o, n := metadata.ResourceData.GetChange("service_plan_id")
				oldPlan, err := parse.ServicePlanID(o.(string))
				if err != nil {
					return err
				}

				newPlan, err := parse.ServicePlanID(n.(string))
				if err != nil {
					return err
				}
				locks.ByID(oldPlan.ID())
				defer locks.UnlockByID(oldPlan.ID())
				locks.ByID(newPlan.ID())
				defer locks.UnlockByID(newPlan.ID())
				if existing.SiteProperties == nil {
					return fmt.Errorf("updating Service Plan for Linux %s: Slot SiteProperties was nil", *id)
				}
				existing.SiteProperties.ServerFarmID = pointer.To(newPlan.ID())
			}

			if metadata.ResourceData.HasChange("enabled") {
				existing.SiteProperties.Enabled = pointer.To(state.Enabled)
			}
			if metadata.ResourceData.HasChange("https_only") {
				existing.SiteProperties.HTTPSOnly = pointer.To(state.HttpsOnly)
			}
			if metadata.ResourceData.HasChange("client_affinity_enabled") {
				existing.SiteProperties.ClientAffinityEnabled = pointer.To(state.ClientAffinityEnabled)
			}
			if metadata.ResourceData.HasChange("client_certificate_enabled") {
				existing.SiteProperties.ClientCertEnabled = pointer.To(state.ClientCertEnabled)
			}
			if metadata.ResourceData.HasChange("client_certificate_mode") {
				existing.SiteProperties.ClientCertMode = web.ClientCertMode(state.ClientCertMode)
			}
			if metadata.ResourceData.HasChange("client_certificate_exclusion_paths") {
				existing.SiteProperties.ClientCertExclusionPaths = pointer.To(state.ClientCertExclusionPaths)
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := expandIdentity(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				existing.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("key_vault_reference_identity_id") {
				existing.KeyVaultReferenceIdentity = pointer.To(state.KeyVaultReferenceIdentityID)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.FromTypedObject(state.Tags)
			}

			if metadata.ResourceData.HasChange("site_config") {
				sc := state.SiteConfig[0]
				siteConfig, err := sc.ExpandForUpdate(metadata, existing.SiteConfig, state.AppSettings)
				if err != nil {
					return fmt.Errorf("expanding Site Config for Linux %s: %+v", id, err)
				}
				existing.SiteConfig = siteConfig
				existing.VnetRouteAllEnabled = existing.SiteConfig.VnetRouteAllEnabled
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				pna := helpers.PublicNetworkAccessEnabled
				if !state.PublicNetworkAccess {
					pna = helpers.PublicNetworkAccessDisabled
				}

				// (@jackofallops) - Values appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
				existing.PublicNetworkAccess = pointer.To(pna)
				existing.SiteConfig.PublicNetworkAccess = existing.PublicNetworkAccess
			}

			if metadata.ResourceData.HasChange("virtual_network_subnet_id") {
				subnetId := metadata.ResourceData.Get("virtual_network_subnet_id").(string)
				if subnetId == "" {
					if _, err := client.DeleteSwiftVirtualNetworkSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName); err != nil {
						return fmt.Errorf("removing `virtual_network_subnet_id` association for %s: %+v", *id, err)
					}
					var empty *string
					existing.SiteProperties.VirtualNetworkSubnetID = empty
				} else {
					existing.SiteProperties.VirtualNetworkSubnetID = pointer.To(subnetId)
				}
			}

			updateFuture, err := client.CreateOrUpdateSlot(ctx, id.ResourceGroup, id.SiteName, existing, id.SlotName)
			if err != nil {
				return fmt.Errorf("updating Linux %s: %+v", id, err)
			}
			if err := updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting to update %s: %+v", id, err)
			}

			// (@jackofallops) - App Settings can clobber logs configuration so must be updated before we send any Log updates
			if metadata.ResourceData.HasChange("app_settings") || metadata.ResourceData.HasChange("site_config.0.health_check_eviction_time_in_min") {
				appSettingsUpdate := helpers.ExpandAppSettingsForUpdate(state.AppSettings)
				appSettingsUpdate.Properties["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"] = pointer.To(strconv.Itoa(state.SiteConfig[0].HealthCheckEvictionTime))

				if _, err := client.UpdateApplicationSettingsSlot(ctx, id.ResourceGroup, id.SiteName, *appSettingsUpdate, id.SlotName); err != nil {
					return fmt.Errorf("updating App Settings for Linux %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("connection_string") {
				connectionStringUpdate := helpers.ExpandConnectionStrings(state.ConnectionStrings)
				if connectionStringUpdate.Properties == nil {
					connectionStringUpdate.Properties = map[string]*web.ConnStringValueTypePair{}
				}
				if _, err := client.UpdateConnectionStringsSlot(ctx, id.ResourceGroup, id.SiteName, *connectionStringUpdate, id.SlotName); err != nil {
					return fmt.Errorf("updating Connection Strings for Linux %s: %+v", id, err)
				}
			}

			updateLogs := false

			if metadata.ResourceData.HasChange("auth_settings") {
				authUpdate := helpers.ExpandAuthSettings(state.AuthSettings)
				if authUpdate.SiteAuthSettingsProperties == nil {
					authUpdate.SiteAuthSettingsProperties = &web.SiteAuthSettingsProperties{
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
				if _, err := client.UpdateAuthSettingsSlot(ctx, id.ResourceGroup, id.SiteName, *authUpdate, id.SlotName); err != nil {
					return fmt.Errorf("updating Auth Settings for Linux %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings_v2") {
				authV2Update := helpers.ExpandAuthV2Settings(state.AuthV2Settings)
				if _, err := client.UpdateAuthSettingsV2Slot(ctx, id.ResourceGroup, id.SiteName, *authV2Update, id.SlotName); err != nil {
					return fmt.Errorf("updating AuthV2 Settings for Linux %s: %+v", id, err)
				}
				updateLogs = true
			}

			if metadata.ResourceData.HasChange("backup") {
				backupUpdate, err := helpers.ExpandBackupConfig(state.Backup)
				if err != nil {
					return fmt.Errorf("expanding backup configuration for Linux %s: %+v", *id, err)
				}
				if backupUpdate.BackupRequestProperties == nil {
					if _, err := client.DeleteBackupConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName); err != nil {
						return fmt.Errorf("removing Backup Settings for Linux %s: %+v", id, err)
					}
				} else {
					if _, err := client.UpdateBackupConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, *backupUpdate, id.SlotName); err != nil {
						return fmt.Errorf("updating Backup Settings for Linux %s: %+v", id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("logs") || updateLogs {
				logsUpdate := helpers.ExpandLogsConfig(state.LogsConfig)
				if logsUpdate.SiteLogsConfigProperties == nil {
					logsUpdate = helpers.DisabledLogsConfig() // The API is update only, so we need to send an update with everything switched of when a user removes the "logs" block
				}
				if _, err := client.UpdateDiagnosticLogsConfigSlot(ctx, id.ResourceGroup, id.SiteName, *logsUpdate, id.SlotName); err != nil {
					return fmt.Errorf("updating Logs Config for Linux %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("storage_account") {
				storageAccountUpdate := helpers.ExpandStorageConfig(state.StorageAccounts)
				if _, err := client.UpdateAzureStorageAccountsSlot(ctx, id.ResourceGroup, id.SiteName, *storageAccountUpdate, id.SlotName); err != nil {
					return fmt.Errorf("updating Storage Accounts for Linux %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("zip_deploy_file") {
				if err = helpers.GetCredentialsAndPublishSlot(ctx, client, id.ResourceGroup, id.SiteName, state.ZipDeployFile, id.SlotName); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
