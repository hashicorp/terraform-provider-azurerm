package logic

import (
	"github.com/hashicorp/go-azure-helpers/framework/identity"
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	logicAppKindWorkflowApp               = "functionapp,workflowapp"
	logicAppKindLinuxContainerWorkflowApp = "functionapp,linux,container,workflowapp"
)

const (
	storagePropName                = "AzureWebJobsStorage"
	functionVersionPropName        = "FUNCTIONS_EXTENSION_VERSION"
	contentSharePropName           = "WEBSITE_CONTENTSHARE"
	contentFileConnStringPropName  = "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"
	appKindPropName                = "APP_KIND"
	appKindPropValue               = "workflowApp"
	storageConnectionFmt           = "DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s"
	extensionBundlePropName        = "AzureFunctionsJobHost__extensionBundle__id"
	extensionBundleName            = "Microsoft.Azure.Functions.ExtensionBundle.Workflows"
	extensionBundleVersionPropName = "AzureFunctionsJobHost__extensionBundle__version"
	webJobsDashboardPropName       = "AzureWebJobsDashboard"
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

type FwLogicAppStandardResourceIdentityModel struct {
	SubscriptionId    string `tfsdk:"subscription_id"`
	ResourceGroupName string `tfsdk:"resource_group_name"`
	Name              string `tfsdk:"name"`
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
type FwLogicAppCORSSettingsModel struct {
	AllowedOrigins     typehelpers.SetValueOf[types.String] `tfsdk:"allowed_origins"`
	SupportCredentials types.Bool                           `tfsdk:"support_credentials"`
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

type FwLogicAppIPRestrictionHeadersModel struct {
	XForwardedHost typehelpers.SetValueOf[types.String]  `tfsdk:"x_forwarded_host"`
	XForwardedFor  typehelpers.SetValueOf[types.String]  `tfsdk:"x_forwarded_for"`
	XAzureFDID     typehelpers.SetValueOf[types.String]  `tfsdk:"x_azure_fdid"`
	XFDHealthProbe typehelpers.ListValueOf[types.String] `tfsdk:"x_fd_health_probe"`
}

type FwLogicAppStandardSiteCredentialsModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type FwLogicAppStandardConnectionStringsModel struct {
	Name  types.String `tfsdk:"name"`
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}
