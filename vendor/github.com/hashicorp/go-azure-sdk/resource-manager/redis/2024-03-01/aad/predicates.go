package aad

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisCacheAccessPolicyOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RedisCacheAccessPolicyOperationPredicate) Matches(input RedisCacheAccessPolicy) bool {

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

type RedisCacheAccessPolicyAssignmentOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RedisCacheAccessPolicyAssignmentOperationPredicate) Matches(input RedisCacheAccessPolicyAssignment) bool {

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
