package storage

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func schemaStorageAccountCorsRule(patchEnabled bool) *pluginsdk.Schema {
	// CorsRule "PATCH" method is only supported by blob
	allowedMethods := []string{
		"DELETE",
		"GET",
		"HEAD",
		"MERGE",
		"POST",
		"OPTIONS",
		"PUT",
	}

	if patchEnabled {
		allowedMethods = append(allowedMethods, "PATCH")
	}

	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 5,
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
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
				"allowed_headers": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 64,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
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
