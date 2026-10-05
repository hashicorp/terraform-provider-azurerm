package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitoringSettingProperties struct {
	AppInsightsAgentVersions      *ApplicationInsightsAgentVersions `json:"appInsightsAgentVersions,omitempty"`
	AppInsightsInstrumentationKey *string                           `json:"appInsightsInstrumentationKey,omitempty"`
	AppInsightsSamplingRate       *float64                          `json:"appInsightsSamplingRate,omitempty"`
	Error                         *Error                            `json:"error,omitempty"`
	ProvisioningState             *MonitoringSettingState           `json:"provisioningState,omitempty"`
	TraceEnabled                  *bool                             `json:"traceEnabled,omitempty"`
}
