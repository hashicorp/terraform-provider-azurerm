// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_kubernetes_cluster

type BackupPolicyKubernatesClusterModel struct {
	Name                         string                 `tfschema:"name"`
	ResourceGroupName            string                 `tfschema:"resource_group_name"`
	VaultName                    string                 `tfschema:"vault_name"`
	BackupRepeatingTimeIntervals []string               `tfschema:"backup_repeating_time_intervals"`
	DefaultRetentionRule         []DefaultRetentionRule `tfschema:"default_retention_rule"`
	RetentionRule                []RetentionRule        `tfschema:"retention_rule"`
	TimeZone                     string                 `tfschema:"time_zone"`
}

func (r Resource) ModelObject() any {
	return &BackupPolicyKubernatesClusterModel{}
}
