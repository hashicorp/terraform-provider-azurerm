package trustedaccess

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrustedAccessRoleOperationPredicate struct {
	Name               *string
	SourceResourceType *string
}

func (p TrustedAccessRoleOperationPredicate) Matches(input TrustedAccessRole) bool {

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.SourceResourceType != nil && (input.SourceResourceType == nil || *p.SourceResourceType != *input.SourceResourceType) {
		return false
	}

	return true
}

type TrustedAccessRoleBindingOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p TrustedAccessRoleBindingOperationPredicate) Matches(input TrustedAccessRoleBinding) bool {

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
