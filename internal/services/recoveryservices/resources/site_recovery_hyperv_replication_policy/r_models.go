// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_hyperv_replication_policy

type ReplicationPolicyHyperVModel struct {
	Name                                          string `tfschema:"name"`
	RecoveryVaultId                               string `tfschema:"recovery_vault_id"`
	RecoveryPointRetentionInHours                 int64  `tfschema:"recovery_point_retention_in_hours"`
	ApplicationConsistentSnapshotFrequencyInHours int64  `tfschema:"application_consistent_snapshot_frequency_in_hours"`
	CopyFrequency                                 int64  `tfschema:"replication_interval_in_seconds"`
}

func (r Resource) ModelObject() any {
	return &ReplicationPolicyHyperVModel{}
}
