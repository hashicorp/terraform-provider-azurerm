package subscriptiondiagnosticsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionDiagnosticSettings struct {
	EventHubAuthorizationRuleId *string                    `json:"eventHubAuthorizationRuleId,omitempty"`
	EventHubName                *string                    `json:"eventHubName,omitempty"`
	Logs                        *[]SubscriptionLogSettings `json:"logs,omitempty"`
	MarketplacePartnerId        *string                    `json:"marketplacePartnerId,omitempty"`
	ServiceBusRuleId            *string                    `json:"serviceBusRuleId,omitempty"`
	StorageAccountId            *string                    `json:"storageAccountId,omitempty"`
	WorkspaceId                 *string                    `json:"workspaceId,omitempty"`
}
