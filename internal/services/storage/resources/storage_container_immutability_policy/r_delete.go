// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_container_immutability_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/immutabilitypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := storageClient.From(metadata.Client).ResourceManager.ImmutabilityPolicies

			id, err := parse.StorageContainerImmutabilityPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			containerId := commonids.NewStorageContainerID(id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.ContainerName)

			resp, err := client.BlobContainersGetImmutabilityPolicy(ctx, containerId, immutabilitypolicies.DefaultBlobContainersGetImmutabilityPolicyOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) || r.isDeleted(resp.Model) {
					return nil
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			options := immutabilitypolicies.BlobContainersDeleteImmutabilityPolicyOperationOptions{
				IfMatch: resp.Model.Etag,
			}

			if _, err := client.BlobContainersDeleteImmutabilityPolicy(ctx, containerId, options); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
