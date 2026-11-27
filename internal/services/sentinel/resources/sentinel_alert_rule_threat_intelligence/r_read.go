// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_threat_intelligence

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
)

func (a Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.AlertRulesClient

			id, err := alertrules.ParseAlertRuleID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if err := helpers.AssertAlertRuleKind(resp.Model, alertrules.AlertRuleKindThreatIntelligence); err != nil {
				return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
			}

			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)

			state := AlertRuleThreatIntelligenceModel{
				Name:        id.RuleId,
				WorkspaceId: workspaceId.ID(),
			}

			if rule, ok := resp.Model.(alertrules.ThreatIntelligenceAlertRule); ok {
				if prop := rule.Properties; prop != nil {
					state.Enabled = prop.Enabled
					state.TemplateName = prop.AlertRuleTemplateName
				}
			}

			return metadata.Encode(&state)
		},
	}
}
