package appservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LinuxWebAppDataSource struct{}

type LinuxWebAppDataSourceModel struct {
	Name                          string                     `tfschema:"name"`
	ResourceGroup                 string                     `tfschema:"resource_group_name"`
	Location                      string                     `tfschema:"location"`
	ServicePlanId                 string                     `tfschema:"service_plan_id"`
	AppSettings                   map[string]string          `tfschema:"app_settings"`
	AuthSettings                  []helpers.AuthSettings     `tfschema:"auth_settings"`
	Backup                        []helpers.Backup           `tfschema:"backup"`
	ClientAffinityEnabled         bool                       `tfschema:"client_affinity_enabled"`
	ClientCertEnabled             bool                       `tfschema:"client_certificate_enabled"`
	ClientCertMode                string                     `tfschema:"client_certificate_mode"`
	Enabled                       bool                       `tfschema:"enabled"`
	HttpsOnly                     bool                       `tfschema:"https_only"`
	KeyVaultReferenceIdentityID   string                     `tfschema:"key_vault_reference_identity_id"`
	LogsConfig                    []helpers.LogsConfig       `tfschema:"logs"`
	MetaData                      map[string]string          `tfschema:"app_metadata"`
	SiteConfig                    []helpers.SiteConfigLinux  `tfschema:"site_config"`
	StickySettings                []helpers.StickySettings   `tfschema:"sticky_settings"`
	StorageAccounts               []helpers.StorageAccount   `tfschema:"storage_account"`
	ConnectionStrings             []helpers.ConnectionString `tfschema:"connection_string"`
	Tags                          map[string]string          `tfschema:"tags"`
	CustomDomainVerificationId    string                     `tfschema:"custom_domain_verification_id"`
	DefaultHostname               string                     `tfschema:"default_hostname"`
	Kind                          string                     `tfschema:"kind"`
	OutboundIPAddresses           string                     `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string                   `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string                     `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string                   `tfschema:"possible_outbound_ip_address_list"`
	SiteCredentials               []helpers.SiteCredential   `tfschema:"site_credential"`
	VirtualNetworkSubnetID        string                     `tfschema:"virtual_network_subnet_id"`
}

var _ sdk.DataSource = LinuxWebAppDataSource{}

func (r LinuxWebAppDataSource) ModelObject() interface{} {
	return &LinuxWebAppDataSourceModel{}
}

func (r LinuxWebAppDataSource) ResourceType() string {
	return "azurerm_linux_web_app"
}

func (r LinuxWebAppDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.WebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r LinuxWebAppDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"app_metadata": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
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

		"key_vault_reference_identity_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

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

		"service_plan_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"site_config": helpers.SiteConfigSchemaLinuxComputed(),

		"storage_account": helpers.StorageAccountSchemaComputed(),

		"sticky_settings": helpers.StickySettingsComputedSchema(),

		"tags": tags.SchemaDataSource(),

		"virtual_network_subnet_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r LinuxWebAppDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var webApp LinuxWebAppDataSourceModel
			if err := metadata.Decode(&webApp); err != nil {
				return err
			}

			id := parse.NewWebAppID(subscriptionId, webApp.ResourceGroup, webApp.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("Linux %s not found", id)
				}
				return fmt.Errorf("retreiving Linux %s: %+v", id, err)
			}

			webAppSiteConfig, err := client.GetConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Site Config for Linux %s: %+v", id, err)
			}

			auth, err := client.GetAuthSettings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Linux %s: %+v", id, err)
			}

			backup, err := client.GetBackupConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				if !utils.ResponseWasNotFound(backup.Response) {
					return fmt.Errorf("reading Backup Settings for Linux %s: %+v", id, err)
				}
			}

			logsConfig, err := client.GetDiagnosticLogsConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Diagnostic Logs information for Linux %s: %+v", id, err)
			}

			appSettings, err := client.ListApplicationSettings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading App Settings for Linux %s: %+v", id, err)
			}

			storageAccounts, err := client.ListAzureStorageAccounts(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Storage Account information for Linux %s: %+v", id, err)
			}

			connectionStrings, err := client.ListConnectionStrings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Linux %s: %+v", id, err)
			}

			stickySettings, err := client.ListSlotConfigurationNames(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Sticky Settings for Linux %s: %+v", id, err)
			}

			siteCredentialsFuture, err := client.ListPublishingCredentials(ctx, id.ResourceGroup, id.SiteName)
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

			webApp.SiteConfig = helpers.FlattenSiteConfigLinux(webAppSiteConfig.SiteConfig, healthCheckCount)

			webApp.StorageAccounts = helpers.FlattenStorageAccounts(storageAccounts)

			webApp.ConnectionStrings = helpers.FlattenConnectionStrings(connectionStrings)

			webApp.StickySettings = helpers.FlattenStickySettings(stickySettings.SlotConfigNames)

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
