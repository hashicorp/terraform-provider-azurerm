// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_kubernetes_cluster

type BackupInstanceKubernatesClusterModel struct {
	Name                       string                       `tfschema:"name"`
	Location                   string                       `tfschema:"location"`
	VaultId                    string                       `tfschema:"vault_id"`
	BackupPolicyId             string                       `tfschema:"backup_policy_id"`
	KubernetesClusterId        string                       `tfschema:"kubernetes_cluster_id"`
	SnapshotResourceGroupName  string                       `tfschema:"snapshot_resource_group_name"`
	BackupDatasourceParameters []BackupDatasourceParameters `tfschema:"backup_datasource_parameters"`
	ProtectionState            string                       `tfschema:"protection_state"`
}

func (r Resource) ModelObject() any {
	return &BackupInstanceKubernatesClusterModel{}
}
