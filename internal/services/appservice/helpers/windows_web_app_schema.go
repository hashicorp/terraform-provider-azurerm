// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
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
					ValidateFunc: validate.ApiID,
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

func (s *SiteConfigWindows) ExpandForCreate(appSettings map[string]string) (*web.SiteConfig, error) {
	expanded := &web.SiteConfig{}

	expanded.AlwaysOn = pointer.To(s.AlwaysOn)
	expanded.AcrUseManagedIdentityCreds = pointer.To(s.UseManagedIdentityACR)
	expanded.AutoHealEnabled = pointer.To(s.AutoHeal)
	expanded.FtpsState = web.FtpsState(s.FtpsState)
	expanded.HTTP20Enabled = pointer.To(s.Http2Enabled)
	expanded.LoadBalancing = web.SiteLoadBalancing(s.LoadBalancing)
	expanded.LocalMySQLEnabled = pointer.To(s.LocalMysql)
	expanded.ManagedPipelineMode = web.ManagedPipelineMode(s.ManagedPipelineMode)
	expanded.MinTLSVersion = web.SupportedTLSVersions(s.MinTlsVersion)
	expanded.RemoteDebuggingEnabled = pointer.To(s.RemoteDebugging)
	expanded.ScmIPSecurityRestrictionsUseMain = pointer.To(s.ScmUseMainIpRestriction)
	expanded.ScmMinTLSVersion = web.SupportedTLSVersions(s.ScmMinTlsVersion)
	expanded.Use32BitWorkerProcess = pointer.To(s.Use32BitWorker)
	expanded.WebSocketsEnabled = pointer.To(s.WebSockets)
	expanded.VirtualApplications = expandVirtualApplications(s.VirtualApplications)
	expanded.VnetRouteAllEnabled = pointer.To(s.VnetRouteAllEnabled)

	if s.ApiManagementConfigId != "" {
		expanded.APIManagementConfig = &web.APIManagementConfig{
			ID: pointer.To(s.ApiManagementConfigId),
		}
	}

	if s.ApiDefinition != "" {
		expanded.APIDefinition = &web.APIDefinitionInfo{
			URL: pointer.To(s.ApiDefinition),
		}
	}

	if s.AppCommandLine != "" {
		expanded.AppCommandLine = pointer.To(s.AppCommandLine)
	}

	expanded.AppSettings = ExpandAppSettingsForCreate(appSettings)

	if len(s.ApplicationStack) == 1 {
		winAppStack := s.ApplicationStack[0]
		if winAppStack.NetFrameworkVersion != "" {
			expanded.NetFrameworkVersion = pointer.To(winAppStack.NetFrameworkVersion)
		}
		if winAppStack.NetCoreVersion != "" {
			expanded.NetFrameworkVersion = pointer.To(winAppStack.NetCoreVersion)
		}
		if winAppStack.PhpVersion != "" {
			if winAppStack.PhpVersion != PhpVersionOff {
				expanded.PhpVersion = pointer.To(winAppStack.PhpVersion)
			} else {
				expanded.PhpVersion = pointer.To("")
			}
		}
		if winAppStack.PythonVersion != "" || winAppStack.Python {
			expanded.PythonVersion = pointer.To(winAppStack.PythonVersion)
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
		}
		if !features.FourPointOhBeta() {
			if winAppStack.DockerContainerName != "" || winAppStack.DockerContainerRegistry != "" || winAppStack.DockerContainerTag != "" {
				if winAppStack.DockerContainerRegistry != "" {
					expanded.WindowsFxVersion = pointer.To(fmt.Sprintf("DOCKER|%s/%s:%s", winAppStack.DockerContainerRegistry, winAppStack.DockerContainerName, winAppStack.DockerContainerTag))
				} else {
					expanded.WindowsFxVersion = pointer.To(fmt.Sprintf("DOCKER|%s:%s", winAppStack.DockerContainerName, winAppStack.DockerContainerTag))
				}
			}
		}

		if winAppStack.DockerImageName != "" {
			expanded.WindowsFxVersion = pointer.To(EncodeDockerFxStringWindows(winAppStack.DockerImageName, winAppStack.DockerRegistryUrl))
			appSettings["DOCKER_REGISTRY_SERVER_URL"] = winAppStack.DockerRegistryUrl
			appSettings["DOCKER_REGISTRY_SERVER_USERNAME"] = winAppStack.DockerRegistryUsername
			appSettings["DOCKER_REGISTRY_SERVER_PASSWORD"] = winAppStack.DockerRegistryPassword

		}

	} else {
		expanded.WindowsFxVersion = pointer.To("")
	}

	if s.ContainerRegistryUserMSI != "" {
		expanded.AcrUserManagedIdentityID = pointer.To(s.ContainerRegistryUserMSI)
	}

	if len(s.DefaultDocuments) != 0 {
		expanded.DefaultDocuments = pointer.To(s.DefaultDocuments)
	}

	if len(s.IpRestriction) != 0 {
		ipRestrictions, err := ExpandIpRestrictions(s.IpRestriction)
		if err != nil {
			return nil, fmt.Errorf("expanding IP Restrictions: %+v", err)
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	if len(s.ScmIpRestriction) != 0 {
		scmIpRestrictions, err := ExpandIpRestrictions(s.ScmIpRestriction)
		if err != nil {
			return nil, fmt.Errorf("expanding SCM IP Restrictions: %+v", err)
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	if s.RemoteDebuggingVersion != "" {
		expanded.RemoteDebuggingVersion = pointer.To(s.RemoteDebuggingVersion)
	}

	if s.HealthCheckPath != "" {
		expanded.HealthCheckPath = pointer.To(s.HealthCheckPath)
	}

	if s.WorkerCount != 0 {
		expanded.NumberOfWorkers = pointer.To(int32(s.WorkerCount))
	}

	if len(s.Cors) != 0 {
		expanded.Cors = ExpandCorsSettings(s.Cors)
	}

	if len(s.AutoHealSettings) != 0 {
		expanded.AutoHealRules = expandAutoHealSettingsWindows(s.AutoHealSettings)
	}

	return expanded, nil
}

func (s *SiteConfigWindows) ExpandForUpdate(metadata sdk.ResourceMetaData, existing *web.SiteConfig, appSettings map[string]string) (*web.SiteConfig, error) {
	expanded := web.SiteConfig{}
	if existing != nil {
		expanded = *existing
	}

	expanded.AlwaysOn = pointer.To(s.AlwaysOn)
	expanded.AcrUseManagedIdentityCreds = pointer.To(s.UseManagedIdentityACR)
	expanded.AutoHealEnabled = pointer.To(s.AutoHeal)
	expanded.HTTP20Enabled = pointer.To(s.Http2Enabled)
	expanded.ScmIPSecurityRestrictionsUseMain = pointer.To(s.ScmUseMainIpRestriction)
	expanded.LocalMySQLEnabled = pointer.To(s.LocalMysql)
	expanded.RemoteDebuggingEnabled = pointer.To(s.RemoteDebugging)
	expanded.Use32BitWorkerProcess = pointer.To(s.Use32BitWorker)
	expanded.WebSocketsEnabled = pointer.To(s.WebSockets)

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.APIManagementConfig = &web.APIManagementConfig{
			ID: pointer.To(s.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.APIDefinition = &web.APIDefinitionInfo{
			URL: pointer.To(s.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = pointer.To(s.AppCommandLine)
	}

	if metadata.ResourceData.HasChange("site_config.0.application_stack") {
		if len(s.ApplicationStack) == 1 {
			winAppStack := s.ApplicationStack[0]
			if winAppStack.NetFrameworkVersion != "" {
				expanded.NetFrameworkVersion = pointer.To(winAppStack.NetFrameworkVersion)
			}
			if winAppStack.NetCoreVersion != "" {
				expanded.NetFrameworkVersion = pointer.To(winAppStack.NetCoreVersion)
			}
			if winAppStack.PhpVersion != "" {
				if winAppStack.PhpVersion != PhpVersionOff {
					expanded.PhpVersion = pointer.To(winAppStack.PhpVersion)
				} else {
					expanded.PhpVersion = pointer.To("")
				}
			}
			if winAppStack.PythonVersion != "" || winAppStack.Python {
				expanded.PythonVersion = pointer.To(winAppStack.PythonVersion)
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
			}
			if !features.FourPointOhBeta() {
				if winAppStack.DockerContainerName != "" || winAppStack.DockerContainerRegistry != "" || winAppStack.DockerContainerTag != "" {
					if winAppStack.DockerContainerRegistry != "" {
						expanded.WindowsFxVersion = pointer.To(fmt.Sprintf("DOCKER|%s/%s:%s", winAppStack.DockerContainerRegistry, winAppStack.DockerContainerName, winAppStack.DockerContainerTag))
					} else {
						expanded.WindowsFxVersion = pointer.To(fmt.Sprintf("DOCKER|%s:%s", winAppStack.DockerContainerName, winAppStack.DockerContainerTag))
					}
				}
			}

			if winAppStack.DockerImageName != "" {
				expanded.WindowsFxVersion = pointer.To(EncodeDockerFxStringWindows(winAppStack.DockerImageName, winAppStack.DockerRegistryUrl))
				appSettings["DOCKER_REGISTRY_SERVER_URL"] = winAppStack.DockerRegistryUrl
				appSettings["DOCKER_REGISTRY_SERVER_USERNAME"] = winAppStack.DockerRegistryUsername
				appSettings["DOCKER_REGISTRY_SERVER_PASSWORD"] = winAppStack.DockerRegistryPassword

			}

		} else {
			expanded.WindowsFxVersion = pointer.To("")
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.virtual_application") {
		expanded.VirtualApplications = expandVirtualApplicationsForUpdate(s.VirtualApplications)
	} else {
		expanded.VirtualApplications = expandVirtualApplications(s.VirtualApplications)
	}

	if metadata.ResourceData.HasChange("site_config.0.container_registry_managed_identity_client_id") {
		expanded.AcrUserManagedIdentityID = pointer.To(s.ContainerRegistryUserMSI)
	}

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = pointer.To(s.DefaultDocuments)
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := ExpandIpRestrictions(s.IpRestriction)
		if err != nil {
			return nil, err
		}

		expanded.IPSecurityRestrictions = ipRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := ExpandIpRestrictions(s.ScmIpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = web.SiteLoadBalancing(s.LoadBalancing)
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = web.ManagedPipelineMode(s.ManagedPipelineMode)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = pointer.To(s.RemoteDebuggingVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		expanded.FtpsState = web.FtpsState(s.FtpsState)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = pointer.To(s.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.worker_count") {
		expanded.NumberOfWorkers = pointer.To(int32(s.WorkerCount))
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTLSVersion = web.SupportedTLSVersions(s.MinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTLSVersion = web.SupportedTLSVersions(s.ScmMinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		cors := ExpandCorsSettings(s.Cors)
		if cors == nil {
			cors = &web.CorsSettings{
				AllowedOrigins: &[]string{},
			}
		}
		expanded.Cors = cors
	}
	if metadata.ResourceData.HasChange("site_config.0.auto_heal_setting") {
		expanded.AutoHealRules = expandAutoHealSettingsWindows(s.AutoHealSettings)
	}

	if metadata.ResourceData.HasChange("site_config.0.vnet_route_all_enabled") {
		expanded.VnetRouteAllEnabled = pointer.To(s.VnetRouteAllEnabled)
	}

	return &expanded, nil
}

func (s *SiteConfigWindows) Flatten(appSiteConfig *web.SiteConfig, currentStack string) error {
	if appSiteConfig != nil {
		s.AlwaysOn = pointer.From(appSiteConfig.AlwaysOn)
		s.AppCommandLine = pointer.From(appSiteConfig.AppCommandLine)
		s.AutoHeal = pointer.From(appSiteConfig.AutoHealEnabled)
		s.AutoHealSettings = flattenAutoHealSettingsWindows(appSiteConfig.AutoHealRules)
		s.ContainerRegistryUserMSI = pointer.From(appSiteConfig.AcrUserManagedIdentityID)
		s.Cors = FlattenCorsSettings(appSiteConfig.Cors)
		s.DetailedErrorLogging = pointer.From(appSiteConfig.DetailedErrorLoggingEnabled)
		s.FtpsState = string(appSiteConfig.FtpsState)
		s.HealthCheckPath = pointer.From(appSiteConfig.HealthCheckPath)
		s.Http2Enabled = pointer.From(appSiteConfig.HTTP20Enabled)
		s.IpRestriction = FlattenIpRestrictions(appSiteConfig.IPSecurityRestrictions)
		s.LoadBalancing = string(appSiteConfig.LoadBalancing)
		s.LocalMysql = pointer.From(appSiteConfig.LocalMySQLEnabled)
		s.ManagedPipelineMode = string(appSiteConfig.ManagedPipelineMode)
		s.MinTlsVersion = string(appSiteConfig.MinTLSVersion)
		s.WorkerCount = int(pointer.From(appSiteConfig.NumberOfWorkers))
		s.RemoteDebugging = pointer.From(appSiteConfig.RemoteDebuggingEnabled)
		s.RemoteDebuggingVersion = strings.ToUpper(pointer.From(appSiteConfig.RemoteDebuggingVersion))
		s.ScmIpRestriction = FlattenIpRestrictions(appSiteConfig.ScmIPSecurityRestrictions)
		s.ScmMinTlsVersion = string(appSiteConfig.ScmMinTLSVersion)
		s.ScmType = string(appSiteConfig.ScmType)
		s.ScmUseMainIpRestriction = pointer.From(appSiteConfig.ScmIPSecurityRestrictionsUseMain)
		s.Use32BitWorker = pointer.From(appSiteConfig.Use32BitWorkerProcess)
		s.UseManagedIdentityACR = pointer.From(appSiteConfig.AcrUseManagedIdentityCreds)
		s.VirtualApplications = flattenVirtualApplications(appSiteConfig.VirtualApplications)
		s.WebSockets = pointer.From(appSiteConfig.WebSocketsEnabled)
		s.VnetRouteAllEnabled = pointer.From(appSiteConfig.VnetRouteAllEnabled)
	}

	if appSiteConfig.APIManagementConfig != nil && appSiteConfig.APIManagementConfig.ID != nil {
		apiId, err := parse.ApiIDInsensitively(*appSiteConfig.APIManagementConfig.ID)
		if err != nil {
			return fmt.Errorf("could not parse API Management ID: %+v", err)
		}
		s.ApiManagementConfigId = apiId.ID()
	}

	if appSiteConfig.APIDefinition != nil && appSiteConfig.APIDefinition.URL != nil {
		s.ApiDefinition = *appSiteConfig.APIDefinition.URL
	}

	if appSiteConfig.DefaultDocuments != nil {
		s.DefaultDocuments = *appSiteConfig.DefaultDocuments
	}

	if appSiteConfig.NumberOfWorkers != nil {
		s.WorkerCount = int(*appSiteConfig.NumberOfWorkers)
	}

	var winAppStack ApplicationStackWindows
	if currentStack == CurrentStackDotNetCore {
		winAppStack.NetCoreVersion = pointer.From(appSiteConfig.NetFrameworkVersion)
	}
	if currentStack == CurrentStackDotNet {
		winAppStack.NetFrameworkVersion = pointer.From(appSiteConfig.NetFrameworkVersion)
	}
	winAppStack.PhpVersion = pointer.From(appSiteConfig.PhpVersion)
	if winAppStack.PhpVersion == "" {
		winAppStack.PhpVersion = PhpVersionOff
	}
	winAppStack.PythonVersion = pointer.From(appSiteConfig.PythonVersion) // This _should_ always be `""`
	winAppStack.Python = currentStack == CurrentStackPython
	winAppStack.JavaVersion = pointer.From(appSiteConfig.JavaVersion)
	switch pointer.From(appSiteConfig.JavaContainer) {
	case JavaContainerTomcat:
		winAppStack.TomcatVersion = *appSiteConfig.JavaContainerVersion
	case JavaContainerEmbeddedServer:
		winAppStack.JavaEmbeddedServer = true
	}

	s.WindowsFxVersion = pointer.From(appSiteConfig.WindowsFxVersion)

	winAppStack.CurrentStack = currentStack

	s.ApplicationStack = []ApplicationStackWindows{winAppStack}

	return nil
}

func (s *SiteConfigWindows) SetHealthCheckEvictionTime(input map[string]string) {
	if v, ok := input["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"]; ok && v != "" {
		// Discarding the error here as an invalid value should result in `0`
		s.HealthCheckEvictionTime, _ = strconv.Atoi(v)
	}
}

func (s *SiteConfigWindows) DecodeDockerAppStack(input map[string]string) {
	applicationStack := ApplicationStackWindows{}
	if len(s.ApplicationStack) == 1 {
		applicationStack = s.ApplicationStack[0]
	}

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

func (s *SiteConfigWindows) DecodeDockerDeprecatedAppStack(input map[string]string, usesDeprecated bool) {
	applicationStack := ApplicationStackWindows{}
	if len(s.ApplicationStack) == 1 {
		applicationStack = s.ApplicationStack[0]
	}

	if !usesDeprecated {
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
		applicationStack.DockerImageName = strings.TrimPrefix(dockerString, registryHost+"/")
	} else {
		parts := strings.Split(strings.TrimPrefix(s.WindowsFxVersion, "DOCKER|"), ":")
		if len(parts) == 2 {
			applicationStack.DockerContainerTag = parts[1]
			path := strings.Split(parts[0], "/")
			if len(path) > 1 {
				applicationStack.DockerContainerRegistry = path[0]
				applicationStack.DockerContainerName = strings.TrimPrefix(parts[0], fmt.Sprintf("%s/", path[0]))
			} else {
				applicationStack.DockerContainerName = path[0]
			}
		}
	}

	s.ApplicationStack = []ApplicationStackWindows{applicationStack}
}

func (s *SiteConfigWindows) ParseNodeVersion(input map[string]string) map[string]string {
	if nodeVer, ok := input["WEBSITE_NODE_DEFAULT_VERSION"]; ok {
		if s.ApplicationStack == nil {
			s.ApplicationStack = make([]ApplicationStackWindows, 0)
			s.ApplicationStack[0] = ApplicationStackWindows{}
		}
		s.ApplicationStack[0].NodeVersion = nodeVer
		delete(input, "WEBSITE_NODE_DEFAULT_VERSION")
	}

	return input
}
