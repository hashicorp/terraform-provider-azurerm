// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replication_policy

type SiteRecoveryReplicationPolicyVmwareModel struct {
	Name                                            string `tfschema:"name"`
	RecoveryVaultID                                 string `tfschema:"recovery_vault_id"`
	ApplicationConsistentSnapshotFrequencyInMinutes int64  `tfschema:"application_consistent_snapshot_frequency_in_minutes"`
	RecoveryPointRetentionInMinutes                 int64  `tfschema:"recovery_point_retention_in_minutes"`
}

func (r Resource) ModelObject() any {
	return &SiteRecoveryReplicationPolicyVmwareModel{}
}
