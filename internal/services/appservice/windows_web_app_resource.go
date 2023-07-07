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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
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

type WindowsWebAppResource struct{}

type WindowsWebAppModel struct {
	Name                          string                      `tfschema:"name"`
	ResourceGroup                 string                      `tfschema:"resource_group_name"`
	Location                      string                      `tfschema:"location"`
	ServicePlanId                 string                      `tfschema:"service_plan_id"`
	AppSettings                   map[string]string           `tfschema:"app_settings"`
	StickySettings                []helpers.StickySettings    `tfschema:"sticky_settings"`
	AuthSettings                  []helpers.AuthSettings      `tfschema:"auth_settings"`
	AuthV2Settings                []helpers.AuthV2Settings    `tfschema:"auth_settings_v2"`
	Backup                        []helpers.Backup            `tfschema:"backup"`
	ClientAffinityEnabled         bool                        `tfschema:"client_affinity_enabled"`
	ClientCertEnabled             bool                        `tfschema:"client_certificate_enabled"`
	ClientCertMode                string                      `tfschema:"client_certificate_mode"`
	ClientCertExclusionPaths      string                      `tfschema:"client_certificate_exclusion_paths"`
	Enabled                       bool                        `tfschema:"enabled"`
	HttpsOnly                     bool                        `tfschema:"https_only"`
	KeyVaultReferenceIdentityID   string                      `tfschema:"key_vault_reference_identity_id"`
	LogsConfig                    []helpers.LogsConfig        `tfschema:"logs"`
	PublicNetworkAccess           bool                        `tfschema:"public_network_access_enabled"`
	SiteConfig                    []helpers.SiteConfigWindows `tfschema:"site_config"`
	StorageAccounts               []helpers.StorageAccount    `tfschema:"storage_account"`
	ConnectionStrings             []helpers.ConnectionString  `tfschema:"connection_string"`
	CustomDomainVerificationId    string                      `tfschema:"custom_domain_verification_id"`
	HostingEnvId                  string                      `tfschema:"hosting_environment_id"`
	DefaultHostname               string                      `tfschema:"default_hostname"`
	Kind                          string                      `tfschema:"kind"`
	OutboundIPAddresses           string                      `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string                    `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string                      `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string                    `tfschema:"possible_outbound_ip_address_list"`
	SiteCredentials               []helpers.SiteCredential    `tfschema:"site_credential"`
	ZipDeployFile                 string                      `tfschema:"zip_deploy_file"`
	Tags                          map[string]string           `tfschema:"tags"`
	VirtualNetworkSubnetID        string                      `tfschema:"virtual_network_subnet_id"`
}

var _ sdk.ResourceWithCustomImporter = WindowsWebAppResource{}

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
			ValidateFunc: validate.ServicePlanID,
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
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(web.ClientCertModeRequired),
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

// TODO - Feature: Deployments (Preview)?
// TODO - Feature: App Insights?

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
			servicePlanClient := metadata.Client.AppService.ServicePlanClient
			aseClient := metadata.Client.AppService.AppServiceEnvironmentClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewWebAppID(subscriptionId, webApp.ResourceGroup, webApp.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Windows %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			availabilityRequest := web.ResourceNameAvailabilityRequest{
				Name: pointer.To(webApp.Name),
				Type: web.CheckNameResourceTypesMicrosoftWebsites,
			}

			servicePlanId, err := parse.ServicePlanID(webApp.ServicePlanId)
			if err != nil {
				return err
			}

			servicePlan, err := servicePlanClient.Get(ctx, servicePlanId.ResourceGroup, servicePlanId.ServerfarmName)
			if err != nil {
				return fmt.Errorf("reading App %s: %+v", servicePlanId, err)
			}
			if ase := servicePlan.HostingEnvironmentProfile; ase != nil {
				// Attempt to check the ASE for the appropriate suffix for the name availability request. Not convinced
				// the `DNSSuffix` field is still valid and possibly should have been deprecated / removed as is legacy
				// setting from ASEv1? Hence the non-fatal approach here.
				nameSuffix := "appserviceenvironment.net"
				if ase.ID != nil {
					aseId, err := parse.AppServiceEnvironmentID(*ase.ID)
					nameSuffix = fmt.Sprintf("%s.%s", aseId.HostingEnvironmentName, nameSuffix)
					if err != nil {
						metadata.Logger.Warnf("could not parse App Service Environment ID determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", webApp.Name, servicePlanId)
					} else {
						existingASE, err := aseClient.Get(ctx, aseId.ResourceGroup, aseId.HostingEnvironmentName)
						if err != nil {
							metadata.Logger.Warnf("could not read App Service Environment to determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", webApp.Name, servicePlanId)
						} else if props := existingASE.AppServiceEnvironment; props != nil && props.DNSSuffix != nil && *props.DNSSuffix != "" {
							nameSuffix = *props.DNSSuffix
						}
					}
				}
				availabilityRequest.Name = pointer.To(fmt.Sprintf("%s.%s", webApp.Name, nameSuffix))
				availabilityRequest.IsFqdn = pointer.To(true)
			}

			checkName, err := client.CheckNameAvailability(ctx, availabilityRequest)
			if err != nil {
				return fmt.Errorf("checking name availability for %s: %+v", id, err)
			}
			if !*checkName.NameAvailable {
				return fmt.Errorf("the Site Name %q failed the availability check: %+v", id.SiteName, *checkName.Message)
			}

			sc := webApp.SiteConfig[0]

			if servicePlan.Sku != nil && servicePlan.Sku.Name != nil {
				if helpers.IsFreeOrSharedServicePlan(*servicePlan.Sku.Name) {
					if sc.AlwaysOn {
						return fmt.Errorf("always_on cannot be set to true when using Free, F1, D1 Sku")
					}
				}
			}

			siteConfig, err := sc.ExpandForCreate(webApp.AppSettings)
			if err != nil {
				return err
			}

			currentStack := ""
			if len(sc.ApplicationStack) == 1 {
				currentStack = sc.ApplicationStack[0].CurrentStack
				if currentStack == helpers.CurrentStackNode || sc.ApplicationStack[0].NodeVersion != "" {
					if webApp.AppSettings == nil {
						webApp.AppSettings = make(map[string]string, 0)
					}
					webApp.AppSettings["WEBSITE_NODE_DEFAULT_VERSION"] = sc.ApplicationStack[0].NodeVersion
				}
			}

			expandedIdentity, err := expandIdentity(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			siteEnvelope := web.Site{
				Location: pointer.To(webApp.Location),
				Tags:     tags.FromTypedObject(webApp.Tags),
				Identity: expandedIdentity,
				SiteProperties: &web.SiteProperties{
					ServerFarmID:          pointer.To(webApp.ServicePlanId),
					Enabled:               pointer.To(webApp.Enabled),
					HTTPSOnly:             pointer.To(webApp.HttpsOnly),
					SiteConfig:            siteConfig,
					ClientAffinityEnabled: pointer.To(webApp.ClientAffinityEnabled),
					ClientCertEnabled:     pointer.To(webApp.ClientCertEnabled),
					ClientCertMode:        web.ClientCertMode(webApp.ClientCertMode),
					VnetRouteAllEnabled:   siteConfig.VnetRouteAllEnabled,
				},
			}

			pna := helpers.PublicNetworkAccessEnabled
			if !webApp.PublicNetworkAccess {
				pna = helpers.PublicNetworkAccessDisabled
			}

			// (@jackofallops) - Values appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
			siteEnvelope.PublicNetworkAccess = pointer.To(pna)
			siteEnvelope.SiteConfig.PublicNetworkAccess = siteEnvelope.PublicNetworkAccess

			if webApp.KeyVaultReferenceIdentityID != "" {
				siteEnvelope.SiteProperties.KeyVaultReferenceIdentity = pointer.To(webApp.KeyVaultReferenceIdentityID)
			}

			if webApp.VirtualNetworkSubnetID != "" {
				siteEnvelope.SiteProperties.VirtualNetworkSubnetID = pointer.To(webApp.VirtualNetworkSubnetID)
			}

			if webApp.ClientCertExclusionPaths != "" {
				siteEnvelope.ClientCertExclusionPaths = pointer.To(webApp.ClientCertExclusionPaths)
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, siteEnvelope)
			if err != nil {
				return fmt.Errorf("creating Windows %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of Windows %s: %+v", id, err)
			}

			metadata.SetID(id)

			if currentStack != "" {
				siteMetadata := web.StringDictionary{Properties: map[string]*string{}}
				siteMetadata.Properties["CURRENT_STACK"] = pointer.To(currentStack)
				if _, err := client.UpdateMetadata(ctx, id.ResourceGroup, id.SiteName, siteMetadata); err != nil {
					return fmt.Errorf("setting Site Metadata for Current Stack on Windows %s: %+v", id, err)
				}
			}

			appSettings := helpers.ExpandAppSettingsForUpdate(webApp.AppSettings)
			if metadata.ResourceData.HasChange("site_config.0.health_check_eviction_time_in_min") {
				appSettings.Properties["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"] = pointer.To(strconv.Itoa(webApp.SiteConfig[0].HealthCheckEvictionTime))
			}

			if appSettings != nil {
				if _, err := client.UpdateApplicationSettings(ctx, id.ResourceGroup, id.SiteName, *appSettings); err != nil {
					return fmt.Errorf("setting App Settings for Windows %s: %+v", id, err)
				}
			}

			stickySettings := helpers.ExpandStickySettings(webApp.StickySettings)

			if stickySettings != nil {
				stickySettingsUpdate := web.SlotConfigNamesResource{
					SlotConfigNames: stickySettings,
				}
				if _, err := client.UpdateSlotConfigurationNames(ctx, id.ResourceGroup, id.SiteName, stickySettingsUpdate); err != nil {
					return fmt.Errorf("updating Sticky Settings for Windows %s: %+v", id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(webApp.AuthSettings)
			if auth.SiteAuthSettingsProperties != nil {
				if _, err := client.UpdateAuthSettings(ctx, id.ResourceGroup, id.SiteName, *auth); err != nil {
					return fmt.Errorf("setting Authorisation Settings for %s: %+v", id, err)
				}
			}

			authv2 := helpers.ExpandAuthV2Settings(webApp.AuthV2Settings)
			if authv2.SiteAuthSettingsV2Properties != nil {
				if _, err = client.UpdateAuthSettingsV2(ctx, id.ResourceGroup, id.SiteName, *authv2); err != nil {
					return fmt.Errorf("updating AuthV2 settings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("logs") {
				logsConfig := helpers.ExpandLogsConfig(webApp.LogsConfig)
				if logsConfig.SiteLogsConfigProperties != nil {
					if _, err := client.UpdateDiagnosticLogsConfig(ctx, id.ResourceGroup, id.SiteName, *logsConfig); err != nil {
						return fmt.Errorf("setting Diagnostic Logs Configuration for Windows %s: %+v", id, err)
					}
				}
			}

			backupConfig, err := helpers.ExpandBackupConfig(webApp.Backup)
			if err != nil {
				return fmt.Errorf("expanding backup configuration for Windows %s: %+v", id, err)
			}

			if backupConfig.BackupRequestProperties != nil {
				if _, err := client.UpdateBackupConfiguration(ctx, id.ResourceGroup, id.SiteName, *backupConfig); err != nil {
					return fmt.Errorf("adding Backup Settings for Windows %s: %+v", id, err)
				}
			}

			storageConfig := helpers.ExpandStorageConfig(webApp.StorageAccounts)
			if storageConfig.Properties != nil {
				if _, err := client.UpdateAzureStorageAccounts(ctx, id.ResourceGroup, id.SiteName, *storageConfig); err != nil {
					if err != nil {
						return fmt.Errorf("setting Storage Accounts for Windows %s: %+v", id, err)
					}
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(webApp.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStrings(ctx, id.ResourceGroup, id.SiteName, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for Windows %s: %+v", id, err)
				}
			}

			if webApp.ZipDeployFile != "" {
				if err = helpers.GetCredentialsAndPublish(ctx, client, id.ResourceGroup, id.SiteName, webApp.ZipDeployFile); err != nil {
					return err
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
			id, err := parse.WebAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			webApp, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				if utils.ResponseWasNotFound(webApp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Windows %s: %+v", id, err)
			}

			// Despite being part of the defined `Get` response model, site_config is always nil so we get it explicitly
			webAppSiteConfig, err := client.GetConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Site Config for Windows %s: %+v", id, err)
			}

			auth, err := client.GetAuthSettings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Windows %s: %+v", id, err)
			}

			var authV2 web.SiteAuthSettingsV2
			if pointer.From(auth.ConfigVersion) == "v2" {
				authV2, err = client.GetAuthSettingsV2(ctx, id.ResourceGroup, id.SiteName)
				if err != nil {
					return fmt.Errorf("reading authV2 settings for Windows %s: %+v", *id, err)
				}
			}

			backup, err := client.GetBackupConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				if !utils.ResponseWasNotFound(backup.Response) {
					return fmt.Errorf("reading Backup Settings for Windows %s: %+v", id, err)
				}
			}

			logsConfig, err := client.GetDiagnosticLogsConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Diagnostic Logs information for Windows %s: %+v", id, err)
			}

			appSettings, err := client.ListApplicationSettings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading App Settings for Windows %s: %+v", id, err)
			}

			stickySettings, err := client.ListSlotConfigurationNames(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Sticky Settings for Linux %s: %+v", id, err)
			}

			storageAccounts, err := client.ListAzureStorageAccounts(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Storage Account information for Windows %s: %+v", id, err)
			}

			connectionStrings, err := client.ListConnectionStrings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Windows %s: %+v", id, err)
			}

			siteCredentialsFuture, err := client.ListPublishingCredentials(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for Windows %s: %+v", id, err)
			}

			if err := siteCredentialsFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for Site Publishing Credential information for Windows %s: %+v", id, err)
			}
			siteCredentials, err := siteCredentialsFuture.Result(*client)
			if err != nil {
				return fmt.Errorf("reading Site Publishing Credential information for Windows %s: %+v", id, err)
			}

			siteMetadata, err := client.ListMetadata(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Site Metadata for Windows %s: %+v", id, err)
			}

			state := WindowsWebAppModel{}
			if props := webApp.SiteProperties; props != nil {
				state = WindowsWebAppModel{
					Name:                          id.SiteName,
					ResourceGroup:                 id.ResourceGroup,
					Location:                      location.NormalizeNilable(webApp.Location),
					AuthSettings:                  helpers.FlattenAuthSettings(auth),
					AuthV2Settings:                helpers.FlattenAuthV2Settings(authV2),
					Backup:                        helpers.FlattenBackupConfig(backup),
					ClientAffinityEnabled:         pointer.From(props.ClientAffinityEnabled),
					ClientCertEnabled:             pointer.From(props.ClientCertEnabled),
					ClientCertMode:                string(props.ClientCertMode),
					ClientCertExclusionPaths:      pointer.From(props.ClientCertExclusionPaths),
					ConnectionStrings:             helpers.FlattenConnectionStrings(connectionStrings),
					CustomDomainVerificationId:    pointer.From(props.CustomDomainVerificationID),
					DefaultHostname:               pointer.From(props.DefaultHostName),
					Enabled:                       pointer.From(props.Enabled),
					HttpsOnly:                     pointer.From(props.HTTPSOnly),
					KeyVaultReferenceIdentityID:   pointer.From(props.KeyVaultReferenceIdentity),
					Kind:                          pointer.From(webApp.Kind),
					LogsConfig:                    helpers.FlattenLogsConfig(logsConfig),
					SiteCredentials:               helpers.FlattenSiteCredentials(siteCredentials),
					StorageAccounts:               helpers.FlattenStorageAccounts(storageAccounts),
					StickySettings:                helpers.FlattenStickySettings(stickySettings.SlotConfigNames),
					OutboundIPAddresses:           pointer.From(props.OutboundIPAddresses),
					OutboundIPAddressList:         strings.Split(pointer.From(props.OutboundIPAddresses), ","),
					PossibleOutboundIPAddresses:   pointer.From(props.PossibleOutboundIPAddresses),
					PossibleOutboundIPAddressList: strings.Split(pointer.From(props.PossibleOutboundIPAddresses), ","),
					PublicNetworkAccess:           !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled),
					Tags:                          tags.ToTypedObject(webApp.Tags),
				}

				serverFarmId, err := parse.ServicePlanID(pointer.From(props.ServerFarmID))
				if err != nil {
					return fmt.Errorf("parsing Service Plan ID for %s: %+v", id, err)
				}

				state.ServicePlanId = serverFarmId.ID()

				if hostingEnv := props.HostingEnvironmentProfile; hostingEnv != nil {
					hostingEnvId, err := parse.AppServiceEnvironmentIDInsensitively(*hostingEnv.ID)
					if err != nil {
						return err
					}
					state.HostingEnvId = hostingEnvId.ID()
				}

				if subnetId := pointer.From(props.VirtualNetworkSubnetID); subnetId != "" {
					state.VirtualNetworkSubnetID = subnetId
				}

			}

			state.AppSettings = helpers.FlattenWebStringDictionary(appSettings)

			currentStack := ""
			currentStackPtr, ok := siteMetadata.Properties["CURRENT_STACK"]
			if ok {
				currentStack = *currentStackPtr
			}

			siteConfig := helpers.SiteConfigWindows{}
			if err := siteConfig.Flatten(webAppSiteConfig.SiteConfig, currentStack); err != nil {
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

			if err := metadata.Encode(&state); err != nil {
				return fmt.Errorf("encoding: %+v", err)
			}

			flattenedIdentity, err := flattenIdentity(webApp.Identity)
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

func (r WindowsWebAppResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := parse.WebAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			deleteMetrics := true // TODO - Look at making this a feature flag?
			deleteEmptyServerFarm := false
			if _, err := client.Delete(ctx, id.ResourceGroup, id.SiteName, &deleteMetrics, &deleteEmptyServerFarm); err != nil {
				return fmt.Errorf("deleting Windows %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r WindowsWebAppResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.WebAppID
}

func (r WindowsWebAppResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			servicePlanClient := metadata.Client.AppService.ServicePlanClient

			id, err := parse.WebAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state WindowsWebAppModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Windows %s: %v", id, err)
			}

			var serviceFarmId string
			servicePlanChange := false
			if existing.SiteProperties.ServerFarmID != nil {
				serviceFarmId = *existing.ServerFarmID
			}
			if metadata.ResourceData.HasChange("service_plan_id") {
				serviceFarmId = state.ServicePlanId
				existing.SiteProperties.ServerFarmID = pointer.To(serviceFarmId)
				servicePlanChange = true
			}
			servicePlanId, err := parse.ServicePlanID(serviceFarmId)
			if err != nil {
				return err
			}

			servicePlan, err := servicePlanClient.Get(ctx, servicePlanId.ResourceGroup, servicePlanId.ServerfarmName)
			if err != nil {
				return fmt.Errorf("reading App %s: %+v", servicePlanId, err)
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

			if metadata.ResourceData.HasChange("virtual_network_subnet_id") {
				subnetId := metadata.ResourceData.Get("virtual_network_subnet_id").(string)
				if subnetId == "" {
					if _, err := client.DeleteSwiftVirtualNetwork(ctx, id.ResourceGroup, id.SiteName); err != nil {
						return fmt.Errorf("removing `virtual_network_subnet_id` association for %s: %+v", *id, err)
					}
					var empty *string
					existing.SiteProperties.VirtualNetworkSubnetID = empty
				} else {
					existing.SiteProperties.VirtualNetworkSubnetID = pointer.To(subnetId)
				}
			}

			currentStack := ""
			sc := state.SiteConfig[0]

			if servicePlan.Sku != nil && servicePlan.Sku.Name != nil {
				if helpers.IsFreeOrSharedServicePlan(*servicePlan.Sku.Name) {
					if sc.AlwaysOn {
						return fmt.Errorf("always_on feature has to be turned off before switching to a free/shared Sku")
					}
				}
			}

			if len(sc.ApplicationStack) == 1 {
				currentStack = sc.ApplicationStack[0].CurrentStack
			}

			if metadata.ResourceData.HasChange("site_config") || servicePlanChange {
				existing.SiteConfig, err = sc.ExpandForUpdate(metadata, existing.SiteConfig, state.AppSettings)
				if err != nil {
					return err
				}
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

			updateFuture, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, existing)
			if err != nil {
				return fmt.Errorf("updating Windows %s: %+v", id, err)
			}
			if err := updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("wating to update %s: %+v", id, err)
			}

			siteMetadata := web.StringDictionary{Properties: map[string]*string{}}
			siteMetadata.Properties["CURRENT_STACK"] = pointer.To(currentStack)
			if _, err := client.UpdateMetadata(ctx, id.ResourceGroup, id.SiteName, siteMetadata); err != nil {
				return fmt.Errorf("setting Site Metadata for Current Stack on Windows %s: %+v", id, err)
			}

			// (@jackofallops) - App Settings can clobber logs configuration so must be updated before we send any Log updates
			if metadata.ResourceData.HasChange("app_settings") || metadata.ResourceData.HasChange("site_config.0.health_check_eviction_time_in_min") {
				appSettingsUpdate := helpers.ExpandAppSettingsForUpdate(state.AppSettings)
				appSettingsUpdate.Properties["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"] = pointer.To(strconv.Itoa(state.SiteConfig[0].HealthCheckEvictionTime))
				if _, err := client.UpdateApplicationSettings(ctx, id.ResourceGroup, id.SiteName, *appSettingsUpdate); err != nil {
					return fmt.Errorf("updating App Settings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("connection_string") {
				connectionStringUpdate := helpers.ExpandConnectionStrings(state.ConnectionStrings)
				if connectionStringUpdate.Properties == nil {
					connectionStringUpdate.Properties = map[string]*web.ConnStringValueTypePair{}
				}
				if _, err := client.UpdateConnectionStrings(ctx, id.ResourceGroup, id.SiteName, *connectionStringUpdate); err != nil {
					return fmt.Errorf("updating Connection Strings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("sticky_settings") {
				emptySlice := make([]string, 0)
				stickySettings := helpers.ExpandStickySettings(state.StickySettings)
				stickySettingsUpdate := web.SlotConfigNamesResource{
					SlotConfigNames: &web.SlotConfigNames{
						AppSettingNames:       &emptySlice,
						ConnectionStringNames: &emptySlice,
					},
				}

				if stickySettings != nil {
					if stickySettings.AppSettingNames != nil {
						stickySettingsUpdate.SlotConfigNames.AppSettingNames = stickySettings.AppSettingNames
					}
					if stickySettings.ConnectionStringNames != nil {
						stickySettingsUpdate.SlotConfigNames.ConnectionStringNames = stickySettings.ConnectionStringNames
					}
				}

				if _, err := client.UpdateSlotConfigurationNames(ctx, id.ResourceGroup, id.SiteName, stickySettingsUpdate); err != nil {
					return fmt.Errorf("updating Sticky Settings for Linux %s: %+v", id, err)
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
				if _, err := client.UpdateAuthSettings(ctx, id.ResourceGroup, id.SiteName, *authUpdate); err != nil {
					return fmt.Errorf("updating Auth Settings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings_v2") {
				authV2Update := helpers.ExpandAuthV2Settings(state.AuthV2Settings)
				if _, err := client.UpdateAuthSettingsV2(ctx, id.ResourceGroup, id.SiteName, *authV2Update); err != nil {
					return fmt.Errorf("updating AuthV2 Settings for Linux %s: %+v", id, err)
				}
				updateLogs = true
			}

			if metadata.ResourceData.HasChange("backup") {
				backupUpdate, err := helpers.ExpandBackupConfig(state.Backup)
				if err != nil {
					return fmt.Errorf("expanding backup configuration for Windows %s: %+v", *id, err)
				}

				if backupUpdate.BackupRequestProperties == nil {
					if _, err := client.DeleteBackupConfiguration(ctx, id.ResourceGroup, id.SiteName); err != nil {
						return fmt.Errorf("removing Backup Settings for Windows %s: %+v", id, err)
					}
				} else {
					if _, err := client.UpdateBackupConfiguration(ctx, id.ResourceGroup, id.SiteName, *backupUpdate); err != nil {
						return fmt.Errorf("updating Backup Settings for Windows %s: %+v", id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("logs") || updateLogs {
				logsUpdate := helpers.ExpandLogsConfig(state.LogsConfig)
				if logsUpdate.SiteLogsConfigProperties == nil {
					logsUpdate = helpers.DisabledLogsConfig() // The API is update only, so we need to send an update with everything switched of when a user removes the "logs" block
				}
				if _, err := client.UpdateDiagnosticLogsConfig(ctx, id.ResourceGroup, id.SiteName, *logsUpdate); err != nil {
					return fmt.Errorf("updating Logs Config for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("storage_account") {
				storageAccountUpdate := helpers.ExpandStorageConfig(state.StorageAccounts)
				if _, err := client.UpdateAzureStorageAccounts(ctx, id.ResourceGroup, id.SiteName, *storageAccountUpdate); err != nil {
					return fmt.Errorf("updating Storage Accounts for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("zip_deploy_file") {
				if err = helpers.GetCredentialsAndPublish(ctx, client, id.ResourceGroup, id.SiteName, state.ZipDeployFile); err != nil {
					return err
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

		id, err := parse.WebAppID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}
		site, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
		if err != nil || site.SiteProperties == nil {
			return fmt.Errorf("reading Windows %s: %+v", id, err)
		}
		props := site.SiteProperties
		if props.ServerFarmID == nil {
			return fmt.Errorf("determining Service Plan ID for Windows %s: %+v", id, err)
		}
		servicePlanId, err := parse.ServicePlanID(*props.ServerFarmID)
		if err != nil {
			return err
		}

		sp, err := servicePlanClient.Get(ctx, servicePlanId.ResourceGroup, servicePlanId.ServerfarmName)
		if err != nil || sp.Kind == nil {
			return fmt.Errorf("reading Service Plan for Windows %s: %+v", id, err)
		}
		if strings.Contains(*sp.Kind, "linux") || strings.Contains(*sp.Kind, "Linux") {
			return fmt.Errorf("specified Service Plan is not a Windows plan")
		}

		return nil
	}
}
