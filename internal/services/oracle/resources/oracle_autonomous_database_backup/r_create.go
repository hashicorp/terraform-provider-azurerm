// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_backup

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabasebackups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existingBackup, err := getBackupFromOCI(ctx, client, dbId, id)
				if err != nil {
					return fmt.Errorf("checking for existing backup: %+v", err)
				}
				if existingBackup != nil {
					return metadata.ResourceRequiresImport(r.ResourceType(), &id)
				}
			}

			param := autonomousdatabasebackups.AutonomousDatabaseBackup{
				Name: pointer.To(model.Name),
				Properties: &autonomousdatabasebackups.AutonomousDatabaseBackupProperties{
					RetentionPeriodInDays: pointer.To(model.RetentionPeriodInDays),
					BackupType:            pointer.ToEnum[autonomousdatabasebackups.AutonomousDatabaseBackupType](model.Type),
				},
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, param, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
