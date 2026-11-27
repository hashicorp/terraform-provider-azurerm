// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_iot

type DataConnectorIOTModel struct {
	Name                    string `tfschema:"name"`
	LogAnalyticsWorkspaceId string `tfschema:"log_analytics_workspace_id"`
	SubscriptionId          string `tfschema:"subscription_id"`
}

func (r Resource) ModelObject() any {
	return &DataConnectorIOTModel{}
}
