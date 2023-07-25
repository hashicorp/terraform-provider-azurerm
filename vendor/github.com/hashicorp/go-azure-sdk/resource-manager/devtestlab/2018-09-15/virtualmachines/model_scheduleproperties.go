package virtualmachines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleProperties struct {
	CreatedDate          *string               `json:"createdDate,omitempty"`
	DailyRecurrence      *DayDetails           `json:"dailyRecurrence,omitempty"`
	HourlyRecurrence     *HourDetails          `json:"hourlyRecurrence,omitempty"`
	NotificationSettings *NotificationSettings `json:"notificationSettings,omitempty"`
	ProvisioningState    *string               `json:"provisioningState,omitempty"`
	Status               *EnableStatus         `json:"status,omitempty"`
	TargetResourceId     *string               `json:"targetResourceId,omitempty"`
	TaskType             *string               `json:"taskType,omitempty"`
	TimeZoneId           *string               `json:"timeZoneId,omitempty"`
	UniqueIdentifier     *string               `json:"uniqueIdentifier,omitempty"`
	WeeklyRecurrence     *WeekDetails          `json:"weeklyRecurrence,omitempty"`
}

func (o *ScheduleProperties) GetCreatedDateAsTime() (*time.Time, error) {
	if o.CreatedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduleProperties) SetCreatedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDate = &formatted
}
