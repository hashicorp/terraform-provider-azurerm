package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterLogAnalyticsProfile struct {
	ApplicationLogs *ClusterLogAnalyticsApplicationLogs `json:"applicationLogs,omitempty"`
	Enabled         bool                                `json:"enabled"`
	MetricsEnabled  *bool                               `json:"metricsEnabled,omitempty"`
}
