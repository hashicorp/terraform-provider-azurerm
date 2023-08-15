// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataSourceKeyVaultAccessPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceKeyVaultAccessPolicyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Key Management",
					"Secret Management",
					"Certificate Management",
					"Key & Secret Management",
					"Key & Certificate Management",
					"Secret & Certificate Management",
					"Key, Secret, & Certificate Management",
				}, false),
			},

			// Computed
			"certificate_permissions": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"key_permissions": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"secret_permissions": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceKeyVaultAccessPolicyRead(d *pluginsdk.ResourceData, _ interface{}) error {
	name := d.Get("name").(string)

	keyPermissions := make([]string, 0)
	if strings.Contains(name, "Key") {
		keyPermissions = append(keyPermissionsManagement(), keyPermissionsRotationPolicy()...)
	}

	secretPermissions := make([]string, 0)
	if strings.Contains(name, "Secret") {
		secretPermissions = secretPermissionsManagement()
	}

	certificatePermissions := make([]string, 0)
	if strings.Contains(name, "Certificate") {
		certificatePermissions = certificatePermissionsManagement()
	}

	d.SetId(name)

	d.Set("key_permissions", keyPermissions)
	d.Set("secret_permissions", secretPermissions)
	d.Set("certificate_permissions", certificatePermissions)

	return nil
}
