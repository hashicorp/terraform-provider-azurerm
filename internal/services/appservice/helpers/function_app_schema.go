// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	apimValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const (
	StorageStringFmt   = "DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s"
	StorageStringFmtKV = "@Microsoft.KeyVault(SecretUri=%s)"
)

type SiteConfigLinuxFunctionApp struct {
	AlwaysOn                      bool                               `tfschema:"always_on"`
	AppCommandLine                string                             `tfschema:"app_command_line"`
	ApiDefinition                 string                             `tfschema:"api_definition_url"`
	ApiManagementConfigId         string                             `tfschema:"api_management_api_id"`
	AppInsightsInstrumentationKey string                             `tfschema:"application_insights_key"` // App Insights Instrumentation Key
	AppInsightsConnectionString   string                             `tfschema:"application_insights_connection_string"`
	AppScaleLimit                 int64                              `tfschema:"app_scale_limit"`
	AppServiceLogs                []FunctionAppAppServiceLogs        `tfschema:"app_service_logs"`
	UseManagedIdentityACR         bool                               `tfschema:"container_registry_use_managed_identity"`
	ContainerRegistryMSI          string                             `tfschema:"container_registry_managed_identity_client_id"`
	DefaultDocuments              []string                           `tfschema:"default_documents"`
	ElasticInstanceMinimum        int64                              `tfschema:"elastic_instance_minimum"`
	Http2Enabled                  bool                               `tfschema:"http2_enabled"`
	IpRestriction                 []IpRestriction                    `tfschema:"ip_restriction"`
	IpRestrictionDefaultAction    string                             `tfschema:"ip_restriction_default_action"`
	LoadBalancing                 string                             `tfschema:"load_balancing_mode"` // TODO - Valid for FunctionApps?
	ManagedPipelineMode           string                             `tfschema:"managed_pipeline_mode"`
	PreWarmedInstanceCount        int64                              `tfschema:"pre_warmed_instance_count"`
	RemoteDebugging               bool                               `tfschema:"remote_debugging_enabled"`
	RemoteDebuggingVersion        string                             `tfschema:"remote_debugging_version"`
	RuntimeScaleMonitoring        bool                               `tfschema:"runtime_scale_monitoring_enabled"`
	ScmIpRestriction              []IpRestriction                    `tfschema:"scm_ip_restriction"`
	ScmIpRestrictionDefaultAction string                             `tfschema:"scm_ip_restriction_default_action"`
	ScmType                       string                             `tfschema:"scm_type"` // Computed?
	ScmUseMainIpRestriction       bool                               `tfschema:"scm_use_main_ip_restriction"`
	Use32BitWorker                bool                               `tfschema:"use_32_bit_worker"`
	WebSockets                    bool                               `tfschema:"websockets_enabled"`
	FtpsState                     string                             `tfschema:"ftps_state"`
	HealthCheckPath               string                             `tfschema:"health_check_path"`
	HealthCheckEvictionTime       int64                              `tfschema:"health_check_eviction_time_in_min"`
	WorkerCount                   int64                              `tfschema:"worker_count"`
	ApplicationStack              []ApplicationStackLinuxFunctionApp `tfschema:"application_stack"`
	MinTlsVersion                 string                             `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion              string                             `tfschema:"scm_minimum_tls_version"`
	Cors                          []CorsSetting                      `tfschema:"cors"`
	DetailedErrorLogging          bool                               `tfschema:"detailed_error_logging_enabled"`
	LinuxFxVersion                string                             `tfschema:"linux_fx_version"`
	VnetRouteAllEnabled           bool                               `tfschema:"vnet_route_all_enabled"` // Not supported in Dynamic plans
}

func SiteConfigSchemaLinuxFunctionApp() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Computed:    true, // Note - several factors change the default for this, so needs to be computed.
					Description: "If this Linux Web App is Always On enabled. Defaults to `false`.",
				},

				"api_management_api_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: apimValidate.ApiID,
					Description:  "The ID of the API Management API for this Linux Function App.",
				},

				"api_definition_url": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					Description:  "The URL of the API definition that describes this Linux Function App.",
				},

				"app_command_line": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The program and any arguments used to launch this app via the command line. (Example `node myapp.js`).",
				},

				"app_scale_limit": {
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: "The number of workers this function app can scale out to. Only applicable to apps on the Consumption and Premium plan.",
					// TODO Validation?
				},

				"application_insights_key": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The Instrumentation Key for connecting the Linux Function App to Application Insights.",
				},

				"application_insights_connection_string": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The Connection String for linking the Linux Function App to Application Insights.",
				},

				"application_stack": linuxFunctionAppStackSchema(),

				"app_service_logs": FunctionAppAppServiceLogsSchema(),

				"container_registry_use_managed_identity": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should connections for Azure Container Registry use Managed Identity.",
				},

				"container_registry_managed_identity_client_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsUUID,
					Description:  "The Client ID of the Managed Service Identity to use for connections to the Azure Container Registry.",
				},

				"default_documents": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of Default Documents for the Linux Web App.",
				},

				"elastic_instance_minimum": {
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: "The number of minimum instances for this Linux Function App. Only affects apps on Elastic Premium plans.",
				},

				"http2_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Specifies if the http2 protocol should be enabled. Defaults to `false`.",
				},

				"ip_restriction": IpRestrictionSchema(),

				"ip_restriction_default_action": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      webapps.DefaultActionAllow,
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForDefaultAction(), false),
				},

				"scm_use_main_ip_restriction": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should the Linux Function App `ip_restriction` configuration be used for the SCM also.",
				},

				"scm_ip_restriction": IpRestrictionSchema(),

				"scm_ip_restriction_default_action": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      webapps.DefaultActionAllow,
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForDefaultAction(), false),
				},

				"load_balancing_mode": { // Supported on Function Apps?
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "LeastRequests",
					ValidateFunc: validation.StringInSlice([]string{
						"LeastRequests", // Service default
						"WeightedRoundRobin",
						"LeastResponseTime",
						"WeightedTotalTraffic",
						"RequestHash",
						"PerSiteRoundRobin",
					}, false),
					Description: "The Site load balancing mode. Possible values include: `WeightedRoundRobin`, `LeastRequests`, `LeastResponseTime`, `WeightedTotalTraffic`, `RequestHash`, `PerSiteRoundRobin`. Defaults to `LeastRequests` if omitted.",
				},

				"managed_pipeline_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(webapps.ManagedPipelineModeIntegrated),
					ValidateFunc: validation.StringInSlice([]string{
						string(webapps.ManagedPipelineModeClassic),
						string(webapps.ManagedPipelineModeIntegrated),
					}, false),
					Description: "The Managed Pipeline mode. Possible values include: `Integrated`, `Classic`. Defaults to `Integrated`.",
				},

				"pre_warmed_instance_count": {
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Computed:    true, // Variable defaults depending on plan etc
					Description: "The number of pre-warmed instances for this function app. Only affects apps on an Elastic Premium plan.",
				},

				"remote_debugging_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should Remote Debugging be enabled. Defaults to `false`.",
				},

				"remote_debugging_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"VS2017",
						"VS2019",
						"VS2022",
					}, false),
					Description: "The Remote Debugging Version. Possible values include `VS2017`, `VS2019`, and `VS2022``",
				},

				"runtime_scale_monitoring_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Description: "Should Functions Runtime Scale Monitoring be enabled.",
				},

				"scm_type": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The SCM Type in use by the Linux Function App.",
				},

				"use_32_bit_worker": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should the Linux Web App use a 32-bit worker.",
				},

				"websockets_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should Web Sockets be enabled. Defaults to `false`.",
				},

				"ftps_state": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(webapps.FtpsStateDisabled),
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForFtpsState(), false),
					Description:  "State of FTP / FTPS service for this function app. Possible values include: `AllAllowed`, `FtpsOnly` and `Disabled`. Defaults to `Disabled`.",
				},

				"health_check_path": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The path to be checked for this function app health.",
					RequiredWith: func() []string {
						if features.FourPointOhBeta() {
							return []string{"site_config.0.health_check_eviction_time_in_min"}
						}
						return []string{}
					}(),
				},

				"health_check_eviction_time_in_min": { // NOTE: Will evict the only node in single node configurations.
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     !features.FourPointOhBeta(),
					ValidateFunc: validation.IntBetween(2, 10),
					RequiredWith: func() []string {
						if features.FourPointOhBeta() {
							return []string{"site_config.0.health_check_path"}
						}
						return []string{}
					}(),
					Description: "The amount of time in minutes that a node is unhealthy before being removed from the load balancer. Possible values are between `2` and `10`. Only valid in conjunction with `health_check_path`",
				},

				"worker_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 100),
					Description:  "The number of Workers for this Linux Function App.",
				},

				"minimum_tls_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(webapps.SupportedTlsVersionsOnePointTwo),
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForSupportedTlsVersions(), false),
					Description:  "The configures the minimum version of TLS required for SSL requests. Possible values include: `1.0`, `1.1`, and  `1.2`. Defaults to `1.2`.",
				},

				"scm_minimum_tls_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(webapps.SupportedTlsVersionsOnePointTwo),
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForSupportedTlsVersions(), false),
					Description:  "Configures the minimum version of TLS required for SSL requests to the SCM site Possible values include: `1.0`, `1.1`, and  `1.2`. Defaults to `1.2`.",
				},

				"cors": CorsSettingsSchema(),

				"vnet_route_all_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should all outbound traffic to have Virtual Network Security Groups and User Defined Routes applied? Defaults to `false`.",
				},

				"detailed_error_logging_enabled": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Is detailed error logging enabled",
				},

				"linux_fx_version": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The Linux FX Version",
				},
			},
		},
	}
}

