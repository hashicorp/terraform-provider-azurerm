package helpers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-03-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	apimValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SiteConfigLinuxWebAppSlot struct {
	AlwaysOn                bool                    `tfschema:"always_on"`
	ApiManagementConfigId   string                  `tfschema:"api_management_api_id"`
	ApiDefinition           string                  `tfschema:"api_definition_url"`
	AppCommandLine          string                  `tfschema:"app_command_line"`
	AutoHeal                bool                    `tfschema:"auto_heal_enabled"`
	AutoHealSettings        []AutoHealSettingLinux  `tfschema:"auto_heal_setting"`
	AutoSwapSlotName        string                  `tfschema:"auto_swap_slot_name"`
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
	WorkerCount             int                     `tfschema:"worker_count"`
	ApplicationStack        []ApplicationStackLinux `tfschema:"application_stack"`
	MinTlsVersion           string                  `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion        string                  `tfschema:"scm_minimum_tls_version"`
	Cors                    []CorsSetting           `tfschema:"cors"`
	DetailedErrorLogging    bool                    `tfschema:"detailed_error_logging_enabled"`
	LinuxFxVersion          string                  `tfschema:"linux_fx_version"`
	VnetRouteAllEnabled     bool                    `tfschema:"vnet_route_all_enabled"`
	// SiteLimits []SiteLimitsSettings `tfschema:"site_limits"` // TODO - New block to (possibly) support? No way to configure this in the portal?
}

func SiteConfigSchemaLinuxWebAppSlot() *pluginsdk.Schema {
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
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
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

				"auto_swap_slot_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

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

type SiteConfigWindowsWebAppSlot struct {
	AlwaysOn                 bool                      `tfschema:"always_on"`
	ApiManagementConfigId    string                    `tfschema:"api_management_api_id"`
	ApiDefinition            string                    `tfschema:"api_definition_url"`
	ApplicationStack         []ApplicationStackWindows `tfschema:"application_stack"`
	AppCommandLine           string                    `tfschema:"app_command_line"`
	AutoHeal                 bool                      `tfschema:"auto_heal_enabled"`
	AutoHealSettings         []AutoHealSettingWindows  `tfschema:"auto_heal_setting"`
	AutoSwapSlotName         string                    `tfschema:"auto_swap_slot_name"`
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
	VirtualApplications      []VirtualApplication      `tfschema:"virtual_application"`
	MinTlsVersion            string                    `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion         string                    `tfschema:"scm_minimum_tls_version"`
	Cors                     []CorsSetting             `tfschema:"cors"`
	DetailedErrorLogging     bool                      `tfschema:"detailed_error_logging_enabled"`
	WindowsFxVersion         string                    `tfschema:"windows_fx_version"`
	VnetRouteAllEnabled      bool                      `tfschema:"vnet_route_all_enabled"`
}

