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
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WindowsFunctionAppDataSource struct{}

type WindowsFunctionAppDataSourceModel struct {
	Name               string `tfschema:"name"`
	ResourceGroup      string `tfschema:"resource_group_name"`
	Location           string `tfschema:"location"`
	ServicePlanId      string `tfschema:"service_plan_id"`
	StorageAccountName string `tfschema:"storage_account_name"`

	StorageAccountKey       string `tfschema:"storage_account_access_key"`
	StorageUsesMSI          bool   `tfschema:"storage_uses_managed_identity"`
	StorageKeyVaultSecretID string `tfschema:"storage_key_vault_secret_id"`

	AppSettings                      map[string]string                      `tfschema:"app_settings"`
	AuthSettings                     []helpers.AuthSettings                 `tfschema:"auth_settings"`
	AuthV2Settings                   []helpers.AuthV2Settings               `tfschema:"auth_settings_v2"`
	Backup                           []helpers.Backup                       `tfschema:"backup"`
	BuiltinLogging                   bool                                   `tfschema:"builtin_logging_enabled"`
	ClientCertEnabled                bool                                   `tfschema:"client_certificate_enabled"`
	ClientCertMode                   string                                 `tfschema:"client_certificate_mode"`
	ClientCertExclusionPaths         string                                 `tfschema:"client_certificate_exclusion_paths"`
	ConnectionStrings                []helpers.ConnectionString             `tfschema:"connection_string"`
	DailyMemoryTimeQuota             int64                                  `tfschema:"daily_memory_time_quota"`
	Enabled                          bool                                   `tfschema:"enabled"`
	FunctionExtensionsVersion        string                                 `tfschema:"functions_extension_version"`
	ForceDisableContentShare         bool                                   `tfschema:"content_share_force_disabled"`
	HttpsOnly                        bool                                   `tfschema:"https_only"`
	PublicNetworkAccess              bool                                   `tfschema:"public_network_access_enabled"`
	PublishingDeployBasicAuthEnabled bool                                   `tfschema:"webdeploy_publish_basic_authentication_enabled"`
	PublishingFTPBasicAuthEnabled    bool                                   `tfschema:"ftp_publish_basic_authentication_enabled"`
	SiteConfig                       []helpers.SiteConfigWindowsFunctionApp `tfschema:"site_config"`
	StickySettings                   []helpers.StickySettings               `tfschema:"sticky_settings"`
	Tags                             map[string]string                      `tfschema:"tags"`
	VirtualNetworkSubnetId           string                                 `tfschema:"virtual_network_subnet_id"`

	CustomDomainVerificationId    string   `tfschema:"custom_domain_verification_id"`
	DefaultHostname               string   `tfschema:"default_hostname"`
	HostingEnvId                  string   `tfschema:"hosting_environment_id"`
	Kind                          string   `tfschema:"kind"`
	OutboundIPAddresses           string   `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string   `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string `tfschema:"possible_outbound_ip_address_list"`

	SiteCredentials []helpers.SiteCredential `tfschema:"site_credential"`

	Identity []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
}

var _ sdk.DataSource = WindowsFunctionAppDataSource{}

func (d WindowsFunctionAppDataSource) ModelObject() interface{} {
	return &WindowsFunctionAppDataSourceModel{}
}

func (d WindowsFunctionAppDataSource) ResourceType() string {
	return "azurerm_windows_function_app"
}

