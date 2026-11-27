// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_microsoft_threat_protection

type DataConnectorMicrosoftThreatProtectionModel struct {
	Name                    string `tfschema:"name"`
	LogAnalyticsWorkspaceId string `tfschema:"log_analytics_workspace_id"`
	TenantId                string `tfschema:"tenant_id"`
}

func (r Resource) ModelObject() any {
	return &DataConnectorMicrosoftThreatProtectionModel{}
}
