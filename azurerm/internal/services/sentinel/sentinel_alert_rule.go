package sentinel

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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
	default:
		return nil
	}
}

func importSentinelAlertRule(expectKind securityinsight.AlertRuleKind) func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
	return func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
		id, err := parse.AlertRuleID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).Sentinel.AlertRulesClient
		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		resp, err := client.Get(ctx, id.ResourceGroup, "Microsoft.OperationalInsights", id.WorkspaceName, id.Name)
		if err != nil {
			return nil, fmt.Errorf("retrieving Sentinel Alert Rule %q: %+v", id, err)
		}

		if err := assertAlertRuleKind(resp.Value, expectKind); err != nil {
			return nil, err
		}
		return []*schema.ResourceData{d}, nil
	}
}

func assertAlertRuleKind(rule securityinsight.BasicAlertRule, expectKind securityinsight.AlertRuleKind) error {
	var kind securityinsight.AlertRuleKind
	switch rule.(type) {
	case securityinsight.FusionAlertRule:
		kind = securityinsight.Fusion
	case securityinsight.MicrosoftSecurityIncidentCreationAlertRule:
		kind = securityinsight.MicrosoftSecurityIncidentCreation
	case securityinsight.ScheduledAlertRule:
		kind = securityinsight.Scheduled
	}
	if expectKind != kind {
		return fmt.Errorf("Sentinel Alert Rule has mismatched kind, expected: %q, got %q", expectKind, kind)
	}
	return nil
}
