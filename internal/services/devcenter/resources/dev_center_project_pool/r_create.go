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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Pools
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DevCenterProjectPoolResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterProjectId, err := projects.ParseProjectID(model.DevCenterProjectId)
			if err != nil {
				return err
			}

			id := pools.NewPoolID(subscriptionId, devCenterProjectId.ResourceGroupName, devCenterProjectId.ProjectName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := pools.Pool{
				Location: location.Normalize(model.Location),
				Properties: &pools.PoolProperties{
					DevBoxDefinitionName:         pointer.To(model.DevBoxDefinitionName),
					NetworkConnectionName:        pointer.To(model.DevCenterAttachedNetworkName),
					LicenseType:                  pointer.To(pools.LicenseTypeWindowsClient),
					ManagedVirtualNetworkRegions: expandDevCenterProjectManagedVirtualNetworkRegions(model.ManagedVirtualNetworkRegions),
					StopOnDisconnect:             expandDevCenterProjectPoolStopOnDisconnect(model.StopOnDisconnectGracePeriodMinutes),
				},
				Tags: pointer.To(model.Tags),
			}

			parameters.Properties.VirtualNetworkType = pointer.To(pools.VirtualNetworkTypeUnmanaged)
			if len(model.ManagedVirtualNetworkRegions) != 0 {
				parameters.Properties.VirtualNetworkType = pointer.To(pools.VirtualNetworkTypeManaged)
			}

			if model.LocalAdministratorEnabled {
				parameters.Properties.LocalAdministrator = pointer.To(pools.LocalAdminStatusEnabled)
			} else {
				parameters.Properties.LocalAdministrator = pointer.To(pools.LocalAdminStatusDisabled)
			}

			parameters.Properties.SingleSignOnStatus = pointer.To(pools.SingleSignOnStatusDisabled)
			if model.SingleSignOnEnabled {
				parameters.Properties.SingleSignOnStatus = pointer.To(pools.SingleSignOnStatusEnabled)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandDevCenterProjectPoolStopOnDisconnect(input int64) *pools.StopOnDisconnectConfiguration {
	if input == 0 {
		return &pools.StopOnDisconnectConfiguration{
			GracePeriodMinutes: pointer.To(int64(0)),
			Status:             pointer.To(pools.StopOnDisconnectEnableStatusDisabled),
		}
	}

	return &pools.StopOnDisconnectConfiguration{
		GracePeriodMinutes: pointer.To(input),
		Status:             pointer.To(pools.StopOnDisconnectEnableStatusEnabled),
	}
}

func expandDevCenterProjectManagedVirtualNetworkRegions(input []string) *[]string {
	result := make([]string, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		result = append(result, location.Normalize(v))
	}

	return &result
}
