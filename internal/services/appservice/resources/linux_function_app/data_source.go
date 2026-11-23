// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package linux_function_app

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

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

	VirtualNetworkBackupRestoreEnabled bool     `tfschema:"virtual_network_backup_restore_enabled"`
	VirtualNetworkSubnetID             string   `tfschema:"virtual_network_subnet_id"`
	CustomDomainVerificationId         string   `tfschema:"custom_domain_verification_id"`
	DefaultHostname                    string   `tfschema:"default_hostname"`
	HostingEnvId                       string   `tfschema:"hosting_environment_id"`
	Kind                               string   `tfschema:"kind"`
	OutboundIPAddresses                string   `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList              []string `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses        string   `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList      []string `tfschema:"possible_outbound_ip_address_list"`
	Usage                              string   `tfschema:"usage"`

	SiteCredentials []helpers.SiteCredential `tfschema:"site_credential"`
}

func (d DataSource) ModelObject() any {
	return &LinuxFunctionAppDataSourceModel{}
}

func (d DataSource) ResourceType() string {
	return "azurerm_linux_function_app"
}

func (d DataSource) Description() string {
	return "Use this data source to access information about an existing Linux Function App."
}

func (d DataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateFunctionAppID
}

func (d DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Linux Function App.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.WebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"service_plan_id": {
			Description: "The ID of the App Service Plan within which this Function App has been created.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"storage_account_name": {
			Description: "The backend Storage Account name used by this Function App.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"storage_account_access_key": {
			Description: "The access key used to access the backend Storage Account for the Function App.",
			Type:        pluginsdk.TypeString,
			Sensitive:   true,
			Computed:    true,
		},

		"storage_uses_managed_identity": {
			Description: "Does the Function App use Managed Identity to access the Storage Account?",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"storage_key_vault_secret_id": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Key Vault Secret ID, including version, that contains the Connection String to connect to the Storage Account for this Function App.",
		},

		"app_settings": {
			Description: "A map of key-value pairs for [App Settings](https://learn.microsoft.com/azure/azure-functions/functions-app-settings) and custom values.",
			Type:        pluginsdk.TypeMap,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"availability": {
			Description: "The current availability state.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"auth_settings": helpers.AuthSettingsSchemaComputed(),

		"auth_settings_v2": helpers.AuthV2SettingsComputedSchema(),

		"backup": helpers.BackupSchemaComputed(),

		"builtin_logging_enabled": {
			Description: "Is built in logging enabled?",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"client_certificate_enabled": {
			Description: "Are Client Certificates enabled?",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"client_certificate_mode": {
			Description: "The mode of the Function App's client certificates requirement for incoming requests.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"client_certificate_exclusion_paths": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "Paths to exclude when using client certificates, separated by `;`.",
		},

		"connection_string": helpers.ConnectionStringSchemaComputed(),

		"daily_memory_time_quota": {
			Description: "The amount of memory in gigabyte-seconds that your application is allowed to consume per day.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"enabled": {
			Description: "Is the Function App enabled?",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"content_share_force_disabled": {
			Description: "Are the settings for linking the Function App to storage suppressed?",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"functions_extension_version": {
			Description: "The runtime version associated with the Function App.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"https_only": {
			Description: "Can the Function App only be accessed via HTTPS?",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

		"site_config": helpers.SiteConfigSchemaLinuxFunctionAppComputed(),

		"tags": commonschema.TagsDataSource(),

		"custom_domain_verification_id": {
			Description: "The identifier used by App Service to perform domain ownership verification via DNS TXT record.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"default_hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hosting_environment_id": {
			Description: "The ID of the App Service Environment used by Function App.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"kind": {
			Description: "The Kind value for this Linux Function App.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"outbound_ip_addresses": {
			Description: "A comma separated list of outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"outbound_ip_address_list": {
			Description: "A list of outbound IP addresses. For example `[\"52.23.25.3\", \"52.143.43.12\"]`.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"possible_outbound_ip_addresses": {
			Description: "A comma separated list of possible outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12,52.143.43.17`. This is a superset of `outbound_ip_addresses`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"possible_outbound_ip_address_list": {
			Description: "A list of possible outbound IP addresses, not all of which are necessarily in use. This is a superset of `outbound_ip_address_list`. For example `[\"52.23.25.3\", \"52.143.43.12\"]`.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"public_network_access_enabled": {
			Description: "Is Public Network Access enabled for this Linux Function App?",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"usage": {
			Description: "The current usage state.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"site_credential": helpers.SiteCredentialSchema(),

		"webdeploy_publish_basic_authentication_enabled": {
			Description: "Are the default WebDeploy Basic Authentication publishing credentials enabled?",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"ftp_publish_basic_authentication_enabled": {
			Description: "Are the default FTP Basic Authentication publishing credentials enabled?",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"sticky_settings": helpers.StickySettingsComputedSchema(),

		"virtual_network_backup_restore_enabled": {
			Description: "Whether backup and restore operations over the linked Virtual Network are enabled.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"virtual_network_subnet_id": {
			Description: "The subnet ID which the Linux Function App is VNet Integrated with.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func (d DataSource) Read() sdk.ResourceFunc {
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
					return fmt.Errorf("the Linux %s was not found", id)
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
					state.Availability = pointer.FromEnum(props.AvailabilityState)
					state.ServicePlanId = pointer.From(props.ServerFarmId)
					state.Enabled = pointer.From(props.Enabled)
					state.ClientCertMode = pointer.FromEnum(props.ClientCertMode)
					state.ClientCertExclusionPaths = pointer.From(props.ClientCertExclusionPaths)
					state.DailyMemoryTimeQuota = pointer.From(props.DailyMemoryTimeQuota)
					state.StickySettings = helpers.FlattenStickySettings(stickySettings.Model.Properties)
					state.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
					state.DefaultHostname = pointer.From(props.DefaultHostName)
					state.Usage = pointer.FromEnum(props.UsageState)
					state.PublicNetworkAccess = !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled)
					state.VirtualNetworkBackupRestoreEnabled = pointer.From(props.VnetBackupRestoreEnabled)

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

		// case "WEBSITES_ENABLE_APP_SERVICE_STORAGE": // TODO - Support this as a configurable bool, default `false` - Ref: https://learn.microsoft.com/azure/app-service/faq-app-service-linux#i-m-using-my-own-custom-container--i-want-the-platform-to-mount-an-smb-share-to-the---home---directory-

		case "APPINSIGHTS_INSTRUMENTATIONKEY":
			m.SiteConfig[0].AppInsightsInstrumentationKey = v

		case "APPLICATIONINSIGHTS_CONNECTION_STRING":
			m.SiteConfig[0].AppInsightsConnectionString = v

		case "AzureWebJobsStorage":
			if strings.HasPrefix(v, "@Microsoft.KeyVault") {
				m.StorageKeyVaultSecretID = strings.TrimPrefix(strings.TrimSuffix(v, ")"), "@Microsoft.KeyVault(")
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
