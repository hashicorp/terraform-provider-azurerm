package networkwatchers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureReachabilityReportParameters struct {
	AzureLocations   *[]string                       `json:"azureLocations,omitempty"`
	EndTime          string                          `json:"endTime"`
	ProviderLocation AzureReachabilityReportLocation `json:"providerLocation"`
	Providers        *[]string                       `json:"providers,omitempty"`
	StartTime        string                          `json:"startTime"`
}

func (o *AzureReachabilityReportParameters) GetEndTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AzureReachabilityReportParameters) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = formatted
}

func (o *AzureReachabilityReportParameters) GetStartTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AzureReachabilityReportParameters) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = formatted
}
