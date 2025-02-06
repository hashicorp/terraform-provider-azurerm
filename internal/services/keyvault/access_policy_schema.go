// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// TODO: in the not-too-distant future we can remove this logic
// however we need to determine the correct value to normalize too

func certificatePermissions() []string {
	return append(certificatePermissionsManagement(), certificatePermissionsPrivileged()...)
}

func certificatePermissionsManagement() []string {
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

func certificatePermissionsPrivileged() []string {
	return []string{
		"Purge",
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
	permissions := keyPermissionsManagement()
	permissions = append(permissions, keyPermissionsCryptographic()...)
	permissions = append(permissions, keyPermissionsPrivileged()...)
	permissions = append(permissions, keyPermissionsRotationPolicy()...)
	return permissions
}

func keyPermissionsManagement() []string {
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

func keyPermissionsCryptographic() []string {
	return []string{
		"Decrypt",
		"Encrypt",
		"UnwrapKey",
		"WrapKey",
		"Verify",
		"Sign",
	}
}

func keyPermissionsPrivileged() []string {
	return []string{
		"Purge",
		"Release",
	}
}

func keyPermissionsRotationPolicy() []string {
	return []string{
		"Rotate",
		"GetRotationPolicy",
		"SetRotationPolicy",
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
	return append(secretPermissionsManagement(), secretPermissionsPrivileged()...)
}

func secretPermissionsManagement() []string {
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

func secretPermissionsPrivileged() []string {
	return []string{
		"Purge",
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
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(certificatePermissions(), false),
		},
	}
}

func schemaKeyPermissions() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(keyPermissions(), false),
		},
	}
}

func schemaSecretPermissions() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(secretPermissions(), false),
		},
	}
}

func schemaStoragePermissions() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(storagePermissions(), false),
		},
	}
}

func expandAccessPolicies(input []interface{}) *[]vaults.AccessPolicyEntry {
	output := make([]vaults.AccessPolicyEntry, 0)

	for _, policySet := range input {
		policyRaw := policySet.(map[string]interface{})

		certificatePermissionsRaw := policyRaw["certificate_permissions"].([]interface{})
		keyPermissionsRaw := policyRaw["key_permissions"].([]interface{})
		secretPermissionsRaw := policyRaw["secret_permissions"].([]interface{})
		storagePermissionsRaw := policyRaw["storage_permissions"].([]interface{})

		policy := vaults.AccessPolicyEntry{
			Permissions: vaults.Permissions{
				Certificates: expandCertificatePermissions(certificatePermissionsRaw),
				Keys:         expandKeyPermissions(keyPermissionsRaw),
				Secrets:      expandSecretPermissions(secretPermissionsRaw),
				Storage:      expandStoragePermissions(storagePermissionsRaw),
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

func flattenAccessPolicies(input *[]vaults.AccessPolicyEntry) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	if input == nil {
		return result
	}

	for _, policy := range *input {
		applicationId := ""
		if policy.ApplicationId != nil {
			applicationId = *policy.ApplicationId
		}

		certs := flattenCertificatePermissions(policy.Permissions.Certificates)
		keys := flattenKeyPermissions(policy.Permissions.Keys)
		secrets := flattenSecretPermissions(policy.Permissions.Secrets)
		storage := flattenStoragePermissions(policy.Permissions.Storage)
		result = append(result, map[string]interface{}{
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

func expandCertificatePermissions(input []interface{}) *[]vaults.CertificatePermissions {
	output := make([]vaults.CertificatePermissions, 0)

	for _, permission := range input {
		output = append(output, vaults.CertificatePermissions(permission.(string)))
	}

	return &output
}

func flattenCertificatePermissions(input *[]vaults.CertificatePermissions) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, certificatePermission := range *input {
			permission := flattenCertificatePermission(string(certificatePermission))
			output = append(output, permission)
		}
	}

	return output
}

func expandKeyPermissions(keyPermissionsRaw []interface{}) *[]vaults.KeyPermissions {
	output := make([]vaults.KeyPermissions, 0)

	for _, permission := range keyPermissionsRaw {
		output = append(output, vaults.KeyPermissions(permission.(string)))
	}
	return &output
}

func flattenKeyPermissions(input *[]vaults.KeyPermissions) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, keyPermission := range *input {
			permission := flattenKeyPermission(string(keyPermission))
			output = append(output, permission)
		}
	}

	return output
}

func expandSecretPermissions(input []interface{}) *[]vaults.SecretPermissions {
	output := make([]vaults.SecretPermissions, 0)

	for _, permission := range input {
		output = append(output, vaults.SecretPermissions(permission.(string)))
	}

	return &output
}

func flattenSecretPermissions(input *[]vaults.SecretPermissions) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, secretPermission := range *input {
			permission := flattenSecretPermission(string(secretPermission))
			output = append(output, permission)
		}
	}

	return output
}

func expandStoragePermissions(input []interface{}) *[]vaults.StoragePermissions {
	output := make([]vaults.StoragePermissions, 0)

	for _, permission := range input {
		output = append(output, vaults.StoragePermissions(permission.(string)))
	}

	return &output
}

func flattenStoragePermissions(input *[]vaults.StoragePermissions) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, storagePermission := range *input {
			permission := flattenStoragePermission(string(storagePermission))
			output = append(output, permission)
		}
	}

	return output
}
