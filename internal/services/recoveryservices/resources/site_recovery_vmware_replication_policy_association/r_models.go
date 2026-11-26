// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replication_policy_association

type SiteRecoveryReplicationPolicyVmwareAssociationModel struct {
	Name                        string `tfschema:"name"`
	RecoveryVaultId             string `tfschema:"recovery_vault_id"`
	RecoveryReplicationPolicyId string `tfschema:"policy_id"`
}

func (s Resource) ModelObject() any {
	return &SiteRecoveryReplicationPolicyVmwareAssociationModel{}
}
