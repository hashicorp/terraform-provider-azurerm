package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitoredResource struct {
	Id                     *string               `json:"id,omitempty"`
	ReasonForLogsStatus    *string               `json:"reasonForLogsStatus,omitempty"`
	ReasonForMetricsStatus *string               `json:"reasonForMetricsStatus,omitempty"`
	SendingLogs            *SendingLogsStatus    `json:"sendingLogs,omitempty"`
	SendingMetrics         *SendingMetricsStatus `json:"sendingMetrics,omitempty"`
}
