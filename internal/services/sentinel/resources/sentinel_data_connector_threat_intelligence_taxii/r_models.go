// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_threat_intelligence_taxii

type DataConnectorThreatIntelligenceTAXIIModel struct {
	Name                    string `tfschema:"name"`
	LogAnalyticsWorkspaceId string `tfschema:"log_analytics_workspace_id"`
	DisplayName             string `tfschema:"display_name"`
	APIRootURL              string `tfschema:"api_root_url"`
	CollectionID            string `tfschema:"collection_id"`
	UserName                string `tfschema:"user_name"`
	Password                string `tfschema:"password"`
	PollingFrequency        string `tfschema:"polling_frequency"`
	LookbackDate            string `tfschema:"lookback_date"`
	TenantId                string `tfschema:"tenant_id"`
}

func (r Resource) ModelObject() any {
	return &DataConnectorThreatIntelligenceTAXIIModel{}
}
