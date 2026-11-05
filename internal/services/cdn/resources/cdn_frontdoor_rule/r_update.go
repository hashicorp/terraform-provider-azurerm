// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (c Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 4 * time.Hour,
		Func: func(ctx context.Context, rmd sdk.ResourceMetaData) error {
			client := rmd.Client.Cdn.FrontDoorRulesClient

			state := CdnFrontDoorRuleResourceModel{}
			if err := rmd.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %w", err)
			}

			id, err := rules.ParseRuleID(rmd.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}
			props := existing.Model.Properties

			if rmd.ResourceData.HasChange("behaviour_on_match") {
				props.MatchProcessingBehavior = pointer.ToEnum[rules.MatchProcessingBehavior](state.BehaviourOnMatch)
			}

			if rmd.ResourceData.HasChange("order") {
				props.Order = pointer.To(state.Order)
			}

			if rmd.ResourceData.HasChange("actions") {
				actions, err := expandCdnFrontDoorRuleActions(state.Actions)
				if err != nil {
					return fmt.Errorf("expanding `actions`: %+v", err)
				}
				props.Actions = &actions
			}

			if rmd.ResourceData.HasChange("conditions") {
				conditions, err := expandCdnFrontDoorRuleConditions(state.Conditions)
				if err != nil {
					return fmt.Errorf("expanding `conditions`: %+v", err)
				}

				if len(conditions) > 10 {
					return fmt.Errorf("expanding `conditions`: configuration file exceeds the maximum of 10 match conditions, got %d", len(conditions))
				}

				props.Conditions = &conditions
			}

			if err := client.CreateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
