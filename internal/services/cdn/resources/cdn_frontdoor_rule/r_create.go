// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (c Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 4 * time.Hour,
		Func: func(ctx context.Context, rmd sdk.ResourceMetaData) error {
			client := rmd.Client.Cdn.FrontDoorRulesClient
			ruleSetsClient := rmd.Client.Cdn.FrontDoorRuleSetsClient

			config := CdnFrontDoorRuleResourceModel{}
			if err := rmd.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %w", err)
			}

			ruleSetId, err := rulesets.ParseRuleSetID(config.CdnFrontDoorRuleSetID)
			if err != nil {
				return err
			}

			id := rules.NewRuleID(ruleSetId.SubscriptionId, ruleSetId.ResourceGroupName, ruleSetId.ProfileName, ruleSetId.RuleSetName, config.Name)

			ruleSet, err := ruleSetsClient.Get(ctx, *ruleSetId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", ruleSetId, err)
			}

			if ruleSet.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", ruleSetId)
			}

			if ruleSet.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", ruleSetId)
			}

			if pointer.From(ruleSet.Model.Properties.BatchMode) {
				return fmt.Errorf("the parent ruleset (%s) was provisioned using batch mode, and individual rules for this cannot be managed by this resource, use `azurerm_cdn_frontdoor_batch_rule_set` instead, or create a non-batch Rule Set with `azurerm_cdn_frontdoor_rule_set`", ruleSetId)
			}

			if !rmd.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				result, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(result.HttpResponse) {
						return fmt.Errorf("checking for existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(result.HttpResponse) {
					return rmd.ResourceRequiresImport("azurerm_cdn_frontdoor_rule", id)
				}
			}

			actions, err := expandCdnFrontDoorRuleActions(config.Actions)
			if err != nil {
				return fmt.Errorf("expanding `actions`: %+v", err)
			}

			conditions, err := expandCdnFrontDoorRuleConditions(config.Conditions)
			if err != nil {
				return fmt.Errorf("expanding `conditions`: %+v", err)
			}

			props := rules.Rule{
				Properties: &rules.RuleProperties{
					Actions:                 &actions,
					Conditions:              &conditions,
					MatchProcessingBehavior: pointer.ToEnum[rules.MatchProcessingBehavior](config.BehaviourOnMatch),
					RuleSetName:             &ruleSetId.RuleSetName,
					Order:                   pointer.To(config.Order),
				},
			}

			if err := client.CreateCallbackThenPoll(ctx, id, props, rmd.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			rmd.SetID(id)

			return nil
		},
	}
}
