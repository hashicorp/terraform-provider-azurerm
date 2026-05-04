package blobauditing

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseBlobAuditingPolicyOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p DatabaseBlobAuditingPolicyOperationPredicate) Matches(input DatabaseBlobAuditingPolicy) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Kind != nil && (input.Kind == nil || *p.Kind != *input.Kind) {
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

type ExtendedDatabaseBlobAuditingPolicyOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ExtendedDatabaseBlobAuditingPolicyOperationPredicate) Matches(input ExtendedDatabaseBlobAuditingPolicy) bool {

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

type ExtendedServerBlobAuditingPolicyOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ExtendedServerBlobAuditingPolicyOperationPredicate) Matches(input ExtendedServerBlobAuditingPolicy) bool {

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

type ServerBlobAuditingPolicyOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ServerBlobAuditingPolicyOperationPredicate) Matches(input ServerBlobAuditingPolicy) bool {

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
