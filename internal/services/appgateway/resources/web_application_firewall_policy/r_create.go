// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_application_firewall_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceWebApplicationFirewallPolicyCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppGateway.WebApplicationFirewallPolicies
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := webapplicationfirewallpolicies.NewApplicationGatewayWebApplicationFirewallPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_web_application_firewall_policy", id.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))
	customRules := d.Get("custom_rules").([]any)
	policySettings := d.Get("policy_settings").([]any)
	managedRules := d.Get("managed_rules").([]any)
	t := d.Get("tags").(map[string]any)

	parameters := webapplicationfirewallpolicies.WebApplicationFirewallPolicy{
		Location: pointer.To(location),
		Properties: &webapplicationfirewallpolicies.WebApplicationFirewallPolicyPropertiesFormat{
			CustomRules:    expandWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(customRules),
			PolicySettings: expandWebApplicationFirewallPolicyPolicySettings(policySettings),
			ManagedRules:   pointer.From(expandWebApplicationFirewallPolicyManagedRulesDefinition(managedRules)),
		},
		Tags: tags.Expand(t),
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceWebApplicationFirewallPolicyRead(d, meta)
}
