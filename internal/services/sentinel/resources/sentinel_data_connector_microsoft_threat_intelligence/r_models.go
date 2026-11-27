// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_microsoft_threat_intelligence

type DataConnectorMicrosoftThreatIntelligenceModel struct {
	Name                                    string `tfschema:"name"`
	WorkspaceId                             string `tfschema:"log_analytics_workspace_id"`
	TenantId                                string `tfschema:"tenant_id"`
	MicrosoftEmergingThreatFeedLookBackDate string `tfschema:"microsoft_emerging_threat_feed_lookback_date"`
}

func (s Resource) ModelObject() any {
	return &DataConnectorMicrosoftThreatIntelligenceModel{}
}
