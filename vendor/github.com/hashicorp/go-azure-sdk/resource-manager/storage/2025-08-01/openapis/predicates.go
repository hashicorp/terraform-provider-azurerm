package openapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedAccountOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p DeletedAccountOperationPredicate) Matches(input DeletedAccount) bool {

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

type SkuInformationOperationPredicate struct {
	ResourceType *string
}

func (p SkuInformationOperationPredicate) Matches(input SkuInformation) bool {

	if p.ResourceType != nil && (input.ResourceType == nil || *p.ResourceType != *input.ResourceType) {
		return false
	}

	return true
}

type UsageOperationPredicate struct {
	CurrentValue *int64
	Limit        *int64
}

func (p UsageOperationPredicate) Matches(input Usage) bool {

	if p.CurrentValue != nil && (input.CurrentValue == nil || *p.CurrentValue != *input.CurrentValue) {
		return false
	}

	if p.Limit != nil && (input.Limit == nil || *p.Limit != *input.Limit) {
		return false
	}

	return true
}
