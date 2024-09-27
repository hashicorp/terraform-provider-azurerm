package providers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProviderOperationPredicate struct {
	Id                 *string
	Namespace          *string
	RegistrationPolicy *string
	RegistrationState  *string
}

func (p ProviderOperationPredicate) Matches(input Provider) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Namespace != nil && (input.Namespace == nil || *p.Namespace != *input.Namespace) {
		return false
	}

	if p.RegistrationPolicy != nil && (input.RegistrationPolicy == nil || *p.RegistrationPolicy != *input.RegistrationPolicy) {
		return false
	}

	if p.RegistrationState != nil && (input.RegistrationState == nil || *p.RegistrationState != *input.RegistrationState) {
		return false
	}

	return true
}

type ProviderPermissionOperationPredicate struct {
	ApplicationId *string
}

func (p ProviderPermissionOperationPredicate) Matches(input ProviderPermission) bool {

	if p.ApplicationId != nil && (input.ApplicationId == nil || *p.ApplicationId != *input.ApplicationId) {
		return false
	}

	return true
}

type ProviderResourceTypeOperationPredicate struct {
	Capabilities      *string
	DefaultApiVersion *string
	ResourceType      *string
}

func (p ProviderResourceTypeOperationPredicate) Matches(input ProviderResourceType) bool {

	if p.Capabilities != nil && (input.Capabilities == nil || *p.Capabilities != *input.Capabilities) {
		return false
	}

	if p.DefaultApiVersion != nil && (input.DefaultApiVersion == nil || *p.DefaultApiVersion != *input.DefaultApiVersion) {
		return false
	}

	if p.ResourceType != nil && (input.ResourceType == nil || *p.ResourceType != *input.ResourceType) {
		return false
	}

	return true
}
