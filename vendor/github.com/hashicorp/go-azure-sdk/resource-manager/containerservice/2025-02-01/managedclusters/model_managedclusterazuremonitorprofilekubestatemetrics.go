package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterAzureMonitorProfileKubeStateMetrics struct {
	MetricAnnotationsAllowList *string `json:"metricAnnotationsAllowList,omitempty"`
	MetricLabelsAllowlist      *string `json:"metricLabelsAllowlist,omitempty"`
}
