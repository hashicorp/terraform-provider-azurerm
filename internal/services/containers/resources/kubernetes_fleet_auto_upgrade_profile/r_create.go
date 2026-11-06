// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_auto_upgrade_profile

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetAutoUpgradeProfilesClient

			var config KubernetesFleetAutoUpgradeProfileResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			fleetId, err := commonids.ParseKubernetesFleetID(config.KubernetesFleetManagerId)
			if err != nil {
				return err
			}

			id := autoupgradeprofiles.NewAutoUpgradeProfileID(fleetId.SubscriptionId, fleetId.ResourceGroupName, fleetId.FleetName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := autoupgradeprofiles.AutoUpgradeProfile{
				Properties: &autoupgradeprofiles.AutoUpgradeProfileProperties{
					Channel:  autoupgradeprofiles.UpgradeChannel(config.Channel),
					Disabled: pointer.To(!config.Enabled),
				},
			}

			if config.UpdateStrategyId != "" {
				payload.Properties.UpdateStrategyId = pointer.To(config.UpdateStrategyId)
			}

			if config.NodeImageSelectionType != "" {
				payload.Properties.NodeImageSelection = &autoupgradeprofiles.AutoUpgradeNodeImageSelection{
					Type: autoupgradeprofiles.AutoUpgradeNodeImageSelectionType(config.NodeImageSelectionType),
				}
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, autoupgradeprofiles.DefaultCreateOrUpdateOperationOptions(), metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}
