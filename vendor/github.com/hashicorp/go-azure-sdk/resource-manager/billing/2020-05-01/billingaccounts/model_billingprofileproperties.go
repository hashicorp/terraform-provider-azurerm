package billingaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BillingProfileProperties struct {
	BillTo                   *AddressDetails           `json:"billTo,omitempty"`
	BillingRelationshipType  *BillingRelationshipType  `json:"billingRelationshipType,omitempty"`
	Currency                 *string                   `json:"currency,omitempty"`
	DisplayName              *string                   `json:"displayName,omitempty"`
	EnabledAzurePlans        *[]AzurePlan              `json:"enabledAzurePlans,omitempty"`
	HasReadAccess            *bool                     `json:"hasReadAccess,omitempty"`
	IndirectRelationshipInfo *IndirectRelationshipInfo `json:"indirectRelationshipInfo,omitempty"`
	InvoiceDay               *int64                    `json:"invoiceDay,omitempty"`
	InvoiceEmailOptIn        *bool                     `json:"invoiceEmailOptIn,omitempty"`
	InvoiceSections          *InvoiceSectionsOnExpand  `json:"invoiceSections,omitempty"`
	PoNumber                 *string                   `json:"poNumber,omitempty"`
	SpendingLimit            *SpendingLimit            `json:"spendingLimit,omitempty"`
	Status                   *BillingProfileStatus     `json:"status,omitempty"`
	StatusReasonCode         *StatusReasonCode         `json:"statusReasonCode,omitempty"`
	SystemId                 *string                   `json:"systemId,omitempty"`
	Tags                     *map[string]string        `json:"tags,omitempty"`
	TargetClouds             *[]TargetCloud            `json:"targetClouds,omitempty"`
}
