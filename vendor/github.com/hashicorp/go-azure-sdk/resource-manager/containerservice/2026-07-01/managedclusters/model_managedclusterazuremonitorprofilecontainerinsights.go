package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterAzureMonitorProfileContainerInsights struct {
	ContainerNetworkLogs             *ContainerNetworkLogs `json:"containerNetworkLogs,omitempty"`
	DisablePrometheusMetricsScraping *bool                 `json:"disablePrometheusMetricsScraping,omitempty"`
	Enabled                          *bool                 `json:"enabled,omitempty"`
	LogAnalyticsWorkspaceResourceId  *string               `json:"logAnalyticsWorkspaceResourceId,omitempty"`
	SyslogPort                       *int64                `json:"syslogPort,omitempty"`
}
