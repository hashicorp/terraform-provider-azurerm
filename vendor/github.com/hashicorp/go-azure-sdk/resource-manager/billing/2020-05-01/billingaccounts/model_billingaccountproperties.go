package billingaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BillingAccountProperties struct {
	AccountStatus            *AccountStatus           `json:"accountStatus,omitempty"`
	AccountType              *AccountType             `json:"accountType,omitempty"`
	AgreementType            *AgreementType           `json:"agreementType,omitempty"`
	BillingProfiles          *BillingProfilesOnExpand `json:"billingProfiles,omitempty"`
	Departments              *[]Department            `json:"departments,omitempty"`
	DisplayName              *string                  `json:"displayName,omitempty"`
	EnrollmentAccounts       *[]EnrollmentAccount     `json:"enrollmentAccounts,omitempty"`
	EnrollmentDetails        *Enrollment              `json:"enrollmentDetails,omitempty"`
	HasReadAccess            *bool                    `json:"hasReadAccess,omitempty"`
	NotificationEmailAddress *string                  `json:"notificationEmailAddress,omitempty"`
	SoldTo                   *AddressDetails          `json:"soldTo,omitempty"`
}
