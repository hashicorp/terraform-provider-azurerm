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
	"github.com/hashicorp/go-azure-sdk/resource-manager/standbypool/2025-03-01/standbyvirtualmachinepools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.StandbyVirtualMachinePoolsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model VirtualMachineScaleSetStandbyPoolModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := standbyvirtualmachinepools.NewStandbyVirtualMachinePoolID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &standbyvirtualmachinepools.StandbyVirtualMachinePoolResource{
				Location: location.Normalize(model.Location),
				Properties: &standbyvirtualmachinepools.StandbyVirtualMachinePoolResourceProperties{
					AttachedVirtualMachineScaleSetId: pointer.To(model.AttachedVirtualMachineScaleSetId),
					ElasticityProfile:                expandStandbyVirtualMachinePoolElasticityProfileModel(model.ElasticityProfile),
					VirtualMachineState:              model.VirtualMachineState,
				},
				Tags: &model.Tags,
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandStandbyVirtualMachinePoolElasticityProfileModel(inputList []VirtualMachineScaleSetStandbyPoolElasticityProfileModel) *standbyvirtualmachinepools.StandbyVirtualMachinePoolElasticityProfile {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := standbyvirtualmachinepools.StandbyVirtualMachinePoolElasticityProfile{
		MaxReadyCapacity: input.MaxReadyCapacity,
		MinReadyCapacity: pointer.To(input.MinReadyCapacity),
	}

	return &output
}
