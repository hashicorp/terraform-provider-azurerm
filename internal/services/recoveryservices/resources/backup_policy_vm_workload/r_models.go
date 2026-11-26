// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_vm_workload

type BackupProtectionPolicyVMWorkloadModel struct {
	Name               string             `tfschema:"name"`
	ResourceGroupName  string             `tfschema:"resource_group_name"`
	RecoveryVaultName  string             `tfschema:"recovery_vault_name"`
	ProtectionPolicies []ProtectionPolicy `tfschema:"protection_policy"`
	Settings           []Settings         `tfschema:"settings"`
	WorkloadType       string             `tfschema:"workload_type"`
}

func (r Resource) ModelObject() any {
	return &BackupProtectionPolicyVMWorkloadModel{}
}
