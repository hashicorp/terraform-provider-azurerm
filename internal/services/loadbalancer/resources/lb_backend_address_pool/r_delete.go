// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_backend_address_pool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerBackendAddressPoolDelete(d *pluginsdk.ResourceData, meta any) error {
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseLoadBalancerBackendAddressPoolID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := loadbalancers.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName)
	loadBalancerID := loadBalancerId.ID()
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	locks.ByName(id.BackendAddressPoolName, backendAddressPoolResourceName)
	defer locks.UnlockByName(id.BackendAddressPoolName, backendAddressPoolResourceName)

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	lb, err := lbClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(lb.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", loadBalancerId, err)
	}

	if lb.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", plbId)
	}
	sku := lb.Model.Sku
	if sku == nil {
		return fmt.Errorf("nil or empty `sku` for Load Balancer %q for Backend Address Pool %q was not found", loadBalancerId, id)
	}

	if *sku.Name == loadbalancers.LoadBalancerSkuNameBasic {
		_, index, exists := FindLoadBalancerBackEndAddressPoolByName(lb.Model, id.BackendAddressPoolName)
		if !exists {
			return nil
		}

		if lb.Model.Properties == nil {
			return fmt.Errorf("retrieving %s: properties was nil", *id)
		}

		backEndPools := *lb.Model.Properties.BackendAddressPools
		backEndPools = append(backEndPools[:index], backEndPools[index+1:]...)
		lb.Model.Properties.BackendAddressPools = &backEndPools

		if err := lbClient.CreateOrUpdateThenPoll(ctx, plbId, *lb.Model); err != nil {
			return fmt.Errorf("updating %s: %+v", loadBalancerId, err)
		}
	} else {
		if err := lbClient.LoadBalancerBackendAddressPoolsDeleteThenPoll(ctx, *id); err != nil {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
