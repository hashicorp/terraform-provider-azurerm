// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	apimValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
)

type SiteConfigWindowsFunctionAppSlot struct {
	AlwaysOn                      bool                                 `tfschema:"always_on"`
	AppCommandLine                string                               `tfschema:"app_command_line"`
	ApiDefinition                 string                               `tfschema:"api_definition_url"`
	ApiManagementConfigId         string                               `tfschema:"api_management_api_id"`
	AppInsightsInstrumentationKey string                               `tfschema:"application_insights_key"` // App Insights Instrumentation Key
	AppInsightsConnectionString   string                               `tfschema:"application_insights_connection_string"`
	AppScaleLimit                 int                                  `tfschema:"app_scale_limit"`
	AppServiceLogs                []FunctionAppAppServiceLogs          `tfschema:"app_service_logs"`
	AutoSwapSlotName              string                               `tfschema:"auto_swap_slot_name"`
	DefaultDocuments              []string                             `tfschema:"default_documents"`
	ElasticInstanceMinimum        int                                  `tfschema:"elastic_instance_minimum"`
	Http2Enabled                  bool                                 `tfschema:"http2_enabled"`
	IpRestriction                 []IpRestriction                      `tfschema:"ip_restriction"`
	LoadBalancing                 string                               `tfschema:"load_balancing_mode"` // TODO - Valid for FunctionApps?
	ManagedPipelineMode           string                               `tfschema:"managed_pipeline_mode"`
	PreWarmedInstanceCount        int                                  `tfschema:"pre_warmed_instance_count"`
	RemoteDebugging               bool                                 `tfschema:"remote_debugging_enabled"`
	RemoteDebuggingVersion        string                               `tfschema:"remote_debugging_version"`
	RuntimeScaleMonitoring        bool                                 `tfschema:"runtime_scale_monitoring_enabled"`
	ScmIpRestriction              []IpRestriction                      `tfschema:"scm_ip_restriction"`
	ScmType                       string                               `tfschema:"scm_type"` // Computed?
	ScmUseMainIpRestriction       bool                                 `tfschema:"scm_use_main_ip_restriction"`
	Use32BitWorker                bool                                 `tfschema:"use_32_bit_worker"`
	WebSockets                    bool                                 `tfschema:"websockets_enabled"`
	FtpsState                     string                               `tfschema:"ftps_state"`
	HealthCheckPath               string                               `tfschema:"health_check_path"`
	HealthCheckEvictionTime       int                                  `tfschema:"health_check_eviction_time_in_min"`
	NumberOfWorkers               int                                  `tfschema:"worker_count"`
	ApplicationStack              []ApplicationStackWindowsFunctionApp `tfschema:"application_stack"`
	MinTlsVersion                 string                               `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion              string                               `tfschema:"scm_minimum_tls_version"`
	Cors                          []CorsSetting                        `tfschema:"cors"`
	DetailedErrorLogging          bool                                 `tfschema:"detailed_error_logging_enabled"`
	WindowsFxVersion              string                               `tfschema:"windows_fx_version"`
	VnetRouteAllEnabled           bool                                 `tfschema:"vnet_route_all_enabled"` // Not supported in Dynamic plans
}

func SiteConfigSchemaWindowsFunctionAppSlot() *pluginsdk.Schema {
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

				"auto_swap_slot_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

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

				"scm_use_main_ip_restriction": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should the Windows Function App `ip_restriction` configuration be used for the SCM also.",
				},

				"scm_ip_restriction": IpRestrictionSchema(),

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
					Default:  string(web.ManagedPipelineModeIntegrated),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ManagedPipelineModeClassic),
						string(web.ManagedPipelineModeIntegrated),
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
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.FtpsStateDisabled),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.FtpsStateAllAllowed),
						string(web.FtpsStateDisabled),
						string(web.FtpsStateFtpsOnly),
					}, false),
					Description: "State of FTP / FTPS service for this function app. Possible values include: `AllAllowed`, `FtpsOnly` and `Disabled`. Defaults to `Disabled`.",
				},

				"health_check_path": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The path to be checked for this function app health.",
				},

				"health_check_eviction_time_in_min": { // NOTE: Will evict the only node in single node configurations.
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntBetween(2, 10),
					Description:  "The amount of time in minutes that a node is unhealthy before being removed from the load balancer. Possible values are between `2` and `10`. Defaults to `10`. Only valid in conjunction with `health_check_path`",
				},

				"worker_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 100),
					Description:  "The number of Workers for this Windows Function App.",
				},

				"minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.SupportedTLSVersionsOneFullStopTwo),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SupportedTLSVersionsOneFullStopZero),
						string(web.SupportedTLSVersionsOneFullStopOne),
						string(web.SupportedTLSVersionsOneFullStopTwo),
					}, false),
					Description: "The configures the minimum version of TLS required for SSL requests. Possible values include: `1.0`, `1.1`, and  `1.2`. Defaults to `1.2`.",
				},

				"scm_minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.SupportedTLSVersionsOneFullStopTwo),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SupportedTLSVersionsOneFullStopZero),
						string(web.SupportedTLSVersionsOneFullStopOne),
						string(web.SupportedTLSVersionsOneFullStopTwo),
					}, false),
					Description: "Configures the minimum version of TLS required for SSL requests to the SCM site Possible values include: `1.0`, `1.1`, and  `1.2`. Defaults to `1.2`.",
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

type SiteConfigLinuxFunctionAppSlot struct {
	AlwaysOn                      bool                               `tfschema:"always_on"`
	AppCommandLine                string                             `tfschema:"app_command_line"`
	ApiDefinition                 string                             `tfschema:"api_definition_url"`
	ApiManagementConfigId         string                             `tfschema:"api_management_api_id"`
	AppInsightsInstrumentationKey string                             `tfschema:"application_insights_key"` // App Insights Instrumentation Key
	AppInsightsConnectionString   string                             `tfschema:"application_insights_connection_string"`
	AppScaleLimit                 int                                `tfschema:"app_scale_limit"`
	AppServiceLogs                []FunctionAppAppServiceLogs        `tfschema:"app_service_logs"`
	AutoSwapSlotName              string                             `tfschema:"auto_swap_slot_name"`
	UseManagedIdentityACR         bool                               `tfschema:"container_registry_use_managed_identity"`
	ContainerRegistryMSI          string                             `tfschema:"container_registry_managed_identity_client_id"`
	DefaultDocuments              []string                           `tfschema:"default_documents"`
	ElasticInstanceMinimum        int                                `tfschema:"elastic_instance_minimum"`
	Http2Enabled                  bool                               `tfschema:"http2_enabled"`
	IpRestriction                 []IpRestriction                    `tfschema:"ip_restriction"`
	LoadBalancing                 string                             `tfschema:"load_balancing_mode"` // TODO - Valid for FunctionApps?
	ManagedPipelineMode           string                             `tfschema:"managed_pipeline_mode"`
	PreWarmedInstanceCount        int                                `tfschema:"pre_warmed_instance_count"`
	RemoteDebugging               bool                               `tfschema:"remote_debugging_enabled"`
	RemoteDebuggingVersion        string                             `tfschema:"remote_debugging_version"`
	RuntimeScaleMonitoring        bool                               `tfschema:"runtime_scale_monitoring_enabled"`
	ScmIpRestriction              []IpRestriction                    `tfschema:"scm_ip_restriction"`
	ScmType                       string                             `tfschema:"scm_type"` // Computed?
	ScmUseMainIpRestriction       bool                               `tfschema:"scm_use_main_ip_restriction"`
	Use32BitWorker                bool                               `tfschema:"use_32_bit_worker"`
	WebSockets                    bool                               `tfschema:"websockets_enabled"`
	FtpsState                     string                             `tfschema:"ftps_state"`
	HealthCheckPath               string                             `tfschema:"health_check_path"`
	HealthCheckEvictionTime       int                                `tfschema:"health_check_eviction_time_in_min"`
	WorkerCount                   int                                `tfschema:"worker_count"`
	ApplicationStack              []ApplicationStackLinuxFunctionApp `tfschema:"application_stack"`
	MinTlsVersion                 string                             `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion              string                             `tfschema:"scm_minimum_tls_version"`
	Cors                          []CorsSetting                      `tfschema:"cors"`
	DetailedErrorLogging          bool                               `tfschema:"detailed_error_logging_enabled"`
	LinuxFxVersion                string                             `tfschema:"linux_fx_version"`
	VnetRouteAllEnabled           bool                               `tfschema:"vnet_route_all_enabled"` // Not supported in Dynamic plans
}

