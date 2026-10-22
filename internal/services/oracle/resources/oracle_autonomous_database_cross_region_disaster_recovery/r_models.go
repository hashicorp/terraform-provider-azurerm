// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_cross_region_disaster_recovery

type AutonomousDatabaseCrossRegionDisasterRecoveryResourceModel struct {
	Name                             string            `tfschema:"name"`
	ResourceGroupName                string            `tfschema:"resource_group_name"`
	Location                         string            `tfschema:"location"`
	DisplayName                      string            `tfschema:"display_name"`
	SourceAutonomousDatabaseId       string            `tfschema:"source_autonomous_database_id"`
	SubnetId                         string            `tfschema:"subnet_id"`
	ReplicateAutomaticBackupsEnabled bool              `tfschema:"replicate_automatic_backups_enabled"`
	Tags                             map[string]string `tfschema:"tags"`
	AutoScalingEnabled               bool              `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled     bool              `tfschema:"auto_scaling_for_storage_enabled"`
	BackupRetentionPeriodInDays      int64             `tfschema:"backup_retention_period_in_days"`
	CharacterSet                     string            `tfschema:"character_set"`
	ComputeCount                     float64           `tfschema:"compute_count"`
	ComputeModel                     string            `tfschema:"compute_model"`
	CustomerContacts                 []string          `tfschema:"customer_contacts"`
	DataStorageSizeInTb              int64             `tfschema:"data_storage_size_in_tb"`
	DatabaseVersion                  string            `tfschema:"database_version"`
	DatabaseWorkload                 string            `tfschema:"database_workload"`
	LicenseModel                     string            `tfschema:"license_model"`
	MtlsConnectionRequired           bool              `tfschema:"mtls_connection_required"`
	NationalCharacterSet             string            `tfschema:"national_character_set"`
	RemoteDisasterRecoveryType       string            `tfschema:"remote_disaster_recovery_type"`
}

func (Resource) ModelObject() any {
	return &AutonomousDatabaseCrossRegionDisasterRecoveryResourceModel{}
}