func SiteConfigSchemaLinuxFunctionAppComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"api_management_api_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"api_definition_url": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"app_command_line": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"app_scale_limit": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"application_insights_key": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"application_insights_connection_string": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"application_stack": linuxFunctionAppStackSchemaComputed(),

				"app_service_logs": FunctionAppAppServiceLogsSchemaComputed(),

				"container_registry_use_managed_identity": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"container_registry_managed_identity_client_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"default_documents": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"elastic_instance_minimum": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"http2_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"ip_restriction": IpRestrictionSchemaComputed(),

				"ip_restriction_default_action": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"scm_ip_restriction": IpRestrictionSchemaComputed(),

				"scm_ip_restriction_default_action": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      webapps.DefaultActionAllow,
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForDefaultAction(), false),
				},

				"load_balancing_mode": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"managed_pipeline_mode": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"pre_warmed_instance_count": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"remote_debugging_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"remote_debugging_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"runtime_scale_monitoring_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"scm_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"use_32_bit_worker": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"websockets_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"ftps_state": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"health_check_path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"health_check_eviction_time_in_min": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"worker_count": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"scm_minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"cors": CorsSettingsSchemaComputed(),

				"vnet_route_all_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"detailed_error_logging_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"linux_fx_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

type SiteConfigWindowsFunctionApp struct {
	AlwaysOn                      bool                                 `tfschema:"always_on"`
	AppCommandLine                string                               `tfschema:"app_command_line"`
	ApiDefinition                 string                               `tfschema:"api_definition_url"`
	ApiManagementConfigId         string                               `tfschema:"api_management_api_id"`
	AppInsightsInstrumentationKey string                               `tfschema:"application_insights_key"` // App Insights Instrumentation Key
	AppInsightsConnectionString   string                               `tfschema:"application_insights_connection_string"`
	AppScaleLimit                 int64                                `tfschema:"app_scale_limit"`
	AppServiceLogs                []FunctionAppAppServiceLogs          `tfschema:"app_service_logs"`
	DefaultDocuments              []string                             `tfschema:"default_documents"`
	ElasticInstanceMinimum        int64                                `tfschema:"elastic_instance_minimum"`
	Http2Enabled                  bool                                 `tfschema:"http2_enabled"`
	IpRestriction                 []IpRestriction                      `tfschema:"ip_restriction"`
	IpRestrictionDefaultAction    string                               `tfschema:"ip_restriction_default_action"`
	LoadBalancing                 string                               `tfschema:"load_balancing_mode"` // TODO - Valid for FunctionApps?
	ManagedPipelineMode           string                               `tfschema:"managed_pipeline_mode"`
	PreWarmedInstanceCount        int64                                `tfschema:"pre_warmed_instance_count"`
	RemoteDebugging               bool                                 `tfschema:"remote_debugging_enabled"`
	RemoteDebuggingVersion        string                               `tfschema:"remote_debugging_version"`
	RuntimeScaleMonitoring        bool                                 `tfschema:"runtime_scale_monitoring_enabled"`
	ScmIpRestriction              []IpRestriction                      `tfschema:"scm_ip_restriction"`
	ScmType                       string                               `tfschema:"scm_type"` // Computed?
	ScmIpRestrictionDefaultAction string                               `tfschema:"scm_ip_restriction_default_action"`
	ScmUseMainIpRestriction       bool                                 `tfschema:"scm_use_main_ip_restriction"`
	Use32BitWorker                bool                                 `tfschema:"use_32_bit_worker"`
	WebSockets                    bool                                 `tfschema:"websockets_enabled"`
	FtpsState                     string                               `tfschema:"ftps_state"`
	HealthCheckPath               string                               `tfschema:"health_check_path"`
	HealthCheckEvictionTime       int64                                `tfschema:"health_check_eviction_time_in_min"`
	NumberOfWorkers               int64                                `tfschema:"worker_count"`
	ApplicationStack              []ApplicationStackWindowsFunctionApp `tfschema:"application_stack"`
	MinTlsVersion                 string                               `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion              string                               `tfschema:"scm_minimum_tls_version"`
	Cors                          []CorsSetting                        `tfschema:"cors"`
	DetailedErrorLogging          bool                                 `tfschema:"detailed_error_logging_enabled"`
	WindowsFxVersion              string                               `tfschema:"windows_fx_version"`
	VnetRouteAllEnabled           bool                                 `tfschema:"vnet_route_all_enabled"` // Not supported in Dynamic plans
}

func SiteConfigSchemaWindowsFunctionApp() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Computed:    true, // Note - several factors change the default for this, so needs to be computed.
					Description: "If this Windows Web App is Always On enabled. Defaults to `false`.",
				},

				"api_management_api_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: apimValidate.ApiID,
					Description:  "The ID of the API Management API for this Windows Function App.",
				},

				"api_definition_url": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					Description:  "The URL of the API definition that describes this Windows Function App.",
				},

				"app_command_line": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The program and any arguments used to launch this app via the command line. (Example `node myapp.js`).",
				},

				"app_scale_limit": {
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: "The number of workers this function app can scale out to. Only applicable to apps on the Consumption and Premium plan.",
					// TODO Validation?
				},

				"application_insights_key": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The Instrumentation Key for connecting the Windows Function App to Application Insights.",
				},

				"application_insights_connection_string": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The Connection String for linking the Windows Function App to Application Insights.",
				},

				"application_stack": windowsFunctionAppStackSchema(),

				"app_service_logs": FunctionAppAppServiceLogsSchema(),

				"default_documents": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of Default Documents for the Windows Web App.",
				},

				"elastic_instance_minimum": {
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: "The number of minimum instances for this Windows Function App. Only affects apps on Elastic Premium plans.",
				},

				"http2_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Specifies if the http2 protocol should be enabled. Defaults to `false`.",
				},

				"ip_restriction": IpRestrictionSchema(),

				"ip_restriction_default_action": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      webapps.DefaultActionAllow,
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForDefaultAction(), false),
				},

				"scm_use_main_ip_restriction": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should the Windows Function App `ip_restriction` configuration be used for the SCM also.",
				},

				"scm_ip_restriction": IpRestrictionSchema(),

				"scm_ip_restriction_default_action": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      webapps.DefaultActionAllow,
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForDefaultAction(), false),
				},

				"load_balancing_mode": { // Supported on Function Apps?
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "LeastRequests",
					ValidateFunc: validation.StringInSlice([]string{
						"LeastRequests", // Service default
						"WeightedRoundRobin",
						"LeastResponseTime",
						"WeightedTotalTraffic",
						"RequestHash",
						"PerSiteRoundRobin",
					}, false),
					Description: "The Site load balancing mode. Possible values include: `WeightedRoundRobin`, `LeastRequests`, `LeastResponseTime`, `WeightedTotalTraffic`, `RequestHash`, `PerSiteRoundRobin`. Defaults to `LeastRequests` if omitted.",
				},

				"managed_pipeline_mode": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(webapps.ManagedPipelineModeIntegrated),
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForManagedPipelineMode(), false),
					Description:  "The Managed Pipeline mode. Possible values include: `Integrated`, `Classic`. Defaults to `Integrated`.",
				},

				"pre_warmed_instance_count": {
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Computed:    true, // Variable defaults depending on plan etc
					Description: "The number of pre-warmed instances for this function app. Only affects apps on an Elastic Premium plan.",
				},

				"remote_debugging_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should Remote Debugging be enabled. Defaults to `false`.",
				},

				"remote_debugging_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"VS2017",
						"VS2019",
						"VS2022",
					}, false),
					Description: "The Remote Debugging Version. Possible values include `VS2017`, `VS2019`, and `VS2022`",
				},

				"runtime_scale_monitoring_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Description: "Should Functions Runtime Scale Monitoring be enabled.",
				},

				"scm_type": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The SCM Type in use by the Windows Function App.",
				},

				"use_32_bit_worker": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Should the Windows Web App use a 32-bit worker.",
				},

				"websockets_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should Web Sockets be enabled. Defaults to `false`.",
				},

				"ftps_state": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(webapps.FtpsStateDisabled),
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForFtpsState(), false),
					Description:  "State of FTP / FTPS service for this function app. Possible values include: `AllAllowed`, `FtpsOnly` and `Disabled`. Defaults to `Disabled`.",
				},

				"health_check_path": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The path to be checked for this function app health.",
					RequiredWith: func() []string {
						if features.FourPointOhBeta() {
							return []string{"site_config.0.health_check_eviction_time_in_min"}
						}
						return []string{}
					}(),
				},

				"health_check_eviction_time_in_min": { // NOTE: Will evict the only node in single node configurations.
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     !features.FourPointOhBeta(),
					ValidateFunc: validation.IntBetween(2, 10),
					RequiredWith: func() []string {
						if features.FourPointOhBeta() {
							return []string{"site_config.0.health_check_path"}
						}
						return []string{}
					}(),
					Description: "The amount of time in minutes that a node is unhealthy before being removed from the load balancer. Possible values are between `2` and `10`. Only valid in conjunction with `health_check_path`",
				},

				"worker_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 100),
					Description:  "The number of Workers for this Windows Function App.",
				},

				"minimum_tls_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(webapps.SupportedTlsVersionsOnePointTwo),
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForSupportedTlsVersions(), false),
					Description:  "The configures the minimum version of TLS required for SSL requests. Possible values include: `1.0`, `1.1`, and  `1.2`. Defaults to `1.2`.",
				},

				"scm_minimum_tls_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(webapps.SupportedTlsVersionsOnePointTwo),
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForSupportedTlsVersions(), false),
					Description:  "Configures the minimum version of TLS required for SSL requests to the SCM site Possible values include: `1.0`, `1.1`, and  `1.2`. Defaults to `1.2`.",
				},

				"cors": CorsSettingsSchema(),

				"vnet_route_all_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should all outbound traffic to have Virtual Network Security Groups and User Defined Routes applied? Defaults to `false`.",
				},

				"detailed_error_logging_enabled": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Is detailed error logging enabled",
				},

				"windows_fx_version": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The Windows FX Version string.",
				},
			},
		},
	}
}

func SiteConfigSchemaWindowsFunctionAppComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"api_management_api_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"api_definition_url": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"app_command_line": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"app_scale_limit": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"application_insights_key": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"application_insights_connection_string": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},

				"application_stack": windowsFunctionAppStackSchemaComputed(),

				"app_service_logs": FunctionAppAppServiceLogsSchemaComputed(),

				"default_documents": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"elastic_instance_minimum": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"http2_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"ip_restriction": IpRestrictionSchemaComputed(),

				"ip_restriction_default_action": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"scm_ip_restriction": IpRestrictionSchemaComputed(),

				"scm_ip_restriction_default_action": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"load_balancing_mode": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"managed_pipeline_mode": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"pre_warmed_instance_count": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"remote_debugging_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"remote_debugging_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"runtime_scale_monitoring_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"scm_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"use_32_bit_worker": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"websockets_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"ftps_state": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"health_check_path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"health_check_eviction_time_in_min": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"worker_count": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"scm_minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"cors": CorsSettingsSchemaComputed(),

				"vnet_route_all_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"detailed_error_logging_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"windows_fx_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

type ApplicationStackLinuxFunctionApp struct {
	// Note - Function Apps differ to Web Apps here. They do not use the named properties in the SiteConfig block and exclusively use the app_settings map
	DotNetVersion         string                   `tfschema:"dotnet_version"`              // Supported values `3.1`, `6.0`, `7.0` and `8.0`.
	DotNetIsolated        bool                     `tfschema:"use_dotnet_isolated_runtime"` // Supported values `true` for `dotnet-isolated`, `false` otherwise
	NodeVersion           string                   `tfschema:"node_version"`                // Supported values `12LTS`, `14LTS`, `16LTS`, `18LTS, `20LTS``
	PythonVersion         string                   `tfschema:"python_version"`              // Supported values `3.12`, `3.11`, `3.10`, `3.9`, `3.8`, `3.7`
	PowerShellCoreVersion string                   `tfschema:"powershell_core_version"`     // Supported values are `7.0`, `7.2`
	JavaVersion           string                   `tfschema:"java_version"`                // Supported values `8`, `11`, `17`
	CustomHandler         bool                     `tfschema:"use_custom_runtime"`          // Supported values `true`
	Docker                []ApplicationStackDocker `tfschema:"docker"`                      // Needs ElasticPremium or Basic (B1) Standard (S 1-3) or Premium(PxV2 or PxV3) LINUX Service Plan
}

type ApplicationStackWindowsFunctionApp struct {
	DotNetVersion         string `tfschema:"dotnet_version"`              // Supported values `v3.0`, `v4.0`, `v6.0`, `v7.0` and `v8.0`
	DotNetIsolated        bool   `tfschema:"use_dotnet_isolated_runtime"` // Supported values `true` for `dotnet-isolated`, `false` otherwise
	NodeVersion           string `tfschema:"node_version"`                // Supported values `12LTS`, `14LTS`, `16LTS`, `18LTS, `20LTS`
	JavaVersion           string `tfschema:"java_version"`                // Supported values `8`, `11`, `17`
	PowerShellCoreVersion string `tfschema:"powershell_core_version"`     // Supported values are `7.0`, `7.2`
	CustomHandler         bool   `tfschema:"use_custom_runtime"`          // Supported values `true`
}

