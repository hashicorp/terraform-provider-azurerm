package registries

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetentionPolicy struct {
	Days            *int64        `json:"days,omitempty"`
	LastUpdatedTime *string       `json:"lastUpdatedTime,omitempty"`
	Status          *PolicyStatus `json:"status,omitempty"`
}

func (o *RetentionPolicy) GetLastUpdatedTimeAsTime() (*time.Time, error) {
	if o.LastUpdatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RetentionPolicy) SetLastUpdatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTime = &formatted
}
