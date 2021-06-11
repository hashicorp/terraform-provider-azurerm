package web

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SchemaWebCorsSettings() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"allowed_origins": {
					Type:     pluginsdk.TypeSet,
					Required: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},
				"support_credentials": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
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
		input := v.(*pluginsdk.Set).List()

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
	result["allowed_origins"] = pluginsdk.NewSet(pluginsdk.HashString, allowedOrigins)

	if input.SupportCredentials != nil {
		result["support_credentials"] = *input.SupportCredentials
	}

	return append(results, result)
}
