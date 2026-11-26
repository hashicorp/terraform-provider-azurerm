// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package lb_rule_test

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LbBackendAddressPoolResource struct{}

func (r LbBackendAddressPoolResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := loadbalancers.ParseLoadBalancerBackendAddressPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	lb, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(lb.HttpResponse) {
			return nil, fmt.Errorf("%s was not found", plbId)
		}
		return nil, fmt.Errorf("retrieving %s: %+v", plbId, err)
	}
	if model := lb.Model; model != nil {
		props := model.Properties
		if props == nil || props.BackendAddressPools == nil || len(*props.BackendAddressPools) == 0 {
			return nil, fmt.Errorf("Backend Pool %q not found in Load Balancer %q (resource group %q)", id.BackendAddressPoolName, id.LoadBalancerName, id.ResourceGroupName)
		}

		found := false
		for _, v := range *props.BackendAddressPools {
			if v.Name != nil && *v.Name == id.BackendAddressPoolName {
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("Backend Pool %q not found in Load Balancer %q (resource group %q)", id.BackendAddressPoolName, id.LoadBalancerName, id.ResourceGroupName)
		}
	}

	return pointer.To(true), nil
}

type LbProbeResource struct{}

func (r LbProbeResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := loadbalancers.ParseProbeID(state.ID)
	if err != nil {
		return nil, err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	lb, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(lb.HttpResponse) {
			return nil, fmt.Errorf("%s was not found", plbId)
		}
		return nil, fmt.Errorf("retrieving %s: %+v", plbId, err)
	}
	props := lb.Model
	if props == nil || props.Properties.Probes == nil || len(*props.Properties.Probes) == 0 {
		return nil, fmt.Errorf("Probe %q not found in Load Balancer %q (resource group %q)", id.ProbeName, id.LoadBalancerName, id.ResourceGroupName)
	}

	found := false
	for _, v := range *props.Properties.Probes {
		if v.Name != nil && *v.Name == id.ProbeName {
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("Probe %q not found in Load Balancer %q (resource group %q)", id.ProbeName, id.LoadBalancerName, id.ResourceGroupName)
	}
	return pointer.To(found), nil
}
