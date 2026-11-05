// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_restore_point

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/restorepointcollections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/restorepoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.RestorePointsClient

			var config VirtualMachineRestorePointResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			collectionId, err := restorepointcollections.ParseRestorePointCollectionID(config.VirtualMachineRestorePointCollectionId)
			if err != nil {
				return err
			}

			id := restorepoints.NewRestorePointID(collectionId.SubscriptionId, collectionId.ResourceGroupName, collectionId.RestorePointCollectionName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id, restorepoints.DefaultGetOperationOptions())
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := restorepoints.RestorePoint{
				Properties: &restorepoints.RestorePointProperties{},
			}

			if config.CrashConsistencyModeEnabled {
				parameters.Properties.ConsistencyMode = pointer.To(restorepoints.ConsistencyModeTypesCrashConsistent)
			}

			if len(config.ExcludedDisks) > 0 {
				excludedDisks := make([]restorepoints.ApiEntityReference, 0)
				for _, diskId := range config.ExcludedDisks {
					excludedDisks = append(excludedDisks, restorepoints.ApiEntityReference{
						Id: pointer.To(diskId),
					})
				}

				parameters.Properties.ExcludeDisks = pointer.To(excludedDisks)
			}

			if err = client.CreateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
