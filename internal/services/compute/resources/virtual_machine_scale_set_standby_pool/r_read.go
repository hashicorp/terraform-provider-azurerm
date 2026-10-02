// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set_standby_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/standbypool/2025-03-01/standbyvirtualmachinepools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.StandbyVirtualMachinePoolsClient

			id, err := standbyvirtualmachinepools.ParseStandbyVirtualMachinePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := VirtualMachineScaleSetStandbyPoolModel{
				Name:              id.StandbyVirtualMachinePoolName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if properties := model.Properties; properties != nil {
					parsedAttachedVirtualMachineScaleSetId, err := virtualmachinescalesets.ParseVirtualMachineScaleSetIDInsensitively(pointer.From(properties.AttachedVirtualMachineScaleSetId))
					if err != nil {
						return fmt.Errorf("parsing `attached_virtual_machine_scale_set_id` for %s: %+v", *id, err)
					}

					state.AttachedVirtualMachineScaleSetId = parsedAttachedVirtualMachineScaleSetId.ID()
					state.ElasticityProfile = flattenStandbyVirtualMachinePoolElasticityProfileModel(properties.ElasticityProfile)
					state.VirtualMachineState = properties.VirtualMachineState
				}

				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenStandbyVirtualMachinePoolElasticityProfileModel(input *standbyvirtualmachinepools.StandbyVirtualMachinePoolElasticityProfile) []VirtualMachineScaleSetStandbyPoolElasticityProfileModel {
	outputList := make([]VirtualMachineScaleSetStandbyPoolElasticityProfileModel, 0)
	if input == nil {
		return outputList
	}

	output := VirtualMachineScaleSetStandbyPoolElasticityProfileModel{
		MaxReadyCapacity: input.MaxReadyCapacity,
		MinReadyCapacity: pointer.From(input.MinReadyCapacity),
	}

	return append(outputList, output)
}
