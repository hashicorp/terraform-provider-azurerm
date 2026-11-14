// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_exascale_database_storage_vault

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/exascaledbstoragevaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ExascaleDbStorageVaults
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ExascaleDatabaseStorageVaultResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := exascaledbstoragevaults.NewExascaleDbStorageVaultID(subscriptionId,
				model.ResourceGroupName,
				model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			param := exascaledbstoragevaults.ExascaleDbStorageVault{
				Name:     pointer.To(model.Name),
				Location: location.Normalize(model.Location),
				Zones:    pointer.To(model.Zones),
				Tags:     pointer.To(model.Tags),
				Properties: &exascaledbstoragevaults.ExascaleDbStorageVaultProperties{
					DisplayName:                   model.DisplayName,
					AdditionalFlashCacheInPercent: pointer.To(model.AdditionalFlashCachePercentage),
					TimeZone:                      pointer.To(model.TimeZone),
					HighCapacityDatabaseStorageInput: exascaledbstoragevaults.ExascaleDbStorageInputDetails{
						TotalSizeInGbs: model.HighCapacityDatabaseStorage[0].TotalSizeInGb,
					},
				},
			}

			if model.Description != "" {
				param.Properties.Description = pointer.To(model.Description)
			}

			if err := client.CreateCallbackThenPoll(ctx, id, param, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
