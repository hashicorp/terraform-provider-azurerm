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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = AutonomousDatabaseBackupResource{}

type AutonomousDatabaseBackupResource struct{}

type AutonomousDatabaseBackupResourceModel struct {
	AutonomousDatabaseId string `tfschema:"autonomous_database_id"`
	Name                 string `tfschema:"name"`

	// Required
	Type                  string `tfschema:"type"`
	RetentionPeriodInDays int64  `tfschema:"retention_period_in_days"`
}

func (AutonomousDatabaseBackupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"autonomous_database_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: autonomousdatabases.ValidateAutonomousDatabaseID,
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
		},
		"retention_period_in_days": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(90, 3650),
		},

		// Optional
		"type": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(autonomousdatabasebackups.AutonomousDatabaseBackupTypeLongTerm),
			ValidateFunc: validation.StringInSlice([]string{
				string(autonomousdatabasebackups.AutonomousDatabaseBackupTypeLongTerm),
			}, false),
		},
	}
}

func (r AutonomousDatabaseBackupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AutonomousDatabaseBackupResource) ModelObject() interface{} {
	return &AutonomousDatabaseBackupResourceModel{}
}

func (r AutonomousDatabaseBackupResource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_backup"
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

			parsedAdbsId, err := autonomousdatabases.ParseAutonomousDatabaseID(model.AutonomousDatabaseId)
			if err != nil {
				return err
			}

			dbId := autonomousdatabases.NewAutonomousDatabaseID(subscriptionId, parsedAdbsId.ResourceGroupName, parsedAdbsId.AutonomousDatabaseName)

			existing, err := dbClient.Get(ctx, dbId)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", dbId, err)
			}

			id := autonomousdatabasebackups.NewAutonomousDatabaseBackupID(
				subscriptionId,
				parsedAdbsId.ResourceGroupName,
				parsedAdbsId.AutonomousDatabaseName,
				model.Name,
			)

			existingBackup, err := getBackupFromOCI(ctx, client, dbId, id)
			if err != nil {
				return fmt.Errorf("checking for existing backup: %+v", err)
			}
			if existingBackup != nil {
				return metadata.ResourceRequiresImport(r.ResourceType(), &id)
			}

			param := autonomousdatabasebackups.AutonomousDatabaseBackup{
				Name: pointer.To(model.Name),
				Properties: &autonomousdatabasebackups.AutonomousDatabaseBackupProperties{
					RetentionPeriodInDays: pointer.To(model.RetentionPeriodInDays),
					BackupType:            pointer.To(autonomousdatabasebackups.AutonomousDatabaseBackupType(model.Type)),
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
				return err
			}

			adbId := autonomousdatabases.NewAutonomousDatabaseID(
				id.SubscriptionId,
				id.ResourceGroupName,
				id.AutonomousDatabaseName,
			)
			backupId := autonomousdatabasebackups.NewAutonomousDatabaseBackupID(id.SubscriptionId, id.ResourceGroupName, id.AutonomousDatabaseName, id.AutonomousDatabaseBackupName)

			backup, err := getBackupFromOCI(ctx, client, adbId, backupId)
			if err != nil {
				return fmt.Errorf("retrieving backup: %+v", err)
			}

			if backup == nil {
				err := metadata.MarkAsGone(id)
				if err != nil {
					return err
				}
				return nil
			}

			state := AutonomousDatabaseBackupResourceModel{
				Name:                 id.AutonomousDatabaseBackupName,
				AutonomousDatabaseId: adbId.ID(),
			}

			if props := backup.Properties; props != nil {
				state.RetentionPeriodInDays = pointer.From(props.RetentionPeriodInDays)
				state.Type = pointer.FromEnum(props.BackupType)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AutonomousDatabaseBackupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if !metadata.ResourceData.HasChange("retention_period_in_days") {
				return nil
			}
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabaseBackups

			id, err := autonomousdatabasebackups.ParseAutonomousDatabaseBackupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			adbId := autonomousdatabases.NewAutonomousDatabaseID(
				id.SubscriptionId,
				id.ResourceGroupName,
				id.AutonomousDatabaseName,
			)
			backupId := autonomousdatabasebackups.NewAutonomousDatabaseBackupID(id.SubscriptionId, id.ResourceGroupName, id.AutonomousDatabaseName, id.AutonomousDatabaseBackupName)

			var model AutonomousDatabaseBackupResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("retrieving %s: %+v", backupId, err)
			}

			_, err = getBackupFromOCI(ctx, client, adbId, backupId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", backupId, err)
			}

			update := autonomousdatabasebackups.AutonomousDatabaseBackupUpdate{
				Properties: &autonomousdatabasebackups.AutonomousDatabaseBackupUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("retention_period_in_days") {
				update.Properties.RetentionPeriodInDays = &model.RetentionPeriodInDays
			}

			if err := client.UpdateThenPoll(ctx, *id, update); err != nil {
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
				return err
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
