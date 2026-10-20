// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_postgresql_flexible_server

type BackupInstancePostgreSQLFlexibleServerModel struct {
	Name            string `tfschema:"name"`
	Location        string `tfschema:"location"`
	VaultId         string `tfschema:"vault_id"`
	BackupPolicyId  string `tfschema:"backup_policy_id"`
	ServerId        string `tfschema:"server_id"`
	ProtectionState string `tfschema:"protection_state"`
}

func (r Resource) ModelObject() any {
	return &BackupInstancePostgreSQLFlexibleServerModel{}
}
