package azure

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SchemaWebCorsSettings() *schema.Schema {
	return &schema.Schema{
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
	}
}

func SchemaWebSourceControl() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"repo_url": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validate.NoEmptyStrings,
				},
				"branch": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validate.NoEmptyStrings,
				},
				"is_manual_integration": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"deployment_rollback_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"is_mercurial": {
					Type:     schema.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

func ExpandWebCorsSettings(input interface{}) web.CorsSettings {
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

func ExpandWebSourceControl(input []interface{}) *web.SiteSourceControlProperties {
	if len(input) == 0 {
		return nil
	}

	sourceControl := input[0].(map[string]interface{})
	sourceControlProperties := web.SiteSourceControlProperties{}

	if v, ok := sourceControl["repo_url"]; ok {
		sourceControlProperties.RepoURL = utils.String(v.(string))
	}

	if v, ok := sourceControl["branch"]; ok {
		sourceControlProperties.Branch = utils.String(v.(string))
	}

	if v, ok := sourceControl["deployment_rollback_enabled"]; ok {
		sourceControlProperties.DeploymentRollbackEnabled = utils.Bool(v.(bool))
	}

	if v, ok := sourceControl["is_manual_integration"]; ok {
		sourceControlProperties.IsManualIntegration = utils.Bool(v.(bool))
	}

	if v, ok := sourceControl["is_mercurial"]; ok {
		sourceControlProperties.IsMercurial = utils.Bool(v.(bool))
	}

	return &sourceControlProperties
}

func FlattenWebCorsSettings(input *web.CorsSettings) []interface{} {
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

func FlattenWebSourceControl(input *web.SiteSourceControlProperties) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	if input.RepoURL != nil {
		result["repo_url"] = *input.RepoURL
	}

	if input.Branch != nil {
		result["branch"] = *input.Branch
	}

	if input.DeploymentRollbackEnabled != nil {
		result["deployment_rollback_enabled"] = *input.DeploymentRollbackEnabled
	}

	if input.IsManualIntegration != nil {
		result["is_manual_integration"] = *input.IsManualIntegration
	}

	if input.IsMercurial != nil {
		result["is_mercurial"] = *input.IsMercurial
	}

	return append(results, result)
}
