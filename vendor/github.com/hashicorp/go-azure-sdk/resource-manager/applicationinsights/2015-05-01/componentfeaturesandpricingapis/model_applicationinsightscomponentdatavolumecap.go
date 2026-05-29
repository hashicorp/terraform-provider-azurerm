package componentfeaturesandpricingapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationInsightsComponentDataVolumeCap struct {
	Cap                                  *float64 `json:"Cap,omitempty"`
	MaxHistoryCap                        *float64 `json:"MaxHistoryCap,omitempty"`
	ResetTime                            *int64   `json:"ResetTime,omitempty"`
	StopSendNotificationWhenHitCap       *bool    `json:"StopSendNotificationWhenHitCap,omitempty"`
	StopSendNotificationWhenHitThreshold *bool    `json:"StopSendNotificationWhenHitThreshold,omitempty"`
	WarningThreshold                     *int64   `json:"WarningThreshold,omitempty"`
}
