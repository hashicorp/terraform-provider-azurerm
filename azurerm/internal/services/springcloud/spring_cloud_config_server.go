package springcloud

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func SchemaConfigServerHttpBasicAuth(conflictsWith ...string) *pluginsdk.Schema {
	s := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"username": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"password": {
					Type:      pluginsdk.TypeString,
					Required:  true,
					Sensitive: true,
				},
			},
		},
	}
	if len(conflictsWith) > 0 {
		s.ConflictsWith = conflictsWith
	}
	return s
}

func SchemaConfigServerSSHAuth(conflictsWith ...string) *pluginsdk.Schema {
	s := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"private_key": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"host_key": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"host_key_algorithm": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"ssh-dss",
						"ssh-rsa",
						"ecdsa-sha2-nistp256",
						"ecdsa-sha2-nistp384",
						"ecdsa-sha2-nistp521",
					}, false),
				},

				"strict_host_key_checking_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},
			},
		},
	}

	if len(conflictsWith) > 0 {
		s.ConflictsWith = conflictsWith
	}
	return s
}

func DataSourceSchemaConfigServerHttpBasicAuth() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"username": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"password": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},
			},
		},
	}
}

func DataSourceSchemaConfigServerSSHAuth() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"private_key": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},
				"host_key": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},
				"host_key_algorithm": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"strict_host_key_checking_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
	}
}
