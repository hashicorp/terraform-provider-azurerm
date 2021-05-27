package common

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-01-15/documentdb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SchemaCorsRule() *pluginsdk.Schema {
	allowedMethods := []string{
		"DELETE",
		"GET",
		"HEAD",
		"MERGE",
		"POST",
		"OPTIONS",
		"PUT",
		"PATCH",
	}

	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"allowed_origins": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 64,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"exposed_headers": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 64,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"allowed_headers": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 64,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"allowed_methods": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 64,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(allowedMethods, false),
					},
				},

				"max_age_in_seconds": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 2000000000),
				},
			},
		},
	}
}

func ExpandCosmosCorsRule(input []interface{}) *[]documentdb.CorsPolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	corsRules := make([]documentdb.CorsPolicy, 0)

	if len(input) == 0 {
		return &corsRules
	}

	for _, attr := range input {
		corsRuleAttr := attr.(map[string]interface{})
		corsRule := documentdb.CorsPolicy{}
		corsRule.AllowedOrigins = utils.String(strings.Join(*utils.ExpandStringSlice(corsRuleAttr["allowed_origins"].([]interface{})), ","))
		corsRule.ExposedHeaders = utils.String(strings.Join(*utils.ExpandStringSlice(corsRuleAttr["exposed_headers"].([]interface{})), ","))
		corsRule.AllowedHeaders = utils.String(strings.Join(*utils.ExpandStringSlice(corsRuleAttr["allowed_headers"].([]interface{})), ","))
		corsRule.AllowedMethods = utils.String(strings.Join(*utils.ExpandStringSlice(corsRuleAttr["allowed_methods"].([]interface{})), ","))
		corsRule.MaxAgeInSeconds = utils.Int64(int64(corsRuleAttr["max_age_in_seconds"].(int)))

		corsRules = append(corsRules, corsRule)
	}

	return &corsRules
}

func FlattenCosmosCorsRule(input *[]documentdb.CorsPolicy) []interface{} {
	corsRules := make([]interface{}, 0)

	if input == nil || len(*input) == 0 {
		return corsRules
	}

	for _, corsRule := range *input {
		var maxAgeInSeconds int

		if corsRule.MaxAgeInSeconds != nil {
			maxAgeInSeconds = int(*corsRule.MaxAgeInSeconds)
		}

		corsRules = append(corsRules, map[string]interface{}{
			"allowed_headers":    flattenCorsProperty(corsRule.AllowedHeaders),
			"allowed_origins":    flattenCorsProperty(corsRule.AllowedOrigins),
			"allowed_methods":    flattenCorsProperty(corsRule.AllowedMethods),
			"exposed_headers":    flattenCorsProperty(corsRule.ExposedHeaders),
			"max_age_in_seconds": maxAgeInSeconds,
		})
	}

	return corsRules
}

func flattenCorsProperty(input *string) []interface{} {
	results := make([]interface{}, 0, len(*input))

	origins := strings.Split(*input, ",")
	for _, origin := range origins {
		results = append(results, origin)
	}

	return results
}
