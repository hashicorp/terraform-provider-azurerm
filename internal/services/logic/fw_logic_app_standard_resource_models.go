package logic

import (
	"context"
	"fmt"
	"math"

	"github.com/hashicorp/go-azure-helpers/framework/identity"
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	logicAppKindWorkflowApp               = "functionapp,workflowapp"
	logicAppKindLinuxContainerWorkflowApp = "functionapp,linux,container,workflowapp"
)

const (
	storagePropName               = "AzureWebJobsStorage"
	functionVersionPropName       = "FUNCTIONS_EXTENSION_VERSION"
	contentSharePropName          = "WEBSITE_CONTENTSHARE"
	contentFileConnStringPropName = "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"
	appKindPropName               = "APP_KIND"
	appKindPropValue              = "workflowApp"
)

// FwLogicAppStandardResourceModel is a temporary struct - please use make generate function after creating/updating the schema function for the resource to correctly populate this and generate all modules for the schema
type FwLogicAppStandardResourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`

	Name                        types.String                                                                 `tfsdk:"name"`
	ResourceGroupName           types.String                                                                 `tfsdk:"resource_group_name"`
	Location                    types.String                                                                 `tfsdk:"location"`
	Identity                    typehelpers.ListNestedObjectValueOf[identity.IdentityModel]                  `tfsdk:"identity"`
	AppServicePlanID            types.String                                                                 `tfsdk:"app_service_plan_id"`
	AppSettings                 typehelpers.MapValueOf[types.String]                                         `tfsdk:"app_settings"`
	UseExtensionBundle          types.Bool                                                                   `tfsdk:"use_extension_bundle"`
	BundleVersion               types.String                                                                 `tfsdk:"bundle_version"`
	ClientAffinityEnabled       types.Bool                                                                   `tfsdk:"client_affinity_enabled"`
	ClientCertificateMode       types.String                                                                 `tfsdk:"client_certificate_mode"`
	ConnectionStrings           typehelpers.SetNestedObjectValueOf[FwLogicAppStandardConnectionStringsModel] `tfsdk:"connection_string"`
	Enabled                     types.Bool                                                                   `tfsdk:"enabled"`
	HTTPSOnly                   types.Bool                                                                   `tfsdk:"https_only"`
	StorageAccountName          types.String                                                                 `tfsdk:"storage_account_name"`
	StorageAccountAccessKey     types.String                                                                 `tfsdk:"storage_account_access_key"`
	PublicNetworkAccess         types.String                                                                 `tfsdk:"public_network_access"`
	StorageAccountShareName     types.String                                                                 `tfsdk:"storage_account_share_name"`
	Version                     types.String                                                                 `tfsdk:"version"`
	VirtualNetworkSubnetID      types.String                                                                 `tfsdk:"virtual_network_subnet_id"`
	CustomDomainVerificationID  types.String                                                                 `tfsdk:"custom_domain_verification_id"`
	DefaultHostname             types.String                                                                 `tfsdk:"default_hostname"`
	Kind                        types.String                                                                 `tfsdk:"kind"`
	OutboundIPAddresses         types.String                                                                 `tfsdk:"outbound_ip_addresses"`
	PossibleOutboundIPAddresses types.String                                                                 `tfsdk:"possible_outbound_ip_addresses"`
	Tags                        typehelpers.MapValueOf[types.String]                                         `tfsdk:"tags"`
	SiteConfig                  typehelpers.ListNestedObjectValueOf[FwLogicAppStandardSiteConfigModel]       `tfsdk:"site_config"`
	SiteCredentials             typehelpers.ListNestedObjectValueOf[FwLogicAppStandardSiteCredentialsModel]  `tfsdk:"site_credential"`
}

type FwLogicAppStandardResourceNativeModel struct {
	ID       string
	Timeouts typehelpers.Timeouts

	Name              string
	ResourceGroupName string
	Location          string
	// Identity                        typehelpers.ListNestedObjectValueOf[identity.Identity]                     `tfsdk:"identity"`
	AppServicePlanID                string
	AppSettings                     map[string]string
	UseExtensionBundle              bool
	BundleVersion                   string
	ClientAffinityEnabled           bool
	ClientCertificateMode           string
	ConnectionStrings               []FwLogicAppStandardConnectionStringsNativeModel
	Enabled                         bool
	HTTPSOnly                       bool
	StorageAccountName              string
	StorageAccountAccessKey         string
	PublicNetworkAccess             string
	StorageAccountShareName         string
	Version                         string
	VirtualNetworkSubnetID          string
	CustomDomainVerificationID      string
	DefaultHostname                 string
	Kind                            string
	OutboundIPAddresses             string
	PossibleOutboundIPAddresses     string
	OutboundIPAddressesList         []string
	PossibleOutboundIPAddressesList []string
	Tags                            map[string]string
	SiteConfig                      []FwLogicAppStandardSiteConfigNativeModel
	SiteCredentials                 []FwLogicAppStandardSiteCredentialsNativeModel
}

type FwLogicAppStandardSiteConfigModel struct {
	AlwaysOn                      types.Bool                                                        `tfsdk:"always_on"`
	CORS                          typehelpers.ListNestedObjectValueOf[FwLogicAppCORSSettingsModel]  `tfsdk:"cors"`
	FTPSState                     types.String                                                      `tfsdk:"ftps_state"`
	HTTP2Enabled                  types.Bool                                                        `tfsdk:"http2_enabled"`
	IPRestriction                 typehelpers.ListNestedObjectValueOf[FwLogicAppIPRestrictionModel] `tfsdk:"ip_restriction"`
	LinuxFXVersion                types.String                                                      `tfsdk:"linux_fx_version"`
	MinTLSVersion                 types.String                                                      `tfsdk:"min_tls_version"`
	PreWarmedInstanceCount        types.Int64                                                       `tfsdk:"pre_warmed_instance_count"`
	SCMIPRestriction              typehelpers.ListNestedObjectValueOf[FwLogicAppIPRestrictionModel] `tfsdk:"scm_ip_restriction"`
	SCMUseMainIPRestriction       types.Bool                                                        `tfsdk:"scm_use_main_ip_restriction"`
	SCMMinTLSVersion              types.String                                                      `tfsdk:"scm_min_tls_version"`
	SCMType                       types.String                                                      `tfsdk:"scm_type"`
	Use32BitWorkerProcesses       types.Bool                                                        `tfsdk:"use_32_bit_worker_process"`
	WebSocketsEnabled             types.Bool                                                        `tfsdk:"websockets_enabled"`
	HealthCheckPath               types.String                                                      `tfsdk:"health_check_path"`
	ElasticInstanceMinimum        types.Int64                                                       `tfsdk:"elastic_instance_minimum"`
	AppScaleLimit                 types.Int64                                                       `tfsdk:"app_scale_limit"`
	RuntimeScaleMonitoringEnabled types.Bool                                                        `tfsdk:"runtime_scale_monitoring_enabled"`
	DotnetFrameworkVersion        types.String                                                      `tfsdk:"dotnet_framework_version"`
	VNETRouteAllEnabled           types.Bool                                                        `tfsdk:"vnet_route_all_enabled"`
	AutoSwapSlotName              types.String                                                      `tfsdk:"auto_swap_slot_name"`
}
type FwLogicAppStandardSiteConfigNativeModel struct {
	AlwaysOn                      bool
	CORS                          []FwLogicAppCORSSettingsNativeModel
	FTPSState                     string
	HTTP2Enabled                  bool
	IPRestriction                 []FwLogicAppIPRestrictionNativeModel
	LinuxFXVersion                string
	MinTLSVersion                 string
	PreWarmedInstanceCount        int64
	SCMIPRestriction              []FwLogicAppIPRestrictionNativeModel
	SCMUseMainIPRestriction       bool
	SCMMinTLSVersion              string
	SCMType                       string
	Use32BitWorkerProcesses       bool
	WebSocketsEnabled             bool
	HealthCheckPath               string
	ElasticInstanceMinimum        int64
	AppScaleLimit                 int64
	RuntimeScaleMonitoringEnabled bool
	DotnetFrameworkVersion        string
	VNETRouteAllEnabled           bool
	AutoSwapSlotName              string
}

type FwLogicAppCORSSettingsModel struct {
	AllowedOrigins     typehelpers.SetValueOf[types.String] `tfsdk:"allowed_origins"`
	SupportCredentials types.Bool                           `tfsdk:"support_credentials"`
}

type FwLogicAppCORSSettingsNativeModel struct {
	AllowedOrigins     []string
	SupportCredentials bool
}

type FwLogicAppIPRestrictionModel struct {
	IPAddress              types.String                                                             `tfsdk:"ip_address"`
	ServiceTag             types.String                                                             `tfsdk:"service_tag"`
	VirtualNetworkSubnetID types.String                                                             `tfsdk:"virtual_network_subnet_id"`
	Name                   types.String                                                             `tfsdk:"name"`
	Priority               types.Int64                                                              `tfsdk:"priority"`
	Action                 types.String                                                             `tfsdk:"action"`
	Headers                typehelpers.ListNestedObjectValueOf[FwLogicAppIPRestrictionHeadersModel] `tfsdk:"headers"`
}

type FwLogicAppIPRestrictionNativeModel struct {
	IPAddress              string
	ServiceTag             string
	VirtualNetworkSubnetID string
	Name                   string
	Priority               int64
	Action                 string
	Headers                []FwLogicAppIPRestrictionHeadersNativeModel
}

type FwLogicAppIPRestrictionHeadersModel struct {
	XForwardedHost typehelpers.SetValueOf[types.String]  `tfsdk:"x_forwarded_host"`
	XForwardedFor  typehelpers.SetValueOf[types.String]  `tfsdk:"x_forwarded_for"`
	XAzureFDID     typehelpers.SetValueOf[types.String]  `tfsdk:"x_azure_fdid"`
	XFDHealthProbe typehelpers.ListValueOf[types.String] `tfsdk:"x_fd_health_probe"`
}

type FwLogicAppIPRestrictionHeadersNativeModel struct {
	XForwardedHost []string
	XForwardedFor  []string
	XAzureFDID     []string
	XFDHealthProbe []string
}

type FwLogicAppStandardSiteCredentialsModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type FwLogicAppStandardSiteCredentialsNativeModel struct {
	Username string
	Password string
}

type FwLogicAppStandardConnectionStringsModel struct {
	Name  types.String `tfsdk:"name"`
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

type FwLogicAppStandardConnectionStringsNativeModel struct {
	Name  string
	Type  string
	value string
}

func ipRestrictionCommonSchema(ctx context.Context) schema.ListNestedBlock {
	return schema.ListNestedBlock{
		CustomType: typehelpers.NewListNestedObjectTypeOf[FwLogicAppIPRestrictionModel](ctx),
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"ip_address": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},

				"service_tag": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},

				"virtual_network_subnet_id": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},

				"name": schema.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},

				"priority": schema.Int64Attribute{
					Optional: true,
					Computed: true,
					Default:  typehelpers.NewWrappedInt64Default(65000),
					Validators: []validator.Int64{
						int64validator.Between(1, math.MaxInt32),
					},
				},

				"action": schema.StringAttribute{
					Default:  typehelpers.NewWrappedStringDefault("Allow"),
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						typehelpers.WrappedStringValidator{
							Func: validation.StringInSlice([]string{
								"Allow",
								"Deny",
							}, false),
						},
					},
				},
			},
			Blocks: map[string]schema.Block{
				"headers": schema.ListNestedBlock{
					CustomType: typehelpers.NewListNestedObjectTypeOf[FwLogicAppIPRestrictionHeadersModel](ctx),
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"x_forwarded_host": schema.SetAttribute{
								CustomType:  typehelpers.NewSetTypeOf[types.String](ctx),
								ElementType: types.StringType,
								Optional:    true,
								Validators: []validator.Set{
									setvalidator.SizeAtMost(8),
								},
							},

							"x_forwarded_for": schema.SetAttribute{
								CustomType:  typehelpers.NewSetTypeOf[types.String](ctx),
								ElementType: types.StringType,
								Optional:    true,
								Validators: []validator.Set{
									setvalidator.All(
										setvalidator.SizeAtMost(8),
										setvalidator.ValueStringsAre(
											typehelpers.WrappedStringValidator{
												Func: validation.IsCIDR,
											},
										),
									),
								},
							},

							"x_azure_fdid": schema.SetAttribute{
								CustomType:  typehelpers.NewSetTypeOf[types.String](ctx),
								ElementType: types.StringType,
								Optional:    true,
								Validators: []validator.Set{
									setvalidator.All(
										setvalidator.SizeAtMost(8),
										setvalidator.ValueStringsAre(
											typehelpers.WrappedStringValidator{
												Func: validation.IsUUID,
											},
										),
									),
								},
							},

							"x_fd_health_probe": schema.ListAttribute{
								CustomType:  typehelpers.NewListTypeOf[types.String](ctx),
								ElementType: types.StringType,
								Optional:    true,
								Validators: []validator.List{
									listvalidator.All(
										listvalidator.SizeAtMost(1),
										listvalidator.ValueStringsAre(
											stringvalidator.OneOf("0", "1"),
										),
									),
								},
							},
						},
					},
				},
			},
		},
	}
}

func assertFwLogicAppStandardResourceModel(input interface{}, response interface{}) *FwLogicAppStandardResourceModel {
	result, ok := input.(*FwLogicAppStandardResourceModel)
	if !ok {
		switch v := response.(type) {
		case *resource.CreateResponse:
			v.Diagnostics.AddError("resource had wrong model type", fmt.Sprintf("got %T", input))
		case *resource.ReadResponse:
			v.Diagnostics.AddError("resource had wrong model type", fmt.Sprintf("got %T", input))
		case *resource.UpdateResponse:
			v.Diagnostics.AddError("resource had wrong model type, ", fmt.Sprintf("got %T", input))
		case *resource.DeleteResponse:
			v.Diagnostics.AddError("resource had wrong model type", fmt.Sprintf("got %T", input))
		}
		return nil
	}

	return result
}
