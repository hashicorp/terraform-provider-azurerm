package saplandscapemonitor

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SapLandscapeMonitorProperties struct {
	Grouping             *SapLandscapeMonitorPropertiesGrouping `json:"grouping,omitempty"`
	ProvisioningState    *SapLandscapeMonitorProvisioningState  `json:"provisioningState,omitempty"`
	TopMetricsThresholds *[]SapLandscapeMonitorMetricThresholds `json:"topMetricsThresholds,omitempty"`
}
