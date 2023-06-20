package querypacks

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogAnalyticsQueryPackProperties struct {
	ProvisioningState *string `json:"provisioningState,omitempty"`
	QueryPackId       *string `json:"queryPackId,omitempty"`
	TimeCreated       *string `json:"timeCreated,omitempty"`
	TimeModified      *string `json:"timeModified,omitempty"`
}

func (o *LogAnalyticsQueryPackProperties) GetTimeCreatedAsTime() (*time.Time, error) {
	if o.TimeCreated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *LogAnalyticsQueryPackProperties) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = &formatted
}

func (o *LogAnalyticsQueryPackProperties) GetTimeModifiedAsTime() (*time.Time, error) {
	if o.TimeModified == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeModified, "2006-01-02T15:04:05Z07:00")
}

func (o *LogAnalyticsQueryPackProperties) SetTimeModifiedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeModified = &formatted
}
