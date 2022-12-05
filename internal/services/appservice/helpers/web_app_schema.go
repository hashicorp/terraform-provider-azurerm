package helpers

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	apimValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SiteConfigWindows struct {
	AlwaysOn                 bool                      `tfschema:"always_on"`
	ApiManagementConfigId    string                    `tfschema:"api_management_api_id"`
	ApiDefinition            string                    `tfschema:"api_definition_url"`
	AppCommandLine           string                    `tfschema:"app_command_line"`
	AutoHeal                 bool                      `tfschema:"auto_heal_enabled"`
	AutoHealSettings         []AutoHealSettingWindows  `tfschema:"auto_heal_setting"`
	UseManagedIdentityACR    bool                      `tfschema:"container_registry_use_managed_identity"`
	ContainerRegistryUserMSI string                    `tfschema:"container_registry_managed_identity_client_id"`
	DefaultDocuments         []string                  `tfschema:"default_documents"`
	Http2Enabled             bool                      `tfschema:"http2_enabled"`
	IpRestriction            []IpRestriction           `tfschema:"ip_restriction"`
	ScmUseMainIpRestriction  bool                      `tfschema:"scm_use_main_ip_restriction"`
	ScmIpRestriction         []IpRestriction           `tfschema:"scm_ip_restriction"`
	LoadBalancing            string                    `tfschema:"load_balancing_mode"`
	LocalMysql               bool                      `tfschema:"local_mysql_enabled"`
	ManagedPipelineMode      string                    `tfschema:"managed_pipeline_mode"`
	RemoteDebugging          bool                      `tfschema:"remote_debugging_enabled"`
	RemoteDebuggingVersion   string                    `tfschema:"remote_debugging_version"`
	ScmType                  string                    `tfschema:"scm_type"`
	Use32BitWorker           bool                      `tfschema:"use_32_bit_worker"`
	WebSockets               bool                      `tfschema:"websockets_enabled"`
	FtpsState                string                    `tfschema:"ftps_state"`
	HealthCheckPath          string                    `tfschema:"health_check_path"`
	HealthCheckEvictionTime  int                       `tfschema:"health_check_eviction_time_in_min"`
	WorkerCount              int                       `tfschema:"worker_count"`
	ApplicationStack         []ApplicationStackWindows `tfschema:"application_stack"`
	VirtualApplications      []VirtualApplication      `tfschema:"virtual_application"`
	MinTlsVersion            string                    `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion         string                    `tfschema:"scm_minimum_tls_version"`
	Cors                     []CorsSetting             `tfschema:"cors"`
	DetailedErrorLogging     bool                      `tfschema:"detailed_error_logging_enabled"`
	WindowsFxVersion         string                    `tfschema:"windows_fx_version"`
	VnetRouteAllEnabled      bool                      `tfschema:"vnet_route_all_enabled"`
	// TODO new properties / blocks
	// SiteLimits []SiteLimitsSettings `tfschema:"site_limits"` // TODO - ASE related for limiting App resource consumption
	// PushSettings - Supported in SDK, but blocked by manual step needed for connecting app to notification hub.
}

type SiteConfigLinux struct {
	AlwaysOn                bool                    `tfschema:"always_on"`
	ApiManagementConfigId   string                  `tfschema:"api_management_api_id"`
	ApiDefinition           string                  `tfschema:"api_definition_url"`
	AppCommandLine          string                  `tfschema:"app_command_line"`
	AutoHeal                bool                    `tfschema:"auto_heal_enabled"`
	AutoHealSettings        []AutoHealSettingLinux  `tfschema:"auto_heal_setting"`
	UseManagedIdentityACR   bool                    `tfschema:"container_registry_use_managed_identity"`
	ContainerRegistryMSI    string                  `tfschema:"container_registry_managed_identity_client_id"`
	DefaultDocuments        []string                `tfschema:"default_documents"`
	Http2Enabled            bool                    `tfschema:"http2_enabled"`
	IpRestriction           []IpRestriction         `tfschema:"ip_restriction"`
	ScmUseMainIpRestriction bool                    `tfschema:"scm_use_main_ip_restriction"`
	ScmIpRestriction        []IpRestriction         `tfschema:"scm_ip_restriction"`
	LoadBalancing           string                  `tfschema:"load_balancing_mode"`
	LocalMysql              bool                    `tfschema:"local_mysql_enabled"`
	ManagedPipelineMode     string                  `tfschema:"managed_pipeline_mode"`
	RemoteDebugging         bool                    `tfschema:"remote_debugging_enabled"`
	RemoteDebuggingVersion  string                  `tfschema:"remote_debugging_version"`
	ScmType                 string                  `tfschema:"scm_type"`
	Use32BitWorker          bool                    `tfschema:"use_32_bit_worker"`
	WebSockets              bool                    `tfschema:"websockets_enabled"`
	FtpsState               string                  `tfschema:"ftps_state"`
	HealthCheckPath         string                  `tfschema:"health_check_path"`
	HealthCheckEvictionTime int                     `tfschema:"health_check_eviction_time_in_min"`
	NumberOfWorkers         int                     `tfschema:"worker_count"`
	ApplicationStack        []ApplicationStackLinux `tfschema:"application_stack"`
	MinTlsVersion           string                  `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion        string                  `tfschema:"scm_minimum_tls_version"`
	Cors                    []CorsSetting           `tfschema:"cors"`
	DetailedErrorLogging    bool                    `tfschema:"detailed_error_logging_enabled"`
	LinuxFxVersion          string                  `tfschema:"linux_fx_version"`
	VnetRouteAllEnabled     bool                    `tfschema:"vnet_route_all_enabled"`
	// SiteLimits []SiteLimitsSettings `tfschema:"site_limits"` // TODO - New block to (possibly) support? No way to configure this in the portal?
}

func SiteConfigSchemaWindows() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"api_management_api_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: apimValidate.ApiID,
				},

				"api_definition_url": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				},

				"application_stack": windowsApplicationStackSchema(),

				"app_command_line": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"auto_heal_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
					RequiredWith: []string{
						"site_config.0.auto_heal_setting",
					},
				},

				"auto_heal_setting": autoHealSettingSchemaWindows(),

				"container_registry_use_managed_identity": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"container_registry_managed_identity_client_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsUUID,
				},

				"default_documents": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"http2_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"ip_restriction": IpRestrictionSchema(),

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"scm_ip_restriction": IpRestrictionSchema(),

				"local_mysql_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"load_balancing_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.SiteLoadBalancingLeastRequests),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SiteLoadBalancingLeastRequests),
						string(web.SiteLoadBalancingWeightedRoundRobin),
						string(web.SiteLoadBalancingLeastResponseTime),
						string(web.SiteLoadBalancingWeightedTotalTraffic),
						string(web.SiteLoadBalancingRequestHash),
						string(web.SiteLoadBalancingPerSiteRoundRobin),
					}, false),
				},

				"managed_pipeline_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.ManagedPipelineModeIntegrated),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ManagedPipelineModeClassic),
						string(web.ManagedPipelineModeIntegrated),
					}, false),
				},

				"remote_debugging_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"remote_debugging_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"VS2017",
						"VS2019",
					}, false),
				},

				"scm_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"use_32_bit_worker": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true, // Variable default value depending on several factors, such as plan type?
				},

				"websockets_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
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
				},

				"health_check_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"health_check_eviction_time_in_min": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(2, 10),
					Description:  "The amount of time in minutes that a node is unhealthy before being removed from the load balancer. Possible values are between `2` and `10`. Defaults to `10`. Only valid in conjunction with `health_check_path`",
				},

				"worker_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 100),
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
				},

				"cors": CorsSettingsSchema(),

				"virtual_application": virtualApplicationsSchema(),

				"vnet_route_all_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should all outbound traffic to have Virtual Network Security Groups and User Defined Routes applied? Defaults to `false`.",
				},

				"detailed_error_logging_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"linux_fx_version": {
					Type:     pluginsdk.TypeString,
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

func SiteConfigSchemaWindowsComputed() *pluginsdk.Schema {
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

				"application_stack": windowsApplicationStackSchemaComputed(),

				"app_command_line": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"auto_heal_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"auto_heal_setting": autoHealSettingSchemaWindowsComputed(),

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

				"http2_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"ip_restriction": IpRestrictionSchemaComputed(),

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"scm_ip_restriction": IpRestrictionSchemaComputed(),

				"local_mysql_enabled": {
					Type:     pluginsdk.TypeBool,
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

				"remote_debugging_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"remote_debugging_version": {
					Type:     pluginsdk.TypeString,
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

				"virtual_application": virtualApplicationsSchemaComputed(),

				"detailed_error_logging_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"windows_fx_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"vnet_route_all_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

func SiteConfigSchemaLinux() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"api_management_api_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: apimValidate.ApiID,
				},

				"api_definition_url": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				},

				"app_command_line": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"application_stack": linuxApplicationStackSchema(),

				"auto_heal_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					RequiredWith: []string{
						"site_config.0.auto_heal_setting",
					},
				},

				"auto_heal_setting": autoHealSettingSchemaLinux(),

				"container_registry_use_managed_identity": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"container_registry_managed_identity_client_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsUUID,
				},

				"default_documents": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"http2_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"ip_restriction": IpRestrictionSchema(),

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"scm_ip_restriction": IpRestrictionSchema(),

				"local_mysql_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"load_balancing_mode": {
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
				},

				"managed_pipeline_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.ManagedPipelineModeIntegrated),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ManagedPipelineModeClassic),
						string(web.ManagedPipelineModeIntegrated),
					}, false),
				},

				"remote_debugging_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"remote_debugging_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"VS2017",
						"VS2019",
					}, false),
				},

				"scm_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"use_32_bit_worker": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"websockets_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
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
				},

				"health_check_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"health_check_eviction_time_in_min": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(2, 10),
					Description:  "The amount of time in minutes that a node is unhealthy before being removed from the load balancer. Possible values are between `2` and `10`. Defaults to `10`. Only valid in conjunction with `health_check_path`",
				},

				"worker_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 100),
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
				},

				"cors": CorsSettingsSchema(),

				"vnet_route_all_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should all outbound traffic to have Virtual Network Security Groups and User Defined Routes applied? Defaults to `false`.",
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

