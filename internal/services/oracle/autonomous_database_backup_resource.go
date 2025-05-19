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
	"log"
	"strings"
	"time"
)

var _ sdk.Resource = AutonomousDatabaseBackupResource{}

type AutonomousDatabaseBackupResource struct{}

type AutonomousDatabaseBackupResourceModel struct {
	Location          string `tfschema:"location"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	DisplayName       string `tfschema:"display_name"`

	// Required
	AutonomousDataBaseID  string `tfschema:"autonomous_database_id"`
	BackupType            string `tfschema:"backup_type"`
	RetentionPeriodInDays int64  `tfschema:"retention_period_in_days"`

	//computed
	AutonomousDataBaseBackupOcid string `tfschema:"autonomous_database_backup_ocid"`
	AutonomousDatabaseOcid       string `tfschema:"autonomous_database_ocid"`
	DbVersion                    string `tfschema:"database_version"`
	BackupSizeInTbs              int64  `tfschema:"database_backup_size_in_tbs"`
	IsAutomatic                  bool   `tfschema:"is_automatic"`
	IsRestorable                 bool   `tfschema:"is_restorable"`
	LifecycleDetails             string `tfschema:"lifecycle_details"`
	LifecycleState               string `tfschema:"lifecycle_state"`
	LicenseModel                 string `tfschema:"license_model"`
	ProvisioningState            string `tfschema:"provisioning_state"`
	TimeAvailableTil             string `tfschema:"time_available_til"`
	TimeEnded                    string `tfschema:"time_ended"`
	TimeStarted                  string `tfschema:"time_started"`
}

func (AutonomousDatabaseBackupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		// Required
		"autonomous_database_id": {
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
	}
}

func (r AutonomousDatabaseBackupResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{

		"location": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"autonomous_database_backup_ocid": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
		"autonomous_database_ocid": {
			Type:     schema.TypeString,
			Computed: true,
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

			autonomousDatabaseId := model.AutonomousDataBaseID

			resourceGroupName, autonomousDatabaseName, err := extractResourceGroupAndNameFromID(autonomousDatabaseId)
			if err != nil {
				return err
			}

			// Check if the autonomous database exists
			dbId := autonomousdatabases.NewAutonomousDatabaseID(subscriptionId, resourceGroupName, autonomousDatabaseName)

			existing, err := dbClient.Get(ctx, dbId)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", dbId, err)
			}

			id := autonomousdatabasebackups.NewAutonomousDatabaseBackupID(
				subscriptionId,
				resourceGroupName,
				autonomousDatabaseName,
				model.DisplayName,
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
					BackupType:            pointer.To(autonomousdatabasebackups.AutonomousDatabaseBackupType(model.BackupType)),
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			log.Printf("[DEBUG] Created resource with ID: %s", id.ID())
			log.Printf("[DEBUG] Resource uses display_name: %s", model.DisplayName)
			log.Printf("[DEBUG] Resource uses autonomous_database_id: %s", model.AutonomousDataBaseID)
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
				fmt.Printf("[DEBUG] Error parsing ID: %s - %+v\n", metadata.ResourceData.Id(), err)
				return fmt.Errorf("parsing ID: %+v", err)
			}
			fmt.Printf("[DEBUG] Reading ID: %s\n", metadata.ResourceData.Id())
			fmt.Printf("[DEBUG] Parsed ID: %+v\n", *id)

			adbId := autonomousdatabases.NewAutonomousDatabaseID(
				id.SubscriptionId,
				id.ResourceGroupName,
				id.AutonomousDatabaseName,
			)

			log.Printf("[DEBUG] Retrieving backups for Autonomous Database %s", adbId.ID())
			resp, err := client.ListByAutonomousDatabase(ctx, autonomousdatabasebackups.AutonomousDatabaseId(adbId))
			if err != nil {
				log.Printf("[ERROR] Failed to list backups: %+v", err)
				return fmt.Errorf("retrieving Autonomous Database Backups: %+v", err)
			}
			log.Printf("[DEBUG] Looking for backup with name: %s", id.AutonomousDatabaseBackupName)
			var backup *autonomousdatabasebackups.AutonomousDatabaseBackup
			if resp.Model != nil {
				log.Printf("[DEBUG] Found %d backups for database", len(*resp.Model))

				for i := range *resp.Model {
					item := (*resp.Model)[i]

					// Log each backup's details
					itemName := "nil"
					if item.Name != nil {
						itemName = *item.Name
					}

					itemDisplayName := "nil"
					if item.Properties != nil && item.Properties.DisplayName != nil {
						itemDisplayName = *item.Properties.DisplayName
					}

					log.Printf("[DEBUG] Backup %d: Name=%s, DisplayName=%s", i, itemName, itemDisplayName)

					// Try matching against both Name and DisplayName
					if (item.Name != nil && *item.Name == id.AutonomousDatabaseBackupName) ||
						(item.Properties != nil && item.Properties.DisplayName != nil && *item.Properties.DisplayName == id.AutonomousDatabaseBackupName) {
						log.Printf("[DEBUG] Found matching backup: %s", itemName)
						backup = &(*resp.Model)[i] // Use direct array access to avoid reference issues
						break
					}
				}
			} else {
				log.Printf("[DEBUG] No backups returned from API (resp.Model is nil)")
			}
			if backup == nil {
				log.Printf("[DEBUG] Resource Autonomous Database Backup %s not found in any of the backups", id.AutonomousDatabaseBackupName)
				metadata.ResourceData.SetId("")
				return nil
			}
			log.Printf("[DEBUG] Successfully found backup %s", *backup.Name)

			// Construct the autonomous database ID from the parsed ID components
			autonomousDatabaseId := adbId.ID()
			log.Printf("[DEBUG] Setting autonomous_database_id to: %s", autonomousDatabaseId)

			state := AutonomousDatabaseBackupResourceModel{
				DisplayName:          id.AutonomousDatabaseBackupName,
				ResourceGroupName:    id.ResourceGroupName,
				AutonomousDataBaseID: autonomousDatabaseId,
			}

			fmt.Printf("[DEBUG] Initial State: %+v\n", state)
			if model := backup.Properties; model != nil {
				state.DisplayName = pointer.From(model.DisplayName)
				state.RetentionPeriodInDays = pointer.From(model.RetentionPeriodInDays)
				state.AutonomousDatabaseOcid = pointer.From(model.AutonomousDatabaseOcid)
				state.AutonomousDataBaseBackupOcid = pointer.From(model.Ocid)
				state.BackupType = string(pointer.From(model.BackupType))
				state.DbVersion = pointer.From(model.DbVersion)
				state.BackupSizeInTbs = int64(pointer.From(model.DatabaseSizeInTbs))
				state.IsAutomatic = pointer.From(model.IsAutomatic)
				state.IsRestorable = pointer.From(model.IsRestorable)
				state.LifecycleDetails = pointer.From(model.LifecycleDetails)
				state.LifecycleState = string(pointer.From(model.LifecycleState))
				state.ProvisioningState = string(pointer.From(model.ProvisioningState))
				state.TimeAvailableTil = pointer.From(model.TimeAvailableTil)
				state.TimeEnded = pointer.From(model.TimeEnded)
				state.TimeStarted = pointer.From(model.TimeStarted)
			}

			log.Printf("[DEBUG] Final state before encoding: %+v", state)
			log.Printf("[DEBUG] autonomous_database_id in final state: %s", state.AutonomousDataBaseID)

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

func extractResourceGroupAndNameFromID(id string) (resourceGroup string, name string, err error) {
	parts := strings.Split(id, "/")
	if len(parts) != 9 {
		return "", "", fmt.Errorf("invalid Autonomous Database ID format: %s", id)
	}

	resourceGroup = parts[4]
	name = parts[8]
	return resourceGroup, name, nil
}
