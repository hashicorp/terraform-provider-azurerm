package appservice

import (
	"context"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WindowsFunctionAppSlotResource struct{}

type WindowsFunctionAppSlotModel struct {
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
	ClientCertEnabled         bool                                   `tfschema:"client_certificate_enabled"`
	ClientCertMode            string                                 `tfschema:"client_certificate_mode"`
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

var _ sdk.ResourceWithUpdate = WindowsFunctionAppSlotResource{}

func (r WindowsFunctionAppSlotResource) ModelObject() interface{} {
	return &WindowsFunctionAppSlotModel{}
}

func (r WindowsFunctionAppSlotResource) ResourceType() string {
	return "azurerm_windows_function_app_slot"
}

func (r WindowsFunctionAppSlotResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	panic("Implement me") // TODO - Add Validation func return here
}

func (r WindowsFunctionAppSlotResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		/*
			TODO - This sections is for configurable items, `Required: true` items first, followed by `Optional: true`,
			both in alphabetical order
		*/
	}
}

func (r WindowsFunctionAppSlotResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		/*
			TODO - This section is for `Computed: true` only items, i.e. useful values that are returned by the
			datasource that can be used as outputs or passed programmatically to other resources or data sources.
		*/
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
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Read Func
			return nil
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
