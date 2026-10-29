// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_backend_address_pool_address

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			lbClient := metadata.Client.LoadBalancers.LoadBalancersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model BackendAddressPoolAddressModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			poolId, err := loadbalancers.ParseLoadBalancerBackendAddressPoolID(model.BackendAddressPoolId)
			if err != nil {
				return err
			}

			locks.ByName(poolId.BackendAddressPoolName, backendAddressPoolResourceName)
			defer locks.UnlockByName(poolId.BackendAddressPoolName, backendAddressPoolResourceName)

			// Backend Addresses can not be created for Basic sku, so we have to check
			plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: subscriptionId, ResourceGroupName: poolId.ResourceGroupName, LoadBalancerName: poolId.LoadBalancerName}
			lb, err := metadata.Client.LoadBalancers.LoadBalancersClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", plbId, err)
			}

			isBasicSku := true
			if lb := lb.Model; lb != nil {
				if lb.Sku != nil && pointer.From(lb.Sku.Name) != loadbalancers.LoadBalancerSkuNameBasic {
					isBasicSku = false
				}
				if isBasicSku {
					return errors.New("backend addresses are not supported on Basic SKU Load Balancers")
				}
				if lb.Sku != nil && pointer.From(lb.Sku.Tier) == loadbalancers.LoadBalancerSkuTierGlobal {
					if model.FrontendIPConfiguration == "" {
						return fmt.Errorf("please set a Regional Backend Address Pool Addresses for the Global load balancer")
					}
				}

				id := parse.NewBackendAddressPoolAddressID(subscriptionId, poolId.ResourceGroupName, poolId.LoadBalancerName, poolId.BackendAddressPoolName, model.Name)
				pool, err := lbClient.LoadBalancerBackendAddressPoolsGet(ctx, *poolId)
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", *poolId, err)
				}

				if pool.Model == nil {
					return fmt.Errorf("retrieving %s: `model` was nil", *poolId)
				}

				if pool.Model.Properties == nil {
					return fmt.Errorf("retrieving %s: `properties` was nil", *poolId)
				}

				if lb.Sku != nil && pointer.From(lb.Sku.Tier) == loadbalancers.LoadBalancerSkuTierRegional {
					if pointer.From(pool.Model.Properties.SyncMode) != "Manual" && (model.IPAddress != "" && model.VirtualNetworkId == "" || model.IPAddress == "" && model.VirtualNetworkId != "") {
						return fmt.Errorf("for regional load balancer, `ip_address` and `virtual_network_id` should be specified when sync mode is not `Manual`")
					}
				}

				addresses := make([]loadbalancers.LoadBalancerBackendAddress, 0)
				if pool.Model.Properties.LoadBalancerBackendAddresses != nil {
					addresses = *pool.Model.Properties.LoadBalancerBackendAddresses
				}

				if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
					for _, address := range addresses {
						if address.Name == nil {
							continue
						}

						if *address.Name == model.Name {
							return metadata.ResourceRequiresImport(r.ResourceType(), id)
						}
					}
				}

				if pointer.From(lb.Sku.Tier) == loadbalancers.LoadBalancerSkuTierGlobal {
					addresses = append(addresses, loadbalancers.LoadBalancerBackendAddress{
						Name: pointer.To(model.Name),
						Properties: &loadbalancers.LoadBalancerBackendAddressPropertiesFormat{
							LoadBalancerFrontendIPConfiguration: &loadbalancers.SubResource{
								Id: pointer.To(model.FrontendIPConfiguration),
							},
						},
					})
				} else {
					address := loadbalancers.LoadBalancerBackendAddress{
						Properties: &loadbalancers.LoadBalancerBackendAddressPropertiesFormat{},
						Name:       pointer.To(model.Name),
					}
					if model.IPAddress != "" {
						address.Properties.IPAddress = pointer.To(model.IPAddress)
					}
					if model.VirtualNetworkId != "" {
						address.Properties.VirtualNetwork = &loadbalancers.SubResource{
							Id: pointer.To(model.VirtualNetworkId),
						}
					}
					addresses = append(addresses, address)
				}

				pool.Model.Properties.LoadBalancerBackendAddresses = &addresses

				// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
				if err := lbClient.LoadBalancerBackendAddressPoolsCreateOrUpdateThenPoll(ctx, *poolId, *pool.Model); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}

				metadata.SetID(id)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
