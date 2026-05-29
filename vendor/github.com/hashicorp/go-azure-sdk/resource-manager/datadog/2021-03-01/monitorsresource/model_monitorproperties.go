package monitorsresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorProperties struct {
	DatadogOrganizationProperties *DatadogOrganizationProperties `json:"datadogOrganizationProperties,omitempty"`
	LiftrResourceCategory         *LiftrResourceCategories       `json:"liftrResourceCategory,omitempty"`
	LiftrResourcePreference       *int64                         `json:"liftrResourcePreference,omitempty"`
	MarketplaceSubscriptionStatus *MarketplaceSubscriptionStatus `json:"marketplaceSubscriptionStatus,omitempty"`
	MonitoringStatus              *MonitoringStatus              `json:"monitoringStatus,omitempty"`
	ProvisioningState             *ProvisioningState             `json:"provisioningState,omitempty"`
	UserInfo                      *UserInfo                      `json:"userInfo,omitempty"`
}
