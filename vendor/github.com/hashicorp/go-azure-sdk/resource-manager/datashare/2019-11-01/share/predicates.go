package share

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProviderShareSubscriptionOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ProviderShareSubscriptionOperationPredicate) Matches(input ProviderShareSubscription) bool {

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

type ShareOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ShareOperationPredicate) Matches(input Share) bool {

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

type ShareSynchronizationListOperationPredicate struct {
	NextLink *string
}

func (p ShareSynchronizationListOperationPredicate) Matches(input ShareSynchronizationList) bool {

	if p.NextLink != nil && (input.NextLink == nil || *p.NextLink != *input.NextLink) {
		return false
	}

	return true
}

type SynchronizationDetailsListOperationPredicate struct {
	NextLink *string
}

func (p SynchronizationDetailsListOperationPredicate) Matches(input SynchronizationDetailsList) bool {

	if p.NextLink != nil && (input.NextLink == nil || *p.NextLink != *input.NextLink) {
		return false
	}

	return true
}
