package subaccount

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitoredResource struct {
	Id                     *string                `json:"id,omitempty"`
	ReasonForLogsStatus    *string                `json:"reasonForLogsStatus,omitempty"`
	ReasonForMetricsStatus *string                `json:"reasonForMetricsStatus,omitempty"`
	SendingLogs            *bool                  `json:"sendingLogs,omitempty"`
	SendingMetrics         *bool                  `json:"sendingMetrics,omitempty"`
	SystemData             *systemdata.SystemData `json:"systemData,omitempty"`
}
