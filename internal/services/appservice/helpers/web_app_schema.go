package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	apimValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SiteConfigWindows struct {
	AlwaysOn                 bool                      `tfschema:"always_on"`
	ApiManagementConfigId    string                    `tfschema:"api_management_api_id"`
	ApiDefinition            string                    `tfschema:"api_definition_url"`
	AppCommandLine           string                    `tfschema:"app_command_line"`
	AutoHeal                 bool                      `tfschema:"auto_heal"`
	AutoHealSettings         []AutoHealSettingWindows  `tfschema:"auto_heal_setting"`
	UseManagedIdentityACR    bool                      `tfschema:"container_registry_use_managed_identity"`
	ContainerRegistryUserMSI string                    `tfschema:"container_registry_managed_identity_client_id"`
	DefaultDocuments         []string                  `tfschema:"default_documents"`
	Http2Enabled             bool                      `tfschema:"http2_enabled"`
	IpRestriction            []IpRestriction           `tfschema:"ip_restriction"`
	ScmUseMainIpRestriction  bool                      `tfschema:"scm_use_main_ip_restriction"`
	ScmIpRestriction         []IpRestriction           `tfschema:"scm_ip_restriction"`
	LoadBalancing            string                    `tfschema:"load_balancing_mode"`
	LocalMysql               bool                      `tfschema:"local_mysql"`
	ManagedPipelineMode      string                    `tfschema:"managed_pipeline_mode"`
	RemoteDebugging          bool                      `tfschema:"remote_debugging"`
	RemoteDebuggingVersion   string                    `tfschema:"remote_debugging_version"`
	ScmType                  string                    `tfschema:"scm_type"`
	Use32BitWorker           bool                      `tfschema:"use_32_bit_worker"`
	WebSockets               bool                      `tfschema:"websockets_enabled"`
	FtpsState                string                    `tfschema:"ftps_state"`
	HealthCheckPath          string                    `tfschema:"health_check_path"`
	HealthCheckEvictionTime  int                       `tfschema:"health_check_eviction_time_in_min"`
	NumberOfWorkers          int                       `tfschema:"number_of_workers"`
	ApplicationStack         []ApplicationStackWindows `tfschema:"application_stack"`
	VirtualApplications      []VirtualApplication      `tfschema:"virtual_application"`
	MinTlsVersion            string                    `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion         string                    `tfschema:"scm_minimum_tls_version"`
	AutoSwapSlotName         string                    `tfschema:"auto_swap_slot_name"`
	Cors                     []CorsSetting             `tfschema:"cors"`
	DetailedErrorLogging     bool                      `tfschema:"detailed_error_logging"`
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
	AutoHeal                bool                    `tfschema:"auto_heal"`
	AutoHealSettings        []AutoHealSettingLinux  `tfschema:"auto_heal_setting"`
	UseManagedIdentityACR   bool                    `tfschema:"container_registry_use_managed_identity"`
	ContainerRegistryMSI    string                  `tfschema:"container_registry_managed_identity_client_id"`
	DefaultDocuments        []string                `tfschema:"default_documents"`
	Http2Enabled            bool                    `tfschema:"http2_enabled"`
	IpRestriction           []IpRestriction         `tfschema:"ip_restriction"`
	ScmUseMainIpRestriction bool                    `tfschema:"scm_use_main_ip_restriction"`
	ScmIpRestriction        []IpRestriction         `tfschema:"scm_ip_restriction"`
	LoadBalancing           string                  `tfschema:"load_balancing_mode"`
	LocalMysql              bool                    `tfschema:"local_mysql"`
	ManagedPipelineMode     string                  `tfschema:"managed_pipeline_mode"`
	RemoteDebugging         bool                    `tfschema:"remote_debugging"`
	RemoteDebuggingVersion  string                  `tfschema:"remote_debugging_version"`
	ScmType                 string                  `tfschema:"scm_type"`
	Use32BitWorker          bool                    `tfschema:"use_32_bit_worker"`
	WebSockets              bool                    `tfschema:"websockets_enabled"`
	FtpsState               string                  `tfschema:"ftps_state"`
	HealthCheckPath         string                  `tfschema:"health_check_path"`
	HealthCheckEvictionTime int                     `tfschema:"health_check_eviction_time_in_min"`
	NumberOfWorkers         int                     `tfschema:"number_of_workers"`
	ApplicationStack        []ApplicationStackLinux `tfschema:"application_stack"`
	MinTlsVersion           string                  `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion        string                  `tfschema:"scm_minimum_tls_version"`
	AutoSwapSlotName        string                  `tfschema:"auto_swap_slot_name"`
	Cors                    []CorsSetting           `tfschema:"cors"`
	DetailedErrorLogging    bool                    `tfschema:"detailed_error_logging"`
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

				"auto_heal": {
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

				"local_mysql": {
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

				"remote_debugging": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"remote_debugging_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
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
					Computed: true, // Variable default value depending on several factors, such as plan type.
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

				"number_of_workers": {
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

				"auto_swap_slot_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					// TODO - Add slot name validation here?
				},

				"detailed_error_logging": {
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

				"auto_heal": {
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

				"local_mysql": {
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

				"remote_debugging": {
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

				"number_of_workers": {
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

				"auto_swap_slot_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"detailed_error_logging": {
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

				"auto_heal": {
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

				"local_mysql": {
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

				"remote_debugging": {
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

				"number_of_workers": {
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

				"auto_swap_slot_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					// TODO - Add slot name validation here?
				},

				"vnet_route_all_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should all outbound traffic to have Virtual Network Security Groups and User Defined Routes applied? Defaults to `false`.",
				},

				"detailed_error_logging": {
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

				"auto_heal": {
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

				"local_mysql": {
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

				"remote_debugging": {
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

				"number_of_workers": {
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

				"auto_swap_slot_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"detailed_error_logging": {
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

type ApplicationStackWindows struct {
	NetFrameworkVersion     string `tfschema:"dotnet_version"`
	PhpVersion              string `tfschema:"php_version"`
	JavaVersion             string `tfschema:"java_version"`
	PythonVersion           string `tfschema:"python_version"`
	NodeVersion             string `tfschema:"node_version"`
	JavaContainer           string `tfschema:"java_container"`
	JavaContainerVersion    string `tfschema:"java_container_version"`
	DockerContainerName     string `tfschema:"docker_container_name"`
	DockerContainerRegistry string `tfschema:"docker_container_registry"`
	DockerContainerTag      string `tfschema:"docker_container_tag"`
	CurrentStack            string `tfschema:"current_stack"`
}

// Version information for the below validations was taken in part from - https://github.com/Azure/app-service-linux-docs/tree/master/Runtime_Support
func windowsApplicationStackSchema() *pluginsdk.Schema {
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
					ValidateFunc: validation.StringInSlice([]string{
						"v2.0",
						"v3.0",
						"v4.0",
						"v5.0",
					}, false),
				},

				"php_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"5.6",
						"7.3",
						"7.4",
					}, false),
				},

				"python_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"2.7",
						"3.4.0",
					}, false),
				},

				"node_version": { // Discarded by service if JavaVersion is specified
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"10.1",   // Linux Only?
						"10.6",   // Linux Only?
						"10.10",  // Linux Only?
						"10.14",  // Linux Only?
						"10-LTS", // Linux Only?
						"12-LTS",
						"14-LTS",
					}, false),
					ConflictsWith: []string{
						"site_config.0.application_stack.0.java_version",
					},
				},

				"java_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					// Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"1.7",
						"1.8",
						"11",
					}, false),
				},

				"java_container": {Type: pluginsdk.TypeString,
					Optional: true,
					// Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"JAVA",
						"JETTY",
						"TOMCAT",
					}, false),
					RequiredWith: []string{
						"site_config.0.application_stack.0.java_container_version",
					},
				},

				"java_container_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					// Computed: true,
					RequiredWith: []string{
						"site_config.0.application_stack.0.java_container",
					},
				},

				"docker_container_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					RequiredWith: []string{
						"site_config.0.application_stack.0.docker_container_registry",
						"site_config.0.application_stack.0.docker_container_tag",
					},
				},

				"docker_container_registry": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					RequiredWith: []string{
						"site_config.0.application_stack.0.docker_container_name",
						"site_config.0.application_stack.0.docker_container_tag",
					},
				},

				"docker_container_tag": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					RequiredWith: []string{
						"site_config.0.application_stack.0.docker_container_name",
						"site_config.0.application_stack.0.docker_container_registry",
					},
				},

				"current_stack": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"dotnet",
						"node",
						"python",
						"php",
						"java",
					}, false),
				},
			},
		},
	}
}

func windowsApplicationStackSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"dotnet_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"php_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"python_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"node_version": { // Discarded by service if JavaVersion is specified
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_container": {Type: pluginsdk.TypeString,
					Computed: true,
				},

				"java_container_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_container_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_container_registry": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_container_tag": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"current_stack": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

type ApplicationStackLinux struct {
	NetFrameworkVersion string `tfschema:"dotnet_version"`
	PhpVersion          string `tfschema:"php_version"`
	PythonVersion       string `tfschema:"python_version"` // Linux Only?
	NodeVersion         string `tfschema:"node_version"`
	JavaVersion         string `tfschema:"java_version"`
	JavaServer          string `tfschema:"java_server"`
	JavaServerVersion   string `tfschema:"java_server_version"`
	DockerImageTag      string `tfschema:"docker_image_tag"`
	DockerImage         string `tfschema:"docker_image"`
	RubyVersion         string `tfschema:"ruby_version"`
}

// version information in the validation here was taken mostly from - `az webapp list-runtimes --linux`
func linuxApplicationStackSchema() *pluginsdk.Schema {
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
					ValidateFunc: validation.StringInSlice([]string{
						"2.1",
						"3.1",
						"5.0",
					}, false),
					ConflictsWith: []string{
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.ruby_version",
						"site_config.0.application_stack.0.java_version",
					},
				},

				"php_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"5.6", // TODO - Remove? 5.6 is available, but deprecated in the service
						"7.2", // TODO - Remove? 7.2 is available, but deprecated in the service
						"7.3",
						"7.4",
					}, false),
					ConflictsWith: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.ruby_version",
						"site_config.0.application_stack.0.java_version",
					},
				},

				"python_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"2.7", // TODO - Remove? 2.7 is available, but deprecated in the service
						"3.6",
						"3.7",
						"3.8",
					}, false),
					ConflictsWith: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.ruby_version",
						"site_config.0.application_stack.0.java_version",
					},
				},

				"node_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"10.1",   // TODO - Remove?  Deprecated
						"10.6",   // TODO - Remove?  Deprecated
						"10.14",  // TODO - Remove?  Deprecated
						"10-lts", // TODO - Remove?  Deprecated
						"12-lts",
						"14-lts",
					}, false),
					ConflictsWith: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.ruby_version",
						"site_config.0.application_stack.0.java_version",
					},
				},

				"ruby_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"2.5",
						"2.6",
					}, false),
					ConflictsWith: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.java_version",
					},
				},

				"java_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty, // There a significant number of variables here, and the versions are not uniformly formatted.
					// TODO - Needs notes in the docs for this to help users navigate the inconsistencies in the service. e.g. jre8 va java8 etc
					ConflictsWith: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.ruby_version",
					},
				},

				"java_server": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"JAVA",
						"TOMCAT",
						"JBOSSEAP",
					}, false),
				},

				"java_server_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"docker_image": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					RequiredWith: []string{
						"site_config.0.application_stack.0.docker_image_tag",
					},
				},

				"docker_image_tag": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					RequiredWith: []string{
						"site_config.0.application_stack.0.docker_image",
					},
				},
			},
		},
	}
}

func linuxApplicationStackSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"dotnet_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"php_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"python_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"node_version": { // Discarded by service if JavaVersion is specified
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"ruby_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_server": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_server_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_image": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_image_tag": {
					Type:     pluginsdk.TypeString,
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
			"site_config.0.auto_heal",
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
			"site_config.0.auto_heal",
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
					//ValidateFunc: // TODO - Time in hh:mm:ss, because why not...
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
					//ValidateFunc: // TODO - Time in hh:mm:ss, because why not...
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
								//ValidateFunc: validation.IsRFC3339Time, // TODO should be hh:mm:ss - This is too loose, need to improve
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
								//ValidateFunc: validation.IsRFC3339Time,
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
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"time_taken": {
								Type:     pluginsdk.TypeString,
								Required: true,
								//ValidateFunc: validation.IsRFC3339Time,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Required: true,
								//ValidateFunc: validation.IsRFC3339Time,
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
								//ValidateFunc: validation.IsRFC3339Time, // TODO should be hh:mm:ss - This is too loose, need to improve?
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
								//ValidateFunc: validation.IsRFC3339Time,
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
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"time_taken": {
								Type:     pluginsdk.TypeString,
								Required: true,
								//ValidateFunc: validation.IsRFC3339Time,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Required: true,
								//ValidateFunc: validation.IsRFC3339Time,
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
				},

				"storage_account_url": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.IsURLWithHTTPS,
				},

				"enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
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
							},

							"frequency_unit": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ValidateFunc: validation.StringInSlice([]string{
									"Day",
									"Hour",
								}, false),
							},

							"keep_at_least_one_backup": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								Default:  false,
							},

							"retention_period_days": {
								Type:         pluginsdk.TypeInt,
								Optional:     true,
								Default:      30,
								ValidateFunc: validation.IntBetween(0, 9999999),
							},

							"start_time": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								Computed: true,
								//DiffSuppressFunc: suppress.RFC3339Time,
								//ValidateFunc:     validation.IsRFC3339Time,
							},

							"last_execution_time": {
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

func BackupSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"storage_account_url": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},

				"enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"schedule": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"frequency_interval": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"frequency_unit": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"keep_at_least_one_backup": {
								Type:     pluginsdk.TypeBool,
								Computed: true,
							},

							"retention_period_days": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"start_time": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"last_execution_time": {
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
					}, true),
				},

				"value": {
					Type:      pluginsdk.TypeString,
					Required:  true,
					Sensitive: true,
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
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"value": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
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

func ExpandSiteConfigWindows(siteConfig []SiteConfigWindows, existing *web.SiteConfig, metadata sdk.ResourceMetaData) (*web.SiteConfig, *string, error) {
	if len(siteConfig) == 0 {
		return nil, nil, nil
	}

	expanded := &web.SiteConfig{}
	if existing != nil {
		expanded = existing
	}

	currentStack := ""

	winSiteConfig := siteConfig[0]

	if metadata.ResourceData.HasChange("site_config.0.always_on") {
		expanded.AlwaysOn = utils.Bool(winSiteConfig.AlwaysOn)
	}

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.APIManagementConfig = &web.APIManagementConfig{
			ID: utils.String(winSiteConfig.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.APIDefinition = &web.APIDefinitionInfo{
			URL: utils.String(winSiteConfig.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = utils.String(winSiteConfig.AppCommandLine)
	}

	if metadata.ResourceData.HasChange("site_config.0.application_stack") {
		winAppStack := winSiteConfig.ApplicationStack[0]
		expanded.NetFrameworkVersion = utils.String(winAppStack.NetFrameworkVersion)
		expanded.PhpVersion = utils.String(winAppStack.PhpVersion)
		expanded.NodeVersion = utils.String(winAppStack.NodeVersion)
		expanded.PythonVersion = utils.String(winAppStack.PythonVersion)
		expanded.JavaVersion = utils.String(winAppStack.JavaVersion)
		expanded.JavaContainer = utils.String(winAppStack.JavaContainer)
		expanded.JavaContainerVersion = utils.String(winAppStack.JavaContainerVersion)
		if winAppStack.DockerContainerName != "" {
			expanded.WindowsFxVersion = utils.String(fmt.Sprintf("DOCKER|%s/%s:%s", winAppStack.DockerContainerRegistry, winAppStack.DockerContainerName, winAppStack.DockerContainerTag))
		}
		currentStack = winAppStack.CurrentStack
	}

	if metadata.ResourceData.HasChange("site_config.0.virtual_application") {
		expanded.VirtualApplications = expandVirtualApplicationsForUpdate(winSiteConfig.VirtualApplications)
	} else {
		expanded.VirtualApplications = expandVirtualApplications(winSiteConfig.VirtualApplications)
	}

	if metadata.ResourceData.HasChange("site_config.0.container_registry_use_managed_identity") {
		expanded.AcrUseManagedIdentityCreds = utils.Bool(winSiteConfig.UseManagedIdentityACR)
	}

	if metadata.ResourceData.HasChange("site_config.0.container_registry_managed_identity_client_id") {
		expanded.AcrUserManagedIdentityID = utils.String(winSiteConfig.ContainerRegistryUserMSI)
	}

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = &winSiteConfig.DefaultDocuments
	}

	if metadata.ResourceData.HasChange("site_config.0.http2_enabled") {
		expanded.HTTP20Enabled = utils.Bool(winSiteConfig.Http2Enabled)
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := ExpandIpRestrictions(winSiteConfig.IpRestriction)
		if err != nil {
			return nil, nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_use_main_ip_restriction") {
		expanded.ScmIPSecurityRestrictionsUseMain = utils.Bool(winSiteConfig.ScmUseMainIpRestriction)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := ExpandIpRestrictions(winSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.local_mysql") {
		expanded.LocalMySQLEnabled = utils.Bool(winSiteConfig.LocalMysql)
	}

	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = web.SiteLoadBalancing(winSiteConfig.LoadBalancing)
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = web.ManagedPipelineMode(winSiteConfig.ManagedPipelineMode)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging") {
		expanded.RemoteDebuggingEnabled = utils.Bool(winSiteConfig.RemoteDebugging)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = utils.String(winSiteConfig.RemoteDebuggingVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.use_32_bit_worker") {
		expanded.Use32BitWorkerProcess = utils.Bool(winSiteConfig.Use32BitWorker)
	}

	if metadata.ResourceData.HasChange("site_config.0.websockets_enabled") {
		expanded.WebSocketsEnabled = utils.Bool(winSiteConfig.WebSockets)
	}

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		expanded.FtpsState = web.FtpsState(winSiteConfig.FtpsState)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = utils.String(winSiteConfig.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.number_of_workers") {
		expanded.NumberOfWorkers = utils.Int32(int32(winSiteConfig.NumberOfWorkers))
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTLSVersion = web.SupportedTLSVersions(winSiteConfig.MinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTLSVersion = web.SupportedTLSVersions(winSiteConfig.ScmMinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_swap_slot_name") {
		expanded.AutoSwapSlotName = utils.String(winSiteConfig.AutoSwapSlotName)
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

	if metadata.ResourceData.HasChange("site_config.0.auto_heal") {
		expanded.AutoHealEnabled = utils.Bool(winSiteConfig.AutoHeal)
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_heal_setting") {
		expanded.AutoHealRules = expandAutoHealSettingsWindows(winSiteConfig.AutoHealSettings)
	}

	if metadata.ResourceData.HasChange("site_config.0.vnet_route_all_enabled") {
		expanded.VnetRouteAllEnabled = utils.Bool(winSiteConfig.VnetRouteAllEnabled)
	}

	return expanded, &currentStack, nil
}

func ExpandSiteConfigLinux(siteConfig []SiteConfigLinux, existing *web.SiteConfig, metadata sdk.ResourceMetaData) (*web.SiteConfig, error) {
	if len(siteConfig) == 0 {
		return nil, nil
	}
	expanded := &web.SiteConfig{}
	if existing != nil {
		expanded = existing
	}

	linuxSiteConfig := siteConfig[0]

	if metadata.ResourceData.HasChange("site_config.0.always_on") {
		expanded.AlwaysOn = utils.Bool(linuxSiteConfig.AlwaysOn)
	}

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.APIManagementConfig = &web.APIManagementConfig{
			ID: utils.String(linuxSiteConfig.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.APIDefinition = &web.APIDefinitionInfo{
			URL: utils.String(linuxSiteConfig.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = utils.String(linuxSiteConfig.AppCommandLine)
	}

	if metadata.ResourceData.HasChange("site_config.0.application_stack") {
		linuxAppStack := linuxSiteConfig.ApplicationStack[0]
		if linuxAppStack.NetFrameworkVersion != "" {
			expanded.LinuxFxVersion = utils.String(fmt.Sprintf("DOTNETCORE|%s", linuxAppStack.NetFrameworkVersion))
		}

		if linuxAppStack.PhpVersion != "" {
			expanded.LinuxFxVersion = utils.String(fmt.Sprintf("PHP|%s", linuxAppStack.PhpVersion))
		}

		if linuxAppStack.NodeVersion != "" {
			expanded.LinuxFxVersion = utils.String(fmt.Sprintf("NODE|%s", linuxAppStack.NodeVersion))
		}

		if linuxAppStack.PythonVersion != "" {
			expanded.LinuxFxVersion = utils.String(fmt.Sprintf("PYTHON|%s", linuxAppStack.PythonVersion))
		}

		if linuxAppStack.JavaServer != "" {
			// (@jackofallops) - Java has some special cases for Java SE when using specific versions of the runtime, resulting in this string
			// being formatted in the form: `JAVA|u242` instead of the standard pattern of `JAVA|u242-java8` for example. This applies to jre8 and java11.
			if linuxAppStack.JavaServer == "JAVA" && linuxAppStack.JavaServerVersion == "" {
				expanded.LinuxFxVersion = utils.String(fmt.Sprintf("%s|%s", linuxAppStack.JavaServer, linuxAppStack.JavaVersion))
			} else {
				expanded.LinuxFxVersion = utils.String(fmt.Sprintf("%s|%s-%s", linuxAppStack.JavaServer, linuxAppStack.JavaServerVersion, linuxAppStack.JavaVersion))
			}
		}

		if linuxAppStack.DockerImage != "" {
			expanded.LinuxFxVersion = utils.String(fmt.Sprintf("DOCKER|%s:%s", linuxAppStack.DockerImage, linuxAppStack.DockerImageTag))
		}
	}

	expanded.AcrUseManagedIdentityCreds = utils.Bool(linuxSiteConfig.UseManagedIdentityACR)

	if metadata.ResourceData.HasChange("site_config.0.container_registry_managed_identity_client_id") {
		expanded.AcrUserManagedIdentityID = utils.String(linuxSiteConfig.ContainerRegistryMSI)
	}

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = &linuxSiteConfig.DefaultDocuments
	}

	if metadata.ResourceData.HasChange("site_config.0.http2_enabled") {
		expanded.HTTP20Enabled = utils.Bool(linuxSiteConfig.Http2Enabled)
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := ExpandIpRestrictions(linuxSiteConfig.IpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	expanded.ScmIPSecurityRestrictionsUseMain = utils.Bool(linuxSiteConfig.ScmUseMainIpRestriction)

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := ExpandIpRestrictions(linuxSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.local_mysql") {
		expanded.LocalMySQLEnabled = utils.Bool(linuxSiteConfig.LocalMysql)
	}

	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = web.SiteLoadBalancing(linuxSiteConfig.LoadBalancing)
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = web.ManagedPipelineMode(linuxSiteConfig.ManagedPipelineMode)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging") {
		expanded.RemoteDebuggingEnabled = utils.Bool(linuxSiteConfig.RemoteDebugging)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = utils.String(linuxSiteConfig.RemoteDebuggingVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.use_32_bit_worker") {
		expanded.Use32BitWorkerProcess = utils.Bool(linuxSiteConfig.Use32BitWorker)
	}

	if metadata.ResourceData.HasChange("site_config.0.websockets_enabled") {
		expanded.WebSocketsEnabled = utils.Bool(linuxSiteConfig.WebSockets)
	}

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		expanded.FtpsState = web.FtpsState(linuxSiteConfig.FtpsState)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = utils.String(linuxSiteConfig.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.number_of_workers") {
		expanded.NumberOfWorkers = utils.Int32(int32(linuxSiteConfig.NumberOfWorkers))
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTLSVersion = web.SupportedTLSVersions(linuxSiteConfig.MinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTLSVersion = web.SupportedTLSVersions(linuxSiteConfig.ScmMinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_swap_slot_name") {
		expanded.AutoSwapSlotName = utils.String(linuxSiteConfig.AutoSwapSlotName)
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

	if metadata.ResourceData.HasChange("site_config.0.auto_heal") {
		expanded.AutoHealEnabled = utils.Bool(linuxSiteConfig.AutoHeal)
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_heal_setting") {
		expanded.AutoHealRules = expandAutoHealSettingsLinux(linuxSiteConfig.AutoHealSettings)
	}

	if metadata.ResourceData.HasChange("site_config.0.vnet_route_all_enabled") {
		expanded.VnetRouteAllEnabled = utils.Bool(linuxSiteConfig.VnetRouteAllEnabled)
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
				SasURL:          utils.String(appLogsBlobs.SasUrl),
				RetentionInDays: utils.Int32(int32(appLogsBlobs.RetentionInDays)),
			}
		}
	}

	if len(logsConfig.HttpLogs) == 1 {
		httpLogs := logsConfig.HttpLogs[0]
		result.HTTPLogs = &web.HTTPLogsConfig{}

		if len(httpLogs.FileSystems) == 1 {
			httpLogFileSystem := httpLogs.FileSystems[0]
			result.HTTPLogs.FileSystem = &web.FileSystemHTTPLogsConfig{
				Enabled:         utils.Bool(true),
				RetentionInMb:   utils.Int32(int32(httpLogFileSystem.RetentionMB)),
				RetentionInDays: utils.Int32(int32(httpLogFileSystem.RetentionDays)),
			}
		}

		if len(httpLogs.AzureBlobStorage) == 1 {
			httpLogsBlobStorage := httpLogs.AzureBlobStorage[0]
			result.HTTPLogs.AzureBlobStorage = &web.AzureBlobStorageHTTPLogsConfig{
				Enabled:         utils.Bool(httpLogsBlobStorage.SasUrl != ""),
				SasURL:          utils.String(httpLogsBlobStorage.SasUrl),
				RetentionInDays: utils.Int32(int32(httpLogsBlobStorage.RetentionInDays)),
			}
		}
	}

	result.DetailedErrorMessages = &web.EnabledConfig{
		Enabled: utils.Bool(logsConfig.DetailedErrorMessages),
	}

	result.FailedRequestsTracing = &web.EnabledConfig{
		Enabled: utils.Bool(logsConfig.FailedRequestTracing),
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
		Enabled:           utils.Bool(backupConfig.Enabled),
		BackupName:        utils.String(backupConfig.Name),
		StorageAccountURL: utils.String(backupConfig.StorageAccountUrl),
		BackupSchedule: &web.BackupSchedule{
			FrequencyInterval:     utils.Int32(int32(backupSchedule.FrequencyInterval)),
			FrequencyUnit:         web.FrequencyUnit(backupSchedule.FrequencyUnit),
			KeepAtLeastOneBackup:  utils.Bool(backupSchedule.KeepAtLeastOneBackup),
			RetentionPeriodInDays: utils.Int32(int32(backupSchedule.RetentionPeriodDays)),
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
			AccountName: utils.String(v.AccountName),
			ShareName:   utils.String(v.ShareName),
			AccessKey:   utils.String(v.AccessKey),
			MountPath:   utils.String(v.MountPath),
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
			Value: utils.String(v.Value),
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
			VirtualPath:    utils.String(v.VirtualPath),
			PhysicalPath:   utils.String(v.PhysicalPath),
			PreloadEnabled: utils.Bool(v.Preload),
		}
		if len(v.VirtualDirectories) > 0 {
			virtualDirs := make([]web.VirtualDirectory, 0)
			for _, d := range v.VirtualDirectories {
				virtualDirs = append(virtualDirs, web.VirtualDirectory{
					VirtualPath:  utils.String(d.VirtualPath),
					PhysicalPath: utils.String(d.PhysicalPath),
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
				VirtualPath:    utils.String("/"),
				PhysicalPath:   utils.String("site\\wwwroot"),
				PreloadEnabled: utils.Bool(true),
			},
		}
	}

	result := make([]web.VirtualApplication, 0)

	for _, v := range virtualApplicationConfig {
		virtualApp := web.VirtualApplication{
			VirtualPath:    utils.String(v.VirtualPath),
			PhysicalPath:   utils.String(v.PhysicalPath),
			PreloadEnabled: utils.Bool(v.Preload),
		}
		if len(v.VirtualDirectories) > 0 {
			virtualDirs := make([]web.VirtualDirectory, 0)
			for _, d := range v.VirtualDirectories {
				virtualDirs = append(virtualDirs, web.VirtualDirectory{
					VirtualPath:  utils.String(d.VirtualPath),
					PhysicalPath: utils.String(d.PhysicalPath),
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

				blobStorage.SasUrl = utils.NormalizeNilableString(appLogs.AzureBlobStorage.SasURL)

				blobStorage.RetentionInDays = int(utils.NormaliseNilableInt32(appLogs.AzureBlobStorage.RetentionInDays))

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
		AlwaysOn:                 utils.NormaliseNilableBool(appSiteConfig.AlwaysOn),
		AppCommandLine:           utils.NormalizeNilableString(appSiteConfig.AppCommandLine),
		AutoHeal:                 utils.NormaliseNilableBool(appSiteConfig.AutoHealEnabled),
		AutoHealSettings:         flattenAutoHealSettingsWindows(appSiteConfig.AutoHealRules),
		ContainerRegistryUserMSI: utils.NormalizeNilableString(appSiteConfig.AcrUserManagedIdentityID),
		DetailedErrorLogging:     utils.NormaliseNilableBool(appSiteConfig.DetailedErrorLoggingEnabled),
		FtpsState:                string(appSiteConfig.FtpsState),
		HealthCheckPath:          utils.NormalizeNilableString(appSiteConfig.HealthCheckPath),
		HealthCheckEvictionTime:  utils.NormaliseNilableInt(healthCheckCount),
		Http2Enabled:             utils.NormaliseNilableBool(appSiteConfig.HTTP20Enabled),
		IpRestriction:            FlattenIpRestrictions(appSiteConfig.IPSecurityRestrictions),
		LoadBalancing:            string(appSiteConfig.LoadBalancing),
		LocalMysql:               utils.NormaliseNilableBool(appSiteConfig.LocalMySQLEnabled),
		ManagedPipelineMode:      string(appSiteConfig.ManagedPipelineMode),
		MinTlsVersion:            string(appSiteConfig.MinTLSVersion),
		NumberOfWorkers:          int(utils.NormaliseNilableInt32(appSiteConfig.NumberOfWorkers)),
		RemoteDebugging:          utils.NormaliseNilableBool(appSiteConfig.RemoteDebuggingEnabled),
		RemoteDebuggingVersion:   strings.ToUpper(utils.NormalizeNilableString(appSiteConfig.RemoteDebuggingVersion)),
		ScmIpRestriction:         FlattenIpRestrictions(appSiteConfig.ScmIPSecurityRestrictions),
		ScmMinTlsVersion:         string(appSiteConfig.ScmMinTLSVersion),
		ScmType:                  string(appSiteConfig.ScmType),
		ScmUseMainIpRestriction:  utils.NormaliseNilableBool(appSiteConfig.ScmIPSecurityRestrictionsUseMain),
		Use32BitWorker:           utils.NormaliseNilableBool(appSiteConfig.Use32BitWorkerProcess),
		UseManagedIdentityACR:    utils.NormaliseNilableBool(appSiteConfig.AcrUseManagedIdentityCreds),
		VirtualApplications:      flattenVirtualApplications(appSiteConfig.VirtualApplications),
		WebSockets:               utils.NormaliseNilableBool(appSiteConfig.WebSocketsEnabled),
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
		siteConfig.NumberOfWorkers = int(*appSiteConfig.NumberOfWorkers)
	}

	var winAppStack ApplicationStackWindows
	winAppStack.NetFrameworkVersion = utils.NormalizeNilableString(appSiteConfig.NetFrameworkVersion)
	winAppStack.PhpVersion = utils.NormalizeNilableString(appSiteConfig.PhpVersion)
	winAppStack.NodeVersion = utils.NormalizeNilableString(appSiteConfig.NodeVersion)
	winAppStack.PythonVersion = utils.NormalizeNilableString(appSiteConfig.PythonVersion)
	winAppStack.JavaVersion = utils.NormalizeNilableString(appSiteConfig.JavaVersion)
	winAppStack.JavaContainer = utils.NormalizeNilableString(appSiteConfig.JavaContainer)
	winAppStack.JavaContainerVersion = utils.NormalizeNilableString(appSiteConfig.JavaContainerVersion)

	siteConfig.WindowsFxVersion = utils.NormalizeNilableString(appSiteConfig.WindowsFxVersion)
	if siteConfig.WindowsFxVersion != "" {
		// Decode the string to docker values
		parts := strings.Split(strings.TrimPrefix(siteConfig.WindowsFxVersion, "DOCKER|"), ":")
		winAppStack.DockerContainerTag = parts[1]
		path := strings.Split(parts[0], "/")
		winAppStack.DockerContainerRegistry = path[0]
		winAppStack.DockerContainerName = strings.TrimPrefix(parts[0], fmt.Sprintf("%s/", path[0]))
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
			siteConfig.Cors = []CorsSetting{cors}
		}
	}

	return []SiteConfigWindows{siteConfig}
}

func FlattenSiteConfigLinux(appSiteConfig *web.SiteConfig, healthCheckCount *int) []SiteConfigLinux {
	if appSiteConfig == nil {
		return nil
	}

	siteConfig := SiteConfigLinux{
		AlwaysOn:                utils.NormaliseNilableBool(appSiteConfig.AlwaysOn),
		AppCommandLine:          utils.NormalizeNilableString(appSiteConfig.AppCommandLine),
		AutoHeal:                utils.NormaliseNilableBool(appSiteConfig.AutoHealEnabled),
		AutoHealSettings:        flattenAutoHealSettingsLinux(appSiteConfig.AutoHealRules),
		ContainerRegistryMSI:    utils.NormalizeNilableString(appSiteConfig.AcrUserManagedIdentityID),
		DetailedErrorLogging:    utils.NormaliseNilableBool(appSiteConfig.DetailedErrorLoggingEnabled),
		Http2Enabled:            utils.NormaliseNilableBool(appSiteConfig.HTTP20Enabled),
		IpRestriction:           FlattenIpRestrictions(appSiteConfig.IPSecurityRestrictions),
		ManagedPipelineMode:     string(appSiteConfig.ManagedPipelineMode),
		ScmType:                 string(appSiteConfig.ScmType),
		FtpsState:               string(appSiteConfig.FtpsState),
		HealthCheckPath:         utils.NormalizeNilableString(appSiteConfig.HealthCheckPath),
		HealthCheckEvictionTime: utils.NormaliseNilableInt(healthCheckCount),
		LoadBalancing:           string(appSiteConfig.LoadBalancing),
		LocalMysql:              utils.NormaliseNilableBool(appSiteConfig.LocalMySQLEnabled),
		MinTlsVersion:           string(appSiteConfig.MinTLSVersion),
		NumberOfWorkers:         int(utils.NormaliseNilableInt32(appSiteConfig.NumberOfWorkers)),
		RemoteDebugging:         utils.NormaliseNilableBool(appSiteConfig.RemoteDebuggingEnabled),
		RemoteDebuggingVersion:  strings.ToUpper(utils.NormalizeNilableString(appSiteConfig.RemoteDebuggingVersion)),
		ScmIpRestriction:        FlattenIpRestrictions(appSiteConfig.ScmIPSecurityRestrictions),
		ScmMinTlsVersion:        string(appSiteConfig.ScmMinTLSVersion),
		ScmUseMainIpRestriction: utils.NormaliseNilableBool(appSiteConfig.ScmIPSecurityRestrictionsUseMain),
		Use32BitWorker:          utils.NormaliseNilableBool(appSiteConfig.Use32BitWorkerProcess),
		UseManagedIdentityACR:   utils.NormaliseNilableBool(appSiteConfig.AcrUseManagedIdentityCreds),
		WebSockets:              utils.NormaliseNilableBool(appSiteConfig.WebSocketsEnabled),
		VnetRouteAllEnabled:     utils.NormaliseNilableBool(appSiteConfig.VnetRouteAllEnabled),
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
			siteConfig.Cors = []CorsSetting{cors}
		}
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

func ExpandAppSettings(settings map[string]string) *web.StringDictionary {
	appSettings := make(map[string]*string)
	for k, v := range settings {
		appSettings[k] = utils.String(v)
	}

	return &web.StringDictionary{
		Properties: appSettings,
	}
}

func FlattenAppSettings(input web.StringDictionary) (map[string]string, *int) {
	maxPingFailures := "WEBSITE_HEALTHCHECK_MAXPINGFAILURE"
	unmanagedSettings := []string{
		"DIAGNOSTICS_AZUREBLOBCONTAINERSASURL",
		"DIAGNOSTICS_AZUREBLOBRETENTIONINDAYS",
		"WEBSITE_HTTPLOGGING_CONTAINER_URL",
		"WEBSITE_HTTPLOGGING_RETENTION_DAYS",
		maxPingFailures,
	}

	var healthCheckCount *int
	appSettings := FlattenWebStringDictionary(input)
	if v, ok := appSettings[maxPingFailures]; ok {
		h, _ := strconv.Atoi(v)
		healthCheckCount = &h
	}

	// Remove the settings the service adds for legacy reasons when logging settings are specified.
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
			VirtualPath:  utils.NormalizeNilableString(v.VirtualPath),
			PhysicalPath: utils.NormalizeNilableString(v.PhysicalPath),
		}
		if preload := v.PreloadEnabled; preload != nil {
			virtualApp.Preload = *preload
		}
		if v.VirtualDirectories != nil && len(*v.VirtualDirectories) > 0 {
			virtualDirs := make([]VirtualDirectory, 0)
			for _, d := range *v.VirtualDirectories {
				virtualDir := VirtualDirectory{
					VirtualPath:  utils.NormalizeNilableString(d.VirtualPath),
					PhysicalPath: utils.NormalizeNilableString(d.PhysicalPath),
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
			Count:        utils.Int32(int32(triggers.Requests[0].Count)),
			TimeInterval: utils.String(triggers.Requests[0].Interval),
		}
	}

	if triggers.PrivateMemoryKB != 0 {
		result.Triggers.PrivateBytesInKB = utils.Int32(int32(triggers.PrivateMemoryKB))
	}

	if len(triggers.StatusCodes) > 0 {
		statusCodeTriggers := make([]web.StatusCodesBasedTrigger, 0)
		statusCodeRangeTriggers := make([]web.StatusCodesRangeBasedTrigger, 0)
		for _, s := range triggers.StatusCodes {
			statusCodeTrigger := web.StatusCodesBasedTrigger{}
			statusCodeRangeTrigger := web.StatusCodesRangeBasedTrigger{}
			parts := strings.Split(s.StatusCodeRange, "-")
			if len(parts) == 2 {
				statusCodeRangeTrigger.StatusCodes = utils.String(s.StatusCodeRange)
				statusCodeRangeTrigger.Count = utils.Int32(int32(s.Count))
				statusCodeRangeTrigger.TimeInterval = utils.String(s.Interval)
				if s.Path != "" {
					statusCodeRangeTrigger.Path = utils.String(s.Path)
				}
				statusCodeRangeTriggers = append(statusCodeRangeTriggers, statusCodeRangeTrigger)
			} else {
				statusCode, err := strconv.Atoi(s.StatusCodeRange)
				if err == nil {
					statusCodeTrigger.Status = utils.Int32(int32(statusCode))
				}
				statusCodeTrigger.Count = utils.Int32(int32(s.Count))
				statusCodeTrigger.TimeInterval = utils.String(s.Interval)
				if s.Path != "" {
					statusCodeTrigger.Path = utils.String(s.Path)
				}
				statusCodeTriggers = append(statusCodeTriggers, statusCodeTrigger)
			}
		}
		result.Triggers.StatusCodes = &statusCodeTriggers
		result.Triggers.StatusCodesRange = &statusCodeRangeTriggers
	}

	action := autoHeal.Actions[0]
	result.Actions.ActionType = web.AutoHealActionType(action.ActionType)
	result.Actions.MinProcessExecutionTime = utils.String(action.MinimumProcessTime)
	if len(action.CustomAction) != 0 {
		customAction := action.CustomAction[0]
		result.Actions.CustomAction = &web.AutoHealCustomAction{
			Exe:        utils.String(customAction.Executable),
			Parameters: utils.String(customAction.Parameters),
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
				Interval: utils.NormalizeNilableString(triggers.Requests.TimeInterval),
			}}
		}

		if privateBytes := triggers.PrivateBytesInKB; privateBytes != nil && *privateBytes != 0 {
			resultTrigger.PrivateMemoryKB = int(*triggers.PrivateBytesInKB)
		}

		statusCodeTriggers := make([]AutoHealStatusCodeTrigger, 0)
		if triggers.StatusCodes != nil {
			for _, s := range *triggers.StatusCodes {
				t := AutoHealStatusCodeTrigger{
					Interval: utils.NormalizeNilableString(s.TimeInterval),
					Path:     utils.NormalizeNilableString(s.Path),
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
					Interval: utils.NormalizeNilableString(s.TimeInterval),
					Path:     utils.NormalizeNilableString(s.Path),
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
				TimeTaken: utils.NormalizeNilableString(triggers.SlowRequests.TimeTaken),
				Interval:  utils.NormalizeNilableString(triggers.SlowRequests.TimeInterval),
				Count:     int(utils.NormaliseNilableInt32(triggers.SlowRequests.Count)),
				Path:      utils.NormalizeNilableString(triggers.SlowRequests.Path),
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
				Executable: utils.NormalizeNilableString(actions.CustomAction.Exe),
				Parameters: utils.NormalizeNilableString(actions.CustomAction.Parameters),
			})
		}

		resultActions := AutoHealActionWindows{
			ActionType:         string(actions.ActionType),
			CustomAction:       customActions,
			MinimumProcessTime: utils.NormalizeNilableString(actions.MinProcessExecutionTime),
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
			Count:        utils.Int32(int32(triggers.Requests[0].Count)),
			TimeInterval: utils.String(triggers.Requests[0].Interval),
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
				statusCodeRangeTrigger.StatusCodes = utils.String(s.StatusCodeRange)
				statusCodeRangeTrigger.Count = utils.Int32(int32(s.Count))
				statusCodeRangeTrigger.TimeInterval = utils.String(s.Interval)
				if s.Path != "" {
					statusCodeRangeTrigger.Path = utils.String(s.Path)
				}
				statusCodeRangeTriggers = append(statusCodeRangeTriggers, statusCodeRangeTrigger)
			} else {
				statusCode, err := strconv.Atoi(s.StatusCodeRange)
				if err == nil {
					statusCodeTrigger.Status = utils.Int32(int32(statusCode))
				}
				statusCodeTrigger.Count = utils.Int32(int32(s.Count))
				statusCodeTrigger.TimeInterval = utils.String(s.Interval)
				if s.Path != "" {
					statusCodeTrigger.Path = utils.String(s.Path)
				}
				statusCodeTriggers = append(statusCodeTriggers, statusCodeTrigger)
			}
		}
		result.Triggers.StatusCodes = &statusCodeTriggers
		result.Triggers.StatusCodesRange = &statusCodeRangeTriggers
	}

	action := autoHeal.Actions[0]
	result.Actions.ActionType = web.AutoHealActionType(action.ActionType)
	result.Actions.MinProcessExecutionTime = utils.String(action.MinimumProcessTime)

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
				Interval: utils.NormalizeNilableString(triggers.Requests.TimeInterval),
			}}
		}

		statusCodeTriggers := make([]AutoHealStatusCodeTrigger, 0)
		if triggers.StatusCodes != nil {
			for _, s := range *triggers.StatusCodes {
				t := AutoHealStatusCodeTrigger{
					Interval: utils.NormalizeNilableString(s.TimeInterval),
					Path:     utils.NormalizeNilableString(s.Path),
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
					Interval: utils.NormalizeNilableString(s.TimeInterval),
					Path:     utils.NormalizeNilableString(s.Path),
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
				TimeTaken: utils.NormalizeNilableString(triggers.SlowRequests.TimeTaken),
				Interval:  utils.NormalizeNilableString(triggers.SlowRequests.TimeInterval),
				Count:     int(utils.NormaliseNilableInt32(triggers.SlowRequests.Count)),
				Path:      utils.NormalizeNilableString(triggers.SlowRequests.Path),
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
			MinimumProcessTime: utils.NormalizeNilableString(actions.MinProcessExecutionTime),
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
				Enabled: utils.Bool(false),
			},
			FailedRequestsTracing: &web.EnabledConfig{
				Enabled: utils.Bool(false),
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
					Enabled: utils.Bool(false),
				},
				AzureBlobStorage: &web.AzureBlobStorageHTTPLogsConfig{
					Enabled: utils.Bool(false),
				},
			},
		},
	}
}
