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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LinuxFunctionAppDataSource struct{}

type LinuxFunctionAppDataSourceModel struct {
	Name               string `tfschema:"name"`
	ResourceGroup      string `tfschema:"resource_group_name"`
	Location           string `tfschema:"location"`
	ServicePlanId      string `tfschema:"service_plan_id"`
	StorageAccountName string `tfschema:"storage_account_name"`

	StorageAccountKey       string `tfschema:"storage_account_access_key"`
	StorageUsesMSI          bool   `tfschema:"storage_uses_managed_identity"` // Storage uses MSI not account key
	StorageKeyVaultSecretID string `tfschema:"storage_key_vault_secret_id"`

	AppSettings                      map[string]string                    `tfschema:"app_settings"`
	AuthSettings                     []helpers.AuthSettings               `tfschema:"auth_settings"`
	AuthV2Settings                   []helpers.AuthV2Settings             `tfschema:"auth_settings_v2"`
	Availability                     string                               `tfschema:"availability"`
	Backup                           []helpers.Backup                     `tfschema:"backup"` // Not supported on Dynamic or Basic plans
	BuiltinLogging                   bool                                 `tfschema:"builtin_logging_enabled"`
	ClientCertEnabled                bool                                 `tfschema:"client_certificate_enabled"`
	ClientCertMode                   string                               `tfschema:"client_certificate_mode"`
	ClientCertExclusionPaths         string                               `tfschema:"client_certificate_exclusion_paths"`
	ConnectionStrings                []helpers.ConnectionString           `tfschema:"connection_string"`
	DailyMemoryTimeQuota             int64                                `tfschema:"daily_memory_time_quota"`
	Enabled                          bool                                 `tfschema:"enabled"`
	FunctionExtensionsVersion        string                               `tfschema:"functions_extension_version"`
	ForceDisableContentShare         bool                                 `tfschema:"content_share_force_disabled"`
	HttpsOnly                        bool                                 `tfschema:"https_only"`
	PublicNetworkAccess              bool                                 `tfschema:"public_network_access_enabled"`
	PublishingDeployBasicAuthEnabled bool                                 `tfschema:"webdeploy_publish_basic_authentication_enabled"`
	PublishingFTPBasicAuthEnabled    bool                                 `tfschema:"ftp_publish_basic_authentication_enabled"`
	SiteConfig                       []helpers.SiteConfigLinuxFunctionApp `tfschema:"site_config"`
	StickySettings                   []helpers.StickySettings             `tfschema:"sticky_settings"`
	Tags                             map[string]string                    `tfschema:"tags"`

	VirtualNetworkSubnetID        string   `tfschema:"virtual_network_subnet_id"`
	CustomDomainVerificationId    string   `tfschema:"custom_domain_verification_id"`
	DefaultHostname               string   `tfschema:"default_hostname"`
	HostingEnvId                  string   `tfschema:"hosting_environment_id"`
	Kind                          string   `tfschema:"kind"`
	OutboundIPAddresses           string   `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string   `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string `tfschema:"possible_outbound_ip_address_list"`
	Usage                         string   `tfschema:"usage"`

	SiteCredentials []helpers.SiteCredential `tfschema:"site_credential"`
}

func (d LinuxFunctionAppDataSource) ModelObject() interface{} {
	return &LinuxFunctionAppDataSourceModel{}
}

func (d LinuxFunctionAppDataSource) ResourceType() string {
	return "azurerm_linux_function_app"
}

func (d LinuxFunctionAppDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateFunctionAppID
}

