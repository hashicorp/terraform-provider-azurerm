// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabasebackups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutonomousDatabaseBackupDataSource struct{}

type AutonomousDatabaseBackupDataModel struct {
	AutonomousDatabaseId string `tfschema:"autonomous_database_id"`
	Name                 string `tfschema:"name"`

	Id                           string  `tfschema:"id"`
	DisplayName                  string  `tfschema:"display_name"`
	Type                         string  `tfschema:"type"`
	RetentionPeriodInDays        int64   `tfschema:"retention_period_in_days"`
	AutonomousDatabaseOcid       string  `tfschema:"autonomous_database_ocid"`
	AutonomousDatabaseBackupOcid string  `tfschema:"autonomous_database_backup_ocid"`
	DatabaseVersion              string  `tfschema:"database_version"`
	BackupSizeInTbs              float64 `tfschema:"database_backup_size_in_tbs"`
	Automatic                    bool    `tfschema:"automatic"`
	Restorable                   bool    `tfschema:"restorable"`
	LifecycleDetails             string  `tfschema:"lifecycle_details"`
	LifecycleState               string  `tfschema:"lifecycle_state"`
	ProvisioningState            string  `tfschema:"provisioning_state"`
	TimeAvailableTil             string  `tfschema:"time_available_til"`
	TimeEnded                    string  `tfschema:"time_ended"`
	TimeStarted                  string  `tfschema:"time_started"`
	Location                     string  `tfschema:"location"`
}

func (a AutonomousDatabaseBackupDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"autonomous_database_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: autonomousdatabases.ValidateAutonomousDatabaseID,
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
		},
	}
}

func (a AutonomousDatabaseBackupDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"retention_period_in_days": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"autonomous_database_ocid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"autonomous_database_backup_ocid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"database_backup_size_in_tbs": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},

		"database_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"automatic": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"restorable": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"lifecycle_details": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"lifecycle_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_available_til": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_ended": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_started": {
			Type:     pluginsdk.TypeString,
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
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabaseBackups
			subscriptionId := metadata.Client.Account.SubscriptionId

			state := AutonomousDatabaseBackupDataModel{}
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parsedAdbsId, err := autonomousdatabases.ParseAutonomousDatabaseID(state.AutonomousDatabaseId)
			if err != nil {
				return err
			}

			dbId := autonomousdatabases.NewAutonomousDatabaseID(subscriptionId, parsedAdbsId.ResourceGroupName, parsedAdbsId.AutonomousDatabaseName)

			id := autonomousdatabasebackups.NewAutonomousDatabaseBackupID(subscriptionId, parsedAdbsId.ResourceGroupName, parsedAdbsId.AutonomousDatabaseName, state.Name)

			resp, err := getBackupFromOCI(ctx, client, dbId, id)
			if err != nil {
				return fmt.Errorf("retrieving backup: %+v", err)
			}

			if resp != nil {
				state.Id = pointer.From(resp.Id)
				if props := resp.Properties; props != nil {
					state.DisplayName = pointer.From(props.DisplayName)
					state.RetentionPeriodInDays = pointer.From(props.RetentionPeriodInDays)
					state.AutonomousDatabaseOcid = pointer.From(props.AutonomousDatabaseOcid)
					state.AutonomousDatabaseBackupOcid = pointer.From(props.Ocid)
					state.Type = pointer.FromEnum(props.BackupType)
					state.DatabaseVersion = pointer.From(props.DbVersion)
					state.BackupSizeInTbs = pointer.From(props.DatabaseSizeInTbs)
					state.Automatic = pointer.From(props.IsAutomatic)
					state.Restorable = pointer.From(props.IsRestorable)
					state.LifecycleDetails = pointer.From(props.LifecycleDetails)
					state.LifecycleState = pointer.FromEnum(props.LifecycleState)
					state.ProvisioningState = pointer.FromEnum(props.ProvisioningState)
					state.TimeAvailableTil = pointer.From(props.TimeAvailableTil)
					state.TimeEnded = pointer.From(props.TimeEnded)
					state.TimeStarted = pointer.From(props.TimeStarted)
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
