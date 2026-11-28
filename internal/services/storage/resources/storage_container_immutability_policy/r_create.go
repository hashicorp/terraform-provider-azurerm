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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ResourceManager.ImmutabilityPolicies

			var model ContainerImmutabilityPolicyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			containerId, err := commonids.ParseStorageContainerID(model.StorageContainerResourceManagerId)
			if err != nil {
				return err
			}

			id := parse.NewStorageContainerImmutabilityPolicyID(containerId.SubscriptionId, containerId.ResourceGroupName, containerId.StorageAccountName, "default", containerId.ContainerName, "default")

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.BlobContainersGetImmutabilityPolicy(ctx, *containerId, immutabilitypolicies.DefaultBlobContainersGetImmutabilityPolicyOperationOptions())
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) && !r.isDeleted(existing.Model) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			input := immutabilitypolicies.ImmutabilityPolicy{
				Properties: immutabilitypolicies.ImmutabilityPolicyProperty{
					AllowProtectedAppendWrites:            pointer.To(model.ProtectedAppendWritesEnabled),
					AllowProtectedAppendWritesAll:         pointer.To(model.ProtectedAppendWritesAllEnabled),
					ImmutabilityPeriodSinceCreationInDays: pointer.To(model.ImmutabilityPeriodInDays),
				},
			}

			resp, err := client.BlobContainersCreateOrUpdateImmutabilityPolicy(ctx, *containerId, input, immutabilitypolicies.DefaultBlobContainersCreateOrUpdateImmutabilityPolicyOperationOptions())
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			// Lock the policy if requested - note that this is a one-way operation that prevents subsequent changes or
			// deletion to the policy, the container it applies to, and the storage account where it resides.
			if model.Locked {
				if resp.Model == nil {
					return fmt.Errorf("preparing to lock %s: model was nil", id)
				}

				options := immutabilitypolicies.BlobContainersLockImmutabilityPolicyOperationOptions{
					IfMatch: resp.Model.Etag,
				}

				if _, err = client.BlobContainersLockImmutabilityPolicy(ctx, *containerId, options); err != nil {
					return fmt.Errorf("locking %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
