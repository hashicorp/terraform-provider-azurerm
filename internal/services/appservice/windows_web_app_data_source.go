// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
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

type WindowsWebAppDataSource struct{}

type WindowsWebAppDataSourceModel struct {
	Name                             string                                     `tfschema:"name"`
	ResourceGroup                    string                                     `tfschema:"resource_group_name"`
	Location                         string                                     `tfschema:"location"`
	ServicePlanId                    string                                     `tfschema:"service_plan_id"`
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
	LogsConfig                       []helpers.LogsConfig                       `tfschema:"logs"`
	PublicNetworkAccess              bool                                       `tfschema:"public_network_access_enabled"`
	PublishingDeployBasicAuthEnabled bool                                       `tfschema:"webdeploy_publish_basic_authentication_enabled"`
	PublishingFTPBasicAuthEnabled    bool                                       `tfschema:"ftp_publish_basic_authentication_enabled"`
	SiteConfig                       []helpers.SiteConfigWindows                `tfschema:"site_config"`
	StickySettings                   []helpers.StickySettings                   `tfschema:"sticky_settings"`
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
	Tags                             map[string]string                          `tfschema:"tags"`
	VirtualNetworkSubnetID           string                                     `tfschema:"virtual_network_subnet_id"`
}

var _ sdk.DataSource = WindowsWebAppDataSource{}

func (d WindowsWebAppDataSource) ModelObject() interface{} {
	return &WindowsWebAppDataSourceModel{}
}

func (d WindowsWebAppDataSource) ResourceType() string {
	return "azurerm_windows_web_app"
}

