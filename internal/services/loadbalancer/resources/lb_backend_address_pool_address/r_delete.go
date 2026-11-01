// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_backend_address_pool_address

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			lbClient := metadata.Client.LoadBalancers.LoadBalancersClient
			id, err := parse.BackendAddressPoolAddressID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.BackendAddressPoolName, backendAddressPoolResourceName)
			defer locks.UnlockByName(id.BackendAddressPoolName, backendAddressPoolResourceName)

			poolId := loadbalancers.NewLoadBalancerBackendAddressPoolID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
			pool, err := lbClient.LoadBalancerBackendAddressPoolsGet(ctx, poolId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", poolId, err)
			}
			if pool.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", poolId)
			}

			if pool.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", poolId)
			}

			timeout, _ := ctx.Deadline()
			lbStatus := &pluginsdk.StateChangeConf{
				Pending:                   []string{string(loadbalancers.ProvisioningStateUpdating)},
				Target:                    []string{string(loadbalancers.ProvisioningStateSucceeded)},
				MinTimeout:                5 * time.Second,
				Refresh:                   loadbalacnerProvisioningStatusRefreshFunc(ctx, lbClient, *id),
				ContinuousTargetOccurence: 10,
				Timeout:                   time.Until(timeout),
			}

			if _, err := lbStatus.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for parent resource loadbalancer status to be ready error: %+v", err)
			}

			addresses := make([]loadbalancers.LoadBalancerBackendAddress, 0)
			if pool.Model.Properties.LoadBalancerBackendAddresses != nil {
				addresses = *pool.Model.Properties.LoadBalancerBackendAddresses
			}

			newAddresses := make([]loadbalancers.LoadBalancerBackendAddress, 0)
			for _, address := range addresses {
				if address.Name == nil {
					continue
				}

				if *address.Name != id.AddressName {
					newAddresses = append(newAddresses, address)
				}
			}

			pool.Model.Properties.LoadBalancerBackendAddresses = &newAddresses

			if err = lbClient.LoadBalancerBackendAddressPoolsCreateOrUpdateThenPoll(ctx, poolId, *pool.Model); err != nil {
				return fmt.Errorf("removing %s: %+v", *id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