type ApplicationStackDocker struct {
	RegistryURL      string `tfschema:"registry_url"`
	RegistryUsername string `tfschema:"registry_username"`
	RegistryPassword string `tfschema:"registry_password"`
	ImageName        string `tfschema:"image_name"`
	ImageTag         string `tfschema:"image_tag"`
}

func linuxFunctionAppStackSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"dotnet_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"3.1",
						"6.0",
						"7.0",
						"8.0",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.powershell_core_version",
						"site_config.0.application_stack.0.docker",
						"site_config.0.application_stack.0.use_custom_runtime",
					},
					Description: "The version of .Net. Possible values are `3.1`, `6.0` and `7.0`",
				},

				"use_dotnet_isolated_runtime": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
					ConflictsWith: []string{
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.powershell_core_version",
						"site_config.0.application_stack.0.docker",
						"site_config.0.application_stack.0.use_custom_runtime",
					},
					Description: "Should the DotNet process use an isolated runtime. Defaults to `false`.",
				},

				"python_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"3.12",
						"3.11",
						"3.10",
						"3.9",
						"3.8",
						"3.7",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.powershell_core_version",
						"site_config.0.application_stack.0.docker",
						"site_config.0.application_stack.0.use_custom_runtime",
					},
					Description: "The version of Python to use. Possible values include `3.12`, `3.11`, `3.10`, `3.9`, `3.8`, and `3.7`.",
				},

				"node_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"12", // Deprecated, and removed from portal, but seemingly accepted by API
						"14",
						"16",
						"18",
						"20",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.powershell_core_version",
						"site_config.0.application_stack.0.docker",
						"site_config.0.application_stack.0.use_custom_runtime",
					},
					Description: "The version of Node to use. Possible values include `12`, `14`, `16`, `18` and `20`",
				},

				"powershell_core_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"7",   // Deprecated / not available in the portal
						"7.2", // preview LTS Support
						"7.4", // current LTS Support
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.powershell_core_version",
						"site_config.0.application_stack.0.docker",
						"site_config.0.application_stack.0.use_custom_runtime",
					},
					Description: "The version of PowerShell Core to use. Possibles values are `7`, `7.2`, and `7.4`",
				},

				"java_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"8",
						"11",
						"17",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.powershell_core_version",
						"site_config.0.application_stack.0.docker",
						"site_config.0.application_stack.0.use_custom_runtime",
					},
					Description: "The version of Java to use. Possible values are `8`, `11`, and `17`",
				},

				"docker": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*schema.Schema{
							"registry_url": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
								Description:  "The URL of the docker registry.",
							},

							"registry_username": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Sensitive:    true,
								ValidateFunc: validation.StringIsNotEmpty,
								Description:  "The username to use for connections to the registry.",
							},

							"registry_password": {
								Type:        pluginsdk.TypeString,
								Optional:    true,
								Sensitive:   true, // Note: whilst it's not a good idea, this _can_ be blank...
								Description: "The password for the account to use to connect to the registry.",
							},

							"image_name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
								Description:  "The name of the Docker image to use.",
							},

							"image_tag": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
								Description:  "The image tag of the image to use.",
							},
						},
					},
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.powershell_core_version",
						"site_config.0.application_stack.0.docker",
						"site_config.0.application_stack.0.use_custom_runtime",
					},
					Description: "A docker block",
				},

				"use_custom_runtime": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.powershell_core_version",
						"site_config.0.application_stack.0.docker",
						"site_config.0.application_stack.0.use_custom_runtime",
					},
				},
			},
		},
	}
}

func linuxFunctionAppStackSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"dotnet_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"use_dotnet_isolated_runtime": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"python_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"node_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"powershell_core_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*schema.Schema{
							"registry_url": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"registry_username": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"registry_password": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"image_name": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"image_tag": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},

				"use_custom_runtime": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

func windowsFunctionAppStackSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"dotnet_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "v4.0",
					ValidateFunc: validation.StringInSlice([]string{
						"v3.0",
						"v4.0",
						"v6.0",
						"v7.0",
						"v8.0",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.powershell_core_version",
						"site_config.0.application_stack.0.use_custom_runtime",
					},
					Description: "The version of .Net. Possible values are `v3.0`, `v4.0`, `v6.0` and `v7.0`",
				},

				"use_dotnet_isolated_runtime": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
					ConflictsWith: []string{
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.powershell_core_version",
						"site_config.0.application_stack.0.use_custom_runtime",
					},
					Description: "Should the DotNet process use an isolated runtime. Defaults to `false`.",
				},

				"node_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"~12",
						"~14",
						"~16",
						"~18",
						"~20",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.powershell_core_version",
						"site_config.0.application_stack.0.use_custom_runtime",
					},
					Description: "The version of Node to use. Possible values include `12`, `14`, `16` and `18`",
				},

				"java_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"1.8",
						"11",
						"17",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.powershell_core_version",
						"site_config.0.application_stack.0.use_custom_runtime",
					},
					Description: "The version of Java to use. Possible values are `1.8`, `11` and `17`",
				},

				"powershell_core_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"7",   // Deprecated / not available in the portal
						"7.2", // preview LTS Support
						"7.4", // current LTS Support
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.powershell_core_version",
						"site_config.0.application_stack.0.use_custom_runtime",
					},
					Description: "The PowerShell Core version to use. Possible values are `7`, `7.2`, and `7.4`",
				},

				"use_custom_runtime": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.powershell_core_version",
						"site_config.0.application_stack.0.use_custom_runtime",
					},
					Description: "Does the Function App use a custom Application Stack?",
				},
			},
		},
	}
}

func windowsFunctionAppStackSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"dotnet_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"use_dotnet_isolated_runtime": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"node_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"powershell_core_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"use_custom_runtime": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

type FunctionAppAppServiceLogs struct {
	DiskQuotaMB         int64 `tfschema:"disk_quota_mb"`
	RetentionPeriodDays int64 `tfschema:"retention_period_days"`
}

func FunctionAppAppServiceLogsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"disk_quota_mb": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      35,
					ValidateFunc: validation.IntBetween(25, 100),
					Description:  "The amount of disk space to use for logs. Valid values are between `25` and `100`.",
				},
				"retention_period_days": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 99999),
					Description:  "The retention period for logs in days. Valid values are between `0` and `99999`. Defaults to `0` (never delete).",
				},
			},
		},
	}
}

func FunctionAppAppServiceLogsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"disk_quota_mb": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
				"retention_period_days": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func ExpandSiteConfigLinuxFunctionApp(siteConfig []SiteConfigLinuxFunctionApp, existing *webapps.SiteConfig, metadata sdk.ResourceMetaData, version string, storageString string, storageUsesMSI bool) (*webapps.SiteConfig, error) {
	if len(siteConfig) == 0 {
		return nil, nil
	}

	expanded := &webapps.SiteConfig{}
	if existing != nil {
		expanded = existing
		// need to zero fxversion to re-calculate based on changes below or removing app_stack doesn't apply
		expanded.LinuxFxVersion = pointer.To("")
	}

	appSettings := make([]webapps.NameValuePair, 0)

	if existing != nil && existing.AppSettings != nil {
		appSettings = *existing.AppSettings
	}

	appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_EXTENSION_VERSION", version, false)

	if storageUsesMSI {
		appSettings = updateOrAppendAppSettings(appSettings, "AzureWebJobsStorage__accountName", storageString, false)
	} else {
		appSettings = updateOrAppendAppSettings(appSettings, "AzureWebJobsStorage", storageString, false)
	}

	linuxSiteConfig := siteConfig[0]

	v := strconv.FormatInt(linuxSiteConfig.HealthCheckEvictionTime, 10)
	if v == "0" || linuxSiteConfig.HealthCheckPath == "" {
		appSettings = updateOrAppendAppSettings(appSettings, "WEBSITE_HEALTHCHECK_MAXPINGFAILURES", v, true)
	} else {
		appSettings = updateOrAppendAppSettings(appSettings, "WEBSITE_HEALTHCHECK_MAXPINGFAILURES", v, false)
	}

	expanded.AlwaysOn = pointer.To(linuxSiteConfig.AlwaysOn)

	if metadata.ResourceData.HasChange("site_config.0.app_scale_limit") {
		expanded.FunctionAppScaleLimit = pointer.To(linuxSiteConfig.AppScaleLimit)
	}

	if linuxSiteConfig.AppInsightsConnectionString == "" {
		appSettings = updateOrAppendAppSettings(appSettings, "APPLICATIONINSIGHTS_CONNECTION_STRING", linuxSiteConfig.AppInsightsConnectionString, true)
	} else {
		appSettings = updateOrAppendAppSettings(appSettings, "APPLICATIONINSIGHTS_CONNECTION_STRING", linuxSiteConfig.AppInsightsConnectionString, false)
	}

	if linuxSiteConfig.AppInsightsInstrumentationKey == "" {
		appSettings = updateOrAppendAppSettings(appSettings, "APPINSIGHTS_INSTRUMENTATIONKEY", linuxSiteConfig.AppInsightsInstrumentationKey, true)
	} else {
		appSettings = updateOrAppendAppSettings(appSettings, "APPINSIGHTS_INSTRUMENTATIONKEY", linuxSiteConfig.AppInsightsInstrumentationKey, false)
	}

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.ApiManagementConfig = &webapps.ApiManagementConfig{
			Id: pointer.To(linuxSiteConfig.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.ApiDefinition = &webapps.ApiDefinitionInfo{
			Url: pointer.To(linuxSiteConfig.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = pointer.To(linuxSiteConfig.AppCommandLine)
	}

	if len(linuxSiteConfig.ApplicationStack) > 0 {
		linuxAppStack := linuxSiteConfig.ApplicationStack[0]
		if linuxAppStack.DotNetVersion != "" {
			if linuxAppStack.DotNetIsolated {
				appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "dotnet-isolated", false)
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("DOTNET-ISOLATED|%s", linuxAppStack.DotNetVersion))
			} else {
				appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "dotnet", false)
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("DOTNET|%s", linuxAppStack.DotNetVersion))
			}
		}

		if linuxAppStack.NodeVersion != "" {
			appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "node", false)
			appSettings = updateOrAppendAppSettings(appSettings, "WEBSITE_NODE_DEFAULT_VERSION", linuxAppStack.NodeVersion, false)
			expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("NODE|%s", linuxAppStack.NodeVersion))
		}

		if linuxAppStack.PythonVersion != "" {
			appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "python", false)
			expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("PYTHON|%s", linuxAppStack.PythonVersion))
		}

		if linuxAppStack.JavaVersion != "" {
			appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "java", false)
			expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("JAVA|%s", linuxAppStack.JavaVersion))
		}

		if linuxAppStack.PowerShellCoreVersion != "" {
			appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "powershell", false)
			expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("POWERSHELL|%s", linuxAppStack.PowerShellCoreVersion))
		}

		if linuxAppStack.CustomHandler {
			appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "custom", false)
			expanded.LinuxFxVersion = pointer.To("") // Custom needs an explicit empty string here
		}

		if linuxAppStack.Docker != nil && len(linuxAppStack.Docker) == 1 {
			dockerConfig := linuxAppStack.Docker[0]
			appSettings = updateOrAppendAppSettings(appSettings, "DOCKER_REGISTRY_SERVER_URL", dockerConfig.RegistryURL, false)
			appSettings = updateOrAppendAppSettings(appSettings, "DOCKER_REGISTRY_SERVER_USERNAME", dockerConfig.RegistryUsername, false)
			appSettings = updateOrAppendAppSettings(appSettings, "DOCKER_REGISTRY_SERVER_PASSWORD", dockerConfig.RegistryPassword, false)
			var dockerUrl string = dockerConfig.RegistryURL
			for _, prefix := range urlSchemes {
				if strings.HasPrefix(dockerConfig.RegistryURL, prefix) {
					dockerUrl = strings.TrimPrefix(dockerConfig.RegistryURL, prefix)
					continue
				}
			}
			expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("DOCKER|%s/%s:%s", dockerUrl, dockerConfig.ImageName, dockerConfig.ImageTag))
		}
	} else {
		appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "", true)
		expanded.LinuxFxVersion = pointer.To("")
	}

	if metadata.ResourceData.HasChange("site_config.0.container_registry_use_managed_identity") {
		expanded.AcrUseManagedIdentityCreds = pointer.To(linuxSiteConfig.UseManagedIdentityACR)
	}

	if metadata.ResourceData.HasChange("site_config.0.vnet_route_all_enabled") {
		expanded.VnetRouteAllEnabled = pointer.To(linuxSiteConfig.VnetRouteAllEnabled)
	}

	if metadata.ResourceData.HasChange("site_config.0.container_registry_managed_identity_client_id") {
		expanded.AcrUserManagedIdentityID = pointer.To(linuxSiteConfig.ContainerRegistryMSI)
	}

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = &linuxSiteConfig.DefaultDocuments
	}

	if metadata.ResourceData.HasChange("site_config.0.http2_enabled") {
		expanded.HTTP20Enabled = pointer.To(linuxSiteConfig.Http2Enabled)
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := ExpandIpRestrictions(linuxSiteConfig.IpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction_default_action") {
		expanded.IPSecurityRestrictionsDefaultAction = pointer.To(webapps.DefaultAction(linuxSiteConfig.IpRestrictionDefaultAction))
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_use_main_ip_restriction") {
		expanded.ScmIPSecurityRestrictionsUseMain = pointer.To(linuxSiteConfig.ScmUseMainIpRestriction)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := ExpandIpRestrictions(linuxSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction_default_action") {
		expanded.ScmIPSecurityRestrictionsDefaultAction = pointer.To(webapps.DefaultAction(linuxSiteConfig.ScmIpRestrictionDefaultAction))
	}

	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = pointer.To(webapps.SiteLoadBalancing(linuxSiteConfig.LoadBalancing))
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = pointer.To(webapps.ManagedPipelineMode(linuxSiteConfig.ManagedPipelineMode))
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_enabled") {
		expanded.RemoteDebuggingEnabled = pointer.To(linuxSiteConfig.RemoteDebugging)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = pointer.To(linuxSiteConfig.RemoteDebuggingVersion)
	}

	expanded.Use32BitWorkerProcess = pointer.To(linuxSiteConfig.Use32BitWorker)

	if metadata.ResourceData.HasChange("site_config.0.websockets_enabled") {
		expanded.WebSocketsEnabled = pointer.To(linuxSiteConfig.WebSockets)
	}

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		expanded.FtpsState = pointer.To(webapps.FtpsState(linuxSiteConfig.FtpsState))
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = pointer.To(linuxSiteConfig.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.worker_count") {
		expanded.NumberOfWorkers = pointer.To(linuxSiteConfig.WorkerCount)
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTlsVersion = pointer.To(webapps.SupportedTlsVersions(linuxSiteConfig.MinTlsVersion))
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTlsVersion = pointer.To(webapps.SupportedTlsVersions(linuxSiteConfig.ScmMinTlsVersion))
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		cors := ExpandCorsSettings(linuxSiteConfig.Cors)
		expanded.Cors = cors
	}

	if metadata.ResourceData.HasChange("site_config.0.pre_warmed_instance_count") {
		expanded.PreWarmedInstanceCount = pointer.To(linuxSiteConfig.PreWarmedInstanceCount)
	}

	if metadata.ResourceData.HasChange("site_config.0.vnet_route_all_enabled") {
		expanded.VnetRouteAllEnabled = pointer.To(linuxSiteConfig.VnetRouteAllEnabled)
	}

	if metadata.ResourceData.HasChange("site_config.0.elastic_instance_minimum") {
		expanded.MinimumElasticInstanceCount = pointer.To(linuxSiteConfig.ElasticInstanceMinimum)
	}

	if metadata.ResourceData.HasChange("site_config.0.runtime_scale_monitoring_enabled") {
		expanded.FunctionsRuntimeScaleMonitoringEnabled = pointer.To(linuxSiteConfig.RuntimeScaleMonitoring)
	}

	expanded.AppSettings = &appSettings

	return expanded, nil
}

