// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_auto_upgrade_profile

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/fleetupdatestrategies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetAutoUpgradeProfilesClient

			id, err := autoupgradeprofiles.ParseAutoUpgradeProfileID(metadata.ResourceData.Id())
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

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *autoupgradeprofiles.AutoUpgradeProfileId, model *autoupgradeprofiles.AutoUpgradeProfile) error {
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}
	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties were nil", id)
	}

	state := KubernetesFleetAutoUpgradeProfileResourceModel{
		Name:                     id.AutoUpgradeProfileName,
		KubernetesFleetManagerId: commonids.NewKubernetesFleetID(id.SubscriptionId, id.ResourceGroupName, id.FleetName).ID(),
		Enabled:                  true,
	}

	props := model.Properties
	state.Channel = string(props.Channel)
	if props.UpdateStrategyId != nil {
		updateStrategyId, err := fleetupdatestrategies.ParseUpdateStrategyIDInsensitively(*props.UpdateStrategyId)
		if err != nil {
			return err
		}
		state.UpdateStrategyId = updateStrategyId.ID()
	}
	if props.Disabled != nil {
		state.Enabled = !*props.Disabled
	}

	if props.NodeImageSelection != nil {
		state.NodeImageSelectionType = string(props.NodeImageSelection.Type)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
