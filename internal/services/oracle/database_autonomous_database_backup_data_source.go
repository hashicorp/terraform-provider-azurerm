// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabasebackups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"time"
)

type AutonomousDatabaseBackupDataSource struct{}

type AutonomousDatabaseBackupDataModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`

	AutonomousDataBaseBackupId   string `tfschema:"autonomous_backup_database_id"`
	AutonomousDataBaseId         string `tfschema:"autonomous_database_id"`
	AutonomousDatabaseOcid       string `tfschema:"autonomous_database_ocid"`
	AutonomousDataBaseBackupOcid string `tfschema:"autonomous_database_backup_ocid"`
	BackupType                   string `tfschema:"backup_type"`
	DbVersion                    string `tfschema:"database_version"`
	DisplayName                  string `tfschema:"display_name"`
	BackupSizeInTbs              int64  `tfschema:"database_backup_size_in_tbs"`
	IsAutomatic                  bool   `tfschema:"is_automatic"`
	IsRestorable                 bool   `tfschema:"is_restorable"`
	LifecycleDetails             string `tfschema:"lifecycle_details"`
	LifecycleState               string `tfschema:"lifecycle_state"`
	LicenseModel                 string `tfschema:"license_model"`
	ProvisioningState            string `tfschema:"provisioning_state"`
	RetentionPeriodInDays        int64  `tfschema:"retention_period_in_days"`
	TimeAvailableTil             string `tfschema:"time_available_til"`
	TimeEnded                    string `tfschema:"time_ended"`
	TimeStarted                  string `tfschema:"time_started"`
}

func (a AutonomousDatabaseBackupDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
		},
	}
}

func (a AutonomousDatabaseBackupDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{

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
		"autonomous_database_backup_ocid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"database_backup_size_in_tbs": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"database_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"backup_type": {
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
	}

}

func (a AutonomousDatabaseBackupDataSource) ModelObject() interface{} {
	return &AutonomousDatabaseBackupDataModel{}
}

func (a AutonomousDatabaseBackupDataSource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_backup"
}

func (a AutonomousDatabaseBackupDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabaseBackups
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state AutonomousDatabaseBackupDataModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Parse the autonomous database ID to get the resource group and database name
			dbId, err := autonomousdatabasebackups.ParseAutonomousDatabaseID(state.AutonomousDataBaseId)
			if err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := autonomousdatabasebackups.NewAutonomousDatabaseBackupID(
				subscriptionId,
				state.ResourceGroupName,
				dbId.AutonomousDatabaseName,
				state.Name,
			)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				state.AutonomousDataBaseBackupId = id.ID()

				properties := model.Properties
				if properties != nil {
					state.AutonomousDatabaseOcid = pointer.From(properties.AutonomousDatabaseOcid)
					state.AutonomousDataBaseBackupOcid = pointer.From(properties.Ocid)
					state.BackupType = string(pointer.From((properties.BackupType)))
					state.DbVersion = pointer.From(properties.DbVersion)
					state.DisplayName = pointer.From(properties.DisplayName)
					state.BackupSizeInTbs = int64(pointer.From(properties.DatabaseSizeInTbs))
					state.IsAutomatic = pointer.From(properties.IsAutomatic)
					state.IsRestorable = pointer.From(properties.IsRestorable)
					state.LifecycleDetails = pointer.From(properties.LifecycleDetails)
					state.LifecycleState = string(pointer.From((properties.LifecycleState)))
					state.ProvisioningState = string(pointer.From((properties.ProvisioningState)))
					state.RetentionPeriodInDays = pointer.From((properties.RetentionPeriodInDays))
					state.TimeAvailableTil = pointer.From(properties.TimeAvailableTil)
					state.TimeEnded = pointer.From(properties.TimeEnded)
					state.TimeStarted = pointer.From(properties.TimeStarted)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)

		},
	}
}

func (a AutonomousDatabaseBackupDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabasebackups.ValidateAutonomousDatabaseBackupID
}
