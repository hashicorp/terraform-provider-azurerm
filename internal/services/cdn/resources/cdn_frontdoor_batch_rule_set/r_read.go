// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_batch_rule_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.FrontDoorRuleSetsClient

			id, err := rulesets.ParseRuleSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func validateCdnFrontDoorBatchRuleDiffQuota(diff *pluginsdk.ResourceDiff) error {
	hashFunc := pluginsdk.HashResource(Resource{}.Arguments()["rule"].Elem.(*pluginsdk.Resource))

	oldRaw, newRaw := diff.GetChange("rule")
	oldRules := cdnFrontDoorBatchRuleSetRules(oldRaw, hashFunc)
	newRules := cdnFrontDoorBatchRuleSetRules(newRaw, hashFunc)

	names := make(map[string]struct{}, len(oldRules)+len(newRules))
	for name := range oldRules {
		names[name] = struct{}{}
	}
	for name := range newRules {
		names[name] = struct{}{}
	}

	changedRules, cacheOperations := 0, 0
	for name := range names {
		oldRule, oldExists := oldRules[name]
		newRule, newExists := newRules[name]

		if oldExists == newExists && oldRule.hash == newRule.hash {
			continue
		}

		changedRules++
		if oldRule.usesCache || newRule.usesCache {
			cacheOperations++
		}
	}

	// NOTE: Front Door Standard/Premium allows up to 100 rules per ruleset. For batch
	// rulesets, internal service guidance states that a rule with a cache-enabled
	// `route_configuration_override_action` consumes two effective rule slots during
	// replacement. This means a batch update can fail even when the final desired
	// ruleset is within the documented 100 rule limit, because the service evaluates the effective
	// change set for the PATCH and returns a 400 Bad Request once that replacement crosses
	// that quota. We validate `changedRules + cacheOperations` here so the plan
	// fails with the same service-side constraint before sending the batch ruleset update.
	if effectiveDiff := changedRules + cacheOperations; effectiveDiff > 100 {
		return fmt.Errorf("the number of changed `rule` blocks exceeds the service-side quota: got changed `%d` rules, (`%d` of which include caching, which counts for 2 rule slots), total `%d` rule slots used, but the maximum allowed is `100`", changedRules, cacheOperations, effectiveDiff)
	}

	return nil
}

type cdnFrontDoorBatchRuleSetRule struct {
	hash      int
	usesCache bool
}

func cdnFrontDoorBatchRuleSetRules(input any, hash pluginsdk.SchemaSetFunc) map[string]cdnFrontDoorBatchRuleSetRule {
	results := make(map[string]cdnFrontDoorBatchRuleSetRule)
	if input == nil {
		return results
	}

	rulesList, ok := input.([]any)
	if !ok {
		return results
	}

	for _, rule := range rulesList {
		r, ok := rule.(map[string]any)
		if !ok {
			continue
		}

		results[r["name"].(string)] = cdnFrontDoorBatchRuleSetRule{
			hash:      hash(r),
			usesCache: cdnFrontDoorBatchRuleSetRuleUsesCache(r),
		}
	}

	return results
}

func cdnFrontDoorBatchRuleSetRuleUsesCache(input map[string]any) bool {
	actionsRaw, ok := input["actions"].([]any)
	if !ok || len(actionsRaw) == 0 || actionsRaw[0] == nil {
		return false
	}

	actions, ok := actionsRaw[0].(map[string]any)
	if !ok {
		return false
	}

	routeOverrides, ok := actions["route_configuration_override"].([]any)
	if !ok || len(routeOverrides) == 0 || routeOverrides[0] == nil {
		return false
	}

	routeOverride, ok := routeOverrides[0].(map[string]any)
	if !ok {
		return false
	}

	cachingList, ok := routeOverride["caching"].([]any)
	if !ok {
		return false
	}

	caching, ok := cachingList[0].(map[string]any)
	if !ok {
		return false
	}

	cacheBehaviour, ok := caching["behaviour"].(string)
	if !ok {
		return false
	}

	return cacheBehaviour != helpers.RuleCacheBehaviorDisabled
}
