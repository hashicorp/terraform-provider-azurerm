package billingaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BillingAccountOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p BillingAccountOperationPredicate) Matches(input BillingAccount) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type InvoiceSectionWithCreateSubPermissionOperationPredicate struct {
	BillingProfileDisplayName *string
	BillingProfileId          *string
	BillingProfileSystemId    *string
	InvoiceSectionDisplayName *string
	InvoiceSectionId          *string
	InvoiceSectionSystemId    *string
}

func (p InvoiceSectionWithCreateSubPermissionOperationPredicate) Matches(input InvoiceSectionWithCreateSubPermission) bool {

	if p.BillingProfileDisplayName != nil && (input.BillingProfileDisplayName == nil || *p.BillingProfileDisplayName != *input.BillingProfileDisplayName) {
		return false
	}

	if p.BillingProfileId != nil && (input.BillingProfileId == nil || *p.BillingProfileId != *input.BillingProfileId) {
		return false
	}

	if p.BillingProfileSystemId != nil && (input.BillingProfileSystemId == nil || *p.BillingProfileSystemId != *input.BillingProfileSystemId) {
		return false
	}

	if p.InvoiceSectionDisplayName != nil && (input.InvoiceSectionDisplayName == nil || *p.InvoiceSectionDisplayName != *input.InvoiceSectionDisplayName) {
		return false
	}

	if p.InvoiceSectionId != nil && (input.InvoiceSectionId == nil || *p.InvoiceSectionId != *input.InvoiceSectionId) {
		return false
	}

	if p.InvoiceSectionSystemId != nil && (input.InvoiceSectionSystemId == nil || *p.InvoiceSectionSystemId != *input.InvoiceSectionSystemId) {
		return false
	}

	return true
}
