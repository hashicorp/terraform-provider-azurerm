// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_backups

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

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

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

func (a DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"autonomous_database_id": {
			Description:  "The azureId of the Autonomous Database for which the backups will be listed.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: autonomousdatabases.ValidateAutonomousDatabaseID,
		},
	}
}

func (a DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"autonomous_database_backups": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description: "The ID of the Autonomous Database Backup.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"display_name": {
						Description: "The user-friendly name of the backup.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"location": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"type": {
						Description: "The type of backup.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"retention_period_in_days": {
						Description: "The retention period in days for the Autonomous Database Backup.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"autonomous_database_ocid": {
						Description: "The OCID of the Autonomous Database OCID.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"autonomous_database_backup_ocid": {
						Description: "The backup OCID.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"database_backup_size_in_tbs": {
						Type:     pluginsdk.TypeFloat,
						Computed: true,
					},

					"database_version": {
						Description: "The Oracle Database version of the Autonomous Database at the time the backup was taken.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"automatic": {
						Description: "Whether the backup is user-initiated or automatic.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"restorable": {
						Description: "Whether the backup can be used to restore the Autonomous Database.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"lifecycle_details": {
						Description: "Information about the current lifecycle state of the backup.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"lifecycle_state": {
						Description: "The current state of the backup.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"provisioning_state": {
						Description: "The current provisioning state of the Autonomous Database Backup.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"time_available_til": {
						Description: "The date and time the backup will become unusable.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"time_ended": {
						Description: "The date and time the backup was completed.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"time_started": {
						Description: "The date and time the backup started.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (a DataSource) ModelObject() any {
	return &AutonomousDatabaseBackupsListDataModel{}
}

func (a DataSource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_backups"
}

func (a DataSource) Description() string {
	return "Use this data source to access information about existing Autonomous Database Backups."
}

func (a DataSource) Read() sdk.ResourceFunc {
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

func (a DataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabasebackups.ValidateAutonomousDatabaseBackupID
}
