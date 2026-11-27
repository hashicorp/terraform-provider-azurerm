// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_container_immutability_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/immutabilitypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ResourceManager.ImmutabilityPolicies

			id, err := parse.StorageContainerImmutabilityPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ContainerImmutabilityPolicyModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			containerId, err := commonids.ParseStorageContainerID(model.StorageContainerResourceManagerId)
			if err != nil {
				return err
			}

			resp, err := client.BlobContainersGetImmutabilityPolicy(ctx, *containerId, immutabilitypolicies.DefaultBlobContainersGetImmutabilityPolicyOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) || r.isDeleted(resp.Model) {
					return nil
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			input := immutabilitypolicies.ImmutabilityPolicy{
				Properties: immutabilitypolicies.ImmutabilityPolicyProperty{
					AllowProtectedAppendWrites:            pointer.To(model.ProtectedAppendWritesEnabled),
					AllowProtectedAppendWritesAll:         pointer.To(model.ProtectedAppendWritesAllEnabled),
					ImmutabilityPeriodSinceCreationInDays: pointer.To(model.ImmutabilityPeriodInDays),
				},
			}

			options := immutabilitypolicies.BlobContainersCreateOrUpdateImmutabilityPolicyOperationOptions{
				IfMatch: resp.Model.Etag,
			}

			updateResp, err := client.BlobContainersCreateOrUpdateImmutabilityPolicy(ctx, *containerId, input, options)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			// Lock the policy if requested - note that this is a one-way operation that prevents subsequent changes or
			// deletion to the policy, the container it applies to, and the storage account where it resides.
			if model.Locked {
				if updateResp.Model == nil {
					return fmt.Errorf("preparing to lock %s: model was nil", id)
				}

				lockOptions := immutabilitypolicies.BlobContainersLockImmutabilityPolicyOperationOptions{
					IfMatch: updateResp.Model.Etag,
				}

				if _, err = client.BlobContainersLockImmutabilityPolicy(ctx, *containerId, lockOptions); err != nil {
					return fmt.Errorf("locking %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
