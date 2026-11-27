// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_metrics

type NextGenerationFirewallMetricsModel struct {
	FirewallID                          string `tfschema:"firewall_id"`
	ApplicationInsightsConnectionString string `tfschema:"application_insights_connection_string"`
	ApplicationInsightsID               string `tfschema:"application_insights_id"`
}

func (r Resource) ModelObject() any {
	return &NextGenerationFirewallMetricsModel{}
}
