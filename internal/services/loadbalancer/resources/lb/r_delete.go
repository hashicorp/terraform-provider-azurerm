// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	loadbalancerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerDelete(d *pluginsdk.ResourceData, meta any) error {
	client := loadbalancerClient.From(meta).LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := loadbalancers.ParseLoadBalancerID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
