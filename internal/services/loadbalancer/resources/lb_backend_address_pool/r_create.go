// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_backend_address_pool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerBackendAddressPoolCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	loadBalancerId, err := loadbalancers.ParseLoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Load Balancer Name and Group: %+v", err)
	}

	name := d.Get("name").(string)
	id := loadbalancers.NewLoadBalancerBackendAddressPoolID(loadBalancerId.SubscriptionId, loadBalancerId.ResourceGroupName, loadBalancerId.LoadBalancerName, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := lbClient.LoadBalancerBackendAddressPoolsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(backendAddressPoolResourceName, id.ID())
			}
		}
	}

	locks.ByName(name, backendAddressPoolResourceName)
	defer locks.UnlockByName(name, backendAddressPoolResourceName)

	locks.ByID(loadBalancerId.ID())
	defer locks.UnlockByID(loadBalancerId.ID())

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: loadBalancerId.SubscriptionId, ResourceGroupName: loadBalancerId.ResourceGroupName, LoadBalancerName: loadBalancerId.LoadBalancerName}
	lb, err := lbClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(lb.HttpResponse) {
			return fmt.Errorf("%s was not found", *loadBalancerId)
		}
		return fmt.Errorf("retrieving %s: %+v", *loadBalancerId, err)
	}

	param := loadbalancers.BackendAddressPool{
		Name: &id.BackendAddressPoolName,
	}

	// Since API version 2020-05-01, there are two ways to CRUD backend address pool - either via the LB endpoint or via the
	// dedicated BAP endpoint. While based on different sku of the LB, users should insist on interacting one of the two endpoints:
	// - Basic sku: interact with LB endpoint for CUD
	// - Standard sku: interact with BAP endpoint for CUD
	// Particularly, the BAP endpoint can be used for R for bot cases.
	// See: https://github.com/Azure/azure-rest-api-specs/issues/11234 for details.
	if lb.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *loadBalancerId)
	}
	sku := lb.Model.Sku
	if sku == nil {
		return fmt.Errorf("nil or empty `sku` for Load Balancer %q for Backend Address Pool %q was not found", *loadBalancerId, id)
	}

	if len(d.Get("tunnel_interface").([]any)) != 0 && *sku.Name != loadbalancers.LoadBalancerSkuNameGateway {
		return fmt.Errorf("only the Gateway (sku) Load Balancer allows IP based Backend Address Pool configuration,"+
			"whilst %q is of sku %s", id, *sku.Name)
	}
	if len(d.Get("tunnel_interface").([]any)) == 0 && *sku.Name == loadbalancers.LoadBalancerSkuNameGateway {
		return fmt.Errorf("`tunnel_interface` is required for %q when sku is set to %s", id, *sku.Name)
	}

	if _, ok := d.GetOk("synchronous_mode"); ok && *sku.Name != loadbalancers.LoadBalancerSkuNameStandard {
		return fmt.Errorf("`synchronous_mode` can set only for Load Balancer with `Standard` SKU")
	}

	if v, ok := d.GetOk("virtual_network_id"); ok {
		param.Properties = &loadbalancers.BackendAddressPoolPropertiesFormat{
			VirtualNetwork: &loadbalancers.SubResource{
				Id: pointer.To(v.(string)),
			},
		}
	}

	if v, ok := d.GetOk("synchronous_mode"); ok {
		if param.Properties == nil {
			param.Properties = &loadbalancers.BackendAddressPoolPropertiesFormat{
				SyncMode: pointer.ToEnum[loadbalancers.SyncMode](v.(string)),
			}
		} else {
			param.Properties.SyncMode = pointer.ToEnum[loadbalancers.SyncMode](v.(string))
		}
	}

	if properties := lb.Model.Properties; properties != nil {
		switch *sku.Name {
		case loadbalancers.LoadBalancerSkuNameBasic:
			if !d.IsNewResource() && d.HasChange("virtual_network_id") {
				return fmt.Errorf("updating the virtual_network_id of Backend Address Pool %q is not allowed for basic (sku) Load Balancer", id)
			}

			// Insert this BAP and update the LB since the dedicated BAP endpoint doesn't work for the Basic sku.
			backendAddressPools := append(*properties.BackendAddressPools, param)
			_, existingPoolIndex, exists := FindLoadBalancerBackEndAddressPoolByName(lb.Model, id.BackendAddressPoolName)
			if exists {
				// this pool is being updated/reapplied remove the old copy from the slice
				backendAddressPools = append(backendAddressPools[:existingPoolIndex], backendAddressPools[existingPoolIndex+1:]...)
			}

			properties.BackendAddressPools = &backendAddressPools

			if d.IsNewResource() {
				if err := lbClient.CreateOrUpdateCallbackThenPoll(ctx, plbId, *lb.Model, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
			} else {
				if err := lbClient.CreateOrUpdateThenPoll(ctx, plbId, *lb.Model); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
			}
		case loadbalancers.LoadBalancerSkuNameStandard:
			if param.Properties == nil {
				param.Properties = &loadbalancers.BackendAddressPoolPropertiesFormat{
					// NOTE: Backend Addresses are managed using `azurerm_lb_backend_pool_address`
				}
			}

			if d.IsNewResource() {
				if err := lbClient.LoadBalancerBackendAddressPoolsCreateOrUpdateCallbackThenPoll(ctx, id, param, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
			} else {
				if err := lbClient.LoadBalancerBackendAddressPoolsCreateOrUpdateThenPoll(ctx, id, param); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
			}
		case loadbalancers.LoadBalancerSkuNameGateway:
			if param.Properties == nil {
				param.Properties = &loadbalancers.BackendAddressPoolPropertiesFormat{}
			}
			param.Properties.TunnelInterfaces = expandGatewayLoadBalancerTunnelInterfaces(d.Get("tunnel_interface").([]any))

			if d.IsNewResource() {
				if err := lbClient.LoadBalancerBackendAddressPoolsCreateOrUpdateCallbackThenPoll(ctx, id, param, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
			} else {
				if err := lbClient.LoadBalancerBackendAddressPoolsCreateOrUpdateThenPoll(ctx, id, param); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
			}
		}

		if d.IsNewResource() {
			d.SetId(id.ID())
			if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
				return err
			}
		}
	}

	return resourceArmLoadBalancerBackendAddressPoolRead(d, meta)
}

func expandGatewayLoadBalancerTunnelInterfaces(input []any) *[]loadbalancers.GatewayLoadBalancerTunnelInterface {
	if len(input) == 0 {
		return nil
	}

	result := make([]loadbalancers.GatewayLoadBalancerTunnelInterface, 0)

	for _, e := range input {
		e := e.(map[string]any)
		result = append(result, loadbalancers.GatewayLoadBalancerTunnelInterface{
			Identifier: pointer.To(int64(e["identifier"].(int))),
			Type:       pointer.ToEnum[loadbalancers.GatewayLoadBalancerTunnelInterfaceType](e["type"].(string)),
			Protocol:   pointer.ToEnum[loadbalancers.GatewayLoadBalancerTunnelProtocol](e["protocol"].(string)),
			Port:       pointer.To(int64(e["port"].(int))),
		})
	}

	return &result
}
