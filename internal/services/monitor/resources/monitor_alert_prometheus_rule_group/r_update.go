// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_alert_prometheus_rule_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2023-03-01/prometheusrulegroupresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.AlertPrometheusRuleGroupClient

			id, err := prometheusrulegroupresources.ParsePrometheusRuleGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model AlertPrometheusRuleGroupResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.PrometheusRuleGroupsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("cluster_name") {
				properties.Properties.ClusterName = pointer.To(model.ClusterName)
			}
			if metadata.ResourceData.HasChange("description") {
				properties.Properties.Description = pointer.To(model.Description)
			}
			if metadata.ResourceData.HasChange("rule_group_enabled") {
				properties.Properties.Enabled = pointer.To(model.RuleGroupEnabled)
			}
			if metadata.ResourceData.HasChange("interval") {
				properties.Properties.Interval = pointer.To(model.Interval)
			}
			if metadata.ResourceData.HasChange("rule") {
				properties.Properties.Rules = expandPrometheusRuleModel(model.Rule, metadata.ResourceData)
			}
			if metadata.ResourceData.HasChange("scopes") {
				properties.Properties.Scopes = model.Scopes
			}
			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = pointer.To(model.Tags)
			}

			if _, err := client.PrometheusRuleGroupsCreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
