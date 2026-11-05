// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor_firewall_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-04-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFrontDoorFirewallPolicyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsPolicyClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := webapplicationfirewallpolicies.NewFrontDoorWebApplicationFirewallPolicyID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.PoliciesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing Front Door Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
				}
			}
			if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
				return tf.ImportAsExistsError("azurerm_frontdoor_firewall_policy", id.ID())
			}
		}
	}

	location := location.Normalize("Global")
	enabled := webapplicationfirewallpolicies.PolicyEnabledStateDisabled
	if d.Get("enabled").(bool) {
		enabled = webapplicationfirewallpolicies.PolicyEnabledStateEnabled
	}
	mode := webapplicationfirewallpolicies.PolicyMode(d.Get("mode").(string))
	redirectUrl := d.Get("redirect_url").(string)
	customBlockResponseStatusCode := d.Get("custom_block_response_status_code").(int)
	customBlockResponseBody := d.Get("custom_block_response_body").(string)
	customRules := d.Get("custom_rule").([]any)
	managedRules := d.Get("managed_rule").([]any)

	t := d.Get("tags").(map[string]any)

	frontdoorWebApplicationFirewallPolicy := webapplicationfirewallpolicies.WebApplicationFirewallPolicy{
		Name:     pointer.To(name),
		Location: pointer.To(location),
		Properties: &webapplicationfirewallpolicies.WebApplicationFirewallPolicyProperties{
			PolicySettings: &webapplicationfirewallpolicies.PolicySettings{
				EnabledState: &enabled,
				Mode:         &mode,
			},
			CustomRules:  expandFrontDoorFirewallCustomRules(customRules),
			ManagedRules: expandFrontDoorFirewallManagedRules(managedRules),
		},
		Tags: tags.Expand(t),
	}

	if redirectUrl != "" {
		frontdoorWebApplicationFirewallPolicy.Properties.PolicySettings.RedirectURL = pointer.To(redirectUrl)
	}
	if customBlockResponseBody != "" {
		frontdoorWebApplicationFirewallPolicy.Properties.PolicySettings.CustomBlockResponseBody = pointer.To(customBlockResponseBody)
	}
	if customBlockResponseStatusCode > 0 {
		frontdoorWebApplicationFirewallPolicy.Properties.PolicySettings.CustomBlockResponseStatusCode = pointer.To(int64(customBlockResponseStatusCode))
	}

	if d.IsNewResource() {
		if err := client.PoliciesCreateOrUpdateCallbackThenPoll(ctx, id, frontdoorWebApplicationFirewallPolicy, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.PoliciesCreateOrUpdateThenPoll(ctx, id, frontdoorWebApplicationFirewallPolicy); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceFrontDoorFirewallPolicyRead(d, meta)
}
