package graph

import (
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"
)

func SchemaAppRoles() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"allowed_member_types": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"description": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"display_name": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"is_enabled": {
					Type:     schema.TypeBool,
					Computed: true,
				},

				"value": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func SchemaOauth2Permissions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"admin_consent_description": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"admin_consent_display_name": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"is_enabled": {
					Type:     schema.TypeBool,
					Computed: true,
				},

				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"user_consent_description": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"user_consent_display_name": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"value": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func FlattenAppRoles(in *[]graphrbac.AppRole) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	appRoles := make([]interface{}, 0)
	for _, role := range *in {
		appRole := make(map[string]interface{})
		if role.ID != nil {
			appRole["id"] = *role.ID
		}
		if role.AllowedMemberTypes != nil {
			appRole["allowed_member_types"] = *role.AllowedMemberTypes
		}
		if role.Description != nil {
			appRole["description"] = *role.Description
		}
		if role.DisplayName != nil {
			appRole["display_name"] = *role.DisplayName
		}
		if role.IsEnabled != nil {
			appRole["is_enabled"] = *role.IsEnabled
		}
		if role.Value != nil {
			appRole["value"] = *role.Value
		}
		appRoles = append(appRoles, appRole)
	}

	return appRoles
}

func FlattenOauth2Permissions(in *[]graphrbac.OAuth2Permission) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0)
	for _, p := range *in {
		permission := make(map[string]interface{})
		if v := p.AdminConsentDescription; v != nil {
			permission["admin_consent_description"] = v
		}
		if v := p.AdminConsentDisplayName; v != nil {
			permission["admin_consent_display_name"] = v
		}
		if v := p.ID; v != nil {
			permission["id"] = v
		}
		if v := p.IsEnabled; v != nil {
			permission["is_enabled"] = *v
		}
		if v := p.Type; v != nil {
			permission["type"] = v
		}
		if v := p.UserConsentDescription; v != nil {
			permission["user_consent_description"] = v
		}
		if v := p.UserConsentDisplayName; v != nil {
			permission["user_consent_display_name"] = v
		}
		if v := p.Value; v != nil {
			permission["value"] = v
		}

		result = append(result, permission)
	}

	return result
}
