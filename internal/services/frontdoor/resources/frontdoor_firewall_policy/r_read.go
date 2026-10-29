// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor_firewall_policy

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-04-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFrontDoorFirewallPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsPolicyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webapplicationfirewallpolicies.ParseFrontDoorWebApplicationFirewallPolicyIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.PoliciesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Front Door Firewall Policy %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.FrontDoorWebApplicationFirewallPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if loc := model.Location; loc != nil {
			d.Set("location", location.Normalize(*loc))
		}
		if properties := model.Properties; properties != nil {
			if policy := properties.PolicySettings; policy != nil {
				if policy.EnabledState != nil {
					d.Set("enabled", *policy.EnabledState == webapplicationfirewallpolicies.PolicyEnabledStateEnabled)
				}
				if policy.Mode != nil {
					d.Set("mode", string(*policy.Mode))
				}
				d.Set("redirect_url", policy.RedirectURL)
				d.Set("custom_block_response_status_code", policy.CustomBlockResponseStatusCode)
				d.Set("custom_block_response_body", policy.CustomBlockResponseBody)
			}

			if err := d.Set("custom_rule", flattenFrontDoorFirewallCustomRules(properties.CustomRules)); err != nil {
				return fmt.Errorf("flattening `custom_rule`: %+v", err)
			}

			if err := d.Set("frontend_endpoint_ids", FlattenFrontendEndpointLinkSlice(properties.FrontendEndpointLinks)); err != nil {
				return fmt.Errorf("flattening `frontend_endpoint_ids`: %+v", err)
			}

			if err := d.Set("managed_rule", flattenFrontDoorFirewallManagedRules(properties.ManagedRules)); err != nil {
				return fmt.Errorf("flattening `managed_rule`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}
