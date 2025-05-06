// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = AutonomousDatabaseBackupResource{}

type AutonomousDatabaseBackupResource struct{}

type AutonomousDatabaseBackupResourceModel struct {
	Location string `tfschema:"location"`
	Name     string `tfschema:"name"`

	// Required
	AutonomousDataBaseId         string  `tfschema:"autonomous_database_id"`
	AutonomousDataBaseBackupId   string  `tfschema:"autonomous_backup_database_id"`
	BackupRetentionPeriodInDays  int64   `tfschema:"backup_retention_period_in_days"`
	CharacterSet                 string  `tfschema:"character_set"`
	ComputeCount                 float64 `tfschema:"compute_count"`
	ComputeModel                 string  `tfschema:"compute_model"`
	DataStorageSizeInTbs         int64   `tfschema:"data_storage_size_in_tbs"`
	DbVersion                    string  `tfschema:"db_version"`
	DbWorkload                   string  `tfschema:"db_workload"`
	DisplayName                  string  `tfschema:"display_name"`
	LicenseModel                 string  `tfschema:"license_model"`
	AutoScalingEnabled           bool    `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled bool    `tfschema:"auto_scaling_for_storage_enabled"`
	MtlsConnectionRequired       bool    `tfschema:"mtls_connection_required"`
	NationalCharacterSet         string  `tfschema:"national_character_set"`
	SubnetId                     string  `tfschema:"subnet_id"`
	VnetId                       string  `tfschema:"virtual_network_id"`

	// Optional
	CustomerContacts []string `tfschema:"customer_contacts"`
}

func (AutonomousDatabaseBackupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		// Required
		"autonomous_database_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"autonomous_backup_database_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"display_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"retention_period_in_days": {
			Type:     schema.TypeInt,
			Required: true,
			Computed: true,
		},

		// Optional

		// Computed
		"autonomous_database_ocid": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"database_backup_size_in_tbs": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"database_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_automatic": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"is_restorable": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"lifecycle_details": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"lifecycle_state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"provisioning_state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ocid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"time_available_til": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"time_ended": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"time_started": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (r AutonomousDatabaseBackupResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AutonomousDatabaseBackupResource) ModelObject() interface{} {
	return &AutonomousDatabaseBackupResource{}
}

func (r AutonomousDatabaseBackupResource) ResourceType() string {
	return "azurerm_autonomous_database_backup"
}

func (r AutonomousDatabaseBackupResource) Create() sdk.ResourceFunc {
	//TODO implement me
	panic("implement me")
}

func (r AutonomousDatabaseBackupResource) Read() sdk.ResourceFunc {
	//TODO implement me
	panic("implement me")
}

func (r AutonomousDatabaseBackupResource) Delete() sdk.ResourceFunc {
	//TODO implement me
	panic("implement me")
}

func (r AutonomousDatabaseBackupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	//TODO implement me
	panic("implement me")
}
