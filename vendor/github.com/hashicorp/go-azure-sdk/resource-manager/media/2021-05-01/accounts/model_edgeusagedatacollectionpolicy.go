package accounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeUsageDataCollectionPolicy struct {
	DataCollectionFrequency           *string                `json:"dataCollectionFrequency,omitempty"`
	DataReportingFrequency            *string                `json:"dataReportingFrequency,omitempty"`
	EventHubDetails                   *EdgeUsageDataEventHub `json:"eventHubDetails,omitempty"`
	MaxAllowedUnreportedUsageDuration *string                `json:"maxAllowedUnreportedUsageDuration,omitempty"`
}