func (d LinuxFunctionAppDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.WebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d LinuxFunctionAppDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"service_plan_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"storage_account_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"storage_account_access_key": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},

		"storage_uses_managed_identity": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"storage_key_vault_secret_id": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Key Vault Secret ID, including version, that contains the Connection String used to connect to the storage account for this Function App.",
		},

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"availability": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"auth_settings": helpers.AuthSettingsSchemaComputed(),

		"auth_settings_v2": helpers.AuthV2SettingsComputedSchema(),

		"backup": helpers.BackupSchemaComputed(),

		"builtin_logging_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"client_certificate_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"client_certificate_mode": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"client_certificate_exclusion_paths": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "Paths to exclude when using client certificates, separated by ;",
		},

		"connection_string": helpers.ConnectionStringSchemaComputed(),

		"daily_memory_time_quota": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"content_share_force_disabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"functions_extension_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"https_only": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

		"site_config": helpers.SiteConfigSchemaLinuxFunctionAppComputed(),

		"tags": tags.SchemaDataSource(),

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

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"usage": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"site_credential": helpers.SiteCredentialSchema(),

		"webdeploy_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"ftp_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"sticky_settings": helpers.StickySettingsComputedSchema(),

		"virtual_network_subnet_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (d LinuxFunctionAppDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 25 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var linuxFunctionApp LinuxFunctionAppDataSourceModel
			if err := metadata.Decode(&linuxFunctionApp); err != nil {
				return err
			}

			id := commonids.NewAppServiceID(subscriptionId, linuxFunctionApp.ResourceGroup, linuxFunctionApp.Name)

			functionApp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(functionApp.HttpResponse) {
					return fmt.Errorf("Linux %s not found", id)
				}
				return fmt.Errorf("reading Linux %s: %+v", id, err)
			}

			appSettingsResp, err := client.ListApplicationSettings(ctx, id)
			if err != nil {
				return fmt.Errorf("reading App Settings for Linux %s: %+v", id, err)
			}

			connectionStrings, err := client.ListConnectionStrings(ctx, id)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Linux %s: %+v", id, err)
			}

			stickySettings, err := client.ListSlotConfigurationNames(ctx, id)
			if err != nil {
				return fmt.Errorf("reading Sticky Settings for Linux %s: %+v", id, err)
			}

			siteCredentials, err := helpers.ListPublishingCredentials(ctx, client, id)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for %s: %+v", id, err)
			}

			auth, err := client.GetAuthSettings(ctx, id)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Linux %s: %+v", id, err)
			}

			var authV2 webapps.SiteAuthSettingsV2
			authV2Resp, err := client.GetAuthSettingsV2(ctx, id)
			if err != nil {
				return fmt.Errorf("reading authV2 settings for Linux %s: %+v", id, err)
			}
			authV2 = *authV2Resp.Model

			backup, err := client.GetBackupConfiguration(ctx, id)
			if err != nil {
				if !response.WasNotFound(backup.HttpResponse) {
					return fmt.Errorf("reading Backup Settings for Linux %s: %+v", id, err)
				}
			}

			logs, err := client.GetDiagnosticLogsConfiguration(ctx, id)
			if err != nil {
				return fmt.Errorf("reading logs configuration for Linux %s: %+v", id, err)
			}

			basicAuthFTP := true
			if basicAuthFTPResp, err := client.GetFtpAllowed(ctx, id); err != nil && basicAuthFTPResp.Model != nil {
				return fmt.Errorf("retrieving state of FTP Basic Auth for %s: %+v", id, err)
			} else if csmProps := basicAuthFTPResp.Model.Properties; csmProps != nil {
				basicAuthFTP = csmProps.Allow
			}

			basicAuthWebDeploy := true
			if basicAuthWebDeployResp, err := client.GetScmAllowed(ctx, id); err != nil && basicAuthWebDeployResp.Model != nil {
				return fmt.Errorf("retrieving state of WebDeploy Basic Auth for %s: %+v", id, err)
			} else if csmProps := basicAuthWebDeployResp.Model.Properties; csmProps != nil {
				basicAuthWebDeploy = csmProps.Allow
			}

			state := LinuxFunctionAppDataSourceModel{
				Name:                             id.SiteName,
				ResourceGroup:                    id.ResourceGroupName,
				PublishingFTPBasicAuthEnabled:    basicAuthFTP,
				PublishingDeployBasicAuthEnabled: basicAuthWebDeploy,
				ConnectionStrings:                helpers.FlattenConnectionStrings(connectionStrings.Model),
				SiteCredentials:                  helpers.FlattenSiteCredentials(siteCredentials),
				AuthSettings:                     helpers.FlattenAuthSettings(auth.Model),
				AuthV2Settings:                   helpers.FlattenAuthV2Settings(authV2),
				Backup:                           helpers.FlattenBackupConfig(backup.Model),
			}

			if model := functionApp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.Kind = pointer.From(model.Kind)

				if props := model.Properties; props != nil {
					state.Availability = string(pointer.From(props.AvailabilityState))
					state.ServicePlanId = pointer.From(props.ServerFarmId)
					state.Enabled = pointer.From(props.Enabled)
					state.ClientCertMode = string(pointer.From(props.ClientCertMode))
					state.ClientCertExclusionPaths = pointer.From(props.ClientCertExclusionPaths)
					state.DailyMemoryTimeQuota = pointer.From(props.DailyMemoryTimeQuota)
					state.StickySettings = helpers.FlattenStickySettings(stickySettings.Model.Properties)
					state.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
					state.DefaultHostname = pointer.From(props.DefaultHostName)
					state.Usage = string(pointer.From(props.UsageState))
					state.PublicNetworkAccess = !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled)

					if hostingEnv := props.HostingEnvironmentProfile; hostingEnv != nil {
						state.HostingEnvId = pointer.From(hostingEnv.Id)
					}

					if v := props.OutboundIPAddresses; v != nil {
						state.OutboundIPAddresses = *v
						state.OutboundIPAddressList = strings.Split(*v, ",")
					}

					if v := props.PossibleOutboundIPAddresses; v != nil {
						state.PossibleOutboundIPAddresses = *v
						state.PossibleOutboundIPAddressList = strings.Split(*v, ",")
					}

					state.HttpsOnly = pointer.From(props.HTTPSOnly)
					state.ClientCertEnabled = pointer.From(props.ClientCertEnabled)
					state.VirtualNetworkSubnetID = pointer.From(props.VirtualNetworkSubnetId)
				}

				configResp, err := client.GetConfiguration(ctx, id)
				if err != nil || configResp.Model == nil {
					return fmt.Errorf("making Read request on AzureRM Function App Configuration %q: %+v", id.SiteName, err)
				}

				siteConfig, err := helpers.FlattenSiteConfigLinuxFunctionApp(configResp.Model.Properties)
				if err != nil {
					return fmt.Errorf("reading Site Config for Linux %s: %+v", id, err)
				}
				state.SiteConfig = []helpers.SiteConfigLinuxFunctionApp{*siteConfig}

				state.unpackLinuxFunctionAppSettings(appSettingsResp.Model, metadata)

				state.SiteConfig[0].AppServiceLogs = helpers.FlattenFunctionAppAppServiceLogs(logs.Model)

				metadata.SetID(id)

				if err := metadata.Encode(&state); err != nil {
					return fmt.Errorf("encoding: %+v", err)
				}

				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				if err := metadata.ResourceData.Set("identity", flattenedIdentity); err != nil {
					return fmt.Errorf("setting `identity`: %+v", err)
				}
			}

			return nil
		},
	}
}

