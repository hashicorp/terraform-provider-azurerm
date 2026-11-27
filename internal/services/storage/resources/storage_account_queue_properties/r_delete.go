// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_queue_properties

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	storageClient2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
)

func (s Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := storageClient2.From(metadata.Client)

			id, err := commonids.ParseStorageAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			account, err := storageClient.GetAccount(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if account == nil {
				return fmt.Errorf("unable to locate %s", *id)
			}

			client, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Queues Client for %s: %v", *id, err)
			}

			if err = client.UpdateServiceProperties(ctx, DefaultValueForAccountQueueProperties()); err != nil {
				return fmt.Errorf("updating Queue Properties for %s: %+v", *id, err)
			}

			return nil
		},
	}
}
