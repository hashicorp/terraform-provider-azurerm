package appservice

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WindowsFunctionAppResource struct{}

type WindowsFunctionAppModel struct {
	Name               string `tfschema:"name"`
	ResourceGroup      string `tfschema:"resource_group_name"`
	Location           string `tfschema:"location"`
	ServicePlanId      string `tfschema:"service_plan_id"`
	StorageAccountName string `tfschema:"storage_account_name"`

	StorageAccountKey string `tfschema:"storage_account_access_key"`
	StorageUsesMSI    bool   `tfschema:"storage_uses_managed_identity"` // Storage uses MSI not account key

	AppSettings               map[string]string                      `tfschema:"app_settings"`
	AuthSettings              []helpers.AuthSettings                 `tfschema:"auth_settings"`
	Backup                    []helpers.Backup                       `tfschema:"backup"` // Not supported on Dynamic or Basic plans
	BuiltinLogging            bool                                   `tfschema:"builtin_logging_enabled"`
	ClientCertEnabled         bool                                   `tfschema:"client_cert_enabled"`
	ClientCertMode            string                                 `tfschema:"client_cert_mode"`
	ConnectionStrings         []helpers.ConnectionString             `tfschema:"connection_string"`
	DailyMemoryTimeQuota      int                                    `tfschema:"daily_memory_time_quota"`
	Enabled                   bool                                   `tfschema:"enabled"`
	FunctionExtensionsVersion string                                 `tfschema:"functions_extension_version"`
	ForceDisableContentShare  bool                                   `tfschema:"content_share_force_disabled"`
	HttpsOnly                 bool                                   `tfschema:"https_only"`
	Identity                  []helpers.Identity                     `tfschema:"identity"`
	SiteConfig                []helpers.SiteConfigWindowsFunctionApp `tfschema:"site_config"`
	Tags                      map[string]string                      `tfschema:"tags"`

	// Computed
	CustomDomainVerificationId    string   `tfschema:"custom_domain_verification_id"`
	DefaultHostname               string   `tfschema:"default_hostname"`
	Kind                          string   `tfschema:"kind"`
	OutboundIPAddresses           string   `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string   `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string `tfschema:"possible_outbound_ip_address_list"`

	SiteCredentials []helpers.SiteCredential `tfschema:"site_credential"`
}

var _ sdk.ResourceWithUpdate = WindowsFunctionAppResource{}

var _ sdk.ResourceWithCustomImporter = WindowsFunctionAppResource{}

var _ sdk.ResourceWithCustomizeDiff = WindowsFunctionAppResource{}

func (r WindowsFunctionAppResource) ModelObject() interface{} {
	return &WindowsFunctionAppModel{}
}

func (r WindowsFunctionAppResource) ResourceType() string {
	return "azurerm_windows_function_app"
}

func (r WindowsFunctionAppResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.FunctionAppID
}

func (r WindowsFunctionAppResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
			Description:  "Specifies the name of the Function App.",
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": location.Schema(),

		"service_plan_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ServicePlanID,
			Description:  "The ID of the App Service Plan within which to create this Function App",
		},

		"storage_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: storageValidate.StorageAccountName,
			Description:  "The backend storage account name which will be used by this Function App.",
		},

		"storage_account_access_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true, // TODO - Uncomment this
			ValidateFunc: validation.NoZeroValues,
			ExactlyOneOf: []string{
				"storage_uses_managed_identity",
				"storage_account_access_key",
			},
			Description: "The access key which will be used to access the storage account for the Function App.",
		},

		"storage_uses_managed_identity": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ExactlyOneOf: []string{
				"storage_uses_managed_identity",
				"storage_account_access_key",
			},
			Description: "Should the Function App use its Managed Identity to access storage",
		},

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Description: "A map of key-value pairs for [App Settings](https://docs.microsoft.com/en-us/azure/azure-functions/functions-app-settings) and custom values.",
		},

		"auth_settings": helpers.AuthSettingsSchema(),

		"backup": helpers.BackupSchema(),

		"builtin_logging_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Should built in logging be enabled. Configures `AzureWebJobsDashboard` app setting based on the configured storage setting",
		},

		"client_cert_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Should the function app use Client Certificates",
		},

		"client_cert_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  web.ClientCertModeOptional,
			ValidateFunc: validation.StringInSlice([]string{
				string(web.ClientCertModeOptional),
				string(web.ClientCertModeRequired),
				string(web.ClientCertModeOptionalInteractiveUser),
			}, false),
			Description: "The mode of the Function App's client certificates requirement for incoming requests. Possible values are `Required`, `Optional`, and `OptionalInteractiveUser` ",
		},

		"connection_string": helpers.ConnectionStringSchema(),

		"daily_memory_time_quota": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
			Description:  "The amount of memory in gigabyte-seconds that your application is allowed to consume per day. Setting this value only affects function apps in Consumption Plans.",
		},

		"enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Is the Windows Function App enabled.",
		},

		"content_share_force_disabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Force disable the content share settings.",
		},

		"functions_extension_version": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "~4",
			Description: "The runtime version associated with the Function App.",
		},

		"https_only": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Can the Function App only be accessed via HTTPS?",
		},

		"identity": helpers.IdentitySchema(),

		"site_config": helpers.SiteConfigSchemaWindowsFunctionApp(),

		"tags": tags.Schema(),
	}
}

