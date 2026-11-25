// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_vault_resource_guard_association

type DataProtectionBackupVaultResourceGuardAssociationModel struct {
	DataProtectionBackupVaultId   string `tfschema:"data_protection_backup_vault_id"`
	DataProtectionResourceGuardId string `tfschema:"data_protection_resource_guard_id"`
}

func (r Resource) ModelObject() any {
	return &DataProtectionBackupVaultResourceGuardAssociationModel{}
}
