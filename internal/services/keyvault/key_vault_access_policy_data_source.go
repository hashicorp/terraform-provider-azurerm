// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2021-10-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
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
	templateManagementPermissions := map[string][]string{}
	if features.FourPointOh() {
		templateManagementPermissions["key"] = []string{
			azure.TitleCase(string(vaults.KeyPermissionsGet)),
			azure.TitleCase(string(vaults.KeyPermissionsList)),
			azure.TitleCase(string(vaults.KeyPermissionsUpdate)),
			azure.TitleCase(string(vaults.KeyPermissionsCreate)),
			azure.TitleCase(string(vaults.KeyPermissionsImport)),
			azure.TitleCase(string(vaults.KeyPermissionsDelete)),
			azure.TitleCase(string(vaults.KeyPermissionsRecover)),
			azure.TitleCase(string(vaults.KeyPermissionsBackup)),
			azure.TitleCase(string(vaults.KeyPermissionsRestore)),
		}
		templateManagementPermissions["secret"] = []string{
			azure.TitleCase(string(vaults.SecretPermissionsGet)),
			azure.TitleCase(string(vaults.SecretPermissionsList)),
			azure.TitleCase(string(vaults.SecretPermissionsSet)),
			azure.TitleCase(string(vaults.SecretPermissionsDelete)),
			azure.TitleCase(string(vaults.SecretPermissionsRecover)),
			azure.TitleCase(string(vaults.SecretPermissionsBackup)),
			azure.TitleCase(string(vaults.SecretPermissionsRestore)),
		}
		templateManagementPermissions["certificate"] = []string{
			azure.TitleCase(string(vaults.CertificatePermissionsGet)),
			azure.TitleCase(string(vaults.CertificatePermissionsList)),
			azure.TitleCase(string(vaults.CertificatePermissionsUpdate)),
			azure.TitleCase(string(vaults.CertificatePermissionsCreate)),
			azure.TitleCase(string(vaults.CertificatePermissionsImport)),
			azure.TitleCase(string(vaults.CertificatePermissionsDelete)),
			azure.TitleCase(string(vaults.CertificatePermissionsManagecontacts)),
			azure.TitleCase(string(vaults.CertificatePermissionsManageissuers)),
			azure.TitleCase(string(vaults.CertificatePermissionsGetissuers)),
			azure.TitleCase(string(vaults.CertificatePermissionsListissuers)),
			azure.TitleCase(string(vaults.CertificatePermissionsSetissuers)),
			azure.TitleCase(string(vaults.CertificatePermissionsDeleteissuers)),
		}
	} else {
		templateManagementPermissions["key"] = []string{
			string(vaults.KeyPermissionsGet),
			string(vaults.KeyPermissionsList),
			string(vaults.KeyPermissionsUpdate),
			string(vaults.KeyPermissionsCreate),
			string(vaults.KeyPermissionsImport),
			string(vaults.KeyPermissionsDelete),
			string(vaults.KeyPermissionsRecover),
			string(vaults.KeyPermissionsBackup),
			string(vaults.KeyPermissionsRestore),
		}
		templateManagementPermissions["secret"] = []string{
			string(vaults.SecretPermissionsGet),
			string(vaults.SecretPermissionsList),
			string(vaults.SecretPermissionsSet),
			string(vaults.SecretPermissionsDelete),
			string(vaults.SecretPermissionsRecover),
			string(vaults.SecretPermissionsBackup),
			string(vaults.SecretPermissionsRestore),
		}
		templateManagementPermissions["certificate"] = []string{
			string(vaults.CertificatePermissionsGet),
			string(vaults.CertificatePermissionsList),
			string(vaults.CertificatePermissionsUpdate),
			string(vaults.CertificatePermissionsCreate),
			string(vaults.CertificatePermissionsImport),
			string(vaults.CertificatePermissionsDelete),
			string(vaults.CertificatePermissionsManagecontacts),
			string(vaults.CertificatePermissionsManageissuers),
			string(vaults.CertificatePermissionsGetissuers),
			string(vaults.CertificatePermissionsListissuers),
			string(vaults.CertificatePermissionsSetissuers),
			string(vaults.CertificatePermissionsDeleteissuers),
		}
	}

	d.SetId(name)

	if strings.Contains(name, "Key") {
		d.Set("key_permissions", templateManagementPermissions["key"])
	}
	if strings.Contains(name, "Secret") {
		d.Set("secret_permissions", templateManagementPermissions["secret"])
	}
	if strings.Contains(name, "Certificate") {
		d.Set("certificate_permissions", templateManagementPermissions["certificate"])
	}

	return nil
}