// updateOrAppendAppSettings is used to modify a collection of webapps.NameValuePair items.
func updateOrAppendAppSettings(input []webapps.NameValuePair, name string, value string, remove bool) []webapps.NameValuePair {
	for k, v := range input {
		if v.Name != nil && *v.Name == name {
			if remove {
				input[k] = input[len(input)-1]
				input[len(input)-1] = webapps.NameValuePair{}
				input = input[:len(input)-1]
			} else {
				input[k] = webapps.NameValuePair{
					Name:  pointer.To(name),
					Value: pointer.To(value),
				}
			}
			return input
		}
	}

	if !remove {
		input = append(input, webapps.NameValuePair{
			Name:  pointer.To(name),
			Value: pointer.To(value),
		})
	}

	return input
}

func ExpandSiteConfigWindowsFunctionApp(siteConfig []SiteConfigWindowsFunctionApp, existing *webapps.SiteConfig, metadata sdk.ResourceMetaData, version string, storageString string, storageUsesMSI bool) (*webapps.SiteConfig, error) {
	if len(siteConfig) == 0 {
		return nil, nil
	}

	expanded := &webapps.SiteConfig{}
	if existing != nil {
		expanded = existing
		// need to zero fxversion to re-calculate based on changes below or removing app_stack doesn't apply
		expanded.WindowsFxVersion = pointer.To("")
	}

	appSettings := make([]webapps.NameValuePair, 0)

	if existing != nil && existing.AppSettings != nil {
		appSettings = *existing.AppSettings
	}

	appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_EXTENSION_VERSION", version, false)

	if storageUsesMSI {
		appSettings = updateOrAppendAppSettings(appSettings, "AzureWebJobsStorage__accountName", storageString, false)
	} else {
		appSettings = updateOrAppendAppSettings(appSettings, "AzureWebJobsStorage", storageString, false)
	}

	windowsSiteConfig := siteConfig[0]

	v := strconv.FormatInt(windowsSiteConfig.HealthCheckEvictionTime, 10)
	if v == "0" || windowsSiteConfig.HealthCheckPath == "" {
		appSettings = updateOrAppendAppSettings(appSettings, "WEBSITE_HEALTHCHECK_MAXPINGFAILURES", v, true)
	} else {
		appSettings = updateOrAppendAppSettings(appSettings, "WEBSITE_HEALTHCHECK_MAXPINGFAILURES", v, false)
	}

	expanded.AlwaysOn = pointer.To(windowsSiteConfig.AlwaysOn)

	if metadata.ResourceData.HasChange("site_config.0.app_scale_limit") {
		expanded.FunctionAppScaleLimit = pointer.To(windowsSiteConfig.AppScaleLimit)
	}

	if windowsSiteConfig.AppInsightsConnectionString == "" {
		appSettings = updateOrAppendAppSettings(appSettings, "APPLICATIONINSIGHTS_CONNECTION_STRING", windowsSiteConfig.AppInsightsConnectionString, true)
	} else {
		appSettings = updateOrAppendAppSettings(appSettings, "APPLICATIONINSIGHTS_CONNECTION_STRING", windowsSiteConfig.AppInsightsConnectionString, false)
	}

	if windowsSiteConfig.AppInsightsInstrumentationKey == "" {
		appSettings = updateOrAppendAppSettings(appSettings, "APPINSIGHTS_INSTRUMENTATIONKEY", windowsSiteConfig.AppInsightsInstrumentationKey, true)
	} else {
		appSettings = updateOrAppendAppSettings(appSettings, "APPINSIGHTS_INSTRUMENTATIONKEY", windowsSiteConfig.AppInsightsInstrumentationKey, false)
	}

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.ApiManagementConfig = &webapps.ApiManagementConfig{
			Id: pointer.To(windowsSiteConfig.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.ApiDefinition = &webapps.ApiDefinitionInfo{
			Url: pointer.To(windowsSiteConfig.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = pointer.To(windowsSiteConfig.AppCommandLine)
	}

	if len(windowsSiteConfig.ApplicationStack) > 0 {
		windowsAppStack := windowsSiteConfig.ApplicationStack[0]
		if windowsAppStack.DotNetVersion != "" {
			if windowsAppStack.DotNetIsolated {
				appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "dotnet-isolated", false)
			} else {
				appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "dotnet", false)
			}
			expanded.NetFrameworkVersion = pointer.To(windowsAppStack.DotNetVersion)
		}

		if windowsAppStack.NodeVersion != "" {
			appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "node", false)
			appSettings = updateOrAppendAppSettings(appSettings, "WEBSITE_NODE_DEFAULT_VERSION", windowsAppStack.NodeVersion, false)
		}

		if windowsAppStack.JavaVersion != "" {
			appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "java", false)
			expanded.JavaVersion = pointer.To(windowsAppStack.JavaVersion)
		}

		if windowsAppStack.PowerShellCoreVersion != "" {
			appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "powershell", false)
			expanded.PowerShellVersion = pointer.To(strings.TrimPrefix(windowsAppStack.PowerShellCoreVersion, "~"))
		}

		if windowsAppStack.CustomHandler {
			appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "custom", false)
			expanded.WindowsFxVersion = pointer.To("") // Custom needs an explicit empty string here
		}
	} else {
		appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "", true)
		expanded.WindowsFxVersion = pointer.To("")
	}

	if metadata.ResourceData.HasChange("site_config.0.vnet_route_all_enabled") {
		expanded.VnetRouteAllEnabled = pointer.To(windowsSiteConfig.VnetRouteAllEnabled)
	}

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = &windowsSiteConfig.DefaultDocuments
	}

	if metadata.ResourceData.HasChange("site_config.0.http2_enabled") {
		expanded.HTTP20Enabled = pointer.To(windowsSiteConfig.Http2Enabled)
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := ExpandIpRestrictions(windowsSiteConfig.IpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction_default_action") {
		expanded.IPSecurityRestrictionsDefaultAction = pointer.To(webapps.DefaultAction(windowsSiteConfig.IpRestrictionDefaultAction))
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_use_main_ip_restriction") {
		expanded.ScmIPSecurityRestrictionsUseMain = pointer.To(windowsSiteConfig.ScmUseMainIpRestriction)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := ExpandIpRestrictions(windowsSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction_default_action") {
		expanded.ScmIPSecurityRestrictionsDefaultAction = pointer.To(webapps.DefaultAction(windowsSiteConfig.ScmIpRestrictionDefaultAction))
	}

	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = pointer.To(webapps.SiteLoadBalancing(windowsSiteConfig.LoadBalancing))
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = pointer.To(webapps.ManagedPipelineMode(windowsSiteConfig.ManagedPipelineMode))
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_enabled") {
		expanded.RemoteDebuggingEnabled = pointer.To(windowsSiteConfig.RemoteDebugging)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = pointer.To(windowsSiteConfig.RemoteDebuggingVersion)
	}

	expanded.Use32BitWorkerProcess = pointer.To(windowsSiteConfig.Use32BitWorker)

	if metadata.ResourceData.HasChange("site_config.0.websockets_enabled") {
		expanded.WebSocketsEnabled = pointer.To(windowsSiteConfig.WebSockets)
	}

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		expanded.FtpsState = pointer.To(webapps.FtpsState(windowsSiteConfig.FtpsState))
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = pointer.To(windowsSiteConfig.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.worker_count") {
		expanded.NumberOfWorkers = pointer.To(windowsSiteConfig.NumberOfWorkers)
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTlsVersion = pointer.To(webapps.SupportedTlsVersions(windowsSiteConfig.MinTlsVersion))
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTlsVersion = pointer.To(webapps.SupportedTlsVersions(windowsSiteConfig.ScmMinTlsVersion))
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		cors := ExpandCorsSettings(windowsSiteConfig.Cors)
		expanded.Cors = cors
	}

	if metadata.ResourceData.HasChange("site_config.0.pre_warmed_instance_count") {
		expanded.PreWarmedInstanceCount = pointer.To(windowsSiteConfig.PreWarmedInstanceCount)
	}

	if metadata.ResourceData.HasChange("site_config.0.vnet_route_all_enabled") {
		expanded.VnetRouteAllEnabled = pointer.To(windowsSiteConfig.VnetRouteAllEnabled)
	}

	if metadata.ResourceData.HasChange("site_config.0.elastic_instance_minimum") {
		expanded.MinimumElasticInstanceCount = pointer.To(windowsSiteConfig.ElasticInstanceMinimum)
	}

	if metadata.ResourceData.HasChange("site_config.0.runtime_scale_monitoring_enabled") {
		expanded.FunctionsRuntimeScaleMonitoringEnabled = pointer.To(windowsSiteConfig.RuntimeScaleMonitoring)
	}

	expanded.AppSettings = &appSettings

	return expanded, nil
}

