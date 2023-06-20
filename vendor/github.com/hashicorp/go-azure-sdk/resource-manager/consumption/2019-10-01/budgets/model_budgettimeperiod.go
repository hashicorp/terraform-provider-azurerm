package budgets

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BudgetTimePeriod struct {
	EndDate   *string `json:"endDate,omitempty"`
	StartDate string  `json:"startDate"`
}

func (o *BudgetTimePeriod) GetEndDateAsTime() (*time.Time, error) {
	if o.EndDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndDate, "2006-01-02T15:04:05Z07:00")
}

func (o *BudgetTimePeriod) SetEndDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndDate = &formatted
}

func (o *BudgetTimePeriod) GetStartDateAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.StartDate, "2006-01-02T15:04:05Z07:00")
}

func (o *BudgetTimePeriod) SetStartDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartDate = formatted
}
