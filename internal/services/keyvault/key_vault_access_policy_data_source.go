package keyvault

import (
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2021-10-01/keyvault" // nolint: staticcheck
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
			azure.TitleCase(string(keyvault.KeyPermissionsGet)),
			azure.TitleCase(string(keyvault.KeyPermissionsList)),
			azure.TitleCase(string(keyvault.KeyPermissionsUpdate)),
			azure.TitleCase(string(keyvault.KeyPermissionsCreate)),
			azure.TitleCase(string(keyvault.KeyPermissionsImport)),
			azure.TitleCase(string(keyvault.KeyPermissionsDelete)),
			azure.TitleCase(string(keyvault.KeyPermissionsRecover)),
			azure.TitleCase(string(keyvault.KeyPermissionsBackup)),
			azure.TitleCase(string(keyvault.KeyPermissionsRestore)),
		}
		templateManagementPermissions["secret"] = []string{
			azure.TitleCase(string(keyvault.SecretPermissionsGet)),
			azure.TitleCase(string(keyvault.SecretPermissionsList)),
			azure.TitleCase(string(keyvault.SecretPermissionsSet)),
			azure.TitleCase(string(keyvault.SecretPermissionsDelete)),
			azure.TitleCase(string(keyvault.SecretPermissionsRecover)),
			azure.TitleCase(string(keyvault.SecretPermissionsBackup)),
			azure.TitleCase(string(keyvault.SecretPermissionsRestore)),
		}
		templateManagementPermissions["certificate"] = []string{
			azure.TitleCase(string(keyvault.CertificatePermissionsGet)),
			azure.TitleCase(string(keyvault.CertificatePermissionsList)),
			azure.TitleCase(string(keyvault.CertificatePermissionsUpdate)),
			azure.TitleCase(string(keyvault.CertificatePermissionsCreate)),
			azure.TitleCase(string(keyvault.CertificatePermissionsImport)),
			azure.TitleCase(string(keyvault.CertificatePermissionsDelete)),
			azure.TitleCase(string(keyvault.CertificatePermissionsManagecontacts)),
			azure.TitleCase(string(keyvault.CertificatePermissionsManageissuers)),
			azure.TitleCase(string(keyvault.CertificatePermissionsGetissuers)),
			azure.TitleCase(string(keyvault.CertificatePermissionsListissuers)),
			azure.TitleCase(string(keyvault.CertificatePermissionsSetissuers)),
			azure.TitleCase(string(keyvault.CertificatePermissionsDeleteissuers)),
		}
	} else {
		templateManagementPermissions["key"] = []string{
			string(keyvault.KeyPermissionsGet),
			string(keyvault.KeyPermissionsList),
			string(keyvault.KeyPermissionsUpdate),
			string(keyvault.KeyPermissionsCreate),
			string(keyvault.KeyPermissionsImport),
			string(keyvault.KeyPermissionsDelete),
			string(keyvault.KeyPermissionsRecover),
			string(keyvault.KeyPermissionsBackup),
			string(keyvault.KeyPermissionsRestore),
		}
		templateManagementPermissions["secret"] = []string{
			string(keyvault.SecretPermissionsGet),
			string(keyvault.SecretPermissionsList),
			string(keyvault.SecretPermissionsSet),
			string(keyvault.SecretPermissionsDelete),
			string(keyvault.SecretPermissionsRecover),
			string(keyvault.SecretPermissionsBackup),
			string(keyvault.SecretPermissionsRestore),
		}
		templateManagementPermissions["certificate"] = []string{
			string(keyvault.CertificatePermissionsGet),
			string(keyvault.CertificatePermissionsList),
			string(keyvault.CertificatePermissionsUpdate),
			string(keyvault.CertificatePermissionsCreate),
			string(keyvault.CertificatePermissionsImport),
			string(keyvault.CertificatePermissionsDelete),
			string(keyvault.CertificatePermissionsManagecontacts),
			string(keyvault.CertificatePermissionsManageissuers),
			string(keyvault.CertificatePermissionsGetissuers),
			string(keyvault.CertificatePermissionsListissuers),
			string(keyvault.CertificatePermissionsSetissuers),
			string(keyvault.CertificatePermissionsDeleteissuers),
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
