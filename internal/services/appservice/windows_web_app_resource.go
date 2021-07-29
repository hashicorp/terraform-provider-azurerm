package appservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-01-15/web"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WindowsWebAppResource struct{}

type WindowsWebAppModel struct {
	Name                          string                   `tfschema:"name"`
	ResourceGroup                 string                   `tfschema:"resource_group_name"`
	Location                      string                   `tfschema:"location"`
	ServicePlanId                 string                   `tfschema:"service_plan_id"`
	AppSettings                   map[string]string        `tfschema:"app_settings"`
	AuthSettings                  []helpers.AuthSettings   `tfschema:"auth_settings"`
	Backup                        []Backup                 `tfschema:"backup"`
	ClientAffinityEnabled         bool                     `tfschema:"client_affinity_enabled"`
	ClientCertEnabled             bool                     `tfschema:"client_cert_enabled"`
	ClientCertMode                string                   `tfschema:"client_cert_mode"`
	Enabled                       bool                     `tfschema:"enabled"`
	HttpsOnly                     bool                     `tfschema:"https_only"`
	Identity                      []helpers.Identity       `tfschema:"identity"`
	LogsConfig                    []LogsConfig             `tfschema:"logs"`
	SiteConfig                    []SiteConfigWindows      `tfschema:"site_config"`
	StorageAccounts               []StorageAccount         `tfschema:"storage_account"`
	ConnectionStrings             []ConnectionString       `tfschema:"connection_string"`
	CustomDomainVerificationId    string                   `tfschema:"custom_domain_verification_id"`
	DefaultHostname               string                   `tfschema:"default_hostname"`
	Kind                          string                   `tfschema:"kind"`
	OutboundIPAddresses           string                   `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string                 `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string                   `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string                 `tfschema:"possible_outbound_ip_address_list"`
	SiteCredentials               []helpers.SiteCredential `tfschema:"site_credential"`
	Tags                          map[string]string        `tfschema:"tags"`
}

var _ sdk.Resource = WindowsWebAppResource{}
var _ sdk.ResourceWithUpdate = WindowsWebAppResource{}

// TODO - Feature: Deployments (Preview)?
// TODO - Feature: App Insights?

func (r WindowsWebAppResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": location.Schema(),

		"service_plan_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ServicePlanID,
		},

		// Optional

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"auth_settings": helpers.AuthSettingsSchema(),

		"backup": backupSchema(),

		"client_affinity_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"client_cert_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"client_cert_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "Required",
			ValidateFunc: validation.StringInSlice([]string{
				string(web.ClientCertModeOptional),
				string(web.ClientCertModeRequired),
			}, false),
		},

		"connection_string": connectionStringSchema(),

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

		"identity": helpers.IdentitySchema(),

		"logs": logsConfigSchema(),

		"site_config": siteConfigSchemaWindows(),

		"storage_account": storageAccountSchema(),

		"tags": tags.Schema(),
	}
}

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
	return WindowsWebAppModel{}
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
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewWebAppID(subscriptionId, webApp.ResourceGroup, webApp.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Windows Web App with %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			availabilityRequest := web.ResourceNameAvailabilityRequest{
				Name: utils.String(webApp.Name),
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
			// TODO - Does this change for Private Link?
			if servicePlan.HostingEnvironmentProfile != nil {
				// TODO - Check this for Gov / Sovereign cloud
				availabilityRequest.Name = utils.String(fmt.Sprintf("%s.%s.appserviceenvironment.net", webApp.Name, servicePlanId.ServerfarmName))
				availabilityRequest.IsFqdn = utils.Bool(true)
			}

			checkName, err := client.CheckNameAvailability(ctx, availabilityRequest)
			if err != nil {
				return fmt.Errorf("checking name availability for %s: %+v", id, err)
			}
			if !*checkName.NameAvailable {
				return fmt.Errorf("the Site Name %q failed the availability check: %+v", id.SiteName, *checkName.Message)
			}

			siteConfig, currentStack, err := expandSiteConfigWindows(webApp.SiteConfig)
			if err != nil {
				return err
			}

			siteEnvelope := web.Site{
				Location: utils.String(webApp.Location),
				Tags:     tags.FromTypedObject(webApp.Tags),
				SiteProperties: &web.SiteProperties{
					ServerFarmID:          utils.String(webApp.ServicePlanId),
					Enabled:               utils.Bool(webApp.Enabled),
					HTTPSOnly:             utils.Bool(webApp.HttpsOnly),
					SiteConfig:            siteConfig,
					ClientAffinityEnabled: utils.Bool(webApp.ClientAffinityEnabled),
					ClientCertEnabled:     utils.Bool(webApp.ClientCertEnabled),
					ClientCertMode:        web.ClientCertMode(webApp.ClientCertMode),
				},
			}

			if identity := helpers.ExpandIdentity(webApp.Identity); identity != nil {
				siteEnvelope.Identity = identity
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, siteEnvelope)
			if err != nil {
				return fmt.Errorf("creating Windows Web App %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of Windows Web App %s: %+v", id, err)
			}

			if currentStack != nil && *currentStack != "" {
				siteMetadata := web.StringDictionary{Properties: map[string]*string{}}
				siteMetadata.Properties["CURRENT_STACK"] = currentStack
				if _, err := client.UpdateMetadata(ctx, id.ResourceGroup, id.SiteName, siteMetadata); err != nil {
					return fmt.Errorf("setting Site Metadata for Current Stack on Windows Web App %s: %+v", id, err)
				}
			}

			metadata.SetID(id)

			appSettings := expandAppSettings(webApp.AppSettings)
			if appSettings != nil {
				if _, err := client.UpdateApplicationSettings(ctx, id.ResourceGroup, id.SiteName, *appSettings); err != nil {
					return fmt.Errorf("setting App Settings for Windows Web App %s: %+v", id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(webApp.AuthSettings)
			if auth != nil {
				if _, err := client.UpdateAuthSettings(ctx, id.ResourceGroup, id.SiteName, *auth); err != nil {
					return fmt.Errorf("setting Authorisation Settings for %s: %+v", id, err)
				}
			}

			logsConfig := expandLogsConfig(webApp.LogsConfig)
			if logsConfig != nil {
				if _, err := client.UpdateDiagnosticLogsConfig(ctx, id.ResourceGroup, id.SiteName, *logsConfig); err != nil {
					return fmt.Errorf("setting Diagnostic Logs Configuration for Windows Web App %s: %+v", id, err)
				}
			}

			backupConfig := expandBackupConfig(webApp.Backup)
			if backupConfig != nil {
				if _, err := client.UpdateBackupConfiguration(ctx, id.ResourceGroup, id.SiteName, *backupConfig); err != nil {
					return fmt.Errorf("adding Backup Settings for Windows Web App %s: %+v", id, err)
				}
			}

			storageConfig := expandStorageConfig(webApp.StorageAccounts)
			if storageConfig != nil {
				if _, err := client.UpdateAzureStorageAccounts(ctx, id.ResourceGroup, id.SiteName, *storageConfig); err != nil {
					if err != nil {
						return fmt.Errorf("setting Storage Accounts for Windows Web App %s: %+v", id, err)
					}
				}
			}

			connectionStrings := expandConnectionStrings(webApp.ConnectionStrings)
			if connectionStrings != nil {
				if _, err := client.UpdateConnectionStrings(ctx, id.ResourceGroup, id.SiteName, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for Windows Web App %s: %+v", id, err)
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
				return fmt.Errorf("reading Windows Web App %s: %+v", id, err)
			}

			if webApp.SiteProperties == nil {
				return fmt.Errorf("reading properties of Windows Web App %s", id)
			}

			// Despite being part of the defined `Get` response model, site_config is always nil so we get it explicitly
			webAppSiteConfig, err := client.GetConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Site Config for Windows Web App %s: %+v", id, err)
			}

			auth, err := client.GetAuthSettings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Windows Web App %s: %+v", id, err)
			}

			backup, err := client.GetBackupConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				if !utils.ResponseWasNotFound(backup.Response) {
					return fmt.Errorf("reading Backup Settings for Windows Web App %s: %+v", id, err)
				}
			}

			logsConfig, err := client.GetDiagnosticLogsConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Diagnostic Logs information for Windows Web App %s: %+v", id, err)
			}

			appSettings, err := client.ListApplicationSettings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading App Settings for Windows Web App %s: %+v", id, err)
			}

			storageAccounts, err := client.ListAzureStorageAccounts(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Storage Account information for Windows Web App %s: %+v", id, err)
			}

			connectionStrings, err := client.ListConnectionStrings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Windows Web App %s: %+v", id, err)
			}

			siteCredentialsFuture, err := client.ListPublishingCredentials(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for Windows Web App %s: %+v", id, err)
			}

			if err := siteCredentialsFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for Site Publishing Credential information for Windows Web App %s: %+v", id, err)
			}
			siteCredentials, err := siteCredentialsFuture.Result(*client)
			if err != nil {
				return fmt.Errorf("reading Site Publishing Credential information for Windows Web App %s: %+v", id, err)
			}

			siteMetadata, err := client.ListMetadata(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Site Metadata for Windows Web App %s: %+v", id, err)
			}

			state := WindowsWebAppModel{
				Name:          id.SiteName,
				ResourceGroup: id.ResourceGroup,
				Location:      location.NormalizeNilable(webApp.Location),
				AppSettings:   flattenAppSettings(appSettings),
				Tags:          tags.ToTypedObject(webApp.Tags),
			}

			webAppProps := webApp.SiteProperties
			if webAppProps.ServerFarmID != nil {
				state.ServicePlanId = *webAppProps.ServerFarmID
			}

			if webAppProps.ClientAffinityEnabled != nil {
				state.ClientAffinityEnabled = *webAppProps.ClientAffinityEnabled
			}

			if webAppProps.ClientCertEnabled != nil {
				state.ClientCertEnabled = *webAppProps.ClientCertEnabled
			}

			if webAppProps.ClientCertMode != "" {
				state.ClientCertMode = string(webAppProps.ClientCertMode)
			}

			if webAppProps.Enabled != nil {
				state.Enabled = *webAppProps.Enabled
			}

			if webAppProps.HTTPSOnly != nil {
				state.HttpsOnly = *webAppProps.HTTPSOnly
			}

			if webAppProps.CustomDomainVerificationID != nil {
				state.CustomDomainVerificationId = *webAppProps.CustomDomainVerificationID
			}

			if webAppProps.DefaultHostName != nil {
				state.DefaultHostname = *webAppProps.DefaultHostName
			}

			if webApp.Kind != nil {
				state.Kind = *webApp.Kind
			}

			if webAppProps.OutboundIPAddresses != nil {
				state.OutboundIPAddresses = *webAppProps.OutboundIPAddresses
				state.OutboundIPAddressList = strings.Split(*webAppProps.OutboundIPAddresses, ",")
			}

			if webAppProps.PossibleOutboundIPAddresses != nil {
				state.PossibleOutboundIPAddresses = *webAppProps.PossibleOutboundIPAddresses
				state.PossibleOutboundIPAddressList = strings.Split(*webAppProps.PossibleOutboundIPAddresses, ",")
			}

			if appAuthSettings := helpers.FlattenAuthSettings(auth); appAuthSettings != nil {
				state.AuthSettings = appAuthSettings
			}

			if appBackupSettings := flattenBackupConfig(backup); appBackupSettings != nil {
				state.Backup = appBackupSettings
			}

			if identity := helpers.FlattenIdentity(webApp.Identity); identity != nil {
				state.Identity = identity
			}

			if logs := flattenLogsConfig(logsConfig); logs != nil {
				state.LogsConfig = logs
			}

			currentStack := ""
			currentStackPtr, ok := siteMetadata.Properties["CURRENT_STACK"]
			if ok {
				currentStack = *currentStackPtr
			}
			if siteConfig := flattenSiteConfigWindows(webAppSiteConfig.SiteConfig, currentStack); siteConfig != nil {
				state.SiteConfig = siteConfig
			}

			if appStorageAccounts := flattenStorageAccounts(storageAccounts); appStorageAccounts != nil {
				state.StorageAccounts = appStorageAccounts
			}

			if appConnectionStrings := flattenConnectionStrings(connectionStrings); appConnectionStrings != nil {
				state.ConnectionStrings = appConnectionStrings
			}

			state.SiteCredentials = helpers.FlattenSiteCredentials(siteCredentials)

			return metadata.Encode(&state)
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

			metadata.Logger.Infof("deleting %s", id)

			deleteMetrics := true // TODO - Look at making this a feature flag?
			deleteEmptyServerFarm := false
			if resp, err := client.Delete(ctx, id.ResourceGroup, id.SiteName, &deleteMetrics, &deleteEmptyServerFarm); err != nil {
				if !utils.ResponseWasNotFound(resp) {
					return fmt.Errorf("deleting Windows Web App %s: %+v", id, err)
				}
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

			id, err := parse.WebAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// TODO - Need locking here when the source control meta resource is added

			var state WindowsWebAppModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			site := web.Site{
				Location: utils.String(state.Location),
				Tags:     tags.FromTypedObject(state.Tags),
				SiteProperties: &web.SiteProperties{
					ServerFarmID:          utils.String(state.ServicePlanId),
					Enabled:               utils.Bool(state.Enabled),
					HTTPSOnly:             utils.Bool(state.HttpsOnly),
					ClientAffinityEnabled: utils.Bool(state.ClientAffinityEnabled),
					ClientCertEnabled:     utils.Bool(state.ClientCertEnabled),
					ClientCertMode:        web.ClientCertMode(state.ClientCertMode),
				},
				Identity: helpers.ExpandIdentity(state.Identity),
			}

			siteConfig, currentStack, err := expandSiteConfigWindows(state.SiteConfig)
			if err != nil {
				return fmt.Errorf("expanding Site Config for Windows Web App %s: %+v", id, err)
			}

			site.SiteConfig = siteConfig
			updateFuture, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, site)
			if err != nil {
				return fmt.Errorf("updating Windows Web App %s: %+v", id, err)
			}
			if err := updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("wating to update %s: %+v", id, err)
			}

			if currentStack != nil && *currentStack != "" {
				siteMetadata := web.StringDictionary{Properties: map[string]*string{}}
				siteMetadata.Properties["CURRENT_STACK"] = currentStack
				if _, err := client.UpdateMetadata(ctx, id.ResourceGroup, id.SiteName, siteMetadata); err != nil {
					return fmt.Errorf("setting Site Metadata for Current Stack on Windows Web App %s: %+v", id, err)
				}
			}

			// (@jackofallops) - App Settings can clobber logs configuration so must be updated before we send any Log updates
			if appSettingsUpdate := expandAppSettings(state.AppSettings); appSettingsUpdate != nil {
				if _, err := client.UpdateApplicationSettings(ctx, id.ResourceGroup, id.SiteName, *appSettingsUpdate); err != nil {
					return fmt.Errorf("updating App Settings for Windows Web App %s: %+v", id, err)
				}
			}

			if connectionStringUpdate := expandConnectionStrings(state.ConnectionStrings); connectionStringUpdate != nil {
				if _, err := client.UpdateConnectionStrings(ctx, id.ResourceGroup, id.SiteName, *connectionStringUpdate); err != nil {
					return fmt.Errorf("updating Connection Strings for Windows Web App %s: %+v", id, err)
				}
			}

			if authUpdate := helpers.ExpandAuthSettings(state.AuthSettings); authUpdate != nil {
				if _, err := client.UpdateAuthSettings(ctx, id.ResourceGroup, id.SiteName, *authUpdate); err != nil {
					return fmt.Errorf("updating Auth Settings for Windows Web App %s: %+v", id, err)
				}
			}

			if backupUpdate := expandBackupConfig(state.Backup); backupUpdate != nil {
				if _, err := client.UpdateBackupConfiguration(ctx, id.ResourceGroup, id.SiteName, *backupUpdate); err != nil {
					return fmt.Errorf("updating Backup Settings for Windows Web App %s: %+v", id, err)
				}
			}

			if logsUpdate := expandLogsConfig(state.LogsConfig); logsUpdate != nil {
				if _, err := client.UpdateDiagnosticLogsConfig(ctx, id.ResourceGroup, id.SiteName, *logsUpdate); err != nil {
					return fmt.Errorf("updating Logs Config for Windows Web App %s: %+v", id, err)
				}
			}

			if storageAccountUpdate := expandStorageConfig(state.StorageAccounts); storageAccountUpdate != nil {
				if _, err := client.UpdateAzureStorageAccounts(ctx, id.ResourceGroup, id.SiteName, *storageAccountUpdate); err != nil {
					return fmt.Errorf("updating Storage Accounts for Windows Web App %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
