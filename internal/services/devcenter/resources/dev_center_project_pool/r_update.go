// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_project_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Pools

			id, err := pools.ParsePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DevCenterProjectPoolResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := pools.PoolUpdate{
				Properties: &pools.PoolUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("dev_box_definition_name") {
				parameters.Properties.DevBoxDefinitionName = pointer.To(model.DevBoxDefinitionName)
			}

			if metadata.ResourceData.HasChange("local_administrator_enabled") {
				if model.LocalAdministratorEnabled {
					parameters.Properties.LocalAdministrator = pointer.To(pools.LocalAdminStatusEnabled)
				} else {
					parameters.Properties.LocalAdministrator = pointer.To(pools.LocalAdminStatusDisabled)
				}
			}

			if metadata.ResourceData.HasChange("dev_center_attached_network_name") {
				parameters.Properties.NetworkConnectionName = pointer.To(model.DevCenterAttachedNetworkName)
				parameters.Properties.VirtualNetworkType = pointer.To(pools.VirtualNetworkTypeUnmanaged)
			}

			if metadata.ResourceData.HasChange("managed_virtual_network_regions") {
				parameters.Properties.ManagedVirtualNetworkRegions = expandDevCenterProjectManagedVirtualNetworkRegions(model.ManagedVirtualNetworkRegions)

				if len(model.ManagedVirtualNetworkRegions) != 0 {
					parameters.Properties.VirtualNetworkType = pointer.To(pools.VirtualNetworkTypeManaged)
				}
			}

			if metadata.ResourceData.HasChange("single_sign_on_enabled") {
				parameters.Properties.SingleSignOnStatus = pointer.To(pools.SingleSignOnStatusDisabled)
				if model.SingleSignOnEnabled {
					parameters.Properties.SingleSignOnStatus = pointer.To(pools.SingleSignOnStatusEnabled)
				}
			}

			if metadata.ResourceData.HasChange("stop_on_disconnect_grace_period_minutes") {
				parameters.Properties.StopOnDisconnect = expandDevCenterProjectPoolStopOnDisconnect(model.StopOnDisconnectGracePeriodMinutes)
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
