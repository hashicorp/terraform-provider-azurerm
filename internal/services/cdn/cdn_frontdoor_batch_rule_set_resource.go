// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name cdn_frontdoor_batch_rule_set -service-package-name cdn -properties "name" -compare-values "subscription_id:cdn_frontdoor_profile_id,resource_group_name:cdn_frontdoor_profile_id,profile_name:cdn_frontdoor_profile_id"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	legacyrulesets "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.ResourceWithCustomizeDiff = CdnFrontDoorBatchRuleSetResource{}
	_ sdk.ResourceWithUpdate        = CdnFrontDoorBatchRuleSetResource{}
	_ sdk.ResourceWithIdentity      = CdnFrontDoorBatchRuleSetResource{}
)

type CdnFrontDoorBatchRuleSetResource struct{}

func (CdnFrontDoorBatchRuleSetResource) Identity() resourceids.ResourceId {
	return &rules.RuleSetId{}
}

func (CdnFrontDoorBatchRuleSetResource) ResourceType() string {
	return "azurerm_cdn_frontdoor_batch_rule_set"
}

func (CdnFrontDoorBatchRuleSetResource) ModelObject() interface{} {
	return &CdnFrontDoorBatchRuleSetModel{}
}

func (CdnFrontDoorBatchRuleSetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return rules.ValidateRuleSetID
}

func (CdnFrontDoorBatchRuleSetResource) Arguments() map[string]*pluginsdk.Schema {
	return cdnFrontDoorBatchRuleSetArguments()
}

func (CdnFrontDoorBatchRuleSetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r CdnFrontDoorBatchRuleSetResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CdnFrontDoorBatchRuleSetModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding diff: %+v", err)
			}

			if err := validateCdnFrontDoorBatchRules(model.Rules); err != nil {
				return err
			}

			if err := validateCdnFrontDoorBatchRuleDiffQuota(metadata.ResourceDiff); err != nil {
				return err
			}

			rawConfig := metadata.ResourceDiff.GetRawConfig()
			if rawConfig.IsNull() {
				return nil
			}

			rulesValue := rawConfig.AsValueMap()["rules"]
			if rulesValue.IsNull() || rulesValue.LengthInt() == 0 {
				return nil
			}

			for _, ruleValue := range rulesValue.AsValueSlice() {
				if ruleValue.IsNull() || !ruleValue.IsKnown() {
					continue
				}

				conditions := ruleValue.AsValueMap()["conditions"]
				if conditions.IsNull() || !conditions.IsKnown() || conditions.LengthInt() == 0 {
					continue
				}

				conditionEntries := conditions.AsValueSlice()
				if len(conditionEntries) == 0 || conditionEntries[0].IsNull() || !conditionEntries[0].IsKnown() {
					continue
				}

				conditionBlock := conditionEntries[0].AsValueMap()
				if err := validateFrontDoorConditionBlocksRequireMatchValues(conditionBlock, []string{"request_scheme_condition", "is_device_condition"}); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func validateCdnFrontDoorBatchRuleDiffQuota(diff *pluginsdk.ResourceDiff) error {
	if diff == nil {
		return nil
	}

	oldRaw, newRaw := diff.GetChange("rules")
	oldRules, err := batchRuleDiffEntries(oldRaw)
	if err != nil {
		return fmt.Errorf("calculating the effective diff for `rules`: %+v", err)
	}
	newRules, err := batchRuleDiffEntries(newRaw)
	if err != nil {
		return fmt.Errorf("calculating the effective diff for `rules`: %+v", err)
	}

	names := make(map[string]struct{}, len(oldRules)+len(newRules))
	for name := range oldRules {
		names[name] = struct{}{}
	}
	for name := range newRules {
		names[name] = struct{}{}
	}

	changedRules := 0
	cacheOperations := 0
	for name := range names {
		oldRule, oldExists := oldRules[name]
		newRule, newExists := newRules[name]

		if oldExists == newExists && oldRule.signature == newRule.signature {
			continue
		}

		changedRules++
		if oldRule.usesCache || newRule.usesCache {
			cacheOperations++
		}
	}

	effectiveDiff := changedRules + cacheOperations
	if effectiveDiff > 100 {
		return fmt.Errorf("the effective diff for `rules` exceeds the service-side quota: got `%d` changed rules and `%d` cache operations, total `%d`, but the maximum allowed is `100`", changedRules, cacheOperations, effectiveDiff)
	}

	return nil
}

type batchRuleDiffEntry struct {
	signature string
	usesCache bool
}

func batchRuleDiffEntries(input interface{}) (map[string]batchRuleDiffEntry, error) {
	results := make(map[string]batchRuleDiffEntry)
	if input == nil {
		return results, nil
	}

	rulesList, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected a list of `rules`, got %T", input)
	}

	for _, item := range rulesList {
		if item == nil {
			continue
		}

		rule, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected a `rules` block, got %T", item)
		}

		name, ok := rule["name"].(string)
		if !ok || name == "" {
			return nil, errors.New("expected each `rules` block to contain a non-empty `name`")
		}

		normalized, err := normalizeBatchRuleDiffValue(rule)
		if err != nil {
			return nil, fmt.Errorf("normalizing `rules[%s]`: %+v", name, err)
		}

		serialized, err := json.Marshal(normalized)
		if err != nil {
			return nil, fmt.Errorf("serializing `rules[%s]`: %+v", name, err)
		}

		results[name] = batchRuleDiffEntry{
			signature: string(serialized),
			usesCache: batchRuleUsesCache(rule),
		}
	}

	return results, nil
}

