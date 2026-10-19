// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_firewall_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2025-03-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorFirewallPolicyCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorFirewallPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := webapplicationfirewallpolicies.NewFrontDoorWebApplicationFirewallPolicyID(subscriptionId, resourceGroup, name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		result, err := client.PoliciesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(result.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(result.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_firewall_policy", id.ID())
		}
	}

	enabled := webapplicationfirewallpolicies.PolicyEnabledStateDisabled

	if d.Get("enabled").(bool) {
		enabled = webapplicationfirewallpolicies.PolicyEnabledStateEnabled
	}

	requestBodyCheck := webapplicationfirewallpolicies.PolicyRequestBodyCheckDisabled

	if d.Get("request_body_check_enabled").(bool) {
		requestBodyCheck = webapplicationfirewallpolicies.PolicyRequestBodyCheckEnabled
	}

	sku := d.Get("sku_name").(string)
	mode := webapplicationfirewallpolicies.PolicyMode(d.Get("mode").(string))
	redirectUrl := d.Get("redirect_url").(string)
	customBlockResponseStatusCode := d.Get("custom_block_response_status_code").(int)
	customBlockResponseBody := d.Get("custom_block_response_body").(string)
	customRules := d.Get("custom_rule").([]any)

	managedRules, err := expandCdnFrontDoorFirewallManagedRules(d.Get("managed_rule").([]any))
	if err != nil {
		return fmt.Errorf("expanding 'managed_rule': %+v", err)
	}

	logScrubbingRules, err := expandCdnFrontDoorFirewallLogScrubbingPolicy(d.Get("log_scrubbing").([]any))
	if err != nil {
		return fmt.Errorf("expanding 'log_scrubbing': %+v", err)
	}

	if sku != string(webapplicationfirewallpolicies.SkuNamePremiumAzureFrontDoor) && managedRules != nil {
		return fmt.Errorf("the 'managed_rule' field is only supported with the 'Premium_AzureFrontDoor' sku, got %q", sku)
	}

	payload := webapplicationfirewallpolicies.WebApplicationFirewallPolicy{
		Location: pointer.To(location.Normalize("Global")),
		Sku: &webapplicationfirewallpolicies.Sku{
			Name: pointer.ToEnum[webapplicationfirewallpolicies.SkuName](sku),
		},
		Properties: &webapplicationfirewallpolicies.WebApplicationFirewallPolicyProperties{
			PolicySettings: &webapplicationfirewallpolicies.PolicySettings{
				EnabledState:     pointer.To(enabled),
				Mode:             pointer.To(mode),
				RequestBodyCheck: pointer.To(requestBodyCheck),
			},
			CustomRules: expandCdnFrontDoorFirewallCustomRules(customRules),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	// NOTE: CAPTCHA and JS Challenge Expiration policy is enabled by default on Premium SKU's with a default of
	// 30 minutes, if it is not in the config set the default and include it in the policy settings
	// payload block...
	if sku == string(webapplicationfirewallpolicies.SkuNamePremiumAzureFrontDoor) {
		// Set the Default values...
		jsChallengeExpirationInMinutes := 30
		captchaExpirationInMinutes := 30

		if v, ok := d.GetOk("js_challenge_cookie_expiration_in_minutes"); ok {
			jsChallengeExpirationInMinutes = v.(int)
		}

		if v, ok := d.GetOk("captcha_cookie_expiration_in_minutes"); ok {
			captchaExpirationInMinutes = v.(int)
		}

		payload.Properties.PolicySettings.JavascriptChallengeExpirationInMinutes = pointer.To(int64(jsChallengeExpirationInMinutes))
		payload.Properties.PolicySettings.CaptchaExpirationInMinutes = pointer.To(int64(captchaExpirationInMinutes))
	}

	if managedRules != nil {
		payload.Properties.ManagedRules = managedRules
	}

	if redirectUrl != "" {
		payload.Properties.PolicySettings.RedirectURL = pointer.To(redirectUrl)
	}

	if customBlockResponseBody != "" {
		payload.Properties.PolicySettings.CustomBlockResponseBody = pointer.To(customBlockResponseBody)
	}

	if customBlockResponseStatusCode > 0 {
		payload.Properties.PolicySettings.CustomBlockResponseStatusCode = pointer.To(int64(customBlockResponseStatusCode))
	}

	if logScrubbingRules != nil {
		payload.Properties.PolicySettings.LogScrubbing = logScrubbingRules
	}

	if err := client.PoliciesCreateOrUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontDoorFirewallPolicyRead(d, meta)
}
