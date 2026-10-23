// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_cosmosdb_account

type BackupInstanceCosmosdbAccountModel struct {
	Name                        string `tfschema:"name"`
	DataProtectionBackupVaultId string `tfschema:"data_protection_backup_vault_id"`
	Location                    string `tfschema:"location"`
	BackupPolicyId              string `tfschema:"backup_policy_cosmosdb_account_id"`
	CosmosdbAccountId           string `tfschema:"cosmosdb_account_id"`
	ProtectionState             string `tfschema:"protection_state"`
}

func (r Resource) ModelObject() any {
	return &BackupInstanceCosmosdbAccountModel{}
}