func SiteConfigSchemaLinuxFunctionAppSlot() *pluginsdk.Schema {
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

				"auto_swap_slot_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

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

				"scm_use_main_ip_restriction": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should the Linux Function App `ip_restriction` configuration be used for the SCM also.",
				},

				"scm_ip_restriction": IpRestrictionSchema(),

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
					Default:  string(web.ManagedPipelineModeIntegrated),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ManagedPipelineModeClassic),
						string(web.ManagedPipelineModeIntegrated),
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
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.FtpsStateDisabled),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.FtpsStateAllAllowed),
						string(web.FtpsStateDisabled),
						string(web.FtpsStateFtpsOnly),
					}, false),
					Description: "State of FTP / FTPS service for this function app. Possible values include: `AllAllowed`, `FtpsOnly` and `Disabled`. Defaults to `Disabled`.",
				},

				"health_check_path": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The path to be checked for this function app health.",
				},

				"health_check_eviction_time_in_min": { // NOTE: Will evict the only node in single node configurations.
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntBetween(2, 10),
					Description:  "The amount of time in minutes that a node is unhealthy before being removed from the load balancer. Possible values are between `2` and `10`. Defaults to `10`. Only valid in conjunction with `health_check_path`",
				},

				"worker_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 100),
					Description:  "The number of Workers for this Linux Function App.",
				},

				"minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.SupportedTLSVersionsOneFullStopTwo),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SupportedTLSVersionsOneFullStopZero),
						string(web.SupportedTLSVersionsOneFullStopOne),
						string(web.SupportedTLSVersionsOneFullStopTwo),
					}, false),
					Description: "The configures the minimum version of TLS required for SSL requests. Possible values include: `1.0`, `1.1`, and  `1.2`. Defaults to `1.2`.",
				},

				"scm_minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.SupportedTLSVersionsOneFullStopTwo),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SupportedTLSVersionsOneFullStopZero),
						string(web.SupportedTLSVersionsOneFullStopOne),
						string(web.SupportedTLSVersionsOneFullStopTwo),
					}, false),
					Description: "Configures the minimum version of TLS required for SSL requests to the SCM site Possible values include: `1.0`, `1.1`, and  `1.2`. Defaults to `1.2`.",
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

