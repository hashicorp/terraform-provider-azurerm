// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dynatrace_tag_rules

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.TagRulesClient
			id, err := tagrules.ParseTagRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state TagRulesResourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := existing.Model

			if metadata.ResourceData.HasChange("metric_rule") {
				model.Properties.MetricRules = ExpandMetricRules(state.MetricRules)
			}

			if metadata.ResourceData.HasChange("log_rule") {
				model.Properties.LogRules = ExpandLogRule(state.LogRules)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
