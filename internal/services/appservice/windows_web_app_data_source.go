package appservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WindowsWebAppDataSource struct{}

type WindowsWebAppDataSourceModel struct {
	Name                          string                      `tfschema:"name"`
	ResourceGroup                 string                      `tfschema:"resource_group_name"`
	Location                      string                      `tfschema:"location"`
	ServicePlanId                 string                      `tfschema:"service_plan_id"`
	AppSettings                   map[string]string           `tfschema:"app_settings"`
	AuthSettings                  []helpers.AuthSettings      `tfschema:"auth_settings"`
	Backup                        []helpers.Backup            `tfschema:"backup"`
	ClientAffinityEnabled         bool                        `tfschema:"client_affinity_enabled"`
	ClientCertEnabled             bool                        `tfschema:"client_certificate_enabled"`
	ClientCertMode                string                      `tfschema:"client_certificate_mode"`
	Enabled                       bool                        `tfschema:"enabled"`
	HttpsOnly                     bool                        `tfschema:"https_only"`
	LogsConfig                    []helpers.LogsConfig        `tfschema:"logs"`
	SiteConfig                    []helpers.SiteConfigWindows `tfschema:"site_config"`
	StorageAccounts               []helpers.StorageAccount    `tfschema:"storage_account"`
	ConnectionStrings             []helpers.ConnectionString  `tfschema:"connection_string"`
	CustomDomainVerificationId    string                      `tfschema:"custom_domain_verification_id"`
	DefaultHostname               string                      `tfschema:"default_hostname"`
	Kind                          string                      `tfschema:"kind"`
	OutboundIPAddresses           string                      `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string                    `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string                      `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string                    `tfschema:"possible_outbound_ip_address_list"`
	SiteCredentials               []helpers.SiteCredential    `tfschema:"site_credential"`
	Tags                          map[string]string           `tfschema:"tags"`
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

		"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
	}
}

func (d WindowsWebAppDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": location.SchemaComputed(),

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

		"site_config": helpers.SiteConfigSchemaWindowsComputed(),

		"storage_account": helpers.StorageAccountSchemaComputed(),

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

			id := parse.NewWebAppID(subscriptionId, webApp.ResourceGroup, webApp.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("Windows %s not found", id)
				}
				return fmt.Errorf("checking for presence of existing Windows %s: %+v", id, err)
			}

			webAppSiteConfig, err := client.GetConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Site Config for Windows %s: %+v", id, err)
			}

			auth, err := client.GetAuthSettings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Windows %s: %+v", id, err)
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

			var healthCheckCount *int
			webApp.AppSettings, healthCheckCount = helpers.FlattenAppSettings(appSettings)
			webApp.Kind = utils.NormalizeNilableString(existing.Kind)
			webApp.Location = location.NormalizeNilable(existing.Location)
			webApp.Tags = tags.ToTypedObject(existing.Tags)
			if props := existing.SiteProperties; props != nil {
				if props.ClientAffinityEnabled != nil {
					webApp.ClientAffinityEnabled = *props.ClientAffinityEnabled
				}
				if props.ClientCertEnabled != nil {
					webApp.ClientCertEnabled = *props.ClientCertEnabled
				}
				webApp.ClientCertMode = string(props.ClientCertMode)
				webApp.CustomDomainVerificationId = utils.NormalizeNilableString(props.CustomDomainVerificationID)
				webApp.DefaultHostname = utils.NormalizeNilableString(props.DefaultHostName)
				if props.Enabled != nil {
					webApp.Enabled = *props.Enabled
				}
				webApp.HttpsOnly = false
				if props.HTTPSOnly != nil {
					webApp.HttpsOnly = *props.HTTPSOnly
				}
				webApp.ServicePlanId = utils.NormalizeNilableString(props.ServerFarmID)
				webApp.OutboundIPAddresses = utils.NormalizeNilableString(props.OutboundIPAddresses)
				webApp.OutboundIPAddressList = strings.Split(webApp.OutboundIPAddresses, ",")
				webApp.PossibleOutboundIPAddresses = utils.NormalizeNilableString(props.PossibleOutboundIPAddresses)
				webApp.PossibleOutboundIPAddressList = strings.Split(webApp.PossibleOutboundIPAddresses, ",")
			}

			webApp.AuthSettings = helpers.FlattenAuthSettings(auth)

			webApp.Backup = helpers.FlattenBackupConfig(backup)

			webApp.LogsConfig = helpers.FlattenLogsConfig(logsConfig)

			currentStack := ""
			currentStackPtr, ok := siteMetadata.Properties["CURRENT_STACK"]
			if ok {
				currentStack = *currentStackPtr
			}
			webApp.SiteConfig = helpers.FlattenSiteConfigWindows(webAppSiteConfig.SiteConfig, currentStack, healthCheckCount)

			webApp.StorageAccounts = helpers.FlattenStorageAccounts(storageAccounts)

			webApp.ConnectionStrings = helpers.FlattenConnectionStrings(connectionStrings)

			webApp.SiteCredentials = helpers.FlattenSiteCredentials(siteCredentials)

			metadata.SetID(id)

			if err := metadata.Encode(&webApp); err != nil {
				return fmt.Errorf("encoding: %+v", err)
			}

			flattenedIdentity, err := flattenIdentity(existing.Identity)
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
