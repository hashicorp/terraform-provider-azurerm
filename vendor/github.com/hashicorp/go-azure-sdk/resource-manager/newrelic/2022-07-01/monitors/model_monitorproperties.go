package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorProperties struct {
	AccountCreationSource         *AccountCreationSource         `json:"accountCreationSource,omitempty"`
	LiftrResourceCategory         *LiftrResourceCategories       `json:"liftrResourceCategory,omitempty"`
	LiftrResourcePreference       *int64                         `json:"liftrResourcePreference,omitempty"`
	MarketplaceSubscriptionId     *string                        `json:"marketplaceSubscriptionId,omitempty"`
	MarketplaceSubscriptionStatus *MarketplaceSubscriptionStatus `json:"marketplaceSubscriptionStatus,omitempty"`
	MonitoringStatus              *MonitoringStatus              `json:"monitoringStatus,omitempty"`
	NewRelicAccountProperties     *NewRelicAccountProperties     `json:"newRelicAccountProperties,omitempty"`
	OrgCreationSource             *OrgCreationSource             `json:"orgCreationSource,omitempty"`
	PlanData                      *PlanData                      `json:"planData,omitempty"`
	ProvisioningState             *ProvisioningState             `json:"provisioningState,omitempty"`
	UserInfo                      *UserInfo                      `json:"userInfo,omitempty"`
}
