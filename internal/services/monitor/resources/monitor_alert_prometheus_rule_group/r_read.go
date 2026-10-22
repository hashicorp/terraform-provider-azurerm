// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_alert_prometheus_rule_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2023-03-01/prometheusrulegroupresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.AlertPrometheusRuleGroupClient

			id, err := prometheusrulegroupresources.ParsePrometheusRuleGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.PrometheusRuleGroupsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := AlertPrometheusRuleGroupResourceModel{
				Name:              id.PrometheusRuleGroupName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if resp.Model != nil {
				state.ClusterName = pointer.From(resp.Model.Properties.ClusterName)
				state.Description = pointer.From(resp.Model.Properties.Description)
				state.Interval = pointer.From(resp.Model.Properties.Interval)
				state.Location = location.Normalize(resp.Model.Location)
				state.Rule = flattenPrometheusRuleModel(&resp.Model.Properties.Rules)
				state.RuleGroupEnabled = pointer.From(resp.Model.Properties.Enabled)
				state.Scopes = resp.Model.Properties.Scopes
				state.Tags = pointer.From(resp.Model.Tags)
			}
			return metadata.Encode(&state)
		},
	}
}
