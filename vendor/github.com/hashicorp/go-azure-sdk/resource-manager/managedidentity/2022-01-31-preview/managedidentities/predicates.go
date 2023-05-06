package managedidentities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureResourceOperationPredicate struct {
	Id                      *string
	Name                    *string
	ResourceGroup           *string
	SubscriptionDisplayName *string
	SubscriptionId          *string
	Type                    *string
}

func (p AzureResourceOperationPredicate) Matches(input AzureResource) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil && *p.Name != *input.Name) {
		return false
	}

	if p.ResourceGroup != nil && (input.ResourceGroup == nil && *p.ResourceGroup != *input.ResourceGroup) {
		return false
	}

	if p.SubscriptionDisplayName != nil && (input.SubscriptionDisplayName == nil && *p.SubscriptionDisplayName != *input.SubscriptionDisplayName) {
		return false
	}

	if p.SubscriptionId != nil && (input.SubscriptionId == nil && *p.SubscriptionId != *input.SubscriptionId) {
		return false
	}

	if p.Type != nil && (input.Type == nil && *p.Type != *input.Type) {
		return false
	}

	return true
}

type FederatedIdentityCredentialOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p FederatedIdentityCredentialOperationPredicate) Matches(input FederatedIdentityCredential) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil && *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil && *p.Type != *input.Type) {
		return false
	}

	return true
}

type IdentityOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p IdentityOperationPredicate) Matches(input Identity) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && *p.Location != input.Location {
		return false
	}

	if p.Name != nil && (input.Name == nil && *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil && *p.Type != *input.Type) {
		return false
	}

	return true
}
