package sentinel

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func alertRuleID(rule securityinsight.BasicAlertRule) *string {
	if rule == nil {
		return nil
	}
	switch rule := rule.(type) {
	case securityinsight.FusionAlertRule:
		return rule.ID
	case securityinsight.MicrosoftSecurityIncidentCreationAlertRule:
		return rule.ID
	case securityinsight.ScheduledAlertRule:
		return rule.ID
	case securityinsight.MLBehaviorAnalyticsAlertRule:
		return rule.ID
	default:
		return nil
	}
}

func importSentinelAlertRule(expectKind securityinsight.AlertRuleKind) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := parse.AlertRuleID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).Sentinel.AlertRulesClient
		resp, err := client.Get(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, id.Name)
		if err != nil {
			return nil, fmt.Errorf("retrieving Sentinel Alert Rule %q: %+v", id, err)
		}

		if err := assertAlertRuleKind(resp.Value, expectKind); err != nil {
			return nil, err
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}

func assertAlertRuleKind(rule securityinsight.BasicAlertRule, expectKind securityinsight.AlertRuleKind) error {
	var kind securityinsight.AlertRuleKind
	switch rule.(type) {
	case securityinsight.MLBehaviorAnalyticsAlertRule:
		kind = securityinsight.AlertRuleKindMLBehaviorAnalytics
	case securityinsight.FusionAlertRule:
		kind = securityinsight.AlertRuleKindFusion
	case securityinsight.MicrosoftSecurityIncidentCreationAlertRule:
		kind = securityinsight.AlertRuleKindMicrosoftSecurityIncidentCreation
	case securityinsight.ScheduledAlertRule:
		kind = securityinsight.AlertRuleKindScheduled
	}
	if expectKind != kind {
		return fmt.Errorf("Sentinel Alert Rule has mismatched kind, expected: %q, got %q", expectKind, kind)
	}
	return nil
}
