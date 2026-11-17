// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_data_lake_storage

type BackupInstanceDataLakeStorageModel struct {
	Name                        string   `tfschema:"name"`
	DataProtectionBackupVaultId string   `tfschema:"data_protection_backup_vault_id"`
	Location                    string   `tfschema:"location"`
	BackupPolicyId              string   `tfschema:"backup_policy_data_lake_storage_id"`
	StorageAccountId            string   `tfschema:"storage_account_id"`
	StorageContainerNames       []string `tfschema:"storage_container_names"`
	ProtectionState             string   `tfschema:"protection_state"`
}

func (r Resource) ModelObject() any {
	return &BackupInstanceDataLakeStorageModel{}
}
