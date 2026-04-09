package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyInsights struct {
	IsEnabled             *bool                                `json:"isEnabled,omitempty"`
	LogAnalyticsResources *FirewallPolicyLogAnalyticsResources `json:"logAnalyticsResources,omitempty"`
	RetentionDays         *int64                               `json:"retentionDays,omitempty"`
}
