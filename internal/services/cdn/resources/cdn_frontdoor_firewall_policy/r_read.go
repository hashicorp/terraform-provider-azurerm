// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_firewall_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2025-03-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorFirewallPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorFirewallPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webapplicationfirewallpolicies.ParseFrontDoorWebApplicationFirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.PoliciesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.FrontDoorWebApplicationFirewallPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if sku := model.Sku; sku != nil {
			d.Set("sku_name", pointer.FromEnum(sku.Name))
		}

		if props := model.Properties; props != nil {
			if err := d.Set("custom_rule", flattenCdnFrontDoorFirewallCustomRules(props.CustomRules)); err != nil {
				return fmt.Errorf("flattening 'custom_rule': %+v", err)
			}

			if err := d.Set("frontend_endpoint_ids", helpers.FlattenFrontendEndpointLinkSlice(props.FrontendEndpointLinks)); err != nil {
				return fmt.Errorf("flattening 'frontend_endpoint_ids': %+v", err)
			}

			if err := d.Set("managed_rule", flattenCdnFrontDoorFirewallManagedRules(props.ManagedRules)); err != nil {
				return fmt.Errorf("flattening 'managed_rule': %+v", err)
			}

			if policy := props.PolicySettings; policy != nil {
				d.Set("enabled", pointer.From(policy.EnabledState) == webapplicationfirewallpolicies.PolicyEnabledStateEnabled)
				d.Set("mode", pointer.FromEnum(policy.Mode))
				d.Set("request_body_check_enabled", pointer.From(policy.RequestBodyCheck) == webapplicationfirewallpolicies.PolicyRequestBodyCheckEnabled)
				d.Set("redirect_url", policy.RedirectURL)
				d.Set("custom_block_response_status_code", int(pointer.From(policy.CustomBlockResponseStatusCode)))
				d.Set("custom_block_response_body", policy.CustomBlockResponseBody)

				// NOTE: `js_challenge_cookie_expiration_in_minutes` and
				// `captcha_cookie_expiration_in_minutes` is only returned
				// for Premium_AzureFrontDoor skus, else they will be 'nil'...
				if policy.JavascriptChallengeExpirationInMinutes != nil {
					d.Set("js_challenge_cookie_expiration_in_minutes", int(pointer.From(policy.JavascriptChallengeExpirationInMinutes)))
				}

				if policy.CaptchaExpirationInMinutes != nil {
					d.Set("captcha_cookie_expiration_in_minutes", int(pointer.From(policy.CaptchaExpirationInMinutes)))
				}

				if err := d.Set("log_scrubbing", flattenCdnFrontDoorFirewallLogScrubbingPolicy(policy.LogScrubbing)); err != nil {
					return fmt.Errorf("flattening 'log_scrubbing': %+v", err)
				}
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
