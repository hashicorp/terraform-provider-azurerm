// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"

	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabasebackups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabases"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutonomousDatabaseBackupDataSource struct{}

type AutonomousDatabaseBackupsDataModel struct {
	AutonomousDatabaseId      string                              `tfschema:"autonomous_database_id"`
	AutonomousDatabaseBackups []AutonomousDatabaseBackupDataModel `tfschema:"autonomous_database_backups"`
}
type AutonomousDatabaseBackupDataModel struct {
	DisplayName                  string  `tfschema:"display_name"`
	BackupType                   string  `tfschema:"backup_type"`
	RetentionPeriodInDays        int64   `tfschema:"retention_period_in_days"`
	AutonomousDatabaseOcid       string  `tfschema:"autonomous_database_ocid"`
	AutonomousDatabaseBackupOcid string  `tfschema:"autonomous_database_backup_ocid"`
	DbVersion                    string  `tfschema:"database_version"`
	BackupSizeInTbs              float64 `tfschema:"database_backup_size_in_tbs"`
	IsAutomatic                  bool    `tfschema:"is_automatic"`
	IsRestorable                 bool    `tfschema:"is_restorable"`
	LifecycleDetails             string  `tfschema:"lifecycle_details"`
	LifecycleState               string  `tfschema:"lifecycle_state"`
	ProvisioningState            string  `tfschema:"provisioning_state"`
	TimeAvailableTil             string  `tfschema:"time_available_til"`
	TimeEnded                    string  `tfschema:"time_ended"`
	TimeStarted                  string  `tfschema:"time_started"`
	Location                     string  `tfschema:"location"`
}

func (a AutonomousDatabaseBackupDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{

		"autonomous_database_id": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func (a AutonomousDatabaseBackupDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"autonomous_database_backups": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"display_name": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"location": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"backup_type": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"retention_period_in_days": {
						Type:     schema.TypeInt,
						Computed: true,
					},

					"autonomous_database_ocid": {
						Type:     schema.TypeString,
						Computed: true,
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
				},
			},
		},
	}
}

func (a AutonomousDatabaseBackupDataSource) ModelObject() interface{} {
	return &AutonomousDatabaseBackupsDataModel{}
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

			state := AutonomousDatabaseBackupsDataModel{
				AutonomousDatabaseBackups: make([]AutonomousDatabaseBackupDataModel, 0),
			}
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parsedAdbsId, err := autonomousdatabases.ParseAutonomousDatabaseID(state.AutonomousDatabaseId)
			if err != nil {
				return fmt.Errorf("decoding id: %+v", err)
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
						backup := AutonomousDatabaseBackupDataModel{
							DisplayName:                  pointer.From(props.DisplayName),
							RetentionPeriodInDays:        pointer.From(props.RetentionPeriodInDays),
							AutonomousDatabaseOcid:       pointer.From(props.AutonomousDatabaseOcid),
							AutonomousDatabaseBackupOcid: pointer.From(props.Ocid),
							BackupType:                   string(pointer.From(props.BackupType)),
							DbVersion:                    pointer.From(props.DbVersion),
							BackupSizeInTbs:              pointer.From(props.DatabaseSizeInTbs),
							IsAutomatic:                  pointer.From(props.IsAutomatic),
							IsRestorable:                 pointer.From(props.IsRestorable),
							LifecycleDetails:             pointer.From(props.LifecycleDetails),
							LifecycleState:               string(pointer.From(props.LifecycleState)),
							ProvisioningState:            string(pointer.From(props.ProvisioningState)),
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

func (a AutonomousDatabaseBackupDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabasebackups.ValidateAutonomousDatabaseBackupID
}
