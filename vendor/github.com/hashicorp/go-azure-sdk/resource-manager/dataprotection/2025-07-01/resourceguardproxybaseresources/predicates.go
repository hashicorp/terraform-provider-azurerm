package resourceguardproxybaseresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceGuardProxyBaseResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ResourceGuardProxyBaseResourceOperationPredicate) Matches(input ResourceGuardProxyBaseResource) bool {

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
