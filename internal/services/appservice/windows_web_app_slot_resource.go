package appservice

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WindowsWebAppSlotResource struct{}

type WindowsWebAppSlotModel struct {
	Name                          string                                `tfschema:"name"`
	AppServiceId                  string                                `tfschema:"app_service_id"`
	AppSettings                   map[string]string                     `tfschema:"app_settings"`
	AuthSettings                  []helpers.AuthSettings                `tfschema:"auth_settings"`
	Backup                        []helpers.Backup                      `tfschema:"backup"`
	ClientAffinityEnabled         bool                                  `tfschema:"client_affinity_enabled"`
	ClientCertEnabled             bool                                  `tfschema:"client_certificate_enabled"`
	ClientCertMode                string                                `tfschema:"client_certificate_mode"`
	ClientCertExclusionPaths      string                                `tfschema:"client_certificate_exclusion_paths"`
	Enabled                       bool                                  `tfschema:"enabled"`
	HttpsOnly                     bool                                  `tfschema:"https_only"`
	KeyVaultReferenceIdentityID   string                                `tfschema:"key_vault_reference_identity_id"`
	LogsConfig                    []helpers.LogsConfig                  `tfschema:"logs"`
	SiteConfig                    []helpers.SiteConfigWindowsWebAppSlot `tfschema:"site_config"`
	StorageAccounts               []helpers.StorageAccount              `tfschema:"storage_account"`
	ConnectionStrings             []helpers.ConnectionString            `tfschema:"connection_string"`
	CustomDomainVerificationId    string                                `tfschema:"custom_domain_verification_id"`
	DefaultHostname               string                                `tfschema:"default_hostname"`
	Kind                          string                                `tfschema:"kind"`
	OutboundIPAddresses           string                                `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string                              `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string                                `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string                              `tfschema:"possible_outbound_ip_address_list"`
	SiteCredentials               []helpers.SiteCredential              `tfschema:"site_credential"`
	ZipDeployFile                 string                                `tfschema:"zip_deploy_file"`
	Tags                          map[string]string                     `tfschema:"tags"`
	VirtualNetworkSubnetID        string                                `tfschema:"virtual_network_subnet_id"`
}

var _ sdk.ResourceWithUpdate = WindowsWebAppSlotResource{}

func (r WindowsWebAppSlotResource) ModelObject() interface{} {
	return &WindowsWebAppSlotModel{}
}

func (r WindowsWebAppSlotResource) ResourceType() string {
	return "azurerm_windows_web_app_slot"
}

func (r WindowsWebAppSlotResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.WebAppSlotID
}

func (r WindowsWebAppSlotResource) Arguments() map[string]*pluginsdk.Schema {
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
			ForceNew:     true,
			ValidateFunc: validate.WebAppID,
		},

		// Optional

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"auth_settings": helpers.AuthSettingsSchema(),

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

		"site_config": helpers.SiteConfigSchemaWindowsWebAppSlot(),

		"storage_account": helpers.StorageAccountSchemaWindows(),

		"zip_deploy_file": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The local path and filename of the Zip packaged application to deploy to this Windows Web App. **Note:** Using this value requires `WEBSITE_RUN_FROM_PACKAGE=1` on the App in `app_settings`.",
		},

		"tags": tags.Schema(),

		"virtual_network_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: networkValidate.SubnetID,
		},
	}
}

func (r WindowsWebAppSlotResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r WindowsWebAppSlotResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var webAppSlot WindowsWebAppSlotModel
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
				return fmt.Errorf("reading parent Windows Web App for %s: %+v", id, err)
			}
			if webApp.Location == nil {
				return fmt.Errorf("could not determine location for %s: %+v", id, err)
			}
			siteProps := webApp.SiteProperties
			if siteProps == nil || siteProps.ServerFarmID == nil {
				return fmt.Errorf("could not determine Service Plan ID for %s: %+v", id, err)
			}

			existing, err := client.GetSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Windows %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			siteConfig, currentStack, err := helpers.ExpandSiteConfigWindowsWebAppSlot(webAppSlot.SiteConfig, nil, metadata)
			if err != nil {
				return err
			}

			if *currentStack == helpers.CurrentStackNode {
				if webAppSlot.AppSettings == nil {
					webAppSlot.AppSettings = make(map[string]string, 0)
				}
				webAppSlot.AppSettings["WEBSITE_NODE_DEFAULT_VERSION"] = webAppSlot.SiteConfig[0].ApplicationStack[0].NodeVersion
			}
			siteConfig.AppSettings = helpers.ExpandAppSettingsForCreate(webAppSlot.AppSettings)

			expandedIdentity, err := expandIdentity(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			siteEnvelope := web.Site{
				Location: webApp.Location,
				Tags:     tags.FromTypedObject(webAppSlot.Tags),
				Identity: expandedIdentity,
				SiteProperties: &web.SiteProperties{
					ServerFarmID:             siteProps.ServerFarmID,
					Enabled:                  utils.Bool(webAppSlot.Enabled),
					HTTPSOnly:                utils.Bool(webAppSlot.HttpsOnly),
					SiteConfig:               siteConfig,
					ClientAffinityEnabled:    utils.Bool(webAppSlot.ClientAffinityEnabled),
					ClientCertEnabled:        utils.Bool(webAppSlot.ClientCertEnabled),
					ClientCertMode:           web.ClientCertMode(webAppSlot.ClientCertMode),
					ClientCertExclusionPaths: utils.String(webAppSlot.ClientCertExclusionPaths),
				},
			}

			if webAppSlot.KeyVaultReferenceIdentityID != "" {
				siteEnvelope.SiteProperties.KeyVaultReferenceIdentity = utils.String(webAppSlot.KeyVaultReferenceIdentityID)
			}

			if webAppSlot.VirtualNetworkSubnetID != "" {
				siteEnvelope.SiteProperties.VirtualNetworkSubnetID = utils.String(webAppSlot.VirtualNetworkSubnetID)
			}

			future, err := client.CreateOrUpdateSlot(ctx, id.ResourceGroup, id.SiteName, siteEnvelope, id.SlotName)
			if err != nil {
				return fmt.Errorf("creating Windows %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of Windows %s: %+v", id, err)
			}

			metadata.SetID(id)

			if currentStack != nil && *currentStack != "" {
				siteMetadata := web.StringDictionary{Properties: map[string]*string{}}
				siteMetadata.Properties["CURRENT_STACK"] = currentStack
				if _, err := client.UpdateMetadataSlot(ctx, id.ResourceGroup, id.SiteName, siteMetadata, id.SlotName); err != nil {
					return fmt.Errorf("setting Site Metadata for Current Stack on Windows %s: %+v", id, err)
				}
			}

			appSettings := helpers.ExpandAppSettingsForUpdate(webAppSlot.AppSettings)
			if metadata.ResourceData.HasChange("site_config.0.health_check_eviction_time_in_min") {
				appSettings.Properties["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"] = utils.String(strconv.Itoa(webAppSlot.SiteConfig[0].HealthCheckEvictionTime))
			}
			if len(appSettings.Properties) > 0 {
				if _, err := client.UpdateApplicationSettingsSlot(ctx, id.ResourceGroup, id.SiteName, *appSettings, id.SlotName); err != nil {
					return fmt.Errorf("setting App Settings for Windows %s: %+v", id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(webAppSlot.AuthSettings)
			if auth.SiteAuthSettingsProperties != nil {
				if _, err := client.UpdateAuthSettingsSlot(ctx, id.ResourceGroup, id.SiteName, *auth, id.SlotName); err != nil {
					return fmt.Errorf("setting Authorisation Settings for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("logs") {
				logsConfig := helpers.ExpandLogsConfig(webAppSlot.LogsConfig)
				if logsConfig.SiteLogsConfigProperties != nil {
					if _, err := client.UpdateDiagnosticLogsConfigSlot(ctx, id.ResourceGroup, id.SiteName, *logsConfig, id.SlotName); err != nil {
						return fmt.Errorf("setting Diagnostic Logs Configuration for Windows %s: %+v", id, err)
					}
				}
			}

			backupConfig := helpers.ExpandBackupConfig(webAppSlot.Backup)
			if backupConfig.BackupRequestProperties != nil {
				if _, err := client.UpdateBackupConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, *backupConfig, id.SlotName); err != nil {
					return fmt.Errorf("adding Backup Settings for Windows %s: %+v", id, err)
				}
			}

			storageConfig := helpers.ExpandStorageConfig(webAppSlot.StorageAccounts)
			if storageConfig.Properties != nil {
				if _, err := client.UpdateAzureStorageAccountsSlot(ctx, id.ResourceGroup, id.SiteName, *storageConfig, id.SlotName); err != nil {
					if err != nil {
						return fmt.Errorf("setting Storage Accounts for Windows %s: %+v", id, err)
					}
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(webAppSlot.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStringsSlot(ctx, id.ResourceGroup, id.SiteName, *connectionStrings, id.SlotName); err != nil {
					return fmt.Errorf("setting Connection Strings for Windows %s: %+v", id, err)
				}
			}

			if webAppSlot.ZipDeployFile != "" {
				if err = helpers.GetCredentialsAndPublishSlot(ctx, client, id.ResourceGroup, id.SiteName, webAppSlot.ZipDeployFile, id.SlotName); err != nil {
					return err
				}
			}

			return nil
		},

		Timeout: 30 * time.Minute,
	}
}

func (r WindowsWebAppSlotResource) Read() sdk.ResourceFunc {
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
				return fmt.Errorf("reading Windows %s: %+v", id, err)
			}

			props := webAppSlot.SiteProperties
			if props == nil {
				return fmt.Errorf("reading properties of Windows %s", id)
			}

			// Despite being part of the defined `Get` response model, site_config is always nil so we get it explicitly
			webAppSiteConfig, err := client.GetConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading Site Config for Windows %s: %+v", id, err)
			}

			auth, err := client.GetAuthSettingsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Windows %s: %+v", id, err)
			}

			backup, err := client.GetBackupConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				if !utils.ResponseWasNotFound(backup.Response) {
					return fmt.Errorf("reading Backup Settings for Windows %s: %+v", id, err)
				}
			}

			logsConfig, err := client.GetDiagnosticLogsConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading Diagnostic Logs information for Windows %s: %+v", id, err)
			}

			appSettings, err := client.ListApplicationSettingsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading App Settings for Windows %s: %+v", id, err)
			}

			storageAccounts, err := client.ListAzureStorageAccountsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading Storage Account information for Windows %s: %+v", id, err)
			}

			connectionStrings, err := client.ListConnectionStringsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Windows %s: %+v", id, err)
			}

			siteCredentialsFuture, err := client.ListPublishingCredentialsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
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

			siteMetadata, err := client.ListMetadataSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading Site Metadata for Windows %s: %+v", id, err)
			}

			state := WindowsWebAppSlotModel{
				Name:                        id.SlotName,
				AppServiceId:                parse.NewWebAppID(id.SubscriptionId, id.ResourceGroup, id.SiteName).ID(),
				AuthSettings:                helpers.FlattenAuthSettings(auth),
				Backup:                      helpers.FlattenBackupConfig(backup),
				ClientAffinityEnabled:       utils.NormaliseNilableBool(props.ClientAffinityEnabled),
				ClientCertEnabled:           utils.NormaliseNilableBool(props.ClientCertEnabled),
				ClientCertMode:              string(props.ClientCertMode),
				ClientCertExclusionPaths:    utils.NormalizeNilableString(props.ClientCertExclusionPaths),
				ConnectionStrings:           helpers.FlattenConnectionStrings(connectionStrings),
				CustomDomainVerificationId:  utils.NormalizeNilableString(props.CustomDomainVerificationID),
				DefaultHostname:             utils.NormalizeNilableString(props.DefaultHostName),
				Enabled:                     utils.NormaliseNilableBool(props.Enabled),
				HttpsOnly:                   utils.NormaliseNilableBool(props.HTTPSOnly),
				KeyVaultReferenceIdentityID: utils.NormalizeNilableString(props.KeyVaultReferenceIdentity),
				Kind:                        utils.NormalizeNilableString(webAppSlot.Kind),
				LogsConfig:                  helpers.FlattenLogsConfig(logsConfig),
				SiteCredentials:             helpers.FlattenSiteCredentials(siteCredentials),
				StorageAccounts:             helpers.FlattenStorageAccounts(storageAccounts),
				Tags:                        tags.ToTypedObject(webAppSlot.Tags),
			}

			if subnetId := utils.NormalizeNilableString(props.VirtualNetworkSubnetID); subnetId != "" {
				state.VirtualNetworkSubnetID = subnetId
			}

			var healthCheckCount *int
			state.AppSettings, healthCheckCount = helpers.FlattenAppSettings(appSettings)

			if v := props.OutboundIPAddresses; v != nil {
				state.OutboundIPAddresses = *v
				state.OutboundIPAddressList = strings.Split(*v, ",")
			}

			if v := props.PossibleOutboundIPAddresses; v != nil {
				state.PossibleOutboundIPAddresses = *v
				state.PossibleOutboundIPAddressList = strings.Split(*v, ",")
			}

			currentStack := ""
			currentStackPtr, ok := siteMetadata.Properties["CURRENT_STACK"]
			if ok {
				currentStack = *currentStackPtr
			}

			state.SiteConfig = helpers.FlattenSiteConfigWindowsAppSlot(webAppSiteConfig.SiteConfig, currentStack, healthCheckCount)

			if nodeVer, ok := state.AppSettings["WEBSITE_NODE_DEFAULT_VERSION"]; ok {
				if nodeVer != "6.9.1" { // Slots appear to have an invalid value for this by default?
					if state.SiteConfig[0].ApplicationStack == nil {
						state.SiteConfig[0].ApplicationStack = make([]helpers.ApplicationStackWindows, 0)
					}
					state.SiteConfig[0].ApplicationStack[0].NodeVersion = nodeVer
				}
				delete(state.AppSettings, "WEBSITE_NODE_DEFAULT_VERSION") // TODO - Allow this to be set directly in app_settings?
			}

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

func (r WindowsWebAppSlotResource) Delete() sdk.ResourceFunc {
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
				return fmt.Errorf("deleting Windows %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r WindowsWebAppSlotResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := parse.WebAppSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// TODO - Need locking here for the source control meta resource?

			var state WindowsWebAppSlotModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.GetSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading Windows %s: %v", id, err)
			}

			if metadata.ResourceData.HasChange("enabled") {
				existing.SiteProperties.Enabled = utils.Bool(state.Enabled)
			}
			if metadata.ResourceData.HasChange("https_only") {
				existing.SiteProperties.HTTPSOnly = utils.Bool(state.HttpsOnly)
			}
			if metadata.ResourceData.HasChange("client_affinity_enabled") {
				existing.SiteProperties.ClientAffinityEnabled = utils.Bool(state.ClientAffinityEnabled)
			}
			if metadata.ResourceData.HasChange("client_certificate_enabled") {
				existing.SiteProperties.ClientCertEnabled = utils.Bool(state.ClientCertEnabled)
			}
			if metadata.ResourceData.HasChange("client_certificate_mode") {
				existing.SiteProperties.ClientCertMode = web.ClientCertMode(state.ClientCertMode)
			}
			if metadata.ResourceData.HasChange("client_certificate_exclusion_paths") {
				existing.SiteProperties.ClientCertExclusionPaths = utils.String(state.ClientCertExclusionPaths)
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := expandIdentity(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				existing.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("key_vault_reference_identity_id") {
				existing.KeyVaultReferenceIdentity = utils.String(state.KeyVaultReferenceIdentityID)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.FromTypedObject(state.Tags)
			}

			currentStack := ""
			stateConfig := state.SiteConfig[0]
			if len(stateConfig.ApplicationStack) == 1 {
				currentStack = stateConfig.ApplicationStack[0].CurrentStack
			}

			if metadata.ResourceData.HasChange("site_config") {
				siteConfig, stack, err := helpers.ExpandSiteConfigWindowsWebAppSlot(state.SiteConfig, existing.SiteConfig, metadata)
				if err != nil {
					return fmt.Errorf("expanding Site Config for Windows %s: %+v", id, err)
				}
				currentStack = *stack
				existing.SiteConfig = siteConfig
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
					existing.SiteProperties.VirtualNetworkSubnetID = utils.String(subnetId)
				}
			}

			updateFuture, err := client.CreateOrUpdateSlot(ctx, id.ResourceGroup, id.SiteName, existing, id.SlotName)
			if err != nil {
				return fmt.Errorf("updating Windows %s: %+v", id, err)
			}
			if err := updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("wating to update %s: %+v", id, err)
			}

			siteMetadata := web.StringDictionary{Properties: map[string]*string{}}
			siteMetadata.Properties["CURRENT_STACK"] = utils.String(currentStack)
			if _, err := client.UpdateMetadataSlot(ctx, id.ResourceGroup, id.SiteName, siteMetadata, id.SlotName); err != nil {
				return fmt.Errorf("setting Site Metadata for Current Stack on Windows %s: %+v", id, err)
			}

			// (@jackofallops) - App Settings can clobber logs configuration so must be updated before we send any Log updates
			if metadata.ResourceData.HasChange("app_settings") || metadata.ResourceData.HasChange("site_config.0.health_check_eviction_time_in_min") || metadata.ResourceData.HasChange("site_config.0.application_stack.0.node_version") {
				appSettingsUpdate := helpers.ExpandAppSettingsForUpdate(state.AppSettings)
				appSettingsUpdate.Properties["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"] = utils.String(strconv.Itoa(state.SiteConfig[0].HealthCheckEvictionTime))
				appSettingsUpdate.Properties["WEBSITE_NODE_DEFAULT_VERSION"] = pointer.To(state.SiteConfig[0].ApplicationStack[0].NodeVersion)
				if _, err := client.UpdateApplicationSettingsSlot(ctx, id.ResourceGroup, id.SiteName, *appSettingsUpdate, id.SlotName); err != nil {
					return fmt.Errorf("updating App Settings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("connection_string") {
				connectionStringUpdate := helpers.ExpandConnectionStrings(state.ConnectionStrings)
				if connectionStringUpdate.Properties == nil {
					connectionStringUpdate.Properties = map[string]*web.ConnStringValueTypePair{}
				}
				if _, err := client.UpdateConnectionStringsSlot(ctx, id.ResourceGroup, id.SiteName, *connectionStringUpdate, id.SlotName); err != nil {
					return fmt.Errorf("updating Connection Strings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings") {
				authUpdate := helpers.ExpandAuthSettings(state.AuthSettings)
				if _, err := client.UpdateAuthSettingsSlot(ctx, id.ResourceGroup, id.SiteName, *authUpdate, id.SlotName); err != nil {
					return fmt.Errorf("updating Auth Settings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("backup") {
				backupUpdate := helpers.ExpandBackupConfig(state.Backup)
				if backupUpdate.BackupRequestProperties == nil {
					if _, err := client.DeleteBackupConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName); err != nil {
						return fmt.Errorf("removing Backup Settings for Windows %s: %+v", id, err)
					}
				} else {
					if _, err := client.UpdateBackupConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, *backupUpdate, id.SlotName); err != nil {
						return fmt.Errorf("updating Backup Settings for Windows %s: %+v", id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("logs") {
				logsUpdate := helpers.ExpandLogsConfig(state.LogsConfig)
				if logsUpdate.SiteLogsConfigProperties == nil {
					logsUpdate = helpers.DisabledLogsConfig() // The API is update only, so we need to send an update with everything switched of when a user removes the "logs" block
				}
				if _, err := client.UpdateDiagnosticLogsConfigSlot(ctx, id.ResourceGroup, id.SiteName, *logsUpdate, id.SlotName); err != nil {
					return fmt.Errorf("updating Logs Config for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("storage_account") {
				storageAccountUpdate := helpers.ExpandStorageConfig(state.StorageAccounts)
				if _, err := client.UpdateAzureStorageAccountsSlot(ctx, id.ResourceGroup, id.SiteName, *storageAccountUpdate, id.SlotName); err != nil {
					return fmt.Errorf("updating Storage Accounts for Windows %s: %+v", id, err)
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
