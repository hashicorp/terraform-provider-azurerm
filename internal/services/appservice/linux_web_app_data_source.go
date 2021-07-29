package appservice

import (
	"context"
	"fmt"
	"strings"
	"time"

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

type LinuxWebAppDataSource struct{}

var _ sdk.DataSource = LinuxWebAppDataSource{}

func (r LinuxWebAppDataSource) ModelObject() interface{} {
	return LinuxWebAppModel{}
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

		"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
	}
}

func (r LinuxWebAppDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": location.SchemaComputed(),

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

		"backup": backupSchemaComputed(),

		"client_affinity_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"client_cert_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"client_cert_mode": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"connection_string": connectionStringSchemaComputed(),

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

		"identity": helpers.IdentitySchemaComputed(),

		"kind": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"logs": logsConfigSchemaComputed(),

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

		"site_config": siteConfigSchemaLinuxComputed(),

		"storage_account": storageAccountSchemaComputed(),

		"tags": tags.SchemaDataSource(),
	}
}

func (r LinuxWebAppDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var webApp LinuxWebAppModel
			if err := metadata.Decode(&webApp); err != nil {
				return err
			}

			id := parse.NewWebAppID(subscriptionId, webApp.ResourceGroup, webApp.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("Linux Web App with %s not found", id)
				}
				return fmt.Errorf("checking for presence of existing Linux Web App with %s: %+v", id, err)
			}

			webAppSiteConfig, err := client.GetConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Site Config for Linux Web App %s: %+v", id, err)
			}

			auth, err := client.GetAuthSettings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Linux Web App %s: %+v", id, err)
			}

			backup, err := client.GetBackupConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				if !utils.ResponseWasNotFound(backup.Response) {
					return fmt.Errorf("reading Backup Settings for Linux Web App %s: %+v", id, err)
				}
			}

			logsConfig, err := client.GetDiagnosticLogsConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Diagnostic Logs information for Linux Web App %s: %+v", id, err)
			}

			appSettings, err := client.ListApplicationSettings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading App Settings for Linux Web App %s: %+v", id, err)
			}

			storageAccounts, err := client.ListAzureStorageAccounts(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Storage Account information for Linux Web App %s: %+v", id, err)
			}

			connectionStrings, err := client.ListConnectionStrings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Linux Web App %s: %+v", id, err)
			}

			siteCredentialsFuture, err := client.ListPublishingCredentials(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for Linux Web App %s: %+v", id, err)
			}

			if err := siteCredentialsFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for Site Publishing Credential information for Linux Web App %s: %+v", id, err)
			}
			siteCredentials, err := siteCredentialsFuture.Result(*client)
			if err != nil {
				return fmt.Errorf("reading Site Publishing Credential information for Linux Web App %s: %+v", id, err)
			}

			webApp.AppSettings = flattenAppSettings(appSettings)
			webApp.Kind = utils.NormalizeNilableString(existing.Kind)
			webApp.Location = location.NormalizeNilable(existing.Location)
			webApp.Tags = tags.ToTypedObject(existing.Tags)
			if props := existing.SiteProperties; props != nil {
				webApp.ClientAffinityEnabled = *props.ClientAffinityEnabled
				webApp.ClientCertEnabled = *props.ClientCertEnabled
				webApp.ClientCertMode = string(props.ClientCertMode)
				webApp.CustomDomainVerificationId = utils.NormalizeNilableString(props.CustomDomainVerificationID)
				webApp.DefaultHostname = utils.NormalizeNilableString(props.DefaultHostName)
				webApp.Enabled = *props.Enabled
				webApp.HttpsOnly = *props.HTTPSOnly
				webApp.ServicePlanId = utils.NormalizeNilableString(props.ServerFarmID)
				webApp.OutboundIPAddresses = utils.NormalizeNilableString(props.OutboundIPAddresses)
				webApp.OutboundIPAddressList = strings.Split(webApp.OutboundIPAddresses, ",")
				webApp.PossibleOutboundIPAddresses = utils.NormalizeNilableString(props.PossibleOutboundIPAddresses)
				webApp.PossibleOutboundIPAddressList = strings.Split(webApp.PossibleOutboundIPAddresses, ",")
			}

			webApp.AuthSettings = helpers.FlattenAuthSettings(auth)

			webApp.Backup = flattenBackupConfig(backup)

			webApp.Identity = helpers.FlattenIdentity(existing.Identity)

			webApp.LogsConfig = flattenLogsConfig(logsConfig)

			webApp.SiteConfig = flattenSiteConfigLinux(webAppSiteConfig.SiteConfig)

			webApp.StorageAccounts = flattenStorageAccounts(storageAccounts)

			webApp.ConnectionStrings = flattenConnectionStrings(connectionStrings)

			webApp.SiteCredentials = helpers.FlattenSiteCredentials(siteCredentials)

			metadata.SetID(id)

			return metadata.Encode(&webApp)
		},
	}
}
