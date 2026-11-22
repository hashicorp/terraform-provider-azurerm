// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_project_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/pools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Pools

			id, err := pools.ParsePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := DevCenterProjectPoolResourceModel{
				Name:               id.PoolName,
				DevCenterProjectId: projects.NewProjectID(id.SubscriptionId, id.ResourceGroupName, id.ProjectName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.DevBoxDefinitionName = pointer.From(props.DevBoxDefinitionName)
					state.LocalAdministratorEnabled = pointer.From(props.LocalAdministrator) == pools.LocalAdminStatusEnabled
					state.DevCenterAttachedNetworkName = pointer.From(props.NetworkConnectionName)
					state.ManagedVirtualNetworkRegions = flattenDevCenterProjectManagedVirtualNetworkRegions(props.ManagedVirtualNetworkRegions)
					state.SingleSignOnEnabled = pointer.From(props.SingleSignOnStatus) == pools.SingleSignOnStatusEnabled
					state.StopOnDisconnectGracePeriodMinutes = flattenDevCenterProjectPoolStopOnDisconnect(props.StopOnDisconnect)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenDevCenterProjectPoolStopOnDisconnect(input *pools.StopOnDisconnectConfiguration) int64 {
	var gracePeriodMinutes int64
	if input == nil || pointer.From(input.Status) == pools.StopOnDisconnectEnableStatusDisabled {
		return gracePeriodMinutes
	}

	return pointer.From(input.GracePeriodMinutes)
}

func flattenDevCenterProjectManagedVirtualNetworkRegions(input *[]string) []string {
	result := make([]string, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, location.Normalize(v))
	}

	return result
}
