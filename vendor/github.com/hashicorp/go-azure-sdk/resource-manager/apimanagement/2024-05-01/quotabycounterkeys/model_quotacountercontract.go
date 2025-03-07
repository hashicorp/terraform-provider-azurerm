package quotabycounterkeys

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuotaCounterContract struct {
	CounterKey      string                               `json:"counterKey"`
	PeriodEndTime   string                               `json:"periodEndTime"`
	PeriodKey       string                               `json:"periodKey"`
	PeriodStartTime string                               `json:"periodStartTime"`
	Value           *QuotaCounterValueContractProperties `json:"value,omitempty"`
}

func (o *QuotaCounterContract) GetPeriodEndTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.PeriodEndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *QuotaCounterContract) SetPeriodEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PeriodEndTime = formatted
}

func (o *QuotaCounterContract) GetPeriodStartTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.PeriodStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *QuotaCounterContract) SetPeriodStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PeriodStartTime = formatted
}
