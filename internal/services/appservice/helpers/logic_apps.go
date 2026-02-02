// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LogicAppSiteConfig struct {
	AlwaysOn                      bool            `tfschema:"always_on"`
	Cors                          []CorsSetting   `tfschema:"cors"`
	FTPSState                     string          `tfschema:"ftps_state"`
	HTTP2Enabled                  bool            `tfschema:"http2_enabled"`
	IpRestriction                 []IpRestriction `tfschema:"ip_restriction"`
	LinuxFxVersion                string          `tfschema:"linux_fx_version"`
	MinTLSVersion                 string          `tfschema:"min_tls_version"`
	PreWarmedInstanceCount        int64           `tfschema:"pre_warmed_instance_count"`
	SCMIPRestriction              []IpRestriction `tfschema:"scm_ip_restriction"`
	SCMUseMainIpRestriction       bool            `tfschema:"scm_use_main_ip_restriction"`
	SCMMinTLSVersion              string          `tfschema:"scm_min_tls_version"`
	SCMType                       string          `tfschema:"scm_type"`
	Use32BitWorkerProcess         bool            `tfschema:"use_32_bit_worker_process"`
	WebSocketsEnabled             bool            `tfschema:"websockets_enabled"`
	HealthCheckPath               string          `tfschema:"health_check_path"`
	ElasticInstanceMinimum        int64           `tfschema:"elastic_instance_minimum"`
	AppScaleLimit                 int64           `tfschema:"app_scale_limit"`
	RuntimeScaleMonitoringEnabled bool            `tfschema:"runtime_scale_monitoring_enabled"`
	DotnetFrameworkVersion        string          `tfschema:"dotnet_framework_version"`
	VNETRouteAllEnabled           bool            `tfschema:"vnet_route_all_enabled"`
	AutoSwapSlotName              string          `tfschema:"auto_swap_slot_name"`

	PublicNetworkAccessEnabled bool `tfschema:"public_network_access_enabled,removedInNextMajorVersion"`
}

func SchemaLogicAppStandardSiteConfig() *pluginsdk.Schema {
	schema := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"cors": CorsSettingsSchema(),

				"ftps_state": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(webapps.FtpsStateAllAllowed),
						string(webapps.FtpsStateDisabled),
						string(webapps.FtpsStateFtpsOnly),
					}, false),
				},

				"http2_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"ip_restriction": IpRestrictionSchema(),

				"linux_fx_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
				},

				"min_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(webapps.SupportedTlsVersionsOnePointTwo),
						string(webapps.SupportedTlsVersionsOnePointThree),
					}, false),
				},

				"pre_warmed_instance_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(0, 20),
				},

				"scm_ip_restriction": IpRestrictionSchema(),

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"scm_min_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(webapps.SupportedTlsVersionsOnePointTwo),
						string(webapps.SupportedTlsVersionsOnePointThree),
					}, false),
				},

				"scm_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(webapps.ScmTypeBitbucketGit),
						string(webapps.ScmTypeBitbucketHg),
						string(webapps.ScmTypeCodePlexGit),
						string(webapps.ScmTypeCodePlexHg),
						string(webapps.ScmTypeDropbox),
						string(webapps.ScmTypeExternalGit),
						string(webapps.ScmTypeExternalHg),
						string(webapps.ScmTypeGitHub),
						string(webapps.ScmTypeLocalGit),
						string(webapps.ScmTypeNone),
						string(webapps.ScmTypeOneDrive),
						string(webapps.ScmTypeTfs),
						string(webapps.ScmTypeVSO),
						string(webapps.ScmTypeVSTSRM),
					}, false),
				},

				"use_32_bit_worker_process": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"websockets_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"health_check_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"elastic_instance_minimum": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(0, 20),
				},

				"app_scale_limit": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntAtLeast(0),
				},

				"runtime_scale_monitoring_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"dotnet_framework_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "v4.0",
					ValidateFunc: validation.StringInSlice([]string{
						"v4.0",
						"v5.0",
						"v6.0",
						"v8.0",
					}, false),
				},

				"vnet_route_all_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
				},

				"auto_swap_slot_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}

	if !features.FivePointOh() {
		schema.Elem.(*pluginsdk.Resource).Schema["public_network_access_enabled"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeBool,
			Optional:   true,
			Computed:   true,
			Deprecated: "the `site_config.public_network_access_enabled` property has been superseded by the `public_network_access` property and will be removed in v5.0 of the AzureRM Provider.",
		}
		schema.Elem.(*pluginsdk.Resource).Schema["scm_min_tls_version"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(webapps.SupportedTlsVersionsOnePointZero),
				string(webapps.SupportedTlsVersionsOnePointOne),
				string(webapps.SupportedTlsVersionsOnePointTwo),
				string(webapps.SupportedTlsVersionsOnePointThree),
			}, false),
		}
		schema.Elem.(*pluginsdk.Resource).Schema["min_tls_version"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(webapps.SupportedTlsVersionsOnePointZero),
				string(webapps.SupportedTlsVersionsOnePointOne),
				string(webapps.SupportedTlsVersionsOnePointTwo),
				string(webapps.SupportedTlsVersionsOnePointThree),
			}, false),
		}
	}

	return schema
}
