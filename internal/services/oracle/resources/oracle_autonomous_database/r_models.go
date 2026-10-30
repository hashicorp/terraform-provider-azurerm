// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/helpers"
)

type AutonomousDatabaseRegularResourceModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`

	// Required
	AdminPassword                string                                  `tfschema:"admin_password"`
	BackupRetentionPeriodInDays  int64                                   `tfschema:"backup_retention_period_in_days"`
	CharacterSet                 string                                  `tfschema:"character_set"`
	ComputeCount                 float64                                 `tfschema:"compute_count"`
	ComputeModel                 string                                  `tfschema:"compute_model"`
	DataStorageSizeInTbs         int64                                   `tfschema:"data_storage_size_in_tbs"`
	DbVersion                    string                                  `tfschema:"db_version"`
	DbWorkload                   string                                  `tfschema:"db_workload"`
	DisplayName                  string                                  `tfschema:"display_name"`
	LicenseModel                 string                                  `tfschema:"license_model"`
	LongTermBackUpSchedule       []helpers.LongTermBackUpScheduleDetails `tfschema:"long_term_backup_schedule"`
	AutoScalingEnabled           bool                                    `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled bool                                    `tfschema:"auto_scaling_for_storage_enabled"`
	MtlsConnectionRequired       bool                                    `tfschema:"mtls_connection_required"`
	NationalCharacterSet         string                                  `tfschema:"national_character_set"`
	SubnetId                     string                                  `tfschema:"subnet_id"`
	VnetId                       string                                  `tfschema:"virtual_network_id"`
	AllowedIps                   []string                                `tfschema:"allowed_ips"`
	Ocid                         string                                  `tfschema:"ocid"`

	// Optional
	CustomerContacts []string `tfschema:"customer_contacts"`
}

func (Resource) ModelObject() any {
	return &AutonomousDatabaseRegularResourceModel{}
}
