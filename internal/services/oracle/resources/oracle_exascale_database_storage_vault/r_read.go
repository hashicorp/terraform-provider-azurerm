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

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := exascaledbstoragevaults.ParseExascaleDbStorageVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Oracle.OracleClient.ExascaleDbStorageVaults
			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ExascaleDatabaseStorageVaultResourceModel{
				Name:              id.ExascaleDbStorageVaultName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := result.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.Zones = pointer.From(model.Zones)
				if props := model.Properties; props != nil {
					state.DisplayName = props.DisplayName
					state.AdditionalFlashCachePercentage = pointer.From(props.AdditionalFlashCacheInPercent)
					state.Description = pointer.From(props.Description)
					state.HighCapacityDatabaseStorage = FlattenHighCapacityDatabaseStorage(props.HighCapacityDatabaseStorage)
					state.TimeZone = pointer.From(props.TimeZone)
				}
			}
			return metadata.Encode(&state)
		},
	}
}
