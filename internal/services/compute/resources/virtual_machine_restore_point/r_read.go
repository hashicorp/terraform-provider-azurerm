// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_restore_point

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/restorepointcollections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/restorepoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.RestorePointsClient

			schema := VirtualMachineRestorePointResourceModel{}

			id, err := restorepoints.ParseRestorePointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, restorepoints.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.Name = id.RestorePointName
				schema.VirtualMachineRestorePointCollectionId = restorepointcollections.NewRestorePointCollectionID(id.SubscriptionId, id.ResourceGroupName, id.RestorePointCollectionName).ID()

				if props := model.Properties; props != nil {
					schema.CrashConsistencyModeEnabled = strings.EqualFold(pointer.FromEnum(props.ConsistencyMode), string(restorepoints.ConsistencyModeTypesCrashConsistent))

					excludedDisksConfig := make([]string, 0)
					if excludedDisks := props.ExcludeDisks; excludedDisks != nil {
						for _, excludedDisk := range *excludedDisks {
							excludedDisksConfig = append(excludedDisksConfig, pointer.From(excludedDisk.Id))
						}
					}
					schema.ExcludedDisks = excludedDisksConfig
				}
			}

			return metadata.Encode(&schema)
		},
	}
}
