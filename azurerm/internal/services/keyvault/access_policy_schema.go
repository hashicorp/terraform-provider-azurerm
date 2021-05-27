package keyvault

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/keyvault/mgmt/2020-04-01-preview/keyvault"
	"github.com/gofrs/uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func certificatePermissions() []string {
	return []string{
		"Backup",
		"Create",
		"Delete",
		"DeleteIssuers",
		"Get",
		"GetIssuers",
		"Import",
		"List",
		"ListIssuers",
		"ManageContacts",
		"ManageIssuers",
		"Purge",
		"Recover",
		"Restore",
		"SetIssuers",
		"Update",
	}
}

func flattenCertificatePermission(input string) string {
	for _, permission := range certificatePermissions() {
		if strings.EqualFold(input, permission) {
			return permission
		}
	}

	return input
}

func keyPermissions() []string {
	return []string{
		"Backup",
		"Create",
		"Decrypt",
		"Delete",
		"Encrypt",
		"Get",
		"Import",
		"List",
		"Purge",
		"Recover",
		"Restore",
		"Sign",
		"UnwrapKey",
		"Update",
		"Verify",
		"WrapKey",
	}
}

func flattenKeyPermission(input string) string {
	for _, permission := range keyPermissions() {
		if strings.EqualFold(input, permission) {
			return permission
		}
	}

	return input
}

func secretPermissions() []string {
	return []string{
		"Backup",
		"Delete",
		"Get",
		"List",
		"Purge",
		"Recover",
		"Restore",
		"Set",
	}
}

func flattenSecretPermission(input string) string {
	for _, permission := range secretPermissions() {
		if strings.EqualFold(input, permission) {
			return permission
		}
	}

	return input
}

func storagePermissions() []string {
	return []string{
		"Backup",
		"Delete",
		"DeleteSAS",
		"Get",
		"GetSAS",
		"List",
		"ListSAS",
		"Purge",
		"Recover",
		"RegenerateKey",
		"Restore",
		"Set",
		"SetSAS",
		"Update",
	}
}

func flattenStoragePermission(input string) string {
	for _, permission := range storagePermissions() {
		if strings.EqualFold(input, permission) {
			return permission
		}
	}

	return input
}

func schemaCertificatePermissions() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Schema{
			Type:             pluginsdk.TypeString,
			ValidateFunc:     validation.StringInSlice(certificatePermissions(), true),
			DiffSuppressFunc: suppress.CaseDifference,
		},
	}
}

func schemaKeyPermissions() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Schema{
			Type:             pluginsdk.TypeString,
			ValidateFunc:     validation.StringInSlice(keyPermissions(), true),
			DiffSuppressFunc: suppress.CaseDifference,
		},
	}
}

func schemaSecretPermissions() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Schema{
			Type:             pluginsdk.TypeString,
			ValidateFunc:     validation.StringInSlice(secretPermissions(), true),
			DiffSuppressFunc: suppress.CaseDifference,
		},
	}
}

func schemaStoragePermissions() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Schema{
			Type:             pluginsdk.TypeString,
			ValidateFunc:     validation.StringInSlice(storagePermissions(), true),
			DiffSuppressFunc: suppress.CaseDifference,
		},
	}
}

