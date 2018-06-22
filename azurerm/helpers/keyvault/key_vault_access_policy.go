package keyvault

import (
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"strings"
)

func KeyPermissionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
			ValidateFunc: validation.StringInSlice([]string{
				string(keyvault.KeyPermissionsBackup),
				string(keyvault.KeyPermissionsCreate),
				string(keyvault.KeyPermissionsDecrypt),
				string(keyvault.KeyPermissionsDelete),
				string(keyvault.KeyPermissionsEncrypt),
				string(keyvault.KeyPermissionsGet),
				string(keyvault.KeyPermissionsImport),
				string(keyvault.KeyPermissionsList),
				string(keyvault.KeyPermissionsPurge),
				string(keyvault.KeyPermissionsRecover),
				string(keyvault.KeyPermissionsRestore),
				string(keyvault.KeyPermissionsSign),
				string(keyvault.KeyPermissionsUnwrapKey),
				string(keyvault.KeyPermissionsUpdate),
				string(keyvault.KeyPermissionsVerify),
				string(keyvault.KeyPermissionsWrapKey),
			}, true),
			DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
		},
	}
}

func SecretPermissionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
			ValidateFunc: validation.StringInSlice([]string{
				string(keyvault.SecretPermissionsBackup),
				string(keyvault.SecretPermissionsDelete),
				string(keyvault.SecretPermissionsGet),
				string(keyvault.SecretPermissionsList),
				string(keyvault.SecretPermissionsPurge),
				string(keyvault.SecretPermissionsRecover),
				string(keyvault.SecretPermissionsRestore),
				string(keyvault.SecretPermissionsSet),
			}, true),
			DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
		},
	}
}

func CertificatePermissionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
			ValidateFunc: validation.StringInSlice([]string{
				string(keyvault.Create),
				string(keyvault.Delete),
				string(keyvault.Deleteissuers),
				string(keyvault.Get),
				string(keyvault.Getissuers),
				string(keyvault.Import),
				string(keyvault.List),
				string(keyvault.Listissuers),
				string(keyvault.Managecontacts),
				string(keyvault.Manageissuers),
				string(keyvault.Purge),
				string(keyvault.Recover),
				string(keyvault.Setissuers),
				string(keyvault.Update),
			}, true),
			DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
		},
	}
}

func FlattenKeyVaultAccessPolicies(policies *[]keyvault.AccessPolicyEntry) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(*policies))

	for _, policy := range *policies {
		policyRaw := make(map[string]interface{})

		keyPermissionsRaw := make([]interface{}, 0, len(*policy.Permissions.Keys))
		for _, keyPermission := range *policy.Permissions.Keys {
			keyPermissionsRaw = append(keyPermissionsRaw, string(keyPermission))
		}

		secretPermissionsRaw := make([]interface{}, 0, len(*policy.Permissions.Secrets))
		for _, secretPermission := range *policy.Permissions.Secrets {
			secretPermissionsRaw = append(secretPermissionsRaw, string(secretPermission))
		}

		policyRaw["tenant_id"] = policy.TenantID.String()
		policyRaw["object_id"] = *policy.ObjectID
		if policy.ApplicationID != nil {
			policyRaw["application_id"] = policy.ApplicationID.String()
		}
		policyRaw["key_permissions"] = keyPermissionsRaw
		policyRaw["secret_permissions"] = secretPermissionsRaw

		if policy.Permissions.Certificates != nil {
			certificatePermissionsRaw := make([]interface{}, 0, len(*policy.Permissions.Certificates))
			for _, certificatePermission := range *policy.Permissions.Certificates {
				certificatePermissionsRaw = append(certificatePermissionsRaw, string(certificatePermission))
			}
			policyRaw["certificate_permissions"] = certificatePermissionsRaw
		}

		result = append(result, policyRaw)
	}

	return result
}

func ExpandKeyVaultAccessPolicyCertificatePermissions(certificatePermissionsRaw []interface{}) *[]keyvault.CertificatePermissions {
	certificatePermissions := []keyvault.CertificatePermissions{}

	for _, permission := range certificatePermissionsRaw {
		certificatePermissions = append(certificatePermissions, keyvault.CertificatePermissions(permission.(string)))
	}
	return &certificatePermissions
}

func ExpandKeyVaultAccessPolicyKeyPermissions(keyPermissionsRaw []interface{}) *[]keyvault.KeyPermissions {
	keyPermissions := []keyvault.KeyPermissions{}

	for _, permission := range keyPermissionsRaw {
		keyPermissions = append(keyPermissions, keyvault.KeyPermissions(permission.(string)))
	}
	return &keyPermissions
}

func ExpandKeyVaultAccessPolicySecretPermissions(secretPermissionsRaw []interface{}) *[]keyvault.SecretPermissions {
	secretPermissions := []keyvault.SecretPermissions{}

	for _, permission := range secretPermissionsRaw {
		secretPermissions = append(secretPermissions, keyvault.SecretPermissions(permission.(string)))
	}
	return &secretPermissions
}

func ignoreCaseDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return strings.ToLower(old) == strings.ToLower(new)
}
