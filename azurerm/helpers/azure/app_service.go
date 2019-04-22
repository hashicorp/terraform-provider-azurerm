package azure

import (
	"log"
	"net"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SchemaAppServiceSiteConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"always_on": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},

				"app_command_line": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"default_documents": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},

				"dotnet_framework_version": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "v4.0",
					ValidateFunc: validation.StringInSlice([]string{
						"v2.0",
						"v4.0",
					}, true),
					DiffSuppressFunc: suppress.CaseDifference,
				},

				"http2_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},

				"ip_restriction": {
					Type:       schema.TypeList,
					Optional:   true,
					Computed:   true,
					ConfigMode: schema.SchemaConfigModeAttr,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ip_address": {
								Type:     schema.TypeString,
								Required: true,
							},
							"subnet_mask": {
								Type:     schema.TypeString,
								Optional: true,
								Default:  "255.255.255.255",
							},
						},
					},
				},

				"java_version": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"1.7",
						"1.8",
						"11",
					}, false),
				},

				"java_container": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"JETTY",
						"TOMCAT",
					}, true),
					DiffSuppressFunc: suppress.CaseDifference,
				},

				"java_container_version": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"local_mysql_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},

				"managed_pipeline_mode": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.Classic),
						string(web.Integrated),
					}, true),
					DiffSuppressFunc: suppress.CaseDifference,
				},

				"php_version": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"5.5",
						"5.6",
						"7.0",
						"7.1",
						"7.2",
					}, false),
				},

				"python_version": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"2.7",
						"3.4",
					}, false),
				},

				"remote_debugging_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},

				"remote_debugging_version": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"VS2012",
						"VS2013",
						"VS2015",
						"VS2017",
					}, true),
					DiffSuppressFunc: suppress.CaseDifference,
				},

				"scm_type": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  string(web.ScmTypeNone),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ScmTypeBitbucketGit),
						string(web.ScmTypeBitbucketHg),
						string(web.ScmTypeCodePlexGit),
						string(web.ScmTypeCodePlexHg),
						string(web.ScmTypeDropbox),
						string(web.ScmTypeExternalGit),
						string(web.ScmTypeExternalHg),
						string(web.ScmTypeGitHub),
						string(web.ScmTypeLocalGit),
						string(web.ScmTypeNone),
						string(web.ScmTypeOneDrive),
						string(web.ScmTypeTfs),
						string(web.ScmTypeVSO),
						// Not in the specs, but is set by Azure Pipelines
						// https://github.com/Microsoft/azure-pipelines-tasks/blob/master/Tasks/AzureRmWebAppDeploymentV4/operations/AzureAppServiceUtility.ts#L19
						// upstream issue: https://github.com/Azure/azure-rest-api-specs/issues/5345
						"VSTSRM",
					}, false),
				},

				"use_32_bit_worker_process": {
					Type:     schema.TypeBool,
					Optional: true,
				},

				"websockets_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},

				"ftps_state": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.AllAllowed),
						string(web.Disabled),
						string(web.FtpsOnly),
					}, false),
				},
				"linux_fx_version": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},

				"min_tls_version": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.OneFullStopZero),
						string(web.OneFullStopOne),
						string(web.OneFullStopTwo),
					}, false),
				},

				"virtual_network_name": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"cors": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"allowed_origins": {
								Type:     schema.TypeSet,
								Required: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"support_credentials": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  false,
							},
						},
					},
				},
			},
		},
	}
}

func SchemaAppServiceDataSourceSiteConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"always_on": {
					Type:     schema.TypeBool,
					Computed: true,
				},

				"app_command_line": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"default_documents": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},

				"dotnet_framework_version": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"http2_enabled": {
					Type:     schema.TypeBool,
					Computed: true,
				},

				"ip_restriction": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ip_address": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"subnet_mask": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},

				"java_version": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"java_container": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"java_container_version": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"local_mysql_enabled": {
					Type:     schema.TypeBool,
					Computed: true,
				},

				"managed_pipeline_mode": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"php_version": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"python_version": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"remote_debugging_enabled": {
					Type:     schema.TypeBool,
					Computed: true,
				},

				"remote_debugging_version": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"scm_type": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"use_32_bit_worker_process": {
					Type:     schema.TypeBool,
					Computed: true,
				},

				"websockets_enabled": {
					Type:     schema.TypeBool,
					Computed: true,
				},

				"ftps_state": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"linux_fx_version": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"min_tls_version": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"virtual_network_name": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"cors": {
					Type:     schema.TypeList,
					Computed: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"allowed_origins": {
								Type:     schema.TypeSet,
								Computed: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"support_credentials": {
								Type:     schema.TypeBool,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func ExpandAppServiceCorsSettings(input interface{}) web.CorsSettings {
	settings := input.([]interface{})
	corsSettings := web.CorsSettings{}

	if len(settings) == 0 {
		return corsSettings
	}

	setting := settings[0].(map[string]interface{})

	if v, ok := setting["allowed_origins"]; ok {
		input := v.(*schema.Set).List()

		allowedOrigins := make([]string, 0)
		for _, param := range input {
			allowedOrigins = append(allowedOrigins, param.(string))
		}

		corsSettings.AllowedOrigins = &allowedOrigins
	}

	if v, ok := setting["support_credentials"]; ok {
		corsSettings.SupportCredentials = utils.Bool(v.(bool))
	}

	return corsSettings
}

func FlattenAppServiceCorsSettings(input *web.CorsSettings) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	allowedOrigins := make([]interface{}, 0)
	if s := input.AllowedOrigins; s != nil {
		for _, v := range *s {
			allowedOrigins = append(allowedOrigins, v)
		}
	}
	result["allowed_origins"] = schema.NewSet(schema.HashString, allowedOrigins)

	if input.SupportCredentials != nil {
		result["support_credentials"] = *input.SupportCredentials
	}

	return append(results, result)
}

func ExpandAppServiceSiteConfig(input interface{}) web.SiteConfig {
	configs := input.([]interface{})
	siteConfig := web.SiteConfig{}

	if len(configs) == 0 {
		return siteConfig
	}

	config := configs[0].(map[string]interface{})

	if v, ok := config["always_on"]; ok {
		siteConfig.AlwaysOn = utils.Bool(v.(bool))
	}

	if v, ok := config["app_command_line"]; ok {
		siteConfig.AppCommandLine = utils.String(v.(string))
	}

	if v, ok := config["default_documents"]; ok {
		input := v.([]interface{})

		documents := make([]string, 0)
		for _, document := range input {
			documents = append(documents, document.(string))
		}

		siteConfig.DefaultDocuments = &documents
	}

	if v, ok := config["dotnet_framework_version"]; ok {
		siteConfig.NetFrameworkVersion = utils.String(v.(string))
	}

	if v, ok := config["java_version"]; ok {
		siteConfig.JavaVersion = utils.String(v.(string))
	}

	if v, ok := config["java_container"]; ok {
		siteConfig.JavaContainer = utils.String(v.(string))
	}

	if v, ok := config["java_container_version"]; ok {
		siteConfig.JavaContainerVersion = utils.String(v.(string))
	}

	if v, ok := config["linux_fx_version"]; ok {
		siteConfig.LinuxFxVersion = utils.String(v.(string))
	}

	if v, ok := config["http2_enabled"]; ok {
		siteConfig.HTTP20Enabled = utils.Bool(v.(bool))
	}

	if v, ok := config["ip_restriction"]; ok {
		ipSecurityRestrictions := v.([]interface{})
		restrictions := make([]web.IPSecurityRestriction, 0)
		for _, ipSecurityRestriction := range ipSecurityRestrictions {
			restriction := ipSecurityRestriction.(map[string]interface{})

			ipAddress := restriction["ip_address"].(string)
			mask := restriction["subnet_mask"].(string)
			// the 2018-02-01 API expects a blank subnet mask and an IP address in CIDR format: a.b.c.d/x
			// so translate the IP and mask if necessary
			restrictionMask := ""
			cidrAddress := ipAddress
			if mask != "" {
				ipNet := net.IPNet{IP: net.ParseIP(ipAddress), Mask: net.IPMask(net.ParseIP(mask))}
				cidrAddress = ipNet.String()
			} else if !strings.Contains(ipAddress, "/") {
				cidrAddress += "/32"
			}

			restrictions = append(restrictions, web.IPSecurityRestriction{
				IPAddress:  &cidrAddress,
				SubnetMask: &restrictionMask,
			})
		}
		siteConfig.IPSecurityRestrictions = &restrictions
	}

	if v, ok := config["local_mysql_enabled"]; ok {
		siteConfig.LocalMySQLEnabled = utils.Bool(v.(bool))
	}

	if v, ok := config["managed_pipeline_mode"]; ok {
		siteConfig.ManagedPipelineMode = web.ManagedPipelineMode(v.(string))
	}

	if v, ok := config["php_version"]; ok {
		siteConfig.PhpVersion = utils.String(v.(string))
	}

	if v, ok := config["python_version"]; ok {
		siteConfig.PythonVersion = utils.String(v.(string))
	}

	if v, ok := config["remote_debugging_enabled"]; ok {
		siteConfig.RemoteDebuggingEnabled = utils.Bool(v.(bool))
	}

	if v, ok := config["remote_debugging_version"]; ok {
		siteConfig.RemoteDebuggingVersion = utils.String(v.(string))
	}

	if v, ok := config["use_32_bit_worker_process"]; ok {
		siteConfig.Use32BitWorkerProcess = utils.Bool(v.(bool))
	}

	if v, ok := config["websockets_enabled"]; ok {
		siteConfig.WebSocketsEnabled = utils.Bool(v.(bool))
	}

	if v, ok := config["scm_type"]; ok {
		siteConfig.ScmType = web.ScmType(v.(string))
	}

	if v, ok := config["ftps_state"]; ok {
		siteConfig.FtpsState = web.FtpsState(v.(string))
	}

	if v, ok := config["min_tls_version"]; ok {
		siteConfig.MinTLSVersion = web.SupportedTLSVersions(v.(string))
	}

	if v, ok := config["virtual_network_name"]; ok {
		siteConfig.VnetName = utils.String(v.(string))
	}

	if v, ok := config["cors"]; ok {
		corsSettings := v.(interface{})
		expand := ExpandAppServiceCorsSettings(corsSettings)
		siteConfig.Cors = &expand
	}

	return siteConfig
}

func FlattenAppServiceSiteConfig(input *web.SiteConfig) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{})

	if input == nil {
		log.Printf("[DEBUG] SiteConfig is nil")
		return results
	}

	if input.AlwaysOn != nil {
		result["always_on"] = *input.AlwaysOn
	}

	if input.AppCommandLine != nil {
		result["app_command_line"] = *input.AppCommandLine
	}

	documents := make([]string, 0)
	if s := input.DefaultDocuments; s != nil {
		documents = *s
	}
	result["default_documents"] = documents

	if input.NetFrameworkVersion != nil {
		result["dotnet_framework_version"] = *input.NetFrameworkVersion
	}

	if input.JavaVersion != nil {
		result["java_version"] = *input.JavaVersion
	}

	if input.JavaContainer != nil {
		result["java_container"] = *input.JavaContainer
	}

	if input.JavaContainerVersion != nil {
		result["java_container_version"] = *input.JavaContainerVersion
	}

	if input.LocalMySQLEnabled != nil {
		result["local_mysql_enabled"] = *input.LocalMySQLEnabled
	}

	if input.HTTP20Enabled != nil {
		result["http2_enabled"] = *input.HTTP20Enabled
	}

	restrictions := make([]interface{}, 0)
	if vs := input.IPSecurityRestrictions; vs != nil {
		for _, v := range *vs {
			block := make(map[string]interface{})
			if ip := v.IPAddress; ip != nil {
				// the 2018-02-01 API uses CIDR format (a.b.c.d/x), so translate that back to IP and mask
				if strings.Contains(*ip, "/") {
					ipAddr, ipNet, _ := net.ParseCIDR(*ip)
					block["ip_address"] = ipAddr.String()
					mask := net.IP(ipNet.Mask)
					block["subnet_mask"] = mask.String()
				} else {
					block["ip_address"] = *ip
				}
			}
			if subnet := v.SubnetMask; subnet != nil {
				block["subnet_mask"] = *subnet
			}
			restrictions = append(restrictions, block)
		}
	}
	result["ip_restriction"] = restrictions

	result["managed_pipeline_mode"] = string(input.ManagedPipelineMode)

	if input.PhpVersion != nil {
		result["php_version"] = *input.PhpVersion
	}

	if input.PythonVersion != nil {
		result["python_version"] = *input.PythonVersion
	}

	if input.RemoteDebuggingEnabled != nil {
		result["remote_debugging_enabled"] = *input.RemoteDebuggingEnabled
	}

	if input.RemoteDebuggingVersion != nil {
		result["remote_debugging_version"] = *input.RemoteDebuggingVersion
	}

	if input.Use32BitWorkerProcess != nil {
		result["use_32_bit_worker_process"] = *input.Use32BitWorkerProcess
	}

	if input.WebSocketsEnabled != nil {
		result["websockets_enabled"] = *input.WebSocketsEnabled
	}

	if input.LinuxFxVersion != nil {
		result["linux_fx_version"] = *input.LinuxFxVersion
	}

	if input.VnetName != nil {
		result["virtual_network_name"] = *input.VnetName
	}

	result["scm_type"] = string(input.ScmType)
	result["ftps_state"] = string(input.FtpsState)
	result["min_tls_version"] = string(input.MinTLSVersion)

	result["cors"] = FlattenAppServiceCorsSettings(input.Cors)

	return append(results, result)
}
