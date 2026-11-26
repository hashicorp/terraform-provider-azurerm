// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_access_policy

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about the permissions from the Management Key Vault Templates.",
		Read:        dataSourceKeyVaultAccessPolicyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The name of the Management Template.",
				Type:        pluginsdk.TypeString,
				Required:    true,
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
				Description: "The certificate permissions for the access policy.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"key_permissions": {
				Description: "The key permissions for the access policy.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"secret_permissions": {
				Description: "The secret permissions for the access policy.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceKeyVaultAccessPolicyRead(d *pluginsdk.ResourceData, _ any) error {
	name := d.Get("name").(string)

	keyPermissions := make([]string, 0)
	if strings.Contains(name, "Key") {
		keyPermissions = append(helpers.KeyPermissionsManagement(), helpers.KeyPermissionsRotationPolicy()...)
	}

	secretPermissions := make([]string, 0)
	if strings.Contains(name, "Secret") {
		secretPermissions = helpers.SecretPermissionsManagement()
	}

	certificatePermissions := make([]string, 0)
	if strings.Contains(name, "Certificate") {
		certificatePermissions = helpers.CertificatePermissionsManagement()
	}

	d.SetId(name)

	d.Set("key_permissions", keyPermissions)
	d.Set("secret_permissions", secretPermissions)
	d.Set("certificate_permissions", certificatePermissions)

	return nil
}
