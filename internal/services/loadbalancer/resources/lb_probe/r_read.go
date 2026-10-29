// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_probe

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerProbeRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseProbeID(d.Id())
	if err != nil {
		return err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	loadBalancer, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(loadBalancer.HttpResponse) {
			d.SetId("")
			log.Printf("[INFO] Load Balancer %q not found. Removing from state", id.LoadBalancerName)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", plbId, err)
	}

	var config *loadbalancers.Probe
	if model := loadBalancer.Model; model != nil {
		var exists bool
		config, _, exists = FindLoadBalancerProbeByName(model, id.ProbeName)
		if !exists {
			d.SetId("")
			log.Printf("[INFO] Load Balancer Probe %q not found. Removing from state", id.ProbeName)
			return nil
		}
	}

	return resourceArmLoadBalancerProbeFlatten(d, id, config)
}

func loadBalancerProbeResourceImporter(_ context.Context, d *pluginsdk.ResourceData, _ any) ([]*pluginsdk.ResourceData, error) {
	id, err := loadbalancers.ParseProbeID(d.Id())
	if err != nil {
		return nil, err
	}

	lbId := loadbalancers.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName)
	d.Set("loadbalancer_id", lbId.ID())

	return []*pluginsdk.ResourceData{d}, nil
}
