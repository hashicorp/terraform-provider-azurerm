// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_backend_address_pool_address

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			lbClient := metadata.Client.LoadBalancers.LoadBalancersClient
			id, err := parse.BackendAddressPoolAddressID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.BackendAddressPoolName, backendAddressPoolResourceName)
			defer locks.UnlockByName(id.BackendAddressPoolName, backendAddressPoolResourceName)

			var model BackendAddressPoolAddressModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			poolId := loadbalancers.NewLoadBalancerBackendAddressPoolID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
			pool, err := lbClient.LoadBalancerBackendAddressPoolsGet(ctx, poolId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if pool.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if pool.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroup, LoadBalancerName: id.LoadBalancerName}
			lb, err := lbClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", plbId, err)
			}

			addresses := make([]loadbalancers.LoadBalancerBackendAddress, 0)
			if backendAddress := pool.Model.Properties.LoadBalancerBackendAddresses; backendAddress != nil {
				addresses = *backendAddress
			}
			index := -1
			for i, address := range addresses {
				if address.Name == nil {
					continue
				}

				if *address.Name == id.AddressName {
					index = i
					break
				}
			}
			if index == -1 {
				return fmt.Errorf("%s was not found", *id)
			}

			if lb.Model != nil && pointer.From(lb.Model.Sku.Tier) == loadbalancers.LoadBalancerSkuTierGlobal {
				addresses[index] = loadbalancers.LoadBalancerBackendAddress{
					Name: pointer.To(model.Name),
					Properties: &loadbalancers.LoadBalancerBackendAddressPropertiesFormat{
						LoadBalancerFrontendIPConfiguration: &loadbalancers.SubResource{
							Id: pointer.To(model.FrontendIPConfiguration),
						},
					},
				}
			} else {
				addresses[index] = loadbalancers.LoadBalancerBackendAddress{
					Properties: &loadbalancers.LoadBalancerBackendAddressPropertiesFormat{
						IPAddress: pointer.To(model.IPAddress),
						VirtualNetwork: &loadbalancers.SubResource{
							Id: pointer.To(model.VirtualNetworkId),
						},
					},
					Name: pointer.To(id.AddressName),
				}
			}

			pool.Model.Properties.LoadBalancerBackendAddresses = &addresses

			timeout, _ := ctx.Deadline()
			lbStatus := &pluginsdk.StateChangeConf{
				Pending:                   []string{string(loadbalancers.ProvisioningStateUpdating)},
				Target:                    []string{string(loadbalancers.ProvisioningStateSucceeded)},
				MinTimeout:                5 * time.Minute,
				PollInterval:              10 * time.Second,
				Refresh:                   loadbalacnerProvisioningStatusRefreshFunc(ctx, lbClient, *id),
				ContinuousTargetOccurence: 10,
				Timeout:                   time.Until(timeout),
			}

			if _, err := lbStatus.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for parent resource loadbalancer status to be ready error: %+v", err)
			}

			if err = lbClient.LoadBalancerBackendAddressPoolsCreateOrUpdateThenPoll(ctx, poolId, *pool.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
