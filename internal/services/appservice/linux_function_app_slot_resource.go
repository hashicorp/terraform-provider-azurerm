package appservice

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
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

type LinuxFunctionAppSlotResource struct{}

type LinuxFunctionAppSlotModel struct {
	Name                          string                                   `tfschema:"name"`
	ResourceGroup                 string                                   `tfschema:"resource_group_name"`
	FunctionAppName               string                                   `tfschema:"function_app_name"`
	Location                      string                                   `tfschema:"location"`
	ServicePlanId                 string                                   `tfschema:"service_plan_id"`
	StorageAccountName            string                                   `tfschema:"storage_account_name"`
	StorageAccountKey             string                                   `tfschema:"storage_account_access_key"`
	StorageUsesMSI                bool                                     `tfschema:"storage_uses_managed_identity"` // Storage uses MSI not account key
	AppSettings                   map[string]string                        `tfschema:"app_settings"`
	AuthSettings                  []helpers.AuthSettings                   `tfschema:"auth_settings"`
	Backup                        []helpers.Backup                         `tfschema:"backup"` // Not supported on Dynamic or Basic plans
	BuiltinLogging                bool                                     `tfschema:"builtin_logging_enabled"`
	ClientCertEnabled             bool                                     `tfschema:"client_certificate_enabled"`
	ClientCertMode                string                                   `tfschema:"client_certificate_mode"`
	ConnectionStrings             []helpers.ConnectionString               `tfschema:"connection_string"`
	DailyMemoryTimeQuota          int                                      `tfschema:"daily_memory_time_quota"` // TODO - Value ignored in for linux apps, even in Consumption plans?
	Enabled                       bool                                     `tfschema:"enabled"`
	FunctionExtensionsVersion     string                                   `tfschema:"functions_extension_version"`
	ForceDisableContentShare      bool                                     `tfschema:"content_share_force_disabled"`
	HttpsOnly                     bool                                     `tfschema:"https_only"`
	Identity                      []helpers.Identity                       `tfschema:"identity"`
	SiteConfig                    []helpers.SiteConfigLinuxFunctionAppSlot `tfschema:"site_config"`
	Tags                          map[string]string                        `tfschema:"tags"`
	CustomDomainVerificationId    string                                   `tfschema:"custom_domain_verification_id"`
	DefaultHostname               string                                   `tfschema:"default_hostname"`
	Kind                          string                                   `tfschema:"kind"`
	OutboundIPAddresses           string                                   `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string                                 `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string                                   `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string                                 `tfschema:"possible_outbound_ip_address_list"`
	SiteCredentials               []helpers.SiteCredential                 `tfschema:"site_credential"`
}

var _ sdk.ResourceWithUpdate = LinuxFunctionAppSlotResource{}

func (r LinuxFunctionAppSlotResource) ModelObject() interface{} {
	return &LinuxFunctionAppSlotModel{}
}

func (r LinuxFunctionAppSlotResource) ResourceType() string {
	return "azurerm_linux_function_app_slot"
}

func (r LinuxFunctionAppSlotResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.FunctionAppSlotID
}

func (r LinuxFunctionAppSlotResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
			Description:  "Specifies the name of the Function App Slot.",
		},

		"function_app_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.WebAppName,
			Description:  "The name of the Windows Function App this Slot is a member of.",
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": location.Schema(),

		"service_plan_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ServicePlanID,
			Description:  "The ID of the App Service Plan within which to create this Function App Slot.",
		},

		"storage_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: storageValidate.StorageAccountName,
			Description:  "The backend storage account name which will be used by this Function App Slot.",
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
			Description: "The access key which will be used to access the storage account for the Function App Slot.",
		},

		"storage_uses_managed_identity": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ExactlyOneOf: []string{
				"storage_uses_managed_identity",
				"storage_account_access_key",
			},
			Description: "Should the Function App Slot use its Managed Identity to access storage.",
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
			Description: "Should built in logging be enabled. Configures `AzureWebJobsDashboard` app setting based on the configured storage setting.",
		},

		"client_certificate_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Should the Function App Slot use Client Certificates.",
		},

		"client_certificate_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  web.ClientCertModeOptional,
			ValidateFunc: validation.StringInSlice([]string{
				string(web.ClientCertModeOptional),
				string(web.ClientCertModeRequired),
				string(web.ClientCertModeOptionalInteractiveUser),
			}, false),
			Description: "The mode of the Function App Slot's client certificates requirement for incoming requests. Possible values are `Required`, `Optional`, and `OptionalInteractiveUser`.",
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
			Description: "Is the Linux Function App Slot enabled.",
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
			Description: "The runtime version associated with the Function App Slot.",
		},

		"https_only": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Can the Function App Slot only be accessed via HTTPS?",
		},

		"identity": helpers.IdentitySchema(),

		"site_config": helpers.SiteConfigSchemaLinuxFunctionAppSlot(),

		"tags": tags.Schema(),
	}
}

func (r LinuxFunctionAppSlotResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r LinuxFunctionAppSlotResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Create Func
			// TODO - Don't forget to set the ID! e.g. metadata.SetID(id)
			return nil
		},
	}
}

func (r LinuxFunctionAppSlotResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 25 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := parse.FunctionAppSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			functionApp, err := client.GetSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				if utils.ResponseWasNotFound(functionApp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Linux %s: %+v", id, err)
			}

			if functionApp.SiteProperties == nil {
				return fmt.Errorf("reading properties of Linux %s", id)
			}
			props := *functionApp.SiteProperties

			appSettingsResp, err := client.ListApplicationSettingsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading App Settings for Linux %s: %+v", id, err)
			}

			connectionStrings, err := client.ListConnectionStringsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Linux %s: %+v", id, err)
			}

			siteCredentialsFuture, err := client.ListPublishingCredentialsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
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

			auth, err := client.GetAuthSettingsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Linux %s: %+v", id, err)
			}

			backup, err := client.GetBackupConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				if !utils.ResponseWasNotFound(backup.Response) {
					return fmt.Errorf("reading Backup Settings for Linux %s: %+v", id, err)
				}
			}

			logs, err := client.GetDiagnosticLogsConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading logs configuration for Linux %s: %+v", id, err)
			}

			state := LinuxFunctionAppSlotModel{
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

			configResp, err := client.GetConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("making Read request on AzureRM Function App Configuration %q: %+v", id.SiteName, err)
			}

			siteConfig, err := helpers.FlattenSiteConfigLinuxFunctionAppSlot(configResp.SiteConfig)
			if err != nil {
				return fmt.Errorf("reading Site Config for Linux %s: %+v", id, err)
			}
			state.SiteConfig = []helpers.SiteConfigLinuxFunctionAppSlot{*siteConfig}

			state.unpackLinuxFunctionAppSettings(appSettingsResp, metadata)

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

func (r LinuxFunctionAppSlotResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := parse.FunctionAppSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting Linux %s", *id)

			deleteMetrics := true
			deleteEmptyServerFarm := false
			if _, err := client.DeleteSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName, &deleteMetrics, &deleteEmptyServerFarm); err != nil {
				return fmt.Errorf("deleting Linux %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r LinuxFunctionAppSlotResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Update Func
			return nil
		},
	}
}

func (m *LinuxFunctionAppSlotModel) unpackLinuxFunctionAppSettings(input web.StringDictionary, metadata sdk.ResourceMetaData) {
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

		case "WEBSITE_NODE_DEFAULT_VERSION": // Note - This is only set if it's not the default of 12, but we collect it from LinuxFxVersion so can discard it here
		case "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING":
			if _, ok := metadata.ResourceData.GetOk("app_settings.WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"); ok {
				appSettings[k] = utils.NormalizeNilableString(v)
			}
		case "WEBSITE_CONTENTSHARE":
			if _, ok := metadata.ResourceData.GetOk("app_settings.WEBSITE_CONTENTSHARE"); ok {
				appSettings[k] = utils.NormalizeNilableString(v)
			}
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

		// case "WEBSITES_ENABLE_APP_SERVICE_STORAGE": // TODO - Support this as a configurable bool, default `false` - Ref: https://docs.microsoft.com/en-us/azure/app-service/faq-app-service-linux#i-m-using-my-own-custom-container--i-want-the-platform-to-mount-an-smb-share-to-the---home---directory-

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

	if dockerSettings.RegistryURL != "" {
		appStack := make([]helpers.ApplicationStackLinuxFunctionApp, 0)
		docker, _ := helpers.DecodeFunctionAppDockerFxString(m.SiteConfig[0].LinuxFxVersion, dockerSettings)
		appStack = append(appStack, helpers.ApplicationStackLinuxFunctionApp{Docker: docker})
		m.SiteConfig[0].ApplicationStack = appStack
	}

	m.AppSettings = appSettings
}