func FlattenSiteConfigLinuxFunctionApp(functionAppSiteConfig *webapps.SiteConfig) (*SiteConfigLinuxFunctionApp, error) {
	if functionAppSiteConfig == nil {
		return nil, fmt.Errorf("flattening site config: SiteConfig was nil")
	}

	result := &SiteConfigLinuxFunctionApp{
		AlwaysOn:                      pointer.From(functionAppSiteConfig.AlwaysOn),
		AppCommandLine:                pointer.From(functionAppSiteConfig.AppCommandLine),
		AppScaleLimit:                 pointer.From(functionAppSiteConfig.FunctionAppScaleLimit),
		ContainerRegistryMSI:          pointer.From(functionAppSiteConfig.AcrUserManagedIdentityID),
		Cors:                          FlattenCorsSettings(functionAppSiteConfig.Cors),
		DetailedErrorLogging:          pointer.From(functionAppSiteConfig.DetailedErrorLoggingEnabled),
		HealthCheckPath:               pointer.From(functionAppSiteConfig.HealthCheckPath),
		Http2Enabled:                  pointer.From(functionAppSiteConfig.HTTP20Enabled),
		IpRestrictionDefaultAction:    string(pointer.From(functionAppSiteConfig.IPSecurityRestrictionsDefaultAction)),
		ScmIpRestrictionDefaultAction: string(pointer.From(functionAppSiteConfig.ScmIPSecurityRestrictionsDefaultAction)),
		LinuxFxVersion:                pointer.From(functionAppSiteConfig.LinuxFxVersion),
		LoadBalancing:                 string(pointer.From(functionAppSiteConfig.LoadBalancing)),
		ManagedPipelineMode:           string(pointer.From(functionAppSiteConfig.ManagedPipelineMode)),
		WorkerCount:                   pointer.From(functionAppSiteConfig.NumberOfWorkers),
		ScmType:                       string(pointer.From(functionAppSiteConfig.ScmType)),
		FtpsState:                     string(pointer.From(functionAppSiteConfig.FtpsState)),
		RuntimeScaleMonitoring:        pointer.From(functionAppSiteConfig.FunctionsRuntimeScaleMonitoringEnabled),
		MinTlsVersion:                 string(pointer.From(functionAppSiteConfig.MinTlsVersion)),
		ScmMinTlsVersion:              string(pointer.From(functionAppSiteConfig.ScmMinTlsVersion)),
		PreWarmedInstanceCount:        pointer.From(functionAppSiteConfig.PreWarmedInstanceCount),
		ElasticInstanceMinimum:        pointer.From(functionAppSiteConfig.MinimumElasticInstanceCount),
		Use32BitWorker:                pointer.From(functionAppSiteConfig.Use32BitWorkerProcess),
		WebSockets:                    pointer.From(functionAppSiteConfig.WebSocketsEnabled),
		ScmUseMainIpRestriction:       pointer.From(functionAppSiteConfig.ScmIPSecurityRestrictionsUseMain),
		UseManagedIdentityACR:         pointer.From(functionAppSiteConfig.AcrUseManagedIdentityCreds),
		RemoteDebugging:               pointer.From(functionAppSiteConfig.RemoteDebuggingEnabled),
		RemoteDebuggingVersion:        strings.ToUpper(pointer.From(functionAppSiteConfig.RemoteDebuggingVersion)),
		VnetRouteAllEnabled:           pointer.From(functionAppSiteConfig.VnetRouteAllEnabled),
	}

	if v := functionAppSiteConfig.ApiDefinition; v != nil && v.Url != nil {
		result.ApiDefinition = *v.Url
	}

	if v := functionAppSiteConfig.ApiManagementConfig; v != nil && v.Id != nil {
		result.ApiManagementConfigId = *v.Id
	}

	if functionAppSiteConfig.IPSecurityRestrictions != nil {
		result.IpRestriction = FlattenIpRestrictions(functionAppSiteConfig.IPSecurityRestrictions)
	}

	if functionAppSiteConfig.ScmIPSecurityRestrictions != nil {
		result.ScmIpRestriction = FlattenIpRestrictions(functionAppSiteConfig.ScmIPSecurityRestrictions)
	}

	if v := functionAppSiteConfig.DefaultDocuments; v != nil {
		result.DefaultDocuments = *v
	}

	var appStack []ApplicationStackLinuxFunctionApp
	if functionAppSiteConfig.LinuxFxVersion != nil {
		decoded, err := DecodeFunctionAppLinuxFxVersion(*functionAppSiteConfig.LinuxFxVersion)
		if err != nil {
			return nil, fmt.Errorf("flattening site config: %s", err)
		}
		appStack = decoded
	}
	result.ApplicationStack = appStack

	return result, nil
}

