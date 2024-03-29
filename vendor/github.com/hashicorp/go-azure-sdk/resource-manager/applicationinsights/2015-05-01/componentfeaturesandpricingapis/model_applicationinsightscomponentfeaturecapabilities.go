package componentfeaturesandpricingapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationInsightsComponentFeatureCapabilities struct {
	AnalyticsIntegration *bool    `json:"AnalyticsIntegration,omitempty"`
	ApiAccessLevel       *string  `json:"ApiAccessLevel,omitempty"`
	ApplicationMap       *bool    `json:"ApplicationMap,omitempty"`
	BurstThrottlePolicy  *string  `json:"BurstThrottlePolicy,omitempty"`
	DailyCap             *float64 `json:"DailyCap,omitempty"`
	DailyCapResetTime    *float64 `json:"DailyCapResetTime,omitempty"`
	LiveStreamMetrics    *bool    `json:"LiveStreamMetrics,omitempty"`
	MetadataClass        *string  `json:"MetadataClass,omitempty"`
	MultipleStepWebTest  *bool    `json:"MultipleStepWebTest,omitempty"`
	OpenSchema           *bool    `json:"OpenSchema,omitempty"`
	PowerBIIntegration   *bool    `json:"PowerBIIntegration,omitempty"`
	ProactiveDetection   *bool    `json:"ProactiveDetection,omitempty"`
	SupportExportData    *bool    `json:"SupportExportData,omitempty"`
	ThrottleRate         *float64 `json:"ThrottleRate,omitempty"`
	TrackingType         *string  `json:"TrackingType,omitempty"`
	WorkItemIntegration  *bool    `json:"WorkItemIntegration,omitempty"`
}
