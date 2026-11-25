// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package linux_web_app_slot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SiteConfigLinuxWebAppSlot struct {
	AlwaysOn                      bool                            `tfschema:"always_on"`
	ApiManagementConfigId         string                          `tfschema:"api_management_api_id"`
	ApiDefinition                 string                          `tfschema:"api_definition_url"`
	AppCommandLine                string                          `tfschema:"app_command_line"`
	AutoHealSettings              []helpers.AutoHealSettingLinux  `tfschema:"auto_heal_setting"`
	AutoSwapSlotName              string                          `tfschema:"auto_swap_slot_name"`
	UseManagedIdentityACR         bool                            `tfschema:"container_registry_use_managed_identity"`
	ContainerRegistryMSI          string                          `tfschema:"container_registry_managed_identity_client_id"`
	DefaultDocuments              []string                        `tfschema:"default_documents"`
	Http2Enabled                  bool                            `tfschema:"http2_enabled"`
	IpRestriction                 []helpers.IpRestriction         `tfschema:"ip_restriction"`
	IpRestrictionDefaultAction    string                          `tfschema:"ip_restriction_default_action"`
	ScmUseMainIpRestriction       bool                            `tfschema:"scm_use_main_ip_restriction"`
	ScmIpRestriction              []helpers.IpRestriction         `tfschema:"scm_ip_restriction"`
	ScmIpRestrictionDefaultAction string                          `tfschema:"scm_ip_restriction_default_action"`
	LoadBalancing                 string                          `tfschema:"load_balancing_mode"`
	LocalMysql                    bool                            `tfschema:"local_mysql_enabled"`
	ManagedPipelineMode           string                          `tfschema:"managed_pipeline_mode"`
	RemoteDebugging               bool                            `tfschema:"remote_debugging_enabled"`
	RemoteDebuggingVersion        string                          `tfschema:"remote_debugging_version"`
	ScmType                       string                          `tfschema:"scm_type"`
	Use32BitWorker                bool                            `tfschema:"use_32_bit_worker"`
	WebSockets                    bool                            `tfschema:"websockets_enabled"`
	FtpsState                     string                          `tfschema:"ftps_state"`
	HealthCheckPath               string                          `tfschema:"health_check_path"`
	HealthCheckEvictionTime       int64                           `tfschema:"health_check_eviction_time_in_min"`
	WorkerCount                   int64                           `tfschema:"worker_count"`
	ApplicationStack              []helpers.ApplicationStackLinux `tfschema:"application_stack"`
	MinTlsVersion                 string                          `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion              string                          `tfschema:"scm_minimum_tls_version"`
	MinTlsCipherSuite             string                          `tfschema:"minimum_tls_cipher_suite"`
	Cors                          []helpers.CorsSetting           `tfschema:"cors"`
	DetailedErrorLogging          bool                            `tfschema:"detailed_error_logging_enabled"`
	LinuxFxVersion                string                          `tfschema:"linux_fx_version"`
	VnetRouteAllEnabled           bool                            `tfschema:"vnet_route_all_enabled"`
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
					ValidateFunc: validation.AsGeneratedID(api.ParseApiIDInsensitively),
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

				"application_stack": helpers.LinuxApplicationStackSchema(),

				"auto_heal_setting": helpers.AutoHealSettingSchemaLinux(),

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
					Computed: true, // azignore:AZS007 - pre-existing violation
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

				"ip_restriction": helpers.IpRestrictionSchema(),

				"ip_restriction_default_action": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      webapps.DefaultActionAllow,
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForDefaultAction(), false),
				},

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"scm_ip_restriction": helpers.IpRestrictionSchema(),

				"scm_ip_restriction_default_action": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      webapps.DefaultActionAllow,
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForDefaultAction(), false),
				},

				"local_mysql_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"load_balancing_mode": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      "LeastRequests",
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForSiteLoadBalancing(), false),
				},

				"managed_pipeline_mode": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(webapps.ManagedPipelineModeIntegrated),
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForManagedPipelineMode(), false),
				},

				"remote_debugging_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"remote_debugging_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true, // azignore:AZS007 - pre-existing violation
					ValidateFunc: validation.StringInSlice([]string{
						"VS2022",
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
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(webapps.FtpsStateDisabled),
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForFtpsState(), false),
				},

				"health_check_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					RequiredWith: []string{"site_config.0.health_check_eviction_time_in_min"},
				},

				"health_check_eviction_time_in_min": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(2, 10),
					RequiredWith: []string{"site_config.0.health_check_path"},
					Description:  "The amount of time in minutes that a node is unhealthy before being removed from the load balancer.",
				},

				"worker_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true, // azignore:AZS007 - pre-existing violation
					ValidateFunc: validation.IntBetween(1, 100),
				},

				"minimum_tls_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(webapps.SupportedTlsVersionsOnePointTwo),
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForSupportedTlsVersions(), false),
				},

				"scm_minimum_tls_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(webapps.SupportedTlsVersionsOnePointTwo),
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForSupportedTlsVersions(), false),
				},

				"minimum_tls_cipher_suite": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForTlsCipherSuites(), false),
					Description:  "Configures the minimum TLS cipher suite for the incoming requests to the Site.",
				},

				"cors": helpers.CorsSettingsSchema(),

				"auto_swap_slot_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"vnet_route_all_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Whether all outbound traffic should have Virtual Network Security Groups and User Defined Routes applied.",
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

func (s *SiteConfigLinuxWebAppSlot) ExpandForCreate(appSettings map[string]string) (*webapps.SiteConfig, error) {
	expanded := &webapps.SiteConfig{}
	expanded.AlwaysOn = pointer.To(s.AlwaysOn)
	expanded.AcrUseManagedIdentityCreds = pointer.To(s.UseManagedIdentityACR)
	expanded.HTTP20Enabled = pointer.To(s.Http2Enabled)
	expanded.ScmIPSecurityRestrictionsUseMain = pointer.To(s.ScmUseMainIpRestriction)
	expanded.LocalMySqlEnabled = pointer.To(s.LocalMysql)
	expanded.LoadBalancing = pointer.ToEnum[webapps.SiteLoadBalancing](s.LoadBalancing)
	expanded.ManagedPipelineMode = pointer.ToEnum[webapps.ManagedPipelineMode](s.ManagedPipelineMode)
	expanded.RemoteDebuggingEnabled = pointer.To(s.RemoteDebugging)
	expanded.Use32BitWorkerProcess = pointer.To(s.Use32BitWorker)
	expanded.WebSocketsEnabled = pointer.To(s.WebSockets)
	expanded.FtpsState = pointer.ToEnum[webapps.FtpsState](s.FtpsState)
	expanded.MinTlsVersion = pointer.ToEnum[webapps.SupportedTlsVersions](s.MinTlsVersion)
	expanded.ScmMinTlsVersion = pointer.ToEnum[webapps.SupportedTlsVersions](s.ScmMinTlsVersion)
	expanded.AutoHealEnabled = pointer.To(false)
	expanded.MinTlsCipherSuite = pointer.ToEnum[webapps.TlsCipherSuites](s.MinTlsCipherSuite)
	expanded.VnetRouteAllEnabled = pointer.To(s.VnetRouteAllEnabled)
	expanded.IPSecurityRestrictionsDefaultAction = pointer.ToEnum[webapps.DefaultAction](s.IpRestrictionDefaultAction)
	expanded.ScmIPSecurityRestrictionsDefaultAction = pointer.ToEnum[webapps.DefaultAction](s.ScmIpRestrictionDefaultAction)

	if s.ApiManagementConfigId != "" {
		expanded.ApiManagementConfig = &webapps.ApiManagementConfig{
			Id: pointer.To(s.ApiManagementConfigId),
		}
	}

	if s.ApiDefinition != "" {
		expanded.ApiDefinition = &webapps.ApiDefinitionInfo{
			Url: pointer.To(s.ApiDefinition),
		}
	}

	if s.AppCommandLine != "" {
		expanded.AppCommandLine = pointer.To(s.AppCommandLine)
	}

	if len(s.ApplicationStack) == 1 {
		linuxAppStack := s.ApplicationStack[0]
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

		if linuxAppStack.PythonVersion != "" {
			expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("PYTHON|%s", linuxAppStack.PythonVersion))
		}

		if linuxAppStack.JavaServer != "" {
			javaString, err := helpers.JavaLinuxFxStringBuilder(linuxAppStack.JavaVersion, linuxAppStack.JavaServer, linuxAppStack.JavaServerVersion)
			if err != nil {
				return nil, fmt.Errorf("could not build linuxFxVersion string: %+v", err)
			}
			expanded.LinuxFxVersion = javaString
		}

		if linuxAppStack.DockerImageName != "" {
			expanded.LinuxFxVersion = pointer.To(helpers.EncodeDockerFxString(linuxAppStack.DockerImageName, linuxAppStack.DockerRegistryUrl))

			if appSettings == nil {
				appSettings = make(map[string]string)
			}

			appSettings["DOCKER_REGISTRY_SERVER_URL"] = linuxAppStack.DockerRegistryUrl
			appSettings["DOCKER_REGISTRY_SERVER_USERNAME"] = linuxAppStack.DockerRegistryUsername
			appSettings["DOCKER_REGISTRY_SERVER_PASSWORD"] = linuxAppStack.DockerRegistryPassword
		}
	}

	expanded.AppSettings = helpers.ExpandAppSettingsForCreate(appSettings)

	if s.AutoSwapSlotName != "" {
		expanded.AutoSwapSlotName = pointer.To(s.AutoSwapSlotName)
	}

	if s.ContainerRegistryMSI != "" {
		expanded.AcrUserManagedIdentityID = pointer.To(s.ContainerRegistryMSI)
	}

	if len(s.DefaultDocuments) != 0 {
		expanded.DefaultDocuments = pointer.To(s.DefaultDocuments)
	}

	if len(s.IpRestriction) != 0 {
		ipRestrictions, err := helpers.ExpandIpRestrictions(s.IpRestriction)
		if err != nil {
			return nil, err
		}

		expanded.IPSecurityRestrictions = ipRestrictions
	}

	if len(s.ScmIpRestriction) != 0 {
		ipRestrictions, err := helpers.ExpandIpRestrictions(s.ScmIpRestriction)
		if err != nil {
			return nil, err
		}

		expanded.ScmIPSecurityRestrictions = ipRestrictions
	}

	if s.RemoteDebuggingVersion != "" {
		expanded.RemoteDebuggingVersion = pointer.To(s.RemoteDebuggingVersion)
	}

	if s.WorkerCount != 0 {
		expanded.NumberOfWorkers = pointer.To(s.WorkerCount)
	}

	if s.HealthCheckPath != "" {
		expanded.HealthCheckPath = pointer.To(s.HealthCheckPath)
	}

	if len(s.Cors) != 0 {
		expanded.Cors = helpers.ExpandCorsSettings(s.Cors)
	}

	if len(s.AutoHealSettings) == 1 {
		expanded.AutoHealEnabled = pointer.To(true)
		expanded.AutoHealRules = helpers.ExpandAutoHealSettingsLinux(s.AutoHealSettings)
	}

	return expanded, nil
}

func (s *SiteConfigLinuxWebAppSlot) ExpandForUpdate(metadata sdk.ResourceMetaData, existing *webapps.SiteConfig, appSettings map[string]string) (*webapps.SiteConfig, error) {
	expanded := *existing

	expanded.AlwaysOn = pointer.To(s.AlwaysOn)
	expanded.AcrUseManagedIdentityCreds = pointer.To(s.UseManagedIdentityACR)
	expanded.HTTP20Enabled = pointer.To(s.Http2Enabled)
	expanded.LocalMySqlEnabled = pointer.To(s.LocalMysql)
	expanded.RemoteDebuggingEnabled = pointer.To(s.RemoteDebugging)
	expanded.ScmIPSecurityRestrictionsUseMain = pointer.To(s.ScmUseMainIpRestriction)
	expanded.Use32BitWorkerProcess = pointer.To(s.Use32BitWorker)
	expanded.WebSocketsEnabled = pointer.To(s.WebSockets)
	expanded.VnetRouteAllEnabled = pointer.To(s.VnetRouteAllEnabled)

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.ApiManagementConfig = &webapps.ApiManagementConfig{
			Id: pointer.To(s.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.ApiDefinition = &webapps.ApiDefinitionInfo{
			Url: pointer.To(s.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = pointer.To(s.AppCommandLine)
	}

	if len(s.ApplicationStack) == 1 {
		linuxAppStack := s.ApplicationStack[0]
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

		if linuxAppStack.PythonVersion != "" {
			expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("PYTHON|%s", linuxAppStack.PythonVersion))
		}

		if linuxAppStack.JavaServer != "" {
			javaString, err := helpers.JavaLinuxFxStringBuilder(linuxAppStack.JavaVersion, linuxAppStack.JavaServer, linuxAppStack.JavaServerVersion)
			if err != nil {
				return nil, fmt.Errorf("could not build linuxFxVersion string: %+v", err)
			}
			expanded.LinuxFxVersion = javaString
		}

		if linuxAppStack.DockerImageName != "" {
			expanded.LinuxFxVersion = pointer.To(helpers.EncodeDockerFxString(linuxAppStack.DockerImageName, linuxAppStack.DockerRegistryUrl))

			if appSettings == nil {
				appSettings = make(map[string]string)
			}

			appSettings["DOCKER_REGISTRY_SERVER_URL"] = linuxAppStack.DockerRegistryUrl
			appSettings["DOCKER_REGISTRY_SERVER_USERNAME"] = linuxAppStack.DockerRegistryUsername
			appSettings["DOCKER_REGISTRY_SERVER_PASSWORD"] = linuxAppStack.DockerRegistryPassword
		}
	} else {
		expanded.LinuxFxVersion = pointer.To("")
	}

	expanded.AppSettings = helpers.ExpandAppSettingsForCreate(appSettings)

	if metadata.ResourceData.HasChange("site_config.0.auto_swap_slot_name") {
		expanded.AutoSwapSlotName = pointer.To(s.AutoSwapSlotName)
	}

	if metadata.ResourceData.HasChange("site_config.0.container_registry_managed_identity_client_id") {
		expanded.AcrUserManagedIdentityID = pointer.To(s.ContainerRegistryMSI)
	}

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = &s.DefaultDocuments
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := helpers.ExpandIpRestrictions(s.IpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction_default_action") {
		expanded.IPSecurityRestrictionsDefaultAction = pointer.ToEnum[webapps.DefaultAction](s.IpRestrictionDefaultAction)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := helpers.ExpandIpRestrictions(s.ScmIpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction_default_action") {
		expanded.ScmIPSecurityRestrictionsDefaultAction = pointer.ToEnum[webapps.DefaultAction](s.ScmIpRestrictionDefaultAction)
	}

	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = pointer.ToEnum[webapps.SiteLoadBalancing](s.LoadBalancing)
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = pointer.ToEnum[webapps.ManagedPipelineMode](s.ManagedPipelineMode)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = pointer.To(s.RemoteDebuggingVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		expanded.FtpsState = pointer.ToEnum[webapps.FtpsState](s.FtpsState)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = pointer.To(s.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.worker_count") {
		expanded.NumberOfWorkers = pointer.To(s.WorkerCount)
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTlsVersion = pointer.ToEnum[webapps.SupportedTlsVersions](s.MinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTlsVersion = pointer.ToEnum[webapps.SupportedTlsVersions](s.ScmMinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_cipher_suite") {
		expanded.MinTlsCipherSuite = pointer.ToEnum[webapps.TlsCipherSuites](s.MinTlsCipherSuite)
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		cors := helpers.ExpandCorsSettings(s.Cors)
		if cors == nil {
			cors = &webapps.CorsSettings{
				AllowedOrigins: &[]string{},
			}
		}
		expanded.Cors = cors
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_heal_setting") {
		expanded.AutoHealEnabled = pointer.To(false)
		if len(s.AutoHealSettings) != 0 {
			expanded.AutoHealEnabled = pointer.To(true)
		}
		expanded.AutoHealRules = helpers.ExpandAutoHealSettingsLinux(s.AutoHealSettings)
	}

	return &expanded, nil
}

func (s *SiteConfigLinuxWebAppSlot) Flatten(appSiteSlotConfig *webapps.SiteConfig) {
	s.AlwaysOn = pointer.From(appSiteSlotConfig.AlwaysOn)
	s.AppCommandLine = pointer.From(appSiteSlotConfig.AppCommandLine)
	s.AutoHealSettings = helpers.FlattenAutoHealSettingsLinux(appSiteSlotConfig.AutoHealRules)
	s.AutoSwapSlotName = pointer.From(appSiteSlotConfig.AutoSwapSlotName)
	s.ContainerRegistryMSI = pointer.From(appSiteSlotConfig.AcrUserManagedIdentityID)
	s.Cors = helpers.FlattenCorsSettings(appSiteSlotConfig.Cors)
	s.DetailedErrorLogging = pointer.From(appSiteSlotConfig.DetailedErrorLoggingEnabled)
	s.Http2Enabled = pointer.From(appSiteSlotConfig.HTTP20Enabled)
	s.IpRestriction = helpers.FlattenIpRestrictions(appSiteSlotConfig.IPSecurityRestrictions)
	s.ManagedPipelineMode = pointer.FromEnum(appSiteSlotConfig.ManagedPipelineMode)
	s.ScmType = pointer.FromEnum(appSiteSlotConfig.ScmType)
	s.FtpsState = pointer.FromEnum(appSiteSlotConfig.FtpsState)
	s.HealthCheckPath = pointer.From(appSiteSlotConfig.HealthCheckPath)
	s.LoadBalancing = pointer.FromEnum(appSiteSlotConfig.LoadBalancing)
	s.LocalMysql = pointer.From(appSiteSlotConfig.LocalMySqlEnabled)
	s.MinTlsVersion = pointer.FromEnum(appSiteSlotConfig.MinTlsVersion)
	s.WorkerCount = pointer.From(appSiteSlotConfig.NumberOfWorkers)
	s.RemoteDebugging = pointer.From(appSiteSlotConfig.RemoteDebuggingEnabled)
	s.RemoteDebuggingVersion = strings.ToUpper(pointer.From(appSiteSlotConfig.RemoteDebuggingVersion))
	s.ScmIpRestriction = helpers.FlattenIpRestrictions(appSiteSlotConfig.ScmIPSecurityRestrictions)
	s.ScmMinTlsVersion = pointer.FromEnum(appSiteSlotConfig.ScmMinTlsVersion)
	s.MinTlsCipherSuite = pointer.FromEnum(appSiteSlotConfig.MinTlsCipherSuite)
	s.ScmUseMainIpRestriction = pointer.From(appSiteSlotConfig.ScmIPSecurityRestrictionsUseMain)
	s.Use32BitWorker = pointer.From(appSiteSlotConfig.Use32BitWorkerProcess)
	s.UseManagedIdentityACR = pointer.From(appSiteSlotConfig.AcrUseManagedIdentityCreds)
	s.WebSockets = pointer.From(appSiteSlotConfig.WebSocketsEnabled)
	s.VnetRouteAllEnabled = pointer.From(appSiteSlotConfig.VnetRouteAllEnabled)
	s.IpRestrictionDefaultAction = pointer.FromEnum(appSiteSlotConfig.IPSecurityRestrictionsDefaultAction)
	s.ScmIpRestrictionDefaultAction = pointer.FromEnum(appSiteSlotConfig.ScmIPSecurityRestrictionsDefaultAction)

	if appSiteSlotConfig.ApiManagementConfig != nil && appSiteSlotConfig.ApiManagementConfig.Id != nil {
		s.ApiManagementConfigId = *appSiteSlotConfig.ApiManagementConfig.Id
	}

	if appSiteSlotConfig.ApiDefinition != nil && appSiteSlotConfig.ApiDefinition.Url != nil {
		s.ApiDefinition = *appSiteSlotConfig.ApiDefinition.Url
	}

	if appSiteSlotConfig.DefaultDocuments != nil {
		s.DefaultDocuments = *appSiteSlotConfig.DefaultDocuments
	}

	if appSiteSlotConfig.LinuxFxVersion != nil {
		var linuxAppStack helpers.ApplicationStackLinux
		s.LinuxFxVersion = *appSiteSlotConfig.LinuxFxVersion
		// Decode the string to docker values
		linuxAppStack = helpers.DecodeApplicationStackLinux(s.LinuxFxVersion)
		s.ApplicationStack = []helpers.ApplicationStackLinux{linuxAppStack}
	}
}

func (s *SiteConfigLinuxWebAppSlot) SetHealthCheckEvictionTime(input map[string]string) {
	if v, ok := input["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"]; ok && v != "" {
		// Discarding the error here as an invalid value should result in `0`
		evictionTime, _ := strconv.Atoi(v)
		s.HealthCheckEvictionTime = int64(evictionTime)
	}
}

func (s *SiteConfigLinuxWebAppSlot) DecodeDockerAppStack(input map[string]string) {
	applicationStack := helpers.ApplicationStackLinux{}
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

	registryHost := helpers.TrimURLScheme(applicationStack.DockerRegistryUrl)
	dockerString := strings.TrimPrefix(s.LinuxFxVersion, "DOCKER|")
	applicationStack.DockerImageName = strings.TrimPrefix(dockerString, registryHost+"/")

	s.ApplicationStack = []helpers.ApplicationStackLinux{applicationStack}
}
