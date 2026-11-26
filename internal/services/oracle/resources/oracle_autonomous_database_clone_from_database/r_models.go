// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_clone_from_database

type AutonomousDatabaseCloneFromDatabaseResourceModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`

	// Required for Clone

	SourceAutonomousDatabaseId string `tfschema:"source_autonomous_database_id"`
	CloneType                  string `tfschema:"clone_type"`

	// Required (inherited from base)

	AdminPassword                string   `tfschema:"admin_password"`
	BackupRetentionPeriodInDays  int64    `tfschema:"backup_retention_period_in_days"`
	CharacterSet                 string   `tfschema:"character_set"`
	ComputeCount                 float64  `tfschema:"compute_count"`
	ComputeModel                 string   `tfschema:"compute_model"`
	DataStorageSizeInTb          int64    `tfschema:"data_storage_size_in_tb"`
	DatabaseVersion              string   `tfschema:"database_version"`
	DatabaseWorkload             string   `tfschema:"database_workload"`
	DisplayName                  string   `tfschema:"display_name"`
	LicenseModel                 string   `tfschema:"license_model"`
	AutoScalingEnabled           bool     `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled bool     `tfschema:"auto_scaling_for_storage_enabled"`
	MtlsConnectionRequired       bool     `tfschema:"mtls_connection_required"`
	NationalCharacterSet         string   `tfschema:"national_character_set"`
	SubnetId                     string   `tfschema:"subnet_id"`
	VnetId                       string   `tfschema:"virtual_network_id"`
	AllowedIpAddresses           []string `tfschema:"allowed_ip_addresses"`

	// Optional for Clone

	CustomerContacts []string `tfschema:"customer_contacts"`
	RefreshableModel string   `tfschema:"refreshable_model"`
}

func (Resource) ModelObject() any {
	return &AutonomousDatabaseCloneFromDatabaseResourceModel{}
}
