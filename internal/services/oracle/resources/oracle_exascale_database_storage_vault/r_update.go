// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_exascale_database_storage_vault

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/exascaledbstoragevaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ExascaleDbStorageVaults
			id, err := exascaledbstoragevaults.ParseExascaleDbStorageVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ExascaleDatabaseStorageVaultResourceModel
			if err = metadata.Decode(&model); err != nil {
				return err
			}

			if _, err = client.Get(ctx, *id); err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			update := &exascaledbstoragevaults.ExascaleDbStorageVaultTagsUpdate{
				Tags: pointer.To(model.Tags),
			}

			if err = client.UpdateThenPoll(ctx, *id, *update); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}
