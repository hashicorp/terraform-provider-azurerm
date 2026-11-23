// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_elastic_san_volume_group

type BackupInstanceElasticSanVolumeGroupModel struct {
	Name                      string `tfschema:"name"`
	Location                  string `tfschema:"location"`
	VaultId                   string `tfschema:"vault_id"`
	BackupPolicyId            string `tfschema:"backup_policy_id"`
	ElasticSanVolumeGroupId   string `tfschema:"elastic_san_volume_group_id"`
	SnapshotResourceGroupName string `tfschema:"snapshot_resource_group_name"`
	VolumeName                string `tfschema:"volume_name"`
	ProtectionState           string `tfschema:"protection_state"`
}

func (r Resource) ModelObject() any {
	return &BackupInstanceElasticSanVolumeGroupModel{}
}
