// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabasebackups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutonomousDatabaseBackupsDataSource struct{}

type AutonomousDatabaseBackupsListDataModel struct {
	AutonomousDatabaseId      string                               `tfschema:"autonomous_database_id"`
	AutonomousDatabaseBackups []AutonomousDatabaseBackupsDataModel `tfschema:"autonomous_database_backups"`
}
type AutonomousDatabaseBackupsDataModel struct {
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

func (a AutonomousDatabaseBackupsDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"autonomous_database_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: autonomousdatabases.ValidateAutonomousDatabaseID,
		},
	}
}

func (a AutonomousDatabaseBackupsDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"autonomous_database_backups": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
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
				},
			},
		},
	}
}

func (a AutonomousDatabaseBackupsDataSource) ModelObject() interface{} {
	return &AutonomousDatabaseBackupsListDataModel{}
}

func (a AutonomousDatabaseBackupsDataSource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_backups"
}

func (a AutonomousDatabaseBackupsDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabaseBackups
			subscriptionId := metadata.Client.Account.SubscriptionId

			state := AutonomousDatabaseBackupsListDataModel{
				AutonomousDatabaseBackups: make([]AutonomousDatabaseBackupsDataModel, 0),
			}
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parsedAdbsId, err := autonomousdatabases.ParseAutonomousDatabaseID(state.AutonomousDatabaseId)
			if err != nil {
				return err
			}
			id := autonomousdatabasebackups.NewAutonomousDatabaseID(subscriptionId, parsedAdbsId.ResourceGroupName, parsedAdbsId.AutonomousDatabaseName)

			resp, err := client.ListByParent(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				for _, element := range *model {
					if props := element.Properties; props != nil {
						backup := AutonomousDatabaseBackupsDataModel{
							Id:                           pointer.From(element.Id),
							DisplayName:                  pointer.From(props.DisplayName),
							RetentionPeriodInDays:        pointer.From(props.RetentionPeriodInDays),
							AutonomousDatabaseOcid:       pointer.From(props.AutonomousDatabaseOcid),
							AutonomousDatabaseBackupOcid: pointer.From(props.Ocid),
							Type:                         pointer.FromEnum(props.BackupType),
							DatabaseVersion:              pointer.From(props.DbVersion),
							BackupSizeInTbs:              pointer.From(props.DatabaseSizeInTbs),
							Automatic:                    pointer.From(props.IsAutomatic),
							Restorable:                   pointer.From(props.IsRestorable),
							LifecycleDetails:             pointer.From(props.LifecycleDetails),
							LifecycleState:               pointer.FromEnum(props.LifecycleState),
							ProvisioningState:            pointer.FromEnum(props.ProvisioningState),
							TimeAvailableTil:             pointer.From(props.TimeAvailableTil),
							TimeEnded:                    pointer.From(props.TimeEnded),
							TimeStarted:                  pointer.From(props.TimeStarted),
						}
						state.AutonomousDatabaseBackups = append(state.AutonomousDatabaseBackups, backup)
					}
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}

func (a AutonomousDatabaseBackupsDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabasebackups.ValidateAutonomousDatabaseBackupID
}
