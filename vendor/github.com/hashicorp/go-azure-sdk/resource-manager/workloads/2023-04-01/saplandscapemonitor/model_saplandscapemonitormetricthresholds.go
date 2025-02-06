package saplandscapemonitor

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SapLandscapeMonitorMetricThresholds struct {
	Green  *float64 `json:"green,omitempty"`
	Name   *string  `json:"name,omitempty"`
	Red    *float64 `json:"red,omitempty"`
	Yellow *float64 `json:"yellow,omitempty"`
}
