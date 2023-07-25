package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrafficAnalyticsConfigurationProperties struct {
	Enabled                  *bool   `json:"enabled,omitempty"`
	TrafficAnalyticsInterval *int64  `json:"trafficAnalyticsInterval,omitempty"`
	WorkspaceId              *string `json:"workspaceId,omitempty"`
	WorkspaceRegion          *string `json:"workspaceRegion,omitempty"`
	WorkspaceResourceId      *string `json:"workspaceResourceId,omitempty"`
}
