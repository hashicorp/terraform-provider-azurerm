package azure

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform/helper/schema"
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