func SiteConfigSchemaLinuxComputed() *pluginsdk.Schema {
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

				"application_stack": linuxApplicationStackSchemaComputed(),

				"auto_heal_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"auto_heal_setting": autoHealSettingSchemaLinuxComputed(),

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

				"http2_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"ip_restriction": IpRestrictionSchemaComputed(),

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"scm_ip_restriction": IpRestrictionSchemaComputed(),

				"local_mysql_enabled": {
					Type:     pluginsdk.TypeBool,
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

				"remote_debugging_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"remote_debugging_version": {
					Type:     pluginsdk.TypeString,
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

				"detailed_error_logging_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"linux_fx_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"vnet_route_all_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

type AutoHealSettingWindows struct {
	Triggers []AutoHealTriggerWindows `tfschema:"trigger"`
	Actions  []AutoHealActionWindows  `tfschema:"action"`
}

type AutoHealSettingLinux struct {
	Triggers []AutoHealTriggerLinux `tfschema:"trigger"`
	Actions  []AutoHealActionLinux  `tfschema:"action"`
}

type AutoHealTriggerWindows struct {
	Requests        []AutoHealRequestTrigger    `tfschema:"requests"`
	PrivateMemoryKB int                         `tfschema:"private_memory_kb"` // Private should be > 102400 KB (100 MB) to 13631488 KB (13 GB), defaults to 0 however and is always present.
	StatusCodes     []AutoHealStatusCodeTrigger `tfschema:"status_code"`       // 0 or more, ranges split by `-`, ranges cannot use sub-status or win32 code
	SlowRequests    []AutoHealSlowRequest       `tfschema:"slow_request"`
}

type AutoHealTriggerLinux struct {
	Requests     []AutoHealRequestTrigger    `tfschema:"requests"`
	StatusCodes  []AutoHealStatusCodeTrigger `tfschema:"status_code"` // 0 or more, ranges split by `-`, ranges cannot use sub-status or win32 code
	SlowRequests []AutoHealSlowRequest       `tfschema:"slow_request"`
}

type AutoHealRequestTrigger struct {
	Count    int    `tfschema:"count"`
	Interval string `tfschema:"interval"`
}

type AutoHealStatusCodeTrigger struct {
	StatusCodeRange string `tfschema:"status_code_range"` // Conflicts with `StatusCode`, `Win32Code`, and `SubStatus` when not a single value...
	SubStatus       int    `tfschema:"sub_status"`
	Win32Status     string `tfschema:"win32_status"`
	Path            string `tfschema:"path"`
	Count           int    `tfschema:"count"`
	Interval        string `tfschema:"interval"` // Format - hh:mm:ss
}

type AutoHealSlowRequest struct {
	TimeTaken string `tfschema:"time_taken"`
	Interval  string `tfschema:"interval"`
	Count     int    `tfschema:"count"`
	Path      string `tfschema:"path"`
}

type AutoHealActionWindows struct {
	ActionType         string                 `tfschema:"action_type"`                    // Enum
	CustomAction       []AutoHealCustomAction `tfschema:"custom_action"`                  // Max: 1, needs `action_type` to be "Custom"
	MinimumProcessTime string                 `tfschema:"minimum_process_execution_time"` // Minimum uptime for process before action will trigger
}

type AutoHealActionLinux struct {
	ActionType         string `tfschema:"action_type"`                    // Enum - Only `Recycle` allowed
	MinimumProcessTime string `tfschema:"minimum_process_execution_time"` // Minimum uptime for process before action will trigger
}

type AutoHealCustomAction struct {
	Executable string `tfschema:"executable"`
	Parameters string `tfschema:"parameters"`
}

func autoHealSettingSchemaWindows() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"trigger": autoHealTriggerSchemaWindows(),

				"action": autoHealActionSchemaWindows(),
			},
		},
		RequiredWith: []string{
			"site_config.0.auto_heal_enabled",
		},
	}
}

func autoHealSettingSchemaWindowsComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"trigger": autoHealTriggerSchemaWindowsComputed(),

				"action": autoHealActionSchemaWindowsComputed(),
			},
		},
	}
}

func autoHealSettingSchemaLinux() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"trigger": autoHealTriggerSchemaLinux(),

				"action": autoHealActionSchemaLinux(),
			},
		},
		RequiredWith: []string{
			"site_config.0.auto_heal_enabled",
		},
	}
}

func autoHealSettingSchemaLinuxComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"trigger": autoHealTriggerSchemaLinuxComputed(),

				"action": autoHealActionSchemaLinuxComputed(),
			},
		},
	}
}

func autoHealActionSchemaWindows() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"action_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.AutoHealActionTypeCustomAction),
						string(web.AutoHealActionTypeLogEvent),
						string(web.AutoHealActionTypeRecycle),
					}, false),
				},

				"custom_action": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"executable": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"parameters": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},

				"minimum_process_execution_time": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					// ValidateFunc: // TODO - Time in hh:mm:ss, because why not...
				},
			},
		},
	}
}

func autoHealActionSchemaWindowsComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"action_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"custom_action": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"executable": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"parameters": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},

				"minimum_process_execution_time": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func autoHealActionSchemaLinux() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"action_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.AutoHealActionTypeRecycle),
					}, false),
				},

				"minimum_process_execution_time": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					// ValidateFunc: // TODO - Time in hh:mm:ss, because why not...
				},
			},
		},
	}
}

func autoHealActionSchemaLinuxComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"action_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"minimum_process_execution_time": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// (@jackofallops) - trigger schemas intentionally left long-hand for now
func autoHealTriggerSchemaWindows() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"requests": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Required: true,
								// ValidateFunc: validation.IsRFC3339Time, // TODO should be hh:mm:ss - This is too loose, need to improve
							},
						},
					},
				},

				"private_memory_kb": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(102400, 13631488),
				},

				"status_code": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"status_code_range": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: nil, // TODO - status code range validation
							},

							"count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Required: true,
								// ValidateFunc: validation.IsRFC3339Time,
							},

							"sub_status": {
								Type:         pluginsdk.TypeInt,
								Optional:     true,
								ValidateFunc: nil, // TODO - no docs on this, needs investigation
							},

							"win32_status": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: nil, // TODO - no docs on this, needs investigation
							},

							"path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},

				"slow_request": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"time_taken": {
								Type:     pluginsdk.TypeString,
								Required: true,
								// ValidateFunc: validation.IsRFC3339Time,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Required: true,
								// ValidateFunc: validation.IsRFC3339Time,
							},

							"count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},

							"path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
}

func autoHealTriggerSchemaWindowsComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"requests": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"count": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},

				"private_memory_kb": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"status_code": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"status_code_range": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"count": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"sub_status": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"win32_status": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"path": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},

				"slow_request": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"time_taken": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"count": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"path": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

// (@jackofallops) - trigger schemas intentionally left long-hand for now
func autoHealTriggerSchemaLinux() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"requests": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Required: true,
								// ValidateFunc: validation.IsRFC3339Time, // TODO should be hh:mm:ss - This is too loose, need to improve?
							},
						},
					},
				},

				"status_code": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"status_code_range": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: nil, // TODO - status code range validation
							},

							"count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Required: true,
								// ValidateFunc: validation.IsRFC3339Time,
							},

							"sub_status": {
								Type:         pluginsdk.TypeInt,
								Optional:     true,
								ValidateFunc: nil, // TODO - no docs on this, needs investigation
							},

							"win32_status": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: nil, // TODO - no docs on this, needs investigation
							},

							"path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},

				"slow_request": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"time_taken": {
								Type:     pluginsdk.TypeString,
								Required: true,
								// ValidateFunc: validation.IsRFC3339Time,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Required: true,
								// ValidateFunc: validation.IsRFC3339Time,
							},

							"count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},

							"path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
}

func autoHealTriggerSchemaLinuxComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"requests": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"count": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},

				"status_code": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"status_code_range": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"count": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"sub_status": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"win32_status": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"path": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},

				"slow_request": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"time_taken": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"count": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"path": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

type VirtualApplication struct {
	VirtualPath        string             `tfschema:"virtual_path"`
	PhysicalPath       string             `tfschema:"physical_path"`
	Preload            bool               `tfschema:"preload"`
	VirtualDirectories []VirtualDirectory `tfschema:"virtual_directory"`
}

type VirtualDirectory struct {
	VirtualPath  string `tfschema:"virtual_path"`
	PhysicalPath string `tfschema:"physical_path"`
}

func virtualApplicationsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"virtual_path": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"physical_path": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"preload": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},

				"virtual_directory": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"virtual_path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"physical_path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
}

func virtualApplicationsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"virtual_path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"physical_path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"preload": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"virtual_directory": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"virtual_path": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"physical_path": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

type StorageAccount struct {
	Name        string `tfschema:"name"`
	Type        string `tfschema:"type"`
	AccountName string `tfschema:"account_name"`
	ShareName   string `tfschema:"share_name"`
	AccessKey   string `tfschema:"access_key"`
	MountPath   string `tfschema:"mount_path"`
}

func StorageAccountSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.AzureStorageTypeAzureBlob),
						string(web.AzureStorageTypeAzureFiles),
					}, false),
				},

				"account_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"share_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"access_key": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"mount_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func StorageAccountSchemaWindows() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.AzureStorageTypeAzureFiles),
					}, false),
				},

				"account_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"share_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"access_key": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"mount_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func StorageAccountSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"account_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"share_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"access_key": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},

				"mount_path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

type Backup struct {
	Name              string           `tfschema:"name"`
	StorageAccountUrl string           `tfschema:"storage_account_url"`
	Enabled           bool             `tfschema:"enabled"`
	Schedule          []BackupSchedule `tfschema:"schedule"`
}

type BackupSchedule struct {
	FrequencyInterval    int    `tfschema:"frequency_interval"`
	FrequencyUnit        string `tfschema:"frequency_unit"`
	KeepAtLeastOneBackup bool   `tfschema:"keep_at_least_one_backup"`
	RetentionPeriodDays  int    `tfschema:"retention_period_days"`
	StartTime            string `tfschema:"start_time"`
	LastExecutionTime    string `tfschema:"last_execution_time"`
}

func BackupSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name which should be used for this Backup.",
				},

				"storage_account_url": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.IsURLWithHTTPS,
					Description:  "The SAS URL to the container.",
				},

				"enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Should this backup job be enabled?",
				},

				"schedule": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"frequency_interval": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntBetween(0, 1000),
								Description:  "How often the backup should be executed (e.g. for weekly backup, this should be set to `7` and `frequency_unit` should be set to `Day`).",
							},

							"frequency_unit": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ValidateFunc: validation.StringInSlice([]string{
									"Day",
									"Hour",
								}, false),
								Description: "The unit of time for how often the backup should take place. Possible values include: `Day` and `Hour`.",
							},

							"keep_at_least_one_backup": {
								Type:        pluginsdk.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "Should the service keep at least one backup, regardless of age of backup. Defaults to `false`.",
							},

							"retention_period_days": {
								Type:         pluginsdk.TypeInt,
								Optional:     true,
								Default:      30,
								ValidateFunc: validation.IntBetween(0, 9999999),
								Description:  "After how many days backups should be deleted.",
							},

							"start_time": {
								Type:        pluginsdk.TypeString,
								Optional:    true,
								Computed:    true,
								Description: "When the schedule should start working in RFC-3339 format.",
								// DiffSuppressFunc: suppress.RFC3339Time,
								// ValidateFunc:     validation.IsRFC3339Time,
							},

							"last_execution_time": {
								Type:        pluginsdk.TypeString,
								Computed:    true,
								Description: "The time the backup was last attempted.",
							},
						},
					},
				},
			},
		},
	}
}

func BackupSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of this Backup.",
				},

				"storage_account_url": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Sensitive:   true,
					Description: "The SAS URL to the container.",
				},

				"enabled": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Is this backup job enabled?",
				},

				"schedule": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"frequency_interval": {
								Type:        pluginsdk.TypeInt,
								Computed:    true,
								Description: "How often the backup should is executed in multiples of the `frequency_unit`.",
							},

							"frequency_unit": {
								Type:        pluginsdk.TypeString,
								Computed:    true,
								Description: "The unit of time for how often the backup takes place.",
							},

							"keep_at_least_one_backup": {
								Type:        pluginsdk.TypeBool,
								Computed:    true,
								Description: "Does the service keep at least one backup, regardless of age of backup.",
							},

							"retention_period_days": {
								Type:        pluginsdk.TypeInt,
								Computed:    true,
								Description: "After how many days are backups deleted.",
							},

							"start_time": {
								Type:        pluginsdk.TypeString,
								Computed:    true,
								Description: "When the schedule should start working in RFC-3339 format.",
							},

							"last_execution_time": {
								Type:        pluginsdk.TypeString,
								Computed:    true,
								Description: "The time the backup was last attempted.",
							},
						},
					},
				},
			},
		},
	}
}

type ConnectionString struct {
	Name  string `tfschema:"name"`
	Type  string `tfschema:"type"`
	Value string `tfschema:"value"`
}

func ConnectionStringSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name which should be used for this Connection.",
				},

				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ConnectionStringTypeAPIHub),
						string(web.ConnectionStringTypeCustom),
						string(web.ConnectionStringTypeDocDb),
						string(web.ConnectionStringTypeEventHub),
						string(web.ConnectionStringTypeMySQL),
						string(web.ConnectionStringTypeNotificationHub),
						string(web.ConnectionStringTypePostgreSQL),
						string(web.ConnectionStringTypeRedisCache),
						string(web.ConnectionStringTypeServiceBus),
						string(web.ConnectionStringTypeSQLAzure),
						string(web.ConnectionStringTypeSQLServer),
					}, false),
					Description: "Type of database. Possible values include: `MySQL`, `SQLServer`, `SQLAzure`, `Custom`, `NotificationHub`, `ServiceBus`, `EventHub`, `APIHub`, `DocDb`, `RedisCache`, and `PostgreSQL`.",
				},

				"value": {
					Type:        pluginsdk.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The connection string value.",
				},
			},
		},
	}
}

func ConnectionStringSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of this Connection.",
				},

				"type": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The type of database.",
				},

				"value": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Sensitive:   true,
					Description: "The connection string value.",
				},
			},
		},
	}
}

type LogsConfig struct {
	ApplicationLogs       []ApplicationLog `tfschema:"application_logs"`
	HttpLogs              []HttpLog        `tfschema:"http_logs"`
	DetailedErrorMessages bool             `tfschema:"detailed_error_messages"`
	FailedRequestTracing  bool             `tfschema:"failed_request_tracing"`
}

type ApplicationLog struct {
	FileSystemLevel  string             `tfschema:"file_system_level"`
	AzureBlobStorage []AzureBlobStorage `tfschema:"azure_blob_storage"`
}

type AzureBlobStorage struct {
	Level           string `tfschema:"level"`
	SasUrl          string `tfschema:"sas_url"`
	RetentionInDays int    `tfschema:"retention_in_days"`
}

type HttpLog struct {
	FileSystems      []LogsFileSystem       `tfschema:"file_system"`
	AzureBlobStorage []AzureBlobStorageHttp `tfschema:"azure_blob_storage"`
}

type AzureBlobStorageHttp struct {
	SasUrl          string `tfschema:"sas_url"`
	RetentionInDays int    `tfschema:"retention_in_days"`
}

type LogsFileSystem struct {
	RetentionMB   int `tfschema:"retention_in_mb"`
	RetentionDays int `tfschema:"retention_in_days"`
}

func LogsConfigSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"application_logs": applicationLogSchema(),

				"http_logs": httpLogSchema(),

				"failed_request_tracing": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"detailed_error_messages": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func LogsConfigSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"application_logs": applicationLogSchemaComputed(),

				"http_logs": httpLogSchemaComputed(),

				"failed_request_tracing": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"detailed_error_messages": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

func applicationLogSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"file_system_level": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.LogLevelError),
						string(web.LogLevelInformation),
						string(web.LogLevelVerbose),
						string(web.LogLevelWarning),
					}, false),
				},

				"azure_blob_storage": appLogBlobStorageSchema(),
			},
		},
	}
}

func applicationLogSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"file_system_level": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"azure_blob_storage": appLogBlobStorageSchemaComputed(),
			},
		},
	}
}

func appLogBlobStorageSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"level": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.LogLevelError),
						string(web.LogLevelInformation),
						string(web.LogLevelVerbose),
						string(web.LogLevelWarning),
					}, false),
				},
				"sas_url": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"retention_in_days": {
					Type:     pluginsdk.TypeInt,
					Required: true,
					// TODO: Validation here?
				},
			},
		},
	}
}

func appLogBlobStorageSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"level": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"sas_url": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},
				"retention_in_days": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func httpLogSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"file_system": httpLogFileSystemSchema(),

				"azure_blob_storage": httpLogBlobStorageSchema(),
			},
		},
	}
}

func httpLogSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"file_system": httpLogFileSystemSchemaComputed(),

				"azure_blob_storage": httpLogBlobStorageSchemaComputed(),
			},
		},
	}
}

func httpLogFileSystemSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"logs.0.http_logs.0.azure_blob_storage"},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"retention_in_mb": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(25, 100),
				},

				"retention_in_days": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(0),
				},
			},
		},
	}
}

func httpLogFileSystemSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"retention_in_mb": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"retention_in_days": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func httpLogBlobStorageSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"logs.0.http_logs.0.file_system"},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"sas_url": {
					Type:      pluginsdk.TypeString,
					Required:  true,
					Sensitive: true,
				},
				"retention_in_days": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntAtLeast(0), // Variable validation here based on the Service Plan SKU
				},
			},
		},
	}
}

func httpLogBlobStorageSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"sas_url": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},
				"retention_in_days": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func ExpandSiteConfigWindows(siteConfig []SiteConfigWindows, existing *web.SiteConfig, metadata sdk.ResourceMetaData, servicePlan web.AppServicePlan) (*web.SiteConfig, *string, error) {
	if len(siteConfig) == 0 {
		return nil, nil, nil
	}

	expanded := &web.SiteConfig{}
	if existing != nil {
		expanded = existing
	}

	winSiteConfig := siteConfig[0]

	currentStack := ""
	if len(winSiteConfig.ApplicationStack) == 1 {
		winAppStack := winSiteConfig.ApplicationStack[0]
		currentStack = winAppStack.CurrentStack
	}

	if servicePlan.Sku != nil && servicePlan.Sku.Name != nil {
		if isFreeOrSharedServicePlan(*servicePlan.Sku.Name) {
			if winSiteConfig.AlwaysOn {
				return nil, nil, fmt.Errorf("always_on cannot be set to true when using Free, F1, D1 Sku")
			}
			if expanded.AlwaysOn != nil && *expanded.AlwaysOn {
				return nil, nil, fmt.Errorf("always_on feature has to be turned off before switching to a free/shared Sku")
			}
		}
	}
	expanded.AlwaysOn = pointer.To(winSiteConfig.AlwaysOn)

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.APIManagementConfig = &web.APIManagementConfig{
			ID: pointer.To(winSiteConfig.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.APIDefinition = &web.APIDefinitionInfo{
			URL: pointer.To(winSiteConfig.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = pointer.To(winSiteConfig.AppCommandLine)
	}

	if metadata.ResourceData.HasChange("site_config.0.application_stack") {
		if len(winSiteConfig.ApplicationStack) == 1 {
			winAppStack := winSiteConfig.ApplicationStack[0]
			// TODO - only one of these should be non-nil?
			if winAppStack.NetFrameworkVersion != "" {
				expanded.NetFrameworkVersion = pointer.To(winAppStack.NetFrameworkVersion)
				if currentStack == "" {
					currentStack = CurrentStackDotNet
				}
			}
			if winAppStack.NetCoreVersion != "" {
				expanded.NetFrameworkVersion = pointer.To(winAppStack.NetFrameworkVersion)
				if currentStack == "" {
					currentStack = CurrentStackDotNetCore
				}
			}
			if winAppStack.NodeVersion != "" {
				// Note: node version is now exclusively controlled via app_setting.WEBSITE_NODE_DEFAULT_VERSION
				if currentStack == "" {
					currentStack = CurrentStackNode
				}
			}
			if winAppStack.PhpVersion != "" {
				if winAppStack.PhpVersion != PhpVersionOff {
					expanded.PhpVersion = pointer.To(winAppStack.PhpVersion)
				} else {
					expanded.PhpVersion = pointer.To("")
				}
				if currentStack == "" {
					currentStack = CurrentStackPhp
				}
			}
			if winAppStack.PythonVersion != "" || winAppStack.Python {
				expanded.PythonVersion = pointer.To(winAppStack.PythonVersion)
				if currentStack == "" {
					currentStack = CurrentStackPython
				}
			}
			if winAppStack.JavaVersion != "" {
				expanded.JavaVersion = pointer.To(winAppStack.JavaVersion)
				if winAppStack.JavaEmbeddedServer {
					expanded.JavaContainer = pointer.To(JavaContainerEmbeddedServer)
					expanded.JavaContainerVersion = pointer.To(JavaContainerEmbeddedServerVersion)
				} else if winAppStack.TomcatVersion != "" {
					expanded.JavaContainer = pointer.To(JavaContainerTomcat)
					expanded.JavaContainerVersion = pointer.To(winAppStack.TomcatVersion)
				} else if winAppStack.JavaContainer != "" {
					expanded.JavaContainer = pointer.To(winAppStack.JavaContainer)
					expanded.JavaContainerVersion = pointer.To(winAppStack.JavaContainerVersion)
				}
				if currentStack == "" {
					currentStack = CurrentStackJava
				}
			}
			if winAppStack.DockerContainerName != "" || winAppStack.DockerContainerRegistry != "" || winAppStack.DockerContainerTag != "" {
				if winAppStack.DockerContainerRegistry != "" {
					expanded.WindowsFxVersion = pointer.To(fmt.Sprintf("DOCKER|%s/%s:%s", winAppStack.DockerContainerRegistry, winAppStack.DockerContainerName, winAppStack.DockerContainerTag))
				} else {
					expanded.WindowsFxVersion = pointer.To(fmt.Sprintf("DOCKER|%s:%s", winAppStack.DockerContainerName, winAppStack.DockerContainerTag))
				}
			}
		} else {
			expanded.WindowsFxVersion = pointer.To("")
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.virtual_application") {
		expanded.VirtualApplications = expandVirtualApplicationsForUpdate(winSiteConfig.VirtualApplications)
	} else {
		expanded.VirtualApplications = expandVirtualApplications(winSiteConfig.VirtualApplications)
	}

	expanded.AcrUseManagedIdentityCreds = pointer.To(winSiteConfig.UseManagedIdentityACR)

	if metadata.ResourceData.HasChange("site_config.0.container_registry_managed_identity_client_id") {
		expanded.AcrUserManagedIdentityID = pointer.To(winSiteConfig.ContainerRegistryUserMSI)
	}

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = &winSiteConfig.DefaultDocuments
	}

	expanded.HTTP20Enabled = pointer.To(winSiteConfig.Http2Enabled)

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := ExpandIpRestrictions(winSiteConfig.IpRestriction)
		if err != nil {
			return nil, nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	expanded.ScmIPSecurityRestrictionsUseMain = pointer.To(winSiteConfig.ScmUseMainIpRestriction)

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := ExpandIpRestrictions(winSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	expanded.LocalMySQLEnabled = pointer.To(winSiteConfig.LocalMysql)

	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = web.SiteLoadBalancing(winSiteConfig.LoadBalancing)
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = web.ManagedPipelineMode(winSiteConfig.ManagedPipelineMode)
	}

	expanded.RemoteDebuggingEnabled = pointer.To(winSiteConfig.RemoteDebugging)

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = pointer.To(winSiteConfig.RemoteDebuggingVersion)
	}

	expanded.Use32BitWorkerProcess = pointer.To(winSiteConfig.Use32BitWorker)

	expanded.WebSocketsEnabled = pointer.To(winSiteConfig.WebSockets)

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		expanded.FtpsState = web.FtpsState(winSiteConfig.FtpsState)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = pointer.To(winSiteConfig.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.worker_count") {
		expanded.NumberOfWorkers = pointer.To(int32(winSiteConfig.WorkerCount))
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTLSVersion = web.SupportedTLSVersions(winSiteConfig.MinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTLSVersion = web.SupportedTLSVersions(winSiteConfig.ScmMinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		cors := ExpandCorsSettings(winSiteConfig.Cors)
		if cors == nil {
			cors = &web.CorsSettings{
				AllowedOrigins: &[]string{},
			}
		}
		expanded.Cors = cors
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_heal_enabled") {
		expanded.AutoHealEnabled = pointer.To(winSiteConfig.AutoHeal)
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_heal_setting") {
		expanded.AutoHealRules = expandAutoHealSettingsWindows(winSiteConfig.AutoHealSettings)
	}

	if metadata.ResourceData.HasChange("site_config.0.vnet_route_all_enabled") {
		expanded.VnetRouteAllEnabled = pointer.To(winSiteConfig.VnetRouteAllEnabled)
	}

	return expanded, &currentStack, nil
}

func ExpandSiteConfigLinux(siteConfig []SiteConfigLinux, existing *web.SiteConfig, metadata sdk.ResourceMetaData, servicePlan web.AppServicePlan) (*web.SiteConfig, error) {
	if len(siteConfig) == 0 {
		return nil, nil
	}
	expanded := &web.SiteConfig{}
	if existing != nil {
		expanded = existing
	}

	linuxSiteConfig := siteConfig[0]

	if servicePlan.Sku != nil && servicePlan.Sku.Name != nil {
		if isFreeOrSharedServicePlan(*servicePlan.Sku.Name) {
			if linuxSiteConfig.AlwaysOn {
				return nil, fmt.Errorf("always_on cannot be set to true when using Free, F1, D1 Sku")
			}
			if expanded.AlwaysOn != nil && *expanded.AlwaysOn {
				return nil, fmt.Errorf("always_on feature has to be turned off before switching to a free/shared Sku")
			}
		}
	}
	expanded.AlwaysOn = pointer.To(linuxSiteConfig.AlwaysOn)

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.APIManagementConfig = &web.APIManagementConfig{
			ID: pointer.To(linuxSiteConfig.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.APIDefinition = &web.APIDefinitionInfo{
			URL: pointer.To(linuxSiteConfig.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = pointer.To(linuxSiteConfig.AppCommandLine)
	}

	if metadata.ResourceData.HasChange("site_config.0.application_stack") {
		if len(linuxSiteConfig.ApplicationStack) == 1 {
			linuxAppStack := linuxSiteConfig.ApplicationStack[0]
			if linuxAppStack.NetFrameworkVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("DOTNETCORE|%s", linuxAppStack.NetFrameworkVersion))
			}

			if linuxAppStack.GoVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("GO|%s", linuxAppStack.GoVersion))
			}

			if linuxAppStack.PhpVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("PHP|%s", linuxAppStack.PhpVersion))
			}

			if linuxAppStack.NodeVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("NODE|%s", linuxAppStack.NodeVersion))
			}

			if linuxAppStack.RubyVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("RUBY|%s", linuxAppStack.RubyVersion))
			}

			if linuxAppStack.PythonVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("PYTHON|%s", linuxAppStack.PythonVersion))
			}

			if linuxAppStack.JavaServer != "" {
				if linuxAppStack.JavaServer == "JAVA" && linuxAppStack.JavaServerVersion == "" {
					expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("%s|%s", linuxAppStack.JavaServer, linuxAppStack.JavaVersion))
				} else {
					expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("%s|%s-%s", linuxAppStack.JavaServer, linuxAppStack.JavaServerVersion, linuxAppStack.JavaVersion))
				}
			}

			if linuxAppStack.DockerImage != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("DOCKER|%s:%s", linuxAppStack.DockerImage, linuxAppStack.DockerImageTag))
			}
		} else {
			expanded.LinuxFxVersion = pointer.To("")
		}
	}

	expanded.AcrUseManagedIdentityCreds = pointer.To(linuxSiteConfig.UseManagedIdentityACR)

	if metadata.ResourceData.HasChange("site_config.0.container_registry_managed_identity_client_id") {
		expanded.AcrUserManagedIdentityID = pointer.To(linuxSiteConfig.ContainerRegistryMSI)
	}

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = &linuxSiteConfig.DefaultDocuments
	}

	expanded.HTTP20Enabled = pointer.To(linuxSiteConfig.Http2Enabled)

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := ExpandIpRestrictions(linuxSiteConfig.IpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	expanded.ScmIPSecurityRestrictionsUseMain = pointer.To(linuxSiteConfig.ScmUseMainIpRestriction)

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := ExpandIpRestrictions(linuxSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	expanded.LocalMySQLEnabled = pointer.To(linuxSiteConfig.LocalMysql)

	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = web.SiteLoadBalancing(linuxSiteConfig.LoadBalancing)
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = web.ManagedPipelineMode(linuxSiteConfig.ManagedPipelineMode)
	}

	expanded.RemoteDebuggingEnabled = pointer.To(linuxSiteConfig.RemoteDebugging)

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = pointer.To(linuxSiteConfig.RemoteDebuggingVersion)
	}

	expanded.Use32BitWorkerProcess = pointer.To(linuxSiteConfig.Use32BitWorker)

	expanded.WebSocketsEnabled = pointer.To(linuxSiteConfig.WebSockets)

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		expanded.FtpsState = web.FtpsState(linuxSiteConfig.FtpsState)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = pointer.To(linuxSiteConfig.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.worker_count") {
		expanded.NumberOfWorkers = pointer.To(int32(linuxSiteConfig.NumberOfWorkers))
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTLSVersion = web.SupportedTLSVersions(linuxSiteConfig.MinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTLSVersion = web.SupportedTLSVersions(linuxSiteConfig.ScmMinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		cors := ExpandCorsSettings(linuxSiteConfig.Cors)
		if cors == nil {
			cors = &web.CorsSettings{
				AllowedOrigins: &[]string{},
			}
		}
		expanded.Cors = cors
	}

	expanded.AutoHealEnabled = pointer.To(linuxSiteConfig.AutoHeal)

	if metadata.ResourceData.HasChange("site_config.0.auto_heal_setting") {
		expanded.AutoHealRules = expandAutoHealSettingsLinux(linuxSiteConfig.AutoHealSettings)
	}

	if metadata.ResourceData.HasChange("site_config.0.vnet_route_all_enabled") {
		expanded.VnetRouteAllEnabled = pointer.To(linuxSiteConfig.VnetRouteAllEnabled)
	}

	return expanded, nil
}

func ExpandLogsConfig(config []LogsConfig) *web.SiteLogsConfig {
	result := &web.SiteLogsConfig{}
	if len(config) == 0 {
		return result
	}

	result.SiteLogsConfigProperties = &web.SiteLogsConfigProperties{}

	logsConfig := config[0]

	if len(logsConfig.ApplicationLogs) == 1 {
		appLogs := logsConfig.ApplicationLogs[0]
		result.SiteLogsConfigProperties.ApplicationLogs = &web.ApplicationLogsConfig{
			FileSystem: &web.FileSystemApplicationLogsConfig{
				Level: web.LogLevel(appLogs.FileSystemLevel),
			},
		}
		if len(appLogs.AzureBlobStorage) == 1 {
			appLogsBlobs := appLogs.AzureBlobStorage[0]
			result.SiteLogsConfigProperties.ApplicationLogs.AzureBlobStorage = &web.AzureBlobStorageApplicationLogsConfig{
				Level:           web.LogLevel(appLogsBlobs.Level),
				SasURL:          pointer.To(appLogsBlobs.SasUrl),
				RetentionInDays: pointer.To(int32(appLogsBlobs.RetentionInDays)),
			}
		}
	}

	if len(logsConfig.HttpLogs) == 1 {
		httpLogs := logsConfig.HttpLogs[0]
		result.HTTPLogs = &web.HTTPLogsConfig{}

		if len(httpLogs.FileSystems) == 1 {
			httpLogFileSystem := httpLogs.FileSystems[0]
			result.HTTPLogs.FileSystem = &web.FileSystemHTTPLogsConfig{
				Enabled:         pointer.To(true),
				RetentionInMb:   pointer.To(int32(httpLogFileSystem.RetentionMB)),
				RetentionInDays: pointer.To(int32(httpLogFileSystem.RetentionDays)),
			}
		}

		if len(httpLogs.AzureBlobStorage) == 1 {
			httpLogsBlobStorage := httpLogs.AzureBlobStorage[0]
			result.HTTPLogs.AzureBlobStorage = &web.AzureBlobStorageHTTPLogsConfig{
				Enabled:         pointer.To(httpLogsBlobStorage.SasUrl != ""),
				SasURL:          pointer.To(httpLogsBlobStorage.SasUrl),
				RetentionInDays: pointer.To(int32(httpLogsBlobStorage.RetentionInDays)),
			}
		}
	}

	result.DetailedErrorMessages = &web.EnabledConfig{
		Enabled: pointer.To(logsConfig.DetailedErrorMessages),
	}

	result.FailedRequestsTracing = &web.EnabledConfig{
		Enabled: pointer.To(logsConfig.FailedRequestTracing),
	}

	return result
}

func ExpandBackupConfig(backupConfigs []Backup) *web.BackupRequest {
	result := &web.BackupRequest{}
	if len(backupConfigs) == 0 {
		return result
	}

	backupConfig := backupConfigs[0]
	backupSchedule := backupConfig.Schedule[0]
	result.BackupRequestProperties = &web.BackupRequestProperties{
		Enabled:           pointer.To(backupConfig.Enabled),
		BackupName:        pointer.To(backupConfig.Name),
		StorageAccountURL: pointer.To(backupConfig.StorageAccountUrl),
		BackupSchedule: &web.BackupSchedule{
			FrequencyInterval:     pointer.To(int32(backupSchedule.FrequencyInterval)),
			FrequencyUnit:         web.FrequencyUnit(backupSchedule.FrequencyUnit),
			KeepAtLeastOneBackup:  pointer.To(backupSchedule.KeepAtLeastOneBackup),
			RetentionPeriodInDays: pointer.To(int32(backupSchedule.RetentionPeriodDays)),
		},
	}

	if backupSchedule.StartTime != "" {
		dateTimeToStart, _ := time.Parse(time.RFC3339, backupSchedule.StartTime)
		result.BackupRequestProperties.BackupSchedule.StartTime = &date.Time{Time: dateTimeToStart}
	}

	return result
}

func ExpandStorageConfig(storageConfigs []StorageAccount) *web.AzureStoragePropertyDictionaryResource {
	storageAccounts := make(map[string]*web.AzureStorageInfoValue)
	result := &web.AzureStoragePropertyDictionaryResource{}
	if len(storageConfigs) == 0 {
		result.Properties = storageAccounts
		return result
	}

	for _, v := range storageConfigs {
		storageAccounts[v.Name] = &web.AzureStorageInfoValue{
			Type:        web.AzureStorageType(v.Type),
			AccountName: pointer.To(v.AccountName),
			ShareName:   pointer.To(v.ShareName),
			AccessKey:   pointer.To(v.AccessKey),
			MountPath:   pointer.To(v.MountPath),
		}
	}

	result.Properties = storageAccounts

	return result
}

func ExpandConnectionStrings(connectionStringsConfig []ConnectionString) *web.ConnectionStringDictionary {
	result := &web.ConnectionStringDictionary{}
	if len(connectionStringsConfig) == 0 {
		return result
	}

	connectionStrings := make(map[string]*web.ConnStringValueTypePair)
	for _, v := range connectionStringsConfig {
		connectionStrings[v.Name] = &web.ConnStringValueTypePair{
			Value: pointer.To(v.Value),
			Type:  web.ConnectionStringType(v.Type),
		}
	}
	result.Properties = connectionStrings

	return result
}

func expandVirtualApplications(virtualApplicationConfig []VirtualApplication) *[]web.VirtualApplication {
	if len(virtualApplicationConfig) == 0 {
		return nil
	}

	result := make([]web.VirtualApplication, 0)

	for _, v := range virtualApplicationConfig {
		virtualApp := web.VirtualApplication{
			VirtualPath:    pointer.To(v.VirtualPath),
			PhysicalPath:   pointer.To(v.PhysicalPath),
			PreloadEnabled: pointer.To(v.Preload),
		}
		if len(v.VirtualDirectories) > 0 {
			virtualDirs := make([]web.VirtualDirectory, 0)
			for _, d := range v.VirtualDirectories {
				virtualDirs = append(virtualDirs, web.VirtualDirectory{
					VirtualPath:  pointer.To(d.VirtualPath),
					PhysicalPath: pointer.To(d.PhysicalPath),
				})
			}
			virtualApp.VirtualDirectories = &virtualDirs
		}

		result = append(result, virtualApp)
	}
	return &result
}

func expandVirtualApplicationsForUpdate(virtualApplicationConfig []VirtualApplication) *[]web.VirtualApplication {
	if len(virtualApplicationConfig) == 0 {
		// to remove this block from the config we need to give the service the original default back, sending an empty struct leaves the previous config in place
		return &[]web.VirtualApplication{
			{
				VirtualPath:    pointer.To("/"),
				PhysicalPath:   pointer.To("site\\wwwroot"),
				PreloadEnabled: pointer.To(true),
			},
		}
	}

	result := make([]web.VirtualApplication, 0)

	for _, v := range virtualApplicationConfig {
		virtualApp := web.VirtualApplication{
			VirtualPath:    pointer.To(v.VirtualPath),
			PhysicalPath:   pointer.To(v.PhysicalPath),
			PreloadEnabled: pointer.To(v.Preload),
		}
		if len(v.VirtualDirectories) > 0 {
			virtualDirs := make([]web.VirtualDirectory, 0)
			for _, d := range v.VirtualDirectories {
				virtualDirs = append(virtualDirs, web.VirtualDirectory{
					VirtualPath:  pointer.To(d.VirtualPath),
					PhysicalPath: pointer.To(d.PhysicalPath),
				})
			}
			virtualApp.VirtualDirectories = &virtualDirs
		}

		result = append(result, virtualApp)
	}
	return &result
}

func FlattenBackupConfig(backupRequest web.BackupRequest) []Backup {
	if backupRequest.BackupRequestProperties == nil {
		return nil
	}
	props := *backupRequest.BackupRequestProperties
	backup := Backup{}
	if props.BackupName != nil {
		backup.Name = *props.BackupName
	}

	if props.StorageAccountURL != nil {
		backup.StorageAccountUrl = *props.StorageAccountURL
	}

	if props.Enabled != nil {
		backup.Enabled = *props.Enabled
	}

	if schedule := props.BackupSchedule; schedule != nil {
		backupSchedule := BackupSchedule{
			FrequencyUnit: string(schedule.FrequencyUnit),
		}
		if schedule.FrequencyInterval != nil {
			backupSchedule.FrequencyInterval = int(*schedule.FrequencyInterval)
		}

		if schedule.KeepAtLeastOneBackup != nil {
			backupSchedule.KeepAtLeastOneBackup = *schedule.KeepAtLeastOneBackup
		}

		if schedule.RetentionPeriodInDays != nil {
			backupSchedule.RetentionPeriodDays = int(*schedule.RetentionPeriodInDays)
		}

		if schedule.StartTime != nil && !schedule.StartTime.IsZero() {
			backupSchedule.StartTime = schedule.StartTime.Format(time.RFC3339)
		}

		if schedule.LastExecutionTime != nil && !schedule.LastExecutionTime.IsZero() {
			backupSchedule.LastExecutionTime = schedule.LastExecutionTime.Format(time.RFC3339)
		}

		backup.Schedule = []BackupSchedule{backupSchedule}
	}

	return []Backup{backup}
}

func FlattenLogsConfig(logsConfig web.SiteLogsConfig) []LogsConfig {
	if logsConfig.SiteLogsConfigProperties == nil {
		return nil
	}
	props := *logsConfig.SiteLogsConfigProperties
	if onlyDefaultLoggingConfig(props) {
		return nil
	}

	logs := LogsConfig{}

	if props.ApplicationLogs != nil {
		appLogs := *props.ApplicationLogs
		applicationLog := ApplicationLog{}

		if appLogs.FileSystem != nil && appLogs.FileSystem.Level != web.LogLevelOff {
			applicationLog.FileSystemLevel = string(appLogs.FileSystem.Level)
			if appLogs.AzureBlobStorage != nil && appLogs.AzureBlobStorage.Level != web.LogLevelOff {
				blobStorage := AzureBlobStorage{
					Level: string(appLogs.AzureBlobStorage.Level),
				}

				blobStorage.SasUrl = pointer.From(appLogs.AzureBlobStorage.SasURL)

				blobStorage.RetentionInDays = int(pointer.From(appLogs.AzureBlobStorage.RetentionInDays))

				applicationLog.AzureBlobStorage = []AzureBlobStorage{blobStorage}
			}
			logs.ApplicationLogs = []ApplicationLog{applicationLog}
		}
	}

	if props.HTTPLogs != nil {
		httpLogs := *props.HTTPLogs
		httpLog := HttpLog{}

		if httpLogs.FileSystem != nil && (httpLogs.FileSystem.Enabled != nil && *httpLogs.FileSystem.Enabled) {
			fileSystem := LogsFileSystem{}
			if httpLogs.FileSystem.RetentionInMb != nil {
				fileSystem.RetentionMB = int(*httpLogs.FileSystem.RetentionInMb)
			}

			if httpLogs.FileSystem.RetentionInDays != nil {
				fileSystem.RetentionDays = int(*httpLogs.FileSystem.RetentionInDays)
			}

			httpLog.FileSystems = []LogsFileSystem{fileSystem}
		}

		if httpLogs.AzureBlobStorage != nil && (httpLogs.AzureBlobStorage.Enabled != nil && *httpLogs.AzureBlobStorage.Enabled) {
			blobStorage := AzureBlobStorageHttp{}
			if httpLogs.AzureBlobStorage.SasURL != nil {
				blobStorage.SasUrl = *httpLogs.AzureBlobStorage.SasURL
			}

			if httpLogs.AzureBlobStorage.RetentionInDays != nil {
				blobStorage.RetentionInDays = int(*httpLogs.AzureBlobStorage.RetentionInDays)
			}

			if blobStorage.RetentionInDays != 0 && blobStorage.SasUrl != "" {
				httpLog.AzureBlobStorage = []AzureBlobStorageHttp{blobStorage}
			}
		}

		if httpLog.FileSystems != nil || httpLog.AzureBlobStorage != nil {
			logs.HttpLogs = []HttpLog{httpLog}
		}
	}

	// logs.DetailedErrorMessages = false
	if props.DetailedErrorMessages != nil && props.DetailedErrorMessages.Enabled != nil {
		logs.DetailedErrorMessages = *props.DetailedErrorMessages.Enabled
	}

	// logs.FailedRequestTracing = false
	if props.FailedRequestsTracing != nil && props.FailedRequestsTracing.Enabled != nil {
		logs.FailedRequestTracing = *props.FailedRequestsTracing.Enabled
	}

	return []LogsConfig{logs}
}

func onlyDefaultLoggingConfig(props web.SiteLogsConfigProperties) bool {
	if props.ApplicationLogs == nil || props.HTTPLogs == nil || props.FailedRequestsTracing == nil || props.DetailedErrorMessages == nil {
		return false
	}
	if props.ApplicationLogs.FileSystem != nil && props.ApplicationLogs.FileSystem.Level != web.LogLevelOff {
		return false
	}
	if props.ApplicationLogs.AzureBlobStorage != nil && props.ApplicationLogs.AzureBlobStorage.Level != web.LogLevelOff {
		return false
	}
	if props.HTTPLogs.FileSystem != nil && props.HTTPLogs.FileSystem.Enabled != nil && (*props.HTTPLogs.FileSystem.Enabled) {
		return false
	}
	if props.HTTPLogs.AzureBlobStorage != nil && props.HTTPLogs.AzureBlobStorage.Enabled != nil && (*props.HTTPLogs.AzureBlobStorage.Enabled) {
		return false
	}
	if props.FailedRequestsTracing.Enabled == nil || *props.FailedRequestsTracing.Enabled {
		return false
	}
	if props.DetailedErrorMessages.Enabled == nil || *props.DetailedErrorMessages.Enabled {
		return false
	}
	return true
}

func FlattenSiteConfigWindows(appSiteConfig *web.SiteConfig, currentStack string, healthCheckCount *int) []SiteConfigWindows {
	if appSiteConfig == nil {
		return nil
	}

	siteConfig := SiteConfigWindows{
		AlwaysOn:                 pointer.From(appSiteConfig.AlwaysOn),
		AppCommandLine:           pointer.From(appSiteConfig.AppCommandLine),
		AutoHeal:                 pointer.From(appSiteConfig.AutoHealEnabled),
		AutoHealSettings:         flattenAutoHealSettingsWindows(appSiteConfig.AutoHealRules),
		ContainerRegistryUserMSI: pointer.From(appSiteConfig.AcrUserManagedIdentityID),
		DetailedErrorLogging:     pointer.From(appSiteConfig.DetailedErrorLoggingEnabled),
		FtpsState:                string(appSiteConfig.FtpsState),
		HealthCheckPath:          pointer.From(appSiteConfig.HealthCheckPath),
		HealthCheckEvictionTime:  pointer.From(healthCheckCount),
		Http2Enabled:             pointer.From(appSiteConfig.HTTP20Enabled),
		IpRestriction:            FlattenIpRestrictions(appSiteConfig.IPSecurityRestrictions),
		LoadBalancing:            string(appSiteConfig.LoadBalancing),
		LocalMysql:               pointer.From(appSiteConfig.LocalMySQLEnabled),
		ManagedPipelineMode:      string(appSiteConfig.ManagedPipelineMode),
		MinTlsVersion:            string(appSiteConfig.MinTLSVersion),
		WorkerCount:              int(pointer.From(appSiteConfig.NumberOfWorkers)),
		RemoteDebugging:          pointer.From(appSiteConfig.RemoteDebuggingEnabled),
		RemoteDebuggingVersion:   strings.ToUpper(pointer.From(appSiteConfig.RemoteDebuggingVersion)),
		ScmIpRestriction:         FlattenIpRestrictions(appSiteConfig.ScmIPSecurityRestrictions),
		ScmMinTlsVersion:         string(appSiteConfig.ScmMinTLSVersion),
		ScmType:                  string(appSiteConfig.ScmType),
		ScmUseMainIpRestriction:  pointer.From(appSiteConfig.ScmIPSecurityRestrictionsUseMain),
		Use32BitWorker:           pointer.From(appSiteConfig.Use32BitWorkerProcess),
		UseManagedIdentityACR:    pointer.From(appSiteConfig.AcrUseManagedIdentityCreds),
		VirtualApplications:      flattenVirtualApplications(appSiteConfig.VirtualApplications),
		WebSockets:               pointer.From(appSiteConfig.WebSocketsEnabled),
		VnetRouteAllEnabled:      pointer.From(appSiteConfig.VnetRouteAllEnabled),
	}

	if appSiteConfig.APIManagementConfig != nil && appSiteConfig.APIManagementConfig.ID != nil {
		siteConfig.ApiManagementConfigId = *appSiteConfig.APIManagementConfig.ID
	}

	if appSiteConfig.APIDefinition != nil && appSiteConfig.APIDefinition.URL != nil {
		siteConfig.ApiDefinition = *appSiteConfig.APIDefinition.URL
	}

	if appSiteConfig.DefaultDocuments != nil {
		siteConfig.DefaultDocuments = *appSiteConfig.DefaultDocuments
	}

	if appSiteConfig.NumberOfWorkers != nil {
		siteConfig.WorkerCount = int(*appSiteConfig.NumberOfWorkers)
	}

	var winAppStack ApplicationStackWindows
	winAppStack.NetFrameworkVersion = pointer.From(appSiteConfig.NetFrameworkVersion)
	if currentStack == CurrentStackDotNetCore {
		winAppStack.NetCoreVersion = pointer.From(appSiteConfig.NetFrameworkVersion)
	}
	winAppStack.PhpVersion = pointer.From(appSiteConfig.PhpVersion)
	if winAppStack.PhpVersion == "" {
		winAppStack.PhpVersion = PhpVersionOff
	}
	winAppStack.NodeVersion = pointer.From(appSiteConfig.NodeVersion)     // TODO - Get from app_settings
	winAppStack.PythonVersion = pointer.From(appSiteConfig.PythonVersion) // This _should_ always be `""`
	winAppStack.JavaVersion = pointer.From(appSiteConfig.JavaVersion)
	winAppStack.JavaContainer = pointer.From(appSiteConfig.JavaContainer)
	winAppStack.JavaContainerVersion = pointer.From(appSiteConfig.JavaContainerVersion)
	if strings.EqualFold(winAppStack.JavaContainer, JavaContainerEmbeddedServer) {
		winAppStack.JavaEmbeddedServer = true
	}

	siteConfig.WindowsFxVersion = pointer.From(appSiteConfig.WindowsFxVersion)
	if siteConfig.WindowsFxVersion != "" {
		// Decode the string to docker values
		parts := strings.Split(strings.TrimPrefix(siteConfig.WindowsFxVersion, "DOCKER|"), ":")
		if len(parts) == 2 {
			winAppStack.DockerContainerTag = parts[1]
			path := strings.Split(parts[0], "/")
			if len(path) > 1 {
				winAppStack.DockerContainerRegistry = path[0]
				winAppStack.DockerContainerName = strings.TrimPrefix(parts[0], fmt.Sprintf("%s/", path[0]))
			} else {
				winAppStack.DockerContainerName = path[0]
			}
		}
	}
	winAppStack.CurrentStack = currentStack

	siteConfig.ApplicationStack = []ApplicationStackWindows{winAppStack}

	if appSiteConfig.Cors != nil {
		cors := CorsSetting{}
		corsSettings := appSiteConfig.Cors
		if corsSettings.SupportCredentials != nil {
			cors.SupportCredentials = *corsSettings.SupportCredentials
		}

		if corsSettings.AllowedOrigins != nil && len(*corsSettings.AllowedOrigins) != 0 {
			cors.AllowedOrigins = *corsSettings.AllowedOrigins
		}
		siteConfig.Cors = []CorsSetting{cors}
	}

	return []SiteConfigWindows{siteConfig}
}

func FlattenSiteConfigLinux(appSiteConfig *web.SiteConfig, healthCheckCount *int) []SiteConfigLinux {
	if appSiteConfig == nil {
		return nil
	}

	siteConfig := SiteConfigLinux{
		AlwaysOn:                pointer.From(appSiteConfig.AlwaysOn),
		AppCommandLine:          pointer.From(appSiteConfig.AppCommandLine),
		AutoHeal:                pointer.From(appSiteConfig.AutoHealEnabled),
		AutoHealSettings:        flattenAutoHealSettingsLinux(appSiteConfig.AutoHealRules),
		ContainerRegistryMSI:    pointer.From(appSiteConfig.AcrUserManagedIdentityID),
		DetailedErrorLogging:    pointer.From(appSiteConfig.DetailedErrorLoggingEnabled),
		Http2Enabled:            pointer.From(appSiteConfig.HTTP20Enabled),
		IpRestriction:           FlattenIpRestrictions(appSiteConfig.IPSecurityRestrictions),
		ManagedPipelineMode:     string(appSiteConfig.ManagedPipelineMode),
		ScmType:                 string(appSiteConfig.ScmType),
		FtpsState:               string(appSiteConfig.FtpsState),
		HealthCheckPath:         pointer.From(appSiteConfig.HealthCheckPath),
		HealthCheckEvictionTime: pointer.From(healthCheckCount),
		LoadBalancing:           string(appSiteConfig.LoadBalancing),
		LocalMysql:              pointer.From(appSiteConfig.LocalMySQLEnabled),
		MinTlsVersion:           string(appSiteConfig.MinTLSVersion),
		NumberOfWorkers:         int(pointer.From(appSiteConfig.NumberOfWorkers)),
		RemoteDebugging:         pointer.From(appSiteConfig.RemoteDebuggingEnabled),
		RemoteDebuggingVersion:  strings.ToUpper(pointer.From(appSiteConfig.RemoteDebuggingVersion)),
		ScmIpRestriction:        FlattenIpRestrictions(appSiteConfig.ScmIPSecurityRestrictions),
		ScmMinTlsVersion:        string(appSiteConfig.ScmMinTLSVersion),
		ScmUseMainIpRestriction: pointer.From(appSiteConfig.ScmIPSecurityRestrictionsUseMain),
		Use32BitWorker:          pointer.From(appSiteConfig.Use32BitWorkerProcess),
		UseManagedIdentityACR:   pointer.From(appSiteConfig.AcrUseManagedIdentityCreds),
		WebSockets:              pointer.From(appSiteConfig.WebSocketsEnabled),
		VnetRouteAllEnabled:     pointer.From(appSiteConfig.VnetRouteAllEnabled),
	}

	if appSiteConfig.APIManagementConfig != nil && appSiteConfig.APIManagementConfig.ID != nil {
		siteConfig.ApiManagementConfigId = *appSiteConfig.APIManagementConfig.ID
	}

	if appSiteConfig.APIDefinition != nil && appSiteConfig.APIDefinition.URL != nil {
		siteConfig.ApiDefinition = *appSiteConfig.APIDefinition.URL
	}

	if appSiteConfig.DefaultDocuments != nil {
		siteConfig.DefaultDocuments = *appSiteConfig.DefaultDocuments
	}

	if appSiteConfig.LinuxFxVersion != nil {
		var linuxAppStack ApplicationStackLinux
		siteConfig.LinuxFxVersion = *appSiteConfig.LinuxFxVersion
		// Decode the string to docker values
		linuxAppStack = decodeApplicationStackLinux(siteConfig.LinuxFxVersion)
		siteConfig.ApplicationStack = []ApplicationStackLinux{linuxAppStack}
	}

	if appSiteConfig.Cors != nil {
		corsSettings := appSiteConfig.Cors
		cors := CorsSetting{}
		if corsSettings.SupportCredentials != nil {
			cors.SupportCredentials = *corsSettings.SupportCredentials
		}

		if corsSettings.AllowedOrigins != nil && len(*corsSettings.AllowedOrigins) != 0 {
			cors.AllowedOrigins = *corsSettings.AllowedOrigins
		}
		siteConfig.Cors = []CorsSetting{cors}
	}

	return []SiteConfigLinux{siteConfig}
}

func FlattenStorageAccounts(appStorageAccounts web.AzureStoragePropertyDictionaryResource) []StorageAccount {
	if len(appStorageAccounts.Properties) == 0 {
		return nil
	}
	var storageAccounts []StorageAccount
	for k, v := range appStorageAccounts.Properties {
		storageAccount := StorageAccount{
			Name: k,
			Type: string(v.Type),
		}
		if v.AccountName != nil {
			storageAccount.AccountName = *v.AccountName
		}

		if v.ShareName != nil {
			storageAccount.ShareName = *v.ShareName
		}

		if v.AccessKey != nil {
			storageAccount.AccessKey = *v.AccessKey
		}

		if v.MountPath != nil {
			storageAccount.MountPath = *v.MountPath
		}

		storageAccounts = append(storageAccounts, storageAccount)
	}

	return storageAccounts
}

func FlattenConnectionStrings(appConnectionStrings web.ConnectionStringDictionary) []ConnectionString {
	if len(appConnectionStrings.Properties) == 0 {
		return nil
	}
	var connectionStrings []ConnectionString
	for k, v := range appConnectionStrings.Properties {
		connectionString := ConnectionString{
			Name: k,
			Type: string(v.Type),
		}
		if v.Value != nil {
			connectionString.Value = *v.Value
		}
		connectionStrings = append(connectionStrings, connectionString)
	}

	return connectionStrings
}

func ExpandAppSettingsForUpdate(settings map[string]string) *web.StringDictionary {
	appSettings := make(map[string]*string)
	for k, v := range settings {
		appSettings[k] = pointer.To(v)
	}

	return &web.StringDictionary{
		Properties: appSettings,
	}
}

func ExpandAppSettingsForCreate(settings map[string]string) *[]web.NameValuePair {
	if len(settings) > 0 {
		result := make([]web.NameValuePair, 0)
		for k, v := range settings {
			result = append(result, web.NameValuePair{
				Name:  pointer.To(k),
				Value: pointer.To(v),
			})
		}
		return &result
	}
	return nil
}

func FlattenAppSettings(input web.StringDictionary) (map[string]string, *int) {
	maxPingFailures := "WEBSITE_HEALTHCHECK_MAXPINGFAILURES"
	unmanagedSettings := []string{
		"DIAGNOSTICS_AZUREBLOBCONTAINERSASURL",
		"DIAGNOSTICS_AZUREBLOBRETENTIONINDAYS",
		"WEBSITE_HTTPLOGGING_CONTAINER_URL",
		"WEBSITE_HTTPLOGGING_RETENTION_DAYS",
		"WEBSITE_VNET_ROUTE_ALL",
		"spring.datasource.password",
		"spring.datasource.url",
		"spring.datasource.username",
		maxPingFailures,
	}

	var healthCheckCount *int
	appSettings := FlattenWebStringDictionary(input)
	if v, ok := appSettings[maxPingFailures]; ok {
		h, _ := strconv.Atoi(v)
		healthCheckCount = &h
	}

	// Remove the settings the service adds for legacy reasons.
	for _, v := range unmanagedSettings { //nolint:typecheck
		delete(appSettings, v)
	}

	return appSettings, healthCheckCount
}

func flattenVirtualApplications(appVirtualApplications *[]web.VirtualApplication) []VirtualApplication {
	if appVirtualApplications == nil || onlyDefaultVirtualApplication(*appVirtualApplications) {
		return nil
	}

	var virtualApplications []VirtualApplication
	for _, v := range *appVirtualApplications {
		virtualApp := VirtualApplication{
			VirtualPath:  pointer.From(v.VirtualPath),
			PhysicalPath: pointer.From(v.PhysicalPath),
		}
		if preload := v.PreloadEnabled; preload != nil {
			virtualApp.Preload = *preload
		}
		if v.VirtualDirectories != nil && len(*v.VirtualDirectories) > 0 {
			virtualDirs := make([]VirtualDirectory, 0)
			for _, d := range *v.VirtualDirectories {
				virtualDir := VirtualDirectory{
					VirtualPath:  pointer.From(d.VirtualPath),
					PhysicalPath: pointer.From(d.PhysicalPath),
				}
				virtualDirs = append(virtualDirs, virtualDir)
			}
			virtualApp.VirtualDirectories = virtualDirs
		}
		virtualApplications = append(virtualApplications, virtualApp)
	}

	return virtualApplications
}

func onlyDefaultVirtualApplication(input []web.VirtualApplication) bool {
	if len(input) > 1 {
		return false
	}
	app := input[0]
	if app.VirtualPath == nil || app.PhysicalPath == nil {
		return false
	}
	if *app.VirtualPath == "/" && *app.PhysicalPath == "site\\wwwroot" && *app.PreloadEnabled && app.VirtualDirectories == nil {
		return true
	}
	return false
}

func expandAutoHealSettingsWindows(autoHealSettings []AutoHealSettingWindows) *web.AutoHealRules {
	if len(autoHealSettings) == 0 {
		return &web.AutoHealRules{}
	}

	result := &web.AutoHealRules{
		Triggers: &web.AutoHealTriggers{},
		Actions:  &web.AutoHealActions{},
	}

	autoHeal := autoHealSettings[0]

	triggers := autoHeal.Triggers[0]
	if len(triggers.Requests) == 1 {
		result.Triggers.Requests = &web.RequestsBasedTrigger{
			Count:        pointer.To(int32(triggers.Requests[0].Count)),
			TimeInterval: pointer.To(triggers.Requests[0].Interval),
		}
	}

	if len(triggers.SlowRequests) == 1 {
		result.Triggers.SlowRequests = &web.SlowRequestsBasedTrigger{
			TimeTaken:    pointer.To(triggers.SlowRequests[0].TimeTaken),
			TimeInterval: pointer.To(triggers.SlowRequests[0].Interval),
			Count:        pointer.To(int32(triggers.SlowRequests[0].Count)),
		}
		if triggers.SlowRequests[0].Path != "" {
			result.Triggers.SlowRequests.Path = pointer.To(triggers.SlowRequests[0].Path)
		}
	}

	if triggers.PrivateMemoryKB != 0 {
		result.Triggers.PrivateBytesInKB = pointer.To(int32(triggers.PrivateMemoryKB))
	}

	if len(triggers.StatusCodes) > 0 {
		statusCodeTriggers := make([]web.StatusCodesBasedTrigger, 0)
		statusCodeRangeTriggers := make([]web.StatusCodesRangeBasedTrigger, 0)
		for _, s := range triggers.StatusCodes {
			statusCodeTrigger := web.StatusCodesBasedTrigger{}
			statusCodeRangeTrigger := web.StatusCodesRangeBasedTrigger{}
			parts := strings.Split(s.StatusCodeRange, "-")
			if len(parts) == 2 {
				statusCodeRangeTrigger.StatusCodes = pointer.To(s.StatusCodeRange)
				statusCodeRangeTrigger.Count = pointer.To(int32(s.Count))
				statusCodeRangeTrigger.TimeInterval = pointer.To(s.Interval)
				if s.Path != "" {
					statusCodeRangeTrigger.Path = pointer.To(s.Path)
				}
				statusCodeRangeTriggers = append(statusCodeRangeTriggers, statusCodeRangeTrigger)
			} else {
				statusCode, err := strconv.Atoi(s.StatusCodeRange)
				if err == nil {
					statusCodeTrigger.Status = pointer.To(int32(statusCode))
				}
				statusCodeTrigger.Count = pointer.To(int32(s.Count))
				statusCodeTrigger.TimeInterval = pointer.To(s.Interval)
				if s.Path != "" {
					statusCodeTrigger.Path = pointer.To(s.Path)
				}
				statusCodeTriggers = append(statusCodeTriggers, statusCodeTrigger)
			}
		}
		result.Triggers.StatusCodes = &statusCodeTriggers
		result.Triggers.StatusCodesRange = &statusCodeRangeTriggers
	}

	action := autoHeal.Actions[0]
	result.Actions.ActionType = web.AutoHealActionType(action.ActionType)
	result.Actions.MinProcessExecutionTime = pointer.To(action.MinimumProcessTime)
	if len(action.CustomAction) != 0 {
		customAction := action.CustomAction[0]
		result.Actions.CustomAction = &web.AutoHealCustomAction{
			Exe:        pointer.To(customAction.Executable),
			Parameters: pointer.To(customAction.Parameters),
		}
	}

	return result
}

func flattenAutoHealSettingsWindows(autoHealRules *web.AutoHealRules) []AutoHealSettingWindows {
	if autoHealRules == nil {
		return nil
	}

	result := AutoHealSettingWindows{}
	// Triggers
	if autoHealRules.Triggers != nil {
		resultTrigger := AutoHealTriggerWindows{}
		triggers := *autoHealRules.Triggers
		if triggers.Requests != nil {
			count := 0
			if triggers.Requests.Count != nil {
				count = int(*triggers.Requests.Count)
			}
			resultTrigger.Requests = []AutoHealRequestTrigger{{
				Count:    count,
				Interval: pointer.From(triggers.Requests.TimeInterval),
			}}
		}

		if privateBytes := triggers.PrivateBytesInKB; privateBytes != nil && *privateBytes != 0 {
			resultTrigger.PrivateMemoryKB = int(*triggers.PrivateBytesInKB)
		}

		statusCodeTriggers := make([]AutoHealStatusCodeTrigger, 0)
		if triggers.StatusCodes != nil {
			for _, s := range *triggers.StatusCodes {
				t := AutoHealStatusCodeTrigger{
					Interval: pointer.From(s.TimeInterval),
					Path:     pointer.From(s.Path),
				}

				if s.Status != nil {
					t.StatusCodeRange = strconv.Itoa(int(*s.Status))
				}

				if s.Count != nil {
					t.Count = int(*s.Count)
				}

				if s.SubStatus != nil {
					t.SubStatus = int(*s.SubStatus)
				}
				statusCodeTriggers = append(statusCodeTriggers, t)
			}
		}
		if triggers.StatusCodesRange != nil {
			for _, s := range *triggers.StatusCodesRange {
				t := AutoHealStatusCodeTrigger{
					Interval: pointer.From(s.TimeInterval),
					Path:     pointer.From(s.Path),
				}
				if s.Count != nil {
					t.Count = int(*s.Count)
				}

				if s.StatusCodes != nil {
					t.StatusCodeRange = *s.StatusCodes
				}
				statusCodeTriggers = append(statusCodeTriggers, t)
			}
		}
		resultTrigger.StatusCodes = statusCodeTriggers

		slowRequestTriggers := make([]AutoHealSlowRequest, 0)
		if triggers.SlowRequests != nil {
			slowRequestTriggers = append(slowRequestTriggers, AutoHealSlowRequest{
				TimeTaken: pointer.From(triggers.SlowRequests.TimeTaken),
				Interval:  pointer.From(triggers.SlowRequests.TimeInterval),
				Count:     int(pointer.From(triggers.SlowRequests.Count)),
				Path:      pointer.From(triggers.SlowRequests.Path),
			})
		}
		resultTrigger.SlowRequests = slowRequestTriggers
		result.Triggers = []AutoHealTriggerWindows{resultTrigger}
	}

	// Actions
	if autoHealRules.Actions != nil {
		actions := *autoHealRules.Actions
		customActions := make([]AutoHealCustomAction, 0)
		if actions.CustomAction != nil {
			customActions = append(customActions, AutoHealCustomAction{
				Executable: pointer.From(actions.CustomAction.Exe),
				Parameters: pointer.From(actions.CustomAction.Parameters),
			})
		}

		resultActions := AutoHealActionWindows{
			ActionType:         string(actions.ActionType),
			CustomAction:       customActions,
			MinimumProcessTime: pointer.From(actions.MinProcessExecutionTime),
		}
		result.Actions = []AutoHealActionWindows{resultActions}
	}

	if result.Actions != nil || result.Triggers != nil {
		return []AutoHealSettingWindows{result}
	}

	return nil
}

func expandAutoHealSettingsLinux(autoHealSettings []AutoHealSettingLinux) *web.AutoHealRules {
	if len(autoHealSettings) == 0 {
		return nil
	}

	result := &web.AutoHealRules{
		Triggers: &web.AutoHealTriggers{},
		Actions:  &web.AutoHealActions{},
	}

	autoHeal := autoHealSettings[0]

	triggers := autoHeal.Triggers[0]
	if len(triggers.Requests) == 1 {
		result.Triggers.Requests = &web.RequestsBasedTrigger{
			Count:        pointer.To(int32(triggers.Requests[0].Count)),
			TimeInterval: pointer.To(triggers.Requests[0].Interval),
		}
	}

	if len(triggers.SlowRequests) == 1 {
		result.Triggers.SlowRequests = &web.SlowRequestsBasedTrigger{
			TimeTaken:    pointer.To(triggers.SlowRequests[0].TimeTaken),
			TimeInterval: pointer.To(triggers.SlowRequests[0].Interval),
			Count:        pointer.To(int32(triggers.SlowRequests[0].Count)),
		}
		if triggers.SlowRequests[0].Path != "" {
			result.Triggers.SlowRequests.Path = pointer.To(triggers.SlowRequests[0].Path)
		}
	}

	if len(triggers.StatusCodes) > 0 {
		statusCodeTriggers := make([]web.StatusCodesBasedTrigger, 0)
		statusCodeRangeTriggers := make([]web.StatusCodesRangeBasedTrigger, 0)
		for _, s := range triggers.StatusCodes {
			statusCodeTrigger := web.StatusCodesBasedTrigger{}
			statusCodeRangeTrigger := web.StatusCodesRangeBasedTrigger{}
			parts := strings.Split(s.StatusCodeRange, "-")
			if len(parts) == 2 {
				statusCodeRangeTrigger.StatusCodes = pointer.To(s.StatusCodeRange)
				statusCodeRangeTrigger.Count = pointer.To(int32(s.Count))
				statusCodeRangeTrigger.TimeInterval = pointer.To(s.Interval)
				if s.Path != "" {
					statusCodeRangeTrigger.Path = pointer.To(s.Path)
				}
				statusCodeRangeTriggers = append(statusCodeRangeTriggers, statusCodeRangeTrigger)
			} else {
				statusCode, err := strconv.Atoi(s.StatusCodeRange)
				if err == nil {
					statusCodeTrigger.Status = pointer.To(int32(statusCode))
				}
				statusCodeTrigger.Count = pointer.To(int32(s.Count))
				statusCodeTrigger.TimeInterval = pointer.To(s.Interval)
				if s.Path != "" {
					statusCodeTrigger.Path = pointer.To(s.Path)
				}
				statusCodeTriggers = append(statusCodeTriggers, statusCodeTrigger)
			}
		}
		result.Triggers.StatusCodes = &statusCodeTriggers
		result.Triggers.StatusCodesRange = &statusCodeRangeTriggers
	}

	action := autoHeal.Actions[0]
	result.Actions.ActionType = web.AutoHealActionType(action.ActionType)
	result.Actions.MinProcessExecutionTime = pointer.To(action.MinimumProcessTime)

	return result
}

func flattenAutoHealSettingsLinux(autoHealRules *web.AutoHealRules) []AutoHealSettingLinux {
	if autoHealRules == nil {
		return nil
	}

	result := AutoHealSettingLinux{}

	// Triggers
	if autoHealRules.Triggers != nil {
		resultTrigger := AutoHealTriggerLinux{}
		triggers := *autoHealRules.Triggers
		if triggers.Requests != nil {
			count := 0
			if triggers.Requests.Count != nil {
				count = int(*triggers.Requests.Count)
			}
			resultTrigger.Requests = []AutoHealRequestTrigger{{
				Count:    count,
				Interval: pointer.From(triggers.Requests.TimeInterval),
			}}
		}

		statusCodeTriggers := make([]AutoHealStatusCodeTrigger, 0)
		if triggers.StatusCodes != nil {
			for _, s := range *triggers.StatusCodes {
				t := AutoHealStatusCodeTrigger{
					Interval: pointer.From(s.TimeInterval),
					Path:     pointer.From(s.Path),
				}

				if s.Status != nil {
					t.StatusCodeRange = strconv.Itoa(int(*s.Status))
				}

				if s.Count != nil {
					t.Count = int(*s.Count)
				}

				if s.SubStatus != nil {
					t.SubStatus = int(*s.SubStatus)
				}
				statusCodeTriggers = append(statusCodeTriggers, t)
			}
		}
		if triggers.StatusCodesRange != nil {
			for _, s := range *triggers.StatusCodesRange {
				t := AutoHealStatusCodeTrigger{
					Interval: pointer.From(s.TimeInterval),
					Path:     pointer.From(s.Path),
				}
				if s.Count != nil {
					t.Count = int(*s.Count)
				}

				if s.StatusCodes != nil {
					t.StatusCodeRange = *s.StatusCodes
				}
				statusCodeTriggers = append(statusCodeTriggers, t)
			}
		}
		resultTrigger.StatusCodes = statusCodeTriggers

		slowRequestTriggers := make([]AutoHealSlowRequest, 0)
		if triggers.SlowRequests != nil {
			slowRequestTriggers = append(slowRequestTriggers, AutoHealSlowRequest{
				TimeTaken: pointer.From(triggers.SlowRequests.TimeTaken),
				Interval:  pointer.From(triggers.SlowRequests.TimeInterval),
				Count:     int(pointer.From(triggers.SlowRequests.Count)),
				Path:      pointer.From(triggers.SlowRequests.Path),
			})
		}
		resultTrigger.SlowRequests = slowRequestTriggers
		result.Triggers = []AutoHealTriggerLinux{resultTrigger}
	}

	// Actions
	if autoHealRules.Actions != nil {
		actions := *autoHealRules.Actions

		result.Actions = []AutoHealActionLinux{{
			ActionType:         string(actions.ActionType),
			MinimumProcessTime: pointer.From(actions.MinProcessExecutionTime),
		}}
	}

	if result.Triggers != nil || result.Actions != nil {
		return []AutoHealSettingLinux{result}
	}

	return nil
}

func DisabledLogsConfig() *web.SiteLogsConfig {
	return &web.SiteLogsConfig{
		SiteLogsConfigProperties: &web.SiteLogsConfigProperties{
			DetailedErrorMessages: &web.EnabledConfig{
				Enabled: pointer.To(false),
			},
			FailedRequestsTracing: &web.EnabledConfig{
				Enabled: pointer.To(false),
			},
			ApplicationLogs: &web.ApplicationLogsConfig{
				FileSystem: &web.FileSystemApplicationLogsConfig{
					Level: web.LogLevelOff,
				},
				AzureBlobStorage: &web.AzureBlobStorageApplicationLogsConfig{
					Level: web.LogLevelOff,
				},
			},
			HTTPLogs: &web.HTTPLogsConfig{
				FileSystem: &web.FileSystemHTTPLogsConfig{
					Enabled: pointer.To(false),
				},
				AzureBlobStorage: &web.AzureBlobStorageHTTPLogsConfig{
					Enabled: pointer.To(false),
				},
			},
		},
	}
}

func isFreeOrSharedServicePlan(inputSKU string) bool {
	result := false
	for _, sku := range freeSkus {
		if inputSKU == sku {
			result = true
		}
	}
	for _, sku := range sharedSkus {
		if inputSKU == sku {
			result = true
		}
	}
	return result
}
