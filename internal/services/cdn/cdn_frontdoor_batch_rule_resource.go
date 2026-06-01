// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name cdn_frontdoor_batch_rule -service-package-name cdn -compare-values "subscription_id:cdn_frontdoor_rule_set_id,resource_group_name:cdn_frontdoor_rule_set_id,profile_name:cdn_frontdoor_rule_set_id,rule_set_name:cdn_frontdoor_rule_set_id"

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	legacyrulesets "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.ResourceWithCustomizeDiff = CdnFrontDoorBatchRuleResource{}
	_ sdk.ResourceWithUpdate        = CdnFrontDoorBatchRuleResource{}
	_ sdk.ResourceWithIdentity      = CdnFrontDoorBatchRuleResource{}
)

type CdnFrontDoorBatchRuleResource struct{}

func (CdnFrontDoorBatchRuleResource) Identity() resourceids.ResourceId {
	return &rules.RuleSetId{}
}

func (CdnFrontDoorBatchRuleResource) ResourceType() string {
	return "azurerm_cdn_frontdoor_batch_rule"
}

func (CdnFrontDoorBatchRuleResource) ModelObject() interface{} {
	return &CdnFrontDoorBatchRuleModel{}
}

func (CdnFrontDoorBatchRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return rules.ValidateRuleSetID
}

func (CdnFrontDoorBatchRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return cdnFrontDoorBatchRuleArguments()
}

func (CdnFrontDoorBatchRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cdn_frontdoor_rule_set_name": batchRuleComputedSchema["cdn_frontdoor_rule_set_name"],
		"id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r CdnFrontDoorBatchRuleResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CdnFrontDoorBatchRuleModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding diff: %+v", err)
			}

			if err := validateCdnFrontDoorBatchRules(model.Rules); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r CdnFrontDoorBatchRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 4 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.FrontDoorRuleSetsClient_v2025_12_01

			var model CdnFrontDoorBatchRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			ruleSetId, err := rules.ParseRuleSetID(model.CdnFrontDoorRuleSetID)
			if err != nil {
				return err
			}

			legacyRuleSetId := legacyrulesets.NewRuleSetID(ruleSetId.SubscriptionId, ruleSetId.ResourceGroupName, ruleSetId.ProfileName, ruleSetId.RuleSetName)
			if err := ensureBatchRuleSetMode(ctx, metadata.Client, legacyRuleSetId); err != nil {
				return err
			}

			existing, err := client.Get(ctx, legacyRuleSetId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", legacyRuleSetId, err)
			}
			if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.Rules != nil && len(*existing.Model.Properties.Rules) > 0 {
				return metadata.ResourceRequiresImport(r.ResourceType(), *ruleSetId)
			}

			payload, err := expandCdnFrontDoorBatchRuleSetPayload(true, model, ruleSetId.RuleSetName)
			if err != nil {
				return err
			}

			if err := client.CreateThenPoll(ctx, legacyRuleSetId, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", legacyRuleSetId, err)
			}

			metadata.SetID(ruleSetId)
			return nil
		},
	}
}

func (r CdnFrontDoorBatchRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.FrontDoorRuleSetsClient_v2025_12_01

			id, err := rules.ParseRuleSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			legacyRuleSetId := legacyrulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)
			resp, err := client.Get(ctx, legacyRuleSetId)
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

func (r CdnFrontDoorBatchRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 4 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.FrontDoorRuleSetsClient_v2025_12_01

			id, err := rules.ParseRuleSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CdnFrontDoorBatchRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			legacyRuleSetId := legacyrulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)
			if err := ensureBatchRuleSetMode(ctx, metadata.Client, legacyRuleSetId); err != nil {
				return err
			}

			payload, err := expandCdnFrontDoorBatchRuleSetPayload(true, model, id.RuleSetName)
			if err != nil {
				return err
			}

			if err := client.CreateThenPoll(ctx, legacyRuleSetId, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r CdnFrontDoorBatchRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.FrontDoorRuleSetsClient_v2025_12_01
			rulesClient := metadata.Client.Cdn.FrontDoorRulesClient_v2025_12_01

			var model CdnFrontDoorBatchRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := rules.ParseRuleSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			legacyRuleSetId := legacyrulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)

			resp, err := client.Get(ctx, legacyRuleSetId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil
				}
				return fmt.Errorf("retrieving %s: %+v", legacyRuleSetId, err)
			}

			batchModeEnabled := true
			if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.BatchMode != nil {
				batchModeEnabled = *resp.Model.Properties.BatchMode
			}

			payload := azuresdkhacks.RuleSet2025{
				Properties: &azuresdkhacks.RuleSetProperties2025{
					BatchMode: &batchModeEnabled,
					Rules:     &[]rules.Rule{},
				},
			}

			// Batch rules are managed through the parent Rule Set PUT API. Deleting this
			// Terraform resource therefore means updating the Rule Set with an empty
			// `rules` collection rather than deleting the parent Rule Set itself.
			if err := client.CreateThenPoll(ctx, legacyRuleSetId, payload); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			for _, batchRule := range model.Rules {
				ruleId := rules.NewRuleID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName, batchRule.Name)
				pollerType := custompollers.NewFrontDoorBatchRuleDeletePoller(rulesClient, ruleId)
				poller := pollers.NewPoller(pollerType, 30*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("waiting for deletion of %s: %+v", ruleId, err)
				}
			}

			return nil
		},
	}
}

func ensureBatchRuleSetMode(ctx context.Context, client *clients.Client, id legacyrulesets.RuleSetId) error {
	resp, err := client.Cdn.FrontDoorRuleSetsClient_v2025_12_01.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || !pointer.From(resp.Model.Properties.BatchMode) {
		return fmt.Errorf("creating or updating a Front Door Batch Rule in %s is not supported unless `batch_mode_enabled` is `true` on the parent Rule Set", id)
	}

	return nil
}