func (d WindowsFunctionAppDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.WebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d WindowsFunctionAppDataSource) Attributes() map[string]*pluginsdk.Schema {
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
			Type:     pluginsdk.TypeString,
			Computed: true,
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

		"site_credential": helpers.SiteCredentialSchema(),

		"webdeploy_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"ftp_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"site_config": helpers.SiteConfigSchemaWindowsFunctionAppComputed(),

		"sticky_settings": helpers.StickySettingsComputedSchema(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

		"tags": tags.SchemaDataSource(),

		"virtual_network_subnet_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (d WindowsFunctionAppDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var functionApp WindowsFunctionAppDataSourceModel
			if err := metadata.Decode(&functionApp); err != nil {
				return err
			}

			baseID := commonids.NewAppServiceID(subscriptionId, functionApp.ResourceGroup, functionApp.Name)
			id, err := commonids.ParseFunctionAppID(baseID.ID())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("Windows %s not found", id)
				}
				return fmt.Errorf("checking for presence of existing Windows %s: %+v", id, err)
			}

			if existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("reading properties of Windows %s", id)
			}

			functionApp.Name = id.SiteName
			functionApp.ResourceGroup = id.ResourceGroupName
			if model := existing.Model; model != nil {
				functionApp.Location = location.Normalize(model.Location)
				functionApp.Tags = pointer.From(model.Tags)
				functionApp.Kind = pointer.From(model.Kind)

				if props := model.Properties; props != nil {
					servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(pointer.From(props.ServerFarmId))
					if err != nil {
						return fmt.Errorf("reading Service Plan Id for %s: %+v", id, err)
					}
					functionApp.ServicePlanId = servicePlanId.ID()
					functionApp.Enabled = utils.NormaliseNilableBool(props.Enabled)
					functionApp.ClientCertMode = string(pointer.From(props.ClientCertMode))
					functionApp.ClientCertExclusionPaths = pointer.From(props.ClientCertExclusionPaths)
					functionApp.DailyMemoryTimeQuota = pointer.From(props.DailyMemoryTimeQuota)
					functionApp.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
					functionApp.DefaultHostname = pointer.From(props.DefaultHostName)
					functionApp.VirtualNetworkSubnetId = pointer.From(props.VirtualNetworkSubnetId)
					functionApp.PublicNetworkAccess = !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled)

					if hostingEnv := props.HostingEnvironmentProfile; hostingEnv != nil {
						functionApp.HostingEnvId = pointer.From(hostingEnv.Id)
					}

					if v := props.OutboundIPAddresses; v != nil {
						functionApp.OutboundIPAddresses = *v
						functionApp.OutboundIPAddressList = strings.Split(*v, ",")
					}

					if v := props.PossibleOutboundIPAddresses; v != nil {
						functionApp.PossibleOutboundIPAddresses = *v
						functionApp.PossibleOutboundIPAddressList = strings.Split(*v, ",")
					}

					functionApp.HttpsOnly = pointer.From(props.HTTPSOnly)
					functionApp.ClientCertEnabled = pointer.From(props.ClientCertEnabled)
				}

				basicAuthFTP := true
				if basicAuthFTPResp, err := client.GetFtpAllowed(ctx, *id); err != nil || basicAuthFTPResp.Model.Properties == nil {
					return fmt.Errorf("retrieving state of FTP Basic Auth for %s: %+v", id, err)
				} else if csmProps := basicAuthFTPResp.Model.Properties; csmProps != nil {
					basicAuthFTP = csmProps.Allow
				}

				basicAuthWebDeploy := true
				if basicAuthWebDeployResp, err := client.GetScmAllowed(ctx, *id); err != nil || basicAuthWebDeployResp.Model.Properties == nil {
					return fmt.Errorf("retrieving state of WebDeploy Basic Auth for %s: %+v", id, err)
				} else if csmProps := basicAuthWebDeployResp.Model.Properties; csmProps != nil {
					basicAuthWebDeploy = csmProps.Allow
				}

				functionApp.PublishingFTPBasicAuthEnabled = basicAuthFTP
				functionApp.PublishingDeployBasicAuthEnabled = basicAuthWebDeploy

				appSettingsResp, err := client.ListApplicationSettings(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading App Settings for Windows %s: %+v", id, err)
				}

				connectionStrings, err := client.ListConnectionStrings(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading Connection String information for Windows %s: %+v", id, err)
				}

				stickySettings, err := client.ListSlotConfigurationNames(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading Sticky Settings for Windows %s: %+v", id, err)
				}

				siteCredentials, err := helpers.ListPublishingCredentials(ctx, client, *id)
				if err != nil {
					return fmt.Errorf("listing Site Publishing Credential information for %s: %+v", id, err)
				}

				auth, err := client.GetAuthSettings(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading Auth Settings for Windows %s: %+v", id, err)
				}

				var authV2 webapps.SiteAuthSettingsV2
				authV2Resp, err := client.GetAuthSettingsV2(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading authV2 settings for Linux %s: %+v", id, err)
				}
				authV2 = *authV2Resp.Model

				backup, err := client.GetBackupConfiguration(ctx, *id)
				if err != nil {
					if !response.WasNotFound(backup.HttpResponse) {
						return fmt.Errorf("reading Backup Settings for Windows %s: %+v", id, err)
					}
				}

				logs, err := client.GetDiagnosticLogsConfiguration(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading logs configuration for Windows %s: %+v", id, err)
				}

				configResp, err := client.GetConfiguration(ctx, *id)
				if err != nil || configResp.Model == nil {
					return fmt.Errorf("making Read request on AzureRM Function App Configuration %q: %+v", id.SiteName, err)
				}

				siteConfig, err := helpers.FlattenSiteConfigWindowsFunctionApp(configResp.Model.Properties)
				if err != nil {
					return fmt.Errorf("reading Site Config for Windows %s: %+v", id, err)
				}

				functionApp.SiteConfig = []helpers.SiteConfigWindowsFunctionApp{*siteConfig}

				functionApp.unpackWindowsFunctionAppSettings(appSettingsResp.Model)

				functionApp.ConnectionStrings = helpers.FlattenConnectionStrings(connectionStrings.Model)

				functionApp.SiteCredentials = helpers.FlattenSiteCredentials(siteCredentials)

				functionApp.AuthSettings = helpers.FlattenAuthSettings(auth.Model)

				functionApp.AuthV2Settings = helpers.FlattenAuthV2Settings(authV2)

				functionApp.Backup = helpers.FlattenBackupConfig(backup.Model)

				functionApp.SiteConfig[0].AppServiceLogs = helpers.FlattenFunctionAppAppServiceLogs(logs.Model)

				functionApp.StickySettings = helpers.FlattenStickySettings(stickySettings.Model.Properties)

				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				functionApp.Identity = pointer.From(flattenedIdentity)

				metadata.SetID(id)

				if err := metadata.Encode(&functionApp); err != nil {
					return fmt.Errorf("encoding: %+v", err)
				}
			}

			return nil
		},
	}
}

