package rulesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RuleSetOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RuleSetOperationPredicate) Matches(input RuleSet) bool {

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

type UsageOperationPredicate struct {
	CurrentValue *int64
	Id           *string
	Limit        *int64
}

func (p UsageOperationPredicate) Matches(input Usage) bool {

	if p.CurrentValue != nil && *p.CurrentValue != input.CurrentValue {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Limit != nil && *p.Limit != input.Limit {
		return false
	}

	return true
}
