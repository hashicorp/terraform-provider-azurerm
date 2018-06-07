package schema

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2016-09-01/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func AppServiceSiteConfigSchema() *schema.Schema {
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
					DiffSuppressFunc: IgnoreCaseDiffSuppressFunc,
				},

				"http2_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},

				"ip_restriction": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
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
					}, false),
				},

				"java_container": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"JETTY",
						"TOMCAT",
					}, true),
					DiffSuppressFunc: IgnoreCaseDiffSuppressFunc,
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
					DiffSuppressFunc: IgnoreCaseDiffSuppressFunc,
				},

				"php_version": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"5.5",
						"5.6",
						"7.0",
						"7.1",
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
					DiffSuppressFunc: IgnoreCaseDiffSuppressFunc,
				},

				"scm_type": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  string(web.ScmTypeNone),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ScmTypeNone),
						string(web.ScmTypeLocalGit),
					}, false),
				},

				"use_32_bit_worker_process": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},

				"websockets_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
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

			restrictions = append(restrictions, web.IPSecurityRestriction{
				IPAddress:  &ipAddress,
				SubnetMask: &mask,
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

	return siteConfig
}

func FlattenAppServiceSiteConfig(input *web.SiteConfig) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{}, 0)

	if input == nil {
		log.Printf("[DEBUG] SiteConfig is nil")
		return results
	}

	if input.AlwaysOn != nil {
		result["always_on"] = *input.AlwaysOn
	}

	if input.DefaultDocuments != nil {
		documents := make([]string, 0)
		for _, document := range *input.DefaultDocuments {
			documents = append(documents, document)
		}

		result["default_documents"] = documents
	}

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
			result := make(map[string]interface{}, 0)
			if ip := v.IPAddress; ip != nil {
				result["ip_address"] = *ip
			}
			if subnet := v.SubnetMask; subnet != nil {
				result["subnet_mask"] = *subnet
			}
			restrictions = append(restrictions, result)
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

	result["scm_type"] = string(input.ScmType)

	return append(results, result)
}
