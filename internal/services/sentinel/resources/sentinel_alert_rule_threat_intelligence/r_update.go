// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_threat_intelligence

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
)

func (a Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config AlertRuleThreatIntelligenceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := sentinelClient.From(metadata.Client).AlertRulesClient

			id, err := alertrules.ParseAlertRuleID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %+v", err)
			}
			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if err := helpers.AssertAlertRuleKind(resp.Model, alertrules.AlertRuleKindThreatIntelligence); err != nil {
				return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
			}

			rule := resp.Model.(alertrules.ThreatIntelligenceAlertRule)

			if rule.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}
			props := rule.Properties

			if metadata.ResourceData.HasChange("enabled") {
				props.Enabled = config.Enabled
			}

			param := alertrules.ThreatIntelligenceAlertRule{
				Properties: &alertrules.ThreatIntelligenceAlertRuleProperties{
					Enabled:               props.Enabled,
					AlertRuleTemplateName: props.AlertRuleTemplateName,
					Severity:              props.Severity,
					DisplayName:           props.DisplayName,
					Description:           props.Description,
					Tactics:               props.Tactics,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %q: %+v", id, err)
			}

			return nil
		},
	}
}
