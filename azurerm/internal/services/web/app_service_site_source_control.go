package web

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func schemaAppServiceSiteSourceControl() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		Optional:      true,
		MaxItems:      1,
		Computed:      true,
		ConflictsWith: []string{"site_config.0.scm_type"},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"repo_url": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					AtLeastOneOf: []string{"source_control.0.repo_url", "source_control.0.branch", "source_control.0.manual_integration",
						"source_control.0.use_mercurial", "source_control.0.rollback_enabled",
					},
				},

				"branch": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					AtLeastOneOf: []string{"source_control.0.repo_url", "source_control.0.branch", "source_control.0.manual_integration",
						"source_control.0.use_mercurial", "source_control.0.rollback_enabled",
					},
				},

				"manual_integration": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
					AtLeastOneOf: []string{"source_control.0.repo_url", "source_control.0.branch", "source_control.0.manual_integration",
						"source_control.0.use_mercurial", "source_control.0.rollback_enabled",
					},
				},

				"use_mercurial": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
					AtLeastOneOf: []string{"source_control.0.repo_url", "source_control.0.branch", "source_control.0.manual_integration",
						"source_control.0.use_mercurial", "source_control.0.rollback_enabled",
					},
				},

				"rollback_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
					AtLeastOneOf: []string{"source_control.0.repo_url", "source_control.0.branch", "source_control.0.manual_integration",
						"source_control.0.use_mercurial", "source_control.0.rollback_enabled",
					},
				},
			},
		},
	}
}

func schemaAppServiceSiteSourceControlDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"repo_url": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"branch": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"manual_integration": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"use_mercurial": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"rollback_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

func expandAppServiceSiteSourceControl(d *pluginsdk.ResourceData) *web.SiteSourceControlProperties {
	sourceControlRaw := d.Get("source_control").([]interface{})
	sourceControl := sourceControlRaw[0].(map[string]interface{})

	result := &web.SiteSourceControlProperties{
		RepoURL:                   utils.String(sourceControl["repo_url"].(string)),
		Branch:                    utils.String(sourceControl["branch"].(string)),
		IsManualIntegration:       utils.Bool(sourceControl["manual_integration"].(bool)),
		IsMercurial:               utils.Bool(sourceControl["use_mercurial"].(bool)),
		DeploymentRollbackEnabled: utils.Bool(sourceControl["rollback_enabled"].(bool)),
	}

	return result
}

func flattenAppServiceSourceControl(input *web.SiteSourceControlProperties) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{})

	if input == nil {
		log.Printf("[DEBUG] SiteSourceControlProperties is nil")
		return results
	}

	if input.RepoURL != nil && *input.RepoURL != "" {
		result["repo_url"] = *input.RepoURL
	}

	if input.Branch != nil && *input.Branch != "" {
		result["branch"] = *input.Branch
	} else {
		result["branch"] = "master"
	}

	result["use_mercurial"] = *input.IsMercurial

	result["manual_integration"] = *input.IsManualIntegration

	result["rollback_enabled"] = *input.DeploymentRollbackEnabled

	return append(results, result)
}
