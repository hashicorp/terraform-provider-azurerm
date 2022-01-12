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

type WindowsFunctionAppSlotResource struct{}

type WindowsFunctionAppSlotModel struct {
	Name                          string                                     `tfschema:"name"`
	ResourceGroup                 string                                     `tfschema:"resource_group_name"`
	FunctionAppName               string                                     `tfschema:"function_app_name"`
	Location                      string                                     `tfschema:"location"`
	ServicePlanId                 string                                     `tfschema:"service_plan_id"`
	StorageAccountName            string                                     `tfschema:"storage_account_name"`
	StorageAccountKey             string                                     `tfschema:"storage_account_access_key"`
	StorageUsesMSI                bool                                       `tfschema:"storage_uses_managed_identity"` // Storage uses MSI not account key
	AppSettings                   map[string]string                          `tfschema:"app_settings"`
	AuthSettings                  []helpers.AuthSettings                     `tfschema:"auth_settings"`
	Backup                        []helpers.Backup                           `tfschema:"backup"` // Not supported on Dynamic or Basic plans
	BuiltinLogging                bool                                       `tfschema:"builtin_logging_enabled"`
	ClientCertEnabled             bool                                       `tfschema:"client_certificate_enabled"`
	ClientCertMode                string                                     `tfschema:"client_certificate_mode"`
	ConnectionStrings             []helpers.ConnectionString                 `tfschema:"connection_string"`
	DailyMemoryTimeQuota          int                                        `tfschema:"daily_memory_time_quota"`
	Enabled                       bool                                       `tfschema:"enabled"`
	FunctionExtensionsVersion     string                                     `tfschema:"functions_extension_version"`
	ForceDisableContentShare      bool                                       `tfschema:"content_share_force_disabled"`
	HttpsOnly                     bool                                       `tfschema:"https_only"`
	Identity                      []helpers.Identity                         `tfschema:"identity"`
	SiteConfig                    []helpers.SiteConfigWindowsFunctionAppSlot `tfschema:"site_config"`
	Tags                          map[string]string                          `tfschema:"tags"`
	CustomDomainVerificationId    string                                     `tfschema:"custom_domain_verification_id"`
	DefaultHostname               string                                     `tfschema:"default_hostname"`
	Kind                          string                                     `tfschema:"kind"`
	OutboundIPAddresses           string                                     `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string                                   `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string                                     `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string                                   `tfschema:"possible_outbound_ip_address_list"`
	SiteCredentials               []helpers.SiteCredential                   `tfschema:"site_credential"`
}

var _ sdk.ResourceWithUpdate = WindowsFunctionAppSlotResource{}

func (r WindowsFunctionAppSlotResource) ModelObject() interface{} {
	return &WindowsFunctionAppSlotModel{}
}

func (r WindowsFunctionAppSlotResource) ResourceType() string {
	return "azurerm_windows_function_app_slot"
}

func (r WindowsFunctionAppSlotResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.FunctionAppSlotID
}

func (r WindowsFunctionAppSlotResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
			Description:  "Specifies the name of the Windows Function App Slot.",
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
			Description: "Is the Windows Function App Slot enabled.",
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

		"site_config": helpers.SiteConfigSchemaWindowsFunctionAppSlot(),

		"tags": tags.Schema(),
	}
}

func (r WindowsFunctionAppSlotResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"custom_domain_verification_id": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "The identifier used by App Service to perform domain ownership verification via DNS TXT record.",
		},

		"default_hostname": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The default hostname of the Windows Function App Slot.",
		},

		"kind": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Kind value for this Windows Function App Slot.",
		},

		"outbound_ip_addresses": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "A comma separated list of outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12`.",
		},

		"outbound_ip_address_list": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Description: "A list of outbound IP addresses. For example `[\"52.23.25.3\", \"52.143.43.12\"]`.",
		},

		"possible_outbound_ip_addresses": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "A comma separated list of possible outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12,52.143.43.17`. This is a superset of `outbound_ip_addresses`. For example `[\"52.23.25.3\", \"52.143.43.12\",\"52.143.43.17\"]`.",
		},

		"possible_outbound_ip_address_list": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Description: "A list of possible outbound IP addresses, not all of which are necessarily in use. This is a superset of `outbound_ip_address_list`. For example `[\"52.23.25.3\", \"52.143.43.12\"]`.",
		},

		"site_credential": helpers.SiteCredentialSchema(),
	}
}

func (r WindowsFunctionAppSlotResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Create Func
			// TODO - Don't forget to set the ID! e.g. metadata.SetID(id)
			return nil
		},
	}
}

func (r WindowsFunctionAppSlotResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 25 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := parse.FunctionAppSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			functionAppSlot, err := client.GetSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				if utils.ResponseWasNotFound(functionAppSlot.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Windows %s: %+v", id, err)
			}

			if functionAppSlot.SiteProperties == nil {
				return fmt.Errorf("reading properties of Windows %s", id)
			}
			props := *functionAppSlot.SiteProperties

			appSettingsResp, err := client.ListApplicationSettingsSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading App Settings for Windows %s: %+v", id, err)
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

			logs, err := client.GetDiagnosticLogsConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("reading logs configuration for Windows %s: %+v", id, err)
			}

			state := WindowsFunctionAppSlotModel{
				Name:                 id.SiteName,
				ResourceGroup:        id.ResourceGroup,
				ServicePlanId:        utils.NormalizeNilableString(props.ServerFarmID),
				Location:             location.NormalizeNilable(functionAppSlot.Location),
				Enabled:              utils.NormaliseNilableBool(functionAppSlot.Enabled),
				ClientCertMode:       string(functionAppSlot.ClientCertMode),
				DailyMemoryTimeQuota: int(utils.NormaliseNilableInt32(props.DailyMemoryTimeQuota)),
				Tags:                 tags.ToTypedObject(functionAppSlot.Tags),
				Kind:                 utils.NormalizeNilableString(functionAppSlot.Kind),
			}

			if identity := helpers.FlattenIdentity(functionAppSlot.Identity); identity != nil {
				state.Identity = identity
			}

			configResp, err := client.GetConfigurationSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
			if err != nil {
				return fmt.Errorf("making Read request on AzureRM Function App Configuration %q: %+v", id.SiteName, err)
			}

			siteConfig, err := helpers.FlattenSiteConfigWindowsFunctionAppSlot(configResp.SiteConfig)
			if err != nil {
				return fmt.Errorf("reading Site Config for Windows %s: %+v", id, err)
			}
			state.SiteConfig = []helpers.SiteConfigWindowsFunctionAppSlot{*siteConfig}

			state.unpackWindowsFunctionAppSettings(appSettingsResp)

			state.ConnectionStrings = helpers.FlattenConnectionStrings(connectionStrings)

			state.SiteCredentials = helpers.FlattenSiteCredentials(siteCredentials)

			state.AuthSettings = helpers.FlattenAuthSettings(auth)

			state.Backup = helpers.FlattenBackupConfig(backup)

			state.SiteConfig[0].AppServiceLogs = helpers.FlattenFunctionAppAppServiceLogs(logs)

			state.HttpsOnly = utils.NormaliseNilableBool(functionAppSlot.HTTPSOnly)
			state.ClientCertEnabled = utils.NormaliseNilableBool(functionAppSlot.ClientCertEnabled)

			return metadata.Encode(&state)
		},
	}
}

func (r WindowsFunctionAppSlotResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Delete Func
			return nil
		},
	}
}

func (r WindowsFunctionAppSlotResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Update Func
			return nil
		},
	}
}

func (m *WindowsFunctionAppSlotModel) unpackWindowsFunctionAppSettings(input web.StringDictionary) {
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
