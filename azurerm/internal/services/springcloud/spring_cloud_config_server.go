package springcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func SchemaConfigServerHttpBasicAuth(conflictsWith ...string) *schema.Schema {
	s := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"username": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"password": {
					Type:      schema.TypeString,
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

func SchemaConfigServerSSHAuth(conflictsWith ...string) *schema.Schema {
	s := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"private_key": {
					Type:         schema.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"host_key": {
					Type:         schema.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"host_key_algorithm": {
					Type:     schema.TypeString,
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
					Type:     schema.TypeBool,
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

func DataSourceSchemaConfigServerHttpBasicAuth() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"username": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"password": {
					Type:      schema.TypeString,
					Computed:  true,
					Sensitive: true,
				},
			},
		},
	}
}

func DataSourceSchemaConfigServerSSHAuth() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"private_key": {
					Type:      schema.TypeString,
					Computed:  true,
					Sensitive: true,
				},
				"host_key": {
					Type:      schema.TypeString,
					Computed:  true,
					Sensitive: true,
				},
				"host_key_algorithm": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"strict_host_key_checking_enabled": {
					Type:     schema.TypeBool,
					Computed: true,
				},
			},
		},
	}
}