func ExpandSiteConfigWindowsFunctionAppSlot(siteConfig []SiteConfigWindowsFunctionAppSlot, existing *web.SiteConfig, metadata sdk.ResourceMetaData, version string, storageString string, storageUsesMSI bool) (*web.SiteConfig, error) {
	if len(siteConfig) == 0 {
		return nil, nil
	}
	expanded := &web.SiteConfig{}
	if existing != nil {
		expanded = existing
		// need to zero fxversion to re-calculate based on changes below or removing app_stack doesn't apply
		expanded.WindowsFxVersion = utils.String("")
	}

	appSettings := make([]web.NameValuePair, 0)

	appSettings = append(appSettings, web.NameValuePair{
		Name:  utils.String("FUNCTIONS_EXTENSION_VERSION"),
		Value: utils.String(version),
	})

	if storageUsesMSI {
		appSettings = append(appSettings, web.NameValuePair{
			Name:  utils.String("AzureWebJobsStorage__accountName"),
			Value: utils.String(storageString),
		})
	} else {
		appSettings = append(appSettings, web.NameValuePair{
			Name:  utils.String("AzureWebJobsStorage"),
			Value: utils.String(storageString),
		})
	}

	windowsSlotSiteConfig := siteConfig[0]

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") || metadata.ResourceData.HasChange("site_config.0.health_check_eviction_time_in_min") {
		v := strconv.Itoa(windowsSlotSiteConfig.HealthCheckEvictionTime)
		if v == "0" || windowsSlotSiteConfig.HealthCheckPath == "" {
			appSettings = updateOrAppendAppSettings(appSettings, "WEBSITE_HEALTHCHECK_MAXPINGFAILURES", v, true)
		} else {
			appSettings = updateOrAppendAppSettings(appSettings, "WEBSITE_HEALTHCHECK_MAXPINGFAILURES", v, false)
		}
	}

	expanded.AlwaysOn = utils.Bool(windowsSlotSiteConfig.AlwaysOn)

	if metadata.ResourceData.HasChange("site_config.0.auto_swap_slot_name") {
		expanded.AutoSwapSlotName = utils.String(windowsSlotSiteConfig.AutoSwapSlotName)
	}

	if metadata.ResourceData.HasChange("site_config.0.app_scale_limit") {
		expanded.FunctionAppScaleLimit = utils.Int32(int32(windowsSlotSiteConfig.AppScaleLimit))
	}

	if windowsSlotSiteConfig.AppInsightsConnectionString != "" {
		appSettings = append(appSettings, web.NameValuePair{
			Name:  utils.String("APPLICATIONINSIGHTS_CONNECTION_STRING"),
			Value: utils.String(windowsSlotSiteConfig.AppInsightsConnectionString),
		})
	}

	if windowsSlotSiteConfig.AppInsightsInstrumentationKey != "" {
		appSettings = append(appSettings, web.NameValuePair{
			Name:  utils.String("APPINSIGHTS_INSTRUMENTATIONKEY"),
			Value: utils.String(windowsSlotSiteConfig.AppInsightsInstrumentationKey),
		})
	}

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.APIManagementConfig = &web.APIManagementConfig{
			ID: utils.String(windowsSlotSiteConfig.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.APIDefinition = &web.APIDefinitionInfo{
			URL: utils.String(windowsSlotSiteConfig.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = utils.String(windowsSlotSiteConfig.AppCommandLine)
	}

	if len(windowsSlotSiteConfig.ApplicationStack) > 0 {
		windowsAppStack := windowsSlotSiteConfig.ApplicationStack[0]
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
			expanded.WindowsFxVersion = utils.String("") // Custom needs an explicit empty string here
		}
	} else {
		appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_WORKER_RUNTIME", "", true)
		expanded.WindowsFxVersion = utils.String("")
	}

	expanded.VnetRouteAllEnabled = utils.Bool(windowsSlotSiteConfig.VnetRouteAllEnabled)

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = &windowsSlotSiteConfig.DefaultDocuments
	}

	expanded.HTTP20Enabled = utils.Bool(windowsSlotSiteConfig.Http2Enabled)

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := ExpandIpRestrictions(windowsSlotSiteConfig.IpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	expanded.ScmIPSecurityRestrictionsUseMain = utils.Bool(windowsSlotSiteConfig.ScmUseMainIpRestriction)

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := ExpandIpRestrictions(windowsSlotSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = web.SiteLoadBalancing(windowsSlotSiteConfig.LoadBalancing)
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = web.ManagedPipelineMode(windowsSlotSiteConfig.ManagedPipelineMode)
	}

	expanded.RemoteDebuggingEnabled = utils.Bool(windowsSlotSiteConfig.RemoteDebugging)

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = utils.String(windowsSlotSiteConfig.RemoteDebuggingVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.runtime_scale_monitoring_enabled") {
		expanded.FunctionsRuntimeScaleMonitoringEnabled = utils.Bool(windowsSlotSiteConfig.RuntimeScaleMonitoring)
	}

	expanded.Use32BitWorkerProcess = utils.Bool(windowsSlotSiteConfig.Use32BitWorker)

	if metadata.ResourceData.HasChange("site_config.0.websockets_enabled") {
		expanded.WebSocketsEnabled = utils.Bool(windowsSlotSiteConfig.WebSockets)
	}

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		expanded.FtpsState = web.FtpsState(windowsSlotSiteConfig.FtpsState)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = utils.String(windowsSlotSiteConfig.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.worker_count") {
		expanded.NumberOfWorkers = utils.Int32(int32(windowsSlotSiteConfig.NumberOfWorkers))
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTLSVersion = web.SupportedTLSVersions(windowsSlotSiteConfig.MinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTLSVersion = web.SupportedTLSVersions(windowsSlotSiteConfig.ScmMinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		cors := ExpandCorsSettings(windowsSlotSiteConfig.Cors)
		expanded.Cors = cors
	}

	if metadata.ResourceData.HasChange("site_config.0.pre_warmed_instance_count") {
		expanded.PreWarmedInstanceCount = utils.Int32(int32(windowsSlotSiteConfig.PreWarmedInstanceCount))
	}

	expanded.AppSettings = &appSettings

	return expanded, nil
}

func FlattenSiteConfigWindowsFunctionAppSlot(functionAppSlotSiteConfig *web.SiteConfig) (*SiteConfigWindowsFunctionAppSlot, error) {
	if functionAppSlotSiteConfig == nil {
		return nil, fmt.Errorf("flattening site config: SiteConfig was nil")
	}

	result := &SiteConfigWindowsFunctionAppSlot{
		AlwaysOn:                utils.NormaliseNilableBool(functionAppSlotSiteConfig.AlwaysOn),
		AppCommandLine:          utils.NormalizeNilableString(functionAppSlotSiteConfig.AppCommandLine),
		AppScaleLimit:           int(utils.NormaliseNilableInt32(functionAppSlotSiteConfig.FunctionAppScaleLimit)),
		AutoSwapSlotName:        utils.NormalizeNilableString(functionAppSlotSiteConfig.AutoSwapSlotName),
		Cors:                    FlattenCorsSettings(functionAppSlotSiteConfig.Cors),
		DetailedErrorLogging:    utils.NormaliseNilableBool(functionAppSlotSiteConfig.DetailedErrorLoggingEnabled),
		HealthCheckPath:         utils.NormalizeNilableString(functionAppSlotSiteConfig.HealthCheckPath),
		Http2Enabled:            utils.NormaliseNilableBool(functionAppSlotSiteConfig.HTTP20Enabled),
		WindowsFxVersion:        utils.NormalizeNilableString(functionAppSlotSiteConfig.WindowsFxVersion),
		LoadBalancing:           string(functionAppSlotSiteConfig.LoadBalancing),
		ManagedPipelineMode:     string(functionAppSlotSiteConfig.ManagedPipelineMode),
		NumberOfWorkers:         int(utils.NormaliseNilableInt32(functionAppSlotSiteConfig.NumberOfWorkers)),
		ScmType:                 string(functionAppSlotSiteConfig.ScmType),
		FtpsState:               string(functionAppSlotSiteConfig.FtpsState),
		RuntimeScaleMonitoring:  utils.NormaliseNilableBool(functionAppSlotSiteConfig.FunctionsRuntimeScaleMonitoringEnabled),
		MinTlsVersion:           string(functionAppSlotSiteConfig.MinTLSVersion),
		ScmMinTlsVersion:        string(functionAppSlotSiteConfig.ScmMinTLSVersion),
		PreWarmedInstanceCount:  int(utils.NormaliseNilableInt32(functionAppSlotSiteConfig.PreWarmedInstanceCount)),
		ElasticInstanceMinimum:  int(utils.NormaliseNilableInt32(functionAppSlotSiteConfig.MinimumElasticInstanceCount)),
		Use32BitWorker:          utils.NormaliseNilableBool(functionAppSlotSiteConfig.Use32BitWorkerProcess),
		WebSockets:              utils.NormaliseNilableBool(functionAppSlotSiteConfig.WebSocketsEnabled),
		ScmUseMainIpRestriction: utils.NormaliseNilableBool(functionAppSlotSiteConfig.ScmIPSecurityRestrictionsUseMain),
		RemoteDebugging:         utils.NormaliseNilableBool(functionAppSlotSiteConfig.RemoteDebuggingEnabled),
		RemoteDebuggingVersion:  strings.ToUpper(utils.NormalizeNilableString(functionAppSlotSiteConfig.RemoteDebuggingVersion)),
		VnetRouteAllEnabled:     utils.NormaliseNilableBool(functionAppSlotSiteConfig.VnetRouteAllEnabled),
	}

	if v := functionAppSlotSiteConfig.APIDefinition; v != nil && v.URL != nil {
		result.ApiDefinition = *v.URL
	}

	if v := functionAppSlotSiteConfig.APIManagementConfig; v != nil && v.ID != nil {
		result.ApiManagementConfigId = *v.ID
	}

	if functionAppSlotSiteConfig.IPSecurityRestrictions != nil {
		result.IpRestriction = FlattenIpRestrictions(functionAppSlotSiteConfig.IPSecurityRestrictions)
	}

	if functionAppSlotSiteConfig.ScmIPSecurityRestrictions != nil {
		result.ScmIpRestriction = FlattenIpRestrictions(functionAppSlotSiteConfig.ScmIPSecurityRestrictions)
	}

	if v := functionAppSlotSiteConfig.DefaultDocuments; v != nil {
		result.DefaultDocuments = *v
	}

	powershellVersion := ""
	if p := functionAppSlotSiteConfig.PowerShellVersion; p != nil {
		powershellVersion = *p
		if powershellVersion == "~7" {
			powershellVersion = "7"
		}
	}

	result.ApplicationStack = []ApplicationStackWindowsFunctionApp{{
		DotNetVersion:         pointer.From(functionAppSlotSiteConfig.NetFrameworkVersion),
		DotNetIsolated:        false, // Note: this is set later from app_settings.FUNCTIONS_WORKER_RUNTIME in unpackWindowsFunctionAppSettings
		NodeVersion:           "",    // Note: this will be set from app_settings later in unpackWindowsFunctionAppSettings
		JavaVersion:           pointer.From(functionAppSlotSiteConfig.JavaVersion),
		PowerShellCoreVersion: powershellVersion,
		CustomHandler:         false, // Note: this is set later from app_settings
	}}

	return result, nil
}

func ExpandSiteConfigLinuxFunctionAppSlot(siteConfig []SiteConfigLinuxFunctionAppSlot, existing *web.SiteConfig, metadata sdk.ResourceMetaData, version string, storageString string, storageUsesMSI bool) (*web.SiteConfig, error) {
	if len(siteConfig) == 0 {
		return nil, nil
	}

	expanded := &web.SiteConfig{}
	if existing != nil {
		expanded = existing
		// need to zero fxversion to re-calculate based on changes below or removing app_stack doesn't apply
		expanded.LinuxFxVersion = utils.String("")
	}

	appSettings := make([]web.NameValuePair, 0)

	appSettings = append(appSettings, web.NameValuePair{
		Name:  utils.String("FUNCTIONS_EXTENSION_VERSION"),
		Value: utils.String(version),
	})

	if storageUsesMSI {
		appSettings = append(appSettings, web.NameValuePair{
			Name:  utils.String("AzureWebJobsStorage__accountName"),
			Value: utils.String(storageString),
		})
	} else {
		appSettings = append(appSettings, web.NameValuePair{
			Name:  utils.String("AzureWebJobsStorage"),
			Value: utils.String(storageString),
		})
	}

	linuxSlotSiteConfig := siteConfig[0]

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") || metadata.ResourceData.HasChange("site_config.0.health_check_eviction_time_in_min") {
		v := strconv.Itoa(linuxSlotSiteConfig.HealthCheckEvictionTime)
		if v == "0" || linuxSlotSiteConfig.HealthCheckPath == "" {
			appSettings = updateOrAppendAppSettings(appSettings, "WEBSITE_HEALTHCHECK_MAXPINGFAILURES", v, true)
		} else {
			appSettings = updateOrAppendAppSettings(appSettings, "WEBSITE_HEALTHCHECK_MAXPINGFAILURES", v, false)
		}
	}

	expanded.AlwaysOn = utils.Bool(linuxSlotSiteConfig.AlwaysOn)

	if metadata.ResourceData.HasChange("site_config.0.auto_swap_slot_name") {
		expanded.AutoSwapSlotName = utils.String(linuxSlotSiteConfig.AutoSwapSlotName)
	}

	if metadata.ResourceData.HasChange("site_config.0.app_scale_limit") {
		expanded.FunctionAppScaleLimit = utils.Int32(int32(linuxSlotSiteConfig.AppScaleLimit))
	}

	if linuxSlotSiteConfig.AppInsightsConnectionString != "" {
		appSettings = append(appSettings, web.NameValuePair{
			Name:  utils.String("APPLICATIONINSIGHTS_CONNECTION_STRING"),
			Value: utils.String(linuxSlotSiteConfig.AppInsightsConnectionString),
		})
	}

	if linuxSlotSiteConfig.AppInsightsInstrumentationKey != "" {
		appSettings = append(appSettings, web.NameValuePair{
			Name:  utils.String("APPINSIGHTS_INSTRUMENTATIONKEY"),
			Value: utils.String(linuxSlotSiteConfig.AppInsightsInstrumentationKey),
		})
	}

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.APIManagementConfig = &web.APIManagementConfig{
			ID: utils.String(linuxSlotSiteConfig.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.APIDefinition = &web.APIDefinitionInfo{
			URL: utils.String(linuxSlotSiteConfig.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = utils.String(linuxSlotSiteConfig.AppCommandLine)
	}

	if len(linuxSlotSiteConfig.ApplicationStack) > 0 {
		linuxAppStack := linuxSlotSiteConfig.ApplicationStack[0]
		if linuxAppStack.DotNetVersion != "" {
			if linuxAppStack.DotNetIsolated {
				appSettings = append(appSettings, web.NameValuePair{
					Name:  utils.String("FUNCTIONS_WORKER_RUNTIME"),
					Value: utils.String("dotnet-isolated"),
				})
				expanded.LinuxFxVersion = utils.String(fmt.Sprintf("DOTNET-ISOLATED|%s", linuxAppStack.DotNetVersion))
			} else {
				appSettings = append(appSettings, web.NameValuePair{
					Name:  utils.String("FUNCTIONS_WORKER_RUNTIME"),
					Value: utils.String("dotnet"),
				})
				expanded.LinuxFxVersion = utils.String(fmt.Sprintf("DOTNET|%s", linuxAppStack.DotNetVersion))
			}
		}

		if linuxAppStack.NodeVersion != "" {
			appSettings = append(appSettings, web.NameValuePair{
				Name:  utils.String("FUNCTIONS_WORKER_RUNTIME"),
				Value: utils.String("node"),
			})
			appSettings = append(appSettings, web.NameValuePair{
				Name:  utils.String("WEBSITE_NODE_DEFAULT_VERSION"),
				Value: utils.String(linuxAppStack.NodeVersion),
			})
			expanded.LinuxFxVersion = utils.String(fmt.Sprintf("NODE|%s", linuxAppStack.NodeVersion))
		}

		if linuxAppStack.PythonVersion != "" {
			appSettings = append(appSettings, web.NameValuePair{
				Name:  utils.String("FUNCTIONS_WORKER_RUNTIME"),
				Value: utils.String("python"),
			})
			expanded.LinuxFxVersion = utils.String(fmt.Sprintf("Python|%s", linuxAppStack.PythonVersion))
		}

		if linuxAppStack.JavaVersion != "" {
			appSettings = append(appSettings, web.NameValuePair{
				Name:  utils.String("FUNCTIONS_WORKER_RUNTIME"),
				Value: utils.String("java"),
			})
			expanded.LinuxFxVersion = utils.String(fmt.Sprintf("Java|%s", linuxAppStack.JavaVersion))
		}

		if linuxAppStack.PowerShellCoreVersion != "" {
			appSettings = append(appSettings, web.NameValuePair{
				Name:  utils.String("FUNCTIONS_WORKER_RUNTIME"),
				Value: utils.String("powershell"),
			})
			expanded.LinuxFxVersion = utils.String(fmt.Sprintf("PowerShell|%s", linuxAppStack.PowerShellCoreVersion))
		}

		if linuxAppStack.CustomHandler {
			appSettings = append(appSettings, web.NameValuePair{
				Name:  utils.String("FUNCTIONS_WORKER_RUNTIME"),
				Value: utils.String("custom"),
			})
			expanded.LinuxFxVersion = utils.String("") // Custom needs an explicit empty string here
		}

		if linuxAppStack.Docker != nil && len(linuxAppStack.Docker) == 1 {
			dockerConfig := linuxAppStack.Docker[0]
			appSettings = append(appSettings, web.NameValuePair{
				Name:  utils.String("DOCKER_REGISTRY_SERVER_URL"),
				Value: utils.String(dockerConfig.RegistryURL),
			})
			appSettings = append(appSettings, web.NameValuePair{
				Name:  utils.String("DOCKER_REGISTRY_SERVER_USERNAME"),
				Value: utils.String(dockerConfig.RegistryUsername),
			})
			appSettings = append(appSettings, web.NameValuePair{
				Name:  utils.String("DOCKER_REGISTRY_SERVER_PASSWORD"),
				Value: utils.String(dockerConfig.RegistryPassword),
			})
			expanded.LinuxFxVersion = utils.String(fmt.Sprintf("DOCKER|%s/%s:%s", dockerConfig.RegistryURL, dockerConfig.ImageName, dockerConfig.ImageTag))
		}
	} else {
		appSettings = append(appSettings, web.NameValuePair{
			Name:  utils.String("FUNCTIONS_WORKER_RUNTIME"),
			Value: utils.String(""),
		})
		expanded.LinuxFxVersion = utils.String("")
	}

	expanded.AcrUseManagedIdentityCreds = utils.Bool(linuxSlotSiteConfig.UseManagedIdentityACR)

	expanded.VnetRouteAllEnabled = utils.Bool(linuxSlotSiteConfig.VnetRouteAllEnabled)

	if metadata.ResourceData.HasChange("site_config.0.container_registry_managed_identity_client_id") {
		expanded.AcrUserManagedIdentityID = utils.String(linuxSlotSiteConfig.ContainerRegistryMSI)
	}

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = &linuxSlotSiteConfig.DefaultDocuments
	}

	expanded.HTTP20Enabled = utils.Bool(linuxSlotSiteConfig.Http2Enabled)

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := ExpandIpRestrictions(linuxSlotSiteConfig.IpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	expanded.ScmIPSecurityRestrictionsUseMain = utils.Bool(linuxSlotSiteConfig.ScmUseMainIpRestriction)

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := ExpandIpRestrictions(linuxSlotSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = web.SiteLoadBalancing(linuxSlotSiteConfig.LoadBalancing)
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = web.ManagedPipelineMode(linuxSlotSiteConfig.ManagedPipelineMode)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_enabled") {
		expanded.RemoteDebuggingEnabled = utils.Bool(linuxSlotSiteConfig.RemoteDebugging)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = utils.String(linuxSlotSiteConfig.RemoteDebuggingVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.runtime_scale_monitoring_enabled") {
		expanded.FunctionsRuntimeScaleMonitoringEnabled = utils.Bool(linuxSlotSiteConfig.RuntimeScaleMonitoring)
	}

	expanded.Use32BitWorkerProcess = utils.Bool(linuxSlotSiteConfig.Use32BitWorker)

	expanded.WebSocketsEnabled = utils.Bool(linuxSlotSiteConfig.WebSockets)

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		expanded.FtpsState = web.FtpsState(linuxSlotSiteConfig.FtpsState)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = utils.String(linuxSlotSiteConfig.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.worker_count") {
		expanded.NumberOfWorkers = utils.Int32(int32(linuxSlotSiteConfig.WorkerCount))
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTLSVersion = web.SupportedTLSVersions(linuxSlotSiteConfig.MinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTLSVersion = web.SupportedTLSVersions(linuxSlotSiteConfig.ScmMinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		cors := ExpandCorsSettings(linuxSlotSiteConfig.Cors)
		expanded.Cors = cors
	}

	if metadata.ResourceData.HasChange("site_config.0.pre_warmed_instance_count") {
		expanded.PreWarmedInstanceCount = utils.Int32(int32(linuxSlotSiteConfig.PreWarmedInstanceCount))
	}

	expanded.AppSettings = &appSettings

	return expanded, nil
}

func FlattenSiteConfigLinuxFunctionAppSlot(functionAppSlotSiteConfig *web.SiteConfig) (*SiteConfigLinuxFunctionAppSlot, error) {
	if functionAppSlotSiteConfig == nil {
		return nil, fmt.Errorf("flattening site config: SiteConfig was nil")
	}

	result := &SiteConfigLinuxFunctionAppSlot{
		AlwaysOn:                utils.NormaliseNilableBool(functionAppSlotSiteConfig.AlwaysOn),
		AppCommandLine:          utils.NormalizeNilableString(functionAppSlotSiteConfig.AppCommandLine),
		AppScaleLimit:           int(utils.NormaliseNilableInt32(functionAppSlotSiteConfig.FunctionAppScaleLimit)),
		AutoSwapSlotName:        utils.NormalizeNilableString(functionAppSlotSiteConfig.AutoSwapSlotName),
		ContainerRegistryMSI:    utils.NormalizeNilableString(functionAppSlotSiteConfig.AcrUserManagedIdentityID),
		Cors:                    FlattenCorsSettings(functionAppSlotSiteConfig.Cors),
		DetailedErrorLogging:    utils.NormaliseNilableBool(functionAppSlotSiteConfig.DetailedErrorLoggingEnabled),
		HealthCheckPath:         utils.NormalizeNilableString(functionAppSlotSiteConfig.HealthCheckPath),
		Http2Enabled:            utils.NormaliseNilableBool(functionAppSlotSiteConfig.HTTP20Enabled),
		LinuxFxVersion:          utils.NormalizeNilableString(functionAppSlotSiteConfig.LinuxFxVersion),
		LoadBalancing:           string(functionAppSlotSiteConfig.LoadBalancing),
		ManagedPipelineMode:     string(functionAppSlotSiteConfig.ManagedPipelineMode),
		WorkerCount:             int(utils.NormaliseNilableInt32(functionAppSlotSiteConfig.NumberOfWorkers)),
		ScmType:                 string(functionAppSlotSiteConfig.ScmType),
		FtpsState:               string(functionAppSlotSiteConfig.FtpsState),
		RuntimeScaleMonitoring:  utils.NormaliseNilableBool(functionAppSlotSiteConfig.FunctionsRuntimeScaleMonitoringEnabled),
		MinTlsVersion:           string(functionAppSlotSiteConfig.MinTLSVersion),
		ScmMinTlsVersion:        string(functionAppSlotSiteConfig.ScmMinTLSVersion),
		PreWarmedInstanceCount:  int(utils.NormaliseNilableInt32(functionAppSlotSiteConfig.PreWarmedInstanceCount)),
		ElasticInstanceMinimum:  int(utils.NormaliseNilableInt32(functionAppSlotSiteConfig.MinimumElasticInstanceCount)),
		Use32BitWorker:          utils.NormaliseNilableBool(functionAppSlotSiteConfig.Use32BitWorkerProcess),
		WebSockets:              utils.NormaliseNilableBool(functionAppSlotSiteConfig.WebSocketsEnabled),
		ScmUseMainIpRestriction: utils.NormaliseNilableBool(functionAppSlotSiteConfig.ScmIPSecurityRestrictionsUseMain),
		UseManagedIdentityACR:   utils.NormaliseNilableBool(functionAppSlotSiteConfig.AcrUseManagedIdentityCreds),
		RemoteDebugging:         utils.NormaliseNilableBool(functionAppSlotSiteConfig.RemoteDebuggingEnabled),
		RemoteDebuggingVersion:  strings.ToUpper(utils.NormalizeNilableString(functionAppSlotSiteConfig.RemoteDebuggingVersion)),
		VnetRouteAllEnabled:     utils.NormaliseNilableBool(functionAppSlotSiteConfig.VnetRouteAllEnabled),
	}

	if v := functionAppSlotSiteConfig.APIDefinition; v != nil && v.URL != nil {
		result.ApiDefinition = *v.URL
	}

	if v := functionAppSlotSiteConfig.APIManagementConfig; v != nil && v.ID != nil {
		result.ApiManagementConfigId = *v.ID
	}

	if functionAppSlotSiteConfig.IPSecurityRestrictions != nil {
		result.IpRestriction = FlattenIpRestrictions(functionAppSlotSiteConfig.IPSecurityRestrictions)
	}

	if functionAppSlotSiteConfig.ScmIPSecurityRestrictions != nil {
		result.ScmIpRestriction = FlattenIpRestrictions(functionAppSlotSiteConfig.ScmIPSecurityRestrictions)
	}

	if v := functionAppSlotSiteConfig.DefaultDocuments; v != nil {
		result.DefaultDocuments = *v
	}

	var appStack []ApplicationStackLinuxFunctionApp
	if functionAppSlotSiteConfig.LinuxFxVersion != nil {
		decoded, err := DecodeFunctionAppLinuxFxVersion(*functionAppSlotSiteConfig.LinuxFxVersion)
		if err != nil {
			return nil, fmt.Errorf("flattening site config: %s", err)
		}
		appStack = decoded
	}
	result.ApplicationStack = appStack

	return result, nil
}
