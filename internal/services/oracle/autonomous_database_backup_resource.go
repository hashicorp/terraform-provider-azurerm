// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabasebackups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"time"
)

var _ sdk.Resource = AutonomousDatabaseBackupResource{}

type AutonomousDatabaseBackupResource struct{}

type AutonomousDatabaseBackupResourceModel struct {
	Location          string `tfschema:"location"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	Name              string `tfschema:"display_name"`

	// Required
	AutonomousDataBaseName       string `tfschema:"autonomous_database_name"`
	AutonomousDataBaseBackupOcid string `tfschema:"autonomous_database_backup_ocid"`
	AutonomousDatabaseOcid       string `tfschema:"autonomous_database_ocid"`
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

func (AutonomousDatabaseBackupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		// Required
		"autonomous_database_name": {
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
			Type:         schema.TypeInt,
			Required:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(1, 60),
		},

		// Optional

		// Computed
		"autonomous_database_backup_ocid": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"autonomous_database_ocid": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"backup_type": {
			Type:     schema.TypeString,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(autonomousdatabasebackups.AutonomousDatabaseBackupTypeFull),
				string(autonomousdatabasebackups.AutonomousDatabaseBackupTypeIncremental),
				string(autonomousdatabasebackups.AutonomousDatabaseBackupTypeLongTerm),
			}, false),
		},
		"database_backup_size_in_tbs": {
			Type:         schema.TypeFloat,
			Computed:     true,
			ValidateFunc: validation.IntBetween(1, 384),
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
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabaseBackups
			dbClient := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AutonomousDatabaseBackupResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding model: %+v", err)
			}

			// Check if the autonomous database exists
			dbId := autonomousdatabases.NewAutonomousDatabaseID(
				subscriptionId,
				model.ResourceGroupName,
				model.AutonomousDataBaseName,
			)

			existing, err := dbClient.Get(ctx, dbId)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", dbId, err)
			}

			id := autonomousdatabasebackups.NewAutonomousDatabaseBackupID(
				subscriptionId,
				model.ResourceGroupName,
				model.AutonomousDataBaseName,
				model.Name,
			)

			resp, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(resp.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), &id)
			}

			param := autonomousdatabasebackups.AutonomousDatabaseBackup{

				Properties: &autonomousdatabasebackups.AutonomousDatabaseBackupProperties{
					DisplayName:           pointer.To(model.DisplayName),
					RetentionPeriodInDays: pointer.To(model.RetentionPeriodInDays),
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AutonomousDatabaseBackupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabaseBackups

			id, err := autonomousdatabasebackups.ParseAutonomousDatabaseBackupID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing ID: %+v", err)
			}

			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := AutonomousDatabaseBackupResourceModel{
				Name:                   id.AutonomousDatabaseBackupName,
				ResourceGroupName:      id.ResourceGroupName,
				AutonomousDataBaseName: id.AutonomousDatabaseName,
			}

			if model := result.Model; model != nil {

				if properties := model.Properties; properties != nil {
					state.DisplayName = pointer.From(properties.DisplayName)
					state.RetentionPeriodInDays = pointer.From(properties.RetentionPeriodInDays)
					state.AutonomousDatabaseOcid = pointer.From(properties.AutonomousDatabaseOcid)
					state.AutonomousDataBaseBackupOcid = pointer.From(properties.Ocid)
					state.BackupType = string(pointer.From(properties.BackupType))
					state.DbVersion = pointer.From(properties.DbVersion)
					state.BackupSizeInTbs = int64(pointer.From(properties.DatabaseSizeInTbs))
					state.IsAutomatic = pointer.From(properties.IsAutomatic)
					state.IsRestorable = pointer.From(properties.IsRestorable)
					state.LifecycleDetails = pointer.From(properties.LifecycleDetails)
					state.LifecycleState = string(pointer.From(properties.LifecycleState))
					state.ProvisioningState = string(pointer.From(properties.ProvisioningState))
					state.TimeAvailableTil = pointer.From(properties.TimeAvailableTil)
					state.TimeEnded = pointer.From(properties.TimeEnded)
					state.TimeStarted = pointer.From(properties.TimeStarted)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AutonomousDatabaseBackupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabaseBackups

			id, err := autonomousdatabasebackups.ParseAutonomousDatabaseBackupID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing ID: %+v", err)
			}

			var model AutonomousDatabaseBackupResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			_, err = client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			update := &autonomousdatabasebackups.AutonomousDatabaseBackupUpdate{
				Properties: &autonomousdatabasebackups.AutonomousDatabaseBackupUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("retention_period_in_days") {
				retentionPeriod := int64(model.RetentionPeriodInDays)
				update.Properties.RetentionPeriodInDays = &retentionPeriod
			}

			if err := client.UpdateThenPoll(ctx, *id, *update); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

func (r AutonomousDatabaseBackupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabaseBackups

			id, err := autonomousdatabasebackups.ParseAutonomousDatabaseBackupID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing ID: %+v", err)
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r AutonomousDatabaseBackupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabasebackups.ValidateAutonomousDatabaseBackupID
}