func (r WindowsFunctionAppResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r WindowsFunctionAppResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var functionApp WindowsFunctionAppModel

			if err := metadata.Decode(&functionApp); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			aseClient := metadata.Client.AppService.AppServiceEnvironmentClient
			servicePlanClient := metadata.Client.AppService.ServicePlanClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewFunctionAppID(subscriptionId, functionApp.ResourceGroup, functionApp.Name)

			servicePlanId, err := parse.ServicePlanID(functionApp.ServicePlanId)
			if err != nil {
				return err
			}

			servicePlan, err := servicePlanClient.Get(ctx, servicePlanId.ResourceGroup, servicePlanId.ServerfarmName)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", servicePlanId, err)
			}

			sendContentSettings := !functionApp.ForceDisableContentShare
			if planSku := servicePlan.Sku; planSku != nil && planSku.Tier != nil {
				switch tier := *planSku.Tier; strings.ToLower(tier) {
				case "dynamic":
				case "elastic":
				case "basic":
					sendContentSettings = false
				case "standard":
					sendContentSettings = false
				case "premiumv2", "premiumv3":
					sendContentSettings = false
				}
			} else {
				return fmt.Errorf("determining plan type for Windows %s: %v", id, err)
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Windows %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			availabilityRequest := web.ResourceNameAvailabilityRequest{
				Name: utils.String(functionApp.Name),
				Type: web.CheckNameResourceTypesMicrosoftWebsites,
			}

			if ase := servicePlan.HostingEnvironmentProfile; ase != nil {
				// Attempt to check the ASE for the appropriate suffix for the name availability request.
				// This varies between internal and external ASE Types, and potentially has other names in other clouds
				// We use the "internal" as the fallback here, if we can read the ASE, we'll get the full one
				nameSuffix := "appserviceenvironment.net"
				if ase.ID != nil {
					aseId, err := parse.AppServiceEnvironmentID(*ase.ID)
					nameSuffix = fmt.Sprintf("%s.%s", aseId.HostingEnvironmentName, nameSuffix)
					if err != nil {
						metadata.Logger.Warnf("could not parse App Service Environment ID determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", functionApp.Name, servicePlanId)
					} else {
						existingASE, err := aseClient.Get(ctx, aseId.ResourceGroup, aseId.HostingEnvironmentName)
						if err != nil {
							metadata.Logger.Warnf("could not read App Service Environment to determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", functionApp.Name, servicePlanId)
						} else if props := existingASE.AppServiceEnvironment; props != nil && props.DNSSuffix != nil && *props.DNSSuffix != "" {
							nameSuffix = *props.DNSSuffix
						}
					}
				}

				availabilityRequest.Name = utils.String(fmt.Sprintf("%s.%s", functionApp.Name, nameSuffix))
				availabilityRequest.IsFqdn = utils.Bool(true)
			}

			checkName, err := client.CheckNameAvailability(ctx, availabilityRequest)
			if err != nil {
				return fmt.Errorf("checking name availability for Windows %s: %+v", id, err)
			}
			if checkName.NameAvailable != nil && !*checkName.NameAvailable {
				return fmt.Errorf("the Site Name %q failed the availability check: %+v", id.SiteName, *checkName.Message)
			}

			storageString := functionApp.StorageAccountName
			if !functionApp.StorageUsesMSI {
				storageString = fmt.Sprintf(helpers.StorageStringFmt, functionApp.StorageAccountName, functionApp.StorageAccountKey, metadata.Client.Account.Environment.StorageEndpointSuffix)
			}
			siteConfig, err := helpers.ExpandSiteConfigWindowsFunctionApp(functionApp.SiteConfig, nil, metadata, functionApp.FunctionExtensionsVersion, storageString, functionApp.StorageUsesMSI)
			if err != nil {
				return fmt.Errorf("expanding site_config for Windows %s: %+v", id, err)
			}

			if functionApp.BuiltinLogging {
				if functionApp.AppSettings == nil {
					functionApp.AppSettings = make(map[string]string)
				}
				if !functionApp.StorageUsesMSI {
					functionApp.AppSettings["AzureWebJobsDashboard"] = storageString
				} else {
					functionApp.AppSettings["AzureWebJobsDashboard__accountName"] = functionApp.StorageAccountName
				}
			}
			if sendContentSettings {
				if functionApp.AppSettings == nil {
					functionApp.AppSettings = make(map[string]string)
				}
				suffix := uuid.New().String()[0:4]
				functionApp.AppSettings["WEBSITE_CONTENTSHARE"] = fmt.Sprintf("%s-%s", strings.ToLower(functionApp.Name), suffix)
				functionApp.AppSettings["WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"] = storageString
			}

			siteConfig.WindowsFxVersion = helpers.EncodeFunctionAppWindowsFxVersion(functionApp.SiteConfig[0].ApplicationStack)
			siteConfig.AppSettings = helpers.MergeUserAppSettings(siteConfig.AppSettings, functionApp.AppSettings)

			siteEnvelope := web.Site{
				Location: utils.String(functionApp.Location),
				Tags:     tags.FromTypedObject(functionApp.Tags),
				Kind:     utils.String("functionapp"),
				Identity: helpers.ExpandIdentity(functionApp.Identity),
				SiteProperties: &web.SiteProperties{
					ServerFarmID:         utils.String(functionApp.ServicePlanId),
					Enabled:              utils.Bool(functionApp.Enabled),
					HTTPSOnly:            utils.Bool(functionApp.HttpsOnly),
					SiteConfig:           siteConfig,
					ClientCertEnabled:    utils.Bool(functionApp.ClientCertEnabled),
					ClientCertMode:       web.ClientCertMode(functionApp.ClientCertMode),
					DailyMemoryTimeQuota: utils.Int32(int32(functionApp.DailyMemoryTimeQuota)),
				},
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, siteEnvelope)
			if err != nil {
				return fmt.Errorf("creating Windows %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of Windows %s: %+v", id, err)
			}

			updateFuture, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, siteEnvelope)
			if err != nil {
				return fmt.Errorf("updating properties of Windows %s: %+v", id, err)
			}
			if err := updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of Windows %s: %+v", id, err)
			}

			backupConfig := helpers.ExpandBackupConfig(functionApp.Backup)
			if backupConfig.BackupRequestProperties != nil {
				if _, err := client.UpdateBackupConfiguration(ctx, id.ResourceGroup, id.SiteName, *backupConfig); err != nil {
					return fmt.Errorf("adding Backup Settings for Windows %s: %+v", id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(functionApp.AuthSettings)
			if auth.SiteAuthSettingsProperties != nil {
				if _, err := client.UpdateAuthSettings(ctx, id.ResourceGroup, id.SiteName, *auth); err != nil {
					return fmt.Errorf("setting Authorisation Settings for Windows %s: %+v", id, err)
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(functionApp.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStrings(ctx, id.ResourceGroup, id.SiteName, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for Windows %s: %+v", id, err)
				}
			}

			if _, ok := metadata.ResourceData.GetOk("site_config.0.app_service_logs"); ok {
				appServiceLogs := helpers.ExpandFunctionAppAppServiceLogs(functionApp.SiteConfig[0].AppServiceLogs)
				if _, err := client.UpdateDiagnosticLogsConfig(ctx, id.ResourceGroup, id.SiteName, appServiceLogs); err != nil {
					return fmt.Errorf("updating App Service Log Settings for %s: %+v", id, err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WindowsFunctionAppResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 25 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := parse.FunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			functionApp, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				if utils.ResponseWasNotFound(functionApp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Windows %s: %+v", id, err)
			}

			if functionApp.SiteProperties == nil {
				return fmt.Errorf("reading properties of Windows %s", id)
			}
			props := *functionApp.SiteProperties

			appSettingsResp, err := client.ListApplicationSettings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading App Settings for Windows %s: %+v", id, err)
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

			logs, err := client.GetDiagnosticLogsConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading logs configuration for Windows %s: %+v", id, err)
			}

			state := WindowsFunctionAppModel{
				Name:                 id.SiteName,
				ResourceGroup:        id.ResourceGroup,
				ServicePlanId:        utils.NormalizeNilableString(props.ServerFarmID),
				Location:             location.NormalizeNilable(functionApp.Location),
				Enabled:              utils.NormaliseNilableBool(functionApp.Enabled),
				ClientCertMode:       string(functionApp.ClientCertMode),
				DailyMemoryTimeQuota: int(utils.NormaliseNilableInt32(props.DailyMemoryTimeQuota)),
				Tags:                 tags.ToTypedObject(functionApp.Tags),
				Kind:                 utils.NormalizeNilableString(functionApp.Kind),
			}

			if identity := helpers.FlattenIdentity(functionApp.Identity); identity != nil {
				state.Identity = identity
			}

			configResp, err := client.GetConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("making Read request on AzureRM Function App Configuration %q: %+v", id.SiteName, err)
			}

			siteConfig, err := helpers.FlattenSiteConfigWindowsFunctionApp(configResp.SiteConfig)
			if err != nil {
				return fmt.Errorf("reading Site Config for Windows %s: %+v", id, err)
			}
			state.SiteConfig = []helpers.SiteConfigWindowsFunctionApp{*siteConfig}

			state.unpackWindowsFunctionAppSettings(appSettingsResp)

			state.ConnectionStrings = helpers.FlattenConnectionStrings(connectionStrings)

			state.SiteCredentials = helpers.FlattenSiteCredentials(siteCredentials)

			state.AuthSettings = helpers.FlattenAuthSettings(auth)

			state.Backup = helpers.FlattenBackupConfig(backup)

			state.SiteConfig[0].AppServiceLogs = helpers.FlattenFunctionAppAppServiceLogs(logs)

			state.HttpsOnly = utils.NormaliseNilableBool(functionApp.HTTPSOnly)
			state.ClientCertEnabled = utils.NormaliseNilableBool(functionApp.ClientCertEnabled)

			return metadata.Encode(&state)
		},
	}
}

func (r WindowsFunctionAppResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := parse.FunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting Windows %s", *id)

			deleteMetrics := true
			deleteEmptyServerFarm := false
			if _, err := client.Delete(ctx, id.ResourceGroup, id.SiteName, &deleteMetrics, &deleteEmptyServerFarm); err != nil {
				return fmt.Errorf("deleting Windows %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r WindowsFunctionAppResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := parse.FunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state WindowsFunctionAppModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Windows %s: %v", id, err)
			}

			// Some service plan updates are allowed - see customiseDiff for exceptions
			if metadata.ResourceData.HasChange("service_plan_id") {
				existing.SiteProperties.ServerFarmID = utils.String(state.ServicePlanId)
			}

			if metadata.ResourceData.HasChange("enabled") {
				existing.SiteProperties.Enabled = utils.Bool(state.Enabled)
			}

			if metadata.ResourceData.HasChange("https_only") {
				existing.SiteProperties.HTTPSOnly = utils.Bool(state.HttpsOnly)
			}

			if metadata.ResourceData.HasChange("client_cert_enabled") {
				existing.SiteProperties.ClientCertEnabled = utils.Bool(state.ClientCertEnabled)
			}

			if metadata.ResourceData.HasChange("client_cert_mode") {
				existing.SiteProperties.ClientCertMode = web.ClientCertMode(state.ClientCertMode)
			}

			if metadata.ResourceData.HasChange("identity") {
				existing.Identity = helpers.ExpandIdentity(state.Identity)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.FromTypedObject(state.Tags)
			}

			storageString := ""
			if !state.StorageUsesMSI {
				storageString = fmt.Sprintf(helpers.StorageStringFmt, state.StorageAccountName, state.StorageAccountKey, metadata.Client.Account.Environment.StorageEndpointSuffix)
			}

			// Note: We process this regardless to give us a "clean" view of service-side app_settings, so we can reconcile the user-defined entries later
			siteConfig, err := helpers.ExpandSiteConfigWindowsFunctionApp(state.SiteConfig, existing.SiteConfig, metadata, state.FunctionExtensionsVersion, storageString, state.StorageUsesMSI)
			if state.BuiltinLogging {
				if state.AppSettings == nil && !state.StorageUsesMSI {
					state.AppSettings = make(map[string]string)
				}
				state.AppSettings["AzureWebJobsDashboard"] = storageString
			}

			if metadata.ResourceData.HasChange("site_config") {
				if err != nil {
					return fmt.Errorf("expanding Site Config for Windows %s: %+v", id, err)
				}
				existing.SiteConfig = siteConfig
			}

			if metadata.ResourceData.HasChange("site_config.0.application_stack") {
				existing.SiteConfig.WindowsFxVersion = helpers.EncodeFunctionAppWindowsFxVersion(state.SiteConfig[0].ApplicationStack)
			}

			existing.SiteConfig.AppSettings = helpers.MergeUserAppSettings(siteConfig.AppSettings, state.AppSettings)

			updateFuture, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, existing)
			if err != nil {
				return fmt.Errorf("updating Windows %s: %+v", id, err)
			}
			if err := updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting to update %s: %+v", id, err)
			}

			if _, err := client.UpdateConfiguration(ctx, id.ResourceGroup, id.SiteName, web.SiteConfigResource{SiteConfig: siteConfig}); err != nil {
				return fmt.Errorf("updating Site Config for Windows %s: %+v", id, err)
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

			if metadata.ResourceData.HasChange("auth_settings") {
				authUpdate := helpers.ExpandAuthSettings(state.AuthSettings)
				if _, err := client.UpdateAuthSettings(ctx, id.ResourceGroup, id.SiteName, *authUpdate); err != nil {
					return fmt.Errorf("updating Auth Settings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("backup") {
				backupUpdate := helpers.ExpandBackupConfig(state.Backup)
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

			if metadata.ResourceData.HasChange("site_config.0.app_service_logs") {
				appServiceLogs := helpers.ExpandFunctionAppAppServiceLogs(state.SiteConfig[0].AppServiceLogs)
				if _, err := client.UpdateDiagnosticLogsConfig(ctx, id.ResourceGroup, id.SiteName, appServiceLogs); err != nil {
					return fmt.Errorf("updating App Service Log Settings for %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r WindowsFunctionAppResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.AppService.WebAppsClient
		servicePlanClient := metadata.Client.AppService.ServicePlanClient

		id, err := parse.FunctionAppID(metadata.ResourceData.Id())
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
		if strings.Contains(strings.ToLower(*sp.Kind), "linux") || !(strings.Contains(strings.ToLower(*sp.Kind), "elastic") || strings.Contains(strings.ToLower(*sp.Kind), "functionapp") || strings.Contains(strings.ToLower(*sp.Kind), "app")) {
			return fmt.Errorf("specified Service Plan is not a Windows Functionapp plan")
		}

		return nil
	}
}

func (r WindowsFunctionAppResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.ServicePlanClient
			rd := metadata.ResourceDiff

			if rd.HasChange("service_plan_id") {
				currentPlanIdRaw, newPlanIdRaw := rd.GetChange("service_plan_id")
				if newPlanIdRaw.(string) == "" {
					// Plans creating a new service_plan inline will be empty as `Computed` known after apply
					return nil
				}
				newPlanId, err := parse.ServicePlanID(newPlanIdRaw.(string))
				if err != nil {
					return fmt.Errorf("reading new plan id %+v", err)
				}

				var currentTierIsDynamic, newTierIsDynamic, newTierIsBasic bool

				newPlan, err := client.Get(ctx, newPlanId.ResourceGroup, newPlanId.ServerfarmName)
				if err != nil {
					return fmt.Errorf("could not read new Service Plan to check tier %s: %+v", newPlanId, err)
				}
				if planSku := newPlan.Sku; planSku != nil {
					if tier := planSku.Tier; tier != nil {
						newTierIsDynamic = strings.EqualFold(*tier, "dynamic")
						newTierIsBasic = strings.EqualFold(*tier, "basic")
					}
				}

				// Service Plans can only be updated in place when both New and Existing are not Dynamic
				if currentPlanIdRaw.(string) != "" {
					currentPlanId, err := parse.ServicePlanID(currentPlanIdRaw.(string))
					if err != nil {
						return fmt.Errorf("reading existing plan id %+v", err)
					}

					currentPlan, err := client.Get(ctx, currentPlanId.ResourceGroup, currentPlanId.ServerfarmName)
					if err != nil {
						return fmt.Errorf("could not read current Service Plan to check tier %s: %+v", currentPlanId, err)
					}

					if planSku := currentPlan.Sku; planSku != nil {
						if tier := planSku.Tier; tier != nil {
							currentTierIsDynamic = strings.EqualFold(*tier, "dynamic")
						}
					}

					if currentTierIsDynamic || newTierIsDynamic {
						if err := rd.ForceNew("service_plan_id"); err != nil {
							return err
						}
					}
				}
				if _, ok := rd.GetOk("backup"); ok && newTierIsDynamic {
					return fmt.Errorf("cannot specify backup configuration for Dynamic tier Service Plans, Standard or higher is required")
				}
				if _, ok := rd.GetOk("backup"); ok && newTierIsBasic {
					return fmt.Errorf("cannot specify backup configuration for Basic tier Service Plans, Standard or higher is required")
				}
			}
			return nil
		},
	}
}

func (m *WindowsFunctionAppModel) unpackWindowsFunctionAppSettings(input web.StringDictionary) {
	if input.Properties == nil {
		return
	}

	appSettings := make(map[string]string)
	var dockerSettings helpers.ApplicationStackDocker
	m.BuiltinLogging = false

	for k, v := range input.Properties {
		switch k {
		case "FUNCTIONS_EXTENSION_VERSION":
			m.FunctionExtensionsVersion = utils.NormalizeNilableString(v)

		case "WEBSITE_NODE_DEFAULT_VERSION": // Note - This is only set if it's not the default of 12, but we collect it from WindowsFxVersion so can discard it here
		case "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING":
		case "WEBSITE_CONTENTSHARE":
		case "WEBSITE_HTTPLOGGING_RETENTION_DAYS":
		case "FUNCTIONS_WORKER_RUNTIME":
			if m.SiteConfig[0].ApplicationStack != nil {
				m.SiteConfig[0].ApplicationStack[0].CustomHandler = strings.EqualFold(*v, "custom")
			}

		case "DOCKER_REGISTRY_SERVER_URL":
			dockerSettings.RegistryURL = utils.NormalizeNilableString(v)

		case "DOCKER_REGISTRY_SERVER_USERNAME":
			dockerSettings.RegistryUsername = utils.NormalizeNilableString(v)

		case "DOCKER_REGISTRY_SERVER_PASSWORD":
			dockerSettings.RegistryPassword = utils.NormalizeNilableString(v)

		case "APPINSIGHTS_INSTRUMENTATIONKEY":
			m.SiteConfig[0].AppInsightsInstrumentationKey = utils.NormalizeNilableString(v)

		case "APPLICATIONINSIGHTS_CONNECTION_STRING":
			m.SiteConfig[0].AppInsightsConnectionString = utils.NormalizeNilableString(v)

		case "AzureWebJobsStorage":
			m.StorageAccountName, m.StorageAccountKey = helpers.ParseWebJobsStorageString(v)

		case "AzureWebJobsDashboard":
			m.BuiltinLogging = true

		case "WEBSITE_HEALTHCHECK_MAXPINGFAILURES":
			i, _ := strconv.Atoi(utils.NormalizeNilableString(v))
			m.SiteConfig[0].HealthCheckEvictionTime = utils.NormaliseNilableInt(&i)

		default:
			appSettings[k] = utils.NormalizeNilableString(v)
		}
	}

	m.AppSettings = appSettings
}
