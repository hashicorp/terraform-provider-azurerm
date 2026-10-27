// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_vault_customer_managed_key

type DataProtectionBackupVaultCustomerManagedKeyModel struct {
	DataProtectionBackupVaultID string `tfschema:"data_protection_backup_vault_id"`
	KeyVaultKeyID               string `tfschema:"key_vault_key_id"`
}

func (r Resource) ModelObject() any {
	return &Resource{}
}
