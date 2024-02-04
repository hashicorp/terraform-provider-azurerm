package billingaccounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Enrollment struct {
	BillingCycle *string             `json:"billingCycle,omitempty"`
	Channel      *string             `json:"channel,omitempty"`
	CountryCode  *string             `json:"countryCode,omitempty"`
	Currency     *string             `json:"currency,omitempty"`
	EndDate      *string             `json:"endDate,omitempty"`
	Language     *string             `json:"language,omitempty"`
	Policies     *EnrollmentPolicies `json:"policies,omitempty"`
	StartDate    *string             `json:"startDate,omitempty"`
	Status       *string             `json:"status,omitempty"`
}

func (o *Enrollment) GetEndDateAsTime() (*time.Time, error) {
	if o.EndDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndDate, "2006-01-02T15:04:05Z07:00")
}

func (o *Enrollment) SetEndDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndDate = &formatted
}

func (o *Enrollment) GetStartDateAsTime() (*time.Time, error) {
	if o.StartDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartDate, "2006-01-02T15:04:05Z07:00")
}

func (o *Enrollment) SetStartDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartDate = &formatted
}