func expandAccessPolicies(input []interface{}) *[]keyvault.AccessPolicyEntry {
	output := make([]keyvault.AccessPolicyEntry, 0)

	for _, policySet := range input {
		policyRaw := policySet.(map[string]interface{})

		certificatePermissionsRaw := policyRaw["certificate_permissions"].([]interface{})
		keyPermissionsRaw := policyRaw["key_permissions"].([]interface{})
		secretPermissionsRaw := policyRaw["secret_permissions"].([]interface{})
		storagePermissionsRaw := policyRaw["storage_permissions"].([]interface{})

		policy := keyvault.AccessPolicyEntry{
			Permissions: &keyvault.Permissions{
				Certificates: expandCertificatePermissions(certificatePermissionsRaw),
				Keys:         expandKeyPermissions(keyPermissionsRaw),
				Secrets:      expandSecretPermissions(secretPermissionsRaw),
				Storage:      expandStoragePermissions(storagePermissionsRaw),
			},
		}

		tenantUUID := uuid.FromStringOrNil(policyRaw["tenant_id"].(string))
		policy.TenantID = &tenantUUID
		objectUUID := policyRaw["object_id"].(string)
		policy.ObjectID = &objectUUID

		if v := policyRaw["application_id"]; v != "" {
			applicationUUID := uuid.FromStringOrNil(v.(string))
			policy.ApplicationID = &applicationUUID
		}

		output = append(output, policy)
	}

	return &output
}

func flattenAccessPolicies(policies *[]keyvault.AccessPolicyEntry) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	if policies == nil {
		return result
	}

	for _, policy := range *policies {
		policyRaw := make(map[string]interface{})

		if tenantId := policy.TenantID; tenantId != nil {
			policyRaw["tenant_id"] = tenantId.String()
		}

		if objectId := policy.ObjectID; objectId != nil {
			policyRaw["object_id"] = *objectId
		}

		if appId := policy.ApplicationID; appId != nil {
			policyRaw["application_id"] = appId.String()
		}

		if permissions := policy.Permissions; permissions != nil {
			certs := flattenCertificatePermissions(permissions.Certificates)
			policyRaw["certificate_permissions"] = certs

			keys := flattenKeyPermissions(permissions.Keys)
			policyRaw["key_permissions"] = keys

			secrets := flattenSecretPermissions(permissions.Secrets)
			policyRaw["secret_permissions"] = secrets

			storage := flattenStoragePermissions(permissions.Storage)
			policyRaw["storage_permissions"] = storage
		}

		result = append(result, policyRaw)
	}

	return result
}

func expandCertificatePermissions(input []interface{}) *[]keyvault.CertificatePermissions {
	output := make([]keyvault.CertificatePermissions, 0)

	for _, permission := range input {
		output = append(output, keyvault.CertificatePermissions(permission.(string)))
	}

	return &output
}

func flattenCertificatePermissions(input *[]keyvault.CertificatePermissions) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, certificatePermission := range *input {
			permission := flattenCertificatePermission(string(certificatePermission))
			output = append(output, permission)
		}
	}

	return output
}

func expandKeyPermissions(keyPermissionsRaw []interface{}) *[]keyvault.KeyPermissions {
	output := make([]keyvault.KeyPermissions, 0)

	for _, permission := range keyPermissionsRaw {
		output = append(output, keyvault.KeyPermissions(permission.(string)))
	}
	return &output
}

func flattenKeyPermissions(input *[]keyvault.KeyPermissions) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, keyPermission := range *input {
			permission := flattenKeyPermission(string(keyPermission))
			output = append(output, permission)
		}
	}

	return output
}

func expandSecretPermissions(input []interface{}) *[]keyvault.SecretPermissions {
	output := make([]keyvault.SecretPermissions, 0)

	for _, permission := range input {
		output = append(output, keyvault.SecretPermissions(permission.(string)))
	}

	return &output
}

func flattenSecretPermissions(input *[]keyvault.SecretPermissions) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, secretPermission := range *input {
			permission := flattenSecretPermission(string(secretPermission))
			output = append(output, permission)
		}
	}

	return output
}

func expandStoragePermissions(input []interface{}) *[]keyvault.StoragePermissions {
	output := make([]keyvault.StoragePermissions, 0)

	for _, permission := range input {
		output = append(output, keyvault.StoragePermissions(permission.(string)))
	}

	return &output
}

func flattenStoragePermissions(input *[]keyvault.StoragePermissions) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, storagePermission := range *input {
			permission := flattenStoragePermission(string(storagePermission))
			output = append(output, permission)
		}
	}

	return output
}
