// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set_standby_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/standbypool/2025-03-01/standbyvirtualmachinepools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.StandbyVirtualMachinePoolsClient

			id, err := standbyvirtualmachinepools.ParseStandbyVirtualMachinePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model VirtualMachineScaleSetStandbyPoolModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("attached_virtual_machine_scale_set_id") {
				properties.Properties.AttachedVirtualMachineScaleSetId = pointer.To(model.AttachedVirtualMachineScaleSetId)
			}

			if metadata.ResourceData.HasChange("elasticity_profile") {
				properties.Properties.ElasticityProfile = expandStandbyVirtualMachinePoolElasticityProfileModel(model.ElasticityProfile)
			}

			if metadata.ResourceData.HasChange("virtual_machine_state") {
				properties.Properties.VirtualMachineState = model.VirtualMachineState
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
