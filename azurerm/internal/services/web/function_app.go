package web

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func schemaAppServiceFunctionAppSiteConfig() *schema.Schema {
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

				"use_32_bit_worker_process": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true, // TODO - toggleable?
				},

				"websockets_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false, // TODO - service defaults to false?
				},

				"linux_fx_version": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},

				"http2_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},

				"ip_restriction": { // TODO - needs extending
					Type:       schema.TypeList,
					Optional:   true,
					Computed:   true,
					ConfigMode: schema.SchemaConfigModeAttr,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ip_address": {
								Type:         schema.TypeString,
								Optional:     true,
								ValidateFunc: validate.CIDR,
							},
							"subnet_id": {
								Type:         schema.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"name": {
								Type:         schema.TypeString,
								Optional:     true,
								Computed:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"priority": {
								Type:         schema.TypeInt,
								Optional:     true,
								Default:      65000,
								ValidateFunc: validation.IntBetween(1, 2147483647),
							},
							"action": {
								Type:     schema.TypeString,
								Default:  "Allow",
								Optional: true,
								ValidateFunc: validation.StringInSlice([]string{
									"Allow",
									"Deny",
								}, false),
							},
						},
					},
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

				"pre_warmed_instance_count": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 10),
				},

				"cors": azure.SchemaWebCorsSettings(),

				// TODO - Add scm_ip_restriction
				//"scm_ip_restriction": {
				//	Type:       schema.TypeList,
				//	Optional:   true,
				//	Computed:   true,
				//	ConfigMode: schema.SchemaConfigModeAttr,
				//	Elem: &schema.Resource{
				//		Schema: map[string]*schema.Schema{
				//			"ip_address": {
				//				Type:         schema.TypeString,
				//				Optional:     true,
				//				ValidateFunc: validate.CIDR,
				//			},
				//			"virtual_network_subnet_id": {
				//				Type:         schema.TypeString,
				//				Optional:     true,
				//				ValidateFunc: validation.StringIsNotEmpty,
				//			},
				//			"name": {
				//				Type:         schema.TypeString,
				//				Optional:     true,
				//				Computed:     true,
				//				ValidateFunc: validation.StringIsNotEmpty,
				//			},
				//			"priority": {
				//				Type:         schema.TypeInt,
				//				Optional:     true,
				//				Default:      65000,
				//				ValidateFunc: validation.IntBetween(1, 2147483647),
				//			},
				//			"action": {
				//				Type:     schema.TypeString,
				//				Optional: true,
				//				Default:  "Allow",
				//				ValidateFunc: validation.StringInSlice([]string{
				//					"Allow",
				//					"Deny",
				//				}, true),
				//			},
				//		},
				//	},
				//},
				//
				"scm_type": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
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
				// TODO - Add scm_use_main_ip_restriction
				//"scm_use_main_ip_restriction": {
				//	Type:     schema.TypeBool,
				//	Optional: true,
				//	Default:  false,
				//},
			},
		},
	}
}