func (m *LinuxFunctionAppDataSourceModel) unpackLinuxFunctionAppSettings(input *webapps.StringDictionary, metadata sdk.ResourceMetaData) {
	if input == nil || input.Properties == nil {
		return
	}

	appSettings := make(map[string]string)
	var dockerSettings helpers.ApplicationStackDocker
	m.BuiltinLogging = false

	for k, v := range *input.Properties {
		switch k {
		case "FUNCTIONS_EXTENSION_VERSION":
			m.FunctionExtensionsVersion = v

		case "WEBSITE_NODE_DEFAULT_VERSION": // Note - This is no longer used in Linux Apps
		case "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING":
			if _, ok := metadata.ResourceData.GetOk("app_settings.WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"); ok {
				appSettings[k] = v
			}
		case "WEBSITE_CONTENTSHARE":
			if _, ok := metadata.ResourceData.GetOk("app_settings.WEBSITE_CONTENTSHARE"); ok {
				appSettings[k] = v
			}
		case "WEBSITE_HTTPLOGGING_RETENTION_DAYS":
		case "FUNCTIONS_WORKER_RUNTIME":
			if len(m.SiteConfig[0].ApplicationStack) > 0 {
				m.SiteConfig[0].ApplicationStack[0].CustomHandler = strings.EqualFold(v, "custom")
			}

		case "DOCKER_REGISTRY_SERVER_URL":
			dockerSettings.RegistryURL = v

		case "DOCKER_REGISTRY_SERVER_USERNAME":
			dockerSettings.RegistryUsername = v

		case "DOCKER_REGISTRY_SERVER_PASSWORD":
			dockerSettings.RegistryPassword = v

		// case "WEBSITES_ENABLE_APP_SERVICE_STORAGE": // TODO - Support this as a configurable bool, default `false` - Ref: https://docs.microsoft.com/en-us/azure/app-service/faq-app-service-linux#i-m-using-my-own-custom-container--i-want-the-platform-to-mount-an-smb-share-to-the---home---directory-

		case "APPINSIGHTS_INSTRUMENTATIONKEY":
			m.SiteConfig[0].AppInsightsInstrumentationKey = v

		case "APPLICATIONINSIGHTS_CONNECTION_STRING":
			m.SiteConfig[0].AppInsightsConnectionString = v

		case "AzureWebJobsStorage":
			if strings.HasPrefix(v, "@Microsoft.KeyVault") {
				trimmed := strings.TrimPrefix(strings.TrimSuffix(v, ")"), "@Microsoft.KeyVault(")
				m.StorageKeyVaultSecretID = trimmed
			} else {
				m.StorageAccountName, m.StorageAccountKey = helpers.ParseWebJobsStorageString(v)
			}

		case "AzureWebJobsDashboard":
			m.BuiltinLogging = true

		case "WEBSITE_HEALTHCHECK_MAXPINGFAILURES":
			i, _ := strconv.Atoi(v)
			m.SiteConfig[0].HealthCheckEvictionTime = int64(i)

		default:
			appSettings[k] = v
		}
	}

	if dockerSettings.RegistryURL != "" {
		appStack := make([]helpers.ApplicationStackLinuxFunctionApp, 0)
		docker, _ := helpers.DecodeFunctionAppDockerFxString(m.SiteConfig[0].LinuxFxVersion, dockerSettings)
		appStack = append(appStack, helpers.ApplicationStackLinuxFunctionApp{Docker: docker})
		m.SiteConfig[0].ApplicationStack = appStack
	}

	m.AppSettings = appSettings
}
