// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseLoadBalancerID(d.Id())
	if err != nil {
		return err
	}

	// Casting because of a bug in provider ID generation due to a data issue that needs investigation.
	resp, err := client.Get(ctx, loadbalancers.ProviderLoadBalancerId(*id), loadbalancers.GetOperationOptions{})
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("updating %s: model was nil", *id)
	}

	model := *resp.Model

	if d.HasChange("frontend_ip_configuration") {
		model.Properties.FrontendIPConfigurations = expandAzureRmLoadBalancerFrontendIpConfigurations(d)
	}

	if d.HasChange("tags") {
		model.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err = client.CreateOrUpdateThenPoll(ctx, loadbalancers.ProviderLoadBalancerId(*id), model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceArmLoadBalancerRead(d, meta)
}
