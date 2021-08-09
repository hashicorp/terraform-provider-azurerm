package appservice

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-01-15/web"
	"github.com/Azure/go-autorest/autorest/date"
	apimValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	msiValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SiteConfigWindows struct {
	AlwaysOn                 bool                      `tfschema:"always_on"`
	ApiManagementConfigId    string                    `tfschema:"api_management_config_id"`
	AppCommandLine           string                    `tfschema:"app_command_line"`
	AutoHeal                 bool                      `tfschema:"auto_heal"`
	AutoHealSettings         []AutoHealSettingWindows  `tfschema:"auto_heal_setting"`
	UseManagedIdentityACR    bool                      `tfschema:"container_registry_use_managed_identity"`
	ContainerRegistryUserMSI string                    `tfschema:"container_registry_managed_identity_id"`
	DefaultDocuments         []string                  `tfschema:"default_documents"`
	Http2Enabled             bool                      `tfschema:"http2_enabled"`
	IpRestriction            []helpers.IpRestriction   `tfschema:"ip_restriction"`
	ScmUseMainIpRestriction  bool                      `tfschema:"scm_use_main_ip_restriction"`
	ScmIpRestriction         []helpers.IpRestriction   `tfschema:"scm_ip_restriction"`
	LoadBalancing            string                    `tfschema:"load_balancing_mode"`
	LocalMysql               bool                      `tfschema:"local_mysql"`
	ManagedPipelineMode      string                    `tfschema:"managed_pipeline_mode"`
	RemoteDebugging          bool                      `tfschema:"remote_debugging"`
	RemoteDebuggingVersion   string                    `tfschema:"remote_debugging_version"`
	ScmType                  string                    `tfschema:"scm_type"`
	Use32BitWorker           bool                      `tfschema:"use_32_bit_worker"`
	WebSockets               bool                      `tfschema:"websockets"`
	FtpsState                string                    `tfschema:"ftps_state"`
	HealthCheckPath          string                    `tfschema:"health_check_path"`
	NumberOfWorkers          int                       `tfschema:"number_of_workers"`
	ApplicationStack         []ApplicationStackWindows `tfschema:"application_stack"`
	VirtualApplications      []VirtualApplication      `tfschema:"virtual_application"`
	MinTlsVersion            string                    `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion         string                    `tfschema:"scm_minimum_tls_version"`
	AutoSwapSlotName         string                    `tfschema:"auto_swap_slot_name"`
	Cors                     []helpers.CorsSetting     `tfschema:"cors"`
	DetailedErrorLogging     bool                      `tfschema:"detailed_error_logging"`
	WindowsFxVersion         string                    `tfschema:"windows_fx_version"`
	// TODO new properties / blocks
	// SiteLimits []SiteLimitsSettings `tfschema:"site_limits"` // TODO - New block to (possibly) support? No way to configure this in the portal?
	// PushSettings - Supported in SDK, but blocked by manual step needed for connecting app to notification hub.
}

type SiteConfigLinux struct {
	AlwaysOn                bool                    `tfschema:"always_on"`
	ApiManagementConfigId   string                  `tfschema:"api_management_config_id"`
	AppCommandLine          string                  `tfschema:"app_command_line"`
	AutoHeal                bool                    `tfschema:"auto_heal"`
	AutoHealSettings        []AutoHealSettingLinux  `tfschema:"auto_heal_setting"`
	UseManagedIdentityACR   bool                    `tfschema:"container_registry_use_managed_identity"`
	ContainerRegistryMSI    string                  `tfschema:"container_registry_managed_identity_id"`
	DefaultDocuments        []string                `tfschema:"default_documents"`
	Http2Enabled            bool                    `tfschema:"http2_enabled"`
	IpRestriction           []helpers.IpRestriction `tfschema:"ip_restriction"`
	ScmUseMainIpRestriction bool                    `tfschema:"scm_use_main_ip_restriction"`
	ScmIpRestriction        []helpers.IpRestriction `tfschema:"scm_ip_restriction"`
	LoadBalancing           string                  `tfschema:"load_balancing_mode"`
	LocalMysql              bool                    `tfschema:"local_mysql"`
	ManagedPipelineMode     string                  `tfschema:"managed_pipeline_mode"`
	RemoteDebugging         bool                    `tfschema:"remote_debugging"`
	RemoteDebuggingVersion  string                  `tfschema:"remote_debugging_version"`
	ScmType                 string                  `tfschema:"scm_type"`
	Use32BitWorker          bool                    `tfschema:"use_32_bit_worker"`
	WebSockets              bool                    `tfschema:"websockets"`
	FtpsState               string                  `tfschema:"ftps_state"`
	HealthCheckPath         string                  `tfschema:"health_check_path"`
	NumberOfWorkers         int                     `tfschema:"number_of_workers"`
	ApplicationStack        []ApplicationStackLinux `tfschema:"application_stack"`
	MinTlsVersion           string                  `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion        string                  `tfschema:"scm_minimum_tls_version"`
	AutoSwapSlotName        string                  `tfschema:"auto_swap_slot_name"`
	Cors                    []helpers.CorsSetting   `tfschema:"cors"`
	DetailedErrorLogging    bool                    `tfschema:"detailed_error_logging"`
	LinuxFxVersion          string                  `tfschema:"linux_fx_version"`
	// SiteLimits []SiteLimitsSettings `tfschema:"site_limits"` // TODO - New block to (possibly) support? No way to configure this in the portal?
}