func SiteConfigSchemaWindowsWebAppSlot() *pluginsdk.Schema {
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

				"auto_swap_slot_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

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

				"windows_fx_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func ExpandSiteConfigLinuxWebAppSlot(siteConfig []SiteConfigLinuxWebAppSlot, existing *web.SiteConfig, metadata sdk.ResourceMetaData) (*web.SiteConfig, error) {
	if len(siteConfig) == 0 {
		return nil, nil
	}
	expanded := &web.SiteConfig{}
	if existing != nil {
		expanded = existing
	}

	linuxSlotSiteConfig := siteConfig[0]

	if metadata.ResourceData.HasChange("site_config.0.always_on") {
		expanded.AlwaysOn = pointer.To(linuxSlotSiteConfig.AlwaysOn)
	}

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.APIManagementConfig = &web.APIManagementConfig{
			ID: pointer.To(linuxSlotSiteConfig.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.APIDefinition = &web.APIDefinitionInfo{
			URL: pointer.To(linuxSlotSiteConfig.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = pointer.To(linuxSlotSiteConfig.AppCommandLine)
	}

	if metadata.ResourceData.HasChange("site_config.0.application_stack") {
		if len(linuxSlotSiteConfig.ApplicationStack) == 1 {
			linuxAppStack := linuxSlotSiteConfig.ApplicationStack[0]
			if linuxAppStack.NetFrameworkVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("DOTNETCORE|%s", linuxAppStack.NetFrameworkVersion))
			}

			if linuxAppStack.PhpVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("PHP|%s", linuxAppStack.PhpVersion))
			}

			if linuxAppStack.NodeVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("NODE|%s", linuxAppStack.NodeVersion))
			}

			if linuxAppStack.PythonVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("PYTHON|%s", linuxAppStack.PythonVersion))
			}

			if linuxAppStack.RubyVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("RUBY|%s", linuxAppStack.RubyVersion))
			}

			if linuxAppStack.JavaServer != "" {
				javaString, err := JavaLinuxFxStringBuilder(linuxAppStack.JavaVersion, linuxAppStack.JavaServer, linuxAppStack.JavaServerVersion)
				if err != nil {
					return nil, fmt.Errorf("could not build linuxFxVersion string: %+v", err)
				}
				expanded.LinuxFxVersion = javaString
			}

			if linuxAppStack.DockerImage != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("DOCKER|%s:%s", linuxAppStack.DockerImage, linuxAppStack.DockerImageTag))
			}
		} else {
			expanded.LinuxFxVersion = pointer.To("")
		}
	}

	expanded.AcrUseManagedIdentityCreds = pointer.To(linuxSlotSiteConfig.UseManagedIdentityACR)

	if metadata.ResourceData.HasChange("site_config.0.container_registry_managed_identity_client_id") {
		expanded.AcrUserManagedIdentityID = pointer.To(linuxSlotSiteConfig.ContainerRegistryMSI)
	}

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = &linuxSlotSiteConfig.DefaultDocuments
	}

	if metadata.ResourceData.HasChange("site_config.0.http2_enabled") {
		expanded.HTTP20Enabled = pointer.To(linuxSlotSiteConfig.Http2Enabled)
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := ExpandIpRestrictions(linuxSlotSiteConfig.IpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	expanded.ScmIPSecurityRestrictionsUseMain = pointer.To(linuxSlotSiteConfig.ScmUseMainIpRestriction)

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := ExpandIpRestrictions(linuxSlotSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.local_mysql_enabled") {
		expanded.LocalMySQLEnabled = pointer.To(linuxSlotSiteConfig.LocalMysql)
	}

	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = web.SiteLoadBalancing(linuxSlotSiteConfig.LoadBalancing)
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = web.ManagedPipelineMode(linuxSlotSiteConfig.ManagedPipelineMode)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_enabled") {
		expanded.RemoteDebuggingEnabled = pointer.To(linuxSlotSiteConfig.RemoteDebugging)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = pointer.To(linuxSlotSiteConfig.RemoteDebuggingVersion)
	}

	expanded.Use32BitWorkerProcess = pointer.To(linuxSlotSiteConfig.Use32BitWorker)

	if metadata.ResourceData.HasChange("site_config.0.websockets_enabled") {
		expanded.WebSocketsEnabled = pointer.To(linuxSlotSiteConfig.WebSockets)
	}

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		expanded.FtpsState = web.FtpsState(linuxSlotSiteConfig.FtpsState)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = pointer.To(linuxSlotSiteConfig.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.worker_count") {
		expanded.NumberOfWorkers = pointer.To(int32(linuxSlotSiteConfig.WorkerCount))
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTLSVersion = web.SupportedTLSVersions(linuxSlotSiteConfig.MinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTLSVersion = web.SupportedTLSVersions(linuxSlotSiteConfig.ScmMinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_swap_slot_name") {
		expanded.AutoSwapSlotName = pointer.To(linuxSlotSiteConfig.AutoSwapSlotName)
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		cors := ExpandCorsSettings(linuxSlotSiteConfig.Cors)
		if cors == nil {
			cors = &web.CorsSettings{
				AllowedOrigins: &[]string{},
			}
		}
		expanded.Cors = cors
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_heal_enabled") {
		expanded.AutoHealEnabled = pointer.To(linuxSlotSiteConfig.AutoHeal)
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_heal_setting") {
		expanded.AutoHealRules = expandAutoHealSettingsLinux(linuxSlotSiteConfig.AutoHealSettings)
	}

	if metadata.ResourceData.HasChange("site_config.0.vnet_route_all_enabled") {
		expanded.VnetRouteAllEnabled = pointer.To(linuxSlotSiteConfig.VnetRouteAllEnabled)
	}

	return expanded, nil
}

func (s *SiteConfigLinuxWebAppSlot) Flatten(appSiteSlotConfig *web.SiteConfig) {
	s.AlwaysOn = pointer.From(appSiteSlotConfig.AlwaysOn)
	s.AppCommandLine = pointer.From(appSiteSlotConfig.AppCommandLine)
	s.AutoHeal = pointer.From(appSiteSlotConfig.AutoHealEnabled)
	s.AutoHealSettings = flattenAutoHealSettingsLinux(appSiteSlotConfig.AutoHealRules)
	s.AutoSwapSlotName = pointer.From(appSiteSlotConfig.AutoSwapSlotName)
	s.ContainerRegistryMSI = pointer.From(appSiteSlotConfig.AcrUserManagedIdentityID)
	s.Cors = FlattenCorsSettings(appSiteSlotConfig.Cors)
	s.DetailedErrorLogging = pointer.From(appSiteSlotConfig.DetailedErrorLoggingEnabled)
	s.Http2Enabled = pointer.From(appSiteSlotConfig.HTTP20Enabled)
	s.IpRestriction = FlattenIpRestrictions(appSiteSlotConfig.IPSecurityRestrictions)
	s.ManagedPipelineMode = string(appSiteSlotConfig.ManagedPipelineMode)
	s.ScmType = string(appSiteSlotConfig.ScmType)
	s.FtpsState = string(appSiteSlotConfig.FtpsState)
	s.HealthCheckPath = pointer.From(appSiteSlotConfig.HealthCheckPath)
	s.LoadBalancing = string(appSiteSlotConfig.LoadBalancing)
	s.LocalMysql = pointer.From(appSiteSlotConfig.LocalMySQLEnabled)
	s.MinTlsVersion = string(appSiteSlotConfig.MinTLSVersion)
	s.WorkerCount = int(pointer.From(appSiteSlotConfig.NumberOfWorkers))
	s.RemoteDebugging = pointer.From(appSiteSlotConfig.RemoteDebuggingEnabled)
	s.RemoteDebuggingVersion = strings.ToUpper(pointer.From(appSiteSlotConfig.RemoteDebuggingVersion))
	s.ScmIpRestriction = FlattenIpRestrictions(appSiteSlotConfig.ScmIPSecurityRestrictions)
	s.ScmMinTlsVersion = string(appSiteSlotConfig.ScmMinTLSVersion)
	s.ScmUseMainIpRestriction = pointer.From(appSiteSlotConfig.ScmIPSecurityRestrictionsUseMain)
	s.Use32BitWorker = pointer.From(appSiteSlotConfig.Use32BitWorkerProcess)
	s.UseManagedIdentityACR = pointer.From(appSiteSlotConfig.AcrUseManagedIdentityCreds)
	s.WebSockets = pointer.From(appSiteSlotConfig.WebSocketsEnabled)
	s.VnetRouteAllEnabled = pointer.From(appSiteSlotConfig.VnetRouteAllEnabled)

	if appSiteSlotConfig.APIManagementConfig != nil && appSiteSlotConfig.APIManagementConfig.ID != nil {
		s.ApiManagementConfigId = *appSiteSlotConfig.APIManagementConfig.ID
	}

	if appSiteSlotConfig.APIDefinition != nil && appSiteSlotConfig.APIDefinition.URL != nil {
		s.ApiDefinition = *appSiteSlotConfig.APIDefinition.URL
	}

	if appSiteSlotConfig.DefaultDocuments != nil {
		s.DefaultDocuments = *appSiteSlotConfig.DefaultDocuments
	}

	if appSiteSlotConfig.LinuxFxVersion != nil {
		var linuxAppStack ApplicationStackLinux
		s.LinuxFxVersion = *appSiteSlotConfig.LinuxFxVersion
		// Decode the string to docker values
		linuxAppStack = decodeApplicationStackLinux(s.LinuxFxVersion)
		s.ApplicationStack = []ApplicationStackLinux{linuxAppStack}
	}

}

func (s *SiteConfigLinuxWebAppSlot) SetHealthCheckEvictionTime(input map[string]string) {
	if v, ok := input["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"]; ok && v != "" {
		// Discarding the error here as an invalid value should result in `0`
		s.HealthCheckEvictionTime, _ = strconv.Atoi(v)
	}
}

func (s *SiteConfigLinuxWebAppSlot) DecodeDockerAppStack(input map[string]string) {
	applicationStack := ApplicationStackLinux{}

	if v, ok := input["DOCKER_REGISTRY_SERVER_URL"]; ok {
		applicationStack.DockerRegistryUrl = v
	}

	if v, ok := input["DOCKER_REGISTRY_SERVER_USERNAME"]; ok {
		applicationStack.DockerRegistryUsername = v
	}

	if v, ok := input["DOCKER_REGISTRY_SERVER_PASSWORD"]; ok {
		applicationStack.DockerRegistryPassword = v
	}

	registryHost := trimURLScheme(applicationStack.DockerRegistryUrl)
	dockerString := strings.TrimPrefix(s.LinuxFxVersion, "DOCKER|")
	applicationStack.DockerImage = strings.TrimPrefix(dockerString, registryHost)

	s.ApplicationStack = []ApplicationStackLinux{applicationStack}
}

func FlattenSiteConfigLinuxWebAppSlot(appSiteSlotConfig *web.SiteConfig, healthCheckCount *int) []SiteConfigLinuxWebAppSlot {
	if appSiteSlotConfig == nil {
		return nil
	}

	siteConfig := SiteConfigLinuxWebAppSlot{
		AlwaysOn:                pointer.From(appSiteSlotConfig.AlwaysOn),
		AppCommandLine:          pointer.From(appSiteSlotConfig.AppCommandLine),
		AutoHeal:                pointer.From(appSiteSlotConfig.AutoHealEnabled),
		AutoHealSettings:        flattenAutoHealSettingsLinux(appSiteSlotConfig.AutoHealRules),
		AutoSwapSlotName:        pointer.From(appSiteSlotConfig.AutoSwapSlotName),
		ContainerRegistryMSI:    pointer.From(appSiteSlotConfig.AcrUserManagedIdentityID),
		Cors:                    FlattenCorsSettings(appSiteSlotConfig.Cors),
		DetailedErrorLogging:    pointer.From(appSiteSlotConfig.DetailedErrorLoggingEnabled),
		Http2Enabled:            pointer.From(appSiteSlotConfig.HTTP20Enabled),
		IpRestriction:           FlattenIpRestrictions(appSiteSlotConfig.IPSecurityRestrictions),
		ManagedPipelineMode:     string(appSiteSlotConfig.ManagedPipelineMode),
		ScmType:                 string(appSiteSlotConfig.ScmType),
		FtpsState:               string(appSiteSlotConfig.FtpsState),
		HealthCheckPath:         pointer.From(appSiteSlotConfig.HealthCheckPath),
		HealthCheckEvictionTime: pointer.From(healthCheckCount),
		LoadBalancing:           string(appSiteSlotConfig.LoadBalancing),
		LocalMysql:              pointer.From(appSiteSlotConfig.LocalMySQLEnabled),
		MinTlsVersion:           string(appSiteSlotConfig.MinTLSVersion),
		WorkerCount:             int(pointer.From(appSiteSlotConfig.NumberOfWorkers)),
		RemoteDebugging:         pointer.From(appSiteSlotConfig.RemoteDebuggingEnabled),
		RemoteDebuggingVersion:  strings.ToUpper(pointer.From(appSiteSlotConfig.RemoteDebuggingVersion)),
		ScmIpRestriction:        FlattenIpRestrictions(appSiteSlotConfig.ScmIPSecurityRestrictions),
		ScmMinTlsVersion:        string(appSiteSlotConfig.ScmMinTLSVersion),
		ScmUseMainIpRestriction: pointer.From(appSiteSlotConfig.ScmIPSecurityRestrictionsUseMain),
		Use32BitWorker:          pointer.From(appSiteSlotConfig.Use32BitWorkerProcess),
		UseManagedIdentityACR:   pointer.From(appSiteSlotConfig.AcrUseManagedIdentityCreds),
		WebSockets:              pointer.From(appSiteSlotConfig.WebSocketsEnabled),
		VnetRouteAllEnabled:     pointer.From(appSiteSlotConfig.VnetRouteAllEnabled),
	}

	if appSiteSlotConfig.APIManagementConfig != nil && appSiteSlotConfig.APIManagementConfig.ID != nil {
		siteConfig.ApiManagementConfigId = *appSiteSlotConfig.APIManagementConfig.ID
	}

	if appSiteSlotConfig.APIDefinition != nil && appSiteSlotConfig.APIDefinition.URL != nil {
		siteConfig.ApiDefinition = *appSiteSlotConfig.APIDefinition.URL
	}

	if appSiteSlotConfig.DefaultDocuments != nil {
		siteConfig.DefaultDocuments = *appSiteSlotConfig.DefaultDocuments
	}

	if appSiteSlotConfig.LinuxFxVersion != nil {
		var linuxAppStack ApplicationStackLinux
		siteConfig.LinuxFxVersion = *appSiteSlotConfig.LinuxFxVersion
		// Decode the string to docker values
		linuxAppStack = decodeApplicationStackLinux(siteConfig.LinuxFxVersion)
		siteConfig.ApplicationStack = []ApplicationStackLinux{linuxAppStack}
	}

	return []SiteConfigLinuxWebAppSlot{siteConfig}
}

func ExpandSiteConfigWindowsWebAppSlot(siteConfig []SiteConfigWindowsWebAppSlot, existing *web.SiteConfig, metadata sdk.ResourceMetaData) (*web.SiteConfig, *string, error) {
	if len(siteConfig) == 0 {
		return nil, nil, nil
	}

	expanded := &web.SiteConfig{}
	if existing != nil {
		expanded = existing
	}

	winSlotSiteConfig := siteConfig[0]

	currentStack := ""

	if len(winSlotSiteConfig.ApplicationStack) == 1 {
		winAppStack := winSlotSiteConfig.ApplicationStack[0]
		currentStack = winAppStack.CurrentStack
	}

	if metadata.ResourceData.HasChange("site_config.0.always_on") {
		expanded.AlwaysOn = pointer.To(winSlotSiteConfig.AlwaysOn)
	}

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.APIManagementConfig = &web.APIManagementConfig{
			ID: pointer.To(winSlotSiteConfig.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.APIDefinition = &web.APIDefinitionInfo{
			URL: pointer.To(winSlotSiteConfig.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = pointer.To(winSlotSiteConfig.AppCommandLine)
	}

	if metadata.ResourceData.HasChange("site_config.0.application_stack") {
		if len(winSlotSiteConfig.ApplicationStack) == 1 {
			winAppStack := winSlotSiteConfig.ApplicationStack[0]

			if winAppStack.NetFrameworkVersion != "" {
				expanded.NetFrameworkVersion = pointer.To(winAppStack.NetFrameworkVersion)
				if currentStack == "" {
					currentStack = CurrentStackDotNet
				}
			}
			if winAppStack.NetCoreVersion != "" {
				expanded.NetFrameworkVersion = pointer.To(winAppStack.NetCoreVersion)
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
				switch {
				case winAppStack.JavaEmbeddedServer:
					expanded.JavaContainer = pointer.To(JavaContainerEmbeddedServer)
					expanded.JavaContainerVersion = pointer.To(JavaContainerEmbeddedServerVersion)

				case winAppStack.TomcatVersion != "":
					expanded.JavaContainer = pointer.To(JavaContainerTomcat)
					expanded.JavaContainerVersion = pointer.To(winAppStack.TomcatVersion)

				case winAppStack.JavaContainer != "":
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
		expanded.VirtualApplications = expandVirtualApplicationsForUpdate(winSlotSiteConfig.VirtualApplications)
	} else {
		expanded.VirtualApplications = expandVirtualApplications(winSlotSiteConfig.VirtualApplications)
	}

	if metadata.ResourceData.HasChange("site_config.0.container_registry_use_managed_identity") {
		expanded.AcrUseManagedIdentityCreds = pointer.To(winSlotSiteConfig.UseManagedIdentityACR)
	}

	if metadata.ResourceData.HasChange("site_config.0.container_registry_managed_identity_client_id") {
		expanded.AcrUserManagedIdentityID = pointer.To(winSlotSiteConfig.ContainerRegistryUserMSI)
	}

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = &winSlotSiteConfig.DefaultDocuments
	}

	if metadata.ResourceData.HasChange("site_config.0.http2_enabled") {
		expanded.HTTP20Enabled = pointer.To(winSlotSiteConfig.Http2Enabled)
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := ExpandIpRestrictions(winSlotSiteConfig.IpRestriction)
		if err != nil {
			return nil, nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_use_main_ip_restriction") {
		expanded.ScmIPSecurityRestrictionsUseMain = pointer.To(winSlotSiteConfig.ScmUseMainIpRestriction)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := ExpandIpRestrictions(winSlotSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.local_mysql_enabled") {
		expanded.LocalMySQLEnabled = pointer.To(winSlotSiteConfig.LocalMysql)
	}

	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = web.SiteLoadBalancing(winSlotSiteConfig.LoadBalancing)
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = web.ManagedPipelineMode(winSlotSiteConfig.ManagedPipelineMode)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_enabled") {
		expanded.RemoteDebuggingEnabled = pointer.To(winSlotSiteConfig.RemoteDebugging)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = pointer.To(winSlotSiteConfig.RemoteDebuggingVersion)
	}

	expanded.Use32BitWorkerProcess = pointer.To(winSlotSiteConfig.Use32BitWorker)

	if metadata.ResourceData.HasChange("site_config.0.websockets_enabled") {
		expanded.WebSocketsEnabled = pointer.To(winSlotSiteConfig.WebSockets)
	}

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		expanded.FtpsState = web.FtpsState(winSlotSiteConfig.FtpsState)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = pointer.To(winSlotSiteConfig.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.worker_count") {
		expanded.NumberOfWorkers = pointer.To(int32(winSlotSiteConfig.WorkerCount))
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTLSVersion = web.SupportedTLSVersions(winSlotSiteConfig.MinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTLSVersion = web.SupportedTLSVersions(winSlotSiteConfig.ScmMinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_swap_slot_name") {
		expanded.AutoSwapSlotName = pointer.To(winSlotSiteConfig.AutoSwapSlotName)
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		cors := ExpandCorsSettings(winSlotSiteConfig.Cors)
		if cors == nil {
			cors = &web.CorsSettings{
				AllowedOrigins: &[]string{},
			}
		}
		expanded.Cors = cors
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_heal_enabled") {
		expanded.AutoHealEnabled = pointer.To(winSlotSiteConfig.AutoHeal)
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_heal_setting") {
		expanded.AutoHealRules = expandAutoHealSettingsWindows(winSlotSiteConfig.AutoHealSettings)
	}

	if metadata.ResourceData.HasChange("site_config.0.vnet_route_all_enabled") {
		expanded.VnetRouteAllEnabled = pointer.To(winSlotSiteConfig.VnetRouteAllEnabled)
	}

	return expanded, &currentStack, nil
}

func (s *SiteConfigWindowsWebAppSlot) Flatten(appSiteSlotConfig *web.SiteConfig, currentStack string) {
	if appSiteSlotConfig == nil {
		return
	}

	s.AlwaysOn = pointer.From(appSiteSlotConfig.AlwaysOn)
	s.AppCommandLine = pointer.From(appSiteSlotConfig.AppCommandLine)
	s.AutoHeal = pointer.From(appSiteSlotConfig.AutoHealEnabled)
	s.AutoHealSettings = flattenAutoHealSettingsWindows(appSiteSlotConfig.AutoHealRules)
	s.AutoSwapSlotName = pointer.From(appSiteSlotConfig.AutoSwapSlotName)
	s.ContainerRegistryUserMSI = pointer.From(appSiteSlotConfig.AcrUserManagedIdentityID)
	s.Cors = FlattenCorsSettings(appSiteSlotConfig.Cors)
	s.DetailedErrorLogging = pointer.From(appSiteSlotConfig.DetailedErrorLoggingEnabled)
	s.FtpsState = string(appSiteSlotConfig.FtpsState)
	s.HealthCheckPath = pointer.From(appSiteSlotConfig.HealthCheckPath)
	s.Http2Enabled = pointer.From(appSiteSlotConfig.HTTP20Enabled)
	s.IpRestriction = FlattenIpRestrictions(appSiteSlotConfig.IPSecurityRestrictions)
	s.LoadBalancing = string(appSiteSlotConfig.LoadBalancing)
	s.LocalMysql = pointer.From(appSiteSlotConfig.LocalMySQLEnabled)
	s.ManagedPipelineMode = string(appSiteSlotConfig.ManagedPipelineMode)
	s.MinTlsVersion = string(appSiteSlotConfig.MinTLSVersion)
	s.WorkerCount = int(pointer.From(appSiteSlotConfig.NumberOfWorkers))
	s.RemoteDebugging = pointer.From(appSiteSlotConfig.RemoteDebuggingEnabled)
	s.RemoteDebuggingVersion = strings.ToUpper(pointer.From(appSiteSlotConfig.RemoteDebuggingVersion))
	s.ScmIpRestriction = FlattenIpRestrictions(appSiteSlotConfig.ScmIPSecurityRestrictions)
	s.ScmMinTlsVersion = string(appSiteSlotConfig.ScmMinTLSVersion)
	s.ScmType = string(appSiteSlotConfig.ScmType)
	s.ScmUseMainIpRestriction = pointer.From(appSiteSlotConfig.ScmIPSecurityRestrictionsUseMain)
	s.Use32BitWorker = pointer.From(appSiteSlotConfig.Use32BitWorkerProcess)
	s.UseManagedIdentityACR = pointer.From(appSiteSlotConfig.AcrUseManagedIdentityCreds)
	s.VirtualApplications = flattenVirtualApplications(appSiteSlotConfig.VirtualApplications)
	s.WebSockets = pointer.From(appSiteSlotConfig.WebSocketsEnabled)
	s.VnetRouteAllEnabled = pointer.From(appSiteSlotConfig.VnetRouteAllEnabled)

	if appSiteSlotConfig.APIManagementConfig != nil && appSiteSlotConfig.APIManagementConfig.ID != nil {
		s.ApiManagementConfigId = *appSiteSlotConfig.APIManagementConfig.ID
	}

	if appSiteSlotConfig.APIDefinition != nil && appSiteSlotConfig.APIDefinition.URL != nil {
		s.ApiDefinition = *appSiteSlotConfig.APIDefinition.URL
	}

	if appSiteSlotConfig.DefaultDocuments != nil {
		s.DefaultDocuments = *appSiteSlotConfig.DefaultDocuments
	}

	if appSiteSlotConfig.NumberOfWorkers != nil {
		s.WorkerCount = int(*appSiteSlotConfig.NumberOfWorkers)
	}

	winAppStack := ApplicationStackWindows{}

	winAppStack.NetFrameworkVersion = pointer.From(appSiteSlotConfig.NetFrameworkVersion)
	if currentStack == CurrentStackDotNetCore {
		winAppStack.NetCoreVersion = pointer.From(appSiteSlotConfig.NetFrameworkVersion)
	}
	winAppStack.PhpVersion = pointer.From(appSiteSlotConfig.PhpVersion)
	if winAppStack.PhpVersion == "" {
		winAppStack.PhpVersion = PhpVersionOff
	}
	winAppStack.PythonVersion = pointer.From(appSiteSlotConfig.PythonVersion) // This _should_ always be `""`
	winAppStack.Python = currentStack == CurrentStackPython
	winAppStack.JavaVersion = pointer.From(appSiteSlotConfig.JavaVersion)
	switch pointer.From(appSiteSlotConfig.JavaContainer) {
	case JavaContainerTomcat:
		winAppStack.TomcatVersion = *appSiteSlotConfig.JavaContainerVersion
	case JavaContainerEmbeddedServer:
		winAppStack.JavaEmbeddedServer = true
	}

	s.WindowsFxVersion = pointer.From(appSiteSlotConfig.WindowsFxVersion)
	if s.WindowsFxVersion != "" {
		// Decode the string to docker values
		parts := strings.Split(strings.TrimPrefix(s.WindowsFxVersion, "DOCKER|"), ":")
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

	s.ApplicationStack = []ApplicationStackWindows{winAppStack}

}

func (s *SiteConfigWindowsWebAppSlot) SetHealthCheckEvictionTime(input map[string]string) {
	if v, ok := input["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"]; ok && v != "" {
		// Discarding the error here as an invalid value should result in `0`
		s.HealthCheckEvictionTime, _ = strconv.Atoi(v)
	}
}

func (s *SiteConfigWindowsWebAppSlot) ParseNodeVersion(input map[string]string) map[string]string {
	if nodeVer, ok := input["WEBSITE_NODE_DEFAULT_VERSION"]; ok {
		if s.ApplicationStack == nil {
			s.ApplicationStack = make([]ApplicationStackWindows, 0)
		}
		s.ApplicationStack[0].NodeVersion = nodeVer
		delete(input, "WEBSITE_NODE_DEFAULT_VERSION")
	}

	return input
}

func (s *SiteConfigWindowsWebAppSlot) DecodeDockerAppStack(input map[string]string) {
	applicationStack := ApplicationStackWindows{}

	if v, ok := input["DOCKER_REGISTRY_SERVER_URL"]; ok {
		applicationStack.DockerRegistryUrl = v
	}

	if v, ok := input["DOCKER_REGISTRY_SERVER_USERNAME"]; ok {
		applicationStack.DockerRegistryUsername = v
	}

	if v, ok := input["DOCKER_REGISTRY_SERVER_PASSWORD"]; ok {
		applicationStack.DockerRegistryPassword = v
	}

	registryHost := trimURLScheme(applicationStack.DockerRegistryUrl)
	dockerString := strings.TrimPrefix(s.WindowsFxVersion, "DOCKER|")
	applicationStack.DockerImageName = strings.TrimPrefix(dockerString, registryHost)

	s.ApplicationStack = []ApplicationStackWindows{applicationStack}
}
