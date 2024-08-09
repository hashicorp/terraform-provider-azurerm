package machines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LicenseProfileArmProductProfileProperties struct {
	BillingEndDate     *string                           `json:"billingEndDate,omitempty"`
	BillingStartDate   *string                           `json:"billingStartDate,omitempty"`
	DisenrollmentDate  *string                           `json:"disenrollmentDate,omitempty"`
	EnrollmentDate     *string                           `json:"enrollmentDate,omitempty"`
	Error              *ErrorDetail                      `json:"error,omitempty"`
	ProductFeatures    *[]ProductFeature                 `json:"productFeatures,omitempty"`
	ProductType        *LicenseProfileProductType        `json:"productType,omitempty"`
	SubscriptionStatus *LicenseProfileSubscriptionStatus `json:"subscriptionStatus,omitempty"`
}

func (o *LicenseProfileArmProductProfileProperties) GetBillingEndDateAsTime() (*time.Time, error) {
	if o.BillingEndDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.BillingEndDate, "2006-01-02T15:04:05Z07:00")
}

func (o *LicenseProfileArmProductProfileProperties) SetBillingEndDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.BillingEndDate = &formatted
}

func (o *LicenseProfileArmProductProfileProperties) GetBillingStartDateAsTime() (*time.Time, error) {
	if o.BillingStartDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.BillingStartDate, "2006-01-02T15:04:05Z07:00")
}

func (o *LicenseProfileArmProductProfileProperties) SetBillingStartDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.BillingStartDate = &formatted
}

func (o *LicenseProfileArmProductProfileProperties) GetDisenrollmentDateAsTime() (*time.Time, error) {
	if o.DisenrollmentDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DisenrollmentDate, "2006-01-02T15:04:05Z07:00")
}

func (o *LicenseProfileArmProductProfileProperties) SetDisenrollmentDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DisenrollmentDate = &formatted
}

func (o *LicenseProfileArmProductProfileProperties) GetEnrollmentDateAsTime() (*time.Time, error) {
	if o.EnrollmentDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EnrollmentDate, "2006-01-02T15:04:05Z07:00")
}

func (o *LicenseProfileArmProductProfileProperties) SetEnrollmentDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EnrollmentDate = &formatted
}
