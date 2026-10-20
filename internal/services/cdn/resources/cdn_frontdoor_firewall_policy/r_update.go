// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_firewall_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2025-03-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorFirewallPolicyUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorFirewallPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webapplicationfirewallpolicies.ParseFrontDoorWebApplicationFirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	result, err := client.PoliciesGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := result.Model

	if model == nil {
		return fmt.Errorf("retrieving %s: 'model' was nil", *id)
	}

	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: 'model.Properties' was nil", *id)
	}

	props := *model.Properties

	// lintignore:R019 // deliberate subset: only the fields feeding PolicySettings; custom_rule, managed_rule and tags are applied in separate branches below
	if d.HasChanges(
		"custom_block_response_body",
		"custom_block_response_status_code",
		"enabled",
		"mode",
		"redirect_url",
		"request_body_check_enabled",
		"js_challenge_cookie_expiration_in_minutes",
		"captcha_cookie_expiration_in_minutes",
		"log_scrubbing",
	) {
		enabled := webapplicationfirewallpolicies.PolicyEnabledStateDisabled
		if d.Get("enabled").(bool) {
			enabled = webapplicationfirewallpolicies.PolicyEnabledStateEnabled
		}

		requestBodyCheck := webapplicationfirewallpolicies.PolicyRequestBodyCheckDisabled
		if d.Get("request_body_check_enabled").(bool) {
			requestBodyCheck = webapplicationfirewallpolicies.PolicyRequestBodyCheckEnabled
		}

		props.PolicySettings = &webapplicationfirewallpolicies.PolicySettings{
			EnabledState:     pointer.To(enabled),
			Mode:             pointer.ToEnum[webapplicationfirewallpolicies.PolicyMode](d.Get("mode").(string)),
			RequestBodyCheck: pointer.To(requestBodyCheck),
		}

		// NOTE: 'captcha_cookie_expiration_in_minutes' and 'js_challenge_cookie_expiration_in_minutes'
		// is only valid for 'Premium_AzureFrontDoor' skus...
		if model.Sku != nil && model.Sku.Name != nil && *model.Sku.Name == webapplicationfirewallpolicies.SkuNamePremiumAzureFrontDoor {
			// Set the Default value...
			jsChallengeExpirationInMinutes := 30
			captchaExpirationInMinutes := 30

			if v, ok := d.GetOk("js_challenge_cookie_expiration_in_minutes"); ok {
				jsChallengeExpirationInMinutes = v.(int)
			}

			if v, ok := d.GetOk("captcha_cookie_expiration_in_minutes"); ok {
				captchaExpirationInMinutes = v.(int)
			}

			props.PolicySettings.JavascriptChallengeExpirationInMinutes = pointer.To(int64(jsChallengeExpirationInMinutes))
			props.PolicySettings.CaptchaExpirationInMinutes = pointer.To(int64(captchaExpirationInMinutes))
		}

		if redirectUrl := d.Get("redirect_url").(string); redirectUrl != "" {
			props.PolicySettings.RedirectURL = pointer.To(redirectUrl)
		}

		if body := d.Get("custom_block_response_body").(string); body != "" {
			props.PolicySettings.CustomBlockResponseBody = pointer.To(body)
		}

		if statusCode := int64(d.Get("custom_block_response_status_code").(int)); statusCode > 0 {
			props.PolicySettings.CustomBlockResponseStatusCode = pointer.To(statusCode)
		}
	}

	if d.HasChange("custom_rule") {
		props.CustomRules = expandCdnFrontDoorFirewallCustomRules(d.Get("custom_rule").([]any))
	}

	if d.HasChange("managed_rule") {
		if model.Sku == nil {
			return fmt.Errorf("retrieving %s: 'model.Sku' was nil", *id)
		}

		managedRules, err := expandCdnFrontDoorFirewallManagedRules(d.Get("managed_rule").([]any))
		if err != nil {
			return fmt.Errorf("expanding managed_rule: %+v", err)
		}

		if pointer.From(model.Sku.Name) != webapplicationfirewallpolicies.SkuNamePremiumAzureFrontDoor && managedRules != nil {
			return fmt.Errorf("the 'managed_rule' field is only supported when using the sku 'Premium_AzureFrontDoor', got %q", pointer.From(model.Sku.Name))
		}

		if managedRules != nil {
			props.ManagedRules = managedRules
		}
	}

	if d.HasChange("log_scrubbing") {
		logScrubbingPolicy, err := expandCdnFrontDoorFirewallLogScrubbingPolicy(d.Get("log_scrubbing").([]any))
		if err != nil {
			return err
		}

		if logScrubbingPolicy != nil {
			props.PolicySettings.LogScrubbing = logScrubbingPolicy
		}
	}

	if d.HasChange("tags") {
		model.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	model.Properties = pointer.To(props)

	if err = client.PoliciesCreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorFirewallPolicyRead(d, meta)
}
