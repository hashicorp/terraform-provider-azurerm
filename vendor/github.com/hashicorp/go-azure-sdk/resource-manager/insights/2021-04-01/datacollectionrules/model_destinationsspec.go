package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DestinationsSpec struct {
	AzureMonitorMetrics *AzureMonitorMetricsDestination `json:"azureMonitorMetrics,omitempty"`
	LogAnalytics        *[]LogAnalyticsDestination      `json:"logAnalytics,omitempty"`
}
