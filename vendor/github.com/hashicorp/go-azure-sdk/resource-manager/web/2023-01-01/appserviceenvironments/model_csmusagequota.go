package appserviceenvironments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CsmUsageQuota struct {
	CurrentValue  *int64             `json:"currentValue,omitempty"`
	Limit         *int64             `json:"limit,omitempty"`
	Name          *LocalizableString `json:"name,omitempty"`
	NextResetTime *string            `json:"nextResetTime,omitempty"`
	Unit          *string            `json:"unit,omitempty"`
}

func (o *CsmUsageQuota) GetNextResetTimeAsTime() (*time.Time, error) {
	if o.NextResetTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.NextResetTime, "2006-01-02T15:04:05Z07:00")
}

func (o *CsmUsageQuota) SetNextResetTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.NextResetTime = &formatted
}
