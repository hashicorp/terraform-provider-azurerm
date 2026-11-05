// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_firewall_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2025-03-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorFirewallPolicyCreate,
		Read:   resourceCdnFrontDoorFirewallPolicyRead,
		Update: resourceCdnFrontDoorFirewallPolicyUpdate,
		Delete: resourceCdnFrontDoorFirewallPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := webapplicationfirewallpolicies.ParseFrontDoorWebApplicationFirewallPolicyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(4 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(4 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Schema: cdnFrontdoorFirewallPolicySchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				currentSku := diff.Get("sku_name").(string)
				standardSku := string(webapplicationfirewallpolicies.SkuNameStandardAzureFrontDoor)
				oldSku, _ := diff.GetChange("sku_name")

				if currentSku == standardSku {
					premiumSku := string(webapplicationfirewallpolicies.SkuNamePremiumAzureFrontDoor)
					managedRules := diff.Get("managed_rule").([]any)
					customRules := expandCdnFrontDoorFirewallCustomRules(diff.Get("custom_rule").([]any))

					// Verify that they are not downgrading the service from Premium SKU -> Standard SKU...
					if oldSku != "" {
						if oldSku.(string) == premiumSku {
							return fmt.Errorf("downgrading from the %q sku to the %q sku is not supported, got %q", premiumSku, standardSku, currentSku)
						}
					}

					// Verify that the Standard SKU is not setting the JSChallenge or Captcha policy...
					if v := diff.Get("js_challenge_cookie_expiration_in_minutes").(int); v > 0 {
						return fmt.Errorf("'js_challenge_cookie_expiration_in_minutes' field is only supported with the %q sku, got %q", premiumSku, currentSku)
					}

					if v := diff.Get("captcha_cookie_expiration_in_minutes").(int); v > 0 {
						return fmt.Errorf("'captcha_cookie_expiration_in_minutes' field is only supported with the %q sku, got %q", premiumSku, currentSku)
					}

					// Verify that the Standard SKU is not using the JSChallenge or CAPTCHA Action type for custom rules...
					if customRules != nil && customRules.Rules != nil {
						for _, v := range *customRules.Rules {
							switch v.Action {
							case webapplicationfirewallpolicies.ActionTypeJSChallenge:
								return fmt.Errorf("'custom_rule' blocks with the 'action' type of 'JSChallenge' are only supported for the %q sku, got action: %q (custom_rule.name: %q, sku_name: %q)", premiumSku, webapplicationfirewallpolicies.ActionTypeJSChallenge, *v.Name, currentSku)
							case webapplicationfirewallpolicies.ActionTypeCAPTCHA:
								return fmt.Errorf("'custom_rule' blocks with the 'action' type of 'CAPTCHA' are only supported for the %q sku, got action: %q (custom_rule.name: %q, sku_name: %q)", premiumSku, webapplicationfirewallpolicies.ActionTypeCAPTCHA, *v.Name, currentSku)
							}
						}
					}

					// Verify that the Standard SKU is not using managed rules...
					if len(managedRules) > 0 {
						return fmt.Errorf("'managed_rule' code block is only supported with the %q sku, got %q", premiumSku, currentSku)
					}
				}

				return nil
			}),

			// Verify that the scrubbing_rule's are valid...
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				if v, ok := diff.GetOk("log_scrubbing"); ok {
					if _, err := expandCdnFrontDoorFirewallLogScrubbingPolicy(v.([]any)); err != nil {
						return err
					}
				}
				return nil
			}),

			// Handle default value reset when field is removed from the configuration
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				rawConfig := diff.GetRawConfig()

				if diff.Get("sku_name").(string) == string(webapplicationfirewallpolicies.SkuNamePremiumAzureFrontDoor) {
					// Force the value to default when removed from config
					if rawConfig.IsNull() || rawConfig.GetAttr("js_challenge_cookie_expiration_in_minutes").IsNull() {
						if diff.Get("js_challenge_cookie_expiration_in_minutes").(int) != 30 {
							if err := diff.SetNew("js_challenge_cookie_expiration_in_minutes", 30); err != nil {
								return fmt.Errorf("setting default for `js_challenge_cookie_expiration_in_minutes`: %+v", err)
							}
						}
					}

					// Force the value to default when removed from config
					if rawConfig.IsNull() || rawConfig.GetAttr("captcha_cookie_expiration_in_minutes").IsNull() {
						if diff.Get("captcha_cookie_expiration_in_minutes").(int) != 30 {
							if err := diff.SetNew("captcha_cookie_expiration_in_minutes", 30); err != nil {
								return fmt.Errorf("setting default for `captcha_cookie_expiration_in_minutes`: %+v", err)
							}
						}
					}
				}

				return nil
			}),
		),
	}
}