func normalizeBatchRuleDiffValue(input interface{}) (interface{}, error) {
	switch value := input.(type) {
	case nil, string, bool, int, int32, int64, float32, float64:
		return value, nil
	case []interface{}:
		results := make([]interface{}, 0, len(value))
		for _, item := range value {
			normalized, err := normalizeBatchRuleDiffValue(item)
			if err != nil {
				return nil, err
			}
			results = append(results, normalized)
		}
		return results, nil
	case map[string]interface{}:
		results := make(map[string]interface{}, len(value))
		for key, item := range value {
			normalized, err := normalizeBatchRuleDiffValue(item)
			if err != nil {
				return nil, err
			}
			results[key] = normalized
		}
		return results, nil
	case *pluginsdk.Set:
		items := value.List()
		normalized := make([]string, 0, len(items))
		for _, item := range items {
			normalizedItem, err := normalizeBatchRuleDiffValue(item)
			if err != nil {
				return nil, err
			}
			serialized, err := json.Marshal(normalizedItem)
			if err != nil {
				return nil, err
			}
			normalized = append(normalized, string(serialized))
		}
		sort.Strings(normalized)
		return normalized, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", input)
	}
}

func batchRuleUsesCache(input map[string]interface{}) bool {
	actionsRaw, ok := input["actions"].([]interface{})
	if !ok || len(actionsRaw) == 0 || actionsRaw[0] == nil {
		return false
	}

	actions, ok := actionsRaw[0].(map[string]interface{})
	if !ok {
		return false
	}

	routeOverrides, ok := actions["route_configuration_override_action"].([]interface{})
	if !ok || len(routeOverrides) == 0 || routeOverrides[0] == nil {
		return false
	}

	routeOverride, ok := routeOverrides[0].(map[string]interface{})
	if !ok {
		return false
	}

	cacheBehavior, ok := routeOverride["cache_behavior"].(string)
	if !ok || cacheBehavior == "" {
		return false
	}

	return cacheBehavior != string(rules.RuleIsCompressionEnabledDisabled)
}

func (r CdnFrontDoorBatchRuleSetResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 4 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			batchModeRuleSetClient := metadata.Client.Cdn.FrontDoorRuleSetsClient_v2025_12_01

			var model CdnFrontDoorBatchRuleSetModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			profileId, err := profiles.ParseProfileID(model.CdnFrontDoorProfileID)
			if err != nil {
				return err
			}

			ruleSetId := rules.NewRuleSetID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, model.Name)
			ruleSetClientId := legacyrulesets.NewRuleSetID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := batchModeRuleSetClient.Get(ctx, ruleSetClientId)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("retrieving %s: %+v", ruleSetClientId, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), ruleSetId)
				}
			}

			payload, err := expandCdnFrontDoorBatchRuleSetPayload(true, model)
			if err != nil {
				return err
			}

			if err := batchModeRuleSetClient.CreateThenPoll(ctx, ruleSetClientId, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", ruleSetClientId, err)
			}

			metadata.SetID(&ruleSetId)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &ruleSetId)
		},
	}
}

func (r CdnFrontDoorBatchRuleSetResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			batchModeRuleSetClient := metadata.Client.Cdn.FrontDoorRuleSetsClient_v2025_12_01

			id, err := rules.ParseRuleSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			ruleSetClientId := legacyrulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)
			resp, err := batchModeRuleSetClient.Get(ctx, ruleSetClientId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.BatchMode == nil || !*resp.Model.Properties.BatchMode {
				return fmt.Errorf("retrieving %s: `batch_mode_enabled` must be `true` on the parent Rule Set", id)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r CdnFrontDoorBatchRuleSetResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 4 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			batchModeRuleSetClient := metadata.Client.Cdn.FrontDoorRuleSetsClient_v2025_12_01

			id, err := rules.ParseRuleSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CdnFrontDoorBatchRuleSetModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			ruleSetClientId := legacyrulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)

			payload, err := expandCdnFrontDoorBatchRuleSetPayload(true, model)
			if err != nil {
				return err
			}

			if err := batchModeRuleSetClient.CreateThenPoll(ctx, ruleSetClientId, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r CdnFrontDoorBatchRuleSetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			batchModeRuleSetClient := metadata.Client.Cdn.FrontDoorRuleSetsClient_v2025_12_01

			id, err := rules.ParseRuleSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			ruleSetClientId := legacyrulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)
			if err := batchModeRuleSetClient.DeleteThenPoll(ctx, ruleSetClientId); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
