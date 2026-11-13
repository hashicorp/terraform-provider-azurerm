package elasticmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitoredResource struct {
	Id                  *string      `json:"id,omitempty"`
	ReasonForLogsStatus *string      `json:"reasonForLogsStatus,omitempty"`
	SendingLogs         *SendingLogs `json:"sendingLogs,omitempty"`
}
