package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocationOperationPredicate struct {
	DisplayName         *string
	Id                  *string
	Name                *string
	RegionalDisplayName *string
	SubscriptionId      *string
}

func (p LocationOperationPredicate) Matches(input Location) bool {

	if p.DisplayName != nil && (input.DisplayName == nil || *p.DisplayName != *input.DisplayName) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.RegionalDisplayName != nil && (input.RegionalDisplayName == nil || *p.RegionalDisplayName != *input.RegionalDisplayName) {
		return false
	}

	if p.SubscriptionId != nil && (input.SubscriptionId == nil || *p.SubscriptionId != *input.SubscriptionId) {
		return false
	}

	return true
}

type SubscriptionOperationPredicate struct {
	AuthorizationSource *string
	DisplayName         *string
	Id                  *string
	SubscriptionId      *string
	TenantId            *string
}

func (p SubscriptionOperationPredicate) Matches(input Subscription) bool {

	if p.AuthorizationSource != nil && (input.AuthorizationSource == nil || *p.AuthorizationSource != *input.AuthorizationSource) {
		return false
	}

	if p.DisplayName != nil && (input.DisplayName == nil || *p.DisplayName != *input.DisplayName) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.SubscriptionId != nil && (input.SubscriptionId == nil || *p.SubscriptionId != *input.SubscriptionId) {
		return false
	}

	if p.TenantId != nil && (input.TenantId == nil || *p.TenantId != *input.TenantId) {
		return false
	}

	return true
}
