// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_threat_intelligence

type AlertRuleThreatIntelligenceModel struct {
	Name         string `tfschema:"name"`
	WorkspaceId  string `tfschema:"log_analytics_workspace_id"`
	TemplateName string `tfschema:"alert_rule_template_guid"`
	Enabled      bool   `tfschema:"enabled"`
}

func (a Resource) ModelObject() any {
	return &AlertRuleThreatIntelligenceModel{}
}
