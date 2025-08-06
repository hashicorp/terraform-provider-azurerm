// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabasebackups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabases"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"log"
	"strings"
	"time"
)

var _ sdk.Resource = AutonomousDatabaseBackupResource{}

type AutonomousDatabaseBackupResource struct{}

type AutonomousDatabaseBackupResourceModel struct {
	AutonomousDatabaseId string `tfschema:"autonomous_database_id"`
	Name                 string `tfschema:"name"`

	// Required
	BackupType            string `tfschema:"backup_type"`
	DisplayName           string `tfschema:"display_name"`
	RetentionPeriodInDays int64  `tfschema:"retention_period_in_days"`
}

func (AutonomousDatabaseBackupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		// Required
		"autonomous_database_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"retention_period_in_days": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(90, 3650),
		},

		// Optional
		"backup_type": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  string(autonomousdatabasebackups.AutonomousDatabaseBackupTypeFull),
			ValidateFunc: validation.StringInSlice([]string{
				string(autonomousdatabasebackups.AutonomousDatabaseBackupTypeFull),
				string(autonomousdatabasebackups.AutonomousDatabaseBackupTypeIncremental),
				string(autonomousdatabasebackups.AutonomousDatabaseBackupTypeLongTerm),
			}, false),
		},
		"display_name": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
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
				return fmt.Errorf("decoding id: %+v", err)
			}

			dbId := autonomousdatabasebackups.NewAutonomousDatabaseID(subscriptionId, parsedAdbsId.ResourceGroupName, parsedAdbsId.AutonomousDatabaseName)

			existing, err := dbClient.Get(ctx, autonomousdatabases.AutonomousDatabaseId(dbId))
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", dbId, err)
			}

			id := autonomousdatabasebackups.NewAutonomousDatabaseBackupID(
				subscriptionId,
				parsedAdbsId.ResourceGroupName,
				parsedAdbsId.AutonomousDatabaseName,
				model.Name,
			)

			existingBackup, err := findBackupByName(ctx, client, autonomousdatabases.AutonomousDatabaseId(dbId), model.Name)
			if err != nil {
				return fmt.Errorf("checking for existing backup: %+v", err)
			}
			if existingBackup != nil {
				return metadata.ResourceRequiresImport(r.ResourceType(), &id)
			}

			param := autonomousdatabasebackups.AutonomousDatabaseBackup{
				Name: pointer.To(model.Name),
				Properties: &autonomousdatabasebackups.AutonomousDatabaseBackupProperties{
					DisplayName:           pointer.To(model.DisplayName),
					RetentionPeriodInDays: pointer.To(model.RetentionPeriodInDays),
					BackupType:            pointer.To(autonomousdatabasebackups.AutonomousDatabaseBackupType(model.BackupType)),
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

			adbId := autonomousdatabases.NewAutonomousDatabaseID(
				id.SubscriptionId,
				id.ResourceGroupName,
				id.AutonomousDatabaseName,
			)

			// Use shared method to find the backup
			backup, err := findBackupByName(ctx, client, adbId, id.AutonomousDatabaseBackupName)
			if err != nil {
				return fmt.Errorf("retrieving backup: %+v", err)
			}

			// Build state from found backup
			state := AutonomousDatabaseBackupResourceModel{
				Name:                 id.AutonomousDatabaseBackupName,
				AutonomousDatabaseId: adbId.ID(),
			}

			if props := backup.Properties; props != nil {
				state.DisplayName = pointer.From(props.DisplayName)
				state.RetentionPeriodInDays = pointer.From(props.RetentionPeriodInDays)
				state.BackupType = string(pointer.From(props.BackupType))
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

func findBackupByName(ctx context.Context, client *autonomousdatabasebackups.AutonomousDatabaseBackupsClient, adbId autonomousdatabases.AutonomousDatabaseId, backupName string) (*autonomousdatabasebackups.AutonomousDatabaseBackup, error) {
	log.Printf("[DEBUG] Looking for backup '%s' in database %s", backupName, adbId.ID())

	resp, err := client.ListByParent(ctx, autonomousdatabasebackups.AutonomousDatabaseId(adbId))
	if err != nil {
		return nil, fmt.Errorf("listing backups for %s: %+v", adbId.ID(), err)
	}

	if resp.Model == nil {
		log.Printf("[DEBUG] No backups found for database %s", adbId.ID())
		return nil, nil
	}

	for i := range *resp.Model {
		backup := (*resp.Model)[i]
		// Match by name
		if backup.Name != nil && *backup.Name == backupName {
			log.Printf("[DEBUG] Found matching backup: %s", *backup.Name)
			return &(*resp.Model)[i], nil
		}
	}
	return nil, nil
}

func (r AutonomousDatabaseBackupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabasebackups.ValidateAutonomousDatabaseBackupID
}

func extractResourceGroupAndNameFromID(id string) (resourceGroup string, name string, err error) {
	parts := strings.Split(id, "/")
	if len(parts) != 9 {
		return "", "", fmt.Errorf("invalid Autonomous Database ID format: %s", id)
	}

	resourceGroup = parts[4]
	name = parts[8]
	return resourceGroup, name, nil
}