func (m *WindowsFunctionAppDataSourceModel) unpackWindowsFunctionAppSettings(input *webapps.StringDictionary) {
	if input == nil || input.Properties == nil {
		return
	}

	appSettings := make(map[string]string)
	var dockerSettings helpers.ApplicationStackDocker
	m.BuiltinLogging = false

	for k, v := range *input.Properties {
		switch k {
		case "FUNCTIONS_EXTENSION_VERSION":
			m.FunctionExtensionsVersion = (v)

		case "WEBSITE_NODE_DEFAULT_VERSION": // Note - This is only set if it's not the default of 12, but we collect it from WindowsFxVersion so can discard it here
		case "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING":
		case "WEBSITE_CONTENTSHARE":
		case "WEBSITE_HTTPLOGGING_RETENTION_DAYS":
		case "FUNCTIONS_WORKER_RUNTIME":
			if len(m.SiteConfig) > 0 && len(m.SiteConfig[0].ApplicationStack) > 0 {
				m.SiteConfig[0].ApplicationStack[0].CustomHandler = strings.EqualFold(v, "custom")
			}

		case "DOCKER_REGISTRY_SERVER_URL":
			dockerSettings.RegistryURL = v

		case "DOCKER_REGISTRY_SERVER_USERNAME":
			dockerSettings.RegistryUsername = v

		case "DOCKER_REGISTRY_SERVER_PASSWORD":
			dockerSettings.RegistryPassword = v

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

	m.AppSettings = appSettings
}
