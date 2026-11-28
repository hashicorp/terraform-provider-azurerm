// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_application_firewall_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceWebApplicationFirewallPolicyUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppGateway.WebApplicationFirewallPolicies
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := webapplicationfirewallpolicies.NewApplicationGatewayWebApplicationFirewallPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}
	if resp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties was nil", id)
	}

	model := resp.Model

	if d.HasChange("custom_rules") {
		model.Properties.CustomRules = expandWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(d.Get("custom_rules").([]any))
	}

	if d.HasChange("policy_settings") {
		model.Properties.PolicySettings = expandWebApplicationFirewallPolicyPolicySettings(d.Get("policy_settings").([]any))
	}

	if d.HasChange("managed_rules") {
		model.Properties.ManagedRules = pointer.From(expandWebApplicationFirewallPolicyManagedRulesDefinition(d.Get("managed_rules").([]any)))
	}

	if d.HasChange("tags") {
		model.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if _, err := client.CreateOrUpdate(ctx, id, *model); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	return resourceWebApplicationFirewallPolicyRead(d, meta)
}
