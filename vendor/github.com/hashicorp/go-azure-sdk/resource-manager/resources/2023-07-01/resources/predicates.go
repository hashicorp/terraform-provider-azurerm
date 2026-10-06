package resources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GenericResourceExpandedOperationPredicate struct {
	ChangedTime       *string
	CreatedTime       *string
	Id                *string
	Kind              *string
	Location          *string
	ManagedBy         *string
	Name              *string
	Properties        *interface{}
	ProvisioningState *string
	Type              *string
}

func (p GenericResourceExpandedOperationPredicate) Matches(input GenericResourceExpanded) bool {

	if p.ChangedTime != nil && (input.ChangedTime == nil || *p.ChangedTime != *input.ChangedTime) {
		return false
	}

	if p.CreatedTime != nil && (input.CreatedTime == nil || *p.CreatedTime != *input.CreatedTime) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Kind != nil && (input.Kind == nil || *p.Kind != *input.Kind) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.ManagedBy != nil && (input.ManagedBy == nil || *p.ManagedBy != *input.ManagedBy) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Properties != nil && (input.Properties == nil || *p.Properties != *input.Properties) {
		return false
	}

	if p.ProvisioningState != nil && (input.ProvisioningState == nil || *p.ProvisioningState != *input.ProvisioningState) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}
