package appservice

import (
	"context"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WindowsWebAppSlotResource struct{}

type WindowsWebAppSlotModel struct {
	Name                          string                                `tfschema:"name"`
	AppServiceName                string                                `tfschema:"app_service_name"`
	ResourceGroup                 string                                `tfschema:"resource_group_name"`
	Location                      string                                `tfschema:"location"`
	ServicePlanId                 string                                `tfschema:"service_plan_id"`
	AppSettings                   map[string]string                     `tfschema:"app_settings"`
	AuthSettings                  []helpers.AuthSettings                `tfschema:"auth_settings"`
	Backup                        []helpers.Backup                      `tfschema:"backup"`
	ClientAffinityEnabled         bool                                  `tfschema:"client_affinity_enabled"`
	ClientCertEnabled             bool                                  `tfschema:"client_certificate_enabled"`
	ClientCertMode                string                                `tfschema:"client_certificate_mode"`
	Enabled                       bool                                  `tfschema:"enabled"`
	HttpsOnly                     bool                                  `tfschema:"https_only"`
	Identity                      []helpers.Identity                    `tfschema:"identity"`
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
	Tags                          map[string]string                     `tfschema:"tags"`
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
		/*
			TODO - This sections is for configurable items, `Required: true` items first, followed by `Optional: true`,
			both in alphabetical order
		*/
	}
}

func (r WindowsWebAppSlotResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		/*
			TODO - This section is for `Computed: true` only items, i.e. useful values that are returned by the
			datasource that can be used as outputs or passed programmatically to other resources or data sources.
		*/
	}
}

func (r WindowsWebAppSlotResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Create Func
			// TODO - Don't forget to set the ID! e.g. metadata.SetID(id)
			return nil
		},
	}
}

func (r WindowsWebAppSlotResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Read Func
			return nil
		},
	}
}

func (r WindowsWebAppSlotResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Delete Func
			return nil
		},
	}
}

func (r WindowsWebAppSlotResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Update Func
			return nil
		},
	}
}
