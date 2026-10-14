// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_service_environment_v3

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/appserviceenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.AppServiceEnvironmentClient
			networksClient := metadata.Client.Network.VirtualNetworks
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AppServiceEnvironmentV3Model
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			subnet, err := commonids.ParseSubnetID(model.SubnetId)
			if err != nil {
				return err
			}

			vnetId := commonids.NewVirtualNetworkID(subnet.SubscriptionId, subnet.ResourceGroupName, subnet.VirtualNetworkName)

			vnet, err := networksClient.Get(ctx, vnetId, virtualnetworks.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving Virtual Network %q (Resource Group %q): %+v", subnet.VirtualNetworkName, subnet.ResourceGroupName, err)
			}
			if vnet.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", subnet)
			}

			vnetLoc := location.NormalizeNilable(vnet.Model.Location)
			if vnetLoc == "" {
				return fmt.Errorf("determining Location from Virtual Network %q (Resource Group %q): `location` was missing", subnet.VirtualNetworkName, subnet.ResourceGroupName)
			}

			id := commonids.NewAppServiceEnvironmentID(subscriptionId, model.ResourceGroup, model.Name)

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

			envelope := expandCreateForAppServiceEnvironmentV3(model, vnetLoc)

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, envelope, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			// Networking config cannot be sent in the initial create and must be updated post-creation.
			aseNetworkConfig := appserviceenvironments.AseV3NetworkingConfiguration{
				Properties: &appserviceenvironments.AseV3NetworkingConfigurationProperties{
					AllowNewPrivateEndpointConnections: pointer.To(model.AllowNewPrivateEndpointConnections),
					RemoteDebugEnabled:                 pointer.To(model.RemoteDebuggingEnabled),
				},
			}

			if _, err := client.UpdateAseNetworkingConfiguration(ctx, id, aseNetworkConfig); err != nil {
				return fmt.Errorf("setting Allow New Private Endpoint Connections on %s: %+v", id, err)
			}

			// Updating Network Config returns quickly, but is actually async on some properties.
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
				Refresh:      checkNetworkConfigUpdate(ctx, client, id, *aseNetworkConfig.Properties),
			}

			if _, err := updateWait.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for Network Update for %s to complete: %+v", id, err)
			}

			return nil
		},
	}
}
