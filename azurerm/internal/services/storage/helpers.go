package storage

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func schemaStorageAccountCorsRule(patchEnabled bool) *schema.Schema {
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

	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 5,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"allowed_origins": {
					Type:     schema.TypeList,
					Required: true,
					MaxItems: 64,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
				"exposed_headers": {
					Type:     schema.TypeList,
					Required: true,
					MaxItems: 64,
					MinItems: 1,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"allowed_headers": {
					Type:     schema.TypeList,
					Required: true,
					MaxItems: 64,
					MinItems: 1,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"allowed_methods": {
					Type:     schema.TypeList,
					Required: true,
					MaxItems: 64,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.StringInSlice(allowedMethods, false),
					},
				},
				"max_age_in_seconds": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 2000000000),
				},
			},
		},
	}
}
