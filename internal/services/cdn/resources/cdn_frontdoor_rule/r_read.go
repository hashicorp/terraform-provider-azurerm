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

func (c Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, rmd sdk.ResourceMetaData) error {
			client := rmd.Client.Cdn.FrontDoorRulesClient

			id, err := rules.ParseRuleID(rmd.ResourceData.Id())
			if err != nil {
				return err
			}

			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return rmd.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := CdnFrontDoorRuleResourceModel{
				Name:                    id.RuleName,
				CdnFrontDoorRuleSetID:   rulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName).ID(),
				CdnFrontDoorRuleSetName: id.RuleSetName,
			}

			if model := result.Model; model != nil {
				if props := model.Properties; props != nil {
					state.BehaviourOnMatch = pointer.FromEnum(props.MatchProcessingBehavior)
					state.Order = pointer.From(props.Order)

					actions, err := flattenCdnFrontDoorRuleActions(props.Actions)
					if err != nil {
						return fmt.Errorf("flattening `actions`: %+v", err)
					}
					state.Actions = actions

					conditions, err := flattenCdnFrontDoorRuleConditions(props.Conditions)
					if err != nil {
						return fmt.Errorf("flattening `conditions`: %+v", err)
					}
					state.Conditions = conditions
				}
			}

			return rmd.Encode(&state)
		},
	}
}