func siteConfigSchemaWindows() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
				},

				"api_management_config_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: apimValidate.ApiManagementID,
				},

				"application_stack": windowsApplicationStackSchema(),

				"app_command_line": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"auto_heal": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
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

				"container_registry_managed_identity_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: msiValidate.UserAssignedIdentityID,
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

				"ip_restriction": helpers.IpRestrictionSchema(),

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"scm_ip_restriction": helpers.IpRestrictionSchema(),

				"local_mysql": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
				},

				"load_balancing_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
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
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ManagedPipelineModeClassic),
						string(web.ManagedPipelineModeIntegrated),
					}, true),
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
					DiffSuppressFunc: suppress.CaseDifference,
				},

				"scm_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"use_32_bit_worker": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"websockets": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"ftps_state": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
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

				"number_of_workers": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 100),
				},

				"minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SupportedTLSVersionsOneFullStopZero),
						string(web.SupportedTLSVersionsOneFullStopOne),
						string(web.SupportedTLSVersionsOneFullStopTwo),
					}, false),
				},

				"scm_minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SupportedTLSVersionsOneFullStopZero),
						string(web.SupportedTLSVersionsOneFullStopOne),
						string(web.SupportedTLSVersionsOneFullStopTwo),
					}, false),
				},

				"cors": helpers.CorsSettingsSchema(),

				"virtual_application": virtualApplicationsSchema(),

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

func siteConfigSchemaWindowsComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"api_management_config_id": {
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

				"container_registry_managed_identity_id": {
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

				"ip_restriction": helpers.IpRestrictionSchemaComputed(),

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"scm_ip_restriction": helpers.IpRestrictionSchemaComputed(),

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

				"websockets": {
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

				"cors": helpers.CorsSettingsSchemaComputed(),

				"virtual_application": virtualApplicationsSchemaComputed(),

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

				"windows_fx_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func siteConfigSchemaLinux() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
				},

				"api_management_config_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: apimValidate.ApiManagementID,
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

				"container_registry_managed_identity_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: msiValidate.UserAssignedIdentityID,
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

				"ip_restriction": helpers.IpRestrictionSchema(),

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"scm_ip_restriction": helpers.IpRestrictionSchema(),

				"local_mysql": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
				},

				"load_balancing_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
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
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ManagedPipelineModeClassic),
						string(web.ManagedPipelineModeIntegrated),
					}, true),
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
					DiffSuppressFunc: suppress.CaseDifference,
				},

				"scm_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"use_32_bit_worker": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"websockets": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"ftps_state": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
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

				"number_of_workers": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 100),
				},

				"minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SupportedTLSVersionsOneFullStopZero),
						string(web.SupportedTLSVersionsOneFullStopOne),
						string(web.SupportedTLSVersionsOneFullStopTwo),
					}, false),
				},

				"scm_minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SupportedTLSVersionsOneFullStopZero),
						string(web.SupportedTLSVersionsOneFullStopOne),
						string(web.SupportedTLSVersionsOneFullStopTwo),
					}, false),
				},

				"cors": helpers.CorsSettingsSchema(),

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

func siteConfigSchemaLinuxComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"api_management_config_id": {
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

				"container_registry_managed_identity_id": {
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

				"ip_restriction": helpers.IpRestrictionSchemaComputed(),

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"scm_ip_restriction": helpers.IpRestrictionSchemaComputed(),

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

				"websockets": {
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

				"cors": helpers.CorsSettingsSchemaComputed(),

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

				"windows_fx_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

type ApplicationStackWindows struct {
	NetFrameworkVersion     string `tfschema:"dotnet_framework_version"`
	PhpVersion              string `tfschema:"php_version"`
	JavaVersion             string `tfschema:"java_version"`
	PythonVersion           string `tfschema:"python_version"` // Linux Only?
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
				"dotnet_framework_version": { // Windows Only
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
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
					Computed: true,
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
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"1.7",
						"1.8",
						"11",
					}, false),
				},

				"java_container": {Type: pluginsdk.TypeString,
					Optional: true,
					Computed: true,
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
					Computed: true,
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
				"dotnet_framework_version": { // Windows Only
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
	NetFrameworkVersion string `tfschema:"dotnet_framework_version"`
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
				"dotnet_framework_version": { // Windows Only
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"2.1",
						"3.1",
						"5.0",
					}, false),
				},

				"php_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"5.6", // TODO - Remove? 5.6 is available, but deprecated in the service
						"7.2", // TODO - Remove? 7.2 is available, but deprecated in the service
						"7.3",
						"7.4",
					}, false),
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
				},

				"node_version": { // Discarded by service if JavaVersion is specified
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
				},

				"java_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty, // There a significant number of variables here, and the versions are not uniformly formatted.
					// TODO - Needs notes in the docs for this to help users navigate the inconsistencies in the service. e.g. jre8 va java8 etc
				},

				"java_server": {Type: pluginsdk.TypeString,
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
				"dotnet_framework_version": { // Windows Only
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
	Win32Status     string `tfschema:"win_32_status"`
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
					Default:      0,
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

							"win_32_status": {
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

							"win_32_status": {
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

							"win_32_status": {
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

							"win_32_status": {
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
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"virtual_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"physical_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"preload": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
				},

				"virtual_directory": {
					Type:     pluginsdk.TypeSet,
					Computed: true,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"virtual_path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Computed:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"physical_path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Computed:     true,
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
		Type:     pluginsdk.TypeSet,
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
					Type:     pluginsdk.TypeSet,
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

func storageAccountSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Computed: true,
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

func storageAccountSchemaComputed() *pluginsdk.Schema {
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

func backupSchema() *pluginsdk.Schema {
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

func backupSchemaComputed() *pluginsdk.Schema {
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

func connectionStringSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Computed: true,
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
					DiffSuppressFunc: suppress.CaseDifference,
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

func connectionStringSchemaComputed() *pluginsdk.Schema {
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
	FileSystems      []FileSystem           `tfschema:"file_system"`
	AzureBlobStorage []AzureBlobStorageHttp `tfschema:"azure_blob_storage"`
}

type AzureBlobStorageHttp struct {
	SasUrl          string `tfschema:"sas_url"`
	RetentionInDays int    `tfschema:"retention_in_days"`
}

type FileSystem struct {
	RetentionMB   int `tfschema:"retention_in_mb"`
	RetentionDays int `tfschema:"retention_in_days"`
}

func logsConfigSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
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

func logsConfigSchemaComputed() *pluginsdk.Schema {
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
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"file_system_level": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "Off",
					ValidateFunc: validation.StringInSlice([]string{
						string(web.LogLevelError),
						string(web.LogLevelInformation),
						string(web.LogLevelOff),
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
						string(web.LogLevelOff),
						string(web.LogLevelVerbose),
						string(web.LogLevelWarning),
					}, false),
				},
				"sas_url": {
					Type:      pluginsdk.TypeString,
					Required:  true,
					Sensitive: true,
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
		Computed: true,
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
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
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
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"sas_url": {
					Type:      pluginsdk.TypeString,
					Required:  true,
					Sensitive: true,
				},
				"retention_in_days": {
					Type:     pluginsdk.TypeInt,
					Required: true,
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

func expandSiteConfigWindows(siteConfig []SiteConfigWindows) (*web.SiteConfig, *string, error) {
	if len(siteConfig) == 0 {
		return nil, nil, nil
	}
	expanded := &web.SiteConfig{}
	currentStack := ""

	winSiteConfig := siteConfig[0]
	expanded.AlwaysOn = utils.Bool(winSiteConfig.AlwaysOn)

	if winSiteConfig.ApiManagementConfigId != "" {
		expanded.APIManagementConfig = &web.APIManagementConfig{
			ID: utils.String(winSiteConfig.ApiManagementConfigId),
		}
	}

	if winSiteConfig.AppCommandLine != "" {
		expanded.AppCommandLine = utils.String(winSiteConfig.AppCommandLine)
	}

	if len(winSiteConfig.ApplicationStack) == 1 {
		winAppStack := winSiteConfig.ApplicationStack[0]
		if winAppStack.NetFrameworkVersion != "" {
			expanded.NetFrameworkVersion = utils.String(winAppStack.NetFrameworkVersion)
		}

		if winAppStack.PhpVersion != "" {
			expanded.PhpVersion = utils.String(winAppStack.PhpVersion)
		}

		if winAppStack.NodeVersion != "" {
			expanded.NodeVersion = utils.String(winAppStack.NodeVersion)
		}

		if winAppStack.PythonVersion != "" {
			expanded.PythonVersion = utils.String(winAppStack.PythonVersion)
		}

		if winAppStack.JavaVersion != "" {
			expanded.JavaVersion = utils.String(winAppStack.JavaVersion)
		}

		if winAppStack.JavaContainer != "" {
			expanded.JavaContainer = utils.String(winAppStack.JavaContainer)
		}

		if winAppStack.JavaContainerVersion != "" {
			expanded.JavaContainerVersion = utils.String(winAppStack.JavaContainerVersion)
		}

		if winAppStack.DockerContainerName != "" {
			expanded.WindowsFxVersion = utils.String(fmt.Sprintf("DOCKER|%s/%s:%s", winAppStack.DockerContainerRegistry, winAppStack.DockerContainerName, winAppStack.DockerContainerTag))
		}
		currentStack = winAppStack.CurrentStack
	}

	if winSiteConfig.VirtualApplications != nil {
		expanded.VirtualApplications = expandVirtualApplications(winSiteConfig.VirtualApplications)
	}

	expanded.AcrUseManagedIdentityCreds = utils.Bool(winSiteConfig.UseManagedIdentityACR)

	expanded.AcrUserManagedIdentityID = utils.String(winSiteConfig.ContainerRegistryUserMSI)

	if len(winSiteConfig.DefaultDocuments) != 0 {
		expanded.DefaultDocuments = &winSiteConfig.DefaultDocuments
	}

	expanded.HTTP20Enabled = utils.Bool(winSiteConfig.Http2Enabled)

	if len(winSiteConfig.IpRestriction) != 0 {
		ipRestrictions, err := helpers.ExpandIpRestrictions(winSiteConfig.IpRestriction)
		if err != nil {
			return nil, nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	expanded.ScmIPSecurityRestrictionsUseMain = utils.Bool(winSiteConfig.ScmUseMainIpRestriction)

	if len(winSiteConfig.ScmIpRestriction) != 0 {
		scmIpRestrictions, err := helpers.ExpandIpRestrictions(winSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	expanded.LocalMySQLEnabled = utils.Bool(winSiteConfig.LocalMysql)

	if winSiteConfig.LoadBalancing != "" {
		expanded.LoadBalancing = web.SiteLoadBalancing(winSiteConfig.LoadBalancing)
	}

	if winSiteConfig.ManagedPipelineMode != "" {
		expanded.ManagedPipelineMode = web.ManagedPipelineMode(winSiteConfig.ManagedPipelineMode)
	}

	expanded.RemoteDebuggingEnabled = utils.Bool(winSiteConfig.RemoteDebugging)

	if winSiteConfig.RemoteDebuggingVersion != "" {
		expanded.RemoteDebuggingVersion = utils.String(winSiteConfig.RemoteDebuggingVersion)
	}

	if winSiteConfig.ScmType != "" {
		expanded.ScmType = web.ScmType(winSiteConfig.ScmType)
	}

	expanded.Use32BitWorkerProcess = utils.Bool(winSiteConfig.Use32BitWorker)

	expanded.WebSocketsEnabled = utils.Bool(winSiteConfig.WebSockets)

	if winSiteConfig.FtpsState != "" {
		expanded.FtpsState = web.FtpsState(winSiteConfig.FtpsState)
	}

	if winSiteConfig.HealthCheckPath != "" {
		expanded.HealthCheckPath = utils.String(winSiteConfig.HealthCheckPath)
	}

	if winSiteConfig.NumberOfWorkers != 0 {
		expanded.NumberOfWorkers = utils.Int32(int32(winSiteConfig.NumberOfWorkers))
	}

	if winSiteConfig.MinTlsVersion != "" {
		expanded.MinTLSVersion = web.SupportedTLSVersions(winSiteConfig.MinTlsVersion)
	}

	if winSiteConfig.ScmMinTlsVersion != "" {
		expanded.ScmMinTLSVersion = web.SupportedTLSVersions(winSiteConfig.ScmMinTlsVersion)
	}

	if winSiteConfig.AutoSwapSlotName != "" {
		expanded.AutoSwapSlotName = utils.String(winSiteConfig.AutoSwapSlotName)
	}

	if len(winSiteConfig.Cors) != 0 {
		expanded.Cors = helpers.ExpandCorsSettings(winSiteConfig.Cors)
	}

	expanded.AutoHealEnabled = utils.Bool(winSiteConfig.AutoHeal)
	if len(winSiteConfig.AutoHealSettings) != 0 {
		expanded.AutoHealRules = expandAutoHealSettingsWindows(winSiteConfig.AutoHealSettings)
	}

	return expanded, &currentStack, nil
}

func expandSiteConfigLinux(siteConfig []SiteConfigLinux) (*web.SiteConfig, error) {
	if len(siteConfig) == 0 {
		return nil, nil
	}
	expanded := &web.SiteConfig{}

	linuxSiteConfig := siteConfig[0]
	expanded.AlwaysOn = utils.Bool(linuxSiteConfig.AlwaysOn)

	if linuxSiteConfig.ApiManagementConfigId != "" {
		expanded.APIManagementConfig = &web.APIManagementConfig{
			ID: utils.String(linuxSiteConfig.ApiManagementConfigId),
		}
	}

	if linuxSiteConfig.AppCommandLine != "" {
		expanded.AppCommandLine = utils.String(linuxSiteConfig.AppCommandLine)
	}

	if len(linuxSiteConfig.ApplicationStack) == 1 {
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

	if linuxSiteConfig.ContainerRegistryMSI != "" {
		expanded.AcrUserManagedIdentityID = utils.String(linuxSiteConfig.ContainerRegistryMSI)
	}

	if len(linuxSiteConfig.DefaultDocuments) != 0 {
		expanded.DefaultDocuments = &linuxSiteConfig.DefaultDocuments
	}

	expanded.HTTP20Enabled = utils.Bool(linuxSiteConfig.Http2Enabled)

	if len(linuxSiteConfig.IpRestriction) != 0 {
		ipRestrictions, err := helpers.ExpandIpRestrictions(linuxSiteConfig.IpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	expanded.ScmIPSecurityRestrictionsUseMain = utils.Bool(linuxSiteConfig.ScmUseMainIpRestriction)

	if len(linuxSiteConfig.ScmIpRestriction) != 0 {
		scmIpRestrictions, err := helpers.ExpandIpRestrictions(linuxSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	expanded.LocalMySQLEnabled = utils.Bool(linuxSiteConfig.LocalMysql)

	if linuxSiteConfig.LoadBalancing != "" {
		expanded.LoadBalancing = web.SiteLoadBalancing(linuxSiteConfig.LoadBalancing)
	}

	if linuxSiteConfig.ManagedPipelineMode != "" {
		expanded.ManagedPipelineMode = web.ManagedPipelineMode(linuxSiteConfig.ManagedPipelineMode)
	}

	if linuxSiteConfig.RemoteDebugging {
		expanded.RemoteDebuggingEnabled = utils.Bool(linuxSiteConfig.RemoteDebugging)
	}

	if linuxSiteConfig.RemoteDebuggingVersion != "" {
		expanded.RemoteDebuggingVersion = utils.String(linuxSiteConfig.RemoteDebuggingVersion)
	}

	if linuxSiteConfig.ScmType != "" {
		expanded.ScmType = web.ScmType(linuxSiteConfig.ScmType)
	}

	expanded.Use32BitWorkerProcess = utils.Bool(linuxSiteConfig.Use32BitWorker)

	expanded.WebSocketsEnabled = utils.Bool(linuxSiteConfig.WebSockets)

	if linuxSiteConfig.FtpsState != "" {
		expanded.FtpsState = web.FtpsState(linuxSiteConfig.FtpsState)
	}

	if linuxSiteConfig.HealthCheckPath != "" {
		expanded.HealthCheckPath = utils.String(linuxSiteConfig.HealthCheckPath)
	}

	if linuxSiteConfig.NumberOfWorkers != 0 {
		expanded.NumberOfWorkers = utils.Int32(int32(linuxSiteConfig.NumberOfWorkers))
	}

	if linuxSiteConfig.MinTlsVersion != "" {
		expanded.MinTLSVersion = web.SupportedTLSVersions(linuxSiteConfig.MinTlsVersion)
	}

	if linuxSiteConfig.ScmMinTlsVersion != "" {
		expanded.ScmMinTLSVersion = web.SupportedTLSVersions(linuxSiteConfig.ScmMinTlsVersion)
	}

	if linuxSiteConfig.AutoSwapSlotName != "" {
		expanded.AutoSwapSlotName = utils.String(linuxSiteConfig.AutoSwapSlotName)
	}

	if len(linuxSiteConfig.Cors) != 0 {
		expanded.Cors = helpers.ExpandCorsSettings(linuxSiteConfig.Cors)
	}

	expanded.AutoHealEnabled = utils.Bool(linuxSiteConfig.AutoHeal)
	if len(linuxSiteConfig.AutoHealSettings) != 0 {
		expanded.AutoHealRules = expandAutoHealSettingsLinux(linuxSiteConfig.AutoHealSettings)
	}

	return expanded, nil
}

func expandLogsConfig(config []LogsConfig) *web.SiteLogsConfig {
	if len(config) == 0 {
		return nil
	}

	siteLogsConfig := &web.SiteLogsConfig{
		SiteLogsConfigProperties: &web.SiteLogsConfigProperties{},
	}
	logsConfig := config[0]

	if len(logsConfig.ApplicationLogs) == 1 {
		appLogs := logsConfig.ApplicationLogs[0]
		siteLogsConfig.SiteLogsConfigProperties.ApplicationLogs = &web.ApplicationLogsConfig{
			FileSystem: &web.FileSystemApplicationLogsConfig{ // TODO - Does this conflict with the use of `AzureBlobStorage` below?
				Level: web.LogLevel(appLogs.FileSystemLevel),
			},
		}
		if len(appLogs.AzureBlobStorage) == 1 {
			appLogsBlobs := appLogs.AzureBlobStorage[0]
			siteLogsConfig.SiteLogsConfigProperties.ApplicationLogs.AzureBlobStorage = &web.AzureBlobStorageApplicationLogsConfig{
				Level:           web.LogLevel(appLogsBlobs.Level),
				SasURL:          utils.String(appLogsBlobs.SasUrl),
				RetentionInDays: utils.Int32(int32(appLogsBlobs.RetentionInDays)),
			}
		}
	}

	if len(logsConfig.HttpLogs) == 1 {
		httpLogs := logsConfig.HttpLogs[0]
		siteLogsConfig.HTTPLogs = &web.HTTPLogsConfig{}

		if len(httpLogs.FileSystems) == 1 {
			httpLogFileSystem := httpLogs.FileSystems[0]
			siteLogsConfig.HTTPLogs.FileSystem = &web.FileSystemHTTPLogsConfig{
				Enabled:         utils.Bool(true),
				RetentionInMb:   utils.Int32(int32(httpLogFileSystem.RetentionMB)),
				RetentionInDays: utils.Int32(int32(httpLogFileSystem.RetentionDays)),
			}
		}

		if len(httpLogs.AzureBlobStorage) == 1 {
			httpLogsBlobStorage := httpLogs.AzureBlobStorage[0]
			siteLogsConfig.HTTPLogs.AzureBlobStorage = &web.AzureBlobStorageHTTPLogsConfig{
				Enabled:         utils.Bool(httpLogsBlobStorage.SasUrl != ""),
				SasURL:          utils.String(httpLogsBlobStorage.SasUrl),
				RetentionInDays: utils.Int32(int32(httpLogsBlobStorage.RetentionInDays)),
			}
		}
	}

	siteLogsConfig.DetailedErrorMessages = &web.EnabledConfig{
		Enabled: utils.Bool(logsConfig.DetailedErrorMessages),
	}

	siteLogsConfig.FailedRequestsTracing = &web.EnabledConfig{
		Enabled: utils.Bool(logsConfig.FailedRequestTracing),
	}

	return siteLogsConfig
}

func expandBackupConfig(backupConfigs []Backup) *web.BackupRequest {
	if len(backupConfigs) == 0 {
		return nil
	}

	backupConfig := backupConfigs[0]
	backupSchedule := backupConfig.Schedule[0]
	backupRequest := &web.BackupRequest{
		BackupRequestProperties: &web.BackupRequestProperties{
			Enabled:           utils.Bool(backupConfig.Enabled),
			BackupName:        utils.String(backupConfig.Name),
			StorageAccountURL: utils.String(backupConfig.StorageAccountUrl),
			BackupSchedule: &web.BackupSchedule{
				FrequencyInterval:     utils.Int32(int32(backupSchedule.FrequencyInterval)),
				FrequencyUnit:         web.FrequencyUnit(backupSchedule.FrequencyUnit),
				KeepAtLeastOneBackup:  utils.Bool(backupSchedule.KeepAtLeastOneBackup),
				RetentionPeriodInDays: utils.Int32(int32(backupSchedule.RetentionPeriodDays)),
			},
		},
	}

	if backupSchedule.StartTime != "" {
		dateTimeToStart, _ := time.Parse(time.RFC3339, backupSchedule.StartTime)
		backupRequest.BackupRequestProperties.BackupSchedule.StartTime = &date.Time{Time: dateTimeToStart}
	}

	return backupRequest
}

func expandStorageConfig(storageConfigs []StorageAccount) *web.AzureStoragePropertyDictionaryResource {
	if len(storageConfigs) == 0 {
		return nil
	}
	storageAccounts := make(map[string]*web.AzureStorageInfoValue)
	for _, v := range storageConfigs {
		storageAccounts[v.Name] = &web.AzureStorageInfoValue{
			Type:        web.AzureStorageType(v.Type),
			AccountName: utils.String(v.AccountName),
			ShareName:   utils.String(v.ShareName),
			AccessKey:   utils.String(v.AccessKey),
			MountPath:   utils.String(v.MountPath),
		}
	}

	return &web.AzureStoragePropertyDictionaryResource{
		Properties: storageAccounts,
	}
}

func expandConnectionStrings(connectionStringsConfig []ConnectionString) *web.ConnectionStringDictionary {
	if len(connectionStringsConfig) == 0 {
		return nil
	}
	connectionStrings := make(map[string]*web.ConnStringValueTypePair)
	for _, v := range connectionStringsConfig {
		connectionStrings[v.Name] = &web.ConnStringValueTypePair{
			Value: utils.String(v.Value),
			Type:  web.ConnectionStringType(v.Type),
		}
	}

	return &web.ConnectionStringDictionary{
		Properties: connectionStrings,
	}
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

func flattenBackupConfig(backupRequest web.BackupRequest) []Backup {
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

	if props.BackupSchedule != nil {
		schedule := *props.BackupSchedule
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

func flattenLogsConfig(logsConfig web.SiteLogsConfig) []LogsConfig {
	if logsConfig.SiteLogsConfigProperties == nil {
		return nil
	}

	logs := LogsConfig{}
	props := *logsConfig.SiteLogsConfigProperties

	if props.ApplicationLogs != nil {
		applicationLog := ApplicationLog{}
		appLogs := *props.ApplicationLogs
		if appLogs.FileSystem != nil {
			applicationLog.FileSystemLevel = string(appLogs.FileSystem.Level)
		}

		if appLogs.AzureBlobStorage != nil {
			blobStorage := AzureBlobStorage{
				Level: string(appLogs.AzureBlobStorage.Level),
			}
			if appLogs.AzureBlobStorage.SasURL != nil {
				blobStorage.SasUrl = *appLogs.AzureBlobStorage.SasURL
			}
			if appLogs.AzureBlobStorage.RetentionInDays != nil {
				blobStorage.RetentionInDays = int(*appLogs.AzureBlobStorage.RetentionInDays)
			}
			applicationLog.AzureBlobStorage = []AzureBlobStorage{blobStorage}
		}
		logs.ApplicationLogs = []ApplicationLog{applicationLog}
	}

	if props.HTTPLogs != nil {
		httpLogs := *props.HTTPLogs
		httpLog := HttpLog{}

		if httpLogs.FileSystem != nil && *httpLogs.FileSystem.Enabled {
			fileSystem := FileSystem{}
			if httpLogs.FileSystem.RetentionInMb != nil {
				fileSystem.RetentionMB = int(*httpLogs.FileSystem.RetentionInMb)
			}

			if httpLogs.FileSystem.RetentionInDays != nil {
				fileSystem.RetentionDays = int(*httpLogs.FileSystem.RetentionInDays)
			}
			httpLog.FileSystems = []FileSystem{fileSystem}
		}

		if httpLogs.AzureBlobStorage != nil {
			blobStorage := AzureBlobStorageHttp{}
			if httpLogs.AzureBlobStorage.SasURL != nil {
				blobStorage.SasUrl = *httpLogs.AzureBlobStorage.SasURL
			}

			if httpLogs.AzureBlobStorage.RetentionInDays != nil {
				blobStorage.RetentionInDays = int(*httpLogs.AzureBlobStorage.RetentionInDays)
			}

			httpLog.AzureBlobStorage = []AzureBlobStorageHttp{blobStorage}
		}

		logs.HttpLogs = []HttpLog{httpLog}
	}

	if props.DetailedErrorMessages != nil {
		logs.DetailedErrorMessages = *props.DetailedErrorMessages.Enabled
	}

	if props.FailedRequestsTracing != nil {
		logs.FailedRequestTracing = *props.FailedRequestsTracing.Enabled
	}

	return []LogsConfig{logs}
}

func flattenSiteConfigWindows(appSiteConfig *web.SiteConfig, currentStack string) []SiteConfigWindows {
	if appSiteConfig == nil {
		return nil
	}

	siteConfig := SiteConfigWindows{
		ManagedPipelineMode: string(appSiteConfig.ManagedPipelineMode),
		ScmType:             string(appSiteConfig.ScmType),
		FtpsState:           string(appSiteConfig.FtpsState),
		MinTlsVersion:       string(appSiteConfig.MinTLSVersion),
		ScmMinTlsVersion:    string(appSiteConfig.ScmMinTLSVersion),
	}

	siteConfig.AlwaysOn = *appSiteConfig.AlwaysOn

	if appSiteConfig.APIManagementConfig != nil && appSiteConfig.APIManagementConfig.ID != nil {
		siteConfig.ApiManagementConfigId = *appSiteConfig.APIManagementConfig.ID
	}

	if appSiteConfig.AppCommandLine != nil {
		siteConfig.AppCommandLine = *appSiteConfig.AppCommandLine
	}

	if appSiteConfig.DefaultDocuments != nil {
		siteConfig.DefaultDocuments = *appSiteConfig.DefaultDocuments
	}

	siteConfig.UseManagedIdentityACR = *appSiteConfig.AcrUseManagedIdentityCreds

	if appSiteConfig.AcrUserManagedIdentityID != nil {
		siteConfig.ContainerRegistryUserMSI = *appSiteConfig.AcrUserManagedIdentityID
	}

	siteConfig.DetailedErrorLogging = *appSiteConfig.DetailedErrorLoggingEnabled

	siteConfig.Http2Enabled = *appSiteConfig.HTTP20Enabled

	if appSiteConfig.IPSecurityRestrictions != nil {
		siteConfig.IpRestriction = helpers.FlattenIpRestrictions(appSiteConfig.IPSecurityRestrictions)
	}

	siteConfig.ScmUseMainIpRestriction = *appSiteConfig.ScmIPSecurityRestrictionsUseMain

	if appSiteConfig.ScmIPSecurityRestrictions != nil {
		siteConfig.ScmIpRestriction = helpers.FlattenIpRestrictions(appSiteConfig.ScmIPSecurityRestrictions)
	}

	siteConfig.LocalMysql = *appSiteConfig.LocalMySQLEnabled

	siteConfig.LoadBalancing = string(appSiteConfig.LoadBalancing)

	if appSiteConfig.RemoteDebuggingEnabled != nil {
		siteConfig.RemoteDebugging = *appSiteConfig.RemoteDebuggingEnabled
	}

	if appSiteConfig.RemoteDebuggingVersion != nil {
		siteConfig.RemoteDebuggingVersion = *appSiteConfig.RemoteDebuggingVersion
	}

	if appSiteConfig.Use32BitWorkerProcess != nil {
		siteConfig.Use32BitWorker = *appSiteConfig.Use32BitWorkerProcess
	}

	siteConfig.WebSockets = *appSiteConfig.WebSocketsEnabled

	if appSiteConfig.HealthCheckPath != nil {
		siteConfig.HealthCheckPath = *appSiteConfig.HealthCheckPath
	}

	if appSiteConfig.NumberOfWorkers != nil {
		siteConfig.NumberOfWorkers = int(*appSiteConfig.NumberOfWorkers)
	}

	var winAppStack ApplicationStackWindows
	if appSiteConfig.NetFrameworkVersion != nil {
		winAppStack.NetFrameworkVersion = *appSiteConfig.NetFrameworkVersion
	}

	if appSiteConfig.PhpVersion != nil {
		winAppStack.PhpVersion = *appSiteConfig.PhpVersion
	}

	if appSiteConfig.NodeVersion != nil {
		winAppStack.NodeVersion = *appSiteConfig.NodeVersion
	}

	if appSiteConfig.PythonVersion != nil {
		winAppStack.PythonVersion = *appSiteConfig.PythonVersion
	}

	if appSiteConfig.JavaVersion != nil {
		winAppStack.JavaVersion = *appSiteConfig.JavaVersion
	}

	if appSiteConfig.JavaContainer != nil {
		winAppStack.JavaContainer = *appSiteConfig.JavaContainer
	}

	if appSiteConfig.JavaContainerVersion != nil {
		winAppStack.JavaContainerVersion = *appSiteConfig.JavaContainerVersion
	}

	if appSiteConfig.WindowsFxVersion != nil {
		siteConfig.WindowsFxVersion = *appSiteConfig.WindowsFxVersion
		// Decode the string to docker values
		parts := strings.Split(strings.TrimPrefix(siteConfig.WindowsFxVersion, "DOCKER|"), ":")
		winAppStack.DockerContainerTag = parts[1]
		path := strings.Split(parts[0], "/")
		winAppStack.DockerContainerRegistry = path[0]
		winAppStack.DockerContainerName = strings.TrimPrefix(parts[0], fmt.Sprintf("%s/", path[0]))
	}
	winAppStack.CurrentStack = currentStack

	siteConfig.ApplicationStack = []ApplicationStackWindows{winAppStack}

	siteConfig.VirtualApplications = flattenVirtualApplications(appSiteConfig.VirtualApplications)

	if appSiteConfig.AutoSwapSlotName != nil {
		siteConfig.AutoSwapSlotName = *appSiteConfig.AutoSwapSlotName
	}

	if appSiteConfig.Cors != nil {
		corsSettings := appSiteConfig.Cors
		cors := helpers.CorsSetting{}
		if corsSettings.SupportCredentials != nil {
			cors.SupportCredentials = *corsSettings.SupportCredentials
		}

		if corsSettings.AllowedOrigins != nil {
			cors.AllowedOrigins = *corsSettings.AllowedOrigins
		}
		siteConfig.Cors = []helpers.CorsSetting{cors}
	}

	siteConfig.AutoHeal = *appSiteConfig.AutoHealEnabled
	siteConfig.AutoHealSettings = flattenAutoHealSettingsWindows(appSiteConfig.AutoHealRules)

	return []SiteConfigWindows{siteConfig}
}

func flattenSiteConfigLinux(appSiteConfig *web.SiteConfig) []SiteConfigLinux {
	// TODO - Make this Linux flavoured...
	if appSiteConfig == nil {
		return nil
	}

	siteConfig := SiteConfigLinux{
		ManagedPipelineMode: string(appSiteConfig.ManagedPipelineMode),
		ScmType:             string(appSiteConfig.ScmType),
		FtpsState:           string(appSiteConfig.FtpsState),
		MinTlsVersion:       string(appSiteConfig.MinTLSVersion),
		ScmMinTlsVersion:    string(appSiteConfig.ScmMinTLSVersion),
	}

	if appSiteConfig.AlwaysOn != nil {
		siteConfig.AlwaysOn = *appSiteConfig.AlwaysOn
	}

	if appSiteConfig.APIManagementConfig != nil && appSiteConfig.APIManagementConfig.ID != nil {
		siteConfig.ApiManagementConfigId = *appSiteConfig.APIManagementConfig.ID
	}

	if appSiteConfig.AppCommandLine != nil {
		siteConfig.AppCommandLine = *appSiteConfig.AppCommandLine
	}

	siteConfig.UseManagedIdentityACR = *appSiteConfig.AcrUseManagedIdentityCreds

	if appSiteConfig.AcrUseManagedIdentityCreds != nil && *appSiteConfig.AcrUseManagedIdentityCreds && appSiteConfig.AcrUserManagedIdentityID != nil {
		siteConfig.ContainerRegistryMSI = *appSiteConfig.AcrUserManagedIdentityID
	}

	if appSiteConfig.DefaultDocuments != nil {
		siteConfig.DefaultDocuments = *appSiteConfig.DefaultDocuments
	}

	siteConfig.DetailedErrorLogging = *appSiteConfig.DetailedErrorLoggingEnabled

	siteConfig.Http2Enabled = *appSiteConfig.HTTP20Enabled

	if appSiteConfig.IPSecurityRestrictions != nil {
		siteConfig.IpRestriction = helpers.FlattenIpRestrictions(appSiteConfig.IPSecurityRestrictions)
	}

	siteConfig.ScmUseMainIpRestriction = *appSiteConfig.ScmIPSecurityRestrictionsUseMain

	if appSiteConfig.ScmIPSecurityRestrictions != nil {
		siteConfig.ScmIpRestriction = helpers.FlattenIpRestrictions(appSiteConfig.ScmIPSecurityRestrictions)
	}

	siteConfig.LocalMysql = *appSiteConfig.LocalMySQLEnabled

	siteConfig.LoadBalancing = string(appSiteConfig.LoadBalancing)

	siteConfig.RemoteDebugging = *appSiteConfig.RemoteDebuggingEnabled

	if appSiteConfig.RemoteDebuggingVersion != nil {
		siteConfig.RemoteDebuggingVersion = *appSiteConfig.RemoteDebuggingVersion
	}

	if appSiteConfig.Use32BitWorkerProcess != nil {
		siteConfig.Use32BitWorker = *appSiteConfig.Use32BitWorkerProcess
	}

	siteConfig.WebSockets = *appSiteConfig.WebSocketsEnabled

	if appSiteConfig.HealthCheckPath != nil {
		siteConfig.HealthCheckPath = *appSiteConfig.HealthCheckPath
	}

	if appSiteConfig.NumberOfWorkers != nil {
		siteConfig.NumberOfWorkers = int(*appSiteConfig.NumberOfWorkers)
	}

	var linuxAppStack ApplicationStackLinux

	if appSiteConfig.LinuxFxVersion != nil {
		siteConfig.LinuxFxVersion = *appSiteConfig.LinuxFxVersion
		// Decode the string to docker values
		linuxAppStack = decodeApplicationStackLinux(siteConfig.LinuxFxVersion)
	}

	siteConfig.ApplicationStack = []ApplicationStackLinux{linuxAppStack}

	if appSiteConfig.LinuxFxVersion != nil {
		siteConfig.LinuxFxVersion = *appSiteConfig.LinuxFxVersion
	}

	if appSiteConfig.AutoSwapSlotName != nil {
		siteConfig.AutoSwapSlotName = *appSiteConfig.AutoSwapSlotName
	}

	if appSiteConfig.Cors != nil {
		corsSettings := appSiteConfig.Cors
		cors := helpers.CorsSetting{}
		if corsSettings.SupportCredentials != nil {
			cors.SupportCredentials = *corsSettings.SupportCredentials
		}

		if corsSettings.AllowedOrigins != nil {
			cors.AllowedOrigins = *corsSettings.AllowedOrigins
		}
		siteConfig.Cors = []helpers.CorsSetting{cors}
	}

	siteConfig.AutoHeal = *appSiteConfig.AutoHealEnabled
	siteConfig.AutoHealSettings = flattenAutoHealSettingsLinux(appSiteConfig.AutoHealRules)

	return []SiteConfigLinux{siteConfig}
}

func flattenStorageAccounts(appStorageAccounts web.AzureStoragePropertyDictionaryResource) []StorageAccount {
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

func flattenConnectionStrings(appConnectionStrings web.ConnectionStringDictionary) []ConnectionString {
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

func expandAppSettings(settings map[string]string) *web.StringDictionary {
	appSettings := make(map[string]*string)
	for k, v := range settings {
		appSettings[k] = utils.String(v)
	}

	return &web.StringDictionary{
		Properties: appSettings,
	}
}

func flattenAppSettings(input web.StringDictionary) map[string]string {
	unmanagedSettings := []string{
		"DIAGNOSTICS_AZUREBLOBCONTAINERSASURL",
		"DIAGNOSTICS_AZUREBLOBRETENTIONINDAYS",
		"WEBSITE_HTTPLOGGING_CONTAINER_URL",
		"WEBSITE_HTTPLOGGING_RETENTION_DAYS",
	}

	appSettings := helpers.FlattenWebStringDictionary(input)

	// Remove the settings the service adds when logging settings are specified.
	for _, v := range unmanagedSettings { //nolint:typecheck
		delete(appSettings, v)
	}

	return appSettings
}

func flattenVirtualApplications(appVirtualApplications *[]web.VirtualApplication) []VirtualApplication {
	if appVirtualApplications == nil {
		return nil
	}

	var virtualApplications []VirtualApplication
	for _, v := range *appVirtualApplications {
		virtualApp := VirtualApplication{
			VirtualPath:  utils.NormalizeNilableString(v.VirtualPath),
			PhysicalPath: utils.NormalizeNilableString(v.PhysicalPath),
			Preload:      *v.PreloadEnabled,
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

func expandAutoHealSettingsWindows(autoHealSettings []AutoHealSettingWindows) *web.AutoHealRules {
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

	resultTrigger := AutoHealTriggerWindows{}
	resultActions := AutoHealActionWindows{}
	// Triggers
	if autoHealRules.Triggers != nil {
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

		if triggers.PrivateBytesInKB != nil {
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

		resultActions = AutoHealActionWindows{
			ActionType:         string(actions.ActionType),
			CustomAction:       customActions,
			MinimumProcessTime: utils.NormalizeNilableString(actions.MinProcessExecutionTime),
		}
	}

	return []AutoHealSettingWindows{{
		Triggers: []AutoHealTriggerWindows{resultTrigger},
		Actions:  []AutoHealActionWindows{resultActions},
	}}
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

	resultTrigger := AutoHealTriggerLinux{}
	resultActions := AutoHealActionLinux{}
	// Triggers
	if autoHealRules.Triggers != nil {
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
	}

	// Actions
	if autoHealRules.Actions != nil {
		actions := *autoHealRules.Actions

		resultActions = AutoHealActionLinux{
			ActionType:         string(actions.ActionType),
			MinimumProcessTime: utils.NormalizeNilableString(actions.MinProcessExecutionTime),
		}
	}

	return []AutoHealSettingLinux{{
		Triggers: []AutoHealTriggerLinux{resultTrigger},
		Actions:  []AutoHealActionLinux{resultActions},
	}}
}
