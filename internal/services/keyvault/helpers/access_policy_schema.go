// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// TODO: in the not-too-distant future we can remove this logic
// however we need to determine the correct value to normalize too

func CertificatePermissions() []string {
	return append(CertificatePermissionsManagement(), CertificatePermissionsPrivileged()...)
}

func CertificatePermissionsManagement() []string {
	return []string{
		"Get",
		"List",
		"Update",
		"Create",
		"Import",
		"Delete",
		"Recover",
		"Backup",
		"Restore",
		"ManageContacts",
		"ManageIssuers",
		"GetIssuers",
		"ListIssuers",
		"SetIssuers",
		"DeleteIssuers",
	}
}

func CertificatePermissionsPrivileged() []string {
	return []string{
		"Purge",
	}
}

func FlattenCertificatePermission(input string) string {
	for _, permission := range CertificatePermissions() {
		if strings.EqualFold(input, permission) {
			return permission
		}
	}

	return input
}

func KeyPermissions() []string {
	permissions := KeyPermissionsManagement()
	permissions = append(permissions, KeyPermissionsCryptographic()...)
	permissions = append(permissions, KeyPermissionsPrivileged()...)
	permissions = append(permissions, KeyPermissionsRotationPolicy()...)
	return permissions
}

func KeyPermissionsManagement() []string {
	return []string{
		"Get",
		"List",
		"Update",
		"Create",
		"Import",
		"Delete",
		"Recover",
		"Backup",
		"Restore",
	}
}

func KeyPermissionsCryptographic() []string {
	return []string{
		"Decrypt",
		"Encrypt",
		"UnwrapKey",
		"WrapKey",
		"Verify",
		"Sign",
	}
}

func KeyPermissionsPrivileged() []string {
	return []string{
		"Purge",
		"Release",
	}
}

func KeyPermissionsRotationPolicy() []string {
	return []string{
		"Rotate",
		"GetRotationPolicy",
		"SetRotationPolicy",
	}
}

func FlattenKeyPermission(input string) string {
	for _, permission := range KeyPermissions() {
		if strings.EqualFold(input, permission) {
			return permission
		}
	}

	return input
}

func SecretPermissions() []string {
	return append(SecretPermissionsManagement(), SecretPermissionsPrivileged()...)
}

func SecretPermissionsManagement() []string {
	return []string{
		"Get",
		"List",
		"Set",
		"Delete",
		"Recover",
		"Backup",
		"Restore",
	}
}

func SecretPermissionsPrivileged() []string {
	return []string{
		"Purge",
	}
}

func FlattenSecretPermission(input string) string {
	for _, permission := range SecretPermissions() {
		if strings.EqualFold(input, permission) {
			return permission
		}
	}

	return input
}

func StoragePermissions() []string {
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

func FlattenStoragePermission(input string) string {
	for _, permission := range StoragePermissions() {
		if strings.EqualFold(input, permission) {
			return permission
		}
	}

	return input
}

func SchemaCertificatePermissions() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(CertificatePermissions(), false),
		},
	}
}

func SchemaKeyPermissions() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(KeyPermissions(), false),
		},
	}
}

func SchemaSecretPermissions() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(SecretPermissions(), false),
		},
	}
}

func SchemaStoragePermissions() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(StoragePermissions(), false),
		},
	}
}

func ExpandAccessPolicies(input []any) *[]vaults.AccessPolicyEntry {
	output := make([]vaults.AccessPolicyEntry, 0)

	for _, policySet := range input {
		policyRaw := policySet.(map[string]any)

		certificatePermissionsRaw := policyRaw["certificate_permissions"].([]any)
		keyPermissionsRaw := policyRaw["key_permissions"].([]any)
		secretPermissionsRaw := policyRaw["secret_permissions"].([]any)
		storagePermissionsRaw := policyRaw["storage_permissions"].([]any)

		policy := vaults.AccessPolicyEntry{
			Permissions: vaults.Permissions{
				Certificates: ExpandCertificatePermissions(certificatePermissionsRaw),
				Keys:         ExpandKeyPermissions(keyPermissionsRaw),
				Secrets:      ExpandSecretPermissions(secretPermissionsRaw),
				Storage:      ExpandStoragePermissions(storagePermissionsRaw),
			},
			ObjectId: policyRaw["object_id"].(string),
			TenantId: policyRaw["tenant_id"].(string),
		}
		if v := policyRaw["application_id"]; v != "" {
			policy.ApplicationId = pointer.To(v.(string))
		}
		output = append(output, policy)
	}

	return &output
}

func FlattenAccessPolicies(input *[]vaults.AccessPolicyEntry) []map[string]any {
	result := make([]map[string]any, 0)

	if input == nil {
		return result
	}

	for _, policy := range *input {
		applicationId := pointer.From(policy.ApplicationId)

		certs := FlattenCertificatePermissions(policy.Permissions.Certificates)
		keys := FlattenKeyPermissions(policy.Permissions.Keys)
		secrets := FlattenSecretPermissions(policy.Permissions.Secrets)
		storage := FlattenStoragePermissions(policy.Permissions.Storage)
		result = append(result, map[string]any{
			"application_id":          applicationId,
			"certificate_permissions": certs,
			"key_permissions":         keys,
			"object_id":               policy.ObjectId,
			"secret_permissions":      secrets,
			"storage_permissions":     storage,
			"tenant_id":               policy.TenantId,
		})
	}

	return result
}

func ExpandCertificatePermissions(input []any) *[]vaults.CertificatePermissions {
	output := make([]vaults.CertificatePermissions, 0)

	for _, permission := range input {
		output = append(output, vaults.CertificatePermissions(permission.(string)))
	}

	return &output
}

func FlattenCertificatePermissions(input *[]vaults.CertificatePermissions) []any {
	output := make([]any, 0)

	if input != nil {
		for _, certificatePermission := range *input {
			permission := FlattenCertificatePermission(string(certificatePermission))
			output = append(output, permission)
		}
	}

	return output
}

func ExpandKeyPermissions(keyPermissionsRaw []any) *[]vaults.KeyPermissions {
	output := make([]vaults.KeyPermissions, 0)

	for _, permission := range keyPermissionsRaw {
		output = append(output, vaults.KeyPermissions(permission.(string)))
	}
	return &output
}

func FlattenKeyPermissions(input *[]vaults.KeyPermissions) []any {
	output := make([]any, 0)

	if input != nil {
		for _, keyPermission := range *input {
			permission := FlattenKeyPermission(string(keyPermission))
			output = append(output, permission)
		}
	}

	return output
}

func ExpandSecretPermissions(input []any) *[]vaults.SecretPermissions {
	output := make([]vaults.SecretPermissions, 0)

	for _, permission := range input {
		output = append(output, vaults.SecretPermissions(permission.(string)))
	}

	return &output
}

func FlattenSecretPermissions(input *[]vaults.SecretPermissions) []any {
	output := make([]any, 0)

	if input != nil {
		for _, secretPermission := range *input {
			permission := FlattenSecretPermission(string(secretPermission))
			output = append(output, permission)
		}
	}

	return output
}

func ExpandStoragePermissions(input []any) *[]vaults.StoragePermissions {
	output := make([]vaults.StoragePermissions, 0)

	for _, permission := range input {
		output = append(output, vaults.StoragePermissions(permission.(string)))
	}

	return &output
}

func FlattenStoragePermissions(input *[]vaults.StoragePermissions) []any {
	output := make([]any, 0)

	if input != nil {
		for _, storagePermission := range *input {
			permission := FlattenStoragePermission(string(storagePermission))
			output = append(output, permission)
		}
	}

	return output
}
