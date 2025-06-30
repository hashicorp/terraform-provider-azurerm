// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
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
	ResourceGroupName      string `tfschema:"resource_group_name"`
	AutonomousDataBaseName string `tfschema:"autonomous_database_name"`
	Name                   string `tfschema:"name"`

	// Required
	BackupType            string `tfschema:"backup_type"`
	DisplayName           string `tfschema:"display_name"`
	RetentionPeriodInDays int64  `tfschema:"retention_period_in_days"`

	//computed
	AutonomousDataBaseBackupOcid string  `tfschema:"autonomous_database_backup_ocid"`
	AutonomousDatabaseOcid       string  `tfschema:"autonomous_database_ocid"`
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
}

func (AutonomousDatabaseBackupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		// Required
		"autonomous_database_name": {
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
		"resource_group_name": commonschema.ResourceGroupName(),

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
	return map[string]*pluginsdk.Schema{

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
			log.Printf("[DEBUG] Created resource with ID: %s", id.ID())
			log.Printf("[DEBUG] Resource uses display_name: %s", model.DisplayName)
			log.Printf("[DEBUG] Resource uses autonomous_database_id: %s", model.AutonomousDataBaseName)
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
			resp, err := client.ListByParent(ctx, autonomousdatabasebackups.AutonomousDatabaseId(adbId))
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
					if item.Id != nil {
						itemName = *item.Id
					}

					log.Printf("[DEBUG] Backup %d: Name=%s", i, itemName)

					// UPDATED: Only compare by Name field
					if item.Id != nil && *item.Id == id.ID() {
						log.Printf("[DEBUG] Found matching backup by name: %s", itemName)
						backup = &(*resp.Model)[i]
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

			state := AutonomousDatabaseBackupResourceModel{
				DisplayName:            id.AutonomousDatabaseBackupName,
				ResourceGroupName:      id.ResourceGroupName,
				AutonomousDataBaseName: id.AutonomousDatabaseName,
			}

			fmt.Printf("[DEBUG] Initial State: %+v\n", state)
			if model := backup.Properties; model != nil {
				state.DisplayName = pointer.From(model.DisplayName)
				state.RetentionPeriodInDays = pointer.From(model.RetentionPeriodInDays)
				state.AutonomousDatabaseOcid = pointer.From(model.AutonomousDatabaseOcid)
				state.AutonomousDataBaseBackupOcid = pointer.From(model.Ocid)
				state.BackupType = string(pointer.From(model.BackupType))
				state.DbVersion = pointer.From(model.DbVersion)
				state.BackupSizeInTbs = pointer.From(model.SizeInTbs)
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
			log.Printf("[DEBUG] autonomous_database_name in final state: %s", state.AutonomousDataBaseName)

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
				Properties: &autonomousdatabasebackups.AutonomousDatabaseBackupProperties{},
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
