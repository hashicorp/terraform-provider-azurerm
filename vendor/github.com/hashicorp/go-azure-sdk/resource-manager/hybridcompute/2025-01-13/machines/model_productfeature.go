package machines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProductFeature struct {
	BillingEndDate     *string                           `json:"billingEndDate,omitempty"`
	BillingStartDate   *string                           `json:"billingStartDate,omitempty"`
	DisenrollmentDate  *string                           `json:"disenrollmentDate,omitempty"`
	EnrollmentDate     *string                           `json:"enrollmentDate,omitempty"`
	Error              *ErrorDetail                      `json:"error,omitempty"`
	Name               *string                           `json:"name,omitempty"`
	SubscriptionStatus *LicenseProfileSubscriptionStatus `json:"subscriptionStatus,omitempty"`
}

func (o *ProductFeature) GetBillingEndDateAsTime() (*time.Time, error) {
	if o.BillingEndDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.BillingEndDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ProductFeature) SetBillingEndDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.BillingEndDate = &formatted
}

func (o *ProductFeature) GetBillingStartDateAsTime() (*time.Time, error) {
	if o.BillingStartDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.BillingStartDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ProductFeature) SetBillingStartDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.BillingStartDate = &formatted
}

func (o *ProductFeature) GetDisenrollmentDateAsTime() (*time.Time, error) {
	if o.DisenrollmentDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DisenrollmentDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ProductFeature) SetDisenrollmentDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DisenrollmentDate = &formatted
}

func (o *ProductFeature) GetEnrollmentDateAsTime() (*time.Time, error) {
	if o.EnrollmentDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EnrollmentDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ProductFeature) SetEnrollmentDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EnrollmentDate = &formatted
}
