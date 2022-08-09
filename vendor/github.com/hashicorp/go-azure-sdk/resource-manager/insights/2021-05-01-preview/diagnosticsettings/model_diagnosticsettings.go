package diagnosticsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticSettings struct {
	EventHubAuthorizationRuleId *string           `json:"eventHubAuthorizationRuleId,omitempty"`
	EventHubName                *string           `json:"eventHubName,omitempty"`
	LogAnalyticsDestinationType *string           `json:"logAnalyticsDestinationType,omitempty"`
	Logs                        *[]LogSettings    `json:"logs,omitempty"`
	MarketplacePartnerId        *string           `json:"marketplacePartnerId,omitempty"`
	Metrics                     *[]MetricSettings `json:"metrics,omitempty"`
	ServiceBusRuleId            *string           `json:"serviceBusRuleId,omitempty"`
	StorageAccountId            *string           `json:"storageAccountId,omitempty"`
	WorkspaceId                 *string           `json:"workspaceId,omitempty"`
}
