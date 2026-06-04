// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name cdn_frontdoor_batch_rule_set -service-package-name cdn -properties "name" -compare-values "subscription_id:cdn_frontdoor_profile_id,resource_group_name:cdn_frontdoor_profile_id,profile_name:cdn_frontdoor_profile_id"

import (
	"context"
	"fmt"
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

			rawConfig := metadata.ResourceData.GetRawConfig()
			if rawConfig.IsNull() {
				return nil
			}

			rulesValue := rawConfig.AsValueMap()["rules"]
			if rulesValue.IsNull() || rulesValue.LengthInt() == 0 {
				return nil
			}

			for _, ruleValue := range rulesValue.AsValueSlice() {
				if ruleValue.IsNull() {
					continue
				}

				conditions := ruleValue.AsValueMap()["conditions"]
				if conditions.IsNull() || conditions.LengthInt() == 0 {
					continue
				}

				conditionBlock := conditions.AsValueSlice()[0].AsValueMap()
				if err := validateFrontDoorConditionBlocksRequireMatchValues(conditionBlock, []string{"request_scheme_condition", "is_device_condition"}); err != nil {
					return err
				}
			}

			return nil
		},
	}
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
			ruleSetResourceId := legacyrulesets.NewRuleSetID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, model.Name)

			existing, err := batchModeRuleSetClient.Get(ctx, ruleSetResourceId)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("retrieving %s: %+v", ruleSetResourceId, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), ruleSetId)
			}

			payload, err := expandCdnFrontDoorBatchRuleSetPayload(true, model)
			if err != nil {
				return err
			}

			if err := batchModeRuleSetClient.CreateThenPoll(ctx, ruleSetResourceId, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", ruleSetResourceId, err)
			}

			metadata.SetID(&ruleSetId)
			return nil
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

			ruleSetResourceId := legacyrulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)
			resp, err := batchModeRuleSetClient.Get(ctx, ruleSetResourceId)
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

			ruleSetResourceId := legacyrulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)

			payload, err := expandCdnFrontDoorBatchRuleSetPayload(true, model)
			if err != nil {
				return err
			}

			if err := batchModeRuleSetClient.CreateThenPoll(ctx, ruleSetResourceId, payload); err != nil {
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

			if err := batchModeRuleSetClient.DeleteThenPoll(ctx, legacyrulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
