// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_container_immutability_policy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/immutabilitypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ResourceManager.ImmutabilityPolicies

			id, err := parse.StorageContainerImmutabilityPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			containerId := commonids.NewStorageContainerID(id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.ContainerName)

			resp, err := client.BlobContainersGetImmutabilityPolicy(ctx, containerId, immutabilitypolicies.DefaultBlobContainersGetImmutabilityPolicyOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) || r.isDeleted(resp.Model) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ContainerImmutabilityPolicyModel{
				StorageContainerResourceManagerId: containerId.ID(),
			}

			if resp.Model != nil {
				props := resp.Model.Properties
				if props.AllowProtectedAppendWrites != nil {
					state.ProtectedAppendWritesEnabled = *props.AllowProtectedAppendWrites
				}
				if props.AllowProtectedAppendWritesAll != nil {
					state.ProtectedAppendWritesAllEnabled = *props.AllowProtectedAppendWritesAll
				}
				if props.ImmutabilityPeriodSinceCreationInDays != nil {
					state.ImmutabilityPeriodInDays = *props.ImmutabilityPeriodSinceCreationInDays
				}
				if props.State != nil {
					state.Locked = *props.State == immutabilitypolicies.ImmutabilityPolicyStateLocked
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r Resource) isDeleted(input *immutabilitypolicies.ImmutabilityPolicy) bool {
	if input == nil {
		return true
	}
	if input.Properties.State != nil && strings.EqualFold(string(*input.Properties.State), "Deleted") {
		return true
	}
	return false
}
