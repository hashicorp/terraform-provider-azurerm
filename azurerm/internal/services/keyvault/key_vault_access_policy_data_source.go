package keyvault

import (
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/keyvault/mgmt/2020-04-01-preview/keyvault"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
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
	templateManagementPermissions := map[string][]string{
		"key": {
			string(keyvault.KeyPermissionsGet),
			string(keyvault.KeyPermissionsList),
			string(keyvault.KeyPermissionsUpdate),
			string(keyvault.KeyPermissionsCreate),
			string(keyvault.KeyPermissionsImport),
			string(keyvault.KeyPermissionsDelete),
			string(keyvault.KeyPermissionsRecover),
			string(keyvault.KeyPermissionsBackup),
			string(keyvault.KeyPermissionsRestore),
		},
		"secret": {
			string(keyvault.SecretPermissionsGet),
			string(keyvault.SecretPermissionsList),
			string(keyvault.SecretPermissionsSet),
			string(keyvault.SecretPermissionsDelete),
			string(keyvault.SecretPermissionsRecover),
			string(keyvault.SecretPermissionsBackup),
			string(keyvault.SecretPermissionsRestore),
		},
		"certificate": {
			string(keyvault.Get),
			string(keyvault.List),
			string(keyvault.Update),
			string(keyvault.Create),
			string(keyvault.Import),
			string(keyvault.Delete),
			string(keyvault.Managecontacts),
			string(keyvault.Manageissuers),
			string(keyvault.Getissuers),
			string(keyvault.Listissuers),
			string(keyvault.Setissuers),
			string(keyvault.Deleteissuers),
		},
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