func (d WindowsWebAppDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.WebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d WindowsWebAppDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"service_plan_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
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

		"client_affinity_enabled": {
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

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"https_only": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

		"kind": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"logs": helpers.LogsConfigSchemaComputed(),

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

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"webdeploy_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"ftp_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"site_config": helpers.SiteConfigSchemaWindowsComputed(),

		"sticky_settings": helpers.StickySettingsComputedSchema(),

		"storage_account": helpers.StorageAccountSchemaComputed(),

		"virtual_network_subnet_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (d WindowsWebAppDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var webApp WindowsWebAppDataSourceModel
			if err := metadata.Decode(&webApp); err != nil {
				return err
			}

			baseId := commonids.NewAppServiceID(subscriptionId, webApp.ResourceGroup, webApp.Name)
			id, err := commonids.ParseWebAppID(baseId.ID())
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

			webAppSiteConfig, err := client.GetConfiguration(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Site Config for Windows %s: %+v", id, err)
			}

			auth, err := client.GetAuthSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Windows %s: %+v", id, err)
			}

			var authV2 webapps.SiteAuthSettingsV2
			if auth.Model != nil && auth.Model.Properties != nil && strings.EqualFold(pointer.From(auth.Model.Properties.ConfigVersion), "v2") {
				authV2Resp, err := client.GetAuthSettingsV2(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading authV2 settings for Linux %s: %+v", *id, err)
				}
				authV2 = *authV2Resp.Model
			}

			backup, err := client.GetBackupConfiguration(ctx, *id)
			if err != nil {
				if !response.WasNotFound(backup.HttpResponse) {
					return fmt.Errorf("reading Backup Settings for Windows %s: %+v", id, err)
				}
			}

			logsConfig, err := client.GetDiagnosticLogsConfiguration(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Diagnostic Logs information for Windows %s: %+v", id, err)
			}

			appSettings, err := client.ListApplicationSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading App Settings for Windows %s: %+v", id, err)
			}

			storageAccounts, err := client.ListAzureStorageAccounts(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Storage Account information for Windows %s: %+v", id, err)
			}

			connectionStrings, err := client.ListConnectionStrings(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Windows %s: %+v", id, err)
			}

			stickySettings, err := client.ListSlotConfigurationNames(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Sticky Settings for Linux %s: %+v", id, err)
			}

			siteCredentials, err := helpers.ListPublishingCredentials(ctx, client, *id)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for %s: %+v", *id, err)
			}

			siteMetadata, err := client.ListMetadata(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Site Metadata for Windows %s: %+v", id, err)
			}

			webApp.AppSettings = helpers.FlattenWebStringDictionary(appSettings.Model)

			if model := existing.Model; model != nil {
				webApp.Kind = pointer.From(model.Kind)
				webApp.Location = location.Normalize(model.Location)
				webApp.Tags = pointer.From(model.Tags)
				if props := model.Properties; props != nil {
					if props.ClientAffinityEnabled != nil {
						webApp.ClientAffinityEnabled = *props.ClientAffinityEnabled
					}
					if props.ClientCertEnabled != nil {
						webApp.ClientCertEnabled = *props.ClientCertEnabled
					}
					webApp.ClientCertMode = string(pointer.From(props.ClientCertMode))
					webApp.ClientCertExclusionPaths = pointer.From(props.ClientCertExclusionPaths)
					webApp.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
					webApp.DefaultHostname = pointer.From(props.DefaultHostName)
					if props.Enabled != nil {
						webApp.Enabled = *props.Enabled
					}
					webApp.HttpsOnly = false
					if props.HTTPSOnly != nil {
						webApp.HttpsOnly = *props.HTTPSOnly
					}
					servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(pointer.From(props.ServerFarmId))
					if err != nil {
						return fmt.Errorf("parsing Service Plan ID for %s, %+v", id, err)
					}
					webApp.ServicePlanId = servicePlanId.ID()
					webApp.OutboundIPAddresses = pointer.From(props.OutboundIPAddresses)
					webApp.OutboundIPAddressList = strings.Split(webApp.OutboundIPAddresses, ",")
					webApp.PossibleOutboundIPAddresses = pointer.From(props.PossibleOutboundIPAddresses)
					webApp.PossibleOutboundIPAddressList = strings.Split(webApp.PossibleOutboundIPAddresses, ",")
					if subnetId := pointer.From(props.VirtualNetworkSubnetId); subnetId != "" {
						webApp.VirtualNetworkSubnetID = subnetId
					}
					if hostingEnv := props.HostingEnvironmentProfile; hostingEnv != nil {
						webApp.HostingEnvId = pointer.From(hostingEnv.Id)
					}
					webApp.PublicNetworkAccess = !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled)
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

				webApp.PublishingFTPBasicAuthEnabled = basicAuthFTP
				webApp.PublishingDeployBasicAuthEnabled = basicAuthWebDeploy

				webApp.AuthSettings = helpers.FlattenAuthSettings(auth.Model)

				webApp.AuthV2Settings = helpers.FlattenAuthV2Settings(authV2)

				webApp.Backup = helpers.FlattenBackupConfig(backup.Model)

				webApp.LogsConfig = helpers.FlattenLogsConfig(logsConfig.Model)

				currentStack := ""

				if sm := siteMetadata.Model; sm != nil {
					p := *sm.Properties
					currentStack = p["CURRENT_STACK"]
				}

				siteConfig := helpers.SiteConfigWindows{}
				err = siteConfig.Flatten(webAppSiteConfig.Model.Properties, currentStack)
				if err != nil {
					return err
				}

				webApp.AppSettings = siteConfig.ParseNodeVersion(webApp.AppSettings)

				siteConfig.SetHealthCheckEvictionTime(webApp.AppSettings)
				if helpers.FxStringHasPrefix(siteConfig.WindowsFxVersion, helpers.FxStringPrefixDocker) {
					siteConfig.DecodeDockerAppStack(webApp.AppSettings)
				}

				webApp.StickySettings = helpers.FlattenStickySettings(stickySettings.Model.Properties)

				webApp.StorageAccounts = helpers.FlattenStorageAccounts(storageAccounts.Model)

				webApp.ConnectionStrings = helpers.FlattenConnectionStrings(connectionStrings.Model)

				webApp.SiteCredentials = helpers.FlattenSiteCredentials(siteCredentials)

				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				webApp.Identity = pointer.From(flattenedIdentity)

				if err = metadata.Encode(&webApp); err != nil {
					return fmt.Errorf("encoding: %+v", err)
				}

			}

			metadata.SetID(id)

			return nil
		},
	}
}
