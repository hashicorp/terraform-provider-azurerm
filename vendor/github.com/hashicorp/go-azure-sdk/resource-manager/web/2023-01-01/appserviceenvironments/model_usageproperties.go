package appserviceenvironments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UsageProperties struct {
	ComputeMode   *ComputeModeOptions `json:"computeMode,omitempty"`
	CurrentValue  *int64              `json:"currentValue,omitempty"`
	DisplayName   *string             `json:"displayName,omitempty"`
	Limit         *int64              `json:"limit,omitempty"`
	NextResetTime *string             `json:"nextResetTime,omitempty"`
	ResourceName  *string             `json:"resourceName,omitempty"`
	SiteMode      *string             `json:"siteMode,omitempty"`
	Unit          *string             `json:"unit,omitempty"`
}

func (o *UsageProperties) GetNextResetTimeAsTime() (*time.Time, error) {
	if o.NextResetTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.NextResetTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UsageProperties) SetNextResetTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.NextResetTime = &formatted
}
