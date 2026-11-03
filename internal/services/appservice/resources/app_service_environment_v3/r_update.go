// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_service_environment_v3

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/appserviceenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.AppServiceEnvironmentClient

			id, err := commonids.ParseAppServiceEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state AppServiceEnvironmentV3Model
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := existing.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("cluster_setting") {
				model.Properties.ClusterSettings = expandClusterSettingsModel(state.ClusterSetting)
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(state.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			aseNetworkConfig := appserviceenvironments.AseV3NetworkingConfiguration{
				Properties: &appserviceenvironments.AseV3NetworkingConfigurationProperties{
					AllowNewPrivateEndpointConnections: pointer.To(state.AllowNewPrivateEndpointConnections),
					RemoteDebugEnabled:                 pointer.To(state.RemoteDebuggingEnabled),
				},
			}

			if _, err := client.UpdateAseNetworkingConfiguration(ctx, *id, aseNetworkConfig); err != nil {
				return fmt.Errorf("setting Allow New Private Endpoint Connections on %s: %+v", id, err)
			}

			// Updating Network Config returns quickly, but is actually async on some properties. e.g. `RemoteDebuggingEnabled`
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("the Network Configuration Update request context had no deadline")
			}

			updateWait := &pluginsdk.StateChangeConf{
				Pending:      []string{"Pending"},
				Target:       []string{"Succeeded"},
				PollInterval: 10 * time.Second,
				Delay:        10 * time.Second,
				Timeout:      time.Until(deadline),
				Refresh:      checkNetworkConfigUpdate(ctx, client, *id, *aseNetworkConfig.Properties),
			}

			if _, err := updateWait.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for Network Update for %s to complete: %+v", *id, err)
			}

			return nil
		},
	}
}
