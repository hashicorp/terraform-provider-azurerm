// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_application_firewall_policy

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appgateway/resources/web_application_firewall_policy/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceWebApplicationFirewallPolicyCreate,
		Read:   resourceWebApplicationFirewallPolicyRead,
		Update: resourceWebApplicationFirewallPolicyUpdate,
		Delete: resourceWebApplicationFirewallPolicyDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&webapplicationfirewallpolicies.ApplicationGatewayWebApplicationFirewallPolicyId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&webapplicationfirewallpolicies.ApplicationGatewayWebApplicationFirewallPolicyId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WebApplicationFirewallPolicyV0ToV1{},
		}),

		Schema: webApplicationFirewallPolicySchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, meta any) error {
			if diff.NewValueKnown("managed_rules.0.managed_rule_set") {
				primaryRuleSetCount := 0
				ruleSetTypesKnown := true
				for i, raw := range diff.Get("managed_rules.0.managed_rule_set").([]any) {
					if !diff.NewValueKnown(fmt.Sprintf("managed_rules.0.managed_rule_set.%d.type", i)) {
						ruleSetTypesKnown = false
						break
					}

					ruleSet := raw.(map[string]any)
					switch ruleSet["type"].(string) {
					case "OWASP", "Microsoft_DefaultRuleSet":
						primaryRuleSetCount++
					}
				}

				if ruleSetTypesKnown && primaryRuleSetCount != 1 {
					return errors.New("`managed_rules.0.managed_rule_set` must contain exactly one primary rule set with `type` of `OWASP` or `Microsoft_DefaultRuleSet`")
				}
			}

			if diff.NewValueKnown("managed_rules.0.exception") {
				for i, raw := range diff.Get("managed_rules.0.exception").([]any) {
					if !diff.NewValueKnown(fmt.Sprintf("managed_rules.0.exception.%d.match_variable", i)) ||
						!diff.NewValueKnown(fmt.Sprintf("managed_rules.0.exception.%d.values", i)) {
						continue
					}

					exception := raw.(map[string]any)
					matchVariable := exception["match_variable"].(string)
					switch matchVariable {
					case string(webapplicationfirewallpolicies.ExceptionEntryMatchVariableRequestHeader), string(webapplicationfirewallpolicies.ExceptionEntryMatchVariableRequestURI):
						if len(exception["values"].([]any)) > 10 {
							return fmt.Errorf("`managed_rules.0.exception.%d.values` cannot contain more than 10 items when `match_variable` is `%s`", i, matchVariable)
						}
					}
				}
			}

			if diff.NewValueKnown("managed_rules.0.exception") && diff.NewValueKnown("managed_rules.0.managed_rule_set") &&
				len(diff.Get("managed_rules.0.exception").([]any)) > 0 {
				for i, raw := range diff.Get("managed_rules.0.managed_rule_set").([]any) {
					if !diff.NewValueKnown(fmt.Sprintf("managed_rules.0.managed_rule_set.%d.type", i)) ||
						!diff.NewValueKnown(fmt.Sprintf("managed_rules.0.managed_rule_set.%d.version", i)) {
						continue
					}

					ruleSet := raw.(map[string]any)
					if ruleSet["type"].(string) == "" || ruleSet["version"].(string) == "" {
						continue
					}
					configuredVersion, err := version.NewVersion(ruleSet["version"].(string))
					if err != nil {
						return fmt.Errorf("parsing `managed_rules.0.managed_rule_set.%d.version`: %+v", i, err)
					}

					switch ruleSet["type"].(string) {
					case "OWASP":
						if configuredVersion.LessThan(version.Must(version.NewVersion("3.2"))) {
							return fmt.Errorf("`managed_rules.0.managed_rule_set.%d.version` must be `3.2` or later for `OWASP` when `managed_rules.0.exception` is configured", i)
						}
					case "Microsoft_DefaultRuleSet":
						if configuredVersion.LessThan(version.Must(version.NewVersion("2.1"))) {
							return fmt.Errorf("`managed_rules.0.managed_rule_set.%d.version` must be `2.1` or later for `Microsoft_DefaultRuleSet` when `managed_rules.0.exception` is configured", i)
						}
					}
				}
			}

			return nil
		}),
	}
}
