package billingaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InvoiceSectionWithCreateSubPermission struct {
	BillingProfileDisplayName      *string                            `json:"billingProfileDisplayName,omitempty"`
	BillingProfileId               *string                            `json:"billingProfileId,omitempty"`
	BillingProfileSpendingLimit    *SpendingLimitForBillingProfile    `json:"billingProfileSpendingLimit,omitempty"`
	BillingProfileStatus           *BillingProfileStatus              `json:"billingProfileStatus,omitempty"`
	BillingProfileStatusReasonCode *StatusReasonCodeForBillingProfile `json:"billingProfileStatusReasonCode,omitempty"`
	BillingProfileSystemId         *string                            `json:"billingProfileSystemId,omitempty"`
	EnabledAzurePlans              *[]AzurePlan                       `json:"enabledAzurePlans,omitempty"`
	InvoiceSectionDisplayName      *string                            `json:"invoiceSectionDisplayName,omitempty"`
	InvoiceSectionId               *string                            `json:"invoiceSectionId,omitempty"`
	InvoiceSectionSystemId         *string                            `json:"invoiceSectionSystemId,omitempty"`
}