func FlattenSiteConfigWindowsFunctionApp(functionAppSiteConfig *webapps.SiteConfig) (*SiteConfigWindowsFunctionApp, error) {
	if functionAppSiteConfig == nil {
		return nil, fmt.Errorf("flattening site config: SiteConfig was nil")
	}

	result := &SiteConfigWindowsFunctionApp{
		AlwaysOn:                      pointer.From(functionAppSiteConfig.AlwaysOn),
		AppCommandLine:                pointer.From(functionAppSiteConfig.AppCommandLine),
		AppScaleLimit:                 pointer.From(functionAppSiteConfig.FunctionAppScaleLimit),
		Cors:                          FlattenCorsSettings(functionAppSiteConfig.Cors),
		DetailedErrorLogging:          pointer.From(functionAppSiteConfig.DetailedErrorLoggingEnabled),
		HealthCheckPath:               pointer.From(functionAppSiteConfig.HealthCheckPath),
		Http2Enabled:                  pointer.From(functionAppSiteConfig.HTTP20Enabled),
		WindowsFxVersion:              pointer.From(functionAppSiteConfig.WindowsFxVersion),
		LoadBalancing:                 string(pointer.From(functionAppSiteConfig.LoadBalancing)),
		ManagedPipelineMode:           string(pointer.From(functionAppSiteConfig.ManagedPipelineMode)),
		NumberOfWorkers:               pointer.From(functionAppSiteConfig.NumberOfWorkers),
		ScmType:                       string(pointer.From(functionAppSiteConfig.ScmType)),
		FtpsState:                     string(pointer.From(functionAppSiteConfig.FtpsState)),
		RuntimeScaleMonitoring:        pointer.From(functionAppSiteConfig.FunctionsRuntimeScaleMonitoringEnabled),
		MinTlsVersion:                 string(pointer.From(functionAppSiteConfig.MinTlsVersion)),
		ScmMinTlsVersion:              string(pointer.From(functionAppSiteConfig.ScmMinTlsVersion)),
		PreWarmedInstanceCount:        pointer.From(functionAppSiteConfig.PreWarmedInstanceCount),
		ElasticInstanceMinimum:        pointer.From(functionAppSiteConfig.MinimumElasticInstanceCount),
		Use32BitWorker:                pointer.From(functionAppSiteConfig.Use32BitWorkerProcess),
		WebSockets:                    pointer.From(functionAppSiteConfig.WebSocketsEnabled),
		ScmUseMainIpRestriction:       pointer.From(functionAppSiteConfig.ScmIPSecurityRestrictionsUseMain),
		RemoteDebugging:               pointer.From(functionAppSiteConfig.RemoteDebuggingEnabled),
		RemoteDebuggingVersion:        strings.ToUpper(pointer.From(functionAppSiteConfig.RemoteDebuggingVersion)),
		VnetRouteAllEnabled:           pointer.From(functionAppSiteConfig.VnetRouteAllEnabled),
		IpRestrictionDefaultAction:    string(pointer.From(functionAppSiteConfig.IPSecurityRestrictionsDefaultAction)),
		ScmIpRestrictionDefaultAction: string(pointer.From(functionAppSiteConfig.ScmIPSecurityRestrictionsDefaultAction)),
	}

	if v := functionAppSiteConfig.ApiDefinition; v != nil && v.Url != nil {
		result.ApiDefinition = *v.Url
	}

	if v := functionAppSiteConfig.ApiManagementConfig; v != nil && v.Id != nil {
		result.ApiManagementConfigId = *v.Id
	}

	if functionAppSiteConfig.IPSecurityRestrictions != nil {
		result.IpRestriction = FlattenIpRestrictions(functionAppSiteConfig.IPSecurityRestrictions)
	}

	if functionAppSiteConfig.ScmIPSecurityRestrictions != nil {
		result.ScmIpRestriction = FlattenIpRestrictions(functionAppSiteConfig.ScmIPSecurityRestrictions)
	}

	if v := functionAppSiteConfig.DefaultDocuments; v != nil {
		result.DefaultDocuments = *v
	}

	powershellVersion := ""
	if p := functionAppSiteConfig.PowerShellVersion; p != nil {
		powershellVersion = *p
		if powershellVersion == "~7" {
			powershellVersion = "7"
		}
	}

	result.ApplicationStack = []ApplicationStackWindowsFunctionApp{{
		DotNetVersion:         pointer.From(functionAppSiteConfig.NetFrameworkVersion),
		DotNetIsolated:        false, // set this later from app_settings
		NodeVersion:           "",    // Need to get this from app_settings later
		JavaVersion:           pointer.From(functionAppSiteConfig.JavaVersion),
		PowerShellCoreVersion: powershellVersion,
		CustomHandler:         false, // set this later from app_settings
	}}

	return result, nil
}

func ParseWebJobsStorageString(input string) (name, key string) {
	if input == "" {
		return
	}

	parts := strings.Split(input, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, "AccountName") {
			name = strings.TrimPrefix(part, "AccountName=")
		}
		if strings.HasPrefix(part, "AccountKey") {
			key = strings.TrimPrefix(part, "AccountKey=")
		}
	}

	return
}

func MergeUserAppSettings(systemSettings *[]webapps.NameValuePair, userSettings map[string]string) *[]webapps.NameValuePair {
	if len(userSettings) == 0 {
		return systemSettings
	}
	combined := *systemSettings
	for k, v := range userSettings {
		// Dedupe, explicit user settings take priority over enumerated, e.g. specifying KeyVault for `AzureWebJobsStorage`
		for i, x := range combined {
			if x.Name != nil && strings.EqualFold(*x.Name, k) {
				copy(combined[i:], combined[i+1:])
				combined = combined[:len(combined)-1]
			}
		}
		combined = append(combined, webapps.NameValuePair{
			Name:  pointer.To(k),
			Value: pointer.To(v),
		})
	}
	return &combined
}

func ExpandFunctionAppAppServiceLogs(input []FunctionAppAppServiceLogs) webapps.SiteLogsConfig {
	if len(input) == 0 {
		return webapps.SiteLogsConfig{
			Properties: &webapps.SiteLogsConfigProperties{
				HTTPLogs: &webapps.HTTPLogsConfig{
					FileSystem: &webapps.FileSystemHTTPLogsConfig{
						Enabled: pointer.To(false),
					},
				},
			},
		}
	}

	config := input[0]
	return webapps.SiteLogsConfig{
		Properties: &webapps.SiteLogsConfigProperties{
			HTTPLogs: &webapps.HTTPLogsConfig{
				FileSystem: &webapps.FileSystemHTTPLogsConfig{
					RetentionInDays: pointer.To(config.RetentionPeriodDays),
					RetentionInMb:   pointer.To(config.DiskQuotaMB),
					Enabled:         pointer.To(true),
				},
			},
		},
	}
}

func FlattenFunctionAppAppServiceLogs(input *webapps.SiteLogsConfig) []FunctionAppAppServiceLogs {
	if input == nil {
		return []FunctionAppAppServiceLogs{}
	}
	if props := input.Properties; props != nil && props.HTTPLogs != nil && props.HTTPLogs.FileSystem != nil && pointer.From(props.HTTPLogs.FileSystem.Enabled) {
		return []FunctionAppAppServiceLogs{{
			DiskQuotaMB:         pointer.From(props.HTTPLogs.FileSystem.RetentionInMb),
			RetentionPeriodDays: pointer.From(props.HTTPLogs.FileSystem.RetentionInDays),
		}}
	}

	return []FunctionAppAppServiceLogs{}
}

func ParseContentSettings(input *webapps.StringDictionary, existing map[string]string) map[string]string {
	if input == nil || input.Properties == nil {
		return nil
	}

	out := existing
	for k, v := range *input.Properties {
		switch k {
		case "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING":
			out[k] = v

		case "WEBSITE_CONTENTSHARE":
			out[k] = v
		}
	}

	return out
}
